package commands

import (
	"fmt"
	"runtime/debug"
)

func InfoCommand(arguments []string) error {
	commitHash, commitTime := getInfo()
	fmt.Println("6502 Profiler")
	fmt.Println("Version 1.3.1")
	fmt.Printf("Commit hash: %s\n", commitHash)
	fmt.Printf("Commit time: %s\n", commitTime)
	fmt.Println("Written by Martin Grap (rmsk2@gmx.de) in 2023")

	return nil
}

func getInfo() (string, string) {
	var hash string
	var time string

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				hash = setting.Value
				continue
			}

			if setting.Key == "vcs.time" {
				time = setting.Value
				continue
			}
		}
	}

	return hash, time
}
