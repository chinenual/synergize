---
layout: default
title: Troubleshooting
description: Serial Port / Hardware Troubleshooting
parent: Serial Port / Hardware Setup
nav_order: 2
permalink: /docs/hardware/troubleshooting
---

## Troubleshooting

* If you have problems connecting, be sure to match the baud rate on your Synergy. Checking this requires opening up and looking at a jumper on the interface daughter board.  Mine was originally set to 9600 (I'm not sure if this was the "factory default" or if it was tweaked by a previous owner). As long as Synergize is configured with the same rate, all is good. I've tested mine at both 9600 and 19,200 baud and things work fine.

* On Linux, be sure your user is a member of the `dialout` group (permissions on the serial port device usually limit access to members of that group).

* One user reports that he had problems (no serial communication at all), but then swapped his USB serial cable for another cable and has had success since.  The original cable was not obviously bad (it worked with some other software he uses), but there was an as-yet-undiagnosed problem when using it with Synergize.  So swapping out cables might be a last resort if you can't make things work otherwise.

### Verbose Logging

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

### Serial Diagnostics (COMTST)

Synergize can test the connection to the Synergy in cooperation with the COMTST onboard test mode of the Synergy. In this test, Synergize sends bytes from 0 through 255 and the Synergy is expected to echo them back.

To run the test, select `Connect->Serial Diagnostics` and follow the instructions on the screen. NOTE: once the Synergy is in this test state, it will remain there until it is power cycled.

This can also be done from the command line:
```
/path/to/Synergize -port /dev/YourSerialPort -baud YourBaud -SERIALVERBOSE -COMTST
```
The bottom row of voice lights is the byte value last received by
the  Synergy and echoed back to the host.


### Serial Diagnostics (LOOPTST)

The command line offers a way to run the Synergy's "LOOP TEST" without crafting a special cable that reverses the transmit and receive lines (essentially Synergize will emulate the special cable in software).

This is essentially the inverse of the COMTST - the Synergy sends bytes and expects to be be able to read the same byte it sent.

Using an unaltered set of serial cables (i.e. a normal null modem cable and a USB serial cable if you use one), run
```
/path/to/Synergize -port /dev/YourSerialPort -baud YourBaud -SERIALVERBOSE -LOOPTST
```
You will be prompted to put the Synergy into the LOOP TEST mode by pressing `RESTORE RESTORE Program #1` on the Synergy's front panel.  If the test is successful, the Synergy will return to its power up state.  If the test fails, the programmer section of the Synergy's panel will flash continuously.
This test will help validate that the Synergy's serial device, your cable(s) and your parity/baud/device configuration are correct.

### Serial Diagnostics (LINKTST)

Finally, Synergize can act as a terminal client for the Synergy's LINKTST. This allows you to type to the terminal and send and recieve bytes to the Synergy's `RESTORE Program #4` serial link test.  The state of the Synergy's serial interface is displayed via voice LED's on the Synergy:
```
/path/to/Synergize -port /dev/YourSerialPort -baud YourBaud -SERIALVERBOSE -LINKTST
```

| LED | Status | Description |
|+--- | +------ |+----   |
| Voice 5  | Clear to Send from Host (CTS) | The host is ready to receive data from Synergy|
| Voice 6  | Break Detect | (not used by Synergy) |
| Voice 7  | Framing error | Last character received did not have a valid stop bit. This used for diagnostics only.|
| Voice 8  | Overrun Error | The host started transmitting before the Synergy read the last byte sent.|
| Voice 9  | Parity Error | Indicates Synergy and host do not have equal parity setups or data is getting scrambled.|
| Voice 10 | Transmitter Enabled | The Synergy enables the transmitter at startup. If this lamp is not on, the interface board is defective or missing and therefore the Synergy|
| Voice 11 | Receive Data Available | This lamp is on when the interface has new data and the Synergy has not read it yet.|
| Voice 12 | Transmitter Buffer Empty | This lamp is on when the Synergy is not sending.|

The bottom row of voice lights is the byte value last received by
the  Synergy and echoed back to the host.
 
