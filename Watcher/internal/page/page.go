package page

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

// TemplateFolder is the folder where the templates are located
var TemplateFolder = "./templates/"

// Page structure that keeps the data of a page, these are used for rendering
// pages.
type Page struct {
	Template  *template.Template
	Context   map[string]interface{}
	Templates []string
}

// NewPage creates a new page in memory, it tries to parse the parent and body
// html templates and if there is an error it returns no page and an error.
// If it was succesfull the function populates Context with an empty
// map[string]interface{}
func NewPage(pages ...string) *Page {
	p := Page{}
	// Add template folder to pages
	for i := 0; i < len(pages); i++ {
		pages[i] = TemplateFolder + pages[i]
	}

	p.Templates = pages

	p.Context = map[string]interface{}{}

	return &p
}

// Parse parses the regsitered pages and creates a template based upon them.
func (p *Page) Parse() {
	var err error
	p.Template, err = template.ParseFiles(p.Templates...)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			fmt.Println("Can't find file " + err.Error())
		}
		fmt.Println("Failed to parse template " + err.Error())
	}
}

// AddContext adds context to the page, this is passed to the templates when
// rendering the page. It takes an key and value to be used as context.
func (p *Page) AddContext(key string, value interface{}) {
	p.Context[key] = value
}

// Render renders the page to the io.Writer that is passed to it. It includes
// the pages Context to load data into.
// After rendering the page it cleans the context for the next request
func (p *Page) Render(w io.Writer) {
	p.Parse()

	err := p.Template.Execute(w, p.Context)
	if err != nil {
		fmt.Println("PageParser", err.Error())
	}
}

// Copy copies the page for thread safe use
func (p *Page) Copy() *Page {
	// Copy the template
	// newTemplate, err := p.Template.Clone()
	// if err != nil {
	// 	fmt.Println("Failed to copy template", err.Error())
	// 	return nil
	// }

	// Create a new page
	newPage := Page{
		nil,
		map[string]interface{}{},
		p.Templates,
	}

	return &newPage
}
