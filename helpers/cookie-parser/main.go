package helper_cookie_parser

import (
	"strings"
)

func ParsingCookie(strCookie string)map[string]string{
	var mapCookie map[string]string = make(map[string]string)
	for _, cookie := range strings.Split(strCookie,";"){
		obj := strings.Split(cookie,"=")
		var key string = obj[0]
		var value string = obj[1]
		mapCookie[key] = value
	}
	return mapCookie
}