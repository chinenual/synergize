---
nav_exclude: true
layout: landing
permalink: /
---
<!--![Go](https://github.com/chinenual/synergize/workflows/Go/badge.svg)-->

<!--# synergize-->

<img src="https://github.com/chinenual/synergize/raw/master/docs/screenshots/logo-for-github.png?raw=true"/>

A portable voice editor and librarian for the DKI Synergy synthesizer

## Features

* Native support for MacOS, Windows and Linux. No need for a Kaypro. No
need to emulate CP/M software.

* Load and create virtual "carts" (CRT files).   No need for a physical
library of carts.

* Load individual "voices" (VCE files).   View their properties in a
manner similar to the original SYNHCS ".DOC" format and the edit screens in SYNHCS.

* Edit all voice parameters (patch routing, envelopes, filters, etc.) and create new VCE files.
Edit existing voices or create your own from scratch.

* Load and save synergy "state" (SYN files).  This preserves
sequencer, portamento, vibrato and other performance customizations.

Note: Synergize does not itself include the original voice libraries -- download a copy from one of the links below.

## Download

Download the release:
[Release Binaries](https://github.com/chinenual/synergize/releases)

See the [User Manual](voice-library.md) sites where you can download the original DK/Mulogix voice library.

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

## Requirements and How to use

See the [User Manual](intro.md) for full details on hardware, OS
requirements and application usage.

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

See [HARDWARE](hardware.md) for some serial port related troubleshooting guidelines.

* On MacOS, if you change serial parameters, you will need to restart the application in order to "reconnect" to the synergy with those parameters. On Windows, you can directly re-connect via the Connect->Connect to Synergy menu.

* Due to the mysteries of serial port communication, attempts to save or load files to the Synergy will sometimes fail (often reporting a TIMEOUT).  If this happens, a second try will usually succeed. Fixes 

* See the
[Release Notes](https://github.com/chinenual/synergize/releases)
for some caveats regarding the Linux builds.


