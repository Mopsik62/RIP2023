package ds

import (
	"gorm.io/datatypes"
)

type Substances struct {
	ID      int    `json:"ID,omitempty" gorm:"primaryKey;AUTO_INCREMENT"`
	Title   string `gorm:"type:varchar(64);not null;unique"`
	Class   string `gorm:"type:varchar(64);not null"`
	Formula string `gorm:"type:varchar(64);not null"`
	Image   string `gorm:"type:bytea;not null"`
	Status  string `gorm:"type:varchar(10);not null"`
}

type Users struct {
	ID            int    `gorm:"primaryKey;AUTO_INCREMENT"`
	Name          string `gorm:"type:varchar(50);not null;unique"`
	Password      string `gorm:"type:varchar(15);not null"`
	Administrator bool   `gorm:"not null"`
}

type Syntheses struct {
	ID                    int `gorm:"primarykey;AUTO_INCREMENT"`
	Name                  string
	Additional_conditions string
	Status                string         `gorm:"not null"`
	Date_created          datatypes.Date `gorm:"not null"`
	Date_processed        datatypes.Date
	Date_finished         datatypes.Date
	Moderator             string `gorm:"not null"`
	User_name             string `gorm:"not null"`
	//Moderator            Users `gorm:"foreignKey:ModeratorID"`
	//User                 Users `gorm:"foreignKey:UserID;not null"`
}
type Synthesis_substance struct {
	Synthesis_ID int    `gorm:"primaryKey"`
	Substance_ID int    `gorm:"primaryKey"`
	Result       string `gorm:"type:varchar(64)"`
	Stage        int
	//Substance   Substances `gorm:"foreignKey:SubstanceID"`
	//Synthesis   Syntheses  `gorm:"foreignKey:SynthesisID"`
}
