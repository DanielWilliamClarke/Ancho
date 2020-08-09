package deveui

import (
	"encoding/hex"
	"math/rand"
	"time"

	mt19937 "github.com/seehuhn/mt19937"
)

type IDevEUIGenerator interface {
	GeneratDevEUI(length int) (string, error)
}

func NewDevEUIGenerator() *DevEUIGenerator {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	return &DevEUIGenerator{rng, make(map[string]bool)}
}

type DevEUIGenerator struct {
	rng   *rand.Rand
	known map[string]bool
}

// GeneratDevEUI generates a unique random hex value of length.
func (d DevEUIGenerator) GeneratDevEUI(length int) (devEUI string, err error) {
	// to reduce the amount of memory allocated divide the length by 2
	// as encoded hex string are of length * 2
	bytes := make([]byte, (length+1)/2)

	_, err = d.rng.Read(bytes)
	if err != nil {
		return "", err
	}

	devEUI = hex.EncodeToString(bytes)[:length]

	// If short code is know, recursively generate again until short code is not known
	shortCode := devEUI[len(devEUI)-5:]
	if d.known[shortCode] {
		devEUI, err = d.GeneratDevEUI(length)
	}
	d.known[shortCode] = true

	return devEUI, err
}
