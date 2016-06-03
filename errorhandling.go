package main

import (
	"encoding/json"
)

//create custom json
func CreateCustomJson(key, values []string) (j []byte) {
	var m map[string]string

	m = make(map[string]string)

	for i,_ := range key {
		m[key[i]] = values[i]
	}

	j, err := json.Marshal(m)

	if err != nil {
		panic("couldn't create the json")
	}
	return
}
