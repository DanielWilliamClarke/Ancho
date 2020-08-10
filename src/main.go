package main

import (
	"fmt"
	"log"

	"dwc.com/ancho/deveui"
)

func main() {
	requestsPerBatch := 10
	totalDevEUIs := 100

	generator := deveui.NewDevEUIGenerator()
	registrationAPI := deveui.NewRegistrationClientAPI()
	batcher := deveui.NewRegistrationBatcher(registrationAPI, requestsPerBatch)

	devEUIs, err := generator.GeneratDevEUIs(totalDevEUIs, 16)
	if err != nil {
		log.Fatalf("cannot generate devEUI: %v", err)
	}

	registeredDevEUIs, errors := batcher.RegisterInParallel(devEUIs)

	fmt.Printf("%d DevEUIs registered successfully ---------------- \n", len(errors))
	for _, err := range errors {
		log.Println(err)
	}
	fmt.Println("---------------------------------------------------")

	fmt.Printf("%d DevEUIs registered successfully ---------------- \n", len(registeredDevEUIs))
	for _, devEUI := range registeredDevEUIs {
		fmt.Println(devEUI)
	}
	fmt.Println("---------------------------------------------------")
}
