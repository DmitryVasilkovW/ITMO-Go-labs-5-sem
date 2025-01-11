//go:build !solution

package structtags

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var fieldMapCache sync.Map

func Unpack(req *http.Request, ptr interface{}) error {
	if err := parseForm(req); err != nil {
		return err
	}

	fieldsMap, err := getFieldMap(ptr)
	if err != nil {
		return err
	}

	if err := populateFields(req.Form, fieldsMap, ptr); err != nil {
		return err
	}

	return nil
}

func parseForm(req *http.Request) error {
	return req.ParseForm()
}

func getFieldMap(ptr interface{}) (map[string]reflect.StructField, error) {
	t := reflect.TypeOf(ptr).Elem()

	fieldsMapInterface, ok := fieldMapCache.Load(t)
	if ok {
		return fieldsMapInterface.(map[string]reflect.StructField), nil
	}

	valueMap := make(map[string]reflect.StructField)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("http")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}
		valueMap[tag] = field
	}

	fieldMapCache.Store(t, valueMap)
	return valueMap, nil
}

func populateFields(formValues map[string][]string, fieldsMap map[string]reflect.StructField, ptr interface{}) error {
	v := reflect.ValueOf(ptr).Elem()

	for name, values := range formValues {
		field, ok := fieldsMap[name]
		if !ok {
			continue
		}

		fieldValue := v.FieldByIndex(field.Index)
		for _, value := range values {
			if err := populate(fieldValue, value); err != nil {
				return fmt.Errorf("%s: %v", name, err)
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		return populateString(v, value)
	case reflect.Int:
		return populateInt(v, value)
	case reflect.Bool:
		return populateBool(v, value)
	case reflect.Slice:
		return populateSlice(v, value)
	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
}

func populateString(v reflect.Value, value string) error {
	v.SetString(value)
	return nil
}

func populateInt(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	v.SetInt(i)
	return nil
}

func populateBool(v reflect.Value, value string) error {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}

func populateSlice(v reflect.Value, value string) error {
	elem := reflect.New(v.Type().Elem()).Elem()

	if err := populate(elem, value); err != nil {
		return err
	}

	v.Set(reflect.Append(v, elem))
	return nil
}
