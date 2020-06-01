* Fixed the voice amplitude proportionality chart - for certain
combinations of center and sensitivity it could show an odd
discontinuous curve.
* Fixed display of Wave Type and Key Prop toggle on the Oscillator table.

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


