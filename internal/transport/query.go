package transport

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// BuildQuery converts a struct (with json tags), [url.Values], or map into [url.Values].
func BuildQuery(params any) (url.Values, error) {
	if params == nil {
		return url.Values{}, nil
	}
	if values, ok := params.(url.Values); ok {
		return cloneValues(values), nil
	}
	if values, ok := params.(*url.Values); ok {
		if values == nil {
			return url.Values{}, nil
		}
		return cloneValues(*values), nil
	}
	switch typed := params.(type) {
	case map[string]string:
		return valuesFromStringMap(typed), nil
	case map[string][]string:
		return valuesFromSliceMap(typed), nil
	case *map[string]string:
		if typed == nil {
			return url.Values{}, nil
		}
		return valuesFromStringMap(*typed), nil
	case *map[string][]string:
		if typed == nil {
			return url.Values{}, nil
		}
		return valuesFromSliceMap(*typed), nil
	}
	return valuesFromStruct(params)
}

func cloneValues(values url.Values) url.Values {
	clone := url.Values{}
	for key, items := range values {
		copied := make([]string, len(items))
		copy(copied, items)
		clone[key] = copied
	}
	return clone
}

func valuesFromStringMap(values map[string]string) url.Values {
	query := url.Values{}
	keys := sortedKeys(values)
	for _, key := range keys {
		query.Set(key, values[key])
	}
	return query
}

func valuesFromSliceMap(values map[string][]string) url.Values {
	query := url.Values{}
	keys := sortedKeys(values)
	for _, key := range keys {
		items := values[key]
		for _, item := range items {
			query.Add(key, item)
		}
	}
	return query
}

func sortedKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func valuesFromStruct(params any) (url.Values, error) {
	value := reflect.ValueOf(params)
	if !value.IsValid() {
		return url.Values{}, nil
	}
	value = reflect.Indirect(value)
	if !value.IsValid() {
		return url.Values{}, nil
	}
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("query params must be a struct, map, or url.Values, got %s", value.Kind())
	}

	values := url.Values{}
	valueType := value.Type()
	for i := range value.NumField() {
		field := valueType.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tag, tagOptions := parseJSONTag(field.Tag.Get("json"))
		if tag == "-" {
			continue
		}
		if tag == "" {
			tag = field.Name
		}

		fieldValue := value.Field(i)
		if tagOptions.omitEmpty && isEmptyValue(fieldValue) {
			continue
		}
		addQueryValue(values, tag, fieldValue)
	}

	return values, nil
}

type jsonTagOptions struct {
	omitEmpty bool
}

func parseJSONTag(tag string) (string, jsonTagOptions) {
	if tag == "" {
		return "", jsonTagOptions{}
	}
	parts := strings.Split(tag, ",")
	name := parts[0]
	options := jsonTagOptions{}
	for _, part := range parts[1:] {
		if part == "omitempty" {
			options.omitEmpty = true
		}
	}
	return name, options
}

func isEmptyValue(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}
	switch value.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return value.Complex() == 0
	case reflect.Interface, reflect.Pointer:
		return value.IsNil()
	case reflect.Struct:
		return value.IsZero()
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return value.IsNil()
	default:
		return false
	}
}

func addQueryValue(values url.Values, key string, value reflect.Value) {
	if !value.IsValid() {
		return
	}
	value = reflect.Indirect(value)
	if !value.IsValid() {
		return
	}

	switch value.Kind() {
	case reflect.Invalid:
		return
	case reflect.Slice, reflect.Array:
		for i := range value.Len() {
			item := value.Index(i)
			addQueryValue(values, key, item)
		}
	case reflect.Map:
		if value.Type().Key().Kind() != reflect.String {
			return
		}
		keys := value.MapKeys()
		sorted := make([]string, 0, len(keys))
		for _, keyValue := range keys {
			sorted = append(sorted, keyValue.String())
		}
		sort.Strings(sorted)
		for _, mapKey := range sorted {
			item := value.MapIndex(reflect.ValueOf(mapKey))
			addQueryValue(values, key, item)
		}
	case reflect.Struct:
		values.Add(key, fmt.Sprintf("%v", value.Interface()))
	case reflect.String:
		values.Add(key, value.String())
	case reflect.Bool:
		values.Add(key, strconv.FormatBool(value.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		values.Add(key, strconv.FormatInt(value.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		values.Add(key, strconv.FormatUint(value.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		values.Add(key, fmt.Sprintf("%v", value.Float()))
	case reflect.Complex64, reflect.Complex128:
		values.Add(key, fmt.Sprintf("%v", value.Complex()))
	case reflect.Interface, reflect.Pointer:
		addQueryValue(values, key, value)
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		values.Add(key, fmt.Sprintf("%v", value.Interface()))
	default:
		values.Add(key, fmt.Sprintf("%v", value.Interface()))
	}
}
