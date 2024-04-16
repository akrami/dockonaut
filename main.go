package main

import (
	"akrami/dockonaut/src/docker"
	"errors"
	"flag"
	"log"
)

func main() {

	_, err := docker.CommandExecute("docker", "version")
	if err != nil {
		panic(errors.New("docker is not running"))
	}

	configFile := flag.String("config", "config.json", "path to your config file")
	flag.Parse()

	repos, errLoad := docker.Load(*configFile)
	if errLoad != nil {
		panic(errLoad)
	}
	switch flag.Arg(0) {
	case "start":
		err := docker.Start(repos)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Started!")

	case "stop":
		errStop := docker.Stop(repos)
		if errStop != nil {
			log.Fatalln(errStop)
		}

	case "restart":
		err := docker.Restart(repos)
		if err != nil {
			log.Fatalln(err)
		}

	case "purge":
		err := docker.Purge(repos)
		if err != nil {
			log.Fatalln(err)
		}

	default:
		log.Fatalln("please provide an start/stop/restart/purge arg")
	}
}
