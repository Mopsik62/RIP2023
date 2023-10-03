package app

import (
	"awesomeProject1/internal/app/dsn"
	"awesomeProject1/internal/app/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Application struct {
	repo repository.Repository
	r    *gin.Engine
}

func New() Application {
	app := Application{}
	//log.Println("Application start3333333!")

	repo, _ := repository.New(dsn.FromEnv())
	log.Println("Application start444444!")
	log.Println(repo)
	app.repo = *repo

	return app

}

func (a *Application) StartServer() {
	//log.Println("Server start up")
	//log.Println("Server start up0")
	a.r = gin.Default()
	log.Println("Server start up1")
	a.r.GET("/", a.loadSubstances)
	a.r.GET("/:substance_name", a.loadSubstance)
	a.r.POST("/delete_substance/:substance_name", a.loadSubstanceChangeVisibility)

	log.Println("Server start up2")
	a.r.LoadHTMLGlob("templates/*.html")
	//r.Static("/image", "./resources/image")
	a.r.Static("/css", "./templates/css")

	a.r.Static("/image", "./resources")

	a.r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
func (a *Application) loadSubstances(c *gin.Context) {
	substance_name := c.Query("substance_name")

	if substance_name == "" {
		all_substances, err := a.repo.GetAllSubstances()

		if err != nil {
			log.Println(err)
			c.Error(err)
		}

		c.HTML(http.StatusOK, "substances.html", gin.H{
			"substances": all_substances,
		})
	} else {
		found_substances, err := a.repo.SearchSubstances(substance_name)

		if err != nil {
			c.Error(err)
			return
		}

		c.HTML(http.StatusOK, "substances.html", gin.H{
			"substances":  found_substances,
			"Search_text": substance_name,
		})
	}
}
func (a *Application) loadSubstance(c *gin.Context) {
	substance_name := c.Param("substance_name")

	if substance_name == "favicon.ico" {
		return
	}

	substance, err := a.repo.GetSubstanceByName(substance_name)

	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "substance.html", gin.H{
		"Title":  substance.Title,
		"Image":  substance.Image,
		"Text":   substance.Text,
		"Number": substance.Number,
	})
}

func (a *Application) loadSubstanceChangeVisibility(c *gin.Context) {
	substance_name := c.Param("substance_name")
	err := a.repo.ChangeSubstanceVisibility(substance_name)

	if err != nil {
		c.Error(err)
	}

	c.Redirect(http.StatusFound, "/")
}
