package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/otofuto/people/pkg/human"
)

var port string

type TempContext struct {
	Message string `json:"message"`
}

func main() {
	_ = godotenv.Load()
	port = os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.Handle("/st/", http.StripPrefix("/st/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", IndexHandle)
	http.HandleFunc("/human/", HumanHandle)

	log.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func IndexHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("template/index.html"))
		if err := temp.Execute(w, TempContext{}); err != nil {
			log.Println(err)
			http.Error(w, "HTTP 500 Internal server error", 500)
			return
		}
	} else {
		http.Error(w, "method not allowed", 405)
	}
}

func HumanHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(32 << 20)
		if isset(r, []string{"text"}) {
			if r.FormValue("text") == "" {
				http.Error(w, "false", 400)
				return
			}
			datas := human.LangSplit(r.FormValue("text"))
			bytes, err := json.Marshal(datas)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			fmt.Fprintf(w, string(bytes))
		} else {
			http.Error(w, "parameter not enough", 400)
		}
	} else {
		http.Error(w, "method not allowed", 405)
	}
}

//GETでは使えない
func isset(r *http.Request, keys []string) bool {
	for _, v := range keys {
		exist := false
		for k, _ := range r.MultipartForm.Value {
			if v == k {
				exist = true
			}
		}
		if !exist {
			return false
		}
	}
	return true
}
