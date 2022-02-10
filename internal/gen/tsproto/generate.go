package tsproto

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/rwlist/gjrpc/internal/gen/astinfo"
	"text/template"

	"github.com/rwlist/gjrpc/internal/gen/protog"
)

func GenerateSource(proto *protog.Protocol) (string, error) {
	w := bytes.NewBuffer(nil)
	_, _ = w.WriteString(tmplHeader)

	_ = w.WriteByte('\n')
	err := apiTemplate.Execute(w, proto)
	if err != nil {
		return "", err
	}

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

	for _, t := range proto.Package.Types {
		if t.Kind != astinfo.Alias {
			continue
		}

		_ = w.WriteByte('\n')
		_, _ = fmt.Fprintf(w, "export type %s = %s", t.Name, convertGoType(t.Alias))
	}

	return w.String(), nil
}

//go:embed header.tstmpl
var tmplHeader string

//go:embed api.tstmpl
var tmplAPI string
var apiTemplate = template.Must(template.New("api").Parse(tmplAPI))

//go:embed service.tstmpl
var tmplService string
var serviceTemplate = template.Must(template.New("service").Parse(tmplService))

//go:embed model.tstmpl
var tmplModel string
var modelTemplate = template.Must(template.New("model").Funcs(funcMap).Parse(tmplModel))
