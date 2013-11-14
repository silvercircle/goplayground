package lib

func (this *Data) DBError(err error) {
    //cnt, _ := LoadTemplate("errors/dberror")
    this.Templates = append(this.Templates, "errors/dberror")
    this.Context["dberror"] = err.Error()
}
