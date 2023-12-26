package handlers

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"ElectricCarsServer/ElectricCarsServer/internal/app/repo"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

const (
	ServerToken = "qwerzxfsdfoiw"
	ServiceUrl  = "http://127.0.0.1:8081/discussion/"
)

func (h *Handler) AssemblyDiscussionStart(c *gin.Context) {
	// принимает заявку и отправляет её в сервис
	var request ds.RequestAsyncService
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("неверный формат"))
		return
	}

	request.Token = ServerToken

	body, _ := json.Marshal(request)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", ServiceUrl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	if resp.StatusCode == 200 {
		c.JSON(http.StatusOK, gin.H{"message": "заявка принята в обработку"})
		return
	}
	c.AbortWithError(http.StatusInternalServerError, errors.New("заявка не принята в обработку"))
}

// ручка вызывается сервисом на python
func (h *Handler) AssemblyDiscussionFinish(c *gin.Context) {
	var request ds.RequestAsyncService
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	if request.Token != ServerToken {
		c.AbortWithError(http.StatusForbidden, errors.New("неверный токен"))
		return
	}

	// сохраняем в базу
	err := h.Repository.SaveAssemblyDiscussion(request)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "данные сохранены"})
}

// AssembliesList godoc
// @Summary      Assembly List
// @Description  Assembly List
// @Tags         Assembly
// @Accept       json
// @Produce      json
// @Param        status query  string  false  "Query string to filter Assemblies by status"
// @Param        start query  string  false  "Query string to filter Assemblies from start date"
// @Param        end query  string  false  "Query string to filter Assemblies to end date"
// @Success      200        {object}  []ds.Assembly
// @Failure      400          {object} error
// @Router       /api/assembly/get-al [get]
func (h *Handler) AssembliesList(ctx *gin.Context) {
	queryStatus, _ := ctx.GetQuery("status")

	queryStart, _ := ctx.GetQuery("start")

	queryEnd, _ := ctx.GetQuery("end")

	assemblies, err := h.Repository.AssembliesList(queryStatus, queryStart, queryEnd, ctx.GetInt(userCtx), ctx.GetBool(adminCtx))

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	h.successHandler(ctx, "assemblies", assemblies)
}

// AssemblyById godoc
// @Summary      Assembly By ID
// @Description  Assembly By ID
// @Tags         Assembly
// @Accept       json
// @Produce      json
// @Param        id   path    int     true        "Assembly ID"
// @Success      200       {object}   []ds.AssemblyDetails
// @Failure      400          {object}  error
// @Failure      500          {object}  error "server error"
// @Router       /api/assembly{id} [get]
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
	autoparts, assembly, errDB := h.Repository.AssemblyByID(uint(assemblyID), ctx.GetInt(userCtx), ctx.GetBool(adminCtx))
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

// UpdateAssembly godoc
// @Summary      Update Assembly by admin
// @Description  Update Assembly by admin
// @Tags         Assembly
// @Accept       json
// @Produce      json
// @Param        input    body    ds.Assembly  true    "updated Assembly"
// @Success      200          {object}  ds.Assembly
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/assembly [put]
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

// FormAssembly godoc
// @Summary      Form Assembly by client
// @Description  Form Assembly by client
// @Tags         Assembly
// @Accept       json
// @Produce      json
// @Param        input    body    ds.AssemblyForm  true    "Form Assembly"
// @Success      200          {object}  ds.Assembly
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/assembly/form [put]
func (h *Handler) FormAssembly(ctx *gin.Context) {
	var formAssembly ds.AssemblyForm
	if err := ctx.BindJSON(&formAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	formAssembly.User_id = uint(ctx.GetInt(userCtx))

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

// CompleteAssembly godoc
// @Summary      Complete Assembly by admin
// @Description  Complete Assembly by admin
// @Tags         Assembly
// @Accept       json
// @Produce      json
// @Param        input    body    ds.AssemblyForm  true    "Complete Assembly"
// @Success      200          {object}  ds.Assembly
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/assembly/complete [put]
func (h *Handler) CompleteAssembly(ctx *gin.Context) {
	var formAssembly ds.AssemblyForm
	if err := ctx.BindJSON(&formAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	formAssembly.User_id = uint(ctx.GetInt(userCtx))

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

	updatedAssembly, err := h.Repository.CompleteAssembly(formAssembly)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	result := ds.AssemblyAdmin{
		Assembly:   updatedAssembly,
		AdminLogin: "Admin",
	}

	h.successHandler(ctx, "completed_assembly", result)
}

// RejectAssembly godoc
// @Summary      Reject Assembly by admin
// @Description  Reject Assembly by admin
// @Tags         Assembly
// @Accept       json
// @Produce      json
// @Param        input    body    ds.AssemblyForm  true    "Reject Assembly"
// @Success      200          {object}  ds.Assembly
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/assembly/reject [put]
func (h *Handler) RejectAssembly(ctx *gin.Context) {
	var formAssembly ds.AssemblyForm
	if err := ctx.BindJSON(&formAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	formAssembly.User_id = uint(ctx.GetInt(userCtx))

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

// DeleteAssembly godoc
// @Summary      Delete Assembly by admin
// @Description  Delete Assembly by admin
// @Tags         Assembly
// @Accept       json
// @Produce      json
// @Param        input    body    ds.AssemblyForm  true    "Delete Assembly"
// @Success      200          {object}  ds.Assembly
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/assembly [delete]
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
