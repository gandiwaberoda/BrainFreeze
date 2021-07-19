package bfvid

import (
	"errors"
	"fmt"
	"strings"
)

type CommandSPOK struct {
	Receiver            string
	Verb                string
	Parameter           []string
	Raw                 string
	Cleaned             string
	ParameterStr        string
	Fulfilment          string
	FulfilmentRaw       string
	FulfilmentParameter []string
}

func ParseCommandSPOK(cmd string) (*CommandSPOK, error) {
	s := cmd
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")

	if strings.Contains(s, ",)") {
		return nil, errors.New("parameter contain \",)\" character")
	}

	level := -1
	level0BracketOpening := -1
	lastParameterPointer := -1
	level0FulfilmentOpening := -1

	foundOpeningBracket := false
	stopCommandParsing := false
	usedParameterStart := make(map[int]bool)

	var paramStr string
	var verb string
	var receiver string
	var rawfulfilment string
	var verbfulfilment string
	var fulfilmentparams []string

	param := make([]string, 0)

	verbStart := 0

	for i, v := range s {
		if v == '(' && !stopCommandParsing {
			level++
			foundOpeningBracket = true
			if level == 0 {
				level0BracketOpening = i
				lastParameterPointer = i + 1
				usedParameterStart[lastParameterPointer] = false
				verb = s[verbStart:i]
			}
		}

		if v == ')' && !stopCommandParsing {
			if level == 0 {
				paramStr = s[level0BracketOpening+1 : i]

				if !usedParameterStart[lastParameterPointer] {
					p := s[lastParameterPointer:i]
					param = append(param, p)
					usedParameterStart[lastParameterPointer] = true
				}
			}
			level--
		}

		if v == ':' && !stopCommandParsing {
			if level == -1 {
				receiver = s[:i]
				verbStart = i + 1
			}
		}

		if v == ',' && level == 0 && !stopCommandParsing {
			if !usedParameterStart[lastParameterPointer] {
				p := s[lastParameterPointer:i]
				param = append(param, p)
				usedParameterStart[lastParameterPointer] = true

				lastParameterPointer = i + 1
				usedParameterStart[lastParameterPointer] = false
			}
		}

		if v == '@' && !stopCommandParsing {
			if level == -1 {
				level0FulfilmentOpening = i + 1
				stopCommandParsing = true
			}
		}
	}

	if level != -1 {
		return nil, errors.New("pair of ( and ) is not match")
	}

	if !foundOpeningBracket {
		verb = s[verbStart:]
	}

	if level0FulfilmentOpening != -1 {
		rawfulfilment = s[level0FulfilmentOpening:]
	}

	if rawfulfilment != "" {
		parseFulfilment, err := ParseCommandSPOK(rawfulfilment)
		if err != nil {
			return nil, errors.New(fmt.Sprint("failed to parse fulfilment:", err))
		}
		verbfulfilment = parseFulfilment.Verb
		fulfilmentparams = parseFulfilment.Parameter
	}

	return &CommandSPOK{
		Receiver:            receiver,
		Parameter:           param,
		Raw:                 cmd,
		Cleaned:             s,
		Verb:                verb,
		ParameterStr:        paramStr,
		FulfilmentRaw:       rawfulfilment,
		Fulfilment:          verbfulfilment,
		FulfilmentParameter: fulfilmentparams,
	}, nil
}
