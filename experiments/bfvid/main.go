package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {
	cmd := `
		all:planned(
			cyan:receive,
			magenta:getball,
			magenta:approach(cyan),
			magenta:lookat(cyan),
			magenta:getball,
			magenta:idle@dur(2500),
			Magenta:passing,
			magenta:approach(fgp),
			magenta:receive,
			cyan:lookat(magenta),
			Cyan:passing,
			magenta:getball,
			cyan:receive,
			magenta:lookat(cyan),
			magenta:getball,
			magenta:idle@dur(2500),
			Magenta:passing,
			cyan:approach(fgp),
			cyan:lookat(fgp),
			cyan:kick,
		)
		`

	//fmt.Println(cmd)
	parsed, err := ParseCommandSPOK(cmd)
	if err != nil {
		fmt.Println("Error parsing command:", err)
	}
	// fmt.Println(parsed)

	fmt.Println("Receiver:", parsed.Receiver)
	fmt.Println("Verb:", parsed.Verb)
	fmt.Println("Parameter:", parsed.Parameter)

	if strings.EqualFold(parsed.Verb, "planned") {
		for i, v := range parsed.Parameter {
			subcmd, err := ParseCommandSPOK(v)
			if err != nil {
				fmt.Println("Error parsing command:", err)
			}
			fmt.Println(fmt.Sprint("Subcmd", i, ":"), v, "=>", subcmd.Verb, "|", subcmd.Parameter)
			// fmt.Println("\tSubcmd Verb", i, ":")
		}
	}

}

type CommandSPOK struct {
	Receiver            string
	Verb                string
	Parameter           []string
	Raw                 string
	Cleaned             string
	ParameterStr        string
	Fulfilment          string
	FulfilmentParameter []string
}

func ParseCommandSPOK(cmd string) (*CommandSPOK, error) {
	s := cmd
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")

	level := -1
	level0BracketOpening := -1
	lastParameterPointer := -1
	foundOpeningBracket := false

	var paramStr string
	var verb string
	var receiver string

	param := make([]string, 0)

	verbStart := 0

	for i, v := range s {
		//fmt.Println(v)
		if v == '(' {
			level++
			foundOpeningBracket = true
			if level == 0 {
				level0BracketOpening = i
				lastParameterPointer = i + 1
				verb = s[verbStart:i]
			}
		}

		if v == ')' {
			if level == 0 {
				paramStr = s[level0BracketOpening+1 : i]
			}
			level--
		}

		if v == ':' {
			if level == -1 {
				receiver = s[:i]
				verbStart = i + 1
			}
		}

		if v == ',' {
			p := s[lastParameterPointer:i]
			param = append(param, p)
			lastParameterPointer = i + 1
		}
	}

	if level != -1 {
		return nil, errors.New("pair of ( and ) is not match")
	}

	if !foundOpeningBracket {
		// verb = s[verbStart:]
	}

	return &CommandSPOK{
		Receiver:     receiver,
		Parameter:    param,
		Raw:          cmd,
		Cleaned:      s,
		Verb:         verb,
		ParameterStr: paramStr,
	}, nil
}
