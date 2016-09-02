package jsonParse

import (
	"encoding/json"
	"strings"
)

func JsonNewDecoder(body []byte) *json.Decoder {
	dec := json.NewDecoder(strings.NewReader(string(body)))
	return dec

}

func JsonUnmarsha(body []byte, v interface{}) {
	var err error
	if err = json.Unmarshal(body, &v); err != nil {
	}
}

func JsonMarsha(v interface{}) []byte {
	b, err := json.Marshal(&v)
	if err != nil {
	}
	return b
}


