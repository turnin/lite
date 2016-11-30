package lite

import (
	"log"
	"os"

	"os/signal"
	"syscall"

	"fmt"

	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
	"github.com/xyproto/otto"
)

var vm *otto.Otto

func init() {
	vm = otto.New()
	handlebarsJS, _ := Asset("handlebars.js")
	_, err := vm.Run(handlebarsJS)
	if err != nil {
		log.Fatal(err)
	}
	// vm.Set("Handlebars", vm.Get("Handlebars"))
}

// WatchAction is the watch command action
func WatchAction(c *cli.Context) error {

	// beforeWatch()

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
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println(Emojitize(":coffee: modified file:" + event.Name))
					precompile(event.Name)
					// log.Println("modified file:", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println(Emojitize(":coffee: delete file:" + event.Name))

					// log.Println("delete file: ", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println(Emojitize(":coffee: delete file:" + event.Name))
					// log.Println("rename file: ", event.Name)
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

func precompile(tmplPath string) error {

	content, _ := ioutil.ReadFile(tmplPath)
	vm.Set("tmpl", string(content))
	jsTmpl, err := vm.Run("Handlebars.precompile(tmpl)")

	if err != nil {
		return err
	}

	fmt.Println(jsTmpl)

	return nil
}
