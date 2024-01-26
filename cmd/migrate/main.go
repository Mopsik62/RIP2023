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
	err := db.AutoMigrate(&ds.User{})

	if err != nil {
		panic(err)
	}
}
