// General purpose "utilities" that act as my own "standard library"
package basic

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

func GetPathParts(path string) []string {
	trimmed := strings.TrimPrefix(path, "/")
	return strings.Split(trimmed, "/")
}

// Takes a tree pointer and a slice of path segments to insert.
// It looks for an existing child with the current segment name; if none is found,
// it creates a new node on the tree. Then it recurses on the remaining segments.
//
// Ex: Generate tree nodes from url segments
// `/app/examples/webpage`     -> {"app", "examples", "webpage"}
// `/app/examples/hello-world` -> {"app", "examples", "hello-world"
// `/auth/login`               -> {"auth", "login"}
//
// =>
//
//	root {
//	    app {
//	        examples {
//	            webpage
//	            hello-world
//	        }
//	    }
//
//	    auth {
//	        login
//	    }
//	}
//
// This function is used in the code generation process to generate `pageinfo` structs
// from all known application page URLs.
func AddStringPartsToTree(tree *Tree, parts []string) {
	if len(parts) == 0 {
		return
	}

	for i, _ := range parts {
		if parts[i] == "" {
			parts[i] = "index"
		}
	}

	if tree.Children == nil {
		tree.Children = new([]Tree)
	}

	// Search for an existing child with the current part's name.
	var child *Tree
	for i := range *tree.Children {
		if (*tree.Children)[i].Name == parts[0] {
			child = &((*tree.Children)[i])
			break
		}
	}

	// If no child is found, create a new one and append it.
	if child == nil {
		newNode := Tree{Name: parts[0]}
		*tree.Children = append(*tree.Children, newNode)
		child = &((*tree.Children)[len(*tree.Children)-1])
	}

	AddStringPartsToTree(child, parts[1:])
}

type Tree struct {
	Name     string
	Children *[]Tree
}

func CapitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	// Convert the first rune to uppercase
	first := []rune(s)[0]
	return string(unicode.ToUpper(first)) + s[1:]
}

func IntAbs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func MakeURLParams(base string, params ...[2]string) string {
	output := base

	for i, v := range params {
		if i == 0 {
			output += "?" + v[0] + "=" + v[1]
		} else {
			output += "&" + v[0] + "=" + v[1]
		}
	}

	return output
}

func SnakeCaseToTitleCase(s string) string {
	parts := strings.Split(s, "_")

	for i, part := range parts {
		parts[i] = CapitalizeFirstLetter(part)
	}

	return strings.Join(parts, " ")
}

func ToString(i interface{}) string {
	v := reflect.ValueOf(i)

	output := ""

	switch v.Kind() {
	case reflect.String:
		output = v.String()
	case reflect.Int:
		output = fmt.Sprintf("%d", v.Int())
	case reflect.Float64:
		output = fmt.Sprintf("%f", v.Float())
	case reflect.Bool:
		output = fmt.Sprintf("%t", v.Bool())
	default:
		output = fmt.Sprintf("%v", i)
	}

	return output
}

func HTMLDateToTime(date string) time.Time {
	t, _ := time.Parse("2006-01-02", date)
	return t
}

func TimeToSqliteString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func SqliteStringToTime(dateTimeStr string) time.Time {
	formats := []string{
		"2006-01-02 15:04:05.000000000",
		"2006-01-02 15:04:05",
	}

	var t time.Time
	var err error

	for _, format := range formats {
		t, err = time.Parse(format, dateTimeStr)
		if err == nil {
			return t
		}
	}

	return t
}

func TimeToTimeString(utcTime time.Time) string {
	return utcTime.Format("03:04 PM")
}

func TimeToString(utcTime time.Time) string {
	return utcTime.Format("01/02/06 03:04 PM")
}

func DateToString(utcTime time.Time) string {
	return utcTime.Format("01/02/06")
}

func StringToDate(ds string) time.Time {
	nt, err := time.Parse("01/02/06", ds)
	if err != nil {
		nt, _ = time.Parse("2006-01-02", ds)
	}
	return nt
}

func Reverse[T comparable](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func IndexOf[T comparable](collection []T, el T) int {
	for i, x := range collection {
		if x == el {
			return i
		}
	}
	return -1
}

func RemoveDuplicates[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GetFirstNChars(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}
