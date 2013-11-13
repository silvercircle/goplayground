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

func Init() {
    var files []os.FileInfo
    files, _ = ioutil.ReadDir(lib.SysConf.Homepath + "/lang/")
    for _, file := range files {
        if !file.IsDir() && strings.ToLower(filepath.Ext(file.Name())) == ".json" {
            id := file.Name()[0 : len(file.Name())-len(filepath.Ext(file.Name()))]
            Languages[id] = ReadLangFile(id)
        }
    }
}

func ReadLangFile(id string) map[string]string {
    file, err := ioutil.ReadFile(lib.SysConf.Homepath + "/lang/" + id + ".json")
    var tmap = make(map[string]string)
    tmap["lang_id"] = id
    if err == nil {
        var objmap map[string]json.RawMessage
        var s string

        json.Unmarshal(file, &objmap)
        for k, _ := range objmap {
            json.Unmarshal(objmap[k], &s)
            tmap[k] = s
        }
    }
    return tmap
}
