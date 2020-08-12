package deveui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

const (
	apiHost string = "https://europe-west1-machinemax-dev-d524.cloudfunctions.net"
)

type IRegistrationClientAPI interface {
	Register(devEUI string) error
}

func NewRegistrationClientAPI() *LoRaWANClientAPI {
	client := &http.Client{}
	return &LoRaWANClientAPI{client, apiHost}
}

type LoRaWANClientAPI struct {
	Client *http.Client
	Url    string
}

type devEUIBody struct {
	DevEUI string `json:"deveui"`
}

// Register registers a given devEUI with the test LoRaWAN api
func (l LoRaWANClientAPI) Register(devEUI string) error {

	requestBody, err := json.Marshal(devEUIBody{DevEUI: devEUI})
	if err != nil {
		log.Printf("Could not marshal deveui into body: %v", err)
		return err
	}

	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/sensor-onboarding-sample", l.Url),
		bytes.NewBuffer(requestBody))

	if err != nil {
		log.Printf("Could not create API request: %v", err)
		return err
	}

	response, err := l.Client.Do(request)
	if err != nil {
		log.Printf("LoRaWAN registration request failed: %v", err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(
			fmt.Sprintf("DevEUI %s already Registered: %d", devEUI, response.StatusCode))
	}

	return nil
}
