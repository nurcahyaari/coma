package model

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ostafen/clover"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ApplicationKey struct {
	Id            string `json:"_id"`
	ApplicationId string `json:"applicationId"`
	StageId       string `json:"stageId"`
	Key           string `json:"key"`
	Salt          string `json:"salt"`
}

func (r *ApplicationKey) GenerateSalt(length int) {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	r.Salt = string(b)
}

func (r *ApplicationKey) GenerateKey() error {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = nil
	claims["stageId"] = r.StageId
	claims["applicationId"] = r.ApplicationId
	tokenStr, err := token.SignedString([]byte(r.Salt))
	if err != nil {
		return err
	}
	r.Key = tokenStr
	return nil
}

func (r ApplicationKey) MapStringInterface() (map[string]interface{}, error) {
	mapStringIntf := make(map[string]interface{})
	j, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, &mapStringIntf)
	if err != nil {
		return nil, err
	}
	return mapStringIntf, nil
}

type ApplicationKeys []ApplicationKey

type FilterApplicationKey struct {
	ApplicationId string
	StageId       string
}

func (f FilterApplicationKey) IsApplicationIdAndStageIdExists() bool {
	return f.ApplicationId != "" && f.StageId != ""
}

func (f FilterApplicationKey) Filter() *clover.Criteria {
	if !f.IsApplicationIdAndStageIdExists() {
		return nil
	}
	criterias := make([]*clover.Criteria, 0)

	if f.ApplicationId != "" {
		criterias = append(criterias, clover.Field("applicationId").Eq(f.ApplicationId))
	}

	if f.StageId != "" {
		criterias = append(criterias, clover.Field("stageId").Eq(f.StageId))
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
