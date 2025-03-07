package generator

import (
	"strings"

	"github.com/sethvargo/go-diceware/diceware"
)

func GeneratePassphrase(words int) string {
	list, _ := diceware.Generate(words)

	return strings.Join(list, " ")
}