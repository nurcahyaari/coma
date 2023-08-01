package entity

import (
	"encoding/json"

	"github.com/ostafen/clover"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `json:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *User) Update(u User) {
	a.Username = u.Username
}

func (a *User) PatchUserPassword(password string) error {
	a.Password = password
	if err := a.HashPassword(); err != nil {
		return err
	}
	return nil
}

func (a *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(a.Password), 12)
	if err != nil {
		return err
	}
	a.Password = string(bytes)
	return nil
}

func (a User) MapStringInterface() (map[string]interface{}, error) {
	mapStringIntf := make(map[string]interface{})
	j, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, &mapStringIntf)
	if err != nil {
		return nil, err
	}
	return mapStringIntf, nil
}

type Users []User

type FilterUser struct {
	Id string
}

func (f *FilterUser) Filter() *clover.Criteria {
	criterias := make([]*clover.Criteria, 0)

	if f.Id != "" {
		criterias = append(criterias, clover.Field("_id").Eq(f.Id))
	}

	filter := &clover.Criteria{}

	if len(criterias) == 0 {
		return nil
	}

	for idx, criteria := range criterias {
		if idx == 0 {
			filter = criteria
			continue
		}

		filter = filter.And(criteria)
	}

	return filter
}
