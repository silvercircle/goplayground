// template system
package lib

import (
	"fmt"
	"text/template"
	"io/ioutil"
	"os"
	"sync"
)

var T *template.Template = nil

var templates_to_load []string

// recursively scan the templates directory and preload all templates found
// called by main on startup once.
// if checkonly == true, the function will only collect the timestamps
// for each template directory 
func DoPreloadTemplates(prefix string, checkonly bool) {
	var files []os.FileInfo

	curdir := SysConf.Homepath + "/templates/" + prefix
	files, _ = ioutil.ReadDir(curdir)

	ft, _ := os.Lstat(curdir)
	timestamp := ft.ModTime().Unix()
	if timestamp > SysConf.TemplateMostRecent {
		SysConf.TemplateMostRecent = timestamp
	}
	var finalPrefix string
	if prefix == "" {
		finalPrefix = ""
	} else {
		finalPrefix = prefix + "/"
	}
	for _, file := range files {
		if file.IsDir() {
			DoPreloadTemplates(finalPrefix + file.Name(), checkonly)
		} else if !checkonly {
			templates_to_load = append(templates_to_load, SysConf.Homepath+"/templates/"+finalPrefix+file.Name())
		}
	}
}

var TemplateLock sync.RWMutex

func PreloadTemplates(prefix string, force bool) {
	TemplateLock.Lock()
	if len(templates_to_load) > 0 {
		templates_to_load = nil
	}
	DoPreloadTemplates(prefix, false)
	fmt.Println("Rescanning templates...")
	TNew := template.Must(template.ParseFiles(templates_to_load...))
	T = TNew
	TNew = nil
	TemplateLock.Unlock()
}
