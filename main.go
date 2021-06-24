package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/otofuto/people/pkg/database"
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
	http.HandleFunc("/drop/", DropHandle)

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
		arr := false
		mode := r.URL.Path[len("/human/"):]
		if strings.HasPrefix(mode, "arr/") {
			arr = true
			mode = mode[len("arr/"):]
		}
		hid, err := strconv.Atoi(mode)
		if err != nil {
			http.Error(w, "human id is not set", 400)
			return
		}
		r.ParseMultipartForm(32 << 20)
		if isset(r, []string{"text", "from"}) && !arr {
			if r.FormValue("text") == "" {
				http.Error(w, "false", 400)
				return
			}
			db := database.Connect()
			defer db.Close()

			datas := human.LangSplit(db, r.FormValue("text"))
			/*err = human.SaveWords(db, hid, datas)
			if err != nil {
				log.Println("main.go HumanHandle(w http.ResponseWriter, r *http.Request)")
				log.Println(err)
				http.Error(w, err.Error(), 500)
				return
			}*/
			talked := ""
			isq := false
			for i := len(datas) - 1; i >= 0; i-- {
				datas[i].Human = hid
				tk1, err := human.ResponseString(db, datas[i])
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				if i > 1 {
					if datas[i].Data == datas[i-2].Data+datas[i-1].Data {
						if tk1.Talked == datas[i].Data {
							continue
						}
					}
				}
				if tk1.Talked == datas[i].Data {
					continue
				}
				talked = tk1.Talked + talked
				if tk1.IsQuestion {
					isq = true
				}
			}
			if isq {
				talked += "？"
			}
			tk := human.Talk{
				Opponent:   r.FormValue("from"),
				Human:      hid,
				Heard:      r.FormValue("text"),
				Talked:     talked,
				IsQuestion: isq,
			}
			if err != nil {
				log.Println("main.go HumanHandle(w http.ResponseWriter, r *http.Request)")
				log.Println(err)
				http.Error(w, err.Error(), 500)
				return
			}
			err = tk.Insert(db)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			resp := struct {
				Sds  []human.StringData `json:"sds"`
				Talk human.Talk         `json:"talk"`
			}{
				Sds:  datas,
				Talk: tk,
			}
			bytes, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			fmt.Fprintf(w, string(bytes))
		} else if arr && isset(r, []string{"text[]"}) {
			for k, v := range r.MultipartForm.Value {
				if "text[]" == k {
					db := database.Connect()
					defer db.Close()

					for _, text := range v {
						datas := human.LangSplit(db, text)
						err = human.SaveWords(db, hid, datas)
						if err != nil {
							log.Println("main.go HumanHandle(w http.ResponseWriter, r *http.Request)")
							log.Println(err)
							http.Error(w, err.Error(), 500)
							return
						}
					}
					break
				}
			}
			fmt.Fprintf(w, "true")
		} else {
			http.Error(w, "parameter not enough", 400)
		}
	} else {
		http.Error(w, "method not allowed", 405)
	}
}

func DropHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method == http.MethodGet {
		filename := r.URL.Path[len("/drop/"):]
		if strings.HasSuffix(filename, "/") {
			filename = filename[:len(filename)-1]
		}
		if filename == "" {
			filename = "index"
		}
		temp := template.Must(template.ParseFiles("template/drop/" + filename + ".html"))
		if err := temp.Execute(w, TempContext{}); err != nil {
			log.Println(err)
			http.Error(w, "HTTP 500 Internal server error", 500)
			return
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
