# synergize

<img src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/logo-for-github.png?raw=true"/>

A portable voice librarian for the DKI Synergy synthesizer

## Features

Native support for both MacOS and Windows. No need for a Kaypro. No
need to emulate CP/M software.

* Load virtual "carts" (CRT files).   No need for a physical
library of carts.

* Load individual "voices" (VCE files).   View their properties in a
manner similar to the original SYNHCS ".DOC" format and some of the edit screens in SYNHCS.

* Load and save synergy "state" (SYN files).  This preserves
sequencer, portamento, vibrato and other performance customizations.

Note: Synergize does not itself include the original voice libraries -- download a copy from one of the links below.

## Download

Download the release:
[Release Binaries](https://github.com/chinenual/synergize/releases/latest)

See links below for sites containing the original DK/Mulogix voice library.

## Issues / Bugs / Feature Requests

Please report problems via a [github issue](https://github.com/chinenual/synergize/issues).  If you don't have a github account, send a bug report email to "support at chinenual.com"

If you are a user or are thinking of becoming a user, please join our Slack workspace [https://synergize-app.slack.com](https://synergize-app.slack.com).  Send an email to "support at chinenual.com" for an invitation to the group.


## Screenshots

<img title="Cartridge voice listing view" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/viewCRT.png?raw=true" width="45%"/>
<img title="Voice main patch view" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/viewVCE_voice.png?raw=true" width="45%"/>
<img title="Voice Envelopes view" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/viewVCE_envs.png?raw=true" width="45%"/>
<img title="Voice Filters view" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/viewVCE_filters.png?raw=true" width="45%"/>
<img title="Voice Key Equalization view" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/viewVCE_keyeq.png?raw=true" width="45%"/>
<img title="Voice Key Proportion view" src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/viewVCE_keyprop.png?raw=true" width="45%"/>

## Requirements

### Synergy II+

Synergize connects to a Synergy II+ via its RS232 serial port.  It is tested with a Synergy running the latest 3.22 firmware - but is probably compatible with earlier versions.

### Serial Cable(s)

The connection from the computer to the Synergy requires a "null modem" cable. I test with a combination of a traditional null modem cable (same one you would use for a direct SYNHCS/Kaypro connection) and an FTDI based USB serial cable.

I've been asked specifically what cables I use. I've been told that FTDI makes the most robust chipset, but I have _NO_ objective reason to prefer one cable over another.  These are not "endorsements" - it's just a note of what "works for me":

* FTDI based USB cable: [Sabrent FTDI USB to Serial](https://www.amazon.com/gp/product/B006AA04K0)
* Null modem serial cable: [C2G 02019 DB25 to DB9 Null Modem cable](https://www.amazon.com/gp/product/B000083K2R/)

### Operating Systems

Synergize has been tested on:

* MacOS 10.14.6 (Mojave)
* Windows 10
* Ubuntu Linux 18.04
* Ubuntu Linux 16.04

## How to use

### Connecting to the Synergy

Set your Serial Port device name and baud rate via the Help->Preferences menu (via the cheesy Help button at the top of the page).  You can then test the connection by selecting the Connect->Connect to Synergy menu.  If successful, Synergize will report the firmware version of the connected Synergy (and display it in the upper left pane of the display).

It is not necessary to explicitly connect in this way -- the first time you invoke a command that needs to communicate with the Synergy, Synergize will initialize the connection implicitly.

### The Library Browser

The left hand pane of the UI is a file browser, allowing you to navigate your Synergy voice library and select SYN, CRT or VCE files to load.  The default location of the library is set via Help->Preferences.

### CRT Viewer

When you load a Cartridge (CRT) file, basic information is displayed (the voice assignments).  You can click on any voice to drill down and see voice documentation.  Or you can click the red  Load CRT button to upload the CRT to the Synergy.

### VCE Viewer

When you load a Voice (VCE) file or drill down to a voice from the CRT viewer, various information about the voice is displayed.  This includes the basic info included in the original "DOC" files, plus screens that replicate various screens from the original SYHNCS software - including frequency and amplitude envelopes, filters, keyboard equalization and proportion curves.

### Returning the Synergy to normal state

Once a CRT is loaded, the Synergy is using its "VRAM" data (its Cartridge button is flashing).   To use the internal voices or a physical cartridge, select Connect->Disable VRAM.

### Diagnostics

Synergize can test the connection to the Synergy in cooperation with the COMTST onboard test mode of the Synergy.  To run that, select Connect->Serial Diagnostics and follow the instructions on the screen.
NOTE: once the Synergy is in this test state, it will remain there until it is power cycled.

## Known Issues

* If you have problems connecting, be sure to match the baud rate on your Synergy. Checking this requires opening up and looking at a jumper on the interface daughter board.  Mine was originally set to 9600 (not sure if this was the "factory default" or if it was tweaked by a previous owner). As long as Synergize is configured with the same rate, all is good. I've tested mine at both 9600 and 19,200 baud and things work fine.

* On Linux, be sure your user is a member of the `dialout` group (permissions on the serial port device usually limit access to members of that group).

* On MacOS, if you change serial parameters, you will need to restart the application in order to "reconnect" to the synergy with those parameters. On Windows, you can directly re-connect via the Connect->Connect to Synergy menu.

* Due to the mysteries of serial port communication, attempts to save or load files to the Synergy will sometimes fail (often reporting a TIMEOUT).  If this happens, a second try will usually succeed. Fixes introduced in 0.2.0 and 1.0.0 greatly improve, but don't completely fix this.

* See the
[Release Notes](https://github.com/chinenual/synergize/releases/latest)
for some caveats regarding the Linux builds.

## Voice Library

Full sets of the DK and Mulogix library CRT and VCE files are also
available via the below links.  The set I'm using includes the Internal voices, the 6 standard Carts and additional "Library" voices -- it is available at:

* [Synergy Voice Library at groups.io](https://groups.io/g/synergy-synth/files/SynergyVoiceLibrary.zip)

## Thank you!

This would not have be possible without access to the excellent
documentation, and well commented firmware and SYNHCS Z80 source code donated to thecommunity by  Stoney Stockell and Mulogix, Inc.  Those are available in several locations:

* [Synergy Facebook Users Group](https://www.facebook.com/groups/synergysynth/)
* [Synergy groups.io group](https://groups.io/g/synergy-synth)
* Aaron Lanterman's [Synergy Preservation Page](https://lanterman.ece.gatech.edu/synergy/)



## TODO

Currently, Synergize can load VCE and CRT files from a preexisting
library and upload them to the Synergy. It can load and save SYN files
to and from the Synergy.  It cannot (yet) edit CRT's or VCE's.
