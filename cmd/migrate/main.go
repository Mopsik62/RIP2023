package main

import (
	"awesomeProject1/internal/app/ds"
	"awesomeProject1/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	MigrateSchema(db)

}
func MigrateSchema(db *gorm.DB) {
	MigrateSubstances(db)
	MigrateSyntheses(db)
	MigrateSynSub(db)
	MigrateUsers(db)
}

func MigrateSubstances(db *gorm.DB) {
	err := db.AutoMigrate(&ds.Substances{})
	if err != nil {
		panic("cant migrate Region to db")
	}
}

func MigrateSyntheses(db *gorm.DB) {
	err := db.AutoMigrate(&ds.Syntheses{})
	if err != nil {
		panic("cant migrate Sythesis to db")
	}
	//log.Println(ds.Synthesis{ID: 1})
	//	log.Println("--------------------------")
}

func MigrateSynSub(db *gorm.DB) {
	err := db.AutoMigrate(&ds.Synthesis_substance{})
	if err != nil {
		panic("cant migrate SynSub to db")
	}

}

func MigrateUsers(db *gorm.DB) {
	err := db.AutoMigrate(&ds.Users{})
	if err != nil {
		panic("cant migrate Users to db")
	}
}
