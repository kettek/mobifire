package settings

import (
	"fmt"
	"reflect"
	"time"
)

var settings map[string]any

func GetWithFallback[V comparable](key string, fallback V) V {
	if v, ok := settings[key]; ok {
		if reflect.TypeOf(fallback) == reflect.TypeOf(v) {
			return v.(V)
		}
	}
	return fallback
}

func Get[V comparable](key string) (V, error) {
	if v, ok := settings[key]; ok {
		if reflect.TypeOf(*new(V)) == reflect.TypeOf(v) {
			return v.(V), nil
		}
	}
	return *new(V), fmt.Errorf("missing key")
}

func Set(key string, value any) {
	settings[key] = value
	Save()
}

func Save() {
	if cancelSave != nil {
		cancelSave <- struct{}{}
		cancelSave = nil
	}
	cancelSave = make(chan struct{}, 1)
	go pendingSave()
}

var cancelSave chan struct{}

func pendingSave() {
	select {
	case <-cancelSave:
		return
	case <-time.After(1 * time.Second):
	}
	if err := save(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	settings = make(map[string]any)
	if err := load(); err != nil {
		panic(err)
	}
}
