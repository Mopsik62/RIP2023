package repository

import (
	"awesomeProject1/internal/app/ds"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}
func (r *Repository) GetAllSubstances(title string, status string) ([]ds.Substances, error) {
	substances := []ds.Substances{}

	var tx *gorm.DB = r.db

	if title != "" {
		tx = tx.Where("title = ?", title)
	}
	if status != "" {
		tx = tx.Where("status = ?", status)
	}

	err := tx.Find(&substances).Error

	if err != nil {
		return nil, err
	}

	return substances, nil
}
func (r *Repository) GetDraft(status string) ([]ds.Syntheses, error) {
	syntheses := []ds.Syntheses{}

	var tx *gorm.DB = r.db

	if status != "" {
		tx = tx.Where("status = ?", status)
	}

	err := tx.Find(&syntheses).Error

	if err != nil {
		return nil, err
	}

	return syntheses, nil
}
func (r *Repository) FindSubstance(title string) ([]ds.Substances, error) {
	substance := []ds.Substances{}
	var tx *gorm.DB = r.db
	if title != "" {
		tx = tx.Where("title = ?", title)
	}
	log.Println("----- " + title)
	err := tx.Find(&substance).Error
	if err != nil {
		return nil, err
	}

	return substance, nil
}
func (r *Repository) GetAllSynthesis() ([]ds.Syntheses, error) {
	var synthesis []ds.Syntheses

	var tx *gorm.DB = r.db.Order("status, date_created")
	err := tx.Find(&synthesis).Error

	if err != nil {
		return nil, err
	}

	return synthesis, nil
}
func (r *Repository) FindSynthesis(title string) ([]ds.Syntheses, error) {
	synthesis := []ds.Syntheses{}
	var tx *gorm.DB = r.db
	if title != "" {
		tx = tx.Where("name = ?", title)
	}

	err := tx.Find(&synthesis).Error
	if err != nil {
		return nil, err
	}

	return synthesis, nil
}
func (r *Repository) CreateSubstance(substance ds.Substances) error {
	return r.db.Create(&substance).Error
}
func (r *Repository) LogicalDeleteSynthesis(synthesis_id int) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ?", synthesis_id).Update("status", "Удалён").Error
}
func (r *Repository) EditSubstance(substance ds.Substances) error {
	return r.db.Model(&ds.Substances{}).Where("title = ?", substance.Title).Updates(substance).Error
}
func (r *Repository) EditSynthesis(synthesis ds.Syntheses) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ?", synthesis.ID).Updates(synthesis).Error
}
func (r *Repository) CreateSynthesisSubstance(synthesis_substance ds.Synthesis_substance) error {
	return r.db.Create(&synthesis_substance).Error
}
func (r *Repository) LogicalDeleteSubstance(substance_id int) error {
	return r.db.Model(&ds.Substances{}).Where("id = ?", substance_id).Update("status", "Удалён").Error
}
func (r *Repository) OrderSynthesis(requestBody ds.OrderSynthesisRequestBody) error {
	user_id := requestBody.User_id
	var substance_first, substance_second int
	substance_first = requestBody.Substance_first
	substance_second = requestBody.Substance_second

	current_date := datatypes.Date(time.Now())

	synthesis := ds.Syntheses{}
	synthesis.Status = "Создана"
	synthesis.Name = "test"
	synthesis.User_ID = user_id
	synthesis.Date_created = current_date

	err := r.db.Omit("moderator_id").Create(&synthesis).Error
	log.Println(synthesis.ID)
	if err != nil {
		return err
	}

	synthesis_substance := ds.Synthesis_substance{}
	synthesis_substance.Synthesis_ID = synthesis.ID
	synthesis_substance.Substance_ID = substance_first
	synthesis_substance.Temperature = "0"
	err = r.CreateSynthesisSubstance(synthesis_substance)
	synthesis_substance.Substance_ID = substance_second
	err = r.CreateSynthesisSubstance(synthesis_substance)
	return err

}
func (r *Repository) FindSubstanceOrder(id string) ([]ds.Synthesis_substance, error) {
	substances := []ds.Synthesis_substance{}
	var tx *gorm.DB = r.db
	if id != "" {
		tx = tx.Where("synthesis_id = ?", id)
	}
	log.Println(id)
	err := tx.Find(&substances).Error
	if err != nil {
		return nil, err
	}

	return substances, nil
}
func (r *Repository) CreateUser(user ds.Users) error {
	return r.db.Create(&user).Error
}
func (r *Repository) OrderAdd(substance ds.Synthesis_substance) error {
	return r.db.Create(&substance).Error
}
