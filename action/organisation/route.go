package organisation

import (
	"github.com/factly/kavach-server/action/organisation/user"
	"github.com/factly/kavach-server/model"
	"github.com/go-chi/chi"
)

type orgWithRole struct {
	model.Organisation
	Permission model.OrganisationUser `json:"permission"`
}

var userkey model.ContextKey = "user"

// Router organisation
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/my", list)
	r.Post("/", create)
	r.Get("/", all)
	r.Route("/{organisation_id}", func(r chi.Router) {
		r.Get("/", details)
		r.Put("/", update)
		r.Delete("/", delete)
		r.Mount("/users", user.Router())
	})

	return r
}
