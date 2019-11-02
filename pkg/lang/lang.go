package lang

type Text struct {
	ID         string
	Definition string
}

type Language struct {
	Code        string
	Definitions []Text
}
type Com struct {
	Component string
	Languages []Language
}
type Dictionary []Com

type Lexicon map[string]string

var dict Dictionary

func ExportLanguage(l string) *Lexicon {
	lex := Lexicon{}
	d := Dictionary{}
	d = append(d, goAppDict())
	for _, c := range d {
		for _, lang := range c.Languages {
			for _, def := range lang.Definitions {
				lex[c.Component+"_"+def.ID] = def.Definition
			}
		}
	}
	return &lex
}

func (l *Lexicon) RenderText(id string) string {
	return (*l)[id]
}
