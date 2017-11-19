package app

//var UsageTmpl = `
//{{- if has_intro -}}
//	{{ intro }}
//{{ end -}}
//Usage:
//	{{ exec_bin }} command [arguments]
//
//{{ if has_runnable . -}}
//The commands are:
//{{range .}}{{if runnable . }}
//	{{.GetName | printf "%-11s"}} {{.GetShort}}{{end}}{{end}}
//
//Use "{{exec_bin}} help [command]" for more information about a command.
//{{ end -}}
//
//{{ if has_not_runnable . }}
//Additional help topics:
//{{range .}}{{if not (runnable .) }}
//	{{.GetName | printf "%-11s"}} {{.GetShort}}{{end}}{{end}}
//
//Use "{{exec_bin}} help [topic]" for more information about that topic.
//{{ end -}}
//`

var UsageTmpl = `
{{- if has_intro -}}
    {{ intro }}
{{ end -}}
Usage:
    {{ exec_bin }} command [arguments]

{{ if .HasCommands true -}}
Available commands:
{{ range $group := .GetGroups true -}}

{{- if $.NotDefaultGroup $group }}

  {{ $group }}
{{- end -}}

{{- range $cmd := $.GetCommands true $group }}
    {{ $cmd.GetName}}{{ space $cmd.GetName $.Max }}{{.GetShort}}
{{- end }}

{{- end }}

Use "{{exec_bin}} help [command]" for more information about a command.
{{- end }}
{{ if .HasCommands false }}
Help topics:
{{ range $group := .GetGroups false -}}

{{- if $.NotDefaultGroup  $group }}{{ $group }}{{ end -}}

{{- range $cmd := $.GetCommands false $group }}
    {{ $cmd.GetName}}{{ space $cmd.GetName $.Max }}{{.GetShort}}
{{- end }}

{{- end }}

Use "{{exec_bin}} help [topic]" for more information about that topic.
{{ end }}`
