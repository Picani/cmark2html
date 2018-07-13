package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aymerick/raymond"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/russross/blackfriday.v2"
)

// EnvironMap converts the result of os.Environ into a map.
func EnvironMap() (mapping map[string]string) {
	mapping = make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		mapping[pair[0]] = pair[1]
	}
	return
}

// initTemplatesDir creates the templates directory if it doesn't exist
// and returns its path.
func initTemplatesDir() (string, error) {
	env := EnvironMap()
	var (
		ok          bool
		templatedir string
	)

	if templatedir, ok = env["XDG_DATA_HOME"]; !ok {
		templatedir = filepath.Join(env["HOME"], ".local/share/cmark2html/")
	}

	if _, err := os.Stat(templatedir); os.IsNotExist(err) {
		err = os.Mkdir(templatedir, 0755)
		if err != nil {
			return "", err
		}
	}

	return templatedir, nil
}

// ListTemplates returns the names of the available templates.
func ListTemplates() ([]string, error) {
	templates := make([]string, 0)
	dir, err := initTemplatesDir()
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".html" {
			name := strings.TrimSuffix(filepath.Base(f.Name()), filepath.Ext(f.Name()))
			templates = append(templates, name)
		}
	}

	return templates, nil
}

// GetTemplate returns the full path for the wanted template.
func GetTemplate(name string) (string, error) {
	dir, err := initTemplatesDir()
	if err != nil {
		return "", err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".html" {
			continue
		}
		basename := strings.TrimSuffix(filepath.Base(f.Name()), filepath.Ext(f.Name()))
		if basename == name {
			return filepath.Join(dir, f.Name()), nil
		}
	}
	return "", errors.New("Template not found")
}

// CompileFile compiles the file mdFile to HTML, inserts its content into
// the "content" tag of the templateFile and writes the result to a file
// with the same name as mdFile but with the extension replaced by .html
func CompileFile(mdFile, templateFile string) error {
	mdContent, err := ioutil.ReadFile(mdFile)
	if err != nil {
		return err
	}
	templateContent, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	htmlContent := blackfriday.Run(mdContent)
	ctx := map[string]interface{}{
		"content": raymond.SafeString(string(htmlContent)),
	}

	final, err := raymond.Render(string(templateContent), ctx)
	if err != nil {
		return err
	}

	basename := strings.TrimSuffix(filepath.Base(mdFile), filepath.Ext(mdFile))
	err = ioutil.WriteFile(basename+".html", []byte(final), 0644)
	if err != nil {
		return err
	}
	return nil
}

var (
	files    = kingpin.Arg("infile.md", "The input CommonMark file(s).").Strings()
	template = kingpin.Flag("template", "The template to use.").Short('t').String()
	list     = kingpin.Flag("list", "List available templates and exit.").Short('l').Bool()
)

func main() {
	kingpin.Version(fmt.Sprintf("cmark2html 1.0.0 - Blackfriday %v", blackfriday.Version))
	kingpin.Parse()

	if !(*list) && *files == nil {
		kingpin.Usage()
		os.Exit(1)
	}

	if *list {
		templates, err := ListTemplates()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, t := range templates {
			fmt.Println(t)
		}
	}

	var (
		tFile string
		err   error
	)

	if *template == "" {
		tFile, err = GetTemplate("default")
	} else {
		tFile, err = GetTemplate(*template)
	}

	if err != nil {
		fmt.Println(err)
	}

	for _, f := range *files {
		err = nil
		err = CompileFile(f, tFile)
		if err != nil {
			fmt.Println(err)
		}
	}
}
