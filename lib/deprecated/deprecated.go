// this is deprecated code, just kept for reference
package deprecated
// a template descriptor node, holding a loaded (and pre-parsed) template, ready for
// rendering.
type TemplateNode struct {
	name     string
	fullname string
	T        *template.Template
	mtime    time.Time // last mtime of the template source file
	// if the timestamp of the source file changes, template will be re-compiled
}

// all templates are kept in a map. The key is the name of the template (e.g. header or sidebar/content)
// which also defines its location in the template folder structure.
var Templates map[string]*TemplateNode

func RenderTemplate(wr io.Writer, name string, data interface{}) error {
	LoadTemplate(name)
	return Templates[name].T.Execute(wr, data)
}

func (s *TemplateNode) RawLoad(t time.Time) error {
	var err error
	fmt.Println("Loading template", s.name)
	s.T = template.New(s.name)
	s.T.ParseFiles(s.fullname)
	if err != nil {
		s.T.Parse("The template " + s.name + " is invalid or does not exist")
		s.mtime = time.Now()
	}
	s.mtime = t
	return err
}

// LoadTemplate loads and prepares a template from a disk file
// The name must be relative to the templates folder of the app.
// For example: LoadTemplate("header") will actually load the template
// in $APPPATH/templates/header.tpl. Consequently, sidebar/top as a name will
// load $APPPATH/templates/sidebar/top.tpl
// the template is parsed and a pointer to the template object is returned. Using
// the Render(..) method on this object, the template can be used for output.
func LoadTemplate(name string) (*template.Template, error) {
	fullname := SysConf.Homepath + "/templates/" + name + ".tpl"
	ft, err := os.Lstat(fullname)
	if err != nil { // the file does not exist, create a "invalid template"
		if val, ok := Templates[name]; ok { // but only, if one for this name tag does not yet exist
			return val.T, nil
		}
		t := new(TemplateNode)
		t.name = name
		t.T = template.New(name)
		t.fullname = fullname
		t.T.Parse(`<div class="red_container smallpadding">The template ` + name + ` is invalid or does not exist</div>`)
		Templates[name] = t
		return t.T, err
	}

	if val, ok := Templates[name]; ok { // template already parsed, check for modified source file
		if ft.ModTime() != val.mtime { // modified? reparse it
			val.RawLoad(ft.ModTime())
		}
		return val.T, nil // and return it
	}
	// creating a new template
	t := new(TemplateNode)
	t.name = name
	t.fullname = fullname
	t.RawLoad(ft.ModTime())
	Templates[name] = t
	return t.T, nil
}
