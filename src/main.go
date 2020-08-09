package main

import (
	"fmt"
	"log"

	"dwc.com/ancho/deveui"
)

func main() {
	generator := deveui.NewDevEUIGenerator()

	totalDevEUIs := 100
	devEUIs := make([]string, totalDevEUIs)
	for i := 0; i < totalDevEUIs; i++ {
		eui, err := generator.GenerateHexCode(16)
		if err != nil {
			log.Fatalf("cannot generate devEUI [ %d: %v ]", i, err)
		}
		devEUIs[i] = eui
		fmt.Println(eui)
	}
}
