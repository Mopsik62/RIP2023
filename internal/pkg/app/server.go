package app

import (
	docs "awesomeProject1/docs"
	"awesomeProject1/internal/app/config"
	"awesomeProject1/internal/app/ds"
	"awesomeProject1/internal/app/dsn"
	"awesomeProject1/internal/app/redis"
	"awesomeProject1/internal/app/repository"
	"awesomeProject1/internal/app/role"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// @BasePath /

type Application struct {
	repo   *repository.Repository
	r      *gin.Engine
	config *config.Config
	redis  *redis.Client
}

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	Login       string `json:"login"`
	Role        string `json:"role"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func New(ctx context.Context) (*Application, error) {
	cfg, err := config.NewConfig(ctx)
	if err != nil {
		return nil, err
	}

	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		return nil, err
	}

	return &Application{
		config: cfg,
		repo:   repo,
		redis:  redisClient,
	}, nil
}

type Singleton struct {
	UserID string
}

var instance *Singleton
var once sync.Once

// GetUserContext возвращает единственный экземпляр UserContext
func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{
			UserID: "Nikita1",
		}
	})
	return instance
}

var instance1 = GetInstance()

func (s *Singleton) GetUserIDAsString() string {
	return s.UserID
}

func (a *Application) StartServer() {
	a.r = gin.Default()

	// swagger
	docs.SwaggerInfo.Title = "One-pot syntheses"
	docs.SwaggerInfo.Description = "API SERVER"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = "/"
	a.r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// registration & etc
	a.r.POST("/login", a.login)
	a.r.POST("/register", a.register)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User, role.Undefined)).GET("substances", a.get_substances)           //(1)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User, role.Undefined)).GET("substances/:substance", a.get_substance) //(2)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).POST("/logout", a.logout)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).GET("syntheses", a.get_syntheses)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).GET("syntheses/:synthesis", a.get_synthesis)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/generate", a.order_synthesis)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/:synthesis/set_substances", a.set_synthesis_substances)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/:synthesis/apply_user", a.apply_synthesis_user) //(5)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/:synthesis/edit", a.edit_synthesis)             //(3)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).DELETE("synthesis_substance/:id1/:id2", a.delete_ss)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).DELETE("syntheses/:synthesis/delete", a.delete_synthesis)            //(6)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/:synthesis/set_synthesis_time", a.set_synthesis_time) //(5)

	//(4)

	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("substances/:substance/edit", a.edit_substance)        //(4)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).DELETE("substances/:substance/delete", a.delete_substance) //(5)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).POST("substances/:substance/add_image", a.add_image)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).POST("substances/add", a.add_substance)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("syntheses/:synthesis/apply", a.apply_synthesis) //(5)

	a.r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}

// @Summary Список субстанций
// @Description Вовзращает все субстанции
// @Tags Субстанции
// @Accept json
// @Produce json
// @Success 200 {} json
// @Param name_pattern query string false "Имя субстанции"
// @Param status query string false "Статус субстанции"
// @Router /substances [get]
func (a *Application) get_substances(c *gin.Context) {
	var name_pattern = c.Query("name_pattern")

	var status = c.Query("status")
	//_roleNumber, _ := c.Get("role")

	//roleString := _roleNumber.(role.Role)
	//	var status = c.Query("status")
	_userUUID, _ := c.Get("userUUID")
	//userUUID := _userUUID.(uuid.UUID)
	var userUUID = uuid.Nil

	//log.Println("dsadas")
	if _userUUID != nil {
		userUUID = _userUUID.(uuid.UUID) //log.Println(userUUID)
		//log.Println("NO NILLLLLLLLLLLLL")

	}
	//log.Println(userUUID)

	UserName, err := a.repo.GetUserNameByID(userUUID)
	//var UserName = instance1.GetUserIDAsString()

	response, err := a.repo.GetAllSubstances(name_pattern, UserName, status)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Одна субстанция
// @Description Возвращает одну субстанцию по имени
// @Tags Субстанции
// @Produce json
// @Param substance path string true "Имя субстанции"
// @Success 200 {object} string "Субстанция"
// @Router /substances/{substance} [get]
func (a *Application) get_substance(c *gin.Context) {
	var title = c.Param("substance")

	found_substance, err := a.repo.FindSubstance(title)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, found_substance)

}

// @Summary Список синтезов
// @Description Получает все синтезы
// @Tags Синтезы
// @Produce json
// @Success 200 {} json
// @Param date1 query string false "Первая дата"
// @Param date2 query string false "Вторая дата"
// @Param status query string false "Статус"
// @Param creator query string false "Создатель"
// @Router /syntheses [get]
func (a *Application) get_syntheses(c *gin.Context) {
	//var status = c.Query("status")

	_roleNumber, _ := c.Get("role")
	//_userUUID, _ := c.Get("userUUID")

	roleNumber := _roleNumber.(role.Role)

	_userUUID, _ := c.Get("userUUID")
	userUUID := _userUUID.(uuid.UUID)

	UserName, err := a.repo.GetUserNameByID(userUUID)

	log.Println("userName= " + UserName)
	log.Println("userrole= " + roleNumber)

	var date1 = c.Query("date1")
	var date2 = c.Query("date2")
	var status = c.Query("status")
	var creator = c.Query("creator")

	found_synthesis, err := a.repo.GetAllSynthesis(date1, date2, status, roleNumber, UserName, creator)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, found_synthesis)

}

// @Summary Один синтез
// @Description Возвращает синтез по ID
// @Tags Синтезы
// @Produce json
// @Success 200 {object} string
// @Param synthesis path string true "Synthesis ID"
// @Router /syntheses/{synthesis} [get]
func (a *Application) get_synthesis(c *gin.Context) {
	var id = c.Param("synthesis")

	found_synthesis, err := a.repo.FindSynthesis(id)
	if err != nil {
		c.Error(err)
		return
	}
	//log.Println(found_synthesis)
	c.JSON(http.StatusOK, found_synthesis)

}

// @Summary      Добавить изображение
// @Description  Добавляет изображение к субстанции + на минио сервер
// @Tags         Субстанции
// @Produce      json
// @Success 201 {object} string "Картинка загружена"
// @Param substance path int true "ID Субстанции"
// @Param file formData file true "Изображение"
// @Router       /substances/{substance}/add_image [post]
func (a *Application) add_image(c *gin.Context) {
	substance_id, err := strconv.Atoi(c.Param("substance"))
	if err != nil {
		c.String(http.StatusBadRequest, "Не получается прочитать ID субстанции")
		log.Println("Не получается прочитать ID субстанции")
		return
	}

	image, header, err := c.Request.FormFile("file")

	if err != nil {
		c.String(http.StatusBadRequest, "Не получается распознать картинку")
		log.Println("Не получается распознать картинку")
		return
	}
	defer image.Close()

	minioClient, err := minio.New("127.0.0.1:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Не получается подключиться к minio")
		log.Println("Не получается подключиться к minio")
		return
	}

	objectName := header.Filename
	_, err = minioClient.PutObject(c.Request.Context(), "substances", objectName, image, header.Size, minio.PutObjectOptions{})

	if err != nil {
		c.String(http.StatusInternalServerError, "Не получилось загрузить картинку в minio")
		log.Println("Не получилось загрузить картинку в minio")
		return
	}
	objectName = "http://127.0.0.1:9000/substances/" + header.Filename

	err = a.repo.SetSubstanceImage(substance_id, objectName)

	if err != nil {
		c.String(http.StatusInternalServerError, "Не получается обновить картинку субстанции")
		log.Println("Не получается обновить картинку субстанции")
		return
	}

	c.String(http.StatusCreated, "Картинка загружена!")

}

// @Summary Добавляет субстанцию
// @Description Создает новую субстанцию из паркметров JSON
// @Tags Субстанции
// @Accept json
// @Produce json
// @Param substance body ds.Substances true "Детали новой субстанции"
// @Success 201 {object} string "Substance created successfully"
// @Router /substances/add [post]
func (a *Application) add_substance(c *gin.Context) {
	var substance ds.Substances

	if err := c.BindJSON(&substance); err != nil {
		c.String(http.StatusBadRequest, "Can't parse substance\n"+err.Error())
		return
	}
	log.Println(substance.Title)
	if substance.Image == "" {
		substance.Image = "http://127.0.0.1:9000/substances/default.jpg"
	}
	err := a.repo.CreateSubstance(substance)

	if err != nil {
		c.String(http.StatusNotFound, "Can't create substance\n"+err.Error())
		return
	}

	c.String(http.StatusCreated, "Substance created successfully")
}

// @Summary Удаляет синтез
// @Description Меняет статус синтеза на "Удалён"
// @Tags Синтезы
// @Produce json
// @Success 200 {object} string "Synthesis was successfully deleted"
// @Param synthesis path int true "ID Синтеза"
// @Router /syntheses/{synthesis}/delete [delete]
func (a *Application) delete_synthesis(c *gin.Context) {
	synthesis_id, _ := strconv.Atoi(c.Param("synthesis"))

	err := a.repo.LogicalDeleteSynthesis(synthesis_id)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusOK, "Synthesis was successfully deleted")
}

// @Summary Удалить субстанцию
// @Description Меняет статус субстанции на "Удалён"
// @Tags Субстанции
// @Produce json
// @Success 200 {object} string "Substance was successfully deleted"
// @Param substance path string true "Имя субстанции"
// @Router /substances/{substance}/delete [delete]
func (a *Application) delete_substance(c *gin.Context) {
	substance_name := c.Param("substance")

	err := a.repo.LogicalDeleteSubstance(substance_name)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusFound, "Substance was successfully deleted")
}

// @Summary      Удалить связь синтеза с субстанцией
// @Description  Ищет связь синтеза с субстанцией и удаляет её
// @Tags         Синтезы
// @Produce      json
// @Success      201  {object}  string
// @Param id1 path int true "ID Синтеза"
// @Param id2 path int true "ID Субстанции"
// @Router       /synthesis_substance/{id1}/{id2} [put]
func (a *Application) delete_ss(c *gin.Context) {
	id1, _ := strconv.Atoi(c.Param("id1"))
	id2, _ := strconv.Atoi(c.Param("id2"))
	err := a.repo.DeleteSS(id1, id2)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "SynthesisSubstance was successfully deleted")
}

// @Summary      Edits Synthesis_Substance
// @Description  Finds Synthesis_Substance by ids and edits it
// @Tags         synthesis_substance
// @Accept json
// @Produce      json
// @Success      302  {object}  string
// @Param id1 path int true "Synthesis id"
// @Param id2 path int true "Substance id"
// @Param ss body ds.Synthesis_substance true "Parameters for ss"
// @Router       /synthesis_substance/{id1}/{id2}/edit [put]
//func (a *Application) edit_ss(c *gin.Context) {
//	id1, _ := strconv.Atoi(c.Param("id1"))
//	id2, _ := strconv.Atoi(c.Param("id2"))
//	var ss ds.Synthesis_substance
//
//	if err := c.BindJSON(&ss); err != nil {
//		c.Error(err)
//		return
//	}
//
//	err := a.repo.EditSS(ss, id1, id2)
//
//	if err != nil {
//		c.Error(err)
//		return
//	}
//
//	c.String(http.StatusCreated, "SynthesisSubstance was successfully edited")
//}

// @Summary      Редактировать субстанцию
// @Description  Ищет субстанцию по имени и меняет её
// @Tags         Субстанции
// @Accept json
// @Produce      json
// @Success      302  {object}  string
// @Param substance body ds.Substances true "Отредактированная субстанция"
// @Param title query string false "Имя субстанции"
// @Router       /substances/{substance}/edit [put]
func (a *Application) edit_substance(c *gin.Context) {
	var substance ds.Substances
	var title = c.Param("substance")
	if err := c.BindJSON(&substance); err != nil {
		c.Error(err)
		return
	}
	log.Println(substance)
	err := a.repo.EditSubstance(substance, title)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "Substance was successfully edited")

}

// @Summary      Edits synthesis
// @Description  Finds synthesis and updates it fields
// @Tags         syntheses
// @Accept json
// @Produce      json
// @Success      302  {object}  string
// @Param synthesis_body body ds.Syntheses true "Edited substance"
// @Param synthesis query int false "Substance name"
// @Router       /syntheses/{synthesis}/edit [put]
//func (a *Application) edit_synthesis(c *gin.Context) {
//	var synthesis_body ds.Syntheses
//	log.Println("edit syntheses by mod")
//	var synthesis = c.Param("synthesis")
//
//	if err := c.BindJSON(&synthesis_body); err != nil {
//		c.Error(err)
//		return
//	}
//
//	err := a.repo.EditSynthesis(synthesis_body, synthesis)
//
//	if err != nil {
//		c.Error(err)
//		return
//	}
//
//	c.String(http.StatusCreated, "Synthesis was successfully edited")
//
//}

// @Summary      Редактировать синтез
// @Description  Ищет синтез по ID
// @Tags         Синтезы
// @Accept json
// @Produce      json
// @Success      200  {object}  string
// @Param synthesis_body body ds.Syntheses true "Отредактированный синтез"
// @Param synthesis query int false "ID Синтеза"
// @Router       /syntheses/{synthesis}/edit [put]
func (a *Application) edit_synthesis(c *gin.Context) {
	var synthesis_body ds.Syntheses
	log.Println("edit synthesis by user")
	var synthesis = c.Param("synthesis")

	if err := c.BindJSON(&synthesis_body); err != nil {
		c.Error(err)
		return
	}

	err := a.repo.EditSynthesis(synthesis_body, synthesis)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusOK, "Synthesis was successfully edited")

}

func (a *Application) set_synthesis_substances(c *gin.Context) {
	var id = c.Param("synthesis")
	var requestBody ds.SetSynthesisSubstancesRequestBody

	idInt, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
		return
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.String(http.StatusBadRequest, "Не получается распознать json запрос")
		return
	}
	requestBody.SynthesisID = idInt

	err = a.repo.SetSynthesisSubstances(requestBody.SynthesisID, requestBody.Substances)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не получилось задать субстанции для заявки\n"+err.Error())
	}

	c.String(http.StatusCreated, "Субстанции заявки успешно заданы!")

}

// @Summary      Поменять статус заявки (синтеза) как модератор
// @Description  Меняет статус заявки на выбранный
// @Tags         Синтезы
// @Accept json
// @Produce      json
// @Success      200  {object}  string
// @Param synthesis_body body ds.ModConfirm true "Статус"
// @Param synthesis path int true "ID Синтеза"
// @Router       /syntheses/{synthesis}/apply [put]
func (a *Application) apply_synthesis(c *gin.Context) {
	//var synthesis_body ds.Syntheses
	var confirm ds.ModConfirm
	var synthesis = c.Param("synthesis")

	_userUUID, ok := c.Get("userUUID")
	if !ok {
		c.String(http.StatusInternalServerError, "You should login first")

		return
	}

	userUUID := _userUUID.(uuid.UUID)
	UserName, err := a.repo.GetUserNameByID(userUUID)

	if err = c.BindJSON(&confirm); err != nil {
		c.Error(err)
		return
	}

	if confirm.Confirm == "True" {
		err := a.repo.ApplySynthesis(synthesis, UserName)
		if err != nil {
			c.Error(err)
			return
		}

	} else if confirm.Confirm == "False" {
		err := a.repo.DenySynthesis(synthesis, UserName)
		if err != nil {
			c.Error(err)
			return
		}
	} else if confirm.Confirm == "End" {
		err := a.repo.EndSynthesis(synthesis, UserName)
		if err != nil {
			c.Error(err)
			return
		}
	}

	c.String(http.StatusOK, "Synthesis was successfully edited")

}

// @Summary      Поменять статус синтезу как лаборант
// @Description  Меняет статус как лаборант
// @Tags         Синтезы
// @Accept json
// @Produce      json
// @Success      201  {object}  string
// @Param synthesis path int true "ID Синтеза"
// @Router       /syntheses/{synthesis}/apply_user [put]
func (a *Application) apply_synthesis_user(c *gin.Context) {
	//var synthesis_body ds.Syntheses

	var synthesis = c.Param("synthesis")

	err := a.repo.ApplySynthesisUser(synthesis, "")

	if err != nil {
		c.Error(err)
		return
	}

	jwtStr := c.GetHeader("Authorization")
	jwtPrefix := "Bearer "
	jwtStr = jwtStr[len(jwtPrefix):]

	url := "http://0.0.0.0:3000/syntheses_time/"

	jsonPayload := []byte(`{"pk": "` + synthesis + `", "token": "` + jwtStr + `"}`)

	_, err = http.Post(url, "application/json",
		bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}

	c.String(http.StatusCreated, "Synthesis was successfully edited")

}

// @Summary      Заказать синтез
// @Description  Создаёт новый/находит существующий синтез и добавляет к нему субстанции
// @Tags Синтезы
// @Accept json
// @Produce      json
// @Success 201 {object} string
// @Param request_body body ds.OrderSynthesisRequestBody true "Параметры заказа"
// @Router /syntheses/generate [put]
func (a *Application) order_synthesis(c *gin.Context) {
	var request_body ds.OrderSynthesisRequestBody

	_userUUID, ok := c.Get("userUUID")
	//log.Println("1")
	if !ok {
		c.String(http.StatusInternalServerError, "You should login first")

		return
	}

	userUUID := _userUUID.(uuid.UUID)
	UserName, err := a.repo.GetUserNameByID(userUUID)
	request_body.User_name = UserName
	if err = c.BindJSON(&request_body); err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, "Cant' parse json")
		return
	}

	substancesList := strings.Split(request_body.Substances, ",")

	var intList []int

	for _, substance := range substancesList {
		// Вызываем функцию GetSubstanceIDByName для каждой субстанции
		substanceID, err := a.repo.GetSubstanceIDByName(substance)
		if err != nil {
			continue
		}
		// Добавляем полученный идентификатор в intList
		log.Println(substanceID)
		intList = append(intList, substanceID)
	}

	//log.Println(ChernId)
	if request_body.Status == "Черновик" {
		//проверка есть ли черновая заявка у пользователя
		ChernId, err := a.repo.CheckForChern(request_body.User_name)
		//=0 когда нет черновой заявки => создаём
		if ChernId != 0 {
			var order ds.Synthesis_substance
			var stage = 1
			order.Synthesis_ID = ChernId
			for _, substanceFirst := range intList {
				// Создаем объект SynthesisSubstance
				order.Substance_ID = substanceFirst
				order.Stage = stage
				stage++
				// Вызываем функцию CreateSynthesisSubstance
				err = a.repo.CreateSynthesisSubstance(order)
				if err != nil {
					fmt.Println("Ошибка при создании Synthesis_Substance:", err)
					// Обработайте ошибку по вашему усмотрению
				}
			}

		} else {
			err = a.repo.OrderSynthesis(request_body)

		}
	} else {
		err = a.repo.OrderSynthesis(request_body)
	}
	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "Can't order synthesis")
		return
	}
	//
	//err = a.repo.OrderSynthesis(request_body)
	//
	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "Can't order synthesis")
		return
	}

	c.String(http.StatusCreated, "Synthesis was successfully ordered")

}

type pingReq struct{}
type pingResp struct {
	Status string `json:"status"`
}

// @Summary Войти в систему
// @Description Возвращает jwt токен
// @Tags Аутентификация
// @Produce json
// @Accept json
// @Success 200 {object} loginResp
// @Param request_body body loginReq true "Тело запроса на вход"
// @Router /login [post]
func (a *Application) login(c *gin.Context) {
	req := &loginReq{}

	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	user, err := a.repo.GetUserByLogin(req.Login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if req.Login == user.Name && user.Pass == req.Password {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &ds.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(3600000000000).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "admin",
			},
			UserUUID: user.UUID,
			Scopes:   []string{}, // test data
			Role:     user.Role,
		})
		//log.Println("token")
		if token == nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))

			return
		}
		//log.Println("token")
		//log.Println(token)

		jwtToken := "test"

		strToken, err := token.SignedString([]byte(jwtToken))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant read str token"))

			return
		}

		c.SetCookie("One-pot-api-token", "Bearer "+strToken, 3600000000000, "", "", true, true)
		//log.Println("role= " + user.Role)
		c.JSON(http.StatusOK, loginResp{
			Login:       user.Name,
			Role:        string(user.Role),
			ExpiresIn:   3600000000000,
			AccessToken: strToken,
			TokenType:   "Bearer",
		})
	}
	c.AbortWithStatus(http.StatusForbidden)
}

type registerReq struct {
	Login    string `json:"login"` // лучше назвать то же самое что login
	Password string `json:"password"`
}

type registerResp struct {
	Ok bool `json:"ok"`
}

// @Summary Зарегистрировать нового пользователя
// @Description Добавляет нового пользователя в БД
// @Tags Аутентификация
// @Produce json
// @Accept json
// @Success 200 {object} registerResp
// @Param request_body body registerReq true "Тело запроса"
// @Router /register [post]
func (a *Application) register(c *gin.Context) {
	req := &registerReq{}
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if req.Password == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Password should not be empty"))
		return
	}
	if req.Login == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Name should not be empty"))
	}

	err = a.repo.Register(&ds.User{
		UUID: uuid.New(),
		Role: role.User,
		Name: req.Login,
		Pass: req.Password,
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &registerResp{
		Ok: true,
	})
}

// @Summary Выйти из системы
// @Details Деактивирует токен пользователя
// @Tags Аутентификация
// @Produce json
// @Accept json
// @Success 200
// @Router /logout [post]
func (a *Application) logout(c *gin.Context) {

	jwtStr, cookieErr := c.Cookie("One-pot-api-token")

	if cookieErr != nil {
		log.Println("ОШИБКА КУКИ")
	}

	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test"), nil
	})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)

		return
	}

	err = a.redis.WriteJWTToBlackList(c.Request.Context(), jwtStr, 3600000000000)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func createSignedTokenString() (string, error) {
	privateKey, err := ioutil.ReadFile("demo.rsa")
	if err != nil {
		return "", fmt.Errorf("error reading private key file: %v\n", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("error parsing RSA private key: %v\n", err)
	}

	token := jwt.New(jwt.SigningMethodRS256)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("error signing token: %v\n", err)
	}

	return tokenString, nil
}

type setSynthesisTimeReq struct {
	synthesisId int
	time        string
}

func (a *Application) set_synthesis_time(c *gin.Context) {
	//req := &setSynthesisTimeReq{}
	//log.Println("Зашло")
	var synthesis_id = c.Param("synthesis")
	var time = c.Query("time")
	//log.Println(time)
	//err := json.NewDecoder(c.Request.Body).Decode(req)
	//if err != nil {
	//	c.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}
	num, err := strconv.Atoi(synthesis_id)
	if err != nil {
		fmt.Println("Ошибка при преобразовании строки в число:", err)
		return
	}

	synthesis := ds.Syntheses{}
	synthesis.ID = num
	synthesis.Time = time

	//log.Println(time)

	err = a.repo.EditSynthesis(synthesis, synthesis_id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

}
