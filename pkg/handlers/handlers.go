package handlers

import (
	"github.com/aweliant/bed-and-breakfast/pkg/config"
	"github.com/aweliant/bed-and-breakfast/pkg/models"
	"github.com/aweliant/bed-and-breakfast/pkg/render"
	"net/http"
)

//使用了 repository pattern
//Repo the repository used by the handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	Cfg *config.AppConfig
	//gonna add more things later
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{a}
}

//NewHandlers ses the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr                              //a string
	m.Cfg.Session.Put(r.Context(), "remote_ip", remoteIP) //add key value pair to session data

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello there~~~"
	remoteIP := m.Cfg.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{StringMap: stringMap})
}
