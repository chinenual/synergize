# Hardware configuration

Synergize relies on a connection to the Synergy through a serial cable - which is unusual in this era of computing.  Many of us have not had to configure a serial cable and think about baud rates for decades (and some perhaps never have).   Here's a description of what needs to be done to connect Synergize to your Synergy.

## Synergy-side configuration

The Synergy's serial port is configured with hardware jumpers on the serial interface daughter board. The settings must match the baud rate setting in Synergize (and in the case of Windows, the flow control settings in Device Manager).

<img title="Serial Jumpers" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/serial-jumpers.png?raw=true" width="85%"/>


## My configuration

My development and testing setup is a mix of direct connections and virtualized environments, so they may not be directly useful to you, but I'm listing them here to give you concrete examples things that work for me.

### MacOS (Mojave)

This is my primary development environment; it's tested natively:

| Environment  		| Device 						| Baud Rate 		| OS Serial Config |
|:----- |:-----:|:-----|:-----
| MacOS (Mojave) 	| /dev/tty.usbserial-AL05OC8S 	| 9600				| n/a

### Windows 10

This is tested via a virtual machine running in Parallels on the mac. Here there are three sets of settings - the Parallels virtual machine, and the windows OS.  Note that I run Synergize at 9600 to match my Synergy, but the virtualized WIndows OS is configured at 19200.  Also note that I'm using the cu variant of the unix serial port rather that the tty variant.


| Environment  		| Device 						| Synergize Baud Rate 	| OS Serial Config |
|:----- |:-----:|:-----|:-----
| Parallels VM		| Serial Device: /dev/cu.usbserial-AL05OC8S 	
| Windows10			| COM1							| 9600 					| Device Manager: 19200, 8, N , 1 


### Linux 64bit

I test the 64bit Linux version via Parallels. Ensure your user is a member of the ``dialout` group.

| Environment  		| Device 						| Synergize Baud Rate 	| OS Serial Config |
|:----- |:-----:|:-----|:-----
| Parallels VM		| Serial Device: /dev/cu.usbserial-AL05OC8S 	
| Ubuntu 18.04			| /dev/ttyS1				| 9600 					| n/a

### Linux 32bit

My 32bit linux environment is raw metal (an old laptop with a real serial port). So no USB serial cable involved - just the null modem cable plugged into the serial port.  Ensure your user is a member of the ``dialout` group.

| Environment  		| Device 						| Synergize Baud Rate 	| OS Serial Config |
|:----- |:-----:|:-----|:-----
| Ubuntu 16.04			| /dev/ttyS1				| 9600 					| n/a

### Serial Cable(s)

The connection from the computer to the Synergy requires a "null modem" cable. I test with a combination of a traditional null modem cable (same one you would use for a direct SYNHCS/Kaypro connection) and an FTDI based USB serial cable.

I've been asked specifically what cables I use. I've been told that FTDI makes the most robust chipset, but I have _NO_ objective reason to prefer one cable over another.  These are not "endorsements" - it's just a note of what "works for me":

* FTDI based USB cable: [Sabrent FTDI USB to Serial](https://www.amazon.com/gp/product/B006AA04K0)
* Null modem serial cable: [C2G 02019 DB25 to DB9 Null Modem cable](https://www.amazon.com/gp/product/B000083K2R/)
