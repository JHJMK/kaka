package handler

import (
	"os"
	"text/template"
)

func GetInstallationPackageBytes(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

var funcMap = template.FuncMap{
	"default": func(defaultValue interface{}, value interface{}) interface{} {
		if value == nil || value == "" {
			return defaultValue
		}
		return value
	},
}
