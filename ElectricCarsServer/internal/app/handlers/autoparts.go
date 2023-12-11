package handlers

import (
	"ElectricCarsServer/ElectricCarsServer/internal/app/ds"
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
)

// AutopartsList godoc
// @Summary      Autoparts List
// @Description  Autoparts List
// @Tags         Autoparts
// @Accept       json
// @Produce      json
// @Param       name query   string  false  "Query string to filter autoparts by name"
// @Success      200        {object}  ds.AutopartList
// @Failure      500          {object}  error
// @Router       /api/autoparts/get-al [get]
func (h *Handler) AutopartsList(ctx *gin.Context) {
	queryBrand, _ := ctx.GetQuery("name")

	autoparts, err := h.Repository.AutopartsList(queryBrand)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "autoparts", autoparts)
}

// AutopartById godoc
// @Summary      Autopart By ID
// @Description  Autopart By ID
// @Tags         Autoparts
// @Accept       json
// @Produce      json
// @Param        id   path    int     true        "Autoaprt ID"
// @Success      200        {object}  ds.Autopart
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/autoparts{id} [get]
func (h *Handler) AutopartById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		err := errors.New("error no get param")
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	autoparts, errBD := h.Repository.AutopartById(uint(id))
	if errBD != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, errBD)
		return
	}

	h.successHandler(ctx, "autopart", autoparts)
}

// DeleteAutopart godoc
// @Summary      Delete Autopart by admin
// @Description  Delete Autopart by admin
// @Tags         Autoparts
// @Accept       json
// @Produce      json
// @Param        id           path        int     true   "Autoaprt ID"
// @Success      200
// @Failure      400          {object}  error
// @Failure      500          {object}  error
// @Router       /api/autoparts{id} [delete]
func (h *Handler) DeleteAutopart(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		err := errors.New("error no get param")
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if id == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}
	if err := h.Repository.DeleteAutopart(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err) // TODO: catch not found
		return
	}

	h.successHandler(ctx, "deleted_id", id)
}

// AddImage godoc
// @Summary      Add Image to autopart by admin
// @Description  Add Image to autopart by admin
// @Tags         Autoparts
// @Accept       multipart/form-data
// @Produce      json
// @Param        id          formData    int     false       "autopart ID"
// @Param        image       formData    file    false       "image"
// @Success      200
// @Failure      400         {object}    error
// @Failure      500         {object}    error
// @Router       /api/autoparts/upload-image [post]
func (h *Handler) AddImage(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	autopartID := ctx.Request.FormValue("autopart_id")

	if autopartID == "" {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}
	if header == nil || header.Size == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, headerNotFound)
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

	// Upload the image to minio server.
	newImageURL, errMinio := h.createImageInMinio(&file, header)
	if errMinio != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, errMinio)
		return
	}
	if err = h.Repository.UpdateAutopartImage(autopartID, newImageURL); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successAddHandler(ctx, "image_url", newImageURL)
}

// AddAutopart godoc
// @Summary      Add Autopart by admin
// @Description  Add Autopart by admin
// @Tags         Autoparts
// @Accept       json
// @Produce      json
// @Param        input    body    ds.Autopart  true    "New autopart"
// @Success      200
// @Failure      400          {object}  error
// @Router       /api/autoparts [post]
func (h *Handler) AddAutopart(ctx *gin.Context) {
	var newAutopart ds.Autopart
	if err := ctx.BindJSON(&newAutopart); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	newAutopart.UserID = uint(ctx.GetInt(userCtx))

	if newAutopart.ID != 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idMustBeEmpty)
		return
	}
	if newAutopart.UserID <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idCantBeEmpty)
		return
	}
	if newAutopart.Name == "" {
		h.errorHandler(ctx, http.StatusBadRequest, autopartNameCannotBeEmpty)
		return
	}
	if newAutopart.Brand == "" {
		h.errorHandler(ctx, http.StatusBadRequest, autopartBrandCannotBeEmpty)
		return
	}
	if newAutopart.Models == "" {
		h.errorHandler(ctx, http.StatusBadRequest, autopartModelsCannotBeEmpty)
		return
	}
	if newAutopart.Year <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, autopartYearCannotBeEmpty)
		return
	}
	if newAutopart.Price <= 0 {
		h.errorHandler(ctx, http.StatusBadRequest, autopartPriceCannotBeEmpty)
		return
	}
	if err := h.Repository.AddAutopart(&newAutopart); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successAddHandler(ctx, "autopart_id", newAutopart.ID)
}

// UpdateAutopart godoc
// @Summary      Update Autopart by admin
// @Description  Update Autopart by admin
// @Tags         Autoparts
// @Accept       json
// @Produce      json
// @Param        input    body    ds.Autopart  true    "updated autopart"
// @Success      200          {object}  ds.Autopart
// @Failure      400          {object}  error
// @Router       /api/autoparts [put]
func (h *Handler) UpdateAutopart(ctx *gin.Context) {
	var updatedAutopart ds.Autopart
	if err := ctx.BindJSON(&updatedAutopart); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updatedAutopart.UserID = uint(ctx.GetInt(userCtx))

	if updatedAutopart.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}
	if err := h.Repository.UpdateAutopart(&updatedAutopart); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successHandler(ctx, "updated_autopart", gin.H{
		"id":            updatedAutopart.ID,
		"autopart_name": updatedAutopart.Name,
		"description":   updatedAutopart.Description,
		"brand":         updatedAutopart.Brand,
		"model":         updatedAutopart.Models,
		"year":          updatedAutopart.Year,
		"image_url":     updatedAutopart.Image,
		"status":        updatedAutopart.Status,
		"price":         updatedAutopart.Price,
	})
}

// AddToAssembly godoc
// @Summary      Add To Assembly by client
// @Description  Add To Assembly by client
// @Tags         Autoparts
// @Accept       json
// @Produce      json
// @Param        input    body    ds.AddToAssemblyID  true    "Add To Assembly autopart"
// @Success      200
// @Failure      400          {object}  error
// @Router       /api/autoparts/add-to-assembly [post]
func (h *Handler) AddToAssembly(ctx *gin.Context) {
	var AddToAssemblyID ds.AddToAssemblyID
	err := ctx.BindJSON(&AddToAssemblyID)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	AddToAssemblyID.User_id = uint(ctx.GetInt(userCtx))

	if AddToAssemblyID.AutopartDetails.Autopart_id <= 0 || AddToAssemblyID.User_id <= 0 {
		err := errors.New("некорректные данные")
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.AddToAssembly(&AddToAssemblyID); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successAddHandler(ctx, "autopart_id", AddToAssemblyID.AutopartDetails.Autopart_id)
}
