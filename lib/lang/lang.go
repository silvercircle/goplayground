package lang

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "testgo/lib"
)

var Languages = map[string]map[string]string{}

// scan the "lang" directory for files in .json format. Each file is considered
// a valid translation
func Init() {
	// first read the default.json (english and reference file)
	Languages["default"] = make(map[string]string)
	ReadLangFile("default", Languages["default"])
    var files []os.FileInfo
    files, _ = ioutil.ReadDir(lib.SysConf.Homepath + "/lang/")
    for _, file := range files {
        if !file.IsDir() && strings.ToLower(filepath.Ext(file.Name())) == ".json" {
            id := file.Name()[0 : len(file.Name())-len(filepath.Ext(file.Name()))]
            if strings.ToLower(id) == "default" {
            	continue							// already done
            }
            // create a copy of the map that holds the default language strings and override them
            // with translations. This ensures there won't be empty strings, just possibly untranslated ones
            Languages[id] = make(map[string]string)
            for k, v := range Languages["default"] {
            	Languages[id][k] = v
            }
            ReadLangFile(id, Languages[id])
        }
    }
}

// read a language file into target
// id is the identifier and, at the same time, the base name of the file holding
// the translation without the .json extension.
func ReadLangFile(id string, target map[string]string) {
    file, err := ioutil.ReadFile(lib.SysConf.Homepath + "/lang/" + id + ".json")
    if err == nil {
        var objmap map[string]json.RawMessage
        var s string

        json.Unmarshal(file, &objmap)
        for k, _ := range objmap {
            json.Unmarshal(objmap[k], &s)
            target[k] = s
        }
    }
}
