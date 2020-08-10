package test

import (
	"log"
	"testing"

	"dwc.com/ancho/deveui"
)

func Test_GenerateUniqueDevEUIs(t *testing.T) {

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

func Test_CanGenerateXDevEUIs(t *testing.T) {
	generator := deveui.NewDevEUIGenerator()
	expectedDevEUILength := 16
	totalDevEUIs := 100

	devEUIs, err := generator.GeneratDevEUIs(totalDevEUIs, expectedDevEUILength)
	if err != nil {
		t.Errorf("cannot generate devEUI: %v", err)
	}

	if len(devEUIs) != totalDevEUIs {
		t.Error("Incorrect number of batches created")
	}
}

func Test_CanSplitIntoBatches(t *testing.T) {
	generator := deveui.NewDevEUIGenerator()
	expectedDevEUILength := 16
	totalDevEUIs := 100
	totalBatches := 10

	devEUIs, err := generator.GeneratDevEUIs(totalDevEUIs, expectedDevEUILength)
	if err != nil {
		log.Fatalf("cannot generate devEUI: %v", err)
	}

	batches := deveui.ChunkDevEUIs(devEUIs, totalBatches)

	if len(batches) != totalBatches {
		t.Error("Incorrect number of batches created")
	}

	for _, batch := range batches {
		if len(batch) != totalDevEUIs/totalBatches {
			t.Error("Incorrect number of items in batch")
		}
	}
}
