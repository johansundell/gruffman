package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/kardianos/service"
)

var logger service.Logger

type program struct {
	exit chan struct{}
}

const appVersionStr = "v1.0"
const filenameSettings = "settings.json"

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		defaultHandler,
	},
}

var router *mux.Router
var settings appSettings

type appSettings struct {
	FileMakerHost     string `json:"filemakerhost"`
	FileMakerUser     string `json:"filemakeruser"`
	FileMakerPassword string `json:"filemakerpassword"`
	WebDomain         string `json:"webdomain"`
	Port              string `json:"port"`
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {
	logger.Infof("I'm running %v.", service.Platform())

	router = NewRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    settings.Port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Println(srv.ListenAndServe())
	}()

	for {
		select {
		case <-p.exit:
			srv.Close()
			return nil
		}
	}
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "gruffman",
		DisplayName: "Gruffman",
		Description: "Gruffman is a mircoservice",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dir, _ := filepath.Split(ex)
	dat, err := ioutil.ReadFile(dir + filenameSettings)
	if err != nil {
		data, _ := json.Marshal(settings)
		ioutil.WriteFile(dir+filenameSettings, data, 0664)
		log.Fatal("settings.json missing, " + err.Error())
	}

	if err := json.Unmarshal(dat, &settings); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<html><body>We are up and running ;)")
	url, err := router.GetRoute("offert").URL("id", "123")
	if err == nil {
		fmt.Fprint(w, "<a href=\""+fixLink(url.RequestURI())+"\">test</a></body></html>")
	}

}

func fixLink(part string) string {
	return settings.WebDomain + part
}

func getDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir, _ := filepath.Split(ex)
	return dir, nil
}
