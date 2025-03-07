package main

import (
	"fmt"
	"os"
)

// set debug constant inside the "config" package
func generateDebugConfig() {
	fmt.Printf("Generating DEBUG/RELEASE config")

	code := METAGEN_AUTO_COMMENT + "\npackage config\n\nconst (\n"

	if envtype == ENVIRONMENT_DEV {
		code += "	DEBUG = true"
	} else {
		code += "	DEBUG = false"
	}

	code += "\n)\n"

	// open file and write code to it
	in := []byte(code)

	err := os.WriteFile("./config/debug.metagen.go", in, 0644)
	handleErr(err)

	printStatus(true)
}
