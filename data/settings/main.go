package settings

import (
	"fmt"
	"time"
)

type Setting struct {
	read  bool
	raw   RawMessage
	value any
}

var settings map[string]Setting

type RawMessage struct {
	unmarshal func(any) error
}

func (msg *RawMessage) UnmarshalYAML(unmarshal func(any) error) error {
	msg.unmarshal = unmarshal
	return nil
}

func (msg *RawMessage) Unmarshal(v any) error {
	return msg.unmarshal(v)
}

func GetWithFallback[V comparable](key string, fallback V) V {
	if v, ok := settings[key]; ok {
		if v.read {
			return v.value.(V)
		}
		r := *new(V)
		v.raw.Unmarshal(&r)
		v.value = r
		v.read = true
		return v.value.(V)
	}
	return fallback
}

func Get[V comparable](key string) (V, error) {
	if v, ok := settings[key]; ok {
		if v.read {
			return v.value.(V), nil
		}
		r := *new(V)
		v.raw.Unmarshal(&r)
		v.value = r
		v.read = true
		return v.value.(V), nil
		/*if reflect.TypeOf(*new(V)) == reflect.TypeOf(v) {
			return v.(V), nil
		}*/
	}
	return *new(V), fmt.Errorf("missing key")
}

func Set(key string, value any) {
	if v, ok := settings[key]; ok {
		v.value = value
		v.read = true
		settings[key] = v
	} else {
		settings[key] = Setting{
			read:  true,
			value: value,
		}
	}
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
	settings = make(map[string]Setting)
	if err := load(); err != nil {
		panic(err)
	}
}
