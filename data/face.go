package data

import (
	"bytes"
	"image"
	_ "image/png"

	"github.com/kettek/termfire/messages"
)

type FaceSet struct {
	Set    int
	Width  int
	Height int
	Anims  map[int]*Anim
	Faces  map[int]*FaceImage
	Names  map[string]int
}

var faceSets = make(map[int]FaceSet)

func AddFaceSet(set int, width, height int) {
	faceSets[set] = FaceSet{
		Set:    set,
		Width:  width,
		Height: height,
		Anims:  make(map[int]*Anim),
		Faces:  make(map[int]*FaceImage),
		Names:  make(map[string]int),
	}
}

func GetFaceSet(set int, num int) (*FaceImage, bool) {
	faceSet, ok := faceSets[set]
	if !ok {
		return nil, false
	}
	if num < 0 || num >= len(faceSet.Faces) {
		return nil, false
	}
	return faceSet.Faces[num], true
}

var anims = make(map[int]*Anim)
var faces = make(map[int]*FaceImage)
var names = make(map[string]int)

var currentFaceSet int

func SetCurrentFaceSet(set int) {
	if _, ok := faceSets[set]; !ok {
		return
	}
	faces = faceSets[set].Faces
	names = faceSets[set].Names
	anims = faceSets[set].Anims
	currentFaceSet = set
}

func CurrentFaceSet() FaceSet {
	return faceSets[currentFaceSet]
}

// FaceImage is a merger of face and image, cuz why not.
type FaceImage struct {
	Num      uint16
	Set      int8
	Width    int
	Height   int
	Data     []byte
	Image    image.Image
	name     string
	Checksum int32
	pending  bool
}

// Name returns the face name. Provides fyne.Resource interface.
func (f *FaceImage) Name() string {
	return f.name
}

// Content returns the image data. Provides fyne.Resource interface.
func (f *FaceImage) Content() []byte {
	return f.Data
}

// GetFace returns a face from the face map.
func GetFace(num int) (*FaceImage, bool) {
	face, ok := faces[num]
	return face, ok
}

// AddFace adds a pending face to the face map. Returns if it exists.
func AddFace(face messages.MessageFace2) bool {
	_, ok := faces[int(face.Num)]
	if ok {
		return true
	}
	faces[int(face.Num)] = &FaceImage{
		Num:      uint16(face.Num),
		Set:      int8(face.SetNum),
		name:     face.Name,
		Checksum: face.Checksum,
		pending:  true,
	}
	names[face.Name] = int(face.Num)
	return false
}

// AddFaceImage adds an image to the face map.
func AddFaceImage(msg messages.MessageImage2) {
	face, ok := faces[int(msg.Face)]
	b := bytes.NewReader(msg.Data)
	img, _, err := image.Decode(b)
	if err != nil {
		panic(err)
	}
	if !ok {
		faces[int(msg.Face)] = &FaceImage{
			Num:     uint16(msg.Face),
			Set:     msg.Set,
			Width:   msg.Width,
			Height:  msg.Height,
			Data:    msg.Data,
			Image:   img,
			pending: false,
		}
		return
	}
	faces[int(msg.Face)] = &FaceImage{
		Num:      uint16(msg.Face),
		Set:      msg.Set,
		Width:    msg.Width,
		Height:   msg.Height,
		Data:     msg.Data,
		Image:    img,
		name:     face.name,
		Checksum: face.Checksum,
		pending:  false,
	}
	names[face.name] = int(msg.Face)
}

type Anim struct {
	Num   int
	Faces []int
}

func AddAnim(msg messages.MessageAnim) {
	anim := &Anim{
		Num:   int(msg.AnimID),
		Faces: make([]int, len(msg.Faces)),
	}
	for i, face := range msg.Faces {
		anim.Faces[i] = int(face)
	}
	anims[int(msg.AnimID)] = anim
}

func GetAnim(num int) *Anim {
	return anims[num]
}
