package prop

import (
	"github.com/magiconair/properties"
	"path/filepath"
)

var isLoad = false
var _prop *properties.Properties

func apply() {
	if !isLoad {
		absPath, _ := filepath.Abs("../godts/config/server.properties")
		// init from a file
		_prop = properties.MustLoadFiles([]string{absPath}, properties.UTF8, true)
		isLoad = true
	}
}

func GetString(key string) string {
	apply()
	return _prop.MustGetString(key)
}

func GetInt(key string, def int) int {
	apply()
	return _prop.GetInt(key, def)
}
