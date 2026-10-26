package handler

import (
	"Museum/internal/service"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AiHandler struct {
	svc *service.AiService
}

func NewAiHandler(svc *service.AiService) *AiHandler {
	return &AiHandler{svc: svc}
}

func (h *AiHandler) HandleText(c *gin.Context) {
	req := c.Query("text")
	c.JSON(http.StatusOK, "прикольчик "+strconv.Itoa(len(req)))
}

func (h *AiHandler) AiResponse(c *gin.Context) {
	req := c.Query("text")
	slog.Debug("aiResponse", "text", req)
	resp, err := h.svc.AiResponse(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"answer": resp})
}

func (h *AiHandler) RecognizeFromTextImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	contentType := mime.TypeByExtension(filepath.Ext(file.Filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	h1 := make(textproto.MIMEHeader)
	h1.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, file.Filename))
	h1.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = io.Copy(part, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	writer.Close()

	//part, err := writer.CreateFormFile("image", file.Filename) // поле "image"
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//_, err = io.Copy(part, src)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//writer.Close() // обязательно, добавляет closing boundary

	// Создаём запрос с правильным Content-Type
	req, err := http.NewRequestWithContext(c.Request.Context(), "POST",
		"http://26.86.186.191:4040/api/v1/analyze/",
		//"https://httpbin.org/post",
		//"http://localhost:8000/upload",
		body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляем
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	//if resp.StatusCode != http.StatusOK {
	//	c.JSON(resp.StatusCode, gin.H{"error": "failed to analyze image"})
	//	return
	//}
	//Обработка ответа
	type Result struct {
		Data struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
			Desc string `json:"ocr_full_text"`
		} `json:"data"`
	}
	var res Result
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	slog.Debug("aiRecognizeFromTextImage", "status", resp.Status)
	slog.Debug("aiRecognizeFromTextImage", "result", res)

	c.JSON(http.StatusOK, gin.H{"name": res.Data.Name.Value, "description": res.Data.Desc})
}
