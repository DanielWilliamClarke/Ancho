package test

import (
	"testing"

	"dwc.com/ancho/deveui"
)

func TestGenerateUniqueHexCodes(t *testing.T) {

	expectedDevEUILength := 16
	totalDevEUIs := 100

	known := make(map[string]bool)
	devEUIs := make([]string, totalDevEUIs)
	for i := 0; i < totalDevEUIs; i++ {

		eui, err := deveui.GenerateHexCode(expectedDevEUILength)
		if err != nil {
			t.Error("cannot generate devEUI")
		}

		shortCode := eui[len(eui)-5:]
		if known[shortCode] {
			t.Error("unique short code detected")
		}
		known[shortCode] = true

		if len(eui) != expectedDevEUILength {
			t.Error("devEUI has incorrect length")
		}

		devEUIs[i] = eui
	}

	if len(devEUIs) != totalDevEUIs {
		t.Error("devEUIs slice has incorrect length")
	}
}
