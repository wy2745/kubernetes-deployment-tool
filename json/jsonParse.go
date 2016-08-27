package jsonParse

import (
	"encoding/json"
	"strings"
	"fmt"
)

func JsonNewDecoder(body []byte) *json.Decoder {
	dec := json.NewDecoder(strings.NewReader(string(body)))
	return dec

}

func JsonUnmarsha(body []byte, v interface{}) {
	if err := json.Unmarshal(body, &v); err != nil {
		fmt.Println(err)
	}
}


