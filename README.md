[![License Apache 2.0](https://img.shields.io/badge/license-Apache%20License%202.0-green.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![PayPal donation](https://img.shields.io/badge/donation-PayPal-red.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=AHWJHJFBAWGL2)
[![Yandex.Money donation](https://img.shields.io/badge/donation-Я.деньги-yellow.svg)](https://money.yandex.ru/embed/small.xml?account=41001158080699&quickpay=small&yamoney-payment-type=on&button-text=01&button-size=l&button-color=orange&targets=%D0%9F%D0%BE%D0%B6%D0%B5%D1%80%D1%82%D0%B2%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D0%B5+%D0%BD%D0%B0+%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D1%8B+%D1%81+%D0%BE%D1%82%D0%BA%D1%80%D1%8B%D1%82%D1%8B%D0%BC+%D0%B8%D1%81%D1%85%D0%BE%D0%B4%D0%BD%D1%8B%D0%BC+%D0%BA%D0%BE%D0%B4%D0%BE%D0%BC&default-sum=100&successURL=)

# What is it?
It is a small utility to convert .BIN snapshots ([BK-0010(01)](http://en.wikipedia.org/wiki/Electronika_BK) emulator format) into sound WAV files which can be played and recognized by real BK-0010 TAP reader.   
The Project is based on [old QBasic based converter project](http://bk-mg.narod.ru/).

# What is BK-0010
[BK-0010](http://en.wikipedia.org/wiki/Electronika_BK) is the most popular soviet 16 bit home computer platform with some PDP-11 compatibility.

# How to load?
 - [The Python script version](https://raw.githubusercontent.com/raydac/bkbin2wav/master/bkbin2wav.py)
 - [The Windows version](https://github.com/raydac/bkbin2wav/releases/download/1.0.1/bkbin2wav-windows386.exe)
 - [The Mac OS version](https://github.com/raydac/bkbin2wav/releases/download/1.0.1/bkbin2wav-macos386.bin)
 - [The Linux version](https://github.com/raydac/bkbin2wav/releases/download/1.0.1/bkbin2wav-linux386.bin)

# Known archives with snapshots
- [Big archive of game snapshot](http://roman-dushkin.narod.ru/bk_games_all.html)
- [Yet another big archive of game snapshots](http://www.bk001x.ru/index/na_bukvu_quot_a_quot/0-184)

# How to use it?
The Script is written in [Python](https://www.python.org/downloads/) so that it is more or less crossplatform one. You should have installed [Python](https://www.python.org/downloads/) on your machine. The Utility has only command line interface, so that format of call is:
```
python bkbin2wav.py -i <input file> [-o <output file>]
```
example for python version:
```
python bkbin2wav.py -i Arkanoid.bin -o Arkanoid.wav
```
example for native version
```
bkbin2wav-windows386.exe -i Arkanoid.bin -o Arkanoid.wav
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
Sometime .BIN files may contain wrong data size value defined in their header, in the case you can use **-f** flag to enforce usage of physical file length instead of the data length defined in the BIN header.
