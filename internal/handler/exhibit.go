package handler

import (
	"Museum/internal/auth"
	"Museum/internal/models"
	"Museum/internal/service"
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExhibitHandler struct {
	svc        *service.ExhibitService
	imagesPath string
}

func NewExhibitHandler(svc *service.ExhibitService, imagesPath string) *ExhibitHandler {
	return &ExhibitHandler{svc: svc, imagesPath: imagesPath}
}

type CreateExhibitRequest struct {
	Category string `form:"category"`
	//Material        *string       `form:"material"`
	Name         string                `form:"name"`
	Description  string                `form:"description"`
	Image        *multipart.FileHeader `form:"image" binding:"required"`
	Length       int32                 `form:"length"`
	Width        int32                 `form:"width"`
	Height       int32                 `form:"height"`
	Weight       int32                 `form:"weight"`
	CreationDate string                `form:"creation_date"`
}

func (h *ExhibitHandler) Create(c *gin.Context) {
	uid := c.GetInt64(auth.UserIdKey)

	var req CreateExhibitRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Image.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл слишком большой"})
		return
	}

	contentType := req.Image.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Только JPEG/PNG"})
		return
	}

	creationDate, err := time.Parse("2006-01-02", req.CreationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты (ожидается YYYY-MM-DD)"})
		return
	}

	ext := filepath.Ext(req.Image.Filename)
	newFilename := uuid.New().String() + ext
	savePath := filepath.Join(h.imagesPath, newFilename)

	if err := c.SaveUploadedFile(req.Image, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить изображение"})
		return
	}
	e, err := h.svc.Create(c.Request.Context(), req.Category, &uid, req.Name, req.Description, savePath,
		req.Length, req.Width, req.Height, req.Weight, creationDate)
	if err != nil {
		os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, e)
}

func (h *ExhibitHandler) GetById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	e, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *ExhibitHandler) List(c *gin.Context) {
	page := 0
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p >= 0 {
		page = p
	}
	size := 20
	if s, err := strconv.Atoi(c.Query("size")); err == nil && s >= 0 {
		size = s
	}
	resp, err := h.svc.List(c.Request.Context(), int64(page), int64(size))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

type ExhibitUpdateRequest struct {
	Id           int64                 `form:"id" binding:"required"`
	Category     string                `form:"category"`
	Name         string                `form:"name"`
	Description  string                `form:"description"`
	Image        *multipart.FileHeader `form:"image"`
	CreationDate string                `form:"creation_date"`
	Length       int32                 `form:"length"`
	Width        int32                 `form:"width"`
	Height       int32                 `form:"height"`
	Weight       int32                 `form:"weight"`
}

func (h *ExhibitHandler) Update(c *gin.Context) {
	uid := c.GetInt64(auth.UserIdKey)
	var req ExhibitUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creationDate, err := time.Parse("2006-01-02", req.CreationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты (ожидается YYYY-MM-DD)"})
		return
	}

	var imgPath = ""
	if req.Image != nil {
		slog.Debug("update exhibit", "image name", req.Image.Filename)
		slog.Debug("update exhibit", "image size", req.Image.Size)
		imgPath1, code, resp := h.storeImage(c, req.Image)
		if imgPath1 == "" {
			c.JSON(code, resp)
			return
		}
		imgPath = imgPath1

	}
	slog.Debug("update exhibit", "new image path", imgPath)
	e := &models.Exhibit{
		ID:              req.Id,
		Category:        req.Category,
		Name:            req.Name,
		Description:     req.Description,
		StorageLocation: imgPath,
		CreationDate:    creationDate,
		Length:          req.Length,
		Width:           req.Width,
		Height:          req.Height,
		Weight:          2,
		Status:          models.ApproveStatus,
	}
	err = h.svc.Update(c.Request.Context(), uid, e)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "exhibit not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, e)
}

func (h *ExhibitHandler) storeImage(c *gin.Context, img *multipart.FileHeader) (string, int, any) {
	if img.Size > 5<<20 {
		return "", http.StatusBadRequest, gin.H{"error": "Файл слишком большой"}
	}

	//contentType := img.Header.Get("Content-Type")
	//if contentType != "image/jpeg" && contentType != "image/png" {
	//	return "", http.StatusBadRequest, gin.H{"error": "Только JPEG/PNG"}
	//}

	ext := filepath.Ext(img.Filename)
	newFilename := uuid.New().String() + ext
	savePath := filepath.Join(h.imagesPath, newFilename)
	slog.Debug("store exhibit", "new image path", savePath)
	if err := c.SaveUploadedFile(img, savePath); err != nil {
		return "", http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить изображение"}
	}

	return savePath, 0, nil
}
