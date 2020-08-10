package test

import (
	"testing"

	"dwc.com/ancho/deveui"
)

func Test_GenerateUniqueDevEUIs(t *testing.T) {

	// The main point of this test is to show that generator.GeneratDevEUI will only create unique devEUIs

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
				t.Fatalf("non unique short code detected: iteration: %d, eui: %s", iter, shortCode)
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
