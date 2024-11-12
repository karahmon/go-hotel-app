package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/karahmon/hotel-bookings/pkg/config"
	"github.com/karahmon/hotel-bookings/pkg/models"
)

var app *config.Appconfig

func NewTemplates(a *config.Appconfig) {
	app = a
}
func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// RenderTemplate renders templates using the template cache
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// Create a Template Cache

	// Get Requested Template from Cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Couldnt get template from template cache")
	}

	// Render the Template
	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to response:", err)
	}
}

// CreateTemplateCache creates a map of cached templates
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)

	// Get all files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// Range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

// func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
// 	parseTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
// 	err := parseTemplate.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println("Error Parsing Template: ", err)
// 		return
// 	}
// }

// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	// check if we already have template in our cache
// 	_, inMap := tc[t]
// 	if !inMap {
// 		//need to create the template
// 		fmt.Println("Creating Template and Adding to Cache")
// 		err = CreateTemplateCache(t)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	} else {
// 		fmt.Println("Using Cached Template")
// 	}

// 	tmpl = tc[t]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func CreateTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.tmpl",
// 	}
// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}
// 	tc[t] = tmpl
// 	return nil
// }
