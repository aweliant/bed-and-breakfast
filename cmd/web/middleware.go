package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"time"
)

func WriteToConsole(next http.Handler) http.Handler {
	//pretty common to name the var next

	//an anonymous function cast to a handlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

//NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		Name:       "",
		Value:      "",
		Path:       "/", //refer the entire site
		Domain:     "",
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     0,
		Secure:     cfg.InProduction, //https or not
		HttpOnly:   true,
		SameSite:   http.SameSiteLaxMode,
		Raw:        "",
		Unparsed:   nil,
	})
	//use cookies to ensure that the token generated is available on a per page basis
	return csrfHandler
}

//SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
