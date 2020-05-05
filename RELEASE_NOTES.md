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


