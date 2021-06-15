// Untuk komunikasi dengan dan atara
// websocket client (plain text), basestation dan robot yang lain
package models

import (
	"encoding/json"
	"errors"
	"strings"
)

type MessageType string

const (
	TELEMETRY MessageType = "TELEMETRY"
	COMMAND   MessageType = "COMMAND"
	REPLY     MessageType = "REPLY"
)

// Message receiver bisa all ALL atau nama dari robot yang ada di yaml
type MessageReceiver string

const (
	ALL MessageReceiver = "ALL"
	ANY MessageReceiver = "ANY"
)

type Intercom struct {
	Kind     MessageType
	Receiver MessageReceiver
	Content  string
}

// Parse json string ke Intercom struct
// json itu case insensitive
func ParseIntercom(raw string) (Intercom, error) {
	var parsed Intercom

	if raw[0] == '{' {
		// json
		err := json.Unmarshal([]byte(raw), &parsed)
		if err != nil {
			return Intercom{}, err
		}
		return parsed, nil
	}

	// command dari websocket client

	// Contoh command (Case insensitive)
	// all getball
	// gandiwa getball
	// aTTacKer kick
	// deFFenDER plannED(goto(10,30) dur(3000))

	splitted := strings.Split(raw, " ")
	if len(splitted) < 2 {
		return Intercom{}, errors.New("wrong command format")
	}

	// Parse KIND, asumsikan semua yang formatnya manusiawi adalah COMMAND
	parsed.Kind = COMMAND

	// Parse receiver
	receiverUpper := strings.ToUpper(splitted[0])
	parsed.Receiver = MessageReceiver(receiverUpper)

	// Parse content (Content adalah sisa setelah 2 tadi)
	contentSubs := raw[len(receiverUpper):]
	contentTrimmed := strings.TrimSpace(contentSubs)
	parsed.Content = contentTrimmed
	return parsed, nil
}
