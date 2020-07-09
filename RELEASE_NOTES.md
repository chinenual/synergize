# 2.0.0-beta2

* Beta test release.
* Adds support for altering the oscillator patch routing.
* Adds support for copying filters
* Adds support for copying envelopes
* A command-line diagnostics executable is included in the Windows
  build. (Mac and Linux builds support these command line diagnostics
  in the main executable).  See the Troubleshooting section in 
  [HARDWARE](docs/HARDWARE.md) for details.
  
# 2.0.0-beta1

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

# 1.0.0

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

# 0.2.0

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

# 0.1.0

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


