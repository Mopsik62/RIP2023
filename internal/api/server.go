package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type Card struct {
	ID     int
	Title  string
	Text   string
	Number string
	Image  string
}

var cards = []Card{
	{1, "Кальций", "Ca", "20", "image/Ca.jpg"},
	{2, "Золото", "Au", "79", "image/Au.jpg"},
	{3, "Свинец", "Pb", "82", "image/Pb.jpg"},
	{4, "Медь", "Cu", "29", "image/Cu.jpg"},
}

func StartServer() {
	log.Println("Server start up")
	r := gin.Default()

	r.LoadHTMLGlob("templates/*html")
	r.Static("/image", "./resources/image")
	r.Static("/css", "./templates/css")

	r.GET("/", loadSubstances)
	r.GET("/:title", loadSubstance)

	r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
func loadSubstances(c *gin.Context) {
	substance_name := c.Query("substance_name")

	if substance_name == "" {
		c.HTML(http.StatusOK, "substances.html", gin.H{
			"cards": cards,
		})
		return
	}

	foundCards := []Card{}
	lowerCardTitle := strings.ToLower(substance_name)
	for i := range cards {
		if strings.Contains(strings.ToLower(cards[i].Title), lowerCardTitle) {
			foundCards = append(foundCards, cards[i])
		}
	}

	c.HTML(http.StatusOK, "substances.html", gin.H{
		"cards":          foundCards,
		"substance_name": substance_name,
	})

}
func loadSubstance(c *gin.Context) {
	title := c.Param("title")

	for i := range cards {
		if cards[i].Title == title {
			c.HTML(http.StatusOK, "substance.html", gin.H{
				"Title":  cards[i].Title,
				"Image":  "../" + cards[i].Image,
				"Text":   cards[i].Text,
				"Number": cards[i].Number,
			})
			return
		}
	}
}
