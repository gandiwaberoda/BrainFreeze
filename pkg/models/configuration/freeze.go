package configuration

import (
	"image"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	frerror "harianugrah.com/brainfreeze/pkg/errors"
)

type RobotColor string

const (
	MAGENTA RobotColor = "MAGENTA"
	CYAN    RobotColor = "CYAN"
)

type RobotConfig struct {
	Name     string     `yaml:"name"`
	Role     RobotType  `yaml:"role"`
	StartPos string     `yaml:"startPos"`
	StartRot int        `yaml:"startRot"`
	Color    RobotColor `yaml:"color"`
}

type CameraConfig struct {
	Src        []string `yaml:"src"`
	SrcForward []string `yaml:"srcForward"`
	UseDshow   bool     `yaml:"useDshow"`

	MidpointX         int     `yaml:"midpointX"`
	MidpointY         int     `yaml:"midpointY"`
	MidpointRad       int     `yaml:"midpointRad"`
	RawWidth          int     `yaml:"rawWidth"`
	RawHeight         int     `yaml:"rawHeight"`
	PostWidth         int     `yaml:"postWidth"`
	PostHeight        int     `yaml:"postHeight"`
	RobFrontOffsetDeg int     `yaml:"robFrontOffsetDeg"`
	TopRobRotPatch    float64 `yaml:"topRobRotPatch"`

	ForWidth      int `yaml:"forWidth"`
	ForHeight     int `yaml:"forHeight"`
	ForPostWidth  int `yaml:"forPostWidth"`
	ForPostHeight int `yaml:"forPostHeight"`
	ForMidX       int `yaml:"forMidX"`

	Midpoint image.Point
}

type ExpirationConfig struct {
	BallExpiration time.Duration `yaml:"ball"`
	MyExpiration   time.Duration `yaml:"my"`
}

type SerialConfig struct {
	Ports      []string `yaml:"ports"`
	ArayaPorts []string `yaml:"arayaPort"`
	CommandHz  int      `yaml:"commandHz"`
}

type SimulatorConfig struct {
	SimpservPort string `yaml:"simpservPort"`
}

type DiagnosticConfig struct {
	TelemetryHz        time.Duration `yaml:"telemetryHz"`
	EnableStream       bool          `yaml:"enableStream"`
	StreamTopProcessed bool          `yaml:"streamTopProcessed"`
	StreamHost         string        `yaml:"streamHost"`
	ShowScreen         bool          `yaml:"showScreen"`
}

type MigraineConfig struct {
	MigraineHz int `yaml:"migraineHz"`
}

type MechanicalConfig struct {
	HorizontalForceRange  int `yaml:"horizontalForceRange"`
	VerticalForceRange    int `yaml:"verticalForceRange"`
	RotationForceMinRange int `yaml:"rotationForceMinRange"`
	RotationForceMaxRange int `yaml:"rotationForceMaxRange"`
}

type CommandParameterConfig struct {
	LookatToleranceDeg  int `yaml:"lookatToleranceDeg"`
	PositionToleranceCm int `yaml:"positionToleranceCm"`

	ApproachDistanceCm int `yaml:"approachDistanceCm"`

	HandlingOnDist int `yaml:"handlingOnDist"`
	RotToMoveDelay int `yaml:"rotToMoveDelay"`
	// OnlyOneDegreeMovement bool `yaml:"onlyOneDegreeMovement"`

	AllowXYTogether    bool `yaml:"allowXYTogether"`
	AllowXYRotTogether bool `yaml:"allowXYRotTogether"`
}

type FulfillmentConfig struct {
	DefaultDurationMs int `yaml:"defaultDurationMs"`
}

type TelepathyConfig struct {
	ChitChatHost []string `yaml:"chitchatHost"`
}

type WandaConfig struct {
	DisableMagentaDetection bool `yaml:"disableMagentaDetection"`
	DisableCyanDetection    bool `yaml:"disableCyanDetection"`

	MinimumHsvArea         float64 `yaml:"minimumHsvArea"`
	MaximumHsvArea         float64 `yaml:"maximumHsvArea"`
	LerpValue              float64 `yaml:"lerpValue"`
	LfFovMin               int     `yaml:"lfFovMin"`
	LfFovMax               int     `yaml:"lfFovMax"`
	WhiteOnGrayVal         int     `yaml:"whiteOnGrayVal"`
	RadiusGoalpostCircular int     `yaml:"radiusGoalpostCircular"`
	DebugShowRadiusLine    bool    `yaml:"debugShowRadiusLine"`
	GoalOffset             float64 `yaml:"goalOffset"`
	StraightMinLength      int     `yaml:"straightMinLength"`
}

type FreezeConfig struct {
	Robot            RobotConfig            `yaml:"robot"`
	Camera           CameraConfig           `yaml:"camera"`
	Expiration       ExpirationConfig       `yaml:"expiration"`
	Diagnostic       DiagnosticConfig       `yaml:"diagnostic"`
	Migraine         MigraineConfig         `yaml:"migraine"`
	Mecha            MechanicalConfig       `yaml:"mecha"`
	Fulfillment      FulfillmentConfig      `yaml:"fulfillment"`
	Telepathy        TelepathyConfig        `yaml:"telepathy"`
	Wanda            WandaConfig            `yaml:"wanda"`
	CommandParameter CommandParameterConfig `yaml:"commandParameter"`
	Serial           SerialConfig           `yaml:"serial"`
	Simulator        SimulatorConfig        `yaml:"simulator"`
}

func LoadStartupConfigByFile(path string) (FreezeConfig, error) {
	conf := FreezeConfig{}

	reader, err := os.Open(path)
	if err != nil {
		return FreezeConfig{}, &frerror.ConfigError{Detail: err.Error()}
	}
	defer reader.Close()

	// b, _ := ioutil.ReadAll(reader)
	// fmt.Print("xx", b)

	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&conf); err != nil {
		return FreezeConfig{}, &frerror.ConfigError{Detail: err.Error()}
	}

	conf.Camera.Midpoint = image.Point{conf.Camera.PostWidth / 2, conf.Camera.PostHeight / 2}

	return conf, nil
}

func LoadStartupConfig() (FreezeConfig, error) {
	return LoadStartupConfigByFile("./config.yaml")
}
