package web

import (
	"akrami/dockonaut/internal/docker"
	"html/template"
	"net/http"
)

func HandleHome(writer http.ResponseWriter, request *http.Request) {
	config, _ := docker.Load("config.json")
	tmpl, _ := template.ParseFiles("tmpl/home.html")
	tmpl.Execute(writer, config)
}
