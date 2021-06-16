package configuration

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	"harianugrah.com/brainfreeze/pkg/errors"
)

type RobotConfig struct {
	Name string    `yaml:"name"`
	Role RobotType `yaml:"role"`
}

type CameraConfig struct {
	Src         []string `yaml:"src"`
	UseDshow    bool     `yaml:"useDshow"`
	MidpointX   int      `yaml:"midpointX"`
	MidpointY   int      `yaml:"midpointY"`
	MidpointRad int      `yaml:"midpointRad"`
	RawWidth    int      `yaml:"rawWidth"`
	RawHeight   int      `yaml:"rawHeight"`
	PostWidth   int      `yaml:"postWidth"`
	PostHeight  int      `yaml:"postHeight"`
}

type ExpirationConfig struct {
	BallExpiration time.Duration `yaml:"ball"`
	MyExpiration   time.Duration `yaml:"my"`
}

type DiagnosticConfig struct {
	TelemetryHz        time.Duration `yaml:"telemetryHz"`
	EnableStream       bool          `yaml:"enableStream"`
	StreamTopProcessed bool          `yaml:"streamTopProcessed"`
	StreamHost         string        `yaml:"streamHost"`
}

type MigraineConfig struct {
	MigraineHz int `yaml:"migraineHz"`
}

type MechanicalConfig struct {
	HorizontalForceRange int `yaml:"horizontalForceRange"`
	VerticalForceRange   int `yaml:"verticalForceRange"`
}

type FulfillmentConfig struct {
	DefaultDurationMs int `yaml:"defaultDurationMs"`
}

type FreezeConfig struct {
	Robot       RobotConfig       `yaml:"robot"`
	Camera      CameraConfig      `yaml:"camera"`
	Expiration  ExpirationConfig  `yaml:"expiration"`
	Diagnostic  DiagnosticConfig  `yaml:"diagnostic"`
	Migraine    MigraineConfig    `yaml:"migraine"`
	Mecha       MechanicalConfig  `yaml:"mecha"`
	Fulfillment FulfillmentConfig `yaml:"fulfillment"`
}

func LoadStartupConfig() (FreezeConfig, error) {
	conf := FreezeConfig{}

	reader, err := os.Open("./config.yaml")
	if err != nil {
		return FreezeConfig{}, &frerror.ConfigError{Detail: err.Error()}
	}
	defer reader.Close()

	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&conf); err != nil {
		return FreezeConfig{}, &frerror.ConfigError{Detail: err.Error()}
	}

	if !validateConfig(conf) {
		log.Fatalln("Startup config tidak beres")
	}

	return conf, nil
}

// Validasi configurasi yang telat dimaut
func validateConfig(conf FreezeConfig) bool {
	return true
}
