package lib

func (this *Data) DBError(err error) {
    cnt, _ := LoadTemplate("errors/dberror")
    this.Templates = append(this.Templates, cnt)
    this.Context["dberror"] = err.Error()
}
