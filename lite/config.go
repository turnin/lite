package lite

import (
	"fmt"
	"io/ioutil"
	"log"
)

import (
	"github.com/BurntSushi/toml"
)

// Config is the configuration of lite
type config struct {
	Title     string
	Templates templates
}

type templates struct {
	Paths      []string
	Precompile string
}

var conf config

// Load is loading the config file
func Load(configPath string) error {
	configData, _ := ioutil.ReadFile(configPath)

	fmt.Println(configData)
	_, err := toml.Decode(string(configData), &conf)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Printf("--- config:\n%v\n\n", conf)

	return nil

}
