package main

import (
   "bytes"
   "fmt"
   "go/parser"
   "go/token"
   "io/ioutil"
   "log"
   "os"
   // "regexp"
   "sort"
   "strings"
   
   "github.com/dave/dst"
   "github.com/dave/dst/decorator"
)

type alias int

var fset = token.NewFileSet()

func (*alias) String(p1, p2 int, bb string) string {
	return ""
}

// func liner(// 	bb []byte) []byte {
// 	changed := true
// 	for changed {
// 		changed = false
// 		rea := regexp.MustCompile("func [(]([*._a-zA-Z0-9]+)[)]")
// 		bb = rea.ReplaceAll(bb, []byte("func(_ $1)"))
// 		reor := regexp.MustCompile(
// 			"func [(]([_a-zA-Z0-9]+[ ][*._a-zA-Z0-9]+)[)]")
// 		bb = reor.ReplaceAll(bb, []byte("func (\n\t$1,\n)"))
// 		reors := regexp.MustCompile(
// 			"func [(]\n\t([_a-zA-Z0-9]+[ ][*._a-zA-Z0-9]+)[)]")
// 		bb = reors.ReplaceAll(bb, []byte("func (\n\t$1,\n)"))
// 		reore := regexp.MustCompile(
// 			"func [(]([_a-zA-Z0-9]+[ ][*._a-zA-Z0-9]+),\n[)]")
// 		bb = reore.ReplaceAll(bb, []byte("func (\n\t$1,\n)"))
// 		rep := regexp.MustCompile(
// 			"(\n[)] [_a-zA-Z0-9]+[(])([_a-zA-Z0-9])",
// 		)
// 		bb = rep.ReplaceAll(bb, []byte("$1\n\t$2"))
// 	}
// 	// testing
// 	return bb
// }
func main() {
	// bb, e := ioutil.ReadFile(os.Args[1])
	// if e != nil {
	// 	panic(e)
	// }
	// bb = liner(bb)
	// if e := ioutil.WriteFile(os.Args[1], bb, 0644); e != nil {
	// 	panic(e)
	// }
	bb, e := ioutil.ReadFile(os.Args[1])
	if e != nil {
		panic(e)
	}
	bb = sorter(bb)
	if e := ioutil.WriteFile(os.Args[1], bb, 0644); e != nil {
		panic(e)
	}
}
func sorter(	bb []byte) []byte {
	ss := string(bb)
	splittedraw := strings.Split(ss, "\n")
	imports := []string{}
	for i, x := range splittedraw {
		if x == "import (" {
			impfound := false
			for j := i; !impfound; j++ {
				imports = append(imports, splittedraw[j])
				if splittedraw[j] == ")" {
					goto imported
				}
			}
		}
	}
imported:
	file, err := decorator.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.ERROR(err)
log.Fatal(err)
	}
	constcounter := 0
	unsortedDecls := make(map[string]dst.Decl)
	for _, decl := range file.Decls {
		switch gen := decl.(type) {
		case *dst.GenDecl:
			switch gen.Tok {
			case token.CONST:
				already := false
				for _, x := range unsortedDecls {
					if decl == x {
						already = true
					}
				}
				if already {
					continue
				}
				unsortedDecls["2const"+fmt.Sprint(constcounter)] = decl
				constcounter++
			case token.VAR:
				for _, y := range gen.Specs {
					if is, ok := y.(*dst.ValueSpec); ok {
						var declName string
						key := "3:"
						for i, z := range is.Names {
							key += fmt.Sprint(z)
							_ = i
							declName = z.Name
						}
						d := dst.Clone(decl).(*dst.GenDecl)
						for i, x := range decl.(*dst.GenDecl).Specs {
							dd := d.Specs[i].(*dst.ValueSpec)
							if dd.Names[0].Name == declName {
								d.Specs = nil
								d.Specs = append(d.Specs, x)
								break
							}
						}
						unsortedDecls[key] = d
					}
				}
			case token.TYPE:
				for _, y := range gen.Specs {
					if is, ok := y.(*dst.TypeSpec); ok {
						var declName string
						key := "1:"
						declName = is.Name.Name
						d := dst.Clone(decl).(*dst.GenDecl)
						for i, x := range decl.(*dst.GenDecl).Specs {
							dd := d.Specs[i].(*dst.TypeSpec)
							if dd.Name.Name == declName {
								key += dd.Name.Name
								d.Specs = nil
								d.Specs = append(d.Specs, x)
								break
							}
						}
						unsortedDecls[key] = d
					}
				}
			}
		case *dst.FuncDecl:
			if fun, ok := decl.(*dst.FuncDecl); ok {
				if fun.Recv != nil {
					for _, x := range fun.Recv.List {
						var typename, star string
						if rt, ok := x.Type.(*dst.StarExpr); ok {
							star = "*"
							typename = fmt.Sprint(rt.X)
						} else {
							typename = fmt.Sprint(x.Type)
						}
						for range x.Names {
							unsortedDecls[fmt.Sprintf("4:(%s%s) %s", star, typename, fun.Name.Name)] = decl
						}
					}
				} else {
					unsortedDecls[fmt.Sprint("4:", fun.Name.Name)] = decl
				}
			}
		}
	}
	var sortedDecls []string
	for i := range unsortedDecls {
		sortedDecls = append(sortedDecls, i)
	}
	var decls []dst.Decl
	sort.Strings(sortedDecls)
	for _, x := range sortedDecls {
		decls = append(decls, unsortedDecls[x])
	}
	file.Decls = decls
	var buf bytes.Buffer
	decorator.Fprint(&buf, file)
	output := buf.String()
	var splitout []string
	splitted := strings.Split(output, "\n")
	packagefound := false
	if len(imports) < 1 {
		packagefound = true
	}
	for _, x := range splitted {
		splitout = append(splitout, x)
		if !packagefound {
			sss := strings.Split(x, " ")
			if sss[0] == "package" {
				splitout = append(splitout, "")
				splitout = append(splitout, imports...)
			}
		}
	}
	joined := []byte(strings.Join(splitout, "\n"))
	return joined
}
