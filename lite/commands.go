package lite

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
	"github.com/xyproto/otto"
	"path/filepath"
)

var adm = `
(function() {
  var template = Handlebars.template, templates = Handlebars.templates = Handlebars.templates || {};
		{{range . }}
			{{. -}}
		{{end}}
})();

`

var sdm = `
  templates["%s"] = template(%s);
`
var tmpl *template.Template

var precompileMap map[string]string

var vm *otto.Otto

func _init() {
	vm = otto.New()
	handlebarsJS, _ := Asset("handlebars.js")
	_, err := vm.Run(handlebarsJS)
	if err != nil {
		log.Fatal(err)
	}
	precompileMap = make(map[string]string)
	tmpl, _ = template.New("tmpl").Parse(adm)

	fmt.Println("w init")

	for _, path := range conf.Templates.Paths {
		filepath.Walk(path, Walkfu)
		if err != nil {
			log.Fatal(err)
		}
	}

	// vm.Set("Handlebars", vm.Get("Handlebars"))
}

// WatchAction is the watch command action
func WatchAction(c *cli.Context) error {

	// beforeWatch()
	_init()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
			  _, name := filepath.Split(event.Name)
				ext := filepath.Ext(name)
				name = name[0:len(name)-len(ext)]

				if event.Op&fsnotify.Write == fsnotify.Write {
					begin := time.Now()
					content, _ := ioutil.ReadFile(event.Name)
					precompile(name, content)

					duration := time.Since(begin)
					fmt.Printf("%s  in %.2fs\n", Emojitize(":coffee: compile file:" + event.Name),  duration.Seconds())

				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {

					begin := time.Now()
					delete(precompileMap, name)
					write()

					duration := time.Since(begin)
					fmt.Printf("%s  in %.2f|s\n", Emojitize(":coffee: compile file:" + event.Name),  duration.Seconds())

				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {

					begin := time.Now()
					delete(precompileMap, name)
					write()

					duration := time.Since(begin)
					fmt.Printf("%s  in %.2fs\n", Emojitize(":coffee: compile file:" + event.Name),  duration.Seconds())

				}
			case _ = <-watcher.Errors:
				// log.Println("error: ", err)
			}
		}
	}()

	fmt.Println(Emojitize(":moon: ready to watch :moon:"))
	for _, path := range conf.Templates.Paths {
		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	<-done
	log.Println("exiting watching")
	return nil

}

func precompile(name string, content []byte) error {

	vm.Set("tmpl", string(content))
	jsTmpl, err := vm.Run("Handlebars.precompile(tmpl)")

	if err != nil {
		return err
	}

	precompileMap[name] = fmt.Sprintf(sdm, name, jsTmpl)

	write()

	return nil
}

func write() {
	file, err := os.OpenFile(conf.Templates.Precompile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		file ,_ = os.Create(conf.Templates.Precompile)
	}
	err = tmpl.Execute(file, precompileMap)
	if err != nil {
		log.Println(err)
	}
}
