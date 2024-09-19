package main

import (
	// "database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"

	// "os"

	// "os"
	// "strconv"

	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	. "github.com/elevin72/click-wars/internal"
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

	linePosition := LinePosition.Load()
	totalHits := TotalHits.Load()
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
	percentage := (float32(LinePosition.Load()) * float32(scalingConstantFloat)) + 50
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

	//TODO:
	// - create envs for table creation, or reuse (in prod). Table location in
	// prod should take into account that file will be in a docker volume
	// - store every click by using long living worker goroutine (pool) which holds db connection
	// - provide API (http) for getting all (filters?) clicks
	// - docker compose it all together
	// - frontend changes to edit animation
	// - fix animation browser bug
	// - build html page with graphs of clicks at /analytics?query=foo

	// os.Remove("./foo.db")

	// db, err := sql.Open("sqlite3", "./foo.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// sqlStmt := `
	// create table foo (id integer not null primary key, name text);
	// delete from foo;
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return
	// }

	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// for i := 0; i < 100; i++ {
	// 	_, err = stmt.Exec(i, fmt.Sprintf("hello world %d", i))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// err = tx.Commit()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rows, err := db.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err = rows.Scan(&id, &name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	server := NewServer()
	go server.Run()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/static/css/percentage.css", percentageOnLoadHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(server, w, r)
	})

	port := "8080"
	fmt.Printf("Server started on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
