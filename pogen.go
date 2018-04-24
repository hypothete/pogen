package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var allTemplates *template.Template

// Page is our representation of saved data
type Page struct {
	Title string
	Body  []byte
}

// TemplatePage just converts the Page body to a string
type TemplatePage struct {
	Title string
	Body  string
}

// IndexTemplateData stores the template data for the front page
type IndexTemplateData struct {
	Pages []string
}

func (p *Page) save() error {
	filename := "pages/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPageData(title string) (*Page, error) {
	filename := "pages/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func (p *Page) toTemplateData() *TemplatePage {
	return &TemplatePage{Title: p.Title, Body: string(p.Body)}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func indexPages() ([]string, error) {
	osFileList, err := ioutil.ReadDir("pages")
	if err != nil {
		return nil, err
	}
	var fileList []string
	for _, f := range osFileList {
		fileName := f.Name()
		fileName = strings.Split(fileName, ".")[0]
		fileList = append(fileList, fileName)
	}
	return fileList, nil
}

func loadTemplates() {
	var allFilenames []string
	files, err := ioutil.ReadDir("./templates")
	check(err)
	for _, file := range files {
		filename := file.Name()
		allFilenames = append(allFilenames, "./templates/"+filename)
	}
	allTemplates, err = template.ParseFiles(allFilenames...)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemp := allTemplates.Lookup("index.html")
	pageIndex, err := indexPages()
	check(err)
	indexData := IndexTemplateData{Pages: pageIndex}
	indexTemp.ExecuteTemplate(w, "index", indexData)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	pathElems := strings.Split(r.URL.Path, "/")
	subpath := pathElems[len(pathElems)-1]
	pageData, err := loadPageData(subpath)
	check(err)
	pageTemp := allTemplates.Lookup("page.html")
	pageTempData := pageData.toTemplateData()
	pageTemp.ExecuteTemplate(w, "page", pageTempData)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bodyData := []byte(r.Form["body"][0])
	titleData := r.Form["title"]
	pageToSave := Page{Body: bodyData, Title: titleData[0]}
	pageToSave.save()
	http.Redirect(w, r, "index.html", http.StatusSeeOther)
}

func serve404(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func main() {
	static := http.Dir("static")

	loadTemplates()

	http.HandleFunc("/index.html", indexHandler)

	http.HandleFunc("/pages/", pageHandler)

	http.HandleFunc("/save", saveHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(static)))

	http.HandleFunc("/favicon.ico", serve404)

	fmt.Println("Starting server on port 3333")

	log.Fatal(http.ListenAndServe(":3333", nil))
}
