package handlers

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) DeleteFromAssembly(ctx *gin.Context) {
	var deleteFromAssembly ds.Autopart_Assembly
	if err := ctx.BindJSON(&deleteFromAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if deleteFromAssembly.AssemblyID <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	if deleteFromAssembly.AutopartID <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	err := h.Repository.DeleteFromAssembly(deleteFromAssembly)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "deleted successfully")
}

func (h *Handler) UpdateCountAutopartAssembly(ctx *gin.Context) {
	var updatedAutopartAssembly ds.Autopart_Assembly
	if err := ctx.BindJSON(&updatedAutopartAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if updatedAutopartAssembly.AssemblyID <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	if updatedAutopartAssembly.AutopartID <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	if updatedAutopartAssembly.Count <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, countInvalid)
		return
	}

	err := h.Repository.UpdateCountAutopartAssembly(updatedAutopartAssembly)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "updated_autopart_assembly", gin.H{
		"factory_id":  updatedAutopartAssembly.AssemblyID,
		"autopart_id": updatedAutopartAssembly.AutopartID,
		"count":       updatedAutopartAssembly.Count,
	})
}
