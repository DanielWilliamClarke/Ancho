package main

import (
	"fmt"
	"log"
	"time"

	"dwc.com/ancho/deveui"
)

func main() {
	// Setup
	maxInFlightRequests := 10
	totalDevEUIs := 100

	generator := deveui.NewDevEUIGenerator(time.Now().UnixNano())
	registrationAPI := deveui.NewRegistrationClientAPI()
	routines := deveui.NewRegistrationRoutines(registrationAPI, generator)
	batcher := deveui.NewRegistrationBatcher(routines, maxInFlightRequests)

	// Register
	registeredDevEUIs, errors := batcher.RegisterInParallel(totalDevEUIs)

	// Output
	fmt.Printf("%d DevEUIs failed ---------------- \n", len(errors))
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
