package main

import (
	"akrami/dockonaut/internal/docker"
	"akrami/dockonaut/internal/engine"
	"akrami/dockonaut/internal/web"
	"errors"
	"flag"
	"log"
	"os"

	"net/http"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gorilla/mux"
)

func main() {
	_, err := docker.CommandExecute("docker", "version")
	if err != nil {
		panic(errors.New("docker is not running"))
	}

	headerChannel := make(chan engine.Header)
	subtextChannel := make(chan engine.Subtext)
	program := tea.NewProgram(engine.InitialModel())

	configFile := flag.String("config", "config.json", "path to your config file")
	daemonFlag := flag.Bool("daemon", false, "run as a web daemon")
	flag.Parse()

	config, errLoad := docker.Load(*configFile)
	if errLoad != nil {
		panic(errLoad)
	}

	if *daemonFlag {
		router := mux.NewRouter()
		router.Use(logMiddleware)
		router.HandleFunc("/", web.HandleHome)
		http.Handle("/", router)

		webErr := http.ListenAndServe(":9090", nil)
		if webErr != nil {
			panic(errors.Join(errors.New("can not start server"), webErr))
		}
	}
	
	go func() {
		for header := range headerChannel {
			program.Send(header)
		}
	}()

	go func() {
		for subtext := range subtextChannel {
			program.Send(subtext)
		}
	}()

	go func() {
		switch flag.Arg(0) {
		case "start":
			if err := config.Start(headerChannel, subtextChannel); err != nil {
				log.Fatalln(err)
			}

		case "stop":
			if err := config.Stop(headerChannel, subtextChannel); err != nil {
				log.Fatalln(err)
			}

		case "restart":
			if err := config.Restart(headerChannel, subtextChannel); err != nil {
				log.Fatalln(err)
			}

		case "purge":
			if err := config.Purge(headerChannel, subtextChannel); err != nil {
				log.Fatalln(err)
			}

		default:
			log.Fatalln("Please provide start/stop/restart/purge as argument")
			os.Exit(1)
		}

		program.Send(engine.Done(true))
		program.Send(tea.Quit())
	}()

	if _, err := program.Run(); err != nil {
		log.Fatalln("Error", err)
		os.Exit(1)
	}
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.RemoteAddr, request.Method, request.RequestURI)
		next.ServeHTTP(writer, request)
	})
}
