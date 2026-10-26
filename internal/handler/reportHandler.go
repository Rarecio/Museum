package handler

import (
	"Museum/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) Report(c *gin.Context) {
	f, err := h.svc.Report(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	buffer, _ := f.WriteToBuffer()
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=report.xlsx")
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())
}
