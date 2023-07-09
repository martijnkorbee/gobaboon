// Render package renders HTML templates found in the rootpath/views folder of the project.
// Go and Jet templating engines can be used. See https://github.com/CloudyKit/jet for full documentation on Jet.
// Expected filenaming for Go templates <template_name>.page.tmpl, for Jet templates <template_name>.jet.
//
// This package has built in CSRF protection with https://github.com/justinas/nosurf.
package render

import (
	"errors"
	"net/http"
	"strings"
	"text/template"

	"github.com/CloudyKit/jet/v6"
	"github.com/justinas/nosurf"
)

type Renderer struct {
	rootpath       string
	templateEngine string
	jetViews       *jet.Set
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
}

// createRenderer
func New(rootpath, engine string, debug bool) *Renderer {
	return &Renderer{
		rootpath:       rootpath,
		templateEngine: engine,
		jetViews:       makeJetSet(rootpath, debug),
	}
}

// Page renders a page with the configured engine
// call the respective methods to render with a specific engine.
func (r *Renderer) Page(w http.ResponseWriter, req *http.Request, view string, variables interface{}, data interface{}) error {

	// default template data
	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	// populate default data
	td.CSRFToken = nosurf.Token(req)

	switch strings.ToLower(r.templateEngine) {
	case "go":
		return r.GoPage(w, req, view, td)
	case "jet":
		return r.JetPage(w, req, view, variables, td)
	default:
		return errors.New("no rendering engine specified")
	}
}

// GoPage renders a template using the GO templating engine.
func (r *Renderer) GoPage(w http.ResponseWriter, req *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(r.rootpath + "/templates/views/" + view + ".page.tmpl")
	if err != nil {
		return err
	}

	if err = tmpl.Execute(w, &data); err != nil {
		return err
	}

	return nil
}

// JetPage renders a template using the Jet templating engine.
func (r *Renderer) JetPage(w http.ResponseWriter, req *http.Request, view string, variables interface{}, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	t, err := r.jetViews.GetTemplate(view + ".jet")
	if err != nil {
		return err
	}

	if err = t.Execute(w, vars, &data); err != nil {
		return err
	}

	return nil
}

// makeJetSet returns a new jet engine
func makeJetSet(rootpath string, debug bool) *jet.Set {
	if debug {
		return jet.NewSet(
			jet.NewOSFileSystemLoader(rootpath+"/templates/views/"),
			jet.InDevelopmentMode(),
		)
	}

	return jet.NewSet(
		jet.NewOSFileSystemLoader(rootpath + "/templates/views/"),
	)
}
