package localconfig

import (
	"os"

	"gopkg.in/yaml.v2"
	"harianugrah.com/brainfreeze/pkg/models"
)

type RobotConfig struct {
	Name string           `yaml:"name"`
	Role models.RobotType `yaml:"role"`
}

type CameraConfig struct {
	Src      int  `yaml:"src"`
	UseDshow bool `yaml:"useDshow"`
}

type FreezeConfig struct {
	Robot  RobotConfig  `yaml:"robot"`
	Camera CameraConfig `yaml:"camera"`
}

func LoadFreezeConfig() (FreezeConfig, error) {
	conf := FreezeConfig{}

	reader, err := os.Open("./config.yaml")
	if err != nil {
		return FreezeConfig{}, err
	}
	defer reader.Close()

	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&conf); err != nil {
		return FreezeConfig{}, err
	}

	return conf, nil
}
