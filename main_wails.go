package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

var wailsApp *application.App
var wailsPrefsWindow *application.WebviewWindow
var wailsAboutWindow *application.WebviewWindow

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func WailsMain() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.

	wailsApp = application.New(application.Options{
		Name:        "Synergize",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&UIService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Synergize",
		Mac: application.MacWindow{
			Backdrop: application.MacBackdropTranslucent,
		},
		BackgroundColour: application.NewRGB(0, 0, 0),
		URL:              "/",
		Width:            990,
		Height:           900,
		DevToolsEnabled:  true,
	})
	wailsPrefsWindow = wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Synergize Preferences",
		URL:    "/prefs.html",
		Height: 680,
		Width:  800,
		Hidden: true,
	})
	wailsAboutWindow = wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:           "About Synergize",
		URL:             "/about.html",
		Height:          470,
		Width:           500,
		InitialPosition: application.WindowCentered,
		Hidden:          true,
	})

	/**
		// Create a goroutine that emits an event containing the current time every second.
		// The frontend can listen to this event and update the UI accordingly.
		go func() {
			for {
				now := time.Now().Format(time.RFC1123)
				app.Event.Emit("time", now)
				time.Sleep(time.Second)
			}
		}()
	**/
	// Run the application. This blocks until the application has been exited.
	err := wailsApp.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
