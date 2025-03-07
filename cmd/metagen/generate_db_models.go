package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func compileJet() {
	fmt.Printf("Compiling Jet generator")

	os.Setenv("CGO_ENABLED", "1")

	cmd := exec.Command("go", "build", "./cmd/jet")
	cmd.Dir = "./tools/jet-2.12.0"
	handleCmdOutput(cmd.CombinedOutput())

}

func generateJetModels() {
	bin := ""
	jetdir := ".jet"

	if runtime.GOOS == "windows" {
		bin = "./tools/jet-2.12.0/jet.exe"
	} else {
		bin = "./tools/jet-2.12.0/jet"
	}

	os.RemoveAll(jetdir)

	// compile bin if not exists
	if _, err := os.Stat(bin); err != nil {
		compileJet()
	}

	fmt.Printf("Generating SQL models (jet)")

	if _, err := os.Stat("passwords.db"); err != nil {
		printStatus(false)
		fmt.Println("\n" + err.Error())
		os.Exit(1)
	}

	databaseType := "sqlite"

	cmd := exec.Command(bin, "-source="+databaseType, "-dsn=file:passwords.db", "-schema=maxwarden", "-path="+jetdir)

	handleCmdOutput(cmd.CombinedOutput())
	printStatus(true)
}
