package utils

import (
	"encoding/json"
	"fmt"
	"gogin/internal/app/myapp02/conf"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func HandleContext(v []byte, bc *conf.Bootstrap) error {
	mapValue := make(map[string]interface{})
	err := json.Unmarshal(v, &mapValue)
	if err != nil {
		panic(err)
	}
	newMap := convertMap(mapValue)
	data, err := marshalJSON(newMap)
	if err != nil {
		panic(err)
	}
	return unmarshalJSON(data, bc)
}

func marshalJSON(v interface{}) ([]byte, error) {
	if m, ok := v.(proto.Message); ok {
		return protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(m)
	}
	return json.Marshal(v)
}

func unmarshalJSON(data []byte, v interface{}) error {
	if m, ok := v.(proto.Message); ok {
		return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, m)
	}
	return json.Unmarshal(data, v)
}

func convertMap(src interface{}) interface{} {
	switch m := src.(type) {
	case map[string]interface{}:
		dst := make(map[string]interface{}, len(m))
		for k, v := range m {
			dst[k] = convertMap(v)
		}
		return dst
	case map[interface{}]interface{}:
		dst := make(map[string]interface{}, len(m))
		for k, v := range m {
			dst[fmt.Sprint(k)] = convertMap(v)
		}
		return dst
	case []interface{}:
		dst := make([]interface{}, len(m))
		for k, v := range m {
			dst[k] = convertMap(v)
		}
		return dst
	case []byte:
		// there will be no binary data in the config data
		return string(m)
	default:
		return src
	}
}
