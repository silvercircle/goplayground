{{define "content"}}
	The content is: {{.content}}<br>
	The link is: <a href="{{.C.testurl}}">{{.L.foo}}</a><br>
	The other link is: {{.C.testlink}}
{{end}}
