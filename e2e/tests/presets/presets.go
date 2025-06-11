package presets

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/welovemedia/ffmate/e2e/tests"
	"github.com/welovemedia/ffmate/internal/dto"
)

var preset dto.Preset
var presets []dto.Preset

func RunTests() error {
	webhookServer, webhookChan := tests.SetupWebhookServer(2)

	tests.RegisterWebhook(webhookServer, dto.PRESET_CREATED)
	tests.RegisterWebhook(webhookServer, dto.PRESET_DELETED)

	if err := testPresetCreation(); err != nil {
		return err
	}
	tests.WaitForWebhook(webhookChan, dto.PRESET_CREATED)

	if err := testPresetList(); err != nil {
		return err
	}

	if err := testPresetDeletion(); err != nil {
		return err
	}
	tests.WaitForWebhook(webhookChan, dto.PRESET_DELETED)

	return nil
}

func testPresetCreation() error {
	resp, err := resty.New().R().
		SetBody(&dto.NewPreset{
			Name:        "Test Preset",
			Description: "Test Preset Description",
			Command:     "-i {input} -c copy {output}",
		}).
		SetResult(&preset).
		Post("http://localhost:3000/api/v1/presets")

	if err != nil || resp.StatusCode() >= 400 {
		return fmt.Errorf("failed to create preset: %v", resp)
	} else {
		fmt.Printf("Created '%s' (uuid: %s)\n", preset.Name, preset.Uuid)
	}

	return nil
}

func testPresetList() error {
	resp, err := resty.New().R().
		SetResult(&presets).
		Get("http://localhost:3000/api/v1/presets")

	if err != nil || resp.StatusCode() >= 400 {
		return fmt.Errorf("failed to list presets: %v", resp)
	} else {
		fmt.Printf("Listed '%d' presets\n", len(presets))
	}

	return nil
}

func testPresetDeletion() error {
	resp, err := resty.New().R().
		Delete("http://localhost:3000/api/v1/presets/" + preset.Uuid)

	if err != nil || resp.StatusCode() >= 400 {
		return fmt.Errorf("failed to delete preset: %v", resp)
	} else {
		fmt.Printf("Deleted '%s' (uuid: %s)\n", preset.Name, preset.Uuid)
	}

	return nil
}
