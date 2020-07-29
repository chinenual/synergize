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

See the [User Manual](voice-library.md) for sites where you can download the original DK/Mulogix voice library.

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


## Known Issues

See [HARDWARE](hardware.md) for some serial port related troubleshooting guidelines.

* On MacOS, if you change serial parameters, you will need to restart the application in order to "reconnect" to the synergy with those parameters. On Windows, you can directly re-connect via the Connect->Connect to Synergy menu.

* Due to the mysteries of serial port communication, attempts to save or load files to the Synergy will sometimes fail (often reporting a TIMEOUT).  If this happens, a second try will usually succeed. Fixes 

* See the
[Release Notes](https://github.com/chinenual/synergize/releases)
for some caveats regarding the Linux builds.


