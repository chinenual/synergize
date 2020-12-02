---
layout: default
title: Release Notes
nav_order: 90
description: "Release Notes"
permalink: /docs/release-notes

---

# Release Notes

## 2.3.0-beta1

* Synergize no longer uses Bonjour/Zeroconf to discover the Synergia
  virtual instrument. A much more reliable mechanism has been
  implemented. Requires version Synergia 0.97 beta or later.
    * Bonjour can still be used to configure the interface to the
      TouchOSC control surface, however this appears to not be
      particularly reliable on Windows - on Windows please manually
      configure your OSC connection.
* Change in behavior: Implements [issue #17](https://github.com/chinenual/synergize/issues/17):
  If a voice is being viewed when enabling Voice Mode, it is loaded
  into the editor rather than the default empty voice. 
* Envelope Graph display and editing:
    * Implements [issue #24](https://github.com/chinenual/synergize/issues/24):
      Points on the Envelope graphs can be dragged directly with a mouse.
      This can be a lot more intuitive than setting curves via time
      deltas.
    * Implements [issue #21](https://github.com/chinenual/synergize/issues/21):
      The time scale on the envelope chart can now be toggled between
      logarithmic and linear scales.
    * Fixes [issue #22](https://github.com/chinenual/synergize/issues/22):
      The implicit start of the amplitude envelope is now displayed as
      part of the amplitude curve on the Envelope graphs.
* Similarly, the Filter, Key Proportion and Key EQ graphs are now
  directly editable via mouse drag.  NOTE: if you drag your mouse too
  quickly, Synergize may not get a mouse event for some indexes in the
  curve (causing some index values to not be changed as you drag the
  mouse across the screen).  If you see that happen, just drag a
  little slower...
* Control Surface issues:
    * Fixes [issue #15](https://github.com/chinenual/synergize/issues/15):
      Tolerates TouchOSC configurations which have the optional
      accelerometer, touch or ping messages enabled.
    * Fixes [issue #18](https://github.com/chinenual/synergize/issues/18):
      Adds address validation for the control surface address. It was
      possible to specify an invalid address without getting any sort of
      error message.
    * Fixes [issue #20](https://github.com/chinenual/synergize/issues/20):
      Mute and solo buttons did not transmit to control surface.
  
## 2.2.0

* Fixes to the serial communications interface that results in far
  fewer problems (often reported as "timeouts" in previous versions of
  Synergize).  This is a _**significant**_ improvement, especially when
  using the iPad control surface interface. All users of older
  versions are encouraged to upgrade.
* Support for connecting to the new [Synergia](https://jariseon.github.io/synergia)
  virtual Synergy -- a chip-level emulation of the Synergy that runs
  as a VST or AU virtual instrument. 
    * During beta testing, we've noticed that Synergize sometimes does
      not discover a running Synergia instance via
      Bonjour/Zeroconf. Often, toggling the VRAM button in Synergia is
      enough to get Synergize's attention. If that does not work,
      however, you can hardcode Synergia's network port (shown in
      green text under the VRAM button) via Synergize's
      Help->Preferences menu.   We are tracking this issue as
      [issue #13](https://github.com/chinenual/synergize/issues/13). 
* Support for Bonjour/Zeroconf automatic network configuration.
    * Synergize advertises itself so that TouchOSC can connect to it with a
      single touch. 
    * It can also use zeroconf to search for available TouchOSC
      devices so you dont have to hardcode their address and port. 
    * It uses zeroconf to search for virtual instruments.
    * Zeroconf is enabled for TouchOSC discovery by default, but can be disabled
      via the Preferences menu. 
* The UI no longer shows the firmware version in the upper left
    corner of the menu bar. Instead it shows the name of the
    connected synergy and control surface.
* The default configuration for the `synergize.log` is much less verbose
* Fixes [issue #3](https://github.com/chinenual/synergize/issues/3):
  removes a limitation on macos: serial port configuration (e.g. baud rate) can be
  changed without requiring the application to be restarted.

## 2.1.0

* Adds support for a Tablet-based TouchOSC external control surface to
  support "GDS like" editing, with sliders for oscillator tuning,
  envelopes, filters, key eq and key proportionality.
* Fixes [issue #8](https://github.com/chinenual/synergize/issues/8):
  tab moves down columns on the Filter, Key Equalization and Key Proportion tables.
* Fixes [issue #9](https://github.com/chinenual/synergize/issues/9):
  the name of "standard patch 6".  SYNHCS called this
  patch `((1~2)~3) + ((1~2)~4) + ((5~6)~7) + ((5~6)~8)`,
  but in fact based on the register assignments it actually uses 
  it should have been
  called `((1+2)~3) + ((1+2)~4) + ((5+6)~7) + ((5+6)~8)`
  (a simple typo of ~ vs. +).  Synergize now uses the corrected text name.
* Fixes [issue #11](https://github.com/chinenual/synergize/issues/11):
  type 1 envelopes sometimes erroneously showed loop points.

## 2.0.0

Synergize now offers parity in functionality vs. SYNHCS (except for
some minor convenience features as documented in the
[TODO](https://github.com/chinenual/synergize/blob/master/TODO.md)).

For those using 1.0.0, this is a major upgrade: Synergize is now an
Editor as well as a Librarian.  For those testing the 2.0.0-beta, it
adds a CRT editor and fixes a number of bugs in the beta
functionality. See the release notes for 2.0.0-beta1 and 2.0.0-beta2
for other news.

* A new [User Manual](https://chinenual.github.io/synergize/docs)
  documenting how to use Synergize, and including some of the theory
  behind the Synergy/GDS voicing architecture and advice on how to get
  the most out of your patches (adapted from the original SYNHCS manual).
* Adds support for editing a CRT (clear voices, add voices and save a new .CRT file).
* Patch routing diagrams generated from the patch table to help
visualize the patch routing.
* Fixes a bug that made the Oscillator SOLO/MUTE functionality
  sometimes silently [sic] fail to alter oscillator audibility.
* Fixes a bug in the envelope copying functionality.
* Fixes a bug in some numeric input edge cases on the envelopes tab.
* Fixes a bug in the patch register editing
* Adds new serial port diagnostic command line option -LINKTST to send
  and receive bytes to the Synergy's `RESTORE Program #4` serial link
  test. See the
  [Troubleshooting](https://chinenual.github.io/synergize/docs/hardware/troubleshooting#serial-diagnostics-linktst)
  section of the manual for details.


## 2.0.0-beta2

* Beta test release.
* Adds support for altering the oscillator patch routing.
* Adds support for copying filters
* Adds support for copying envelopes
* A command-line diagnostics executable is included in the Windows
  build. (Mac and Linux builds support these command line diagnostics
  in the main executable).  See the Troubleshooting section in 
  [HARDWARE](https://chinenual.github.io/synergize/docs/hardware) for details.
  
## 2.0.0-beta1

* Beta test release.
* This release introduces voice editing functionality.  The editor is
  modeled on, but is not exactly the same as the Kaypro SYNHCS
  software:
    * Unlike SYNHCS, all editing is done in the application -
      not via buttons and knobs on the Synergy. Changes made in the editor
      are immediately transmitted to the Synergy, so you can hear the
      effect of the parameter change in real time.
	* In some cases (e.g. Vibrato, Timbre, Amplitude and Transpose
      settings), the Synergy front panel can also be used to alter the
      values - however if used, the values are not reflected in the
      Synergize display - only in the Synergy's VRAM. 
    * Most parameter editing is one-for-one the same as in SYNHCS,
      except for the Oscillator Mute and Solo features.  Instead of
      adopting the Ensemble/Group concepts from SYNHCS, Synergize uses
      Mute and Solo buttons modelled on DAW conventions (e.g. Logic
      Pro, Ableton Live - hopefully this is intuitive).
* SYNHCS functionality that is not implemented:
    * Not yet implemented, but planned:
        * This version only supports the 10 "factory" patch types; "user
	      defined" patch routing is not yet handled.
		* No support yet to copy another oscillator's
          parameters/filters, or import filters from another voice.
        * No support yet for creating new CRT's from existing VCE's.
          This will be added, but the main point of this beta is to
          get user feedback on the basic "editing" functionality and
          catch any errors in the way parameters are displayed on
          screen vs. how they affect a Synergy voice.
        * The only way to "clear" the voicing system is to turn off
          Voicing mode and restart.
    * Functionality not planned to be supported (use a github issue or
      start a discussion in the Synergize Slack group if you want to
      make a case that these should be implemented):
	    * Performance Program controller (sequence of performance
 		  parameters for live performance support).
	    * Support for creating or viewing .DOC files
		* Support for printing voicing parameters
        * Support for changing portamento settings via Synergize.
          Since these values are not stored in the voice file (VCE or
          CRT), you must change them via the Synergy front panel.
		* "Scaling" the envelope and filter displays.  In Synergize,
          these auto-scale - there's no support for zooming in or out.
* Other known limitations with this "beta" release:
    * The menu and buttons are subject to change and improvement.
    * There is an as-yet undiagnosed bug in the CRT parser: the
       VCART6.CRT can't be loaded.
	* When you "save" a voice, you are saving what you hear in the
      Synergy. If there is a serial communication error during
      editing, the values displayed on the screen may not reflect what
      has been loaded into the Synergy. I will add some sort of
      mechanism to re-sync the display with the state of the
      Synergy. But for now, if you get a timeout during editing, just
      be aware that the screen no longer displays reality. 
*  Fixes some miscellaneous bugs found while implementing the editor:
    * Fixed the voice amplitude proportionality chart - for certain
      combinations of center and sensitivity it could show an odd
      discontinuous curve.
    * Fixed display of Wave Type and Key Prop toggle on the Oscillator table.
    * Fixed display of Sustain and Loop points.

## 1.0.0

* Voice settings (envelopes, filters, key equalization, key
  proportionality) are displayed in tabular and graphical formats
  inspired by the SYNHCS editor.
* Default window size is a little taller to accomodate new voice
  settings screens.
* Tweaked timeout settings to reduce spurious communications timeouts.
* Linux caveats:
     * Linux testing remains problematic. 
        * My virtualized amd64 environment continues to show problems
          with the Save SYN functionality (i.e. downloads from the
          Synergy show problems, uploads to it are reliable).  I am
          suspicious of the virtualized serial port config, so can't
          be sure there's really a problem.
        * So, I've recycled an old laptop with an old 32-bit version of
          Ubuntu and tested the 386 version: it's as reliable as the
          Mac and Windows ports.   Further evidence that the 64-bit
          version might also be OK and that my test setup is to blame.
        * I have no ability to test the ARM port at all - but know
          that there is interest in running on raspberry pi, so I'm
          building  it in hopes it works for those that try it.
	* The UI dialog notification glitch noted in release 0.2.0
       remains; as noted before, it is annoying but harmless.
	* On Linux, file dialogs often open _behind_ the main window.

## 0.2.0

* Beta test release.
* Fixed serial I/O bugs; communications are much more reliable.
* Adds preliminary support for Linux
    * Note: the tar.gz contains only an executable, no meta files for desktop menu setup.
    * The Linux UI exposes a bug in the underlying UI framework: when
	  invoking a file dialog, you may see a brief popup dialog saying
      that the framework "is ready".  This is harmless, but annoying.
	  Mac and Windows do not share this problem.
	* Serial communications behavior on Linux is suspect. In my test setup,
	  loading data to the Synergy is quite reliable, but I get frequent
	  errors when attempting to fetch data (the Save Synergy State
	  functionality).  I'm not sure if this is due to a problem with the
	  program or with the hardware in my test setup, so I'm
	  releasing this with these caveats for anyone willing to help out
	  testing on Linux.

## 0.1.0

* Beta test release.
* Supports both Windows and MacOS.
    * Note: the mac release is not "signed", so you will need to
      explicitly allow the app from an "unidentified developer" to run on
      your machine (see https://support.apple.com/guide/mac-help/open-a-mac-app-from-an-unidentified-developer-mh40616/mac).
* Basic librarian functionality:
    * Save and Load "Synergy State" (.SYN) files
    * View and Load "Cartridge" (.CRT) files, including voice details
    * View "Voice" (.VCE) files
    * Disable VRAM
    * Run basic serial port diagnostic test ("Restore+Program 4" AKA COMTST)


