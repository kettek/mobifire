package play

import "github.com/kettek/termfire/messages"

type Object struct {
	messages.ItemObject
	examineInfo  string
	container    *Object
	containerTag int32
}

func (o Object) IsContainer() bool {
	return o.Type >= 50 && o.Type <= 99
}

var objects = make(map[int32]*Object) // I guess it's okay to use a map.

func AddObject(io messages.ItemObject) *Object {
	obj := &Object{ItemObject: io}
	objects[io.Tag] = obj
	return obj
}

func GetObject(tag int32) *Object {
	return objects[tag]
}

func RemoveObject(tag int32) {
	delete(objects, tag)
}
