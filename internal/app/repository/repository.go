package repository

import (
	"awesomeProject1/internal/app/ds"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strings"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	log.Println("Application start55555!")
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetProductByID(id int) (*ds.Substances, error) {
	product := &ds.Substances{}
	log.Println("Application start66666!")
	err := r.db.First(product, "id = ?", id).Error // find product with id = 1
	if err != nil {
		return nil, err
	}

	return product, nil
}
func (r *Repository) GetAllSubstances() ([]ds.Substances, error) {
	substances := []ds.Substances{}

	err := r.db.Order("id ASC").Find(&substances, "status = ?", "Активно").Error //err := r.db.First(&substance, "title = ?", name).Error
	//err := r.db.First(&substance, "title = ?", name).Error
	if err != nil {
		return nil, err
	}

	return substances, nil
}

func (r *Repository) SearchSubstances(substance_name string) ([]ds.Substances, error) {
	substances := []ds.Substances{}

	all_substances, all_substances_err := r.GetAllSubstances()

	if all_substances_err != nil {
		return nil, all_substances_err
	}

	for i := range all_substances {
		if strings.Contains(strings.ToLower(all_substances[i].Title), strings.ToLower(substance_name)) {
			substances = append(substances, all_substances[i])
		}
	}

	return substances, nil
}
func (r *Repository) GetSubstanceByName(name string) (*ds.Substances, error) {
	substance := &ds.Substances{}
	err := r.db.First(substance, "title = ?", name).Error
	if err != nil {
		return nil, err
	}

	return substance, nil
}
func (r *Repository) ChangeSubstanceVisibility(substance_name string) error {
	substance, err := r.GetSubstanceByName(substance_name)

	if err != nil {
		log.Println(err)
		return err
	}

	new_status := ""

	if substance.Status == "Активно" {
		new_status = "Неактивно"
	} else {
		new_status = "Активно"
	}

	return r.db.Model(&ds.Substances{}).Where("title = ?", substance_name).Update("status", new_status).Error
}
