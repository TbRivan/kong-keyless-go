package helper_cookie_parser

import (
	"strings"
)

func ParsingCookie(strCookie string)map[string]string{
	cookieMap := make(map[string]string)

	pairs := strings.Split(strCookie, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		cookie := strings.SplitN(pair, "=", 2)
		if len(cookie) == 2 {
			name := strings.TrimSpace(cookie[0])
			value := strings.TrimSpace(cookie[1])
			cookieMap[name] = value
		}
	}
	return cookieMap
}