package app

var UsageTmpl = `
{{- if has_intro -}}
	{{ intro }}
{{ end -}}
Usage:
	{{ exec_bin }} command [arguments]

{{ if has_runnable . -}}
The commands are:
{{range .}}{{if runnable . }}
	{{.GetName | printf "%-11s"}} {{.GetShort}}{{end}}{{end}}

Use "{{exec_bin}} help [command]" for more information about a command.
{{ end -}}

{{ if has_not_runnable . }}
Additional help topics:
{{range .}}{{if not (runnable .) }}
	{{.GetName | printf "%-11s"}} {{.GetShort}}{{end}}{{end}}

Use "{{exec_bin}} help [topic]" for more information about that topic.
{{ end -}}
`
