package helper

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

//HELPER
func AmIReceiver(identifier string, m *configuration.FreezeConfig) bool {
	_amIReceiver := false

	myReceiverTag := []string{string(models.ALL), string(m.Robot.Name), string(m.Robot.Role), string(m.Robot.Color)}

	for _, v := range myReceiverTag {
		// Case insensitive
		if strings.EqualFold(identifier, v) {
			_amIReceiver = true
		}
	}
	return _amIReceiver
}
