package handler

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/role"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func ParseDateString(dateString string) (time.Time, error) {
	format := "2006-01-02 15:04:05"
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})`)
	matches := re.FindStringSubmatch(dateString)
	if len(matches) < 2 {
		return time.Time{}, nil
	}
	parsedTime, err := time.Parse(format, matches[1])
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

// AssemblyList godoc
// @Summary      Get list of tender requests
// @Description  Retrieves a list of tender requests based on the provided parameters
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        status      query  string    false  "Tender request status"
// @Param        start  query  string    false  "Start date in the format '2006-01-02T15:04:05Z'"
// @Param        end    query  string    false  "End date in the format '2006-01-02T15:04:05Z'"
// @Success      200  {object}  []ds.Tender
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /api/tenders [get]
func (h *Handler) AssemblyList(ctx *gin.Context) {
	userID, existsUser := ctx.Get("user_id")
	userRole, existsRole := ctx.Get("user_role")
	if !existsUser || !existsRole {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	switch userRole {
	case role.Buyer:
		h.assemblyByUserId(ctx, fmt.Sprintf("%d", userID))
		return
	default:
		break
	}

	queryStatus := ctx.Query("status_id")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" {
		startDateStr = "0001-01-01 00:00:00"
	}
	if endDateStr == "" {
		endDateStr = time.Now().Add(time.Hour * 24).String()
	}

	startDate, errStart := ParseDateString(startDateStr + " 00:00:00")
	endDate, errEnd := ParseDateString(endDateStr + " 00:00:00")
	h.Logger.Info(startDate, endDate)
	if errEnd != nil || errStart != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("incorrect `start_date` or `end_date`"))
		return
	}

	assemblies, err := h.Repository.AssembliesList(queryStatus, startDate, endDate)

	if err != nil {
		h.errorHandler(ctx, http.StatusNoContent, err)
		return
	}
	h.successHandler(ctx, "assemblies", assemblies)
}

func (h *Handler) assemblyByUserId(ctx *gin.Context, userID string) {
	assemblies, errDB := h.Repository.AssemblyByUserID(userID)
	if errDB != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, errDB)
		return
	}

	h.successHandler(ctx, "assemblies", assemblies)
}

// GetAssemblyById godoc
// @Summary      Get tender request by ID
// @Description  Retrieves a tender request with the given ID
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Tender Request ID"
// @Success      200  {object}  ds.TenderDetails
// @Failure      400  {object}  error
// @Router       /api/tenders/{id} [get]
func (h *Handler) GetAssemblyById(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	assembly, err := h.Repository.AssemblyByID(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	h.successHandler(c, "assembly", assembly)
}

// UpdateAssembly godoc
// @Summary      Update Tender by admin
// @Description  Update Tender by admin
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        input    body    ds.Tender  true    "updated Assembly"
// @Success      200          {object}  nil
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/tenders [put]
func (h *Handler) UpdateAssembly(ctx *gin.Context) {
	userID, existsUser := ctx.Get("user_id")
	userRole, existsRole := ctx.Get("user_role")
	if !existsUser || !existsRole {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	var updatedAssembly ds.UpdateAssembly
	if err := ctx.BindJSON(&updatedAssembly); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if updatedAssembly.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("id некоректен"))
		return
	}

	var updated ds.Assembly
	updated.ID = updatedAssembly.ID
	updated.Name = updatedAssembly.Name

	assembly, err := h.Repository.AssemblyModel(updated.ID)
	assembly.Name = updatedAssembly.Name

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("assembly with `id` = %d not found", assembly.ID))
		return
	}

	if assembly.UserID != userID && userRole == role.Buyer {
		h.errorHandler(ctx, http.StatusForbidden, errors.New("you cannot change the assembly if it's not yours"))
		return
	}

	if err := h.Repository.UpdateAssembly(assembly); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "updated_assembly", gin.H{
		"id":              updatedAssembly.ID,
		"assembly_name":   updatedAssembly.Name,
		"creation_date":   assembly.CreationDate,
		"completion_date": assembly.CompletionDate,
		"formation_date":  assembly.FormationDate,
		"user_id":         assembly.UserID,
		"status":          assembly.Status,
	})
}

// FormAssemblyRequest godoc
// @Summary      Form Company by client
// @Description  Form Company by client
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Tender form ID"
// @Success      200          {object}  ds.TenderDetails
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/tenders/form [put]
func (h *Handler) FormAssemblyRequest(c *gin.Context) {
	userID, existsUser := c.Get("user_id")
	if !existsUser {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	err, _ := h.Repository.FormAssemblyRequestByID(userID.(uint))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

// UpdateStatusAssemblyRequest godoc
// @Summary      Update transaction request status by ID
// @Description  Updates the status of a transaction request with the given ID on "завершен"/"отклонен"
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Request ID"
// @Param        input    body    ds.NewStatus  true    "update status"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /tenders/updateStatus [put]
func (h *Handler) UpdateStatusAssemblyRequest(c *gin.Context) {
	var status ds.NewStatus
	if err := c.BindJSON(&status); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	userIDStr, existsUser := c.Get("user_id")
	if !existsUser {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}
	userID := userIDStr.(uint)

	if status.Status != "отклонен" && status.Status != "завершен" {
		h.errorHandler(c, http.StatusBadRequest, errors.New("статус можно поменять только на 'отклонен' и 'завершен'"))
	}

	if err := h.Repository.FinishRejectHelper(status.Status, status.AssemblyID, userID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

// DeleteAutopartFromRequest godoc
// @Summary      Delete company from request
// @Description  Deletes a company from a request based on the user ID and company ID
// @Tags         Tender_Company
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "company ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  error
// @Router       /api/tender-request-company [delete]
func (h *Handler) DeleteAutopartFromRequest(c *gin.Context) {
	var body struct {
		ID int `json:"id"`
	}

	if err := c.BindJSON(&body); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if body.ID == 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	err := h.Repository.DeleteAutopartFromRequest(body.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	h.successHandler(c, "deleted_autopart_assembly", body.ID)
}

// DeleteAssembly godoc
// @Summary      Delete tender request by user ID
// @Description  Deletes a tender request for the given user ID
// @Tags         Tenders
// @Accept       json
// @Produce      json
// @Param        user_id  path  int  true  "User ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/tenders [delete]
func (h *Handler) DeleteAssembly(c *gin.Context) {
	userID, existsUser := c.Get("user_id")
	userRole, existsRole := c.Get("user_role")
	if !existsUser || !existsRole {
		h.errorHandler(c, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	var request struct {
		ID uint `json:"id"`
	}

	if err := c.BindJSON(&request); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if request.ID == 0 {
		h.errorHandler(c, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	assembly, err := h.Repository.AssemblyModel(request.ID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, fmt.Errorf("assembly with `id` = %d not found", assembly.ID))
		return
	}

	if assembly.UserID != userID && userRole == role.Buyer {
		h.errorHandler(c, http.StatusForbidden, errors.New("you are not the creator. you can't delete a assembly"))
		return
	}

	err = h.Repository.DeleteAssemblyByID(request.ID)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	h.successHandler(c, "assembly_id", request.ID)
}

// UpdateTenderCompany godoc
// @Summary      Update money Tender Company
// @Description  Update money Tender Company by client
// @Tags         Tender_Company
// @Accept       json
// @Produce      json
// @Param        input    	  body    ds.TenderCompany true    "Update money Tender Company"
// @Success      200          {object} map[string]string "update"
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/tender-request-company [put]
func (h *Handler) UpdateAssemblyAutopart(c *gin.Context) {
	//var TenderCompany ds.TenderCompany
	var AssemblyAutopartU ds.AssemblyAutopartUpdate
	if err := c.BindJSON(&AssemblyAutopartU); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	err := h.Repository.UpdateAssemblyAutopart(AssemblyAutopartU.ID, AssemblyAutopartU.Count)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, "update")
}
