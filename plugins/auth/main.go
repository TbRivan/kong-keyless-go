package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	helpers "github.com/farhanswitch/kong-keyless/helpers/buffer-formatter"
)

type ValidateAccessResponse struct {
	StatusCode int32 `json:"statusCode"`
	Details    ValidateAccessDetails
}
type ValidateAccessDetails struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type AuthProvider struct {
	authURL    string
}

var provider *AuthProvider

func (ap AuthProvider) AuthRequest(endpoint string, token string) error {
	// Encode to JSON
	payloadBody := helpers.EncodeBuffer([]byte(fmt.Sprintf(`{"path":"%s"}`, endpoint)))
	req, err := http.NewRequest(http.MethodPost, ap.authURL, payloadBody)
	if err != nil {

		log.Printf("Error: %v\n", err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	// Lakukan HTTP Request ke service Autentikasi
	res, err := client.Do(req)
	if err != nil {

		log.Printf("Error: %v\n", err)
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var accessResponse ValidateAccessResponse
	err = json.Unmarshal(body, &accessResponse)
	if err != nil {
		return err
	}
	if accessResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", "Access Forbidden")
	}
	return nil
}

func AuthProviderFactory(authURL string) AuthProvider {
	if provider == nil {
		provider = &AuthProvider{authURL}
	}
	return *provider
}
