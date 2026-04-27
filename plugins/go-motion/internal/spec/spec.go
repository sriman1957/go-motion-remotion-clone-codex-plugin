package spec

type Composition struct {
	Title          string  `json:"title"`
	FPS            int     `json:"fps"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	DurationFrames int     `json:"durationFrames"`
	Scenes         []Scene `json:"scenes"`
}

type Scene struct {
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	Headline       string `json:"headline"`
	Body           string `json:"body"`
	DurationFrames int    `json:"durationFrames"`
}
