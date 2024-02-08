package repository

import (
	"awesomeProject1/internal/app/ds"
	"awesomeProject1/internal/app/role"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
func (r *Repository) GetUserByID(id uuid.UUID) (*ds.User, error) {
	user := &ds.User{}

	err := r.db.First(user, "UUID = ?", id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetSubstanceIDByName(name string) (int, error) {
	substance := &ds.Substances{}

	err := r.db.First(substance, "Title = ?", name).Error
	if err != nil {
		return 0, err
	}

	return substance.ID, nil
}

func (r *Repository) GetUserNameByID(id uuid.UUID) (string, error) {
	user := &ds.User{}

	err := r.db.First(user, "UUID = ?", id).Error
	if err != nil {

		return "", err
	}

	return user.Name, nil
}

func (r *Repository) GetUserID(name string) (uuid.UUID, error) {
	user := &ds.User{}

	err := r.db.First(user, "name = ?", name).Error
	if err != nil {

		return uuid.Nil, err
	}

	return user.UUID, nil
}

func (r *Repository) GetUserRole(name string) (role.Role, error) {
	user := &ds.User{}

	err := r.db.First(user, "name = ?", name).Error
	if err != nil {
		return role.Undefined, err
	}
	return user.Role, nil
}

func (r *Repository) Register(user *ds.User) error {
	if user.UUID == uuid.Nil {
		user.UUID = uuid.New()
	}

	return r.db.Create(user).Error
}

func (r *Repository) GetAllSubstances(name_pattern string, user_id string, status string) (ds.ResponseData, error) {
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
		synthesis1 := synthesis.ID
		//log.Println("СИнтезис ид = ")
		//log.Println(substanceID)
		if user_id != "00000000-0000-0000-0000-000000000000" {
			responseData.SynthesesChern = synthesis1

		} else {
			responseData.SynthesesChern = 0
		}
	}
	tx = r.db
	if name_pattern != "" {
		tx = tx.Where("title like ?", "%"+name_pattern+"%")

	}
	//log.Println(responseData.SynthesesIDs)

	if status != "All" {
		tx = tx.Where("status = ?", "Активно")
	} //тут скорее всего 4-я сломается т.к. добавил статус пустой

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
func (r *Repository) FindSubstance(title string) (ds.Substances, error) {
	substance := ds.Substances{}
	var tx *gorm.DB = r.db
	if title != "" {
		tx = tx.Where("title = ?", title)
	}
	log.Println("----- " + title)
	err := tx.Find(&substance).Error
	if err != nil {
		return substance, err
	}

	return substance, nil
}
func (r *Repository) GetAllSynthesis(date1 string, date2 string, status string, roleNumber role.Role, UserName string, creator string) ([]ds.Syntheses, error) {
	synthesis := []ds.Syntheses{}

	var tx *gorm.DB = r.db
	//log.Println(date)
	if date1 != "" {
		tx = tx.Where("date_created >= ?", date1)
	}

	if date2 != "" {
		tx = tx.Where("date_created <= ?", date2)
	}
	if roleNumber != role.User {
		tx = tx.Where("status NOT IN (?, ?)", "Удалён", "Черновик") //тут скорее всего 4-я сломается т.к. добавил статус пустой
	}
	//tx = tx.Where("Status NOT IN (?, ?)", "Удалён", "Черновик")
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	if creator != "" {
		tx = tx.Where("user_name=?", creator)
	}
	if roleNumber == role.User {
		tx = tx.Where("user_name = ?", UserName)
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
	Answer.Time = synthesis.Time
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
func (r *Repository) CheckForChern(user_name string) (int, error) {
	syntheses := []ds.Syntheses{}

	var tx *gorm.DB = r.db
	tx = tx.Where("status = ? AND user_name = ?", "Черновик", user_name)

	err := tx.Find(&syntheses).Error

	var ChernId int
	for _, synthesis := range syntheses {
		ChernId = synthesis.ID
	}

	//tx = r.db
	//tx = tx.Where("id IN (?)", substanceIDs)
	//err = tx.Find(&substances).Error
	if err != nil {
		return -1, err
	}

	return ChernId, nil
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
	return r.db.Model(&ds.Syntheses{}).Where("id= ?", id).Updates(synthesis).Error
}

func (r *Repository) GenerateSynthesis(synthesis ds.Syntheses, id string) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ? AND status = ?", id, "Черновик").Updates(synthesis).Error
}
func (r *Repository) ApplySynthesis(id string, userName string) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ?", id).Updates(map[string]interface{}{"status": "В работе", "date_processed": datatypes.Date(time.Now()), "moderator": userName}).Error
}
func (r *Repository) DenySynthesis(id string, userName string) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ?", id).Updates(map[string]interface{}{"status": "Отклонён", "date_finished": datatypes.Date(time.Now()), "moderator": userName}).Error
}

func (r *Repository) EndSynthesis(id string, userName string) error {
	return r.db.Model(&ds.Syntheses{}).Where("id = ?", id).Updates(map[string]interface{}{"status": "Завершён", "date_finished": datatypes.Date(time.Now()), "moderator": userName}).Error
}

func (r *Repository) CreateSynthesisSubstance(synthesis_substance ds.Synthesis_substance) error {
	return r.db.Create(&synthesis_substance).Error
}
func (r *Repository) SetSubstanceImage(id int, image string) error {
	return r.db.Model(&ds.Substances{}).Where("id = ?", id).Update("image", image).Error
}
func (r *Repository) LogicalDeleteSubstance(substance_name string) error {
	return r.db.Model(&ds.Substances{}).Where("title = ?", substance_name).Updates(map[string]interface{}{"status": "Удалён", "additional_conditions": "", "date_finished": datatypes.Date(time.Now())}).Error
}
func (r *Repository) DeleteSS(id1 int, id2 int) error {
	return r.db.Model(&ds.Synthesis_substance{}).Where("synthesis_id = ? AND substance_id = ?", id1, id2).Delete(&ds.Synthesis_substance{}).Error
}
func (r *Repository) EditSS(ss ds.Synthesis_substance, id1 int, id2 int) error {
	return r.db.Model(&ds.Synthesis_substance{}).Where("synthesis_id = ? AND substance_id = ?", id1, id2).Updates(ss).Error
}
func (r *Repository) OrderSynthesis(requestBody ds.OrderSynthesisRequestBody) error {
	user_id := requestBody.User_name
	substancesStr := requestBody.Substances
	substancesList := strings.Split(substancesStr, ",")

	var intList []int

	for _, substance := range substancesList {
		// Вызываем функцию GetSubstanceIDByName для каждой субстанции
		substanceID, err := r.GetSubstanceIDByName(substance)
		if err != nil {
			continue
		}
		// Добавляем полученный идентификатор в intList
		log.Println(substanceID)
		intList = append(intList, substanceID)
	}

	current_date := datatypes.Date(time.Now())

	synthesis := ds.Syntheses{}
	synthesis.Status = requestBody.Status
	synthesis.Additional_conditions = requestBody.Additional_conditions
	synthesis.Name = ""
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

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	user := &ds.User{}

	err := r.db.First(user, "name = ?", login).Error
	if err != nil {
		return nil, err
	}
	log.Println("User= ")
	log.Println(user)
	return user, nil
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

func (r *Repository) SetSynthesisSubstances(synthesisID int, substances string) error {

	substancesList := strings.Split(substances, ",")
	var intList []int

	for _, substance := range substancesList {
		// Вызываем функцию GetSubstanceIDByName для каждой субстанции
		substanceID, err := r.GetSubstanceIDByName(substance)
		if err != nil {
			continue
		}
		// Добавляем полученный идентификатор в intList
		log.Println(substanceID)
		intList = append(intList, substanceID)
	}

	synthesis_substance := ds.Synthesis_substance{}

	synthesis_substance.Synthesis_ID = synthesisID

	ids := []ds.Synthesis_substance{}

	var substanceIDs []int
	var tx *gorm.DB = r.db
	tx = tx.Where("synthesis_id = ?", synthesisID)
	err := tx.Find(&ids).Error
	for _, synthesisSubstance := range ids {
		substanceID := synthesisSubstance.Substance_ID
		substanceIDs = append(substanceIDs, substanceID)
	}
	for _, substanceID := range substanceIDs {
		r.DeleteSS(synthesisID, substanceID)
	}

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

	return nil
}
