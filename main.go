package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/matishsiao/goInfo"

	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/chinenual/synergize/data"
	"github.com/chinenual/synergize/logger"
	"github.com/chinenual/synergize/osc"
	"github.com/chinenual/synergize/synio"
	"github.com/chinenual/synergize/zeroconf"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

// Vars injected via ldflags by bundler
var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

// Application Vars
var (
	port              = flag.String("port", getDefaultPort(), "the serial device")
	baud              = flag.Uint("baud", getDefaultBaud(), "the serial baud rate")
	vst               = flag.Uint("vst", 0, "port for the VST instrument")
	serial            = flag.Bool("serial", false, "use the serial port even if a VST is available")
	record            = flag.String("RECORD", "", "capture bytes to <record>.in and <record>.out")
	uitest            = flag.Int("UITEST", 0, "alter startup to support automated testing (specifies listening port)")
	loglevel          = flag.String("LOGLEVEL", "INFO", "Set log level to DEBUG")
	provisionOnly     = flag.Bool("PROVISION", false, "run the provisioner and then exit")
	serialVerboseFlag = flag.Bool("SERIALVERBOSE", false, "Show each byte operation through the serial port")
	verboseOscIn      = flag.Bool("OSCINVERBOSE", false, "Show OSC input events")
	verboseOscOut     = flag.Bool("OSCOUTVERBOSE", false, "Show OSC output events")
	mockSynio         = flag.Bool("MOCKSYNIO", false, "Mock the Synergy I/O for testing")
	comtst            = flag.Bool("COMTST", false, "run command line diagnostics rather than the GUI")
	looptst           = flag.Bool("LOOPTST", false, "run command line diagnostics rather than the GUI")
	linktst           = flag.Bool("LINKTST", false, "run command line diagnostics rather than the GUI")
	savevce           = flag.String("SAVEVCE", "", "save the Synergy VRAM state to the named VCE file ")
	loadvce           = flag.String("LOADVCE", "", "load the named VCE file into Synergy")
	loadcrt           = flag.String("LOADCRT", "", "load the named CRT file into Synergy")
	savesyn           = flag.String("SAVESYN", "", "save the Synergy state to the named SYN file")
	loadsyn           = flag.String("LOADSYN", "", "load the named SYN file into Synergy")
	synver            = flag.Bool("SYNVER", false, "Print the firmware version of the connected Synergy")
	//	rawlog            = flag.Bool("RAWLOG", false, "Turn off timestamps to make logs easier to compare")
	//midiproxy = flag.Bool("MIDIPROXY", false, "present a MIDI interface and use serial IO to control the Synergy")

	w       *astilectron.Window
	about_w *astilectron.Window
	prefs_w *astilectron.Window
	a       *astilectron.Astilectron
	//l          logger.Logger //*log.Logger
	AppVersion string
	OsVersion  string
)

func getDefaultBaud() uint {
	// FIXME: loads the prefs twice - harmless, but annoying
	prefsLoadPreferences()

	if prefsUserPreferences.SerialBaud != 0 {
		return prefsUserPreferences.SerialBaud
	}
	return 9600
}

func getDefaultPort() string {
	// FIXME: loads the prefs twice - harmless, but annoying
	prefsLoadPreferences()

	if prefsUserPreferences.SerialPort != "" {
		return prefsUserPreferences.SerialPort
	}
	if runtime.GOOS == "darwin" {
		files, _ := filepath.Glob("/dev/tty.usbserial*")
		for _, f := range files {
			return f
		}
	} else if runtime.GOOS == "linux" {
		files, _ := filepath.Glob("/dev/ttyUSB*")
		for _, f := range files {
			return f
		}
		// if no USB serial, assume /dev/ttyS0
		return "/dev/ttyS0"

	} else {
		// windows
		return "COM1"
	}
	return ""
}

func setVersion() {
	// convert the BuiltAt string to something more useful as a version id.
	// BuiltAt looks like:
	//    "2020-04-07 09:56:14.790283 -0400 EDT m=+9.882658457"
	timestamp := strings.Split(BuiltAt, " ")[0]
	// now:   "2020-04-07"
	timestamp = strings.ReplaceAll(timestamp, "-", "")
	// now:   "20200407"
	AppVersion = Version + " (" + timestamp + ")"
	// now:   " 0.1.0 (20200407)"
	data.SetAppVersion(AppVersion)
}

// platform specific config to ensure logs and preferences go to reasonable locations
func getWorkingDirectory() (path string) {
	// don't do this if we are running from the source tree
	_, err := os.Stat("bundler.json")
	if !os.IsNotExist(err) {
		// running from source directory
		path = "."
		return
	}
	_, err = os.Stat("../bundler.json")
	if !os.IsNotExist(err) {
		// running from uitest directory
		path = "."
		return
	}

	path, _ = os.UserConfigDir()
	path = path + "/Synergize"

	// create it if necessary
	_ = os.MkdirAll(path, os.ModePerm)
	return
}

func mainInit() {
	var level logger.Level
	var levelMsg = ""
	switch *loglevel {
	case "DEBUG":
		level = logger.LevelDebug
	case "INFO":
		level = logger.LevelInfo
	case "WARN":
		level = logger.LevelWarn
	case "ERROR":
		level = logger.LevelWarn
	default:
		level = logger.LevelInfo
		levelMsg = fmt.Sprintf("Invalid value for --LOGLEVEL (%s).  Defaulting to INFO\n", *loglevel)
	}
	logger.Init(getWorkingDirectory()+"/synergize.log", level)
	if levelMsg != "" {
		logger.Warn(levelMsg)
	}
	setVersion()
	logger.Infof("Running app version %s\n", AppVersion)
	// log some info about the operating system.  This uses uname on linux and macos, and ver on windows, so info is limited.
	// for windows, version reports via ver:
	//    Core:"10.0.18363.1016",
	//      the 18363 corresponds to the Build version
	// for macos, it reports the darwin version:
	//    Core:"18.8.0"
	//      18.* = Mojave
	//      19.* = Catalina
	//      20.* = Big Sure
	// for linux, it reports the kernel
	//    Core:"5.4.0-52-generic"
	gi := goInfo.GetInfo()
	OsVersion = gi.Core
	logger.Infof("Operating environment: %#v\n", gi)
}

func refreshNavPane(path string) {
	if err := bootstrap.SendMessage(w, "explore", path, func(m *bootstrap.MessageIn) {}); err != nil {
		logger.Error(fmt.Errorf("sending refreshNav event failed: %w", err))
	}
}

func recordIo(f func(string) error, arg string) (err error) {
	// nil means use "preferences" config
	prefsUserPreferences.UseSerial = true
	if err = ConnectToSynergy(nil); err != nil {
		return
	}
	if *record != "" {
		logger.Infof(" --- start recording\n")
		synio.Conn().StartRecord()
	}
	if err = f(arg); err != nil {
		return
	}
	if *record != "" {
		logger.Infof(" --- end recording\n")
		in, out := synio.Conn().GetRecord()
		err = ioutil.WriteFile(*record+".in", in, 0644)
		if err != nil {
			logger.Errorf("ERROR: failed to write in bytes: %v\n", err)
			return
		}
		err = ioutil.WriteFile(*record+".out", out, 0644)
		if err != nil {
			logger.Errorf("ERROR: failed to write out bytes: %v\n", err)
			return
		}
		logger.Infof("Saved %d bytes to %s and %d bytes to %s\n", len(in), *record+".in", len(out), *record+".out")
		return
	}
	return
}

func main() {
	// Parse flags
	flag.Parse()
	mainInit()

	// if we read something different off the command line, set it
	// (but dont persist it) in the preferences object:
	prefsUserPreferences.SerialPort = *port
	prefsUserPreferences.SerialBaud = *baud

	logger.Infof("Default serial device is %s at %d baud\n",
		prefsUserPreferences.SerialPort,
		prefsUserPreferences.SerialBaud)

	//	if *rawlog {
	//		l.SetFlags(0)
	//		log.SetFlags(0)
	//	}

	var err error
	{
		var code = 0
		// run the command line tests instead of the Electron app:
		if *synver {
			// nil means use "preferences" config
			prefsUserPreferences.UseSerial = true
			if err = ConnectToSynergy(nil); err != nil {
				return
			}
			if err = diagInitAndPrintFirmwareID(); err != nil {
				code = 1
				logger.Error(err)
			}
			os.Exit(code)
		} else if *comtst {
			// nil means use "preferences" config
			prefsUserPreferences.UseSerial = true
			if err = ConnectToSynergy(nil); err != nil {
				return
			}
			diagCOMTST()
			os.Exit(0)
		} else if *looptst {
			// nil means use "preferences" config
			prefsUserPreferences.UseSerial = true
			if err = ConnectToSynergy(nil); err != nil {
				return
			}
			diagLOOPTST()
			os.Exit(0)
		} else if *linktst {
			// nil means use "preferences" config
			prefsUserPreferences.UseSerial = true
			if err = ConnectToSynergy(nil); err != nil {
				return
			}
			diagLINKTST()
			os.Exit(0)
		} else if *savevce != "" {
			if err = recordIo(diagSaveVCE, *savevce); err != nil {
				code = 1
				logger.Error(err)
			}
			os.Exit(code)
		} else if *loadvce != "" {
			if err = recordIo(diagLoadVCE, *loadvce); err != nil {
				code = 1
				logger.Error(err)
			}
			os.Exit(code)
		} else if *loadcrt != "" {
			if err = recordIo(diagLoadCRT, *loadcrt); err != nil {
				code = 1
				logger.Error(err)
			}
			os.Exit(code)
		} else if *savesyn != "" {
			if err = recordIo(diagSaveSYN, *savesyn); err != nil {
				code = 1
				logger.Error(err)
			}
			os.Exit(code)
		} else if *loadsyn != "" {
			if err = recordIo(diagLoadSYN, *loadsyn); err != nil {
				code = 1
				logger.Error(err)
			}
			os.Exit(code)
		}
	}

	if prefsUserPreferences.UseOsc {
		// always advertise on zeroconf even if we're not going to listen for connections
		if err := zeroconf.StartServer(prefsUserPreferences.OscPort, prefsSynergyName()); err != nil {
			logger.Errorf("could not start zeroconf: %v\n", err)
		}
		defer zeroconf.CloseServer()
	}
	if prefsUserPreferences.UseOsc && prefsUserPreferences.OscAutoConfig {
		zeroconf.StartListener()
	}

	macOSMenus := []*astilectron.MenuItemOptions{{
		Label: astikit.StrPtr("Synergize"),
		SubMenu: []*astilectron.MenuItemOptions{
			{
				Label: astikit.StrPtr("About Synergize"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					if err := bootstrap.SendMessage(about_w, "setVersion", AppVersion, func(m *bootstrap.MessageIn) {}); err != nil {
						logger.Errorf("sending about event failed: %w", err)
					}
					about_w.Show()
					return
				},
			},
			{Type: astilectron.MenuItemTypeSeparator},
			{
				Label:       astikit.StrPtr("Preferences..."),
				Accelerator: astilectron.NewAccelerator("CommandOrControl+P"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					prefs_w.Show()
					return
				},
			},
			{Role: astilectron.MenuItemRoleServices},
			{Type: astilectron.MenuItemTypeSeparator},
			{
				// Override the "Hide Electron" label
				Label: astikit.StrPtr("Hide Synergize"),
				Role:  astilectron.MenuItemRoleHide,
			},
			{Role: astilectron.MenuItemRoleHideOthers},
			{Role: astilectron.MenuItemRoleUnhide},
			{Type: astilectron.MenuItemTypeSeparator},
			{
				// Override the "Quit Electron" label
				Label: astikit.StrPtr("Quit Synergize"),
				Role:  astilectron.MenuItemRoleQuit,
			},
		},
	},

		{Role: astilectron.MenuItemRoleEditMenu},
		{Role: astilectron.MenuItemRoleWindowMenu},
		/*
			Label: astikit.StrPtr("Edit"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Role: astilectron.MenuItemRoleUndo},
				{Role: astilectron.MenuItemRoleRedo},
				{ Type: astilectron.MenuItemTypeSeparator },
				{Role: astilectron.MenuItemRoleCut},
				{Role: astilectron.MenuItemRoleCopy},
				{Role: astilectron.MenuItemRolePaste},
				{Role: astilectron.MenuItemRoleSelectAll},`
			},
		*/
	}

	var menuOptions = []*astilectron.MenuItemOptions{}
	if runtime.GOOS == "darwin" {
		menuOptions = macOSMenus
	} else {
		// leave empty for windows and linux
		// FIXME: empty menus causes the bootstrap to crash - so add the menus as a workaround
		//		menuOptions = macOSMenus
	}

	var executer = astilectron.DefaultExecuter
	var acceptTimeout = astilectron.DefaultAcceptTCPTimeout
	var adapter bootstrap.AstilectronAdapter = nil
	var astiPort = 0

	if *uitest != 0 {
		astiPort = *uitest

		executer = func(l astikit.SeverityLogger, a *astilectron.Astilectron, cmd *exec.Cmd) (err error) {
			l.Infof("======= NOT STARTING CMD %s\n", strings.Join(cmd.Args, " "))
			return
		}

		acceptTimeout = time.Minute * 3

		adapter = func(a *astilectron.Astilectron) {
			logger.Infof("======= In UI Test adapter - suppressing executor\n")
			a.SetExecuter(executer)
		}
	}

	defer func() {
		fmt.Printf("Close Event.\n")
		if err = osc.Quit(); err != nil {
			logger.Error(err)
		}
	}()

	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		Adapter:  adapter,
		AstilectronOptions: astilectron.Options{
			TCPPort:            &astiPort,
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,

			AcceptTCPTimeout: acceptTimeout,

			//ElectronSwitches: []string{
			//	"enable-logging",
			//	"no-sandbox",
			//	"remote-debugging-port", "9315",
			//	"host-rules", "MAP * 127.0.0.1",
			//},
		},
		Debug:       prefsUserPreferences.HTTPDebug,
		Logger:      logger.GetLogger(),
		MenuOptions: menuOptions,
		OnWait: func(as *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			if *provisionOnly {
				logger.Infof("Provisioning completed. Exiting.\n")
				// Quit causes a segv.  Exiting without it sometimes leaves a dialog open
				// the later is more compatible with github CI
				//a.Quit();
				os.Exit(0)
			}
			a = as
			w = ws[0]
			about_w = ws[1]
			prefs_w = ws[2]

			if err = osc.OscRegisterBridge(w); err != nil {
				logger.Error(err)
			}

			// Need to explicitly intercept Closed event on the main
			// window since the about window is never closed - only hidden.
			w.On(astilectron.EventNameWindowEventClosed, func(e astilectron.Event) (deleteListener bool) {
				a.Quit()
				return true
			})
			return nil
		},
		RestoreAssets: RestoreAssets,

		// Suppress warnings about "signal urgent I/O condition"
		// (https://github.com/asticode/go-astilectron/issues/239):
		IgnoredSignals: ignoredSignals,

		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("black"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(900),
				Width:           astikit.IntPtr(990),
				WebPreferences: &astilectron.WebPreferences{
					EnableRemoteModule: astikit.BoolPtr(true),
				},
			},
		}, {
			Homepage:       "about.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("black"),
				Center:          astikit.BoolPtr(true),
				Show:            astikit.BoolPtr(false),
				Height:          astikit.IntPtr(470),
				Width:           astikit.IntPtr(500),
				Custom: &astilectron.WindowCustomOptions{
					HideOnClose: astikit.BoolPtr(true),
				},
				WebPreferences: &astilectron.WebPreferences{
					EnableRemoteModule: astikit.BoolPtr(true),
				},
			},
		}, {
			Homepage:       "prefs.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#ccc"),
				Center:          astikit.BoolPtr(true),
				Show:            astikit.BoolPtr(false),
				Height:          astikit.IntPtr(680),
				Width:           astikit.IntPtr(800),
				Custom: &astilectron.WindowCustomOptions{
					HideOnClose: astikit.BoolPtr(true),
				},
				WebPreferences: &astilectron.WebPreferences{
					EnableRemoteModule: astikit.BoolPtr(true),
				},
			},
		}},
	}); err != nil {
		logger.Errorf("running bootstrap failed: %w", err)
	}
}
