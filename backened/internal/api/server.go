package api

import (
	backend "backened"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func StartServer() {
	r := gin.Default()

	r.Static("/static", "./static")

	autoparts := []backend.Autoparts{
		{1, "Прицепное устройство (фаркоп)", "Прицепное устройство Bosal, под фаркоп, в сборе с розеткой, Tesla Model X.", "autopart1.jpg", []string{"1027581-00-D", "1046032-00-C", "10057987-001"}},
		{2, "Полуось задняя", "Полуось задняя, лев./прав., Tesla Model X.", "autopart2.jpg", []string{"1285161-00-C", "1027161-00-C", "1046032-00-C"}},
		{3, "Модуль управления двери передней левой", "Блок управления двери передней левой, Tesla Model X.", "autopart3.jpg", []string{"1037222-00-C", "1751222-00-D", "170347162-00-A"}},
		{4, "Вентилятор печки салона", "Вентилятор климатической установки (печки) салона, Tesla Model 3.", "autopart4.jpg", []string{"1107669-00-B", " 1099999-00-H", "30231C"}},
		{5, "Тепловой насос", "Тепловой насос в сборе (теплообменник, бачок расширительный, трубки фреона, клапана, электропроводка), Tesla Model 3, Y.", "autopart5.jpg", []string{"1547595-96-F", ", 1523000-96-E", "1523001-00-D"}},
	}

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		r.SetHTMLTemplate(template.Must(template.ParseFiles("./templates/index.html")))
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Autoparts": autoparts,
		})
	})

	r.GET("/autoparts/:id", func(c *gin.Context) {
		r.SetHTMLTemplate(template.Must(template.ParseFiles("./templates/info.html")))
		id := c.Param("id")
		var selectedautopart backend.Autoparts
		for _, autopart := range autoparts {
			if strconv.Itoa(autopart.ID) == id {
				selectedautopart = autopart
				break
			}
		}

		c.HTML(http.StatusOK, "info.html", gin.H{
			"Autoparts": selectedautopart,
		})
	})
	r.GET("/search", func(c *gin.Context) {
		searchQuery := c.Query("search")

		filteredautoparts := []backend.Autoparts{}
		for _, autopart := range autoparts {
			if strings.Contains(strings.ToLower(autopart.Name), strings.ToLower(searchQuery)) {
				filteredautoparts = append(filteredautoparts, autopart)
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Autoparts": filteredautoparts,
		})

	})

	r.Run(":8080")
}
