package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type IndexData struct {
	Title      string
	Adsense_ID string
}

var (
	IP            = flag.String("ip", "0.0.0.0", "Which ip to listen on")
	Port          = flag.Int("port", 3333, "Server listen port")
	AID           = flag.String("googleid", "", "Google Adsense ID")
	Log           = flag.String("log", "website.log", "File to log to")
	Help          = flag.Bool("help", false, "Show help")
	Foreground    = flag.Bool("f", false, "Start in foreground mode, don't log to file")
	IndexTemplate *template.Template
	err           error
)

func index(w http.ResponseWriter, r *http.Request) {
	data := IndexData{
		Title:      "Joe Reid",
		Adsense_ID: *AID,
	}

	IndexTemplate.Execute(w, data)

	log.Println(r.Method, r.RemoteAddr, r.URL.Path, r.UserAgent())
}

func main() {
	r := mux.NewRouter()

	// catch any requests for "real" files (css, js, images, etc)
	r.PathPrefix("/assets/").Handler(http.FileServer(http.Dir("static/")))

	r.HandleFunc("/", index)

	addr := fmt.Sprintf("%s:%d", *IP, *Port)

	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Printf("Attempting to listen on %s\n", addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.Parse()

	if *Help {
		flag.Usage()
		os.Exit(1)
	}

	if !*Foreground {
		logfile, err := os.OpenFile(*Log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		log.SetOutput(logfile)
	}

	IndexTemplate, err = template.ParseFiles("template/index.tmpl")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
