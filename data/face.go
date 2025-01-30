package data

import (
	"github.com/kettek/termfire/messages"
)

// FaceImage is a merger of face and image, cuz why not.
type FaceImage struct {
	Num      uint16
	Set      uint8
	Width    int
	Height   int
	Data     []byte
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

var faces = make(map[int]FaceImage)

// GetFace returns a face from the face map.
func GetFace(num int) (FaceImage, bool) {
	face, ok := faces[num]
	return face, ok
}

// AddFace adds a pending face to the face map. Returns if it exists.
func AddFace(face messages.MessageFace2) bool {
	_, ok := faces[int(face.Num)]
	if ok {
		return true
	}
	faces[int(face.Num)] = FaceImage{
		Num:      uint16(face.Num),
		Set:      uint8(face.SetNum),
		name:     face.Name,
		Checksum: face.Checksum,
		pending:  true,
	}
	return false
}

// AddFaceImage adds an image to the face map.
func AddFaceImage(image messages.MessageImage2) {
	face, ok := faces[int(image.Face)]
	if !ok {
		faces[int(image.Face)] = FaceImage{
			Num:     uint16(image.Face),
			Set:     image.Set,
			Width:   image.Width,
			Height:  image.Height,
			Data:    image.Data,
			pending: false,
		}
		return
	}
	faces[int(image.Face)] = FaceImage{
		Num:      uint16(image.Face),
		Set:      image.Set,
		Width:    image.Width,
		Height:   image.Height,
		Data:     image.Data,
		name:     face.name,
		Checksum: face.Checksum,
		pending:  false,
	}
}
