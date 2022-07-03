package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/aweliant/bed-and-breakfast/pkg/config"
	"github.com/aweliant/bed-and-breakfast/pkg/handlers"
	"github.com/aweliant/bed-and-breakfast/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var cfg config.AppConfig
var session *scs.SessionManager

func main() {
	//http.ResponseWriter 是一种 io.Writer
	//This is a handler func
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	n, err := fmt.Fprintf(w, "Hello world!")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println("Bytes written:" + string(n))
	//})
	cfg.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true //behaviour after close window
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = cfg.InProduction //cookie is encrypted or not. if yes, https. if no, localhost ok
	cfg.Session = session
	//if cfg.UseCache {
	tc, err := render.GenerateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	cfg.TemplateCache = tc
	cfg.UseCache = true
	render.NewTemplates(&cfg)
	//}
	repo := handlers.NewRepo(&cfg)
	handlers.NewHandlers(repo)

	//这个叫set up routes
	//http.HandleFunc("/", handlers.Repo.Home) //why not just repo.Home
	//http.HandleFunc("/about", handlers.Repo.About)

	//使用pat包进行route
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&cfg),
	}

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	//_ = http.ListenAndServe(portNumber, nil)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
