package main 

import (
	"api"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"os"
	"text/template"
)

type StaticResource struct {
	val      []byte
	dataType string
}

type Page struct {
	Html     string
	JsPath   string
	CssPath  string
	PageName string
}

var Pages map[string]string
var Static map[string]StaticResource

func init() {
	Static = getStaticResources()

	if Static == nil {
		Static = make(map[string]StaticResource)
	}

	Pages = make(map[string]string)

	pagePaths := directoryContentsHierarchy(PAGES_DIRECTORY)

	baseTemplate := getBaseTemplate()

	generatePages(pagePaths, baseTemplate)
}

func getBaseTemplate() *template.Template {
	templatePage, _ := ioutil.ReadFile("pages/" + BASE_TEMPLATE + ".html")
	
	pageTemplate, _ := template.New("pageTemplate").Parse(string(templatePage))

	return pageTemplate
}

func generatePages(paths []DirectoryStructure, parent *template.Template) {
	for _, path := range paths {

		page := renderPage(path)

		buf := new(bytes.Buffer)
		err := parent.Execute(buf, page)

		if err != nil {
			panic(err)
		}

		pageText := buf.String()
		Pages[path.url] = pageText

		fmt.Printf("Rendered Page: %s\n", path.url)

		var newParentTemplate *template.Template
		newParentTemplate, err = template.New("pageTemplate").Parse(pageText)

		if err != nil {
			panic(err)
		}

		generatePages(path.children, newParentTemplate)
	}
}

func filePath(path DirectoryStructure, extension string) string {
	return path.path + path.name + "." + extension
}

func staticFileUrl(path DirectoryStructure, extension string) string {
	return STATIC_HEADER + "/" + path.url + path.name + "." + extension
}

func renderPage(path DirectoryStructure) (page Page) {

	htmlPath := filePath(path, HTML_EXTENSION)
	jsPath := filePath(path, JS_EXTENSION)
	cssPath := filePath(path, CSS_EXTENSION)

	jsUrl := staticFileUrl(path, JS_EXTENSION)
	cssUrl := staticFileUrl(path, CSS_EXTENSION)

	html, err := ioutil.ReadFile(htmlPath)

	if err != nil {
		log.Fatal("File System ", path, " not setup:", err)
	}

	Static[jsUrl] = loadStaticResource(jsPath, JS_CONTEXT, false)
	Static[cssUrl] = loadStaticResource(cssPath, CSS_CONTEXT, false)

	var pageTitle []byte
	pageTitle, _ = ioutil.ReadFile(path.path + "/" + path.name + ".txt")

	return Page{
		string(html),
		jsUrl,
		cssUrl,
		string(pageTitle),
	}
}

func loadStaticResource(filePath string, contextType string, required bool) StaticResource {
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		if required {
			panic(err)
		} else {
			return StaticResource{}
		}
	}

	return StaticResource{content, contextType}
}

func getStaticResources() (map[string]StaticResource) {
	resources := make(map[string]StaticResource)

	scripts := directoryContents(SCRIPTS_DIRECTORY)
	images := directoryContents(IMAGES_DIRECTORY)
	stylesheets := directoryContents(CSS_DIRECTORY)
	pdfs := directoryContents(PDF_DIRECTORY)

	for _, scriptPath := range scripts {
		resources["static/" + scriptPath] = loadStaticResource(scriptPath, JS_CONTEXT, true)
	}

	for _, imagePath := range images {
		resources["static/" + imagePath] = loadStaticResource(imagePath, IMAGE_CONTEXT, true)
	}

	for _, stylesheetPath := range stylesheets {
		resources["static/" + stylesheetPath] = loadStaticResource(stylesheetPath, CSS_CONTEXT, true)
	}

	for _, pdfPath := range pdfs {
		resources["pdf/" + pdfPath] = loadStaticResource(pdfPath, PDF_CONTEXT, true)
	}

	return resources
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request: %s\n", r.URL.Path)

	if r.URL.Path == "/api" {
		api.ContactAPI(w, r)
		return
	}

	showPage(w, r)
}

func showPage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	splitUrl := strings.Split(path, "/")

	if path == "" {
		if html, ok := Pages["index/"]; ok {
			fmt.Fprintf(w, html)
			return
		}
	} else if path == "contact" {
		queryValues := r.URL.Query()

		messageArray, _ := queryValues["message"]

		message := "\n"

		for _, part := range messageArray {
			message += part
		}

		file, _ := os.OpenFile("contact.txt", os.O_APPEND|os.O_WRONLY, 0644)

		defer file.Close()

		fmt.Fprintf(file, message)

		fmt.Fprintf(w, "")
	}

	if splitUrl[0] == "static" {
		if val, ok := Static[path]; ok {
			w.Header().Set("content-type", val.dataType)

			length, _ := w.Write(val.val)
			w.Header().Set("content-length", strconv.Itoa(length))
			return
		}
	}

	if len(splitUrl) > 1 {
		path = strings.Join(splitUrl[:len(splitUrl)-1], "/")
	}

	if path[:len(path) - 1] != "/" {
		path += "/"
	}

	if html, ok := Pages[path]; ok {
		fmt.Fprintf(w, html)
		return
	}

	fmt.Fprintf(w, Pages["404"])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
