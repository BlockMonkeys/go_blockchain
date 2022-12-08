package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"pkg/blockchain"
)

const (
	PORT        string = ":5000"
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(w http.ResponseWriter, r *http.Request) {
	data := homeData{"MY HOME", nil}
	templates.ExecuteTemplate(w, "home", data)
}

func add(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "add", nil)
	case "POST":
		blockchain.Blockchain().AddBlock()
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(port int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Println("Listening On http://localhost:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
