package canvas

// Rotation represents a screen rotation, using macOS "terminology"
type Rotation int

const (
	// Rotation0 respresents the standard rotation
	Rotation0 Rotation = iota * 90
	// Rotation90 respresents a 90° rotation (typically portrait)
	Rotation90
	// Rotation180 respresents a 180° rotation (typically landscape flipped)
	Rotation180
	// Rotation270 respresents a 270° rotation (typically portrait flipped)
	Rotation270
)

// Canvas represents current settings of a canvas
// a canvas can be displayed on multiple clients (with identical settings, basically only for testing purposes)
type Canvas struct {
	Key       string   `json:"key"`
	Rotation  Rotation `json:"rotation"`
	Artbundle string   `json:"artbundle"` // the path to the currently displayed Artbundle, from the artwork folder
}
