package app

import (
	docs "awesomeProject1/docs"
	"awesomeProject1/internal/app/config"
	"awesomeProject1/internal/app/ds"
	"awesomeProject1/internal/app/dsn"
	"awesomeProject1/internal/app/redis"
	"awesomeProject1/internal/app/repository"
	"awesomeProject1/internal/app/role"
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
	"gorm.io/datatypes"
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
	log.Println("Server start up1")
	//#Услуги - GET список(1), GET одна запись(2), POST добавление(3), PUT изменение(4),
	//DELETE удаление(5), POST добавление в заявку (6)(объединил с сформированием заявки(synthesis(4))
	a.r.GET("substances", a.get_substances)           //(1)
	a.r.GET("substances/:substance", a.get_substance) //(2)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("substances/add", a.add_substance) //(3)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("substances/:substance/edit", a.edit_substance)        //(4)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).DELETE("substances/:substance/delete", a.delete_substance) //(5)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).POST("substances/:substance/add_image", a.add_image)

	//Заявки - GET список(1),
	//GET одна запись (2), PUT изменение(3),
	//PUT сформировать создателем(4), PUT завершить/отклонить модератором(5), DELETE удаление(6)
	//a.r.GET("syntheses", a.get_syntheses)                              //(1)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).GET("syntheses/:synthesis", a.get_synthesis)                   //(2)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("syntheses/:synthesis/edit", a.edit_synthesis)             //(3)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/generate", a.order_synthesis)                   //(4)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("syntheses/:synthesis/apply", a.apply_synthesis)           //(5)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/:synthesis/apply_user", a.apply_synthesis_user) //(5)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).DELETE("syntheses/:synthesis/delete", a.delete_synthesis) //(6)

	//м-м - DELETE удаление из заявки(1), PUT изменение количества/значения в м-м(2)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).DELETE("synthesis_substance/:id1/:id2", a.delete_ss) //(1)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("synthesis_substance/:id1/:id2/edit", a.edit_ss) //(2)

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
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).POST("/logout", a.logout)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).GET("syntheses", a.get_syntheses)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).GET("syntheses/:synthesis", a.get_synthesis)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/generate", a.order_synthesis)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/:synthesis/apply_user", a.apply_synthesis_user) //(5)
	//(4)

	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("substances/:substance/edit", a.edit_substance)        //(4)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).DELETE("substances/:substance/delete", a.delete_substance) //(5)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).POST("substances/:substance/add_image", a.add_image)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("substances/add", a.add_substance)

	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).GET("syntheses", a.get_syntheses)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).GET("syntheses/:synthesis", a.get_synthesis)                   //(2)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("syntheses/:synthesis/edit", a.edit_synthesis) //(3)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/generate", a.order_synthesis)                   //(4)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("syntheses/:synthesis/apply", a.apply_synthesis) //(5)
	//a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin, role.User)).PUT("syntheses/:synthesis/apply_user", a.apply_synthesis_user) //(5)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).DELETE("syntheses/:synthesis/delete", a.delete_synthesis) //(6)

	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).DELETE("synthesis_substance/:id1/:id2", a.delete_ss) //(1)
	a.r.Use(a.WithAuthCheck(role.Moderator, role.Admin)).PUT("synthesis_substance/:id1/:id2/edit", a.edit_ss) //(2)

	a.r.Use(a.WithAuthCheck(role.Admin)).GET("/ping", a.Ping)

	a.r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}

// @Summary Get all existing substances
// @Description Returns all existing substances
// @Tags substances
// @Accept json
// @Produce json
// @Success 200 {} json
// @Param name_pattern query string false "Substances name pattern"
// @Param title query string false "Substances title"
// @Router /substances [get]
func (a *Application) get_substances(c *gin.Context) {
	var name_pattern = c.Query("name_pattern")
	var title = c.Query("title")
	//	var status = c.Query("status")
	var user_id = instance1.GetUserIDAsString()

	response, err := a.repo.GetAllSubstances(title, name_pattern, user_id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      Get substance
// @Description  Returns substance with given name
// @Tags         substances
// @Produce      json
// @Param substance path string true "Substances name"
// @Success      200  {object}  string
// @Router       /substances/{substance} [get]
func (a *Application) get_substance(c *gin.Context) {
	var title = c.Param("substance")

	found_substance, err := a.repo.FindSubstance(title)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, found_substance)

}

// @Summary      Get syntheses
// @Description  Returns list of all syntheses
// @Tags         syntheses
// @Produce      json
// @Success 200 {} json
// @Param date1 query string false "Substances oldest date"
// @Param date2 query string false "Substances newest date"
// @Param status query string false "Substances status"
// @Router       /syntheses [get]
func (a *Application) get_syntheses(c *gin.Context) {
	//var status = c.Query("status")

	_roleNumber, _ := c.Get("role")
	_userUUID, _ := c.Get("userUUID")

	roleNumber := _roleNumber.(role.Role)

	userUUID := _userUUID.(uuid.UUID)

	UserName, err := a.repo.GetUserNameByID(userUUID)

	var date1 = c.Query("date1")
	var date2 = c.Query("date2")
	var status = c.Query("status")

	found_synthesis, err := a.repo.GetAllSynthesis(date1, date2, status, roleNumber, UserName)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, found_synthesis)

}

// @Summary      Get synthesis
// @Description  Returns synthesis with given id
// @Tags         syntheses
// @Produce      json
// @Success      302 {object}  string
// @Param synthesis path string true "Substances name"
// @Router       /syntheses/{synthesis} [get]
func (a *Application) get_synthesis(c *gin.Context) {
	var id = c.Param("synthesis")

	found_synthesis, err := a.repo.FindSynthesis(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, found_synthesis)

}

// @Summary      Adds image
// @Description  Adds image to substance + minio server
// @Tags         substances
// @Produce      json
// @Success      302  {object}  string
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

// @Summary      Adds substance to database
// @Description  Creates a new substance with parameters, specified in json
// @Tags substances
// @Accept json
// @Produce      json
// @Param substance body ds.Substances true "New substance's details"
// @Success      201  {object}  string "Substance created successfully"
// @Router       /substances/add [post]
func (a *Application) add_substance(c *gin.Context) {
	var substance ds.Substances

	if err := c.BindJSON(&substance); err != nil {
		c.String(http.StatusBadRequest, "Can't parse substance\n"+err.Error())
		return
	}

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

// @Summary      Deletes synthesis
// @Description  Changes synthesis status to "Удалён"
// @Tags         syntheses
// @Produce      json
// @Success      302  {object}  string
// @Param synthesis_id path int true "Synthesis id"
// @Router       /syntheses/{synthesis}/delete [delete]
func (a *Application) delete_synthesis(c *gin.Context) {
	synthesis_id, _ := strconv.Atoi(c.Param("synthesis"))

	err := a.repo.LogicalDeleteSynthesis(synthesis_id)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusFound, "Synthesis was successfully deleted")
}

// @Summary      Deletes substance
// @Description  Finds substance by name and changes its status to "Удалён"
// @Tags         substances
// @Accept json
// @Produce      json
// @Success      302  {object}  string
// @Param substance_name path string true "Substances name"
// @Router       /substances/{substance}/delete [delete]
func (a *Application) delete_substance(c *gin.Context) {
	substance_name := c.Param("substance")

	err := a.repo.LogicalDeleteSubstance(substance_name)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusFound, "Substance was successfully deleted")
}

// @Summary      Deletes Synthesis_Substance
// @Description  Finds Synthesis_Substance by ids and remove it
// @Tags         synthesis_substance
// @Produce      json
// @Success      201  {object}  string
// @Param id1 path int true "Synthesis id"
// @Param id2 path int true "Substance id"
// @Router       /synthesis_substance/{id1}/{id2} [delete]
func (a *Application) delete_ss(c *gin.Context) {
	id1, _ := strconv.Atoi(c.Param("id1"))
	id2, _ := strconv.Atoi(c.Param("id2"))
	err := a.repo.DeleteSS(id1, id2)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusFound, "SynthesisSubstance was successfully deleted")
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
func (a *Application) edit_ss(c *gin.Context) {
	id1, _ := strconv.Atoi(c.Param("id1"))
	id2, _ := strconv.Atoi(c.Param("id2"))
	var ss ds.Synthesis_substance

	if err := c.BindJSON(&ss); err != nil {
		c.Error(err)
		return
	}

	err := a.repo.EditSS(ss, id1, id2)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "SynthesisSubstance was successfully edited")
}

// @Summary      Edits substance
// @Description  Finds substance by name and updates its fields
// @Tags         substances
// @Accept json
// @Produce      json
// @Success      302  {object}  string
// @Param substance body ds.Substances true "Edited substance"
// @Param title query string false "Substance name"
// @Router       /substances/{substance}/edit [put]
func (a *Application) edit_substance(c *gin.Context) {
	var substance ds.Substances
	var title = c.Param("substance")
	if err := c.BindJSON(&substance); err != nil {
		c.Error(err)
		return
	}

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
func (a *Application) edit_synthesis(c *gin.Context) {
	var synthesis_body ds.Syntheses

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

	c.String(http.StatusCreated, "Synthesis was successfully edited")

}

// @Summary      Changes synthesis status as moderator
// @Description  Changes synthesis status to any available status
// @Tags         syntheses
// @Accept json
// @Produce      json
// @Success      201  {object}  string
// @Param synthesis_body body ds.Syntheses true "Syntheses body"
// @Param synthesis path int true "Synthesis id"
// @Router       /syntheses/{synthesis}/apply [put]
func (a *Application) apply_synthesis(c *gin.Context) {
	var synthesis_body ds.Syntheses

	var synthesis = c.Param("synthesis")

	if err := c.BindJSON(&synthesis_body); err != nil {
		c.Error(err)
		return
	}

	switch synthesis_body.Status {
	case "В работе":
		synthesis_body.Date_processed = datatypes.Date(time.Now())
	case "Отклонена":
		synthesis_body.Date_processed = datatypes.Date(time.Now())
		synthesis_body.Date_finished = datatypes.Date(time.Now())
	case "Завершена":
		synthesis_body.Date_finished = datatypes.Date(time.Now())
	}

	err := a.repo.ApplySynthesis(synthesis_body, synthesis)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "Synthesis was successfully edited")

}

// @Summary      Changes synthesis status as user
// @Description  Changes synthesis status as allowed to user
// @Tags         syntheses
// @Accept json
// @Produce      json
// @Success      201  {object}  string
// @Param synthesis_body body ds.Syntheses true "Syntheses body"
// @Param synthesis path int true "Synthesis id"
// @Router       /syntheses/{synthesis}/apply_user [put]
func (a *Application) apply_synthesis_user(c *gin.Context) {
	var synthesis_body ds.Syntheses

	var synthesis = c.Param("synthesis")

	if err := c.BindJSON(&synthesis_body); err != nil {
		c.Error(err)
		return
	}

	err := a.repo.ApplySynthesis(synthesis_body, synthesis)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "Synthesis was successfully edited")

}

// @Summary      Order synthesis
// @Description  Creates a new/ find existing synthesis and adds current substances in it
// @Tags syntheses
// @Accept json
// @Produce      json
// @Success      302  {object}  string
// @Param request_body body ds.OrderSynthesisRequestBody true "Ordering request parameters"
// @Router       /syntheses/generate [put]
func (a *Application) order_synthesis(c *gin.Context) {
	var request_body ds.OrderSynthesisRequestBody

	_userUUID, ok := c.Get("userUUID")

	if !ok {
		c.String(http.StatusInternalServerError, "You should login first")

		return
	}

	userUUID := _userUUID.(uuid.UUID)
	UserName, err := a.repo.GetUserNameByID(userUUID)
	request_body.User_name = UserName
	if err := c.BindJSON(&request_body); err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, "Cant' parse json")
		return
	}

	substancesList := strings.Split(request_body.Substances, ",")
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

	//log.Println(ChernId)
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
	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "Can't order synthesis")
		return
	}
	//
	err = a.repo.OrderSynthesis(request_body)
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

// @Summary      Show hello text
// @Description  very very friendly response
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  pingResp
// @Router       /ping/{name} [get]
func (a *Application) Ping(gCtx *gin.Context) {
	name := gCtx.Param("name")
	gCtx.String(http.StatusOK, "Hello %s", name)
}

// @Summary Login into system
// @Description Returns your token
// @Tags auth
// @Produce json
// @Accept json
// @Success 200 {object} loginResp
// @Param request_body body loginReq true "Login request body"
// @Router /login [post]
func (a *Application) login(c *gin.Context) {
	//log.Println("i am here")
	req := &loginReq{}

	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	//log.Println(req.Login)

	user, err := a.repo.GetUserByLogin(req.Login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	log.Println(req.Login + " and" + user.Name)
	log.Println(req.Password + " and" + user.Pass)

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
		log.Println("token")
		if token == nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))

			return
		}
		log.Println("token")
		log.Println(token)

		jwtToken := "test"

		strToken, err := token.SignedString([]byte(jwtToken))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant read str token"))

			return
		}

		c.SetCookie("One-pot-api-token", "Bearer "+strToken, 3600000000000, "", "", true, true)

		c.JSON(http.StatusOK, loginResp{
			ExpiresIn:   3600000000000,
			AccessToken: strToken,
			TokenType:   "Bearer",
		})
	}
	//log.Println("DSADASDSADASDSAD")
	c.AbortWithStatus(http.StatusForbidden)
}

type registerReq struct {
	Name string `json:"name"` // лучше назвать то же самое что login
	Pass string `json:"pass"`
}

type registerResp struct {
	Ok bool `json:"ok"`
}

// @Summary register a new user
// @Description adds a new user to the database
// @Tags auth
// @Produce json
// @Accept json
// @Success 200 {object} registerResp
// @Param request_body body registerReq true "Request body"
// @Router /register [post]
func (a *Application) register(c *gin.Context) {
	req := &registerReq{}
	err := json.NewDecoder(c.Request.Body).Decode(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if req.Pass == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Password should not be empty"))
		return
	}
	if req.Name == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Name should not be empty"))
	}

	err = a.repo.Register(&ds.User{
		UUID: uuid.New(),
		Role: role.Undefined,
		Name: req.Name,
		Pass: generateHashString(req.Pass),
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &registerResp{
		Ok: true,
	})
}

// @Summary Logout
// @Details Deactivates user's current token
// @Tags auth
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
