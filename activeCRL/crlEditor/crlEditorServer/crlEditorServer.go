// Copyright 2017 Paul C. Brown. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"html/template"
	"io/ioutil"

	//	"log"
	"net/http"
	"regexp"
)

type page struct {
	Title string
	Body  []byte
}

var root = "../src/github.com/pbrown12303/activeCRL/activeCRL/"

var templates = template.Must(template.ParseFiles(root+"crlEditor/tmpl/index.html", root+"crlEditor/tmpl/displayTrace.html", root+"crlEditor/tmpl/sandbox.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := &page{Title: "CRL Editor"}
	renderTemplate(w, "index", p)
}

func displayTraceHandler(w http.ResponseWriter, r *http.Request) {
	p := &page{Title: "Display Trace"}
	renderTemplate(w, "displayTrace", p)
}

func sandboxHandler(w http.ResponseWriter, r *http.Request) {
	p := &page{Title: "Sandbox"}
	renderTemplate(w, "sandbox", p)
}

func loadPage(title string) (*page, error) {
	filename := root + "crlEditor/data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &page{Title: title, Body: body}, nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here we will extract the page title from the Request,
		// and call the provided handler 'fn'
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/index/", indexHandler)
	http.HandleFunc("/displayTrace/", displayTraceHandler)
	http.HandleFunc("/sandbox/", sandboxHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(root+"crlEditor/js"))))
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir(root+"crlEditor/icons"))))
	http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *page) save() error {
	filename := root + "crlEditor/data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}
