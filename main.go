package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	uploadsPath := filepath.Join(".", "uploads")
	err := os.MkdirAll(uploadsPath, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		files, err := ioutil.ReadDir("uploads")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sort.Slice(files, func(a, b int) bool {
			return files[a].ModTime().After(files[b].ModTime())
		})

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

			tempFile, err := ioutil.TempFile("uploads", "*.webp")
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

			http.Redirect(w, r, "/", http.StatusFound)
		}
	})

	uploads := http.FileServer(http.Dir("./uploads"))
	public := http.FileServer(http.Dir("./public"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", uploads))
	http.Handle("/public/", http.StripPrefix("/public/", public))

	log.Fatal(http.ListenAndServe(":1337", nil))
	fmt.Print("server running")
}
