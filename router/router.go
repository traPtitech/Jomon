package router

import (
	"github.com/gorilla/sessions"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
)

type Handlers struct {
	Repo         model.Repository
	Service      service.Service
	SessionName  string
	SessionStore sessions.Store
}
