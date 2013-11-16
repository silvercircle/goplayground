{{define "test/test"}}
<div class="red_container mediumpadding">
	{{/* range over Context["D"]["user"] and call test/bit template to output each item */}}
	{{/* Context["D"]["user"] can be an array or a map */}}
	{{range .D.user}}
		{{/* The dot after the template name is important, otherwise test/bit would be called with nil data, thus not producing any output */}}
		{{template "test/bit" .}}
	{{end}}
</div>
{{end}}
