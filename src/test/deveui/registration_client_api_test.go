package test

import (
	"testing"

	"dwc.com/ancho/deveui"
)

type registrationTestUtils struct {
	api       deveui.IRegistrationClientAPI
	generator deveui.IDevEUIGenerator
}

func SetupRegistrationTests() registrationTestUtils {
	return registrationTestUtils{
		api:       deveui.NewRegistrationClientAPI(),
		generator: deveui.NewDevEUIGenerator(),
	}
}

func Test_Integration_HitRegistrationEndpoint(t *testing.T) {

	// the aim of this test is to ensure that the client api can hit the endpoint

	utils := SetupRegistrationTests()

	devEUI, err := utils.generator.GeneratDevEUI(16)
	if err != nil {
		t.Errorf("Cannot generate test devEUI: %v", err)
	}

	err = utils.api.Register(devEUI)
	if err != nil {
		t.Errorf("Test registration failed: %v", err)
	}
}

func Test_Integration_Are422ErrorsRandom(t *testing.T) {

	// the aim of this test was to confirm my suspicions that the API was no reliable

	utils := SetupRegistrationTests()
	devEUI, _ := utils.generator.GeneratDevEUI(16)

	iterations := 20
	for i := 0; i < iterations; i++ {
		err := utils.api.Register(devEUI)
		if err != nil {
			t.Logf("FAIL: DevEUI %s already registered", devEUI)
		} else {
			t.Logf("SUCCESS: %s", devEUI)
		}
	}
}
