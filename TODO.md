# For MVP

* review copyright language for DK/Mulogix references

disable VRAM function
* connect remaining functions to main buttons
* remove old main menu items except for About and Preferences ?

* re-connect after serial port or baud preferences change
* update left-browser pain after library path preference change

* seg faults if not connected
* refresh connection status after an IPC call (which might implicitly connect)

* windows file browser needs to be able to change drives
drives

** find platform specific default locations for log (currently just CWD)

* remove Debug option from menus (enable via preferences)

* windows variant for main menus

* DONE - auto select a serial port
* connect all functionality into GUI
* DONE - serial timeout behavior - review goroutines
* DONE - opcode "drain" input logic

* Test on windows

* display spinner without shifting content

* packaging
** DMG for Macos
** MSI package on Windows

* simple website landing page

* github "release"

* unit tests

* error handling - do reasonable things if synergy not connected/io fails
