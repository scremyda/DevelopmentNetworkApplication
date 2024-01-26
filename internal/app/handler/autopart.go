package handler

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/utils"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AutopartList godoc
// @Summary      Companies List
// @Description  Companies List
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        name query   string  false  "Query string to filter companies by name"
// @Success      200          {object}  ds.CompanyList
// @Failure      500          {object}  error
// @Router       /api/companies [get]
func (h *Handler) AutopartsList(ctx *gin.Context) {
	queryText, _ := ctx.GetQuery("autopart_name")
	autoparts, err := h.Repository.AutopartsList(queryText)
	if err != nil {
		h.errorHandler(ctx, http.StatusNoContent, err)
		return
	}
	userID, existsUser := ctx.Get("user_id")
	var draftIdRes uint = 0
	if existsUser {
		basketId, errBask := h.Repository.GetAssemblyDraftID(userID.(uint))
		if errBask != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errBask)
			return
		}
		draftIdRes = basketId
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"autoparts": autoparts,
		"draft_id":  draftIdRes,
	})
}

// GetAutopartById godoc
// @Summary      Company By ID
// @Description  Company By ID
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id   path    int     true        "Companies ID"
// @Success      200          {object}  ds.Company
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/companies/{id} [get]
func (h *Handler) GetAutopartById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id")[:], 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
	}

	autopart, err := h.Repository.GetAutopartById(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "autopart", autopart)
}

// DeleteAutopart godoc
// @Summary      Delete company by ID
// @Description  Deletes a company with the given ID
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Company ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /api/companies [delete]
func (h *Handler) DeleteAutopart(ctx *gin.Context) {
	var request struct {
		ID string `json:"id"`
	}
	if err := ctx.BindJSON(&request); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id, err2 := strconv.Atoi(request.ID)
	if err2 != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err2)
		return
	}
	if id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}

	url := h.Repository.DeleteAutopartImage(uint(id))

	if len(url) != 0 {
		err := h.DeleteImage(utils.ExtractObjectNameFromUrl(url))
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	}

	err := h.Repository.DeleteAutopart(uint(id))

	if gorm.IsRecordNotFoundError(err) {
		h.errorHandler(ctx, http.StatusBadRequest, err)
	} else if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
	}

	h.successHandler(ctx, "deleted_id", id)
}

func (h *Handler) AssemblyCurrent(ctx *gin.Context) {
	userID, existsUser := ctx.Get("user_id")
	if !existsUser {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("not fount `user_id` or `user_role`"))
		return
	}

	assemblies, errDB := h.Repository.AssemblyDraftId(userID.(uint))
	if errDB != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, errDB)
		return
	}

	h.successHandler(ctx, "tenders", assemblies)
}

// AddAutopart godoc
// @Summary      Add new company
// @Description  Add a new company with image, name, IIN
// @Tags         Companies
// @Accept       multipart/form-data
// @Produce      json
// @Param        image formData file true "Company image"
// @Param        name formData string true "Company name"
// @Param        description formData string false "Company description"
// @Param        IIN formData integer true "Company IIN"
// @Success      201  {string}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /api/companies [post]
func (h *Handler) AddAutopart(ctx *gin.Context) {
	var newAutopart ds.Autopart

	newAutopart.AutopartName = ctx.Request.FormValue("autopart_name")
	if newAutopart.AutopartName == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("имя автозапчасти не может быть пустой"))
		return
	}

	price := ctx.Request.FormValue("price")
	if len(price) != 0 {
		floatValue, err := strconv.ParseFloat(price, 64)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		newAutopart.Price = floatValue
	} else {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("цена автозапчасти не может быть пустой"))
		return
	}
	//if price == "" {
	//	h.errorHandler(ctx, http.StatusBadRequest, errors.New("цена не может быть пустой"))
	//	return
	//}

	newAutopart.Description = ctx.Request.FormValue("description")
	if newAutopart.Description == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("описание автозапчасти не может быть пустой"))
		return
	}
	//if newAutopart.Description == "" {
	//	h.errorHandler(ctx, http.StatusBadRequest, errors.New("описание не может быть пустой"))
	//	return
	//}

	//year := ctx.Request.FormValue("year")
	//intValue, err := strconv.Atoi(year)
	//if err != nil {
	//	h.errorHandler(ctx, http.StatusInternalServerError, err)
	//}
	//newAutopart.Year = intValue

	newAutopart.Status = ctx.Request.FormValue("status")

	file, header, err := ctx.Request.FormFile("image_url")
	if err != http.ErrMissingFile && err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при загрузке изображения"})
		return
	}
	if err == nil {
		if newAutopart.ImageURL, err = h.SaveImage(ctx.Request.Context(), file, header); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при сохранении изображения"})
			return
		}
	}

	create_id, err := h.Repository.AddAutopart(&newAutopart)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successAddHandler(ctx, "autopart_id", create_id)
}

// UpdateAutopart godoc
// @Summary      Update company by ID
// @Description  Updates a company with the given ID
// @Tags         Companies
// @Accept       multipart/form-data
// @Produce      json
// @Param        id          path        int     true        "ID"
// @Param        name        formData    string  false       "name"
// @Param        description formData    string  false       "description"
// @Param        IIN         formData    string  false       "IIN"
// @Param        image       formData    file    false       "image"
// @Success      200         {object}    map[string]any
// @Failure      400         {object}    error
// @Router       /api/companies/ [put]
func (h *Handler) UpdateAutopart(ctx *gin.Context) {
	var updatedAutopart ds.Autopart
	if err := ctx.BindJSON(&updatedAutopart); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if updatedAutopart.ImageURL != "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New(`image_url must be empty`))
		return
	}
	if updatedAutopart.Status != "действует" && updatedAutopart.Status != "удален" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New(`status_id может быть только действует или удален`))
		return
	}

	if updatedAutopart.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}
	if err := h.Repository.UpdateAutopart(&updatedAutopart); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successHandler(ctx, "updated_autopart", gin.H{
		"id":            updatedAutopart.ID,
		"autopart_name": updatedAutopart.AutopartName,
		"description":   updatedAutopart.Description,
		"image_url":     updatedAutopart.ImageURL,
		"status":        updatedAutopart.Status,
		"price":         updatedAutopart.Price,
		"year":          updatedAutopart.Year,
	})
}

func (h *Handler) AddImage(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	autopartID := ctx.Request.FormValue("autopart_id")

	if autopartID == "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("param `id` not found"))
		return
	}
	if header == nil || header.Size == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("no file uploaded"))
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	defer func(file multipart.File) {
		errLol := file.Close()
		if errLol != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errLol)
			return
		}
	}(file)

	ID, _ := strconv.Atoi(autopartID)
	url := h.Repository.DeleteAutopartImage(uint(ID))

	if len(url) != 0 {
		err := h.DeleteImage(utils.ExtractObjectNameFromUrl(url))
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	}

	var imageURL string
	if imageURL, err = h.SaveImage(ctx.Request.Context(), file, header); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	var updatedAutopart ds.Autopart
	updatedAutopart.ID = uint(ID)
	updatedAutopart.ImageURL = imageURL
	if err := h.Repository.UpdateAutopart(&updatedAutopart); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successAddHandler(ctx, "image_url", imageURL)
}

// AddAutopartToRequest godoc
// @Summary      Add company to request
// @Description  Adds a company to a tender request
// @Tags         Companies
// @Accept       json
// @Produce      json
// @Param        threatId  path  int  true  "Threat ID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Router       /companies/request [post]
func (h *Handler) AddAutopartToRequest(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("user_id not found"))
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		h.errorHandler(ctx, http.StatusUnauthorized, errors.New("`user_id` must be uint number"))
		return
	}

	var request ds.AddToAutopartID
	if err := ctx.BindJSON(&request); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if request.AutopartID == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "услуга не может быть пустой"})
		return
	}

	draftID, err := h.Repository.AddAutopartToDraft(request.AutopartID, userIDUint, request.Count)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	h.successHandler(ctx, "id", draftID)
}
