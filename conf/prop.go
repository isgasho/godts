package conf

import (
	"github.com/magiconair/properties"
	"path/filepath"
)

var isLoad = false

type Prop struct {
	FileNames []string
	context   *properties.Properties
}

func (prop *Prop) apply() {
	if !isLoad {
		absPath, _ := filepath.Abs("../godts/config/server.properties")
		// init from a file
		prop.context = properties.MustLoadFiles([]string{absPath}, properties.UTF8, true)
		isLoad = true
	}
}

func (prop *Prop) GetString(key string) string {
	prop.apply()
	return prop.context.MustGetString(key)
}

func (prop *Prop) GetInt(key string, def int) int {
	prop.apply()
	return prop.context.GetInt(key, def)
}
