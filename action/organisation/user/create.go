package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/factly/kavach-server/util"

	"github.com/factly/kavach-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

type invite struct {
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

type role struct {
	Members []string `json:"members"`
}

// create return all user in organisation
func create(w http.ResponseWriter, r *http.Request) {
	organisationID := chi.URLParam(r, "organisation_id")
	orgID, err := strconv.Atoi(organisationID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	var currentUID int
	currentUID, err = strconv.Atoi(r.Header.Get("X-User"))

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	// Check if logged in user is owner
	err = util.CheckOwner(uint(currentUID), uint(orgID))

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	// FindOrCreate invitee
	req := invite{}
	json.NewDecoder(r.Body).Decode(&req)

	validationError := validationx.Check(req)
	if err != nil {
		errorx.Render(w, validationError)
		return
	}

	invitee := model.User{}

	model.DB.FirstOrCreate(&invitee, &model.User{
		Email: req.Email,
	})

	// Check if invitee already exist in organisation
	var totPermissions int
	permission := &model.OrganisationUser{}
	permission.OrganisationID = uint(orgID)
	permission.UserID = invitee.ID

	model.DB.Model(&model.OrganisationUser{}).Where(permission).Count(&totPermissions)

	if totPermissions != 0 {
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	if req.Role == "owner" {
		/* creating policy for admins */
		reqRole := &model.Role{}
		reqRole.Members = []string{fmt.Sprint(invitee.ID)}

		util.UpdateKetoRole(w, "/engines/acp/ory/regex/roles/roles:org:"+fmt.Sprint(orgID)+":admin/members", reqRole)
	}

	// Add user into organisation
	permission.OrganisationID = uint(orgID)
	permission.UserID = invitee.ID
	permission.Role = req.Role

	err = model.DB.Model(&model.OrganisationUser{}).Create(&permission).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	result := &userWithPermission{}

	result.User = invitee
	result.Permission = *permission

	renderx.JSON(w, http.StatusCreated, result)
}
