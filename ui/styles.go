package ui

import (
	"maxwarden/basic"
	"maxwarden/security"
	"strings"

	. "maragu.dev/gomponents"
)

// @Macro - At runtime this inserts an attribute component with the hash of the input
// At compile time, the input is collected, and used to generate a global css file.
// WARNING: THE PREPROCESSOR EXPANSION ONLY WORKS WITH *STRING LITERAL* INPUTS
// ATTEMPTING TO USE THIS WITH DYNAMIC INPUTS WILL FAIL COMPILATION
//
// For generating "dynamic" styles, use the If() / IfElse() component, and multiple calls
// to `InlineStyle`.
//
// Ex:
//
//	func MyComponent(cond bool) Node {
//	    return IfElse(cond, InlineStyle("color: green;"), InlineStyle("color: red;"))
//	}
//
// THE FOLLOWING DOES NOT WORK:
//
//	func MyInvalidComponent(cond bool) Node {
//	    var css string
//
//	    if cond {
//	        css = "color: green;"
//	    } else {
//	        css = "color: red;"
//	    }
//
//	    return InlineStyle(css) <--- THIS IS AN ERROR! THE PREPROCESSOR HAS NO IDEA WHAT THE VALUE OF `css` IS,
//	                                 SINCE IT IS NOT KNOWN AT COMPILE TIME!
//	}
func InlineStyle(input string) Node {
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\t", "")

	s, _ := security.HighwayHash58(input)
	s = basic.GetFirstNChars(s, 8)

	return Attr("__inlinecss_" + s)
}
