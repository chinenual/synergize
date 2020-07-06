# Hardware configuration

Synergize relies on a connection to the Synergy through a serial cable - which is unusual in this era of computing.  Many of us have not had to configure a serial cable and think about baud rates for decades (and some perhaps never have).   Here's a description of what needs to be done to connect Synergize to your Synergy.

## Synergy-side configuration

### Baud Rate
The Synergy's serial port baud rate is configured with a hardware jumper on the serial interface daughter board. The settings must match the baud rate setting in Synergize.

<img title="Serial Jumpers" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/serial-jumpers.png?raw=true" width="85%"/>

### Data bits

The Synergy always communicates with 8 bits.

### Parity

By default, the Synergy is configured for "no parity" and that is how I have tested the serial connections.  However, this is configurable via buttons on the front panel of the Synergy:
Press RESTORE twice and the LED's in the Channel Assign button display the parity settings.  The Channel Assign button can be used to alter the parity setting while in this mode.

| Channel Assign | Parity |
|+----- |  |
| Both On | No Parity |
| Left On | Odd Parity |
| Right On | Even Parity |
| Neither On | Error detected on serial board |

The Synergy always uses one start and stop bit.

### Flow Control

The Synergy uses DTR and CTS ("hardware flow control").  

## My configuration

My development and testing setup is a mix of direct connections and virtualized environments, so they may not be directly useful to you, but I'm listing them here to give you concrete examples things that work for me.

### Synergy

My Synergy is configured at 9600 baud and no parity as shown above.  I've tested successfully with 19200 but my day to day use is at 9600 baud.

### Serial Cable(s)

The connection from the computer to the Synergy requires a
"[null modem](https://en.wikipedia.org/wiki/Null_modem)" cable. I test with a combination of a traditional null modem cable (same type you would use for a direct SYNHCS/Kaypro connection) and an FTDI based USB serial cable.

I've been asked specifically what cables I use. I've been told that FTDI makes the most robust chipset, but I have _NO_ objective reason to prefer one cable over another.  These are not "endorsements" - it's just a note of what "works for me":

* FTDI based USB cable: [Sabrent FTDI USB to Serial](https://www.amazon.com/gp/product/B006AA04K0)
* Null modem serial cable: [C2G 02019 DB25 to DB9 Null Modem cable](https://www.amazon.com/gp/product/B000083K2R/)

### MacOS (Mojave)

This is my primary development environment; it's tested natively:

| Environment  		| Device 						| Baud Rate 		| OS Serial Config |
|:----- |:-----:|:-----|:-----
| MacOS (Mojave) 	| `/dev/tty.usbserial-AL05OC8S` 	| `9600`				| n/a

### Windows 10

This is tested via a virtual machine running in Parallels on the Mac. Here there are three sets of settings - the Parallels virtual machine, and the Windows OS and Synergize itself.  Note that I run Synergize at 9600 to match my Synergy, but the virtualized Windows OS is configured at 19200.  Also note that I'm using the cu variant of the unix serial port rather that the tty variant.


| Environment  		| Device 						| Synergize Baud Rate 	| OS Serial Config |
|:----- |:-----:|:-----|:-----
| Parallels VM		| Serial Device: `/dev/cu.usbserial-AL05OC8S` 	
| Windows10			| `COM1`							| 9600 					| Device Manager: `19200 baud, 8 bits, No parity , 1 stop bit, hardware flow control`


### Linux 64bit

I test the 64bit Linux version via Parallels. Ensure your user is a member of the `dialout` group.

| Environment  		| Device 						| Synergize Baud Rate 	| OS Serial Config |
|:----- |:-----:|:-----|:-----
| Parallels VM		| Serial Device: `/dev/cu.usbserial-AL05OC8S` 	
| Ubuntu 18.04			| `/dev/ttyS1`				| `9600` 					| n/a

### Linux 32bit

My 32bit linux environment is raw metal (an old laptop with a real serial port). So no USB serial cable involved - just the null modem cable plugged into the serial port.  Ensure your user is a member of the `dialout` group.

| Environment  		| Device 						| Synergize Baud Rate 	| OS Serial Config |
|:----- |:-----:|:-----|:-----
| Ubuntu 16.04			| `/dev/ttyS1`				| `9600` 					| n/a

## Troubleshooting

* If you have problems connecting, be sure to match the baud rate on your Synergy. Checking this requires opening up and looking at a jumper on the interface daughter board.  Mine was originally set to 9600 (I'm not sure if this was the "factory default" or if it was tweaked by a previous owner). As long as Synergize is configured with the same rate, all is good. I've tested mine at both 9600 and 19,200 baud and things work fine.

* On Linux, be sure your user is a member of the `dialout` group (permissions on the serial port device usually limit access to members of that group).

* One user reports that he had problems (no serial communication at all), but then swapped his USB serial cable for another cable and has had success since.  The original cable was not obviously bad (it worked with some other software he uses), but there was an as-yet-undiagnosed problem when using it with Synergize.  So swapping out cables might be a last resort if you can't make things work otherwise.

* You can enable very verbose logging to help diagnose serial port configuration issues.  Try the following from a command line:
```
/path/to/Synergize -port /dev/YourSerialPort -baud YourBaud -SERIALVERBOSE -SYNVER
```
e.g.,
```
/path/to/Synergize -port /dev/ttyS1 -baud 19200 -SERIALVERBOSE -SYNVER
```
this will print some detailed log statements to the terminal which may shed light on the problem.
*NOTE:* On Windows, a separate executable named `Synergize-cmd.exe` must be used -- the main `Synergize.exe` does not support command line diagnostics.
* The command line also offers a way to run the Synergy's "LOOP TEST" without crafting a special cable that reverses the transmit and recieve lines (essentially Synergize will emulate the special cable in software).  Using an unaltered set of serial cables (i.e. a null modem cable and a USB serial cable if you use one), run
```
/path/to/Synergize -port /dev/YourSerialPort -baud YourBaud -SERIALVERBOSE -LOOPTST
```
You will be prompted to put the Synergy into the LOOP TEST mode by pressing "`RESTORE RESTORE Program #1`" on the Synergy's front panel.  If the test is successful, the Synergy will return to its power up state.  If the test fails, the programmer section of the Synergy's panel will flash continuously.
This test will help validate that the Synergy's serial device, your cable(s) and your parity/baud/device configuration are correct.
