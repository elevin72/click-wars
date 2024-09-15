package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"

	// "os"
	// "strconv"

	"net/http"
	"strings"
)

const (
	LONG  byte = 0
	SHORT byte = 1
)

//go:embed static/*
var static embed.FS

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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(static, "static/index.html")
	if err != nil {
		log.Print(err)
	}

	linePosition := linePosition.Load()
	totalHits := totalHits.Load()
	i := IndexParams{
		LinePosition: linePosition,
		TotalHits:    totalHits,
		LeftCount:    float32(totalHits+linePosition) / 2,
		RightCount:   float32(totalHits-linePosition) / 2,
	}
	t.Execute(w, i)
}

type IndexParams struct {
	LinePosition int32
	TotalHits    int32
	LeftCount    float32
	RightCount   float32
}

type Widths struct {
	LeftWidth  float32
	RightWidth float32
}

func calcPercentageWidths() Widths {
	// scalingConstant, ok := os.LookupEnv("SCALE")
	// if !ok {
	// 	scalingConstant = "1"
	// }
	// scalingConstantFloat, err := strconv.ParseFloat(scalingConstant, 32)
	// if err != nil {
	// 	fmt.Print(err)
	// }
	scalingConstantFloat := 2
	percentage := (float32(linePosition.Load()) * float32(scalingConstantFloat)) + 50
	return Widths{
		LeftWidth:  percentage,
		RightWidth: 100 - percentage,
	}
}

func percentageOnLoadHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(static, "static/css/percentage.css.tmpl")
	if err != nil {
		log.Print(err)
	}
	w.Header().Set("Content-Type", "text/css")
	widths := calcPercentageWidths()
	t.Execute(w, widths)
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
