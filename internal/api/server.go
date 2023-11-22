package app

import (
	"awesomeProject1/internal/app/ds"
	"awesomeProject1/internal/app/dsn"
	"awesomeProject1/internal/app/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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

func (a *Application) StartServer() {
	a.r = gin.Default()
	log.Println("Server start up1")

	a.r.GET("substances", a.get_substances)
	a.r.GET("substance/:substance", a.get_substance)
	a.r.GET("syntheses", a.get_syntheses)
	a.r.GET("synthesis/:synthesis", a.get_synthesis)

	a.r.PUT("user/add", a.add_user)
	a.r.PUT("substance/add", a.add_substance)
	a.r.PUT("substance/delete/:substance_id", a.delete_substance)
	a.r.PUT("synthesis/delete/:synthesis_id", a.delete_synthesis)
	a.r.PUT("substance/edit", a.edit_substance)
	a.r.PUT("synthesis/edit", a.edit_synthesis)

	a.r.GET("synthesis/draft", a.synthesis_draft)

	a.r.GET("order/:synthesis_id", a.order_find_substances)

	a.r.PUT("order", a.order_synthesis)
	a.r.PUT("order/add", a.order_add)

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
	//var name_pattern = c.Query("name_pattern")
	var title = c.Query("title")
	var status = c.Query("status")

	substances, err := a.repo.GetAllSubstances(title, status)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, substances)
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

	found_synthesis, err := a.repo.GetAllSynthesis()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, found_synthesis)

}
func (a *Application) get_synthesis(c *gin.Context) {
	var title = c.Param("synthesis")

	found_synthesis, err := a.repo.FindSynthesis(title)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusFound, found_synthesis)

}
func (a *Application) add_substance(c *gin.Context) {
	var substance ds.Substances

	if err := c.BindJSON(&substance); err != nil {
		c.String(http.StatusBadRequest, "Can't parse substance\n"+err.Error())
		return
	}

	err := a.repo.CreateSubstance(substance)

	if err != nil {
		c.String(http.StatusNotFound, "Can't create substance\n"+err.Error())
		return
	}

	c.String(http.StatusCreated, "Substance created successfully")
}
func (a *Application) delete_synthesis(c *gin.Context) {
	synthesis_id, _ := strconv.Atoi(c.Param("synthesis_id"))

	err := a.repo.LogicalDeleteSynthesis(synthesis_id)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusFound, "Synthesis was successfully deleted")
}
func (a *Application) delete_substance(c *gin.Context) {
	substance_id, _ := strconv.Atoi(c.Param("substance_id"))

	err := a.repo.LogicalDeleteSubstance(substance_id)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusFound, "Substance was successfully deleted")
}
func (a *Application) edit_substance(c *gin.Context) {
	var substance ds.Substances

	if err := c.BindJSON(&substance); err != nil {
		c.Error(err)
		return
	}

	err := a.repo.EditSubstance(substance)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "Substance was successfully edited")

}
func (a *Application) edit_synthesis(c *gin.Context) {
	var synthesis ds.Syntheses

	if err := c.BindJSON(&synthesis); err != nil {
		c.Error(err)
		return
	}

	err := a.repo.EditSynthesis(synthesis)

	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusCreated, "Synthesis was successfully edited")

}
func (a *Application) order_synthesis(c *gin.Context) {
	var request_body ds.OrderSynthesisRequestBody

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
