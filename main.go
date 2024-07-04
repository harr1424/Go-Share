package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("uploads")
	if err != nil {
		http.Error(w, "Error finding available files", http.StatusInternalServerError)
		return
	}

	data := struct {
		Files []os.DirEntry
	}{
		Files: files,
	}

	templates.ExecuteTemplate(w, "index.html", data)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		localFile, err := os.Create(filepath.Join("uploads", handler.Filename))
		if err != nil {
			http.Error(w, "Unable to create a copy of the file on the server", http.StatusInternalServerError)
			return
		}
		defer localFile.Close()

		if _, err := io.Copy(localFile, file); err != nil {
			http.Error(w, "Unable to copy file to server", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	encodedFilename := r.URL.Path[len("/download/"):]

	filename, err := url.QueryUnescape(encodedFilename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filepath := filepath.Join("uploads", filename)
	http.ServeFile(w, r, filepath)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new", uploadHandler)
	http.HandleFunc("/download/", fileHandler)

	fmt.Println("Starting server at https://localhost:8443")
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
