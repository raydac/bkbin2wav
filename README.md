# What is it?
It is a small utility to convert .BIN snapshots ([BK-0010(01)](http://en.wikipedia.org/wiki/Electronika_BK) emulator format) into sound WAV files which can be played and recognized by real BK-0010 TAP reader.   
The Project is based on [old QBasic based converter project](http://bk-mg.narod.ru/).

# What is BK-0010
[BK-0010](http://en.wikipedia.org/wiki/Electronika_BK) is the most popular soviet 16 bit home computer platform with some PDP-11 compatibility.

# How to use it?
The Script is written in [Python](https://www.python.org/downloads/) so that it is more or less crossplatform one. You should have installed [Python](https://www.python.org/downloads/) on your machine. The Utility has only command line interface, so that format of call is:
```
python bkbin2wav.py -i <input file> [-o <output file>]
```
example:
```
python bkbin2wav.py -i Arkanoid.bin -o Arkanoid.wav
```
if to start the script without parameters, then it will print allowed CLI flags
```
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
```
Sometime .BIN files contain wrong size value in their header, in the case you can use **-f** option to enforce usage of physical file length instead of the header length value.

# Donation   
If you like the software you can make some donation to the author   
[![https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=AHWJHJFBAWGL2)
