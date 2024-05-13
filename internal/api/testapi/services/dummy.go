package services

import (
	"fmt"

	. "github.com/cosys-io/cosys/internal/.gen/models"
	"github.com/cosys-io/cosys/internal/cosys"
)

var DummyService = cosys.NewService(map[string]cosys.Function{
	"findUserByName":  findUserByName,
	"findActiveUsers": findActiveUsers,
})

func findUserByName(cs cosys.Cosys) cosys.ServiceFunction {
	return func(name string) (*User, error) {
		es, err := cs.EntityService()
		if err != nil {
			return nil, err
		}

		users, err := es.FindMany("testapi::users", cosys.ESParam().Filter(Users.Name.Eq(name)))
		if err != nil {
			return nil, err
		}

		if len(users) == 0 {
			return nil, fmt.Errorf("user not found")
		}
		if len(users) > 1 {
			return nil, fmt.Errorf("multiple users found")
		}

		userAsserted, ok := users[0].(*User)
		if !ok {
			return nil, fmt.Errorf("error")
		}

		return userAsserted, nil
	}
}

func findActiveUsers(cs cosys.Cosys) cosys.ServiceFunction {
	return func() ([]*User, error) {
		es, err := cs.EntityService()
		if err != nil {
			return nil, err
		}

		users, err := es.FindMany("testapi::users", cosys.ESParam().Filter(Users.Active))
		if err != nil {
			return nil, err
		}

		usersAsserted := []*User{}
		for _, user := range users {
			userAsserted, ok := user.(*User)
			if !ok {
				return nil, fmt.Errorf("error")
			}

			usersAsserted = append(usersAsserted, userAsserted)
		}

		return usersAsserted, nil
	}
}
