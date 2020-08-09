package test

import (
	"testing"

	"dwc.com/ancho/deveui"
)

func TestGenerateUniqueHexCodes(t *testing.T) {

	generator := deveui.NewDevEUIGenerator()
	iterations := 1000
	expectedDevEUILength := 16
	totalDevEUIs := 100

	for iter := 0; iter < iterations; iter++ {
		known := make(map[string]bool)
		devEUIs := make([]string, totalDevEUIs)

		for index := 0; index < totalDevEUIs; index++ {

			eui, err := generator.GeneratDevEUI(expectedDevEUILength)
			if err != nil {
				t.Errorf("cannot generate devEUI: %v", err)
			}

			shortCode := eui[len(eui)-5:]
			if known[shortCode] {
				t.Fatalf("unique short code detected: iteration: %d, eui: %d", iter, index)
			}
			known[shortCode] = true

			if len(eui) != expectedDevEUILength {
				t.Error("devEUI has incorrect length")
			}

			devEUIs[index] = eui
		}

		if len(devEUIs) != totalDevEUIs {
			t.Error("devEUIs slice has incorrect length")
		}
	}
}
