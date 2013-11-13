{{define "errors/dberror"}}
<div class="cat_bar2">
		<h3>{{.L.dberror_title}}</h3>
</div>
<div class="blue_container cleantop mediumpadding">
		{{.dberror}}
</div>
{{end}}