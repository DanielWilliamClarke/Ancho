package deveui

import (
	"encoding/hex"
	"math/rand"
	"sync"
)

type IDevEUIGenerator interface {
	GenerateDevEUI(length int) (string, error)
}

func NewDevEUIGenerator(seed int64) *DevEUIGenerator {
	rand.Seed(seed)
	return &DevEUIGenerator{&sync.Map{}}
}

type DevEUIGenerator struct {
	known *sync.Map
}

// GenerateDevEUI generates a unique random hex value of length.
func (d DevEUIGenerator) GenerateDevEUI(length int) (devEUI string, err error) {
	// to reduce the amount of memory allocated divide the length by 2
	// as encoded hex string are of length * 2
	bytes := make([]byte, (length+1)/2)

	_, err = rand.Read(bytes)
	if err != nil {
		return "", err
	}

	devEUI = hex.EncodeToString(bytes)[:length]

	// If short code is known, recursively generate again until short code is not known
	shortCode := devEUI[len(devEUI)-5:]
	if _, ok := d.known.Load(shortCode); ok {
		devEUI, err = d.GenerateDevEUI(length)
	}
	d.known.Store(shortCode, true)

	return devEUI, err
}
