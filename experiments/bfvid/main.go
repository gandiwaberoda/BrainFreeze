package main

import (
	"fmt"

	"harianugrah.com/brainfreeze/pkg/bfvid"
)

func main() {
	// cmd := `
	// 	all:planned(
	// 		cyan:receive,
	// 		magenta:getball,
	// 		magenta:approach(cyan),
	// 		magenta:lookat(cyan),
	// 		magenta:getball,
	// 		magenta:idle@dur(2500),
	// 		Magenta:passing,
	// 		magenta:approach(fgp),
	// 		magenta:receive,
	// 		cyan:lookat(magenta),
	// 		Cyan:passing,
	// 		cyan:dribble(
	// 			approach(fgp),
	// 			lookat(magenta),
	// 			kick
	// 		)@complex,
	// 		magenta:getball,
	// 		cyan:receive,
	// 		magenta:lookat(cyan),
	// 		magenta:getball,
	// 		magenta:idle@dur(2500),
	// 		Magenta:passing,
	// 		cyan:approach(fgp),
	// 		cyan:lookat(fgp),
	// 		cyan:kick
	// 	)@and(
	// 		glanced(fgp), lostball,dist(fgp,300)
	// 	)
	// 	`
	cmd := "all:idle@hold"

	parsed, err := bfvid.ParseCommandSPOK(cmd)
	if err != nil {
		panic(fmt.Sprint("error parsing command:", err))
	}

	fmt.Println("Receiver:", parsed.Receiver)
	fmt.Println("Verb:", parsed.Verb)
	fmt.Println("Parameter:", parsed.Parameter)
	fmt.Println("Cleaned:", parsed.Cleaned)
	// fmt.Println("Raw Fulfilment:", parsed.FulfilmentRaw)
	fmt.Println("Fulfilment Verb:", parsed.Fulfilment)
	fmt.Println("Fulfilment Parameter:", parsed.FulfilmentParameter)

	// if strings.EqualFold(parsed.Verb, "planned") {
	// 	for i, v := range parsed.Parameter {
	// 		subcmd, err := ParseCommandSPOK(v)
	// 		if err != nil {
	// 			fmt.Println("Error parsing command:", err)
	// 		}
	// 		fmt.Println(fmt.Sprint("Subcmd", i, ":"), v, "=>", subcmd.Verb, "|", subcmd.Parameter)
	// 		// fmt.Println("\tSubcmd Verb", i, ":")
	// 	}
	// }

}
