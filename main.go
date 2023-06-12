package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/farhanswitch/kong-keyless/auth"
	cookieparser "github.com/farhanswitch/kong-keyless/helpers/cookie-parser"
)
type Config struct {
	AuthUrl string
	RefUrl string
}

func New() interface{} {
	return &Config{}
}

func (c Config) Access (kong *pdk.PDK){
	cookie, err := kong.Request.GetHeader("cookie")
	log.Println(cookie)
	if err != nil{
		log.Printf("There is no cookie")
		kong.Response.Exit(http.StatusForbidden,`{"status":"Access Denied"}`, map[string][]string{
			"Content-Type":{"application/json"},
		})
		return
	}
	var mapCookie map[string]string = cookieparser.ParsingCookie(cookie)
	XACS, ok := mapCookie["XACS"]
	XFRE, ok2 := mapCookie["XFRE"]
	log.Printf("XACS: %s\n",XACS)
	log.Printf("XFRE: %s\n",XFRE)
	authProvider := auth.AuthProviderFactory(c.AuthUrl, c.RefUrl)
	log.Printf("AuthProvider: %v\n", authProvider)
	if !ok {
		log.Printf("There is no ASRF token")
		if !ok2 {
			kong.Response.Exit(http.StatusForbidden,`{"status":"Access Denied"}`, map[string][]string{
				"Content-Type":{"application/json"},
			})
			return

		}
		res, err := authProvider.RefreshToken(XFRE)
		if err != nil {
			kong.Response.Exit(http.StatusForbidden,`{"status":"Access Denied"}`, map[string][]string{
				"Content-Type":{"application/json"},
			})
			return
		}
		fmt.Println(res)
		cookieXACS := http.Cookie{
			Name: "XACS",
			Value: res.Access_Token,
			Path: "/",
			MaxAge: int(res.Expires_In) * 1000,
			Secure: true,
		}
		cookieXFRE := http.Cookie{
			Name: "XFRE",
			Value: res.Refresh_Token,
			Path: "/",
			MaxAge: 5 * 1000 * 3600,
			Secure: true,
		}
		XACS = res.Access_Token
		kong.Response.SetHeader("Set-Cookie", cookieXACS.String())
		kong.Response.SetHeader("Set-Cookie", cookieXFRE.String())

	}

	 path, err := kong.Request.GetPath()
	 log.Printf("Path: %s\n", path)
	 if err != nil {
		log.Printf("Cannot get Destination URL")
		kong.Response.Exit(http.StatusForbidden,`{"status":"Access Denied"}`, map[string][]string{
			"Content-Type":{"application/json"},
		})
		return
	 }
	 
	 err = authProvider.AuthRequest(path, XACS)
	 if err != nil {
		log.Printf("Auth Provider Failed")
		kong.Response.Exit(http.StatusForbidden,`{"status":"Access Denied"}`, map[string][]string{
			"Content-Type":{"application/json"},
		})
		return
	 }

	 
	}
func main() {
	Version := "1.0"
	Priority := 1000
	_ = server.StartServer(New, Version, Priority)
}