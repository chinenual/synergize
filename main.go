package main

import (
	"flag"
	"fmt"
	"log"
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
	debug = flag.Bool("d", true, "enables the debug mode")
	w        *astilectron.Window
	about_w  *astilectron.Window
	AppVersion string
)

func setVersion() {
	// convert the BuiltAt string to something more useful as a version id.
	// BuiltAt looks like:
	//    "2020-04-07 09:56:14.790283 -0400 EDT m=+9.882658457"
	AppVersion = strings.Split(BuiltAt,".")[0];
	// now:   "2020-04-07 09:56:14"
	AppVersion = strings.ReplaceAll(AppVersion, ":","")
	AppVersion = strings.ReplaceAll(AppVersion, "-","")
	AppVersion = strings.ReplaceAll(AppVersion, " ","")
	// now:   "20200407095614"
}

func main() {
	setVersion()
	
	// Parse flags
	flag.Parse()

	// Create logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags() | log.Lshortfile)
	
	// Run bootstrapls
	
	l.Printf("Running app version %s\n", AppVersion)

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
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						err := bootstrap.SendMessage(about_w, "setVersion", AppVersion, func(m *bootstrap.MessageIn){}) 
						if err != nil {
							l.Println(fmt.Errorf("sending about event failed: %w", err))
						}
						about_w.Show()
						return
					},
				},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
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
				Width:           astikit.IntPtr(1024),
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
			},
		}},
	}); err != nil {
		l.Fatal(fmt.Errorf("running bootstrap failed: %w", err))
	}
}
