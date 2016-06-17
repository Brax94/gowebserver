package main

import (
    "io"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "html/template"
    "io/ioutil"
)

func hello(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Hello world!")
}

func to_roman(n int)  string {
    var s string
	while (n > 0 && n < 4)
	{
		s = s + "I"
		n = n - 1
	}
	return s
}

type romanGenerator int
func (n romanGenerator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ascii_num := r.URL.Path[7:]
    i, err := strconv.Atoi(ascii_num)
    if err != nil {
        log.Print(err)
    }
    fmt.Fprintf(w, "Here's your number: %s\n", to_roman(i))
}

type Page struct {
    Title string
    Body  []byte
}
func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func inputHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/input/"):]
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    t, _ := template.ParseFiles("ui.html")
    t.Execute(w, p)
}

func convertHandler(w http.ResponseWriter, r *http.Request){
    number := r.FormValue("quantity")
    http.Redirect(w, r, "/roman/" + number, http.StatusFound)
}

func main() {
    http.Handle("/roman/", romanGenerator(1))
    http.HandleFunc("/", hello)
    http.HandleFunc("/input/", inputHandler)
    http.HandleFunc("/convert/", convertHandler)

    err := http.ListenAndServe(":8000", nil)
    log.Fatal(err)
}
