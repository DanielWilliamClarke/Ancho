package deveui

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type IRegistrationRoutines interface {
	RunBatch(config ParallelRegistrationConfig)
	Observe(registeredChannel chan string, requiredDevEUIS int, shutdownChannel chan struct{})
	GracefulShutdown(sigChannel chan os.Signal, shutdownChannel chan struct{})
	CleanUp(sigChannel chan os.Signal)
}

func NewRegistrationRoutines(api IRegistrationClientAPI, generator IDevEUIGenerator) *RegistrationRoutines {
	return &RegistrationRoutines{
		api:       api,
		generator: generator,
	}
}

type RegistrationRoutines struct {
	api       IRegistrationClientAPI
	generator IDevEUIGenerator
}

type ParallelRegistrationConfig struct {
	RegisteredChannel chan string
	RequiredDevEUIS   int
	ShutdownChannel   chan struct{}
	ErrorChannel      chan error
	WaitGroup         *sync.WaitGroup
	GoroutineIndex    int
}

// RunBatch goroutine that runs until its expected number of DevEUis have been registered, can be shutdown using SIGINT
func (r RegistrationRoutines) RunBatch(config ParallelRegistrationConfig) {
	fmt.Printf("Starting registration batch: %d \n", config.GoroutineIndex)
	defer config.WaitGroup.Done()

	registeredDevEUIs := 0

	for registeredDevEUIs < config.RequiredDevEUIS {
		select {
		case <-config.ShutdownChannel:
			log.Printf("Shutdown registration batch: %d \n", config.GoroutineIndex)
			return
		default:
			devEUI, err := r.generator.GenerateDevEUI(16)
			if err != nil {
				config.ErrorChannel <- err
				return
			}
			err = r.api.Register(devEUI)
			if err != nil {
				config.ErrorChannel <- err
			} else {
				config.RegisteredChannel <- devEUI
				registeredDevEUIs++
			}
		}
	}

	fmt.Printf("Registration batch %d complete! \n", config.GoroutineIndex)
}

// Observe reports total number of DevEUIs registered by all goroutines, can be shutdown using SIGINT
func (r RegistrationRoutines) Observe(registeredChannel chan string, requiredDevEUIS int, shutdownChannel chan struct{}) {
	for len(registeredChannel) < requiredDevEUIS {
		select {
		case <-shutdownChannel:
			fmt.Println("Shutdown obervation \n")
			return
		default:
			time.Sleep(1 * time.Second)
			fmt.Printf("Total DevEUIs registered: %d \n", len(registeredChannel))
		}
	}
	fmt.Println("Registration complete")
}

// GracefulShutdown blocks in its goroutine until SIGINT has been triggered, closing the shutdown channel to gracefully end registration if required
func (r RegistrationRoutines) GracefulShutdown(sigChannel chan os.Signal, shutdownChannel chan struct{}) {
	<-sigChannel           // received SIGINT or SIGTERM
	close(shutdownChannel) // Signal all goroutines to shutdown
	fmt.Println("Quit signal received, gracefully shutdown registration...")
}

// CleanUp closes the SIGINT signal channel on defer
func (r RegistrationRoutines) CleanUp(sigChannel chan os.Signal) {
	fmt.Println("Cleaning Up")
	close(sigChannel)
}
