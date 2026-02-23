package helpers

import "os"

func PrintInfo(message string) {
	println(message)
}

func PrintWarn(message string) {
	println("Warn: " + message)
}

func PrintError(message string) {
	println("Error: " + message)
	os.Exit(1)
}
