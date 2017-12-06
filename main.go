package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

type Page struct {
	Title string
	Body  []byte
}

var client redis.Client
var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	return string(b)
}

// SaveNewValue returns for success
func saveNewValue(key string, value string) {
	err := client.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

// GetValue finds the shortURL string in redis and returns the long url
func getValue(shortURL string) string {
	val, err := client.Get(shortURL).Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Returning value " + val)
	return val
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
	fmt.Println(title)

	if title == "" {
		p, _ := loadPage("index")
		fmt.Fprintf(w, "%s", p.Body)
	} else if title != "favicon.ico" {
		longURL := getValue(title)
		http.Redirect(w, r, longURL, 301)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		randString := randSeq(8)
		go saveNewValue(randString, r.Form["url"][0])
		htmlString := "<h1>Your ShortURL</h1>" + "<a href=\"" + "http://localhost:8080/" + randString + "\">" + "http://localhost:8080/" + randString + "</a>"
		fmt.Fprintf(w, "%s", htmlString)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func shrinkURL(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.PostFormValue("url"))
}

func main() {
	client = *redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/shrink", postHandler)
	http.ListenAndServe(":8080", mux)
}
