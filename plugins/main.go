package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/farhanswitch/kong-keyless/auth"
)

type CustomResponse struct {
	StatusCode int    `json:"statusCode"`
	Details    Details `json:"details"`
}

type Details struct {
	Message string `json:"message"`
}

type Config struct {
	AuthUrl string
}

func New() interface{} {
	return &Config{}
}

func (c Config) Access(kong *pdk.PDK) {
	// Baca Access Token.
	// Access token akan ditempatkan di custom header yang bernama "xacs"
	// JIka berbeda, silahkan disesuaikan
	xacs, err := kong.Request.GetHeader("authorization")

	responseUnauthorized := CustomResponse{
		StatusCode: http.StatusUnauthorized,
		Details: Details{
			Message: "Login first!",
		},
	}
	jsonResponseUnauthorized, err := json.Marshal(responseUnauthorized)

	// Jika ada error saat baca custom header
	if err != nil {
		log.Printf("There is no cookie")
		// Berikan response 401
		kong.Response.Exit(http.StatusOK, string(jsonResponseUnauthorized), map[string][]string{
			"Content-Type": {"application/json"},
		})
		return
	}

	// Buat instance dari Auth Provider
	authProvider := auth.AuthProviderFactory(c.AuthUrl)
	log.Printf("AuthProvider: %v\n", authProvider)
	// Jika Access Token nya kosong
	if xacs == "" {
		log.Printf("There is no Access token")
	}

	responseServerError := CustomResponse{
		StatusCode: http.StatusInternalServerError,
		Details: Details{
			Message: "Something went wrong!",
		},
	}
	jsonResponseServerError, err := json.Marshal(responseServerError)

	// Dapatkan URL yang akan dituju
	path, err := kong.Request.GetPath()
	log.Printf("Path: %s\n", path)
	if err != nil {
		log.Printf("Cannot get Destination URL")
		kong.Response.Exit(http.StatusOK, string(jsonResponseServerError), map[string][]string{
			"Content-Type": {"application/json"},
		})
		return
	}

	responseForbidden := CustomResponse{
		StatusCode: http.StatusForbidden,
		Details: Details{
			Message: "Access Denied!",
		},
	}
	jsonResponseForbidden, err := json.Marshal(responseForbidden)

	// Hit Service Authentikasi untuk verify apakah user dengan token tersebut diperbolehkan access ke resource yang akan dituju
	err = authProvider.AuthRequest(path, xacs)
	if err != nil {
		log.Printf("Auth Provider Failed\n")
		kong.Response.Exit(http.StatusOK, string(jsonResponseForbidden), map[string][]string{
			"Content-Type": {"application/json"},
		})
		return
	}

}
func main() {
	Version := "1.0"
	Priority := 1000
	_ = server.StartServer(New, Version, Priority)
}
