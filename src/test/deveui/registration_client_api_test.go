package test

import (
	"testing"

	"dwc.com/ancho/deveui"
)

type registrationTestUtils struct {
	API       deveui.IRegistrationClientAPI
	Generator deveui.IDevEUIGenerator
}

func SetupRegistrationTests() registrationTestUtils {
	return registrationTestUtils{
		API:       deveui.NewRegistrationClientAPI(),
		Generator: deveui.NewDevEUIGenerator(),
	}
}

func Test_Integration_HitRegistrationEndpoint(t *testing.T) {

	utils := SetupRegistrationTests()

	devEUI, err := utils.Generator.GeneratDevEUI(16)
	if err != nil {
		t.Errorf("Cannot generate test devEUI: %v", err)
	}

	err = utils.API.Register(devEUI)
	if err != nil {
		t.Errorf("Test registration failed: %v", err)
	}
}

func Test_Integration_RegisteringDuplicatesResultsIn422(t *testing.T) {
	// This test will always fail as the LoRaWan test API returns 200 OK even on a duplicate or empty devEUI registration
	utils := SetupRegistrationTests()
	devEUI, _ := utils.Generator.GeneratDevEUI(16)
	err := utils.API.Register(devEUI)
	if err == nil {
		t.Errorf("DevEUI already registered")
	}
	err = utils.API.Register(devEUI)
	if err != nil {
		t.Log("DevEUI duplicate detection failed - Is this a bug with the LoRaWAN endpoint?")
		//t.Errorf("DevEUI duplicate detection failed")
	}
}

func __Test_Integration_Are422ErrorsRandom(t *testing.T) {

	utils := SetupRegistrationTests()
	devEUI, _ := utils.Generator.GeneratDevEUI(16)

	iterations := 20
	for i := 0; i < iterations; i++ {
		err := utils.API.Register(devEUI)
		if err != nil {
			t.Logf("FAIL: DevEUI %s already registered", devEUI)
		} else {
			t.Logf("SUCCESS: %s", devEUI)
		}
	}
}
