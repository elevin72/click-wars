package routes

import (
	// "database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"

	"io/fs"
	"net/http"
	"strings"

	. "github.com/elevin72/click-wars/internal/server"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed frontend/static
var frontend embed.FS

func RootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(frontend, "frontend/static/index.html")
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

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(path, "css") {
		w.Header().Set("Content-Type", "text/css")
	} else {
		fmt.Println("something else")
	}
	static, err := fs.Sub(frontend, "frontend")
	if err != nil {
		fmt.Print(err)
	}
	http.FileServer(http.FS(static)).ServeHTTP(w, r)
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

func PercentageOnLoadHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(frontend, "frontend/static/css/percentage.css.tmpl")
	if err != nil {
		log.Print(err)
	}
	w.Header().Set("Content-Type", "text/css")
	widths := calcPercentageWidths()
	t.Execute(w, widths)
}

func MakeWebsocketHandler(server *Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWs(server, w, r)
	}
}
