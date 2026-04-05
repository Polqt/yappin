package util

import (
	"net/http"
)

func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	env := GetEnv("ENVIRONMENT", "development")
	cookieDomain := GetEnv("COOKIE_DOMAIN", "")

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	}

	if env == "production" {
		if cookieDomain != "" {
			cookie.Domain = cookieDomain
		}
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(w, cookie)
}

func ClearSecureCookie(w http.ResponseWriter, name string) {
	env := GetEnv("ENVIRONMENT", "development")
	cookieDomain := GetEnv("COOKIE_DOMAIN", "")

	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	if env == "production" {
		if cookieDomain != "" {
			cookie.Domain = cookieDomain
		}
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}

	http.SetCookie(w, cookie)
}
