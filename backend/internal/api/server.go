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
		{1, "Зарядный кабель", "Прицепное устройство Bosal, под фаркоп, в сборе с розеткой, Tesla Model X.", "image11.jpeg", []string{"1027581-00-D", "1046032-00-C", "10057987-001"}},
		{2, "Ротор электромотора", "Ротор электромотора, лев./прав., Tesla Model X.", "image22.jpeg", []string{"1285161-00-C", "1027161-00-C", "1046032-00-C"}},
		{3, "Электродвигатель", "Электродвигатель в сборе (теплообменник, бачок расширительный, трубки фреона, клапана, электропроводка), Tesla Model 3, Y.", "image33.jpg", []string{"1547595-96-F", ", 1523000-96-E", "1523001-00-D"}},
	}

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/index.tmpl")))
		searchQuery := c.Query("search")
		filteredautoparts := []backend.Autoparts{}
		if searchQuery != "" {
			for _, autopart := range autoparts {
				if strings.Contains(strings.ToLower(autopart.Name), strings.ToLower(searchQuery)) {
					filteredautoparts = append(filteredautoparts, autopart)
				}
			}

			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"Autoparts":   filteredautoparts,
				"SearchValue": searchQuery, // добавлено здесь
			})

		} else {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"Autoparts": autoparts,
			})
		}
	})

	r.GET("/autoparts/:id", func(c *gin.Context) {
		r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/info.tmpl")))
		id := c.Param("id")
		var selectedautopart backend.Autoparts
		for _, autopart := range autoparts {
			if strconv.Itoa(autopart.ID) == id {
				selectedautopart = autopart
				break
			}
		}

		c.HTML(http.StatusOK, "info.tmpl", gin.H{
			"Autoparts": selectedautopart,
		})
	})

	r.Run(":8080")
}
