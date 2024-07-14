package main

import (
	"akrami/dockonaut/internal/docker"
	"akrami/dockonaut/internal/engine"
	"errors"
	"flag"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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
	flag.Parse()

	config, errLoad := docker.Load(*configFile)
	if errLoad != nil {
		panic(errLoad)
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
