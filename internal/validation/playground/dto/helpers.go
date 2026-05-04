package dto

import "reflect"

func copyStruct(src, dst interface{}) {
	if dst == nil || src == nil {
		panic("dst or src is nil")
	}

	if reflect.ValueOf(dst).Kind() != reflect.Pointer {
		panic("dst must be a pointer")
	}

	vSrc := reflect.ValueOf(src)

	if vSrc.Kind() != reflect.Struct {
		panic("src must be a struct")
	}

	vDst := reflect.ValueOf(dst).Elem()

	if vDst.Kind() != reflect.Struct {
		panic("dst must be a struct")
	}

	for i := 0; i < vDst.NumField(); i++ {
		dstField := vDst.Field(i)

		if !dstField.CanSet() {
			continue
		}

		srcField := vSrc.FieldByName(vDst.Type().Field(i).Name)

		if !srcField.IsValid() {
			continue
		}

		if dstField.Type() == srcField.Type() {
			dstField.Set(srcField)
		} else if srcField.Type().ConvertibleTo(dstField.Type()) {
			dstField.Set(srcField.Convert(dstField.Type()))
		} else if dstField.Kind() == reflect.Struct && srcField.Kind() == reflect.Struct {
			copyStruct(srcField.Interface(), dstField.Addr().Interface())
		}
	}
}
