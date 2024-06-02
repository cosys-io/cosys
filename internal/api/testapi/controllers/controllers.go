package controllers

import (
	"github.com/cosys-io/cosys/internal/apis"
)

var Controllers = map[string]*apis.Controller{
	"users": UsersController,
}
