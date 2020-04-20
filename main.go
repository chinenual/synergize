package main

import (
        "encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	
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
	comtst = flag.Bool("COMTST", false, "run command line diagnostics rather than the GUI")
	looptst = flag.Bool("LOOPTST", false, "run command line diagnostics rather than the GUI")
	loadvce = flag.String("LOADVCE", "", "load the named VCE file into Synergy")
	loadcrt = flag.String("LOADCRT", "", "load the named CRT file into Synergy")
	savesyn = flag.String("SAVESYN", "", "save the Synergy state to the named SYN file")
	loadsyn = flag.String("LOADSYN", "", "load the named SYN file into Synergy")
	synver  = flag.Bool("SYNVER", false, "Print the firmware version of the connected Synergy")
	
	debug = flag.Bool("d", true, "enables the debug mode")
	w        *astilectron.Window
	about_w  *astilectron.Window
	l 	 *log.Logger
	AppVersion string
)

func setVersion() {
	// convert the BuiltAt string to something more useful as a version id.
	// BuiltAt looks like:
	//    "2020-04-07 09:56:14.790283 -0400 EDT m=+9.882658457"
	timestamp := strings.Split(BuiltAt," ")[0];
	// now:   "2020-04-07"
	timestamp = strings.ReplaceAll(timestamp, "-","")
	// now:   "20200407"
	AppVersion = Version + " (" + timestamp + ")"
	// now:   " 0.1.0 (20200407)"
}

func init() {
	setVersion()

	// Create logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags() | log.Lshortfile)
	
	l.Printf("Running app version %s\n", AppVersion)
	l.Printf("Default serial device is %s\n", defaultPort)
	
}

func main() {	
	// Parse flags
	flag.Parse()

	// run the command line tests instead of the Electron app:
	if *synver {
		diagInitAndPrintFirmwareID();
		os.Exit(0);
	} else if *comtst {
		diagCOMTST();
		os.Exit(0);
	} else if *looptst {
		diagLOOPTST();
		os.Exit(0);
	} else if *loadvce != "" {
		diagLoadVCE(*loadvce);
		os.Exit(0);
	} else if *loadcrt != "" {
		diagLoadCRT(*loadcrt);
		os.Exit(0);
	} else if *savesyn != "" {
		diagSaveSYN(*savesyn);
		os.Exit(0);
	} else if *loadsyn != "" {
		diagLoadSYN(*loadsyn);
		os.Exit(0);
	}

	// Run bootstrapls
	


	macOSMenus := []*astilectron.MenuItemOptions{{
			Label: astikit.StrPtr("Synergize"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About Synergize"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						err := bootstrap.SendMessage(about_w, "setVersion", AppVersion, func(m *bootstrap.MessageIn){}) 
						if err != nil {
							l.Println(fmt.Errorf("sending about event failed: %w", err))
						}
						about_w.Show()
						return
					},
				},
				{ Type: astilectron.MenuItemTypeSeparator },
				{
					Label: astikit.StrPtr("Preferences..."),
					Accelerator: astilectron.NewAccelerator("CommandOrControl+P"),
				},
				{Role: astilectron.MenuItemRoleServices},
				{ Type: astilectron.MenuItemTypeSeparator },
				{
					// Override the "Hide Electron" label
					Label: astikit.StrPtr("Hide Synergize"),
					Role: astilectron.MenuItemRoleHide,
				},
				{Role: astilectron.MenuItemRoleHideOthers},
				{Role: astilectron.MenuItemRoleUnhide},
				{ Type: astilectron.MenuItemTypeSeparator },
				{
					// Override the "Quit Electron" label
					Label: astikit.StrPtr("Quit Synergize"),
					Role: astilectron.MenuItemRoleQuit,
				},
			},
		},{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("Open File..."),
					Accelerator: astilectron.NewAccelerator("CommandOrControl+O"),
                                        OnClick: func(e astilectron.Event) (deleteListener bool) {
                                                if err := bootstrap.SendMessage(w, "fileDialog", "", func(m *bootstrap.MessageIn) {
                                                        // Unmarshal payload
                                                        var s []string
                                                        if err := json.Unmarshal(m.Payload, &s); err != nil {
                                                                l.Println(fmt.Errorf("unmarshaling payload failed: %s : %w", m.Payload, err))
                                                                return
                                                        }
                                                        l.Printf("fileDialog payload is %s!\n", s)

							vce,_ := vceReadFile(s[0]);
							if err := bootstrap.SendMessage(w, "viewVCE", vce, func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								if err := json.Unmarshal(m.Payload, &s); err != nil {
									l.Println(fmt.Errorf("unmarshaling payload failed: %s : %w", m.Payload, err))
									return
								}
								l.Printf("viewVCE payload is %s!\n", s)							
							}); err != nil {
								l.Println(fmt.Errorf("sending viewVCE event failed: %w", err))
							}
						}); err != nil {
							l.Println(fmt.Errorf("sending fileDialog event failed: %w", err))
						}
						return
					},
				},
			},
		},
		{
			Label: astikit.StrPtr("Diagnostics"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("Sanity Test..."),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if err := bootstrap.SendMessage(w, "runDiag", nil, func(m *bootstrap.MessageIn) {
							// Unmarshal payload
							var s string
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								l.Println(fmt.Errorf("unmarshaling payload failed: %s : %w", m.Payload, err))
								return
							}
							l.Printf("diag payload is %s!\n", s)							
						}); err != nil {
							l.Println(fmt.Errorf("sending viewVCE event failed: %w", err))
						}
						return
					},
				},
			},
		},
	}
	menuOptions := macOSMenus
	
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug:  *debug,
		Logger: l,
		MenuOptions: menuOptions,
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
                        w = ws[0]
			about_w = ws[1]
                        return nil
                },
		RestoreAssets: RestoreAssets,
		
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("black"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(700),
				Width:           astikit.IntPtr(700),
			},
		},{
			Homepage:       "about.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("black"),
				Center:          astikit.BoolPtr(true),
				Show:            astikit.BoolPtr(false),
				Height:          astikit.IntPtr(300),
				Width:           astikit.IntPtr(400),
				Custom: &astilectron.WindowCustomOptions{
					HideOnClose:	astikit.BoolPtr(true),
				},
			},
		}},
	}); err != nil {
		l.Fatal(fmt.Errorf("running bootstrap failed: %w", err))
	}
}
