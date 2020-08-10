package deveui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	apiHost string = "https://europe-west1-machinemax-dev-d524.cloudfunctions.net"
)

type IRegistrationClientAPI interface {
	Register(devEUI string) (bool, error)
}

func NewRegistrationAPI() *LoRaWANClientAPI {
	client := &http.Client{}
	return &LoRaWANClientAPI{client}
}

type LoRaWANClientAPI struct {
	client *http.Client
}

type devEUIBody struct {
	devEUI string `json:"deveui"`
}

// Register registers a given devEUI with the test LoRaWAN api
func (l LoRaWANClientAPI) Register(devEUI string) (bool, error) {

	requestBody, err := json.Marshal(devEUIBody{devEUI: devEUI})
	if err != nil {
		log.Printf("Could not marshal deveui into body: %v", err)
		return false, err
	}

	request, err := http.NewRequest("POST",
		fmt.Sprintf("%s/sensor-onboarding-sample", apiHost),
		bytes.NewBuffer(requestBody))

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Cache-Control", "private, no-store, max-age=0")

	if err != nil {
		log.Printf("Could not create API request: %v", err)
		return false, err
	}

	response, err := l.client.Do(request)
	if err != nil {
		log.Printf("LoRaWAN registration request failed: %v", err)
		return false, err
	}

	success := response.StatusCode == http.StatusOK

	return success, nil
}
