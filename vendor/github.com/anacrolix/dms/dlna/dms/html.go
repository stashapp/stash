package dms

import (
	"html/template"
)

var rootTmpl *template.Template

func init() {
	rootTmpl = template.Must(template.New("root").Parse(
		`<form method="post">
			Path: <input type="text"
				name="path"
				{{if .Readonly}} readonly="readonly"{{end}}
				value="{{.Path}}"
			/>
			<input type="submit" value="Update"{{if .Readonly}} disabled="disabled"{{end}}/>
		</form>`))
}
