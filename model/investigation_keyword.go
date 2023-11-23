package model

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Strings []string

func (s *Strings) Scan(src interface{}) (err error) {
	var synonyms []string
	switch src.(type) {
	case string:
		if src.(string) == "{}" {
			synonyms = []string{}
			break
		}
		// convert this string "{word1,word2,word3}" to []string
		//remove first and last character
		src = src.(string)[1 : len(src.(string))-1]
		// split by comma
		synonyms = strings.Split(src.(string), ",")
	case []byte:
		src = string(src.([]byte))[1 : len(src.([]byte))-1]
		synonyms = strings.Split(string(src.([]byte)), ",")
	case nil:
		synonyms = []string{}
	default:
		return errors.New("Incompatible types")
	}
	if err != nil {
		return
	}
	*s = synonyms
	return nil
}

type InvestigationKeyword struct {
	Id              uuid.UUID `db:"id" json:"id"`
	UserId          uuid.UUID `db:"user_id" json:"userId"`
	InvestigationId uuid.UUID `db:"investigation_id" json:"investigationId"`
	Word            string    `db:"word" json:"word"`
	Synonyms        Strings   `db:"synonyms" json:"synonyms"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
}

func NewInvestigationKeyword(userId uuid.UUID, investigationId uuid.UUID, word string, synonyms []string) *InvestigationKeyword {
	return &InvestigationKeyword{
		Id:              uuid.New(),
		UserId:          userId,
		InvestigationId: investigationId,
		Word:            word,
		Synonyms:        synonyms,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (i *InvestigationKeyword) Update(word string, synonyms []string) {
	i.Word = word
	i.Synonyms = synonyms
	i.UpdatedAt = time.Now()
}
