package astinfo

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type Package struct {
	Build         *build.Package
	PkgImportPath string
	PkgName       string
	Types         map[string]*Type
}

func ParseDir(path string) (*Package, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	buildPkg, err := build.ImportDir(path, build.ImportComment)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("required exactly 1 package in directory, found %v", len(pkgs))
	}

	for _, pkg := range pkgs {
		return ParsePkg(pkg, buildPkg)
	}

	return nil, fmt.Errorf("unreachable")
}

func ParsePkg(pkg *ast.Package, buildPkg *build.Package) (*Package, error) {
	res := &Package{
		Build:   buildPkg,
		PkgName: pkg.Name,
		Types:   map[string]*Type{},
	}

	importPath, err := findPackagePath(buildPkg.Dir)
	if err != nil {
		return nil, err
	}
	res.PkgImportPath = importPath

	for _, file := range pkg.Files {
		err := parseFile(res, file)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func parseFile(pkg *Package, src *ast.File) error {
	for _, decl := range src.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok { //nolint:exhaustive
			case token.TYPE:
				for _, spec := range d.Specs {
					if s, ok := spec.(*ast.TypeSpec); ok {
						err := parseType(pkg, s, d.Doc)
						if err != nil {
							return err
						}
					}
				}
			default:
				// other tokens are not interesting
			}
		}
	}

	return nil
}

func parseType(pkg *Package, s *ast.TypeSpec, doc *ast.CommentGroup) error {
	t := &Type{
		Name:        s.Name.Name,
		Annotations: parseAnnotations(doc),
	}
	var err error

	switch expr := s.Type.(type) {
	case *ast.StructType:
		t.Kind = Struct
		t.Fields, err = parseFields(expr.Fields)
		if err != nil {
			return err
		}
	case *ast.InterfaceType:
		t.Kind = Interface
		t.Methods, err = parseMethods(expr.Methods)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown type %v", expr)
	}

	if _, ok := pkg.Types[t.Name]; ok {
		return fmt.Errorf("duplicate type %s", t.Name)
	}
	pkg.Types[t.Name] = t

	return nil
}

func parseMethods(methods *ast.FieldList) ([]Method, error) {
	if methods == nil {
		return nil, nil
	}

	var res []Method

	for _, f := range methods.List {
		if len(f.Names) != 1 {
			return nil, fmt.Errorf("invalid name in func %v", f)
		}
		name := f.Names[0].Name

		fun, ok := f.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		args, err := parseFields(fun.Params)
		if err != nil {
			return nil, err
		}

		results, err := parseFields(fun.Results)
		if err != nil {
			return nil, err
		}

		method := Method{
			Name:        name,
			Params:      args,
			Results:     results,
			Annotations: parseAnnotations(f.Doc),
		}
		res = append(res, method)
	}

	return res, nil
}

func parseFields(fields *ast.FieldList) ([]Field, error) {
	if fields == nil {
		return nil, nil
	}

	var res []Field

	for _, f := range fields.List {
		if len(f.Names) > 1 {
			return nil, fmt.Errorf("ambigious name in field %v", f)
		}

		name := ""
		if len(f.Names) == 1 {
			name = f.Names[0].Name
		}

		var typeName string

		switch t := f.Type.(type) {
		case *ast.Ident:
			typeName = t.Name
		case *ast.SelectorExpr:
			typeSel := t.Sel.Name
			pkgIdent, ok := t.X.(*ast.Ident)
			if !ok {
				return nil, fmt.Errorf("expected typename in field %s", name)
			}
			typeName = pkgIdent.Name + "." + typeSel
		default:
			return nil, fmt.Errorf("expected typename in field %s", name)
		}

		field := Field{
			Name:        name,
			Type:        typeName,
			Annotations: parseAnnotations(f.Doc),
		}
		res = append(res, field)
	}

	return res, nil
}

func parseAnnotations(doc *ast.CommentGroup) []Annotation {
	if doc == nil {
		return nil
	}

	var annos []Annotation

	for _, comment := range doc.List {
		const slashes = "//"
		const prefix = slashes + annotationPrefix
		if !strings.HasPrefix(comment.Text, prefix) {
			continue
		}

		raw := strings.TrimPrefix(comment.Text, slashes)
		tokens := strings.Split(raw, " ")

		if len(tokens) < 1 {
			// possibly invalid annotation found
			continue
		}

		annos = append(annos, Annotation{
			Key:    tokens[0],
			Values: tokens[1:],
		})
	}

	return annos
}
