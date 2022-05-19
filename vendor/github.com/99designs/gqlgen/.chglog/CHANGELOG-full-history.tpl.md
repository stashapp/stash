{{- $repourl := $.Info.RepositoryURL -}}
# CHANGELOG
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<a name="unreleased"></a>
## [Unreleased]({{ .Info.RepositoryURL }}/compare/{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}...HEAD)

{{ if .Unreleased.NoteGroups }}
{{ range .Unreleased.NoteGroups -}}
### {{ .Title }}
{{ range .Notes -}}
{{ .Body }}
{{ end -}}  <!-- end of Notes -->
{{ end -}} <!-- end of NoteGroups -->
{{ end -}}  <!-- end of if -->
{{ range .Unreleased.CommitGroups }}
{{ range .Commits -}}

{{- /** Remove markdown urls when there's a pull request linked and replace it with a tag **/ -}}
{{- $subject := (regexReplaceAll `URL` (regexReplaceAll `\[(.*)(\d\d)\]\(.*?\)` .Subject "<a href=\"URL/pull/${2}\">${1}${2}</a>") $repourl) -}}
{{- /** Filter out refs mentioned in the title **/ -}}
{{- $list := (list) -}}
{{- range $idx, $ref := .Refs -}}
{{- if not (regexMatch $ref.Ref $subject) -}}
{{ $list = append $list $ref }}
{{- end -}}
{{- end -}}
{{- /** end custom variables **/ -}}

{{ if .TrimmedBody -}}<dl><dd><details><summary>{{ else -}}- {{ end -}}
<a href="{{$repourl}}/commit/{{.Hash.Long}}"><tt>{{.Hash.Short}}</tt></a> {{ $subject }}
{{- if $list -}}
{{ printf " %s " "(closes"}}
{{- range $idx, $ref := $list -}}{{ if $idx }}, {{ end -}}
<a href="{{ $repourl }}/issues/{{ $ref.Ref}}"> #{{ $ref.Ref}}</a>{{ end }})
{{- end -}}
{{ if .TrimmedBody -}}</summary>{{ printf "\n\n%s\n\n" .TrimmedBody }}</details></dd></dl>{{ end }}

{{ end }} <!-- end of Commits -->
{{ end -}} <!-- end of CommitGroups -->

{{- if .Versions }}
{{ range .Versions -}}
<a name="{{ .Tag.Name }}"></a>
## {{ if .Tag.Previous }}[{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}[{{ .Tag.Name }}](https://github.com/99designs/gqlgen/releases/tag/{{ .Tag.Name }}){{ end }} - {{ datetime "2006-01-02" .Tag.Date }}
{{- if .CommitGroups -}}
{{ range .CommitGroups -}}

### {{ .Title }}
{{ range .Commits -}}
{{- /** Remove markdown urls when there's a pull request linked and replace it with a tag **/ -}}
{{- $subject := (regexReplaceAll `URL` (regexReplaceAll `\[(.*)(\d\d)\]\(.*?\)` .Subject "<a href=\"URL/pull/${2}\">${1}${2}</a>") $repourl) -}}
{{- /** Filter out refs mentioned in the title **/ -}}
{{- $list := (list) -}}
{{- range $idx, $ref := .Refs -}}
{{- if not (regexMatch $ref.Ref $subject) -}}
{{ $list = append $list $ref }}
{{- end -}}
{{- end -}}
{{- /** end custom varaibles **/ -}}

{{ if .TrimmedBody -}}<dl><dd><details><summary>{{ else -}}- {{ end -}}
<a href="{{$repourl}}/commit/{{.Hash.Long}}"><tt>{{.Hash.Short}}</tt></a> {{ $subject }}
{{- if $list -}}
{{ printf " %s " "(closes"}}
{{- range $idx, $ref := $list -}}{{ if $idx }}, {{ end -}}
<a href="{{ $repourl }}/issues/{{ $ref.Ref}}"> #{{ $ref.Ref}}</a>{{ end }})
{{- end -}}

- {{ if .Type }}**{{ .Type }}:** {{ end }}{{ if .Subject }}{{ .Subject }}{{ else }}{{ .Header }}{{ end }}
{{ end }} <!-- end of Commits -->
{{ end -}} <!-- end of CommitGroups -->
{{ else }}
{{ range .Commits -}}

{{- /** Remove markdown urls when there's a pull request linked and replace it with a tag **/ -}}
{{- $subject := (regexReplaceAll `URL` (regexReplaceAll `\[(.*)(\d\d)\]\(.*?\)` .Subject "<a href=\"URL/pull/${2}\">${1}${2}</a>") $repourl) -}}
{{- /** Filter out refs mentioned in the title **/ -}}
{{- $list := (list) -}}
{{- range $idx, $ref := .Refs -}}
{{- if not (regexMatch $ref.Ref $subject) -}}
{{ $list = append $list $ref }}
{{- end -}}
{{- end -}}
{{- /** end custom variables **/ -}}

{{ if .TrimmedBody -}}<dl><dd><details><summary>{{ else -}}- {{ end -}}
<a href="{{$repourl}}/commit/{{.Hash.Long}}"><tt>{{.Hash.Short}}</tt></a> {{ $subject }}
{{- if $list -}}
{{ printf " %s " "(closes"}}
{{- range $idx, $ref := $list -}}{{ if $idx }}, {{ end -}}
<a href="{{ $repourl }}/issues/{{ $ref.Ref}}"> #{{ $ref.Ref}}</a>{{ end }})
{{- end -}}
{{ if .TrimmedBody -}}</summary>{{ printf "\n\n%s\n\n" .TrimmedBody }}</details></dd></dl>{{ end }}

{{ end }} <!-- end of Commits -->
{{ end -}} <!-- end of Else -->

{{ if .NoteGroups }}
{{ range .NoteGroups -}}
### {{ .Title }}
{{ range .Notes -}}
{{ .Body }}
{{ end -}} <!-- end of Notes -->
{{ end -}} <!-- end of NoteGroups -->
{{ end -}} <!-- end of If NoteGroups -->
{{ end -}} <!-- end of Versions -->
{{ end -}} <!-- end of If Versions -->
