package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) AutopartsList(ctx *gin.Context) {
	searchQuery := ctx.Query("search")
	if searchQuery == "" {
		autoparts, err := h.Repository.Searchautopart(searchQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"Autoparts": autoparts,
		})
	} else {
		filteredautoparts, err := h.Repository.Searchautopart(searchQuery)
		if err != nil {

		}
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"Autoparts":   filteredautoparts,
			"SearchQuery": searchQuery,
		})

	}
}

func (h *Handler) AutopartById(ctx *gin.Context) {
	id := ctx.Param("id")
	autoparts, err := h.Repository.AutopartById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.HTML(http.StatusOK, "info.html", gin.H{
		"Autoparts": autoparts,
	})
}

func (h *Handler) Deleteautopart(ctx *gin.Context) {
	id := ctx.Param("id")
	h.Repository.Deleteautopart(id)
	ctx.Redirect(http.StatusFound, "/")
}
