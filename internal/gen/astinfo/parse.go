package astinfo

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Package struct {
	Build         *build.Package
	PkgImportPath string
	PkgName       string
	Types         map[string]*TypeDecl
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
		return nil, errors.Errorf("required exactly 1 package in directory, found %v", len(pkgs))
	}

	for _, unresolvedPkg := range pkgs {
		// TODO: provide importer and universe to handle resolving imports
		pkg, _ := ast.NewPackage(fset, unresolvedPkg.Files, nil, nil)
		return ParsePkg(pkg, buildPkg)
	}

	return nil, errors.Errorf("unreachable")
}

func ParsePkg(pkg *ast.Package, buildPkg *build.Package) (*Package, error) {
	res := &Package{
		Build:   buildPkg,
		PkgName: pkg.Name,
		Types:   map[string]*TypeDecl{},
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
	t := &TypeDecl{
		Name:        s.Name.Name,
		Annotations: parseAnnotations(doc),
	}
	var err error

	switch expr := s.Type.(type) {
	case *ast.StructType:
		t.Kind = Struct
		t.Fields, err = parseFields(pkg, expr.Fields)
		if err != nil {
			return err
		}
	case *ast.InterfaceType:
		t.Kind = Interface
		t.Methods, err = parseMethods(pkg, expr.Methods)
		if err != nil {
			return err
		}
	case *ast.Ident:
		t.Kind = Alias
		t.Alias, err = parseTypeRef(pkg, expr)
		if err != nil {
			return err
		}
	default:
		// TODO: parse type ref here also?
		return errors.Errorf("unknown type %v", expr)
	}

	if _, ok := pkg.Types[t.Name]; ok {
		return errors.Errorf("duplicate type %s", t.Name)
	}
	pkg.Types[t.Name] = t

	return nil
}

func parseMethods(pkg *Package, methods *ast.FieldList) ([]Method, error) {
	if methods == nil {
		return nil, nil
	}

	var res []Method

	for _, f := range methods.List {
		if len(f.Names) != 1 {
			return nil, errors.Errorf("invalid name in func %v", f)
		}
		name := f.Names[0].Name

		fun, ok := f.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		args, err := parseFields(pkg, fun.Params)
		if err != nil {
			return nil, err
		}

		results, err := parseFields(pkg, fun.Results)
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

func parseFields(pkg *Package, fields *ast.FieldList) ([]Field, error) {
	if fields == nil {
		return nil, nil
	}

	var res []Field

	for _, f := range fields.List {
		if len(f.Names) > 1 {
			return nil, errors.Errorf("ambigious name in field %v", f)
		}

		name := ""
		if len(f.Names) == 1 {
			name = f.Names[0].Name
		}

		typeRef, err := parseTypeRef(pkg, f.Type)
		if err != nil {
			return nil, err
		}

		field := Field{
			Name:        name,
			Type:        typeRef,
			Annotations: parseAnnotations(f.Doc),
		}
		res = append(res, field)
	}

	return res, nil
}

func parseTypeRef(pkg *Package, expr ast.Expr) (*TypeRef, error) {
	switch t := expr.(type) {
	case *ast.Ident:
		name := t.Name
		if prim := IsPrimitive(name); prim != nil {
			return &TypeRef{
				RefKind:   RefPrimitive,
				Primitive: prim,
				Name:      name,
			}, nil
		}
		return &TypeRef{
			RefKind:     RefRef,
			Name:        name,
			ExternalPkg: "",
		}, nil
	case *ast.SelectorExpr:
		typeSel := t.Sel.Name
		pkgIdent, ok := t.X.(*ast.Ident)
		if !ok {
			// TODO: is there anything else?
			return nil, errors.Errorf("expected typename in field %#v", t)
		}
		return &TypeRef{
			RefKind:     RefRef,
			Name:        typeSel,
			ExternalPkg: pkgIdent.Name, // TODO: resolve full package name
		}, nil
	case *ast.StarExpr:
		nxt, err := parseTypeRef(pkg, t.X)
		if err != nil {
			return nil, err
		}
		if nxt.IsPointer {
			nxt.AdditionalPointers++
		} else {
			nxt.IsPointer = true
		}
		return nxt, nil
	case *ast.FuncType:
		// TODO: handle params and results, or at least save them in Embedded
		return &TypeRef{
			RefKind:  RefEmbedded,
			Embedded: &Embedded{Kind: Func},
			Name:     "func()()",
		}, nil
	case *ast.ArrayType:
		// t.Len, as well as arrays (not slices) are not supported
		nxt, err := parseTypeRef(pkg, t.Elt)
		if err != nil {
			return nil, err
		}
		return &TypeRef{
			RefKind:   RefSlice,
			ValueType: nxt,
			Name:      "[]",
		}, nil
	case *ast.InterfaceType:
		return &TypeRef{
			RefKind:  RefEmbedded,
			Embedded: &Embedded{Kind: Interface},
			Name:     "interface{}",
		}, nil
	case *ast.MapType:
		mapKey, err := parseTypeRef(pkg, t.Key)
		if err != nil {
			return nil, err
		}
		mapValue, err := parseTypeRef(pkg, t.Value)
		if err != nil {
			return nil, err
		}
		return &TypeRef{
			RefKind:   RefMap,
			KeyType:   mapKey,
			ValueType: mapValue,
			Name:      "map[][]",
		}, nil
	default:
		return nil, errors.Errorf("expected typename in field %#v, pkg=%s", t, pkg.PkgName)
	}
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
