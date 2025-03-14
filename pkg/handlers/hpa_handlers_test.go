package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BedrockStreaming/prescaling-exporter/pkg/services"
)

// MockHPAService est un mock du service HPA pour les tests
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
		t.Error("NewHPAHandlers devrait retourner une instance non nil")
	}
}

func TestHPAHandlers_CheckHPAMinimums(t *testing.T) {
	// Test avec tous les HPA respectant leurs minimums
	t.Run("all_hpas_meet_minimum", func(t *testing.T) {
		// Créer un mock du service HPA qui retourne un statut où tous les HPA respectent leurs minimums
		mockService := &MockHPAService{
			hpaMinimumStatus: &services.HPAMinimumStatus{
				AllHPAsMeetMinimum: true,
				HPAInfos:           []services.HPAInfo{},
			},
			err: nil,
		}

		// Créer le handler avec le mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Créer une requête HTTP factice
		req, err := http.NewRequest("GET", "/api/v1/hpa-check", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Créer un ResponseRecorder pour enregistrer la réponse
		rr := httptest.NewRecorder()

		// Appeler le handler
		handler.CheckHPAMinimums(rr, req)

		// Vérifier le code de statut
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler a retourné un code de statut incorrect: got %v want %v",
				status, http.StatusOK)
		}

		// Vérifier le corps de la réponse
		var response services.HPAMinimumStatus
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Impossible de désérialiser la réponse: %v", err)
		}

		if !response.AllHPAsMeetMinimum {
			t.Errorf("AllHPAsMeetMinimum incorrect dans la réponse: got %v, want %v", response.AllHPAsMeetMinimum, true)
		}
	})

	// Test avec certains HPA ne respectant pas leurs minimums
	t.Run("some_hpas_below_minimum", func(t *testing.T) {
		// Créer un mock du service HPA qui retourne un statut où certains HPA ne respectent pas leurs minimums
		mockService := &MockHPAService{
			hpaMinimumStatus: &services.HPAMinimumStatus{
				AllHPAsMeetMinimum: false,
				HPAInfos: []services.HPAInfo{
					{
						Namespace:       "default",
						Name:            "test-hpa",
						TargetKind:      "Deployment",
						TargetName:      "test-deployment",
						HpaMinReplicas:  2,
						MinReplicasAnnotation: 3,
						MeetsMinimum:    false,
					},
				},
			},
			err: nil,
		}

		// Créer le handler avec le mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Créer une requête HTTP factice
		req, err := http.NewRequest("GET", "/api/v1/hpa-check", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Créer un ResponseRecorder pour enregistrer la réponse
		rr := httptest.NewRecorder()

		// Appeler le handler
		handler.CheckHPAMinimums(rr, req)

		// Vérifier le code de statut
		if status := rr.Code; status != http.StatusServiceUnavailable {
			t.Errorf("handler a retourné un code de statut incorrect: got %v want %v",
				status, http.StatusServiceUnavailable)
		}

		// Vérifier le corps de la réponse
		var response services.HPAMinimumStatus
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Impossible de désérialiser la réponse: %v", err)
		}

		if response.AllHPAsMeetMinimum {
			t.Errorf("AllHPAsMeetMinimum incorrect dans la réponse: got %v, want %v", response.AllHPAsMeetMinimum, false)
		}

		if len(response.HPAInfos) != 1 {
			t.Errorf("HPAInfos devrait contenir 1 élément, mais en contient %d", len(response.HPAInfos))
		}
	})

	// Test avec erreur
	t.Run("error", func(t *testing.T) {
		// Créer un mock du service HPA qui retourne une erreur
		mockService := &MockHPAService{
			hpaMinimumStatus: nil,
			err:              errors.New("erreur de test"),
		}

		// Créer le handler avec le mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Créer une requête HTTP factice
		req, err := http.NewRequest("GET", "/api/v1/hpa-check", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Créer un ResponseRecorder pour enregistrer la réponse
		rr := httptest.NewRecorder()

		// Appeler le handler
		handler.CheckHPAMinimums(rr, req)

		// Vérifier le code de statut
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler a retourné un code de statut incorrect: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func TestHPAHandlers_GetPlatformScalingStatus(t *testing.T) {
	// Test avec la plateforme correctement scalée
	t.Run("platform_scaled", func(t *testing.T) {
		// Créer un mock du service HPA qui retourne un statut où la plateforme est correctement scalée
		mockService := &MockHPAService{
			platformScalingStatus: &services.PlatformScalingStatus{
				IsPlatformScaled: true,
			},
			err: nil,
		}

		// Créer le handler avec le mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Créer une requête HTTP factice
		req, err := http.NewRequest("GET", "/api/v1/platform-scaling", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Créer un ResponseRecorder pour enregistrer la réponse
		rr := httptest.NewRecorder()

		// Appeler le handler
		handler.GetPlatformScalingStatus(rr, req)

		// Vérifier le code de statut
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler a retourné un code de statut incorrect: got %v want %v",
				status, http.StatusOK)
		}

		// Vérifier le corps de la réponse
		var response services.PlatformScalingStatus
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Impossible de désérialiser la réponse: %v", err)
		}

		if !response.IsPlatformScaled {
			t.Errorf("IsPlatformScaled incorrect dans la réponse: got %v, want %v", response.IsPlatformScaled, true)
		}
	})

	// Test avec la plateforme non correctement scalée
	t.Run("platform_not_scaled", func(t *testing.T) {
		// Créer un mock du service HPA qui retourne un statut où la plateforme n'est pas correctement scalée
		mockService := &MockHPAService{
			platformScalingStatus: &services.PlatformScalingStatus{
				IsPlatformScaled: false,
			},
			err: nil,
		}

		// Créer le handler avec le mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Créer une requête HTTP factice
		req, err := http.NewRequest("GET", "/api/v1/platform-scaling", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Créer un ResponseRecorder pour enregistrer la réponse
		rr := httptest.NewRecorder()

		// Appeler le handler
		handler.GetPlatformScalingStatus(rr, req)

		// Vérifier le code de statut
		if status := rr.Code; status != http.StatusServiceUnavailable {
			t.Errorf("handler a retourné un code de statut incorrect: got %v want %v",
				status, http.StatusServiceUnavailable)
		}

		// Vérifier le corps de la réponse
		var response services.PlatformScalingStatus
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Impossible de désérialiser la réponse: %v", err)
		}

		if response.IsPlatformScaled {
			t.Errorf("IsPlatformScaled incorrect dans la réponse: got %v, want %v", response.IsPlatformScaled, false)
		}
	})

	// Test avec erreur
	t.Run("error", func(t *testing.T) {
		// Créer un mock du service HPA qui retourne une erreur
		mockService := &MockHPAService{
			platformScalingStatus: nil,
			err:                   errors.New("erreur de test"),
		}

		// Créer le handler avec le mock
		handler := &HPAHandlers{
			hpaService: mockService,
		}

		// Créer une requête HTTP factice
		req, err := http.NewRequest("GET", "/api/v1/platform-scaling", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Créer un ResponseRecorder pour enregistrer la réponse
		rr := httptest.NewRecorder()

		// Appeler le handler
		handler.GetPlatformScalingStatus(rr, req)

		// Vérifier le code de statut
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler a retourné un code de statut incorrect: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}
