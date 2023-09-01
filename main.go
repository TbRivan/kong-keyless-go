package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/farhanswitch/kong-keyless/auth"
)

type Config struct {
	AuthUrl string
	RefUrl  string
}

func New() interface{} {
	return &Config{}
}

func (c Config) Access(kong *pdk.PDK) {
	// Baca Access Token.
	// Access token akan ditempatkan di custom header yang bernama "xacs"
	// JIka berbeda, silahkan disesuaikan
	xacs, err := kong.Request.GetHeader("xacs")

	// Jika ada error saat baca custom header
	if err != nil {
		log.Printf("There is no cookie")
		// Berikan response 401
		kong.Response.Exit(http.StatusUnauthorized, `{"status":"Login first!"}`, map[string][]string{
			"Content-Type": {"application/json"},
		})
		return
	}

	// Baca Refresh Token.
	// Refresh token akan ditempatkan di custom header yang bernama "xfre"
	// JIka berbeda, silahkan disesuaikan
	xfre, err := kong.Request.GetHeader("xfre")
	// Jika ada error saat baca custom header

	if err != nil {
		// Berikan response 401
		log.Printf("There is no cookie")
		kong.Response.Exit(http.StatusUnauthorized, `{"status":"Login Unauthorized"}`, map[string][]string{
			"Content-Type": {"application/json"},
		})
		return
	}

	// Buat instance dari Auth Provider
	authProvider := auth.AuthProviderFactory(c.AuthUrl, c.RefUrl)
	log.Printf("AuthProvider: %v\n", authProvider)
	// Jika Access Token nya kosong
	if xacs == "" {
		log.Printf("There is no ASRF token")
		if xfre == "" {
			kong.Response.Exit(http.StatusUnauthorized, `{"status":"Login first!"}`, map[string][]string{
				"Content-Type": {"application/json"},
			})
			return

		}
		// Coba dapatkan Access Token baru berdasarkan Refresh Token
		res, err := authProvider.RefreshToken(xfre)

		if err != nil {
			log.Printf("ERR REFRESH: %v\n", err)
			kong.Response.Exit(http.StatusForbidden, `{"status":"Access Denied. Cannot Refresh"}`, map[string][]string{
				"Content-Type": {"application/json"},
			})
			return
		}

		// Set Custom Header untuk Response
		kong.Response.SetHeader("xacs", res.Details.AccessToken)
		kong.Response.SetHeader("xacs-exp", fmt.Sprintf("%d", int(res.Details.ExpiresIn)*1000))
		kong.Response.SetHeader("xfre", res.Details.RefreshToken)
		kong.Response.SetHeader("xfre-exp", fmt.Sprintf("%d", int(res.Details.RefreshExpiresIn)*1000))
		kong.Response.SetHeader("Access-Control-Expose-Headers", "xacs,xfre,xacs-exp,xfre-exp")

		xacs = res.Details.AccessToken

		reqHeaders, _ := kong.Request.GetHeaders(50)
		reqHeaders["xacs"] = []string{res.Details.AccessToken}
		kong.ServiceRequest.SetHeaders(reqHeaders)

	}

	// Dapatkan URL yang akan dituju
	path, err := kong.Request.GetPath()
	log.Printf("Path: %s\n", path)
	if err != nil {
		log.Printf("Cannot get Destination URL")
		kong.Response.Exit(http.StatusInternalServerError, `{"status":"Something went wrong!"}`, map[string][]string{
			"Content-Type": {"application/json"},
		})
		return
	}

	// Hit Service Authentikasi untuk verify apakah user dengan token tersebut diperbolehkan access ke resource yang akan dituju
	err = authProvider.AuthRequest(path, xacs)
	if err != nil {
		log.Printf("Auth Provider Failed\n")
		kong.Response.Exit(http.StatusForbidden, `{"status":"Access Denied!"}`, map[string][]string{
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
