package ui

import (
	"fmt"
	. "maxwarden/basic"
	"maxwarden/config"
	"strings"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// This component generates a <script> tag that automatically includes a blank scope { }.
// Avoid using the javascript `var` keyword inside the scope as it will be placed in the global
// scope.
//
// WARNING: DO NOT PASS USER INPUT TO THIS COMPONENT!!
// THIS COULD BE A ONE WAY TICKET TO XSS IF USED IMPROPERLY!
func InlineScript(script string) Node {
	// minify in production only
	if !config.DEBUG {
		script = strings.ReplaceAll(script, "\n", " ")
		script = strings.ReplaceAll(script, "\t", "")
	}

	return Script(Raw(`{` + script + `}`))
}

// Wrapper function around `InlineScript` that allows for printf style format strings + data
//
// WARNING: DO NOT PASS USER INPUT TO THIS COMPONENT!!
// THIS COULD BE A ONE WAY TICKET TO XSS IF USED IMPROPERLY!
func InlineScriptf(scriptFormat string, items ...interface{}) Node {
	return InlineScript(fmt.Sprintf(scriptFormat, items...))
}

// Convert Go arrays to a string containing a javascript array.
// Useful for injecting data into a dynamically generated script.
func MakeJsArray[T any](list []T) string {
	var out_string string

	out_string += "["

	for i, v := range list {
		if i == len(list)-1 {
			out_string += fmt.Sprintf("'%s'", ToString(v))
		} else {
			out_string += fmt.Sprintf("'%s',", ToString(v))
		}
	}

	out_string += "]"

	return out_string
}
