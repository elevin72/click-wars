package main

import (
	// "database/sql"

	"fmt"
	"log"

	// "os"

	// "os"
	// "strconv"

	"net/http"

	_ "github.com/mattn/go-sqlite3"

	. "github.com/elevin72/click-wars/internal/routes"
	. "github.com/elevin72/click-wars/internal/server"
)

func main() {

	server := NewServer()
	go server.Run()

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/static/", StaticHandler)
	http.HandleFunc("/static/css/percentage.css", PercentageOnLoadHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(server, w, r)
	})

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
