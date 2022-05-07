package weapp

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type container struct {
	constructor string
	sync.Map
}

func newContainer() *container {
	container := &container{}
	container.constructor = "Init"
	return container
}

func (container *container) SetConstructor(funcName string) {
	container.constructor = funcName
}

func (container *container) Store(key string, instance any, construct func(instance any)) error {
	refType := reflect.TypeOf(instance)
	refValue := reflect.ValueOf(instance)
	isStruct := refType.Kind() == reflect.Struct
	isStructPtr := refType.Kind() == reflect.Ptr && refType.Elem().Kind() == reflect.Struct
	if !isStruct && !isStructPtr {
		return errors.New("instance is not struct or pointer of struct")
	}
	initMethod, initExists := reflect.TypeOf(instance).MethodByName(container.constructor)
	if isStructPtr {
		refType = refType.Elem()
		refValue = refValue.Elem()
	}
	fieldCount := refType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := refType.Field(i)
		if tag, ok := field.Tag.Lookup("inject"); ok {
			fieldTypeName := field.Type.Name()
			fieldValue := refValue.Field(i)
			if fieldValue.CanSet() {
				if (field.Type.Kind() == reflect.Interface || field.Type.Kind() == reflect.Ptr) && !fieldValue.IsNil() {
					continue
				} else if !reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(field.Type).Interface()) {
					continue
				}
				var fieldInstance reflect.Value
				if field.Type.Kind() == reflect.Map {
					fieldInstance = reflect.MakeMap(field.Type)
				} else if field.Type.Kind() == reflect.Slice {
					fieldInstance = reflect.MakeSlice(field.Type, 0, 0)
				} else if field.Type.Kind() == reflect.Struct {
					fieldInstance = reflect.Indirect(reflect.New(field.Type))
				} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
					if tag != "" {
						if object, err := container.LoadOrStore(tag, reflect.New(field.Type.Elem()).Interface()); err != nil {
							return err
						} else {
							fieldInstance = reflect.ValueOf(object)
						}
					} else {
						fieldValue := reflect.New(field.Type.Elem()).Interface()
						if err := container.Store(fieldTypeName, fieldValue, nil); err == nil {
							fieldInstance = reflect.ValueOf(fieldValue)
						} else {
							return err
						}
					}
				} else if field.Type.Kind() == reflect.Interface {
					if tag != "" {
						if object, ok := container.Load(tag); ok {
							if reflect.TypeOf(object).Implements(field.Type) {
								fieldInstance = reflect.ValueOf(object)
							} else {
								return fmt.Errorf("instance `%s` is not implements interface `%s`", tag, field.Type.Name())
							}
						} else {
							return fmt.Errorf("instance `%s` is not provided", tag)
						}
					} else {
						found := false
						container.Map.Range(func(key any, value any) bool {
							if reflect.TypeOf(value).Implements(field.Type) {
								fieldInstance = reflect.ValueOf(value)
								found = true
								return false
							}
							return true
						})
						if !found {
							return fmt.Errorf("no provider is implements interface `%s`", field.Type.Name())
						}
					}
				} else {
					continue
				}
				fieldValue.Set(fieldInstance)
			} else {
				return errors.New("inject can only work on an exported field")
			}
		}
	}
	if construct != nil {
		construct(instance)
	} else if initExists {
		initMethod.Func.Call([]reflect.Value{refValue.Addr()})
	}
	container.Map.Store(key, instance)
	return nil
}

func (container *container) Load(key string) (any, bool) {
	return container.Map.Load(key)
}

func (container *container) LoadOrStore(key string, instance any) (any, error) {
	object, ok := container.Load(key)
	if ok {
		return object, nil
	}
	if err := container.Store(key, instance, nil); err != nil {
		return nil, err
	}
	return instance, nil
}
