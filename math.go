package toolkit

import (
	"math/rand"
	"sync"
	"time"
)

type randomizer struct {
	sync.Mutex
	r *rand.Rand
}

var r *randomizer

func (r *randomizer) Intn(limit int) int {
	defer r.Unlock()
	r.Lock()
	return r.r.Intn(limit)
}

func initRandomSource() {
	if r == nil {
		src := rand.NewSource(time.Now().UnixNano())
		r = new(randomizer)
		r.r = rand.New(src)
	}
}

func RandInt(limit int) int {
	initRandomSource()
	return r.Intn(limit)
}
