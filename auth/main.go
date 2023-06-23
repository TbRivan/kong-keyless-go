package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	helpers "github.com/farhanswitch/kong-keyless/helpers/buffer-formatter"
)

type DataResponse struct{
	StatusCode int32 `json:"statusCode"`
	Details RefreshTokenResponse
}
type RefreshTokenResponse struct {
	
	AccessToken string `json:"access_token"`
	ExpiresIn int32 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	RefreshExpiresIn int32 `json:"refresh_expires_in"`
	Message string `json:"message"`
	
}
type ValidateAccessResponse struct{
	StatusCode int32 `json:"statusCode"`
	Details ValidateAccessDetails
}
type ValidateAccessDetails struct{
	Success bool `json:"success"`
	Message string `json:"message"`
}
type AuthProvider struct{
	authURL string
	refreshURL string
}
var provider *AuthProvider
func (ap AuthProvider) AuthRequest(endpoint string, token string) error{
	payloadBody := helpers.EncodeBuffer([]byte(fmt.Sprintf(`{"path":"%s","token":"%s"}`,endpoint,token)))
	req, err := http.NewRequest(http.MethodPost, ap.authURL,payloadBody)
	if err != nil {
		
		log.Printf("Error: %v\n", err)
		return err
	}
	req.Header.Add("Content-Type","application/json")
	client := &http.Client{}
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
	if accessResponse.StatusCode != http.StatusOK{
		return fmt.Errorf("%s","Access Forbidden")
	}
	return nil
}
func(ap AuthProvider) RefreshToken(token string) (DataResponse, error){
	payloadBody := helpers.EncodeBuffer([]byte(fmt.Sprintf(`{"token":"%s"}`,token)))
	req, err := http.NewRequest(http.MethodPost,ap.refreshURL,payloadBody)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return DataResponse{}, err
	}
	req.Header.Add("Content-Type","application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return DataResponse{},err
	}
	defer res.Body.Close()
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return DataResponse{},err
	}
	var refResponse DataResponse
	err = json.Unmarshal(body, &refResponse)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return DataResponse{},err
	}
	if refResponse.StatusCode != http.StatusOK{
		return DataResponse{}, fmt.Errorf("access forbidden. status code: %d", refResponse.StatusCode)
	}
	return refResponse, nil

}

func AuthProviderFactory(authURL string, refreshURL string) AuthProvider{
	if provider == nil {
		provider = &AuthProvider{authURL, refreshURL}
	}
	return *provider
}