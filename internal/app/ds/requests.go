package ds

import (
	"time"
)

type GetSubstancesRequestBody struct {
	Name   string
	Status string
}

type ModConfirm struct {
	Confirm string
}

type OrderSynthesisRequestBody struct {
	User_name             string
	Substances            string `json:"substances"`
	Additional_conditions string
	Status                string
}

type SetSynthesisSubstancesRequestBody struct {
	SynthesisID int
	Substances  string
}

type SynthesesOne struct {
	ID                    int `gorm:"primarykey;AUTO_INCREMENT"`
	Name                  string
	Additional_conditions string
	Status                string
	Date_created          time.Time `gorm:"not null"`
	Date_processed        time.Time
	Date_finished         time.Time
	Moderator             string
	User_name             string
	Time                  string
	Substances            []Substances
}
type ResponseData struct {
	SynthesesChern int
	Substances     []Substances
}
