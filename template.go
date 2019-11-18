// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"html/template"
	"io"
	"strings"
)

// buildTemplate builds the template structure and saves it
// in the *Server object
func (s *Server) buildTemplate() {
	const tpl = `<!DOCTYPE html>
<html>
<head>
{{- $base := printf "%s/%s" .Base .Pkg -}}
{{- $host := printf "%s%s" .VcsHost .Pkg -}}
{{- $import := printf "%s %s %s" $base .VcsType $host -}}
{{- $redirect := .Pkg -}}
{{- if ne .Redirect .VcsHost -}}
{{- $redirect = printf "%s%s" .Redirect .Request -}}
{{- else -}}
{{- $redirect = printf "%s%s" .Redirect .Pkg -}}
{{- end  }}
	<meta charset="UTF-8">
	<meta name="go-import" content="{{ $import }}">
{{- if or .Dir .File -}}
{{- $dir := (printf "%s/%s" $host .Dir) -}}{{- if not .Dir -}}{{- $dir = "_" -}}{{- end -}}
{{- $file := (printf "%s/%s" $host .File) -}}{{- if not .File -}}{{- $file = "_" -}}{{- end -}}
{{- $source := (printf "%s %s %s %s" $base $host $dir $file)  }}
	<meta name="go-source" content="{{ $source }}">
{{- end  }}
	<meta http-equiv="refresh" content="0;url={{$redirect}}">
</head>
<body>
<p>Redirecting to <a href="{{$redirect}}">{{$redirect}}</a></p>
</body>
</html>`

	s.template = template.Must(template.New("vanity").Parse(tpl))
}

func (s *Server) serveMeta(w io.Writer, req string) {
	pkg := repoBase(req)

	tplData := struct {
		Base     string
		Pkg      string
		VcsHost  string
		VcsType  string
		Redirect string
		Request  string
		Dir      string
		File     string
	}{
		Base:     s.base,
		Pkg:      pkg,
		VcsHost:  s.repo.root,
		VcsType:  s.repo.vcsType,
		Redirect: s.redirect,
		Request:  strings.TrimPrefix(req, "/"),
		Dir:      s.repo.dirFormat,
		File:     s.repo.fileFormat,
	}

	s.template.Execute(w, tplData)
}
