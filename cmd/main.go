package main

import (
	// "database/sql"

	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/elevin72/click-wars/internal/routes"
	"github.com/elevin72/click-wars/internal/server"
)

func main() {

	server := server.New()
	go server.Run()

	http.HandleFunc("/", routes.RootHandler)
	http.HandleFunc("/static/", routes.StaticHandler)
	http.HandleFunc("/static/css/percentage.css", routes.PercentageOnLoadHandler)
	http.HandleFunc("/ws", routes.MakeWebsocketHandler(server))

	port := "8080"
	fmt.Printf("Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

//TODO:
// - create envs for table creation, or reuse (in prod). Table location in
// prod should take into account that file will be in a docker volume
// - store every click by using long living worker goroutine (pool) which holds db connection
// - provide API (http) for getting all (filters?) clicks
// - docker compose it all together
// - frontend changes to edit animation
// - fix animation browser bug
// - build html page with graphs of clicks at /analytics?query=foo
