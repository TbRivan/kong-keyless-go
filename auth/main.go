package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	helpers "github.com/farhanswitch/kong-keyless/helpers/buffer-formatter"
)
type RefreshTokenResponse struct {
	Access_Token string `json:"access_token"`
	Expires_In int32 `json:"expires_in"`
	Refresh_Token string `json:"refresh_token"`
	Refresh_Expires_In int32 `json:"refresh_expires_in"`
	Token_Type string 
	Not_Before_Policy string `json:"not-before-policy"`
	Session_State string
	Scope string
}
type AuthProvider struct{
	authURL string
	refreshURL string
}
var provider *AuthProvider
func (ap AuthProvider) AuthRequest(endpoint string, token string) error{
	payloadBody := helpers.EncodeBuffer([]byte(fmt.Sprintf(`{"endpoint":"%s","token":"%s"}`,endpoint,token)))
	fmt.Println(payloadBody)
	req, err := http.NewRequest(http.MethodPost, ap.authURL,payloadBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type","application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK{
		return fmt.Errorf("%s","Access Forbidden")
	}
	return nil
}
func(ap AuthProvider) RefreshToken(token string) (RefreshTokenResponse, error){
	payloadBody := helpers.EncodeBuffer([]byte(fmt.Sprintf(`{"token":"%s"}`,token)))
	fmt.Println(payloadBody)
	req, err := http.NewRequest(http.MethodPost,ap.refreshURL,payloadBody)
	if err != nil {
		return RefreshTokenResponse{}, err
	}
	req.Header.Add("Content-Type","application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return RefreshTokenResponse{},err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RefreshTokenResponse{},err
	}
	var refResponse RefreshTokenResponse
	fmt.Println(refResponse)
	err = json.Unmarshal(body, &refResponse)
	if err != nil {
		return RefreshTokenResponse{},err
	}
	return refResponse, nil

}

func AuthProviderFactory(authURL string, refreshURL string) AuthProvider{
	if provider == nil {
		provider = &AuthProvider{authURL, refreshURL}
	}
	return *provider
}