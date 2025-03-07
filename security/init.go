package security

import "github.com/microcosm-cc/bluemonday"

var SanitizationPolicy *bluemonday.Policy

func Init() {
	SanitizationPolicy = bluemonday.UGCPolicy()
}
