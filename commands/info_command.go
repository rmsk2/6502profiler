package commands

import "fmt"

func InfoCommand(arguments []string) error {
	fmt.Println("6502 Profiler")
	fmt.Println("Version 0.9.5")
	fmt.Println("Written by Martin Grap (rmsk2@gmx.de) in 2023")

	return nil
}
