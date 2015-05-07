package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"unicode"
)

// Template for generated assets source file.
const tpl = `package assets

var {{.Name}} = map[string]string {
{{range $file, $content := .Assets}}
	"{{$file}}": ` + "`" + `{{$content}}` + "`" + `,
{{end}}
}
`

var dir string

// Data for template rendering.
var tplData struct {
	Name   string
	Assets map[string]string
}

func main() {
	dir = getWorkingDir()

	tplData.Name = getName(dir)
	tplData.Assets = make(map[string]string)

	filepath.Walk(dir, walk)
	t := template.Must(template.New(tplData.Name).Parse(tpl))

	out := getOutput(dir)
	t.Execute(out, tplData)
	log.Println("assets file hase been generated:", out.Name())
}

// Get working dir path from CLI arguments.
func getWorkingDir() string {
	var dir string

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		dir = args[0]
	} else {
		dir = "."
	}

	dir = path.Clean(dir)
	dir, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}

	fi, err := os.Stat(dir)
	if err != nil {
		panic(err)
	}

	if !fi.IsDir() {
		panic("is not directory")
	}

	return dir
}

// Get assets name. Get last end directory name and put it first character
// to the uppercase.
func getName(dir string) string {
	name := []rune(path.Base(dir))
	name[0] = unicode.ToUpper(name[0])
	return string(name)
}

// Create file to write generated output.
func getOutput(n string) *os.File {
	out, err := os.Create(n + ".go")
	if err != nil {
		log.Fatalln(err)
	}
	return out
}

// Walk by assets' directory and generated map of assets.
func walk(path string, fi os.FileInfo, err error) error {
	if err != nil {
		log.Println(err)
		return err
	}

	if !fi.IsDir() {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
			return err
		}

		k, err := filepath.Rel(dir, path)
		if err != nil {
			log.Println(err)
			return err
		}

		// TODO: Escape back quotes.
		tplData.Assets[k] = string(content)
	}

	return nil
}
