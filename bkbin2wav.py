#!/usr/bin/python
# The Utility allows to convert BK-0010 .BIN snapshots into WAV format to be loaded through cable to real device
#
# URL: https://github.com/raydac/bkbin2wav
# Author: Igor Maznitsa (http://www.igormaznitsa.com)
# Version: 1.01
#
# License: Apache License 2.0 (https://www.apache.org/licenses/LICENSE-2.0)

__author__ = 'Igor Maznitsa (http://www.igormaznitsa.com)'
__version__ = '1.0.1'
__projecturl__ = 'https://github.com/raydac/bkbin2wav'

import wave
import getopt
import sys
import os
import struct

SND_PARTS = [
    bytearray([0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40]),
    bytearray(
        [0x80, 0xa0, 0xb7, 0xc0, 0xb7, 0xa0, 0x80, 0x5f, 0x48, 0x3f, 0x48, 0x5f, 0x80, 0xb7, 0xb7, 0x80, 0x48, 0x48]),
    bytearray(
        [0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40,
         0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40,
         0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40, 0x80, 0xbf, 0xbf, 0x80, 0x40, 0x40]),
    bytearray(
        [0x80, 0x90, 0x9d, 0xa4, 0xa6, 0xa9, 0xa9, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0x80, 0x6f, 0x62, 0x5b, 0x57, 0x56,
         0x55, 0x55, 0x55, 0x55, 0x55, 0x6e, 0x80, 0x9a, 0xa2, 0xa6, 0xa7, 0xa9, 0x80, 0x6f, 0x63, 0x5c, 0x59, 0x59,
         0x80, 0xb7, 0xb7, 0x80, 0x48, 0x48]),
    bytearray(
        [0x80, 0x8f, 0xa8, 0xb5, 0xbc, 0xbf, 0xc1, 0xc1, 0xc2, 0xc2, 0xc2, 0xc2, 0xc2, 0xc1, 0xc1, 0xc1, 0xc1, 0xc1,
         0xc0, 0xc0, 0xc0, 0xc0, 0xc0, 0xbf, 0xbf, 0xbe, 0xbe, 0xbe, 0xbd, 0xbc, 0xb2, 0x80, 0x59, 0x48, 0x40, 0x3b,
         0x39, 0x38, 0x37, 0x37, 0x37, 0x37, 0x37, 0x38, 0x38, 0x38, 0x38, 0x39, 0x39, 0x39, 0x39, 0x39, 0x3a, 0x3a,
         0x3a, 0x3a, 0x3a, 0x3b, 0x3b, 0x3b, 0x3b, 0x3c, 0x3e, 0x7a])
]

SIGNAL_RESET = 0
SIGNAL_SET = 1
SIGNAL_START_SYNCHRO = 2
SIGNAL_SYNCHRO = 3
SIGNAL_END_MARKER = 4


def amplify():
    """Extends audio diapason in the WAV signals to 1.0 instead of 0.5"""
    min_l = 256
    max_l = 0
    for arr in SND_PARTS:
        for b in arr:
            if b < min_l:
                min_l = b
            if b > max_l:
                max_l = b

    max_l -= 128
    min_l -= 128

    c_max = 128.0 / max_l
    c_min = -127.0 / min_l

    coeff = min(c_max, c_min)

    for arr in SND_PARTS:
        for i in range(len(arr)):
            arr[i] = max(0, min(255, int(round((arr[i] - 128) * coeff)) + 128))


def play(i, arr, num=1):
    """Save wav sequence for num-times into the result"""
    for c in range(num):
        for x in SND_PARTS[i]:
            arr.append(x)


def byte2arr(b, arr):
    """Append wav representation of a sound to array"""
    for i in range(0, 8):
        for a in SND_PARTS[(b >> i) & 1]:
            arr.append(a)


def txt2arr(txt, arr):
    """Write text chars as bytes into sound array (with 16 char restrictions)"""
    ws = 0
    if len(txt) < 16:
        ws = 16 - len(txt)
    if len(txt) > 16:
        txt = txt[0:16]
    for c in txt:
        byte2arr(ord(c), arr)

    while ws != 0:
        byte2arr(ord(' '), arr)
        ws -= 1


def short2arr(num, arr):
    """Write a short (0...FFFF) value into sound array)"""
    byte2arr(num & 0xFF, arr)
    byte2arr((num >> 8) & 0xFF, arr)


def init_wav(turbo_mode, filename, data_length):
    r = wave.open(filename, 'w')
    r.setnchannels(1)
    r.setsampwidth(1)
    if turbo_mode:
        r.setframerate(22050)
    else:
        r.setframerate(11025)
    r.setnframes(data_length)
    r.setcomptype('NONE', 'PCM')
    return r


def calc_crc(data):
    a = 0
    for i in data:
        a += i
        if a > 0xFFFF:
            a -= 0xFFFF
    return a


def make_bk_wav(turbo_mode, filename, start, name, data, data_length):
    sound_buffer = []

    play(SIGNAL_START_SYNCHRO, sound_buffer, 512)
    play(SIGNAL_SYNCHRO, sound_buffer)

    play(SIGNAL_START_SYNCHRO, sound_buffer)
    play(SIGNAL_SYNCHRO, sound_buffer)

    short2arr(start, sound_buffer)
    short2arr(data_length, sound_buffer)
    txt2arr(name, sound_buffer)

    play(SIGNAL_START_SYNCHRO, sound_buffer)
    play(SIGNAL_SYNCHRO, sound_buffer)

    for i in range(data_length):
        byte2arr(data[i], sound_buffer)

    tap_check_code = calc_crc(data)
    short2arr(tap_check_code, sound_buffer)

    play(SIGNAL_END_MARKER, sound_buffer)
    play(SIGNAL_START_SYNCHRO, sound_buffer, 64 if turbo_mode else 32)
    play(SIGNAL_SYNCHRO, sound_buffer)

    wav_file = init_wav(turbo_mode, filename, len(sound_buffer))
    try:
        wav_file.writeframes(bytearray(sound_buffer))
    finally:
        wav_file.close()
    return tap_check_code


def read_short(f):
    v = struct.unpack('BB', f.read(2))
    return (v[1] << 8) | v[0]


def read_bin(bin_file, enforce_physical_length):
    physical_length = os.path.getsize(bin_file) - 4

    if physical_length < 0:
        print("Wrong BIN snapshot format, its size less than 4")
        exit(1)

    f = open(bin_file, 'rb')
    with open(bin_file, 'rb') as f:
        start_address = read_short(f)
        bin_data_length = read_short(f)
        data_length = bin_data_length
        if enforce_physical_length:
            data_length = physical_length

        if data_length > physical_length:
            print("Wrong data size length in BIN snapshot %d" % (data_length))
            exit(1)
        arr = bytearray(f.read(data_length))

        return start_address, data_length, bin_data_length, arr,


def header():
    print("""
    BKBIN2WAV allows to convert .BIN snapshots (for BK-0010(01) emulators) into WAV format.

    Project page : %s
          Author : %s
         Version : %s

    It is a converter of .BIN files (a snapshot format for BK-0010(01) emulators) into sound WAV files which compatible with real BK-0010 TAP reader
    """ % (__projecturl__, __author__, __version__))


def help():
    print("""    Usage:

    bkbin2wav -i <binfile> [-a] [-o <wavfile>] [-n <name>] [-s addr] [-t]

    Command line options:
        -h          Print help
        -f          Use file size instead of .BIN header size field value
        -a          Amplify the audio signal in the result WAV file
        -i <file>   The BIN file to be converted
        -o <file>   The Result WAV file (by default the BIN file name with WAV extension)
        -n <name>   The Name of the file in the TAP header (must be less or equals 16 chars)
        -s <addr>   The Start address for the TAP header (by default the start address from the BIN will be used)
        -t          Use the double frequency "turbo" mode
    """)


header()

opts, args = getopt.getopt(sys.argv[1:], "atfhi:o:n:s:",
                           ['amplify', 'turbo', 'help', 'input', 'output', 'name', 'start'])

infile = ''
outfile = ''
bin_start_address = -1
name = ''
turbo = False
amplifier_flag = False
enforce_physical_size = False

for o, a in opts:
    if o in ('-h', '--help', '-?', '--h', '-help'):
        help()
        exit(1)
    elif o == '-i':
        infile = a
    elif o == '-s':
        bin_start_address = int(a) & 0xFFFF
    elif o == '-a':
        amplifier_flag = True
        amplify()
    elif o == '-o':
        outfile = a
    elif o == '-n':
        name = a
    elif o == '-t':
        turbo = True
    elif o == '-f':
        enforce_physical_size = True
    else:
        help()
        exit(1)

if len(opts) == 0:
    help()
    exit(1)

if infile == '':
    print("The Input BIN file is undefined")
    exit(1)

if not os.path.isfile(infile):
    print("Can't find the BIN file %s" % (infile))
    exit(1)

if outfile == '':
    outfile = os.path.abspath(os.path.dirname(infile)) + os.sep + os.path.basename(infile) + '.wav'

if name == '':
    name = os.path.splitext(os.path.basename(infile))[0].upper()[:16]

bin_file_size_without_header = os.path.getsize(infile) - 4;
bin_start, bin_length, bin_hdr_len, bin_read_length = read_bin(infile, enforce_physical_size)

if bin_start_address == -1:
    bin_start_address = bin_start

if enforce_physical_size:
    print(
        "Detected flag to enforce physical file size (size defined inside of .BIN is %s byte(s), real size is %s byte(s))" % (
        bin_hdr_len, bin_file_size_without_header))
else:
    if bin_hdr_len != bin_file_size_without_header:
        print(
            "Warning! Detected different size defined in BIN header, use -f to use file size instead of header size (%s != %s)" % (
            bin_hdr_len, bin_file_size_without_header))

if bin_start_address != bin_start:
    print("Warning! The Start address has been changed from %d(%s) to %d(%s)" % (
    bin_start, oct(bin_start), bin_start_address, oct(bin_start_address)))

print("""
   Input file : %s
  Output file : %s
         Name : %s
Start address : %d (%s)
   Turbo mode : %s
    Amplifier : %s
""" % (infile, outfile, name, bin_start_address, oct(bin_start_address), ('ON' if turbo else 'OFF'),
       ('ON' if amplifier_flag else 'OFF')))

print("""    Fields of the BIN file header:
       Start  : %d (%s)
       Length : %d (%s)
""" % (bin_start, oct(bin_start), bin_length, oct(bin_length)))

crc_code = make_bk_wav(turbo, outfile, bin_start_address, name, bin_read_length, bin_length)

print("The WAV file has been saved successfully as '%s', the checksum is %s" % (outfile, hex(crc_code)))