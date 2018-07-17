package main

type Check struct {
	Err         error   `json:"-"`
	StatusCode  int     `json:"statusCode"`
	Bytes       int     `json:"bytes"`
	TimeElapsed float64 `json:"timeElapsed"`
}
