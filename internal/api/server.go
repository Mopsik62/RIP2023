package app

import (
	"awesomeProject1/internal/app/ds"
	"awesomeProject1/internal/app/dsn"
	"awesomeProject1/internal/app/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Application struct {
	repo repository.Repository
	r    *gin.Engine
}

func New() Application {
	app := Application{}

	repo, _ := repository.New(dsn.FromEnv())

	app.repo = *repo

	return app

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
	//DELETE удаление(5), POST добавление в заявку (6)
	a.r.GET("substances", a.get_substances)                        //(1)
	a.r.GET("substances/:substance", a.get_substance)              //(2)
	a.r.POST("substances/add", a.add_substance)                    //(3)
	a.r.PUT("substances/:substance/edit", a.edit_substance)        //(4)
	a.r.DELETE("substances/:substance/delete", a.delete_substance) //(5)
	a.r.POST("substances/order_add", a.order_add)                  //(6)
	//Заявки - GET список(1),
	//GET одна запись (2), PUT изменение(3),
	//PUT сформировать создателем(4), PUT завершить/отклонить модератором(5), DELETE удаление(6)
	a.r.GET("syntheses", a.get_syntheses)                         //(1)
	a.r.GET("syntheses/:synthesis", a.get_synthesis)              //(2)
	a.r.PUT("syntheses/:synthesis/edit", a.edit_synthesis)        //(3)
	a.r.PUT("syntheses/generate", a.order_synthesis)              //(4)
	a.r.PUT("syntheses/:synthesis/apply", a.apply_synthesis)      //(5)
	a.r.DELETE("syntheses/:synthesis/delete", a.delete_synthesis) //(6)
	//a.r.GET("syntheses/:synthesis/substances", a.get_SubBySyn)

	//м-м - DELETE удаление из заявки(1), PUT изменение количества/значения в м-м(2)
	a.r.DELETE("synthesis_substance/:id1/:id2", a.delete_ss) //(1)
	a.r.PUT("synthesis_substance/:id1/:id2/edit", a.edit_ss) //(2)

	//a.r.GET("synthesis/draft", a.synthesis_draft)
	//a.r.GET("order/:synthesis_id", a.order_find_substances)
	//a.r.PUT("user/add", a.add_user)
	//a.r.PUT("order", a.order_synthesis)
	//a.r.PUT("order/add", a.order_add)

	a.r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
func (a *Application) synthesis_draft(c *gin.Context) {
	//var name_pattern = c.Query("name_pattern")
	var status = "Черновик"

	syntheses, err := a.repo.GetDraft(status)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, syntheses)
}
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

	c.JSON(http.StatusFound, response)
}
func (a *Application) get_substance(c *gin.Context) {
	var title = c.Param("substance")

	found_substance, err := a.repo.FindSubstance(title)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, found_substance)

}
func (a *Application) get_syntheses(c *gin.Context) {
	//var status = c.Query("status")
	var date = c.Query("date")
	var status = c.Query("status")

	found_synthesis, err := a.repo.GetAllSynthesis(date, status)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, found_synthesis)

}
func (a *Application) get_synthesis(c *gin.Context) {
	var id = c.Param("synthesis")

	found_synthesis, err := a.repo.FindSynthesis(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, found_synthesis)

}

//	func (a *Application) get_SubBySyn(c *gin.Context) {
//		var id = c.Param("synthesis")
//
//		found_substances, err := a.repo.FindSubBySyn(id)
//		if err != nil {
//			c.Error(err)
//			return
//		}
//
//		c.JSON(http.StatusFound, found_substances)
//
// }
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
func (a *Application) delete_synthesis(c *gin.Context) {
	synthesis_id, _ := strconv.Atoi(c.Param("synthesis"))

	err := a.repo.LogicalDeleteSynthesis(synthesis_id)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusFound, "Synthesis was successfully deleted")
}
func (a *Application) delete_substance(c *gin.Context) {
	substance_name := c.Param("substance")

	err := a.repo.LogicalDeleteSubstance(substance_name)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusFound, "Substance was successfully deleted")
}
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

	c.String(http.StatusFound, "SynthesisSubstance was successfully edited")
}
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
func (a *Application) edit_synthesis(c *gin.Context) {
	var synthesis ds.Syntheses

	var id = c.Param("synthesis")

	if err := c.BindJSON(&synthesis); err != nil {
		c.Error(err)
		return
	}

	err := a.repo.EditSynthesis(synthesis, id)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "Synthesis was successfully edited")

}
func (a *Application) generate_synthesis(c *gin.Context) {
	var synthesis ds.Syntheses

	var id = c.Param("synthesis")
	//	var title = c.Param("substance")

	if err := c.BindJSON(&synthesis); err != nil {
		c.Error(err)
		return
	}

	synthesis.Date_created = datatypes.Date(time.Now())
	err := a.repo.GenerateSynthesis(synthesis, id)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "Synthesis was successfully edited")

}
func (a *Application) apply_synthesis(c *gin.Context) {
	var synthesis ds.Syntheses

	var id = c.Param("synthesis")

	if err := c.BindJSON(&synthesis); err != nil {
		c.Error(err)
		return
	}

	switch synthesis.Status {
	case "В работе":
		synthesis.Date_processed = datatypes.Date(time.Now())
	case "Отклонена":
		synthesis.Date_processed = datatypes.Date(time.Now())
		synthesis.Date_finished = datatypes.Date(time.Now())
	case "Завершена":
		synthesis.Date_finished = datatypes.Date(time.Now())
	}

	//Date_processed = datatypes.Date(time.Now())
	//synthesis.Date_processed = datatypes.Date(time.Now())
	err := a.repo.ApplySynthesis(synthesis, id)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "Synthesis was successfully edited")

}
func (a *Application) order_synthesis(c *gin.Context) {
	var request_body ds.OrderSynthesisRequestBody
	request_body.User_id = instance1.GetUserIDAsString()
	if err := c.BindJSON(&request_body); err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, "Cant' parse json")
		return
	}

	err := a.repo.OrderSynthesis(request_body)

	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "Can't order synthesis")
		return
	}

	c.String(http.StatusCreated, "Synthesis was successfully ordered")

}
func (a *Application) order_find_substances(c *gin.Context) {
	var synthesis_id = c.Param("synthesis_id")

	found_substances, err := a.repo.FindSubstanceOrder(synthesis_id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, found_substances)

}
func (a *Application) add_user(c *gin.Context) {
	var user ds.Users

	if err := c.BindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, "Can't parse user\n"+err.Error())
		return
	}

	err := a.repo.CreateUser(user)

	if err != nil {
		c.String(http.StatusNotFound, "Can't create user\n"+err.Error())
		return
	}

	c.String(http.StatusCreated, "User created successfully")
}
func (a *Application) order_add(c *gin.Context) {
	var order ds.Synthesis_substance

	if err := c.BindJSON(&order); err != nil {
		c.String(http.StatusBadRequest, "Can't parse order\n"+err.Error())
		return
	}

	err := a.repo.OrderAdd(order)

	if err != nil {
		c.String(http.StatusNotFound, "Can't add substance\n"+err.Error())
		return
	}

	c.String(http.StatusCreated, "Add substance")
}
