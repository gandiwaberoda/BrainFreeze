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
	GAMESTATE MessageType = "GAMESTATE"
	COMMAND   MessageType = "COMMAND"
	REPLY     MessageType = "REPLY"
)

// Message receiver bisa all ALL atau nama dari robot yang ada di yaml
type MessageReceiver string

const (
	ALL      MessageReceiver = "ALL"
	ROBOT    MessageReceiver = "ROBOT"
	ATTACKER MessageReceiver = "ATTACKER"
	DEFENDER MessageReceiver = "DEFENDER"
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

	// TODO: Sanitized input later
	sanitized := raw

	if sanitized[0] == '{' {
		// json
		err := json.Unmarshal([]byte(sanitized), &parsed)
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

	splitted := strings.Split(sanitized, "/")
	if len(splitted) < 2 {
		return Intercom{}, errors.New("wrong command format")
	}

	// Parse KIND, asumsikan semua yang formatnya manusiawi adalah COMMAND
	parsed.Kind = COMMAND

	// Parse receiver
	receiverUpper := strings.ToUpper(splitted[0])
	parsed.Receiver = MessageReceiver(receiverUpper)

	// Parse content (Content adalah sisa setelah 2 tadi)
	contentSubs := sanitized[len(receiverUpper):]
	contentNoDelim := contentSubs[1:]
	contentTrimmed := strings.TrimSpace(contentNoDelim)
	parsed.Content = contentTrimmed
	return parsed, nil
}

func (s Intercom) AsBytes() ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, errors.New("gagal get bytes")
	}
	return b, nil
}

func (s Intercom) AsJson() (string, error) {
	b, err := s.AsBytes()
	if err != nil {
		return "", err
	}
	jsonMsg := string(b)
	return jsonMsg, nil
}
