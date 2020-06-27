package user

import (
	"github.com/factly/kavach-server/model"
	"github.com/go-chi/chi"
)

type userWithPermission struct {
	model.User
	Permission model.OrganizationUser `json:"permission"`
}

// Router organization
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", create)
	r.Route("/{user_id}", func(r chi.Router) {
		r.Delete("/", delete)
	})

	return r
}
