package binding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type jsonBinding struct {
	DisallowUnknownFields bool
	IsValidate            bool
}

func (j *jsonBinding) Name() string {
	return "json"
}

func (j *jsonBinding) Bind(r *http.Request, obj any) error {
	body := r.Body
	if body == nil {
		return errors.New("invalid request")
	}
	decoder := json.NewDecoder(body)
	if j.DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if j.IsValidate {
		err := validateParam(obj, decoder)
		if err != nil {
			return err
		}
	} else {
		err := decoder.Decode(obj)
		if err != nil {
			return err
		}
	}
	return validate(obj)
}

func validateParam(obj any, decoder *json.Decoder) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer {
		return errors.New("this argument must have a pointer type")
	}
	elem := v.Elem().Interface()
	of := reflect.ValueOf(elem)

	switch of.Kind() {
	case reflect.Struct:
		return checkParam(of, obj, decoder)
	case reflect.Slice, reflect.Array:
		elem := of.Type().Elem()
		if elem.Kind() == reflect.Struct {
			return checkParamSlice(elem, obj, decoder)
		}
	default:
		decoder.Decode(obj)
	}

	return nil
}

func checkParamSlice(of reflect.Type, obj any, decoder *json.Decoder) error {
	mapValue := make([]map[string]any, 0)
	_ = decoder.Decode(&mapValue)
	for i := 0; i < of.NumField(); i++ {
		field := of.Field(i)
		name := field.Name
		jsonName := field.Tag.Get("json")
		if jsonName != "" {
			name = jsonName
		}
		required := field.Tag.Get("sonata")
		for _, v := range mapValue {
			value := v[name]
			if value == nil && required == "required" {
				return fmt.Errorf("field [%s] is not exist", jsonName)
			}
		}

		println(name)
	}
	b, _ := json.Marshal(mapValue)
	_ = json.Unmarshal(b, obj)
	return nil
}

func checkParam(of reflect.Value, obj any, decoder *json.Decoder) error {
	mapValue := make(map[string]any)
	_ = decoder.Decode(&mapValue)
	for i := 0; i < of.NumField(); i++ {
		field := of.Type().Field(i)
		name := field.Name
		jsonName := field.Tag.Get("json")
		if jsonName != "" {
			name = jsonName
		}
		required := field.Tag.Get("sonata")
		value := mapValue[name]
		if value == nil && required == "required" {
			return fmt.Errorf("field [%s] is not exist", jsonName)
		}
		println(name)
	}
	b, _ := json.Marshal(mapValue)
	_ = json.Unmarshal(b, obj)
	return nil
}
