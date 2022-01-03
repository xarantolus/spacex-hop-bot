package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Twitter struct {
		AccessToken       string `yaml:"access_token"`
		AccessTokenSecret string `yaml:"access_secret"`
		APIKey            string `yaml:"api_key"`
		APISecretKey      string `yaml:"api_secret"`
	} `yaml:"twitter"`

	Lists struct {
		MainStarshipListID int64   `yaml:"main_starship_list"`
		IgnoredListIDs     []int64 `yaml:"ignored_lists"`
	} `yaml:"lists"`
}

func (c Config) IgnoredListsMapping() (mapping map[int64]bool) {
	mapping = make(map[int64]bool)

	for _, lid := range c.Lists.IgnoredListIDs {
		mapping[lid] = true
	}

	return
}

func Parse(filename string) (c Config, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		return
	}

	return
}
