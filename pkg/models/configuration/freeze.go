package configuration

import (
	"image"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	frerror "harianugrah.com/brainfreeze/pkg/errors"
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
	Midpoint    image.Point
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

type TelepathyConfig struct {
	ChitChatHost string `yaml:"chitchatHost"`
}

type WandaConfig struct {
	MinimumHsvArea float64 `yaml:"minimumHsvArea"`
}

type FreezeConfig struct {
	Robot       RobotConfig       `yaml:"robot"`
	Camera      CameraConfig      `yaml:"camera"`
	Expiration  ExpirationConfig  `yaml:"expiration"`
	Diagnostic  DiagnosticConfig  `yaml:"diagnostic"`
	Migraine    MigraineConfig    `yaml:"migraine"`
	Mecha       MechanicalConfig  `yaml:"mecha"`
	Fulfillment FulfillmentConfig `yaml:"fulfillment"`
	Telepathy   TelepathyConfig   `yaml:"telepathy"`
	Wanda       WandaConfig       `yaml:"wanda"`
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

	conf.Camera.Midpoint = image.Point{conf.Camera.PostWidth / 2, conf.Camera.PostHeight / 2}

	return conf, nil
}
