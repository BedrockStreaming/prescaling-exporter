package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BedrockStreaming/prescaling-exporter/pkg/services"
)

// MockHPAService is a mock of the HPA service for testing
type MockHPAService struct {
	hpaMinimumStatus      *services.HPAMinimumStatus
	platformScalingStatus *services.PlatformScalingStatus
	err                   error
}

func (m *MockHPAService) CheckHPAMinimums() (*services.HPAMinimumStatus, error) {
	return m.hpaMinimumStatus, m.err
}

func (m *MockHPAService) GetPlatformScalingStatus() (*services.PlatformScalingStatus, error) {
	return m.platformScalingStatus, m.err
}

func TestNewHPAHandlers(t *testing.T) {
	mockService := &MockHPAService{}
	handlers := NewHPAHandlers(mockService)

	if handlers == nil {
		t.Error("NewHPAHandlers should return a non-nil instance")
	}
}

func TestHPAHandlers_CheckHPAMinimums(t *testing.T) {
	// Test with all HPAs meeting their minimums
	t.Run("all_hpas_meet_minimum", func(t *testing.T) {
		// Create a mock of the HPA service that returns a status where all HPAs meet their minimums
		mockService := &MockHPAService{
			hpaMinimumStatus: &services.HPAMinimumStatus{
				AllHPAsMeetMinimum: true,
				HPAInfos:           []services.HPAInfo{},
			},
			err: nil,
		}

		// Create the handler with the mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Create a fake HTTP request
		req, err := http.NewRequest("GET", "/api/v1/hpas", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Call the handler
		handler.CheckHPAMinimums(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned incorrect status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body
		var response services.HPAMinimumStatus
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Unable to deserialize response: %v", err)
		}

		if !response.AllHPAsMeetMinimum {
			t.Errorf("Incorrect AllHPAsMeetMinimum in response: got %v, want %v", response.AllHPAsMeetMinimum, true)
		}
	})

	// Test with some HPAs not meeting their minimums
	t.Run("some_hpas_below_minimum", func(t *testing.T) {
		// Create a mock of the HPA service that returns a status where some HPAs don't meet their minimums
		mockService := &MockHPAService{
			hpaMinimumStatus: &services.HPAMinimumStatus{
				AllHPAsMeetMinimum: false,
				HPAInfos: []services.HPAInfo{
					{
						Namespace:             "default",
						Name:                  "test-hpa",
						TargetKind:            "Deployment",
						TargetName:            "test-deployment",
						HpaMinReplicas:        2,
						MinReplicasAnnotation: 3,
						MeetsMinimum:          false,
					},
				},
			},
			err: nil,
		}

		// Create the handler with the mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Create a fake HTTP request
		req, err := http.NewRequest("GET", "/api/v1/hpas", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Call the handler
		handler.CheckHPAMinimums(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusServiceUnavailable {
			t.Errorf("handler returned incorrect status code: got %v want %v",
				status, http.StatusServiceUnavailable)
		}

		// Check the response body
		var response services.HPAMinimumStatus
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Unable to deserialize response: %v", err)
		}

		if response.AllHPAsMeetMinimum {
			t.Errorf("Incorrect AllHPAsMeetMinimum in response: got %v, want %v", response.AllHPAsMeetMinimum, false)
		}

		if len(response.HPAInfos) != 1 {
			t.Errorf("HPAInfos should contain 1 element, but contains %d", len(response.HPAInfos))
		}
	})

	// Test with error
	t.Run("error", func(t *testing.T) {
		// Create a mock of the HPA service that returns an error
		mockService := &MockHPAService{
			hpaMinimumStatus: nil,
			err:              errors.New("test error"),
		}

		// Create the handler with the mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Create a fake HTTP request
		req, err := http.NewRequest("GET", "/api/v1/hpas", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Call the handler
		handler.CheckHPAMinimums(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned incorrect status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func TestHPAHandlers_GetPlatformScalingStatus(t *testing.T) {
	// Test with the platform correctly scaled
	t.Run("platform_scaled", func(t *testing.T) {
		// Create a mock of the HPA service that returns a status where the platform is correctly scaled
		mockService := &MockHPAService{
			platformScalingStatus: &services.PlatformScalingStatus{
				IsPlatformScaled: true,
			},
			err: nil,
		}

		// Create the handler with the mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Create a fake HTTP request
		req, err := http.NewRequest("GET", "/api/v1/hpas/check", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Call the handler
		handler.GetPlatformScalingStatus(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned incorrect status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body
		var response services.PlatformScalingStatus
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Unable to deserialize response: %v", err)
		}

		if !response.IsPlatformScaled {
			t.Errorf("Incorrect IsPlatformScaled in response: got %v, want %v", response.IsPlatformScaled, true)
		}
	})

	// Test with the platform not correctly scaled
	t.Run("platform_not_scaled", func(t *testing.T) {
		// Create a mock of the HPA service that returns a status where the platform is not correctly scaled
		mockService := &MockHPAService{
			platformScalingStatus: &services.PlatformScalingStatus{
				IsPlatformScaled: false,
			},
			err: nil,
		}

		// Create the handler with the mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Create a fake HTTP request
		req, err := http.NewRequest("GET", "/api/v1/hpas/check", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Call the handler
		handler.GetPlatformScalingStatus(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusServiceUnavailable {
			t.Errorf("handler returned incorrect status code: got %v want %v",
				status, http.StatusServiceUnavailable)
		}

		// Check the response body
		var response services.PlatformScalingStatus
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Unable to deserialize response: %v", err)
		}

		if response.IsPlatformScaled {
			t.Errorf("Incorrect IsPlatformScaled in response: got %v, want %v", response.IsPlatformScaled, false)
		}
	})

	// Test with error
	t.Run("error", func(t *testing.T) {
		// Create a mock of the HPA service that returns an error
		mockService := &MockHPAService{
			platformScalingStatus: nil,
			err:                   errors.New("test error"),
		}

		// Create the handler with the mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Create a fake HTTP request
		req, err := http.NewRequest("GET", "/api/v1/hpas/check", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response
		rr := httptest.NewRecorder()

		// Call the handler
		handler.GetPlatformScalingStatus(rr, req)

		// Check the status code
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned incorrect status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}
