package handlers

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"ElectricCarsServer/ElectricCarsServer/internal/app/repo"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) AssembliesList(ctx *gin.Context) {
	queryStatus, _ := ctx.GetQuery("status")

	queryStart, _ := ctx.GetQuery("start")

	queryEnd, _ := ctx.GetQuery("end")

	assemblies, err := h.Repository.AssembliesList(queryStatus, queryStart, queryEnd)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	h.successHandler(ctx, "assemblies", assemblies)
}

func (h *Handler) AssemblyById(ctx *gin.Context) {
	assemblyStringID := ctx.Param("id")
	if assemblyStringID == "" {
		err := errors.New("error no get param")
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

func (h *Handler) FormAssembly(ctx *gin.Context) {
	var formAssembly ds.AssemblyForm
	if err := ctx.BindJSON(&formAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if formAssembly.User_id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	if formAssembly.Factory_id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	err := h.Repository.SameUser(formAssembly.User_id, formAssembly.Factory_id)
	if err != nil {
		if errors.Is(err, repo.NotSameUser) {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	updatedAssembly, err := h.Repository.FormAssembly(formAssembly)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "formed_assembly", updatedAssembly)
}

func (h *Handler) CompleteRejectAssembly(ctx *gin.Context) {
	var formAssembly ds.AssemblyForm
	if err := ctx.BindJSON(&formAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if formAssembly.User_id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	if formAssembly.Factory_id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	isModerator, err := h.Repository.UserAdmin(formAssembly.User_id)
	if err != nil {
		if errors.Is(err, repo.UserNotFound) {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if isModerator == false {
		h.errorHandler(ctx, http.StatusBadRequest, userIsNotModerator)
		return
	}

	updatedAssembly, err := h.Repository.CompleteRejectAssembly(formAssembly)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "completed/rejected_assembly", updatedAssembly)
}

func (h *Handler) RejectAssembly(ctx *gin.Context) {
	var formAssembly ds.AssemblyForm
	if err := ctx.BindJSON(&formAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if formAssembly.User_id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	if formAssembly.Factory_id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	isModerator, err := h.Repository.UserAdmin(formAssembly.User_id)
	if err != nil {
		if errors.Is(err, repo.UserNotFound) {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if isModerator == false {
		h.errorHandler(ctx, http.StatusBadRequest, userIsNotModerator)
		return
	}

	updatedAssembly, err := h.Repository.RejectAssembly(formAssembly)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "rejected_assembly", updatedAssembly)
}

func (h *Handler) DeleteAssembly(ctx *gin.Context) {
	var formAssembly ds.AssemblyForm
	if err := ctx.BindJSON(&formAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if formAssembly.User_id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	if formAssembly.Factory_id <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}

	isModerator, err := h.Repository.UserAdmin(formAssembly.User_id)
	if err != nil {
		if errors.Is(err, repo.UserNotFound) {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	errSameUser := h.Repository.SameUser(formAssembly.User_id, formAssembly.Factory_id)
	if errSameUser == nil || isModerator == true {
		updatedAssembly, err := h.Repository.DeleteAssembly(formAssembly)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		h.successHandler(ctx, "deleted_assembly", updatedAssembly)
	}
	if errSameUser != nil {
		if errors.Is(errSameUser, repo.NotSameUser) {
			h.errorHandler(ctx, http.StatusBadRequest, errSameUser)
			return
		}
		h.errorHandler(ctx, http.StatusInternalServerError, errSameUser)
		return
	}

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
