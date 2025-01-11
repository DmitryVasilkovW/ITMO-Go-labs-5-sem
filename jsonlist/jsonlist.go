//go:build !solution

package jsonlist

import (
	"encoding/json"
	"io"
	"reflect"
)

func ValidateSlice(slice interface{}) (reflect.Value, error) {
	sliceVal := reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		return reflect.Value{}, &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}
	return sliceVal, nil
}

func MarshalElement(elem interface{}) ([]byte, error) {
	return json.Marshal(elem)
}

func WriteSeparator(w io.Writer, isFirst bool) error {
	if !isFirst {
		if _, err := w.Write([]byte(" ")); err != nil {
			return err
		}
	}
	return nil
}

func WriteBytes(w io.Writer, bytes []byte) error {
	_, err := w.Write(bytes)
	return err
}

func Marshal(w io.Writer, slice interface{}) error {
	sliceVal, err := ValidateSlice(slice)
	if err != nil {
		return err
	}

	length := sliceVal.Len()
	for i := 0; i < length; i++ {
		elem := sliceVal.Index(i).Interface()
		bytes, err := MarshalElement(elem)
		if err != nil {
			return err
		}

		if err := WriteSeparator(w, i == 0); err != nil {
			return err
		}
		if err := WriteBytes(w, bytes); err != nil {
			return err
		}
	}
	return nil
}

func ValidateSlicePointer(slicePointer interface{}) (reflect.Value, reflect.Type, error) {
	slicePtrVal := reflect.ValueOf(slicePointer)
	if slicePtrVal.Kind() != reflect.Ptr || slicePtrVal.Elem().Kind() != reflect.Slice {
		return reflect.Value{}, nil, &json.UnsupportedTypeError{Type: reflect.TypeOf(slicePointer)}
	}

	sliceVal := slicePtrVal.Elem()
	elemType := sliceVal.Type().Elem()
	return sliceVal, elemType, nil
}

func DecodeElement(decoder *json.Decoder, elemType reflect.Type) (reflect.Value, error) {
	elemPtr := reflect.New(elemType).Interface()
	if err := decoder.Decode(elemPtr); err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(elemPtr).Elem(), nil
}

func Unmarshal(r io.Reader, slicePointer interface{}) error {
	sliceVal, elemType, err := ValidateSlicePointer(slicePointer)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(r)
	for {
		elem, err := DecodeElement(decoder, elemType)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, elem))
	}

	return nil
}
