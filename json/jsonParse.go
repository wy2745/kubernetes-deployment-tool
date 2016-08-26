package jsonParse

import (
	"encoding/json"
	"strings"
	"fmt"
	classType "../type"
)
func JsonNewDecoder(body []byte) *json.Decoder{
	dec := json.NewDecoder(strings.NewReader(string(body)))
	return dec

}

func JsonUnmarsha(body []byte) classType.NodeList{
	var v classType.NodeList
	if err := json.Unmarshal(body, &v); err != nil {
		fmt.Println(err)
		return v
	}
	return v
}


