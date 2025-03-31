package etc

import "encoding/json"

func MarshalJSON(v interface{}) []byte {
	jsonByte, _ := json.Marshal(v)

	return jsonByte
}
