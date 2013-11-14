// template system
package lib

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

var T *template.Template = nil

var templates_to_load []string

// recursively scan the templates directory and preload all templates found
// called by main on startup once.
func DoPreloadTemplates(prefix string) {
	var files []os.FileInfo

	curdir := SysConf.Homepath + "/templates/" + prefix
	files, _ = ioutil.ReadDir(curdir)

	ft, _ := os.Lstat(curdir)
	if ft.ModTime().Unix() > SysConf.TemplateMostRecent {
		SysConf.TemplateMostRecent = ft.ModTime().Unix()
	}
	var finalPrefix string
	if prefix == "" {
		finalPrefix = ""
	} else {
		finalPrefix = prefix + "/"
	}
	for _, file := range files {
		if file.IsDir() {
			DoPreloadTemplates(finalPrefix + file.Name())
		} else {
			//ext := filepath.Ext(file.Name())
			//basename := file.Name()[0 : len(file.Name())-len(ext)]
			//LoadTemplate(finalPrefix + basename)
			templates_to_load = append(templates_to_load, SysConf.Homepath+"/templates/"+finalPrefix+file.Name())
		}
	}
}

func PreloadTemplates(prefix string, force bool) {
	if len(templates_to_load) > 0 {
		templates_to_load = nil
	}
	DoPreloadTemplates(prefix)
	fmt.Println("Rescanning templates...")
	if T != nil {
		T = nil
	}
	T = template.Must(template.ParseFiles(templates_to_load...))
}
