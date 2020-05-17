package user

import (
	"net/http"
	"strconv"

	"github.com/factly/identity/model"
	"github.com/factly/identity/util/render"
)

// create create organization
func me(w http.ResponseWriter, r *http.Request) {

	userID, _ := strconv.Atoi(r.Header.Get("X-User"))

	me := &model.User{}
	me.ID = uint(userID)

	err := model.DB.Model(&model.User{}).First(&me).Error

	if err != nil {
		return
	}

	render.JSON(w, http.StatusOK, me)
}