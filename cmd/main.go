package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"

	"net/http"
	"strings"
)

const (
	LONG  byte = 0
	SHORT byte = 1
)

//go:embed static/*
var static embed.FS

func rootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(static, "static/index.html")
	if err != nil {
		log.Print(err)
	}

	t.Execute(w, struct {
		LinePosition int32
		TotalHits    int32
	}{
		LinePosition: linePosition.Load(),
		TotalHits:    totalHits.Load(),
	})
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(path, "css") {
		w.Header().Set("Content-Type", "text/css")
	} else {
		fmt.Println("something else")
	}

	http.FileServer(http.FS(static)).ServeHTTP(w, r)
}

func percentageOnLoadHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(static, "static/css/percentage.css.tmpl")
	if err != nil {
		log.Print(err)
	}
	w.Header().Set("Content-Type", "text/css")
	percentage := float32(linePosition.Load()/-2) + 50
	t.Execute(w, struct {
		LeftSide  float32
		RightSide float32
	}{
		LeftSide:  100 - percentage,
		RightSide: percentage,
	})
}

func main() {

	server := newServer()
	go server.run()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/static/css/percentage.css", percentageOnLoadHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(server, w, r)
	})

	port := "8080"
	fmt.Printf("Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
