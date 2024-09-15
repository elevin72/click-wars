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
	fmt.Println("hello world!!!")
	t, err := template.ParseFS(static, "static/index.html")
	if err != nil {
		log.Print(err)
	}

	// get state ???
	t.Execute(w, struct {
		LinePosition int32
		TotalHits    int32
	}{
		LinePosition: 0,
		TotalHits:    0,
	})
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else {
		w.Header().Set("Content-Type", "text/css")
	}
	// http.StripPrefix("/static/", http.FileServer(http.FS(static))).ServeHTTP(w, r)
	http.FileServer(http.FS(static)).ServeHTTP(w, r)
}

// func percentageOnLoadHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Printf("hit percentage.css\n")
// 	t, err := template.ParseFiles("../frontend/static/css/percentage.css.tmpl")
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	percentage := float32(LINE_POSITION/-2) + 50
// 	w.Header().Set("Content-Type", "text/css")
// 	t.Execute(w, struct {
// 		LeftSide  float32
// 		RightSide float32
// 	}{
// 		LeftSide:  100 - percentage,
// 		RightSide: percentage,
// 	})
// }

func main() {

	server := newServer()
	go server.run()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(server, w, r)
	})
	// http.HandleFunc("/static/css/percentage.css", percentageOnLoadHandler)

	port := "8080"
	fmt.Printf("Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
