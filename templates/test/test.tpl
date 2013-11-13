{{define "test/test"}}
<div class="red_container mediumpadding">
	{{range .D.user}}
		{{.NAME}} - {{.ID}} - {{.EMAIL}}<br>
	{{end}}
</div>
{{end}}
