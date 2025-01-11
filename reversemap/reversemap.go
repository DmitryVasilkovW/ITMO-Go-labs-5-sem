//go:build !solution

package reversemap

import (
	"reflect"
)

func ValidateMap(forward interface{}) reflect.Value {
	forwardValue := reflect.ValueOf(forward)

	if forwardValue.Kind() != reflect.Map {
		panic("🐓!")
	}

	forwardType := forwardValue.Type()
	if forwardType.Key().Kind() != reflect.String && forwardType.Key().Kind() != reflect.Int {
		panic("🐓!!")
	}
	if forwardType.Elem().Kind() != reflect.String && forwardType.Elem().Kind() != reflect.Int {
		panic("🐓!!!")
	}

	return forwardValue
}

func CreateReverseMapType(forwardType reflect.Type) reflect.Type {
	return reflect.MapOf(forwardType.Elem(), forwardType.Key())
}

func PopulateReverseMap(forwardValue reflect.Value, reverse reflect.Value) {
	for _, key := range forwardValue.MapKeys() {
		reverse.SetMapIndex(forwardValue.MapIndex(key), key)
	}
}

func ReverseMap(forward interface{}) interface{} {
	forwardValue := ValidateMap(forward)
	forwardType := forwardValue.Type()

	reverseType := CreateReverseMapType(forwardType)
	reverse := reflect.MakeMap(reverseType)

	PopulateReverseMap(forwardValue, reverse)
	return reverse.Interface()
}
