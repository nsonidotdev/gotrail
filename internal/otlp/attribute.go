package otlp

import (
	"encoding/base64"
	"reflect"
	"regexp"
	"strconv"
)

type attrValue interface {
	isAttrValue()
}

type attr struct {
	Key   string    `json:"key"`
	Value attrValue `json:"value"`
}

type attrEmptyValue struct{}

type attrStringValue struct {
	StringValue string `json:"stringValue"`
}

type attrIntValue struct {
	IntValue string `json:"intValue"`
}

type attrDoubleValue struct {
	DoubleValue float64 `json:"doubleValue"`
}

type attrBoolValue struct {
	BoolValue bool `json:"boolValue"`
}

type attrBytesValue struct {
	// Base64 encoded string
	BytesValue string `json:"bytesValue"`
}

type attrArrayValue struct {
	ArrayValue attrArrayValueInner `json:"arrayValue"`
}

type attrKVListValue struct {
	KVListValue attrKVListValueInner `json:"kvlistValue"`
}

type attrArrayValueInner struct {
	Values []attrValue `json:"values"`
}

type attrKVListValueInner struct {
	Values []attrKVListValueEntry `json:"values"`
}

type attrKVListValueEntry struct {
	Key   string    `json:"key"`
	Value attrValue `json:"value"`
}

const otlpAttrStructKey = "gtattr"

func (attrStringValue) isAttrValue() {}
func (attrIntValue) isAttrValue()    {}
func (attrDoubleValue) isAttrValue() {}
func (attrBoolValue) isAttrValue()   {}
func (attrBytesValue) isAttrValue()  {}
func (attrArrayValue) isAttrValue()  {}
func (attrKVListValue) isAttrValue() {}
func (attrEmptyValue) isAttrValue()  {}

func serializeAttributes(attributes map[string]any) []*attr {
	result := make([]*attr, 0, len(attributes))
	for key, value := range attributes {
		if !isKeyValid(key) {
			continue
		}

		parsed := serializeAttributeValue(value)

		if _, ok := parsed.(attrEmptyValue); !ok {
			attribute := &attr{Key: key, Value: parsed}
			result = append(result, attribute)
		}
	}
	return result
}

func serializeAttributeValue(value any) attrValue {
	if value == nil {
		return attrEmptyValue{}
	}

	typeof := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	kind := typeof.Kind()

	if kind == reflect.Pointer {
		if v.IsNil() {
			return attrEmptyValue{}
		}

		v = v.Elem()
	}

	switch kind {
	case reflect.Invalid:
		return attrEmptyValue{}

	// String
	case reflect.String:
		stringValue := v.String()
		return attrStringValue{StringValue: stringValue}

	// Boolean
	case reflect.Bool:
		boolValue := v.Bool()
		return attrBoolValue{BoolValue: boolValue}

	// Integer
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue := v.Int()
		return attrIntValue{IntValue: strconv.FormatInt(intValue, 10)}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		intValue := int64(v.Uint())
		return attrIntValue{IntValue: strconv.FormatInt(intValue, 10)}

	// Double
	case reflect.Float32, reflect.Float64:
		doubleValue := v.Float()
		return attrDoubleValue{DoubleValue: doubleValue}

	// Array and Bytes
	case reflect.Slice, reflect.Array:
		// Bytes
		if v.Type().Elem().Kind() == reflect.Uint8 {
			bytes := v.Interface().([]byte)
			if len(bytes) == 0 {
				return attrEmptyValue{}
			}

			return attrBytesValue{
				BytesValue: base64.StdEncoding.EncodeToString(bytes),
			}
		}

		// Array
		arrayValue := make([]attrValue, 0, v.Len())
		for i := range v.Len() {
			rawValue := v.Index(i).Interface()
			value := serializeAttributeValue(rawValue)
			if isEmptyAttr(value) {
				continue
			}
			arrayValue = append(arrayValue, value)
		}

		if len(arrayValue) == 0 {
			return attrEmptyValue{}
		}

		if hasMixedAttrs(arrayValue) {
			for i, attr := range arrayValue {
				if !isPrimitiveAttr(attr) {
					continue
				}

				arrayValue[i] = transformToStringValue(attr)
			}
		}

		return attrArrayValue{ArrayValue: attrArrayValueInner{Values: arrayValue}}

	// KV
	case reflect.Map:
		entries := make([]attrKVListValueEntry, 0, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			rawKey := iter.Key()
			if rawKey.Kind() != reflect.String {
				continue
			}

			key := rawKey.String()
			if !isKeyValid(key) {
				continue
			}

			rawValue := iter.Value().Interface()
			value := serializeAttributeValue(rawValue)
			if isEmptyAttr(value) {
				continue
			}

			entries = append(entries, attrKVListValueEntry{Key: key, Value: value})
		}

		if len(entries) == 0 {
			return attrEmptyValue{}
		}

		return attrKVListValue{KVListValue: attrKVListValueInner{Values: entries}}

	case reflect.Struct:
		entries := make([]attrKVListValueEntry, 0, v.NumField())
		for i := range v.NumField() {
			tagValue, ok := typeof.Field(i).Tag.Lookup(otlpAttrStructKey)
			if !ok {
				continue
			}

			if !isKeyValid(tagValue) {
				continue
			}

			rawValue := v.Field(i).Interface()
			value := serializeAttributeValue(rawValue)

			entries = append(entries, attrKVListValueEntry{Key: tagValue, Value: value})
		}

		if len(entries) == 0 {
			return attrEmptyValue{}
		}

		return attrKVListValue{KVListValue: attrKVListValueInner{Values: entries}}
	}

	return attrEmptyValue{}
}

func isEmptyAttr(attr attrValue) bool {
	_, ok := attr.(attrEmptyValue)
	return ok
}

func isPrimitiveAttr(value attrValue) bool {
	switch value.(type) {
	case attrStringValue, attrIntValue, attrDoubleValue, attrBoolValue:
		return true
	}

	return false
}

func hasMixedAttrs(attrs []attrValue) bool {
	if len(attrs) <= 1 {
		return false
	}

	var firstPrimtiveAttr reflect.Type
	for _, attr := range attrs {
		if !isPrimitiveAttr(attr) {
			continue
		}

		currentType := reflect.TypeOf(attr)

		if firstPrimtiveAttr == nil {
			firstPrimtiveAttr = currentType
			continue
		}

		if currentType != firstPrimtiveAttr {
			return true
		}
	}

	return false
}

func transformToStringValue(attr attrValue) attrStringValue {
	switch v := attr.(type) {
	case attrStringValue:
		return v
	case attrIntValue:
		return attrStringValue{
			StringValue: v.IntValue,
		}
	case attrDoubleValue:
		return attrStringValue{
			StringValue: strconv.FormatFloat(v.DoubleValue, 'f', -1, 64),
		}
	case attrBoolValue:
		return attrStringValue{
			StringValue: strconv.FormatBool(v.BoolValue),
		}
	}
	return attrStringValue{StringValue: ""}
}

func isKeyValid(key string) bool {
	pattern := `^[a-zA-Z][a-zA-Z0-9_.\-]*$`
	matched, err := regexp.MatchString(pattern, key)
	if err != nil {
		return false
	}

	return matched
}
