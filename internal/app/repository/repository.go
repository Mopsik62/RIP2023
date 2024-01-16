package repository

import (
	"awesomeProject1/internal/app/ds"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
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
func (r *Repository) GetAllSubstances(title string, name_pattern string, user_id string) (ds.ResponseData, error) {
	substances := []ds.Substances{}
	responseData := ds.ResponseData{}
	var tx *gorm.DB = r.db

	ids := []ds.Syntheses{}
	//var synthesesIDs []int

	tx = tx.Where("user_name = ?", user_id)
	tx = tx.Where("status = ?", "Черновик")
	log.Println(tx)
	err := tx.Find(&ids).Error

	//	log.Println(user_id)

	for _, synthesis := range ids {
		synthesisAll := synthesis
		//log.Println("СИнтезис ид = ")
		//log.Println(substanceID)
		responseData.SynthesesChern = append(responseData.SynthesesChern, synthesisAll)
	}
	tx = r.db
	if name_pattern != "" {
		tx = tx.Where("title like ?", "%"+name_pattern+"%")

	}
	//log.Println(responseData.SynthesesIDs)
	if title != "" {
		tx = tx.Where("title = ?", title)
	}

	tx = tx.Where("status = ?", "Активно") //тут скорее всего 4-я сломается т.к. добавил статус пустой

	err = tx.Find(&substances).Error
	responseData.Substances = substances
	//log.Println(substances)
	if err != nil {
		return responseData, err
	}

	return responseData, nil
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
func (r *Repository) GetAllSynthesis(date string, status string) ([]ds.Syntheses, error) {
	var synthesis []ds.Syntheses
	var tx *gorm.DB = r.db
	log.Println(date)
	if date == "True" {
		tx = r.db.Order("date_created")
	}

	tx = tx.Where("status NOT IN (?, ?)", "Удалён", "Черновик") //тут скорее всего 4-я сломается т.к. добавил статус пустой
	//tx = tx.Where("Status NOT IN (?, ?)", "Удалён", "Черновик")
	if status != "" {
		tx = tx.Where("status = ?", status)
	}

	err := tx.Find(&synthesis).Error

	if err != nil {
		return nil, err
	}

	return synthesis, nil
}
func (r *Repository) FindSynthesis(id string) (ds.SynthesesOne, error) {
	synthesis := ds.Syntheses{}
	Answer := ds.SynthesesOne{}
	var tx *gorm.DB = r.db
	tx = tx.Where("ID = ?", id)

	err := tx.Find(&synthesis).Error
	if err != nil {
		return ds.SynthesesOne{}, err
	}

	ids := []ds.Synthesis_substance{}

	var substanceIDs []int
	tx = r.db
	tx = tx.Where("synthesis_id = ?", id)
	err = tx.Find(&ids).Error
	for _, synthesisSubstance := range ids {
		substanceID := synthesisSubstance.Substance_ID
		substanceIDs = append(substanceIDs, substanceID)
	}

	substances := []ds.Substances{}
	tx = r.db
	tx = tx.Where("id IN (?)", substanceIDs)
	err = tx.Find(&substances).Error
	if err != nil {
		return ds.SynthesesOne{}, err
	}
	//substancesAll := ds.Substances{}
	//tx = r.db
	//tx = tx.Where("id = ?", id)
	//err = tx.Find(&substancesAll).Error
	//for _, synthesisSubstance := range ids {
	//	substanceID := synthesisSubstance.Substance_ID
	//	substanceIDs = append(substanceIDs, substanceID)
	//}
	Answer.ID = synthesis.ID
	Answer.Date_processed = synthesis.Date_processed
	Answer.Name = synthesis.Name
	Answer.User_name = synthesis.User_name
	Answer.Date_created = synthesis.Date_created
	Answer.Status = synthesis.Status
	Answer.Date_processed = synthesis.Date_processed
	Answer.Date_finished = synthesis.Date_finished
	Answer.Moderator = synthesis.Moderator
	Answer.Additional_conditions = synthesis.Additional_conditions
	Answer.Substances = substances
	return Answer, nil
}
func (r *Repository) FindSubBySyn(id string) ([]ds.Substances, error) {
	substances := []ds.Substances{}
	ids := []ds.Synthesis_substance{}
	var substanceIDs []int

	var tx *gorm.DB = r.db
	tx = tx.Where("synthesis_id = ?", id)
	err := tx.Find(&ids).Error
	for _, synthesisSubstance := range ids {
		substanceID := synthesisSubstance.Substance_ID

		substanceIDs = append(substanceIDs, substanceID)
	}

	//tx = r.db
	//tx = tx.Where("id IN (?)", substanceIDs)
	//err = tx.Find(&substances).Error
	if err != nil {
		return nil, err
	}

	return substances, nil
}
func (r *Repository) CreateSubstance(substance ds.Substances) error {

	return r.db.Create(&substance).Error
}
func (r *Repository) LogicalDeleteSynthesis(synthesis_id int) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ?", synthesis_id).Update("status", "Удалён").Error
}
func (r *Repository) EditSubstance(substance ds.Substances, title string) error {
	return r.db.Model(&ds.Substances{}).Where("title = ?", title).Updates(substance).Error
}
func (r *Repository) EditSynthesis(synthesis ds.Syntheses, id string) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ?", id).Updates(synthesis).Error
}

func (r *Repository) GenerateSynthesis(synthesis ds.Syntheses, id string) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ? AND status = ?", id, "Черновик").Updates(synthesis).Error
}
func (r *Repository) ApplySynthesis(synthesis ds.Syntheses, id string) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ?", id).Updates(synthesis).Error
}
func (r *Repository) CreateSynthesisSubstance(synthesis_substance ds.Synthesis_substance) error {
	return r.db.Create(&synthesis_substance).Error
}
func (r *Repository) LogicalDeleteSubstance(substance_name string) error {
	return r.db.Model(&ds.Substances{}).Where("title = ?", substance_name).Updates(map[string]interface{}{"status": "Удалён", "additional_conditions": "dads", "date_finished": datatypes.Date(time.Now())}).Error
}
func (r *Repository) DeleteSS(id1 int, id2 int) error {
	return r.db.Model(&ds.Synthesis_substance{}).Where("synthesis_id = ? AND substance_id = ?", id1, id2).Delete(&ds.Synthesis_substance{}).Error
}
func (r *Repository) EditSS(ss ds.Synthesis_substance, id1 int, id2 int) error {
	return r.db.Model(&ds.Synthesis_substance{}).Where("synthesis_id = ? AND substance_id = ?", id1, id2).Updates(ss).Error
}
func (r *Repository) OrderSynthesis(requestBody ds.OrderSynthesisRequestBody) error {
	user_id := requestBody.User_id
	substancesStr := requestBody.Substances
	substancesList := strings.Split(substancesStr, ",")

	var intList []int

	// Проходим по каждой строке в списке substancesList
	for _, str := range substancesList {
		// Преобразуем строку в целое число
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("Ошибка преобразования строки в число:", err)
			// Обработайте ошибку по вашему усмотрению
			continue
		}

		// Добавляем целое число в список intList
		intList = append(intList, num)
	}

	current_date := datatypes.Date(time.Now())

	synthesis := ds.Syntheses{}
	synthesis.Status = "Черновик"
	synthesis.Name = "test"
	synthesis.User_name = user_id
	synthesis.Date_created = current_date
	//synthesis.Date_processed = nil
	err := r.db.Omit("moderator", "date_processed", "date_finished").Create(&synthesis).Error
	//err := r.db.Omit("moderator").Create(&synthesis).Error
	//log.Println(synthesis.ID)
	if err != nil {
		return err
	}

	synthesis_substance := ds.Synthesis_substance{}
	synthesis_substance.Synthesis_ID = synthesis.ID
	var stage = 1
	// Проходим по каждому числу в списке
	for _, substanceFirst := range intList {
		// Создаем объект SynthesisSubstance
		synthesis_substance.Substance_ID = substanceFirst
		synthesis_substance.Stage = stage
		stage++
		// Вызываем функцию CreateSynthesisSubstance
		err = r.CreateSynthesisSubstance(synthesis_substance)
		if err != nil {
			fmt.Println("Ошибка при создании Synthesis_Substance:", err)
			// Обработайте ошибку по вашему усмотрению
		}
	}
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
