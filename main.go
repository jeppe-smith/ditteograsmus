package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		files, err := ioutil.ReadDir("public")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		paths := []string{}

		for _, file := range files {
			paths = append(paths, file.Name())
		}

		if err := tmpl.Execute(w, paths); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("upload.html")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := tmpl.Execute(w, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		if r.Method == http.MethodPost {
			r.ParseMultipartForm(10 << 20)

			file, _, err := r.FormFile("file")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			tempFile, err := ioutil.TempFile("public", "*.webp")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer tempFile.Close()

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = tempFile.Write(fileBytes)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "http://localhost:1337", http.StatusFound)
		}
	})

	public := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", public))

	log.Fatal(http.ListenAndServe(":1337", nil))
	fmt.Print("server running")
}
