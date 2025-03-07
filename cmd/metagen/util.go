package main

import (
	"fmt"
	"go/ast"
	. "maxwarden/basic"
	"net/url"
	"os"
	"reflect"
	"regexp"
)

func printStatus(b bool) {
	var status string

	if b {
		status = "SUCCESS"
	} else {
		status = "FAILED"
	}

	fmt.Printf("... %s\n", status)
}

func handleCmdOutput(out []byte, err error) {
	if err != nil {
		fmt.Printf("\n%s\n", out)
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}

func handleErr(err error) {
	if err != nil {
		printStatus(false)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// convert an sqlite DSN to a file name
func parseSQLiteFilename(dsn string) (string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	// For simple file paths
	if u.Scheme == "" {
		return u.Path, nil
	}

	// For more complex DSNs
	if u.Scheme == "file" {
		return u.Opaque, nil
	}

	return "", fmt.Errorf("invalid DSN format: %s", dsn)
}

// given an ast.Decl, and destination struct, look at the struct for any
// boolean fields with the struct tag `Note`. If valid notes are found in the given Decl doc string
// set the tagged booleans to true on the input struct.
// file is the file that the decl originated from.
func parseNotesFromDocComment(decl ast.Decl, file *os.File, dest any) error {
	re := regexp.MustCompile(`@(\w+)`)

	var identifier string
	var docstring string

	var docNotes []string
	var validNotes []string

	// check if decl is a function
	if funcDecl, ok := decl.(*ast.FuncDecl); ok {
		identifier = funcDecl.Name.Name
		docstring = funcDecl.Doc.Text()
	}

	matches := re.FindAllStringSubmatch(docstring, -1)

	for _, match := range matches {
		// match[1] contains the first capture group (the word after '@')
		docNotes = append(docNotes, match[1])
	}

	v := reflect.ValueOf(dest).Elem()
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("\ndestination struct must be a pointer to a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		// Look for the `Note` tag in the struct field
		if field.Tag == "note:\"true\"" {
			if Contains(docNotes, field.Name) {
				v.Field(i).SetBool(true)
			}

			validNotes = append(validNotes, field.Name)
		}
	}

	for _, v := range docNotes {
		if !Contains(validNotes, v) {
			return fmt.Errorf("\n`%s`: Unknown note `@%s`, Identifier: `%s`\n\tValid values are: %v", file.Name(), v, identifier, validNotes)
		}
	}

	return nil
}
