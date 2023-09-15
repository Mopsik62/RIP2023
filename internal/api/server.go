package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type Card struct {
	ID    int
	Title string
	Text  string
	Image string
}

var cards = []Card{
	{1, "Кальций", "", "image/Ca.jpg"},
	{2, "Золото", "", "image/Au.jpg"},
	{3, "Свинья", "", "image/Pb.jpg"},
	{4, "Медь", "", "image/Cu.jpg"},
}

//type CardPage struct {
//	Title string
//	Text  string
//	Image string
//}

func StartServer() {
	log.Println("Server start up")
	r := gin.Default()
	//cards := []Card{
	//	{1, "Кальций", "", "image/Ca.jpg"},
	//	{2, "Золото", "", "image/Au.jpg"},
	//	{3, "Свинец", "", "image/Pb.jpg"},
	//	{4, "Медь", "", "image/Cu.jpg"},
	//}
	//cardPages := []CardPage{
	//	{"Кальций", "Здесь вы сможете создать заявку на синтез Кальция", "../image/Ca.jpg"},
	//	{"Золото", "Здесь вы сможете создать заявку на синтез Золота", "../image/Au.jpg"},
	//	{"Свинец", "Здесь вы сможете создать заявку на синтез Свинца", "../image/Pb.jpg"},
	//	{"Натрий", "Здесь вы сможете создать заявку на синтез Меди", "../image/Cu.jpg"},
	//}

	//idToPage := make(map[int]CardPage)
	//idToPage[1] = cardPages[0]
	//idToPage[2] = cardPages[1]
	//idToPage[3] = cardPages[2]
	//idToPage[4] = cardPages[3]
	//
	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "pong",
	//	})
	//})
	r.LoadHTMLGlob("templates/*html")
	r.Static("/image", "./resources/image")
	r.Static("/css", "./templates/css")

	r.GET("/", loadHome)
	r.GET("/:title", loadPage)
	//r.GET("/home", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "index.html", gin.H{
	//		"cards": cards,
	//	})
	//})
	//
	//r.GET("/home/:id", func(c *gin.Context) {
	//	id, err := strconv.Atoi(c.Param("id"))
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	cardPage := idToPage[id]
	//
	//	c.HTML(http.StatusOK, "page.html", gin.H{
	//		"Title": cardPage.Title,
	//		"Image": cardPage.Image,
	//		"Text":  cardPage.Text,
	//	})
	//})
	//r.GET("/search", func(c *gin.Context) {
	//	card_title := c.Query("card_title")
	//
	//	if card_title == "" {
	//		c.Redirect(http.StatusFound, "/home")
	//	}
	//
	//	foundCards := []Card{}
	//	lowerCardTitle := strings.ToLower(card_title)
	//	for i := range cards {
	//		if strings.Contains(strings.ToLower(cards[i].Title), lowerCardTitle) {
	//			foundCards = append(foundCards, cards[i])
	//		}
	//	}
	//
	//	c.HTML(http.StatusOK, "search.html", gin.H{
	//		"cards":       foundCards,
	//		"Amount":      len(foundCards),
	//		"Search_text": card_title,
	//	})
	//
	//})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
func loadHome(c *gin.Context) {
	card_title := c.Query("card_title")

	if card_title == "" {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"cards": cards,
		})
		return
	}

	foundCards := []Card{}
	lowerCardTitle := strings.ToLower(card_title)
	for i := range cards {
		if strings.Contains(strings.ToLower(cards[i].Title), lowerCardTitle) {
			foundCards = append(foundCards, cards[i])
		}
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"cards": foundCards,
	})

}
func loadPage(c *gin.Context) {
	title := c.Param("title")

	for i := range cards {
		if cards[i].Title == title {
			c.HTML(http.StatusOK, "page.html", gin.H{
				"Title": cards[i].Title,
				"Image": "../" + cards[i].Image,
				"Text":  cards[i].Text,
			})
			return
		}
	}
}
