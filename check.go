package main

import (
	"time"
)

type check struct {
	statusCode    string
	nameLookup    time.Duration
	connect       time.Duration
	contentFetch  time.Duration
	contentSizeKB uint16
}
