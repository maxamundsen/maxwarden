package main

import (
	"fmt"
	"os"
	"path/filepath"
	"maxwarden/basic"
	"maxwarden/security"
	"regexp"
	"strings"
)

var dedupMap = make(map[string]bool)

// Walk the FS tree and search for `.go` files containing calls to `InlineStyle()`
// Collect the inputs to each call (they must be string literals), expand shorthand macros, and
func generateInlineStyles() {
	fmt.Printf("Compiling Inline Styles")

	inlineStyleRegex := regexp.MustCompile("InlineStyle\\((((?:'[^']*')|(?:\"[^\"]*\")|(`(?:[^`]|[\r\n])*?`)))\\)")

	// these are the directories that get scanned
	var dirs = [2]string{"handlers", "ui"}

	var matches []string

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(pathStr string, info os.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".go") {
				pathStr = filepath.ToSlash(pathStr)
				handleErr(err)

				file, err := os.Open(pathStr)
				handleErr(err)

				defer file.Close()

				content, contentErr := os.ReadFile(pathStr)
				if contentErr != nil {
					return contentErr
				}

				s := string(content)

				// only pre-process page & component files
				inlineStyleRegex.ReplaceAllStringFunc(s, func(match string) string {
					submatch := inlineStyleRegex.FindStringSubmatch(match)

					if len(submatch) > 1 {
						submatch[1] = strings.TrimPrefix(submatch[1], "\"")
						submatch[1] = strings.TrimPrefix(submatch[1], "'")
						submatch[1] = strings.TrimPrefix(submatch[1], "`")
						submatch[1] = strings.TrimSuffix(submatch[1], "\"")
						submatch[1] = strings.TrimSuffix(submatch[1], "'")
						submatch[1] = strings.TrimSuffix(submatch[1], "`")
						submatch[1] = strings.ReplaceAll(submatch[1], "\n", " ")
						submatch[1] = strings.ReplaceAll(submatch[1], "\t", "")

						rawInput := submatch[1] // this is the input before we process it

						// Skip duplicates
						_, found := dedupMap[rawInput]
						if !found {
							cssHash, _ := security.HighwayHash58(rawInput)
							cssHash = basic.GetFirstNChars(cssHash, 8)

							// Expand custom css macros (see comments below for details)
							submatch[1] = expandMe(submatch[1], cssHash)
							submatch[1] = expandMedia(submatch[1])
							submatch[1] = expandColor(submatch[1])
							submatch[1] = expandSpacing(submatch[1])

							matches = append(matches, submatch[1])
							dedupMap[rawInput] = true
						}
					}

					// We aren't actually replacing anything in the src file, we just needed to iterate over the regex matches
					return ""
				})
			}

			return nil
		})

		handleErr(err)
	}

	outputCSS := "/* " + METAGEN_AUTO_COMMENT + " */\n"
	outputCSS += strings.Join(matches, "")
	os.WriteFile("./wwwroot/css/style.metagen.css", []byte(outputCSS), 0664)

	printStatus(true)
}

// Expand spacing macro
// Ex:
//
//	padding: $5;
//
// => padding: calc(var(--spacing) * 5);
func expandSpacing(input string) string {
	re := regexp.MustCompile(`\$([0-9]+(?:\.[0-9]+)?)`)

	transformed := re.ReplaceAllStringFunc(input, func(match string) string {
		if len(re.FindStringSubmatch(match)) >= 2 {
			number := re.FindStringSubmatch(match)[1]
			return fmt.Sprintf("calc(var(--spacing) * %s)", number)
		}

		return ""
	})

	return transformed
}

// Expand color macros:
func expandColor(input string) string {
	re := regexp.MustCompile(`\$color\((.*?)(?:\/(0*(?:[1-9][0-9]?|100)))?\)`)

	transformed := re.ReplaceAllStringFunc(input, func(match string) string {
		if len(re.FindStringSubmatch(match)) == 3 {
			if re.FindStringSubmatch(match)[2] == "" {
				return fmt.Sprintf("var(--color-%s)", re.FindStringSubmatch(match)[1])
			} else {
				return fmt.Sprintf("oklch(from var(--color-%s) l c h / %s%%)", re.FindStringSubmatch(match)[1], re.FindStringSubmatch(match)[2])
			}
		}

		return ""
	})

	return transformed
}

// Expand shorthand media queries.
// Ex:
//
//	media md { ... }
//
// => media screen and (min-width: 768px) { ... }
func expandMedia(input string) string {
	input = strings.ReplaceAll(input, "$dark", "(prefers-color-scheme: dark)")
	input = strings.ReplaceAll(input, "$light", "(prefers-color-scheme: light)")
	input = strings.ReplaceAll(input, "$xs-", "screen and (max-width: 639px)")
	input = strings.ReplaceAll(input, "$sm-", "screen and (max-width: 767px)")
	input = strings.ReplaceAll(input, "$md-", "screen and (max-width: 1023px)")
	input = strings.ReplaceAll(input, "$lg-", "screen and (max-width: 1279px)")
	input = strings.ReplaceAll(input, "$xl-", "screen and (max-width: 1535px)")
	input = strings.ReplaceAll(input, "$sm", "screen and (min-width: 640px)")
	input = strings.ReplaceAll(input, "$md", "screen and (min-width: 768px)")
	input = strings.ReplaceAll(input, "$lg", "screen and (min-width: 1024px)")
	input = strings.ReplaceAll(input, "$xl", "screen and (min-width: 1280px)")
	input = strings.ReplaceAll(input, "$xx", "screen and (min-width: 1536px)")

	return input
}

// Expand "$me" macro with inline style attribute
// Ex:
//
//	$me { ... }
//
// => [__inlinecss_{REPLACEMENT_ID}] { ... }
func expandMe(input string, replacementId string) string {
	return strings.ReplaceAll(input, "$me", fmt.Sprintf("[__inlinecss_%s]", replacementId))
}
