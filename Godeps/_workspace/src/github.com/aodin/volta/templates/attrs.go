package templates

import (
	"encoding/json"
	"html/template"
)

type Attrs map[string]interface{}

func (a Attrs) Merge(b map[string]interface{}) {
	for key, value := range b {
		a[key] = value
	}
}

func AsJSON(key string, value interface{}) (attrs Attrs) {
	attrs = Attrs{}
	b, err := json.Marshal(value)
	if err != nil {
		// TODO panic on error?
		return
	}
	attrs[key] = template.JS(b)
	return
}
