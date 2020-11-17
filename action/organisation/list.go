package organisation

import (
	"net/http"
	"strconv"

	"github.com/factly/kavach-server/model"
	"github.com/factly/x/renderx"
)

// list - Get all organisations
// @Summary Show all organisations
// @Description Get all organisations
// @Tags Organisation
// @ID get-all-my-organisations
// @Produce  json
// @Param X-User header string true "User ID"
// @Success 200 {array} []orgWithRole
// @Router /organisations/my [get]
func list(w http.ResponseWriter, r *http.Request) {
	organisationUser := make([]model.OrganisationUser, 0)

	userID, err := strconv.Atoi(r.Header.Get("X-User"))

	if err != nil {
		renderx.JSON(w, http.StatusBadRequest, nil)
		return
	}

	model.DB.Model(&model.OrganisationUser{}).Where(&model.OrganisationUser{
		UserID: uint(userID),
	}).Preload("Organisation").Preload("Organisation.Medium").Find(&organisationUser)

	result := make([]orgWithRole, 0)

	for _, each := range organisationUser {
		if each.Organisation != nil {
			eachOrg := orgWithRole{}
			eachOrg.Organisation = *each.Organisation
			eachOrg.Permission = each

			eachOrg.Permission.Organisation = nil

			result = append(result, eachOrg)
		}

	}

	renderx.JSON(w, http.StatusOK, result)
}
