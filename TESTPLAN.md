# UI tests

* remove preferences.json
* start the app
     * library should be in user's home dir
* open preference menu
	* serial port should be platform-specific default
	* serial baud should be 9600
	* library should be empty
* change library path and SAVE
    * preferences.json should be written
    * nav pane should show new directory
* Quit App
* start the App
* open Preferences menu
	* verify that the settings were preserved
* Help->About
	* should show correct date and version
* turn on synergy
* Connect->Connect to Synergy
	* firmware version should be displayed
* view a CRT
	* voices should display in table
* click a voice
	* voice details and sensativity graph should display
* Back to CRT
* Load CRT
	* sanity check at least one of the voices corresponds to the CRT
	* Synergy Cartridge light should be blinking
* Disable VRAM
	* Synergy Cartridge light stops blinking
	* sanity check voices are "internal"
* Save SYN
	* verify no errors
* Load SYN
	* verify no errors 
* Connect->Serial Diagnostics
	* turn on COMTST on Synergy (Restore+Program 4)
	* run the test
* power cycle the synergy

	
# Commnd line options

* -SYNVER
* -SAVESYN
* -LOADSYN
* -LOADCRT
* -LOADVCE
* -COMTST
* -LOOPTST

Run the integration tests: make itest

Run SAVE/LOAD options in a loop to try to catch intermittent errors.  Pass count of 1 for a quick sanity test:

```
testbin/testmac [count]
testbin/testwindows [count]
testbin/testlinux [count]
```
