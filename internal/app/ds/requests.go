package ds

import "gorm.io/datatypes"

type GetSubstancesRequestBody struct {
	Name   string
	Status string
}
type OrderSynthesisRequestBody struct {
	User_name  string
	Substances string
}
type SynthesesOne struct {
	ID                    int `gorm:"primarykey;AUTO_INCREMENT"`
	Name                  string
	Additional_conditions string
	Status                string
	Date_created          datatypes.Date `gorm:"not null"`
	Date_processed        datatypes.Date
	Date_finished         datatypes.Date
	Moderator             string
	User_name             string
	Substances            []Substances
}
type ResponseData struct {
	SynthesesChern int
	Substances     []Substances
}
