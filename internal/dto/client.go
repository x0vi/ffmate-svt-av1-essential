package dto

type Client struct {
	Version string `json:"version"`
	FFmpeg  string `json:"ffmpeg"`
	Os      string `json:"os"`
	Arch    string `json:"arch"`
}
