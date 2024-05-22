# Synergize

When built for development, use cntl-D to toggle Chrome Developer Tools.

# SYNHCS

## Setup

* compile interceptty 
* install DOSBox
* install Z80 emulator Z80MU
* setup DOSBox config. add: 
   ```
   serial1=directserial realport:ttys002
   mount C "/Users/tynor/Documents/DOSBox/SYNHCS"
   mount A "/Users/tynor/Documents/DOSBox/SYNHCS/SYNHCS/CARTS"
   mount B "/Users/tynor/Documents/DOSBox/SYNHCS/SYNHCS/VOICES/Internal"
   C:
   ```

## run

* run interceptty and take note of device it allocates.  after it starts ls -l /tmp/tty.proxy to find the pty to use for dosbox.  Edit the DOSBox preferences accordingly
   
   ```
   ~/src/3rdParty/interceptty/interceptty -v /dev/cu.usbserial-AL05OC8S /tmp/tty.proxy
cc```

* Start DOSBox
* Start Z80MU
* Start SYNHCS
   
   
