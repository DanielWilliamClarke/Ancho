package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber"
	"github.com/gofiber/logger"

	"dwc.com/ancho/deveui"
)

func main() {

	responseCache, err := deveui.NewResponseCache()
	if err != nil {
		log.Fatalf("Cannot create cache: %v", err)
	}

	maxInFlightRequests := 10
	totalDevEUIs := 100
	registrationAPI := deveui.NewRegistrationClientAPI()

	app := fiber.New()
	api := app.Group("/v1/api", logger.New())

	api.Put("/register", func(c *fiber.Ctx) {

		// Check Cache
		key := c.Get("Idempotency-key")
		if cache, err := responseCache.Load(key); err == nil {
			err = c.Status(200).JSON(cache.Payload)
			if err != nil {
				log.Printf("Could not send payload data: %v", err)
			}
			return
		}
		// Store empty payload to ensure any requests that occur before this one completes,
		// does not try to trigger more registrations with the same key
		responseCache.Store(key, deveui.DevEUIPayload{
			DevEUIs: []string{},
		})

		// Setup
		generator := deveui.NewDevEUIGenerator(time.Now().UnixNano())
		routines := deveui.NewRegistrationRoutines(registrationAPI, generator)
		batcher := deveui.NewRegistrationBatcher(routines, maxInFlightRequests)

		// Register
		registeredDevEUIs, errors := batcher.RegisterInParallel(totalDevEUIs)

		// Log Errors
		fmt.Printf("%d DevEUIs failed ---------------- \n", len(errors))
		for _, err := range errors {
			log.Println(err)
		}
		fmt.Println("---------------------------------------------------")

		// Package payload
		payload := deveui.DevEUIPayload{
			DevEUIs: registeredDevEUIs,
		}

		// Store payload in cache
		err = responseCache.Store(key, payload)
		if err != nil {
			log.Printf("Could not store idempotent payload data: %v", err)
		}

		// Render output to client
		err = c.Status(200).JSON(payload)
		if err != nil {
			log.Printf("Could not send payload data: %v", err)
		}
	})

	port := 3000
	err = app.Listen(port)
	if err != nil {
		log.Printf("Could not start api server on port: %d -> %v", port, err)
	}
}
