package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/ng-vu/githubx/backend/pkg/github"
)

type Config struct {
	GitHub  github.Config  `yaml:"github"`
	DBFile  string         `yaml:"db_file"`
	LogFile string         `yaml:"log_file"`
	Listen  string         `yaml:"listen"`
	OrgRepo github.OrgRepo `yaml:",inline"`
}

func Default() Config {
	cfg := Config{
		Listen: ":8080",
	}
	return cfg
}

func Load(filename string) (Config, error) {
	cfg := Default()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}
