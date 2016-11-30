package lite

import "io/ioutil"

// Config is the configuration of lite
type Config struct {
	TmplConfig struct {
		Paths      []string `yaml:"paths,flow"`
		Precompile string   `yaml:"precompile"`
	}
}

var config *Config

// Load is loading the config file
func Load(configPath string) error {
	// configData, _ := ioutil.ReadFile("./config/config.yaml")
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	return nil

}
