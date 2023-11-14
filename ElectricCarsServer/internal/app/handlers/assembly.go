package handlers

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) AssembliesList(ctx *gin.Context) {
	assemblies, err := h.Repository.AssembliesList()

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	h.successHandler(ctx, "assemblies", assemblies)
}

func (h *Handler) AssemblyById(ctx *gin.Context) {
	assemblyStringID := ctx.Query("assembly")
	if assemblyStringID == "" {
		err := errors.New("error no query")
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	assemblyID, err := strconv.Atoi(assemblyStringID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	autoparts, assembly, errDB := h.Repository.AssemblyByID(uint(assemblyID))
	if errDB != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, errDB)
		return
	}
	assemblyDetails := ds.AssemblyDetails{
		Assembly:  assembly,
		Autoparts: autoparts,
	}

	h.successHandler(ctx, "assembly", assemblyDetails)
}

//func (h *Handler) AddAssembly(ctx *gin.Context) {
//	var assembly ds.Assembly
//	if err := ctx.BindJSON(&assembly); err != nil {
//		h.errorHandler(ctx, http.StatusBadRequest, err)
//		return
//	}
//	if assembly.ID != 0 {
//		h.errorHandler(ctx, http.StatusBadRequest, idMustBeEmpty)
//		return
//	}
//	if err := h.Repository.AddAssembly(&assembly); err != nil {
//		h.errorHandler(ctx, http.StatusInternalServerError, err)
//		return
//	}
//
//	h.successAddHandler(ctx, "assembly_id", assembly.ID)
//}

func (h *Handler) DeleteAssembly(ctx *gin.Context) {
	var request struct {
		ID uint `json:"id"`
	}
	if err := ctx.BindJSON(&request); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if request.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}
	if err := h.Repository.DeleteAssembly(request.ID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "assembly_id", request.ID)
}

func (h *Handler) UpdateAssembly(ctx *gin.Context) {
	var updatedAssembly ds.Assembly
	if err := ctx.BindJSON(&updatedAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if updatedAssembly.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}
	if err := h.Repository.UpdateAssembly(&updatedAssembly); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "updated_assembly", gin.H{
		"id":            updatedAssembly.ID,
		"assembly_name": updatedAssembly.Name,
		"date_created":  updatedAssembly.DateStart,
		"date_end":      updatedAssembly.DateEnd,
		// "image_url":     updatedAssembly.ImageURL,
		"status":      updatedAssembly.Status,
		"description": updatedAssembly.Description,
	})
}

//func (h *Handler) DeleteAssembly(ctx *gin.Context) {
//	var request struct {
//		ID uint `json:"id"`
//	}
//	if err := ctx.BindJSON(&request); err != nil {
//		h.errorHandler(ctx, http.StatusBadRequest, err)
//		return
//	}
//	if request.ID == 0 {
//		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
//		return
//	}
//	if err := h.Repository.DeleteAssembly(request.ID); err != nil {
//		h.errorHandler(ctx, http.StatusInternalServerError, err)
//		return
//	}
//
//	h.successHandler(ctx, "assembly_id", request.ID)
//}
