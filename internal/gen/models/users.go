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
	*models.ModelSchema

	Id     *models.IntAttribute
	Active *models.BoolAttribute
	Name   *models.StringAttribute
}

var Users = UsersModel{
	models.NewModelSchema("collectionType", "users", "Users", "user", "users", ""),

	models.NewIntAttribute("id", "Id"),
	models.NewBoolAttribute("active", "Active"),
	models.NewStringAttribute("name", "Name"),
}

func (u UsersModel) Name_() string {
	return "users"
}

func (u UsersModel) New_() models.Entity {
	return &User{
		0,
		true,
		"",
	}
}

func (u UsersModel) All_() []models.Attribute {
	return []models.Attribute{
		Users.Id,
		Users.Active,
		Users.Name,
	}
}

func (u UsersModel) Id_() *models.IntAttribute {
	return Users.Id
}
