# 0.2.0

* Fixed a serial IO bug that makes communications more reliable.

* Adds preliminary support for Linux
    * Note: the tar.gz contains only an executable, not meta files for desktop menu setup.

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
