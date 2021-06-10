package toolkit

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type randomizer struct {
	sync.Mutex
	r *rand.Rand
}

func (r *randomizer) Intn(limit int) int {
	defer r.Unlock()
	r.Lock()
	return r.r.Intn(limit)
}

var (
	once sync.Once
	r    *randomizer
)

func initRandomSource() {
	once.Do(func() {
		src := rand.NewSource(time.Now().UnixNano())
		r = new(randomizer)
		r.r = rand.New(src)
	})
}

func RandInt(limit int) int {
	initRandomSource()
	return r.Intn(limit)
}

func RandFloat(limit int, decimal int) float64 {
	flim := float64(limit)
	fdec := float64(decimal)
	initRandomSource()
	powerLimit := int(flim * math.Pow(10, fdec))
	randPower := r.Intn(powerLimit)
	return float64(randPower) / math.Pow(10, fdec)
}

func Div(f1, f2 float64) float64 {
	if f2 == 0 {
		return 0
	}

	return f1 / f2
}
