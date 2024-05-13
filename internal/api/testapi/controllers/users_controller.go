package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	. "github.com/cosys-io/cosys/internal/.gen/models"
	"github.com/cosys-io/cosys/internal/apis"
	"github.com/cosys-io/cosys/internal/cosys"
)

var UsersController = apis.NewController(map[string]apis.Action{
	"findOne": findOneUser,
	"create":  createUser,
	"update":  updateUser,
	"delete":  deleteUser,
})

func findOneUser(cs cosys.Cosys, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value("query_params").([]string)
		id, _ := strconv.Atoi(ctx[0])

		es, err := cs.EntityService()
		if err != nil {
			InternalErrorHandler(w)
			return
		}

		user, err := es.FindOne("testapi::users", id, cosys.ESParam())
		if err != nil {
			WriteError(w, "Could not find user.", 400)
			return
		}

		WriteJSON(w, user, 200)
	}
}

func createUser(cs cosys.Cosys, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createUserRequest := &struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		}{
			false,
			"",
		}

		if err := json.NewDecoder(r.Body).Decode(createUserRequest); err != nil {
			WriteError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		user := &User{
			Id:     0,
			Active: createUserRequest.Active,
			Name:   createUserRequest.Name,
		}

		es, err := cs.EntityService()
		if err != nil {
			InternalErrorHandler(w)
			return
		}

		newUser, err := es.Create("testapi::users", user, cosys.ESParam())
		if err != nil {
			WriteError(w, "Could not create user.", 400)
			return
		}

		WriteJSON(w, newUser, 200)
	}
}

func updateUser(cs cosys.Cosys, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value("query_params").([]string)
		id, _ := strconv.Atoi(ctx[0])

		updateUserRequest := &struct {
			Active bool   `json:"active"`
			Name   string `json:"name"`
		}{
			false,
			"",
		}

		if err := json.NewDecoder(r.Body).Decode(updateUserRequest); err != nil {
			WriteError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		user := &User{
			Id:     0,
			Active: updateUserRequest.Active,
			Name:   updateUserRequest.Name,
		}

		es, err := cs.EntityService()
		if err != nil {
			InternalErrorHandler(w)
			return
		}

		newUser, err := es.Update("testapi::users", user, id, cosys.ESParam().SetField(Users.Active, Users.Name))
		if err != nil {
			WriteError(w, "Could not update user.", 400)
			return
		}

		WriteJSON(w, newUser, 200)
	}
}

func deleteUser(cs cosys.Cosys, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context().Value("query_params").([]string)
		id, _ := strconv.Atoi(ctx[0])

		es, err := cs.EntityService()
		if err != nil {
			InternalErrorHandler(w)
			return
		}

		oldUser, err := es.Delete("testapi::users", id, cosys.ESParam())
		if err != nil {
			WriteError(w, "Could not delete user.", 400)
			return
		}

		WriteJSON(w, oldUser, 200)
	}
}
