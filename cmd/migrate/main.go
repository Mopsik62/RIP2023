package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"awesomeProject1/internal/app/ds"
	"awesomeProject1/internal/app/dsn"
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
}

func MigrateSubstances(db *gorm.DB) {
	err := db.AutoMigrate(&ds.Substances{})
	if err != nil {
		panic("cant migrate Region to db")
	}
}
