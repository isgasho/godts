package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var isInit = false

type Yml struct {
	context *YmlContext
}

func (yml *Yml) apply() {
	if !isInit {
		buf, err := ioutil.ReadFile("../godts/config/application.yml")
		if err != nil {
			panic(err)
		}

		var context YmlContext

		err = yaml.Unmarshal(buf, &context)
		if err != nil {
			panic(err)
		}
		log.Printf("d: %+v", context)

		yml.context = &context
		isInit = true
	}
}

func (yml *Yml) GetYmlContext() *YmlContext {
	yml.apply()
	return yml.context
}
