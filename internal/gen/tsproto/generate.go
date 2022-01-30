package tsproto

import (
	"bytes"
	_ "embed"
	"github.com/rwlist/gjrpc/internal/gen/protog"
	"text/template"
)

func GenerateSource(proto *protog.Protocol) (string, error) {
	w := bytes.NewBuffer(nil)
	_, _ = w.WriteString(tmplHeader)

	for _, s := range proto.Services {
		_ = w.WriteByte('\n')
		err := serviceTemplate.Execute(w, s)
		if err != nil {
			return "", err
		}
	}

	for _, m := range proto.Models {
		_ = w.WriteByte('\n')
		err := modelTemplate.Execute(w, m)
		if err != nil {
			return "", err
		}
	}

	return w.String(), nil
}

//go:embed header.tstmpl
var tmplHeader string

//go:embed service.tstmpl
var tmplService string
var serviceTemplate = template.Must(template.New("service").Parse(tmplService))

//go:embed model.tstmpl
var tmplModel string
var modelTemplate = template.Must(template.New("model").Funcs(funcMap).Parse(tmplModel))
