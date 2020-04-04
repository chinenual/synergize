package main

import (
	"github.com/zserge/webview"

//	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var appName = "Synergize"
var windowWidth, windowHeight = 1024, 700

var rootDir string                           // current directory


func init() {
	var err error
	rootDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	// get source locations in log
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// channel to get the web prefix
	prefixChannel := make(chan string)
	// run the web server in a separate goroutine
	go app(prefixChannel)
	prefix := <- prefixChannel

	// create a web view
	debug := true
	w := webview.New(debug)
	defer w.Destroy()

	w.SetTitle(appName);
	w.SetSize(windowWidth, windowHeight, webview.HintNone);
	w.Navigate(prefix + "/index")
	w.Run();

}


func app(prefixChannel chan string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/index", indexHandler);
	mux.HandleFunc("/view", viewHandler);
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(rootDir+"/static"))))
//	mux.HandleFunc("/start", start)
//	mux.HandleFunc("/frame", getFrame)
	//	mux.HandleFunc("/key", captureKeys)
	// get an ephemeral port, so we're guaranteed not to conflict with anything else
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	portAddress := listener.Addr().String()
	prefixChannel <- "http://" + portAddress
	listener.Close()
	server := &http.Server{
		Addr:    portAddress,
		Handler: mux,
	}
	server.ListenAndServe()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
 log.Println("indexHandler")
	t,err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println("err: ", err)
	}
	t.Execute(w, nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
 log.Println("viewHandler")
	t,err := template.ParseFiles("template/view.html")
	if err != nil {
		log.Println("err: ", err)
	}
	var vce VCE
	vce,err = ReadVCEFile("VOICES/CARLOS1/CLIK.VCE")
	if err != nil {
		log.Println("err: ", err)
	}
 log.Println("viewHandler ", VCEToJSON(vce))

/*	b,_ := json.Marshal(vce)
	json := string(b)

	data := struct {
		Json string
	}{ Json: json}
*/
	t.Execute(w, vce)
}
