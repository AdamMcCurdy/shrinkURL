package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	// Import the Radix.v2 redis package.
	"github.com/mediocregopher/radix.v2/redis"
)

type Page struct {
	Title string
	Body  []byte
}

func dialRedis(port int) {
	p := strconv.Itoa(port)
	conn, err := redis.Dial("tcp", "localhost:"+p)
	if err != nil {
		log.Fatal(err)
	}
	// Importantly, use defer to ensure the connection is always properly
	// closed before exiting the main() function.
	defer conn.Close()
}

func (p *Page) save() error {
	filename := p.Title + ".html"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]	
	if title == "" {
        p, _ := loadPage("test")
        fmt.Fprintf(w, "%s", p.Body)        		
	}
}

func shrinkURL(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
    fmt.Fprintln(w, r.PostFormValue("url"))
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/shrink", shrinkURL)
	http.ListenAndServe(":8080", nil)
}
