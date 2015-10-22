package toolkit

import (
	"math/rand"
	"time"
)

var r *rand.Rand

func initRandomSource() {
	if r == nil {
		src := rand.NewSource(time.Now().UnixNano())
		r = rand.New(src)
	}
}

func RandInt(limit int) int {
	initRandomSource()
	return r.Intn(limit)
}
