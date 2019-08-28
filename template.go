// Copyright 2019 Nirenjan Krishnan. All rights reserved.

package vanity

import (
	"io"
	"text/template"
)

// buildTemplate builds the template structure and saves it
// in the *Server object
func (s *Server) buildTemplate() {
	const tpl = `<!DOCTYPE html>
<html>
<head>
{{- $import := printf "%s/%s %s %s/%s" .Base .Pkg .VcsType .VcsHost .Pkg -}}
{{- $redirect := (printf "%s/%s" .Redirect .Pkg) -}}
{{- if ne .Redirect .VcsHost -}}
{{- $redirect = (printf "%s%s" .Redirect .Request) -}}
{{- end -}}
	<meta charset="UTF-8">
	<meta name="go-import" content="{{ $import }}">
	<meta http-equiv="refresh" content="0;url={{$redirect}}">
</head>
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
		Request:  req,
		Dir:      s.repo.dirFormat,
		File:     s.repo.fileFormat,
	}

	s.template.Execute(w, tplData)
}
