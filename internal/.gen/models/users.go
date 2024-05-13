// AUTO-GENERATED
// DO NOT CHANGE

package models

import (
	"github.com/cosys-io/cosys/internal/models"
)

type User struct {
	Id     int    `json:"id"`
	Active bool   `json:"active"`
	Name   string `json:"name"`
}

type UsersModel struct {
	Id     *models.IntAttribute
	Active *models.BoolAttribute
	Name   *models.StringAttribute
}

var Users = UsersModel{
	models.NewIntAttribute("id", "Id"),
	models.NewBoolAttribute("active", "Active"),
	models.NewStringAttribute("name", "Name"),
}

func (u UsersModel) Model_Name() string {
	return "users"
}

func (u UsersModel) Model_New() models.Entity {
	return &User{
		0,
		true,
		"",
	}
}

func (u UsersModel) Model_All() []models.IAttribute {
	return []models.IAttribute{
		Users.Id,
		Users.Active,
		Users.Name,
	}
}

func (u UsersModel) Model_Id() *models.IntAttribute {
	return Users.Id
}
