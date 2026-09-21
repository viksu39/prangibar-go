package controllers

import (
	"net/http"
	"strconv"
	"time"

	"prangibar-go/config"
	"prangibar-go/models"

	"github.com/gin-gonic/gin"
)

type LogsController struct{}

func NewLogsController() *LogsController {
	return &LogsController{}
}

// @Summary      List log aktivitas API
// @Tags         Logs
// @Security     BearerAuth
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        limit      query     int     false  "Items per page"  default(50)
// @Param        method     query     string  false  "Filter by HTTP method"
// @Param        statusCode query     int     false  "Filter by status code"
// @Param        adminId    query     int     false  "Filter by admin ID"
// @Success      200        {object}  map[string]interface{}
// @Router       /logs [get]
func (ctrl *LogsController) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	method := c.Query("method")
	statusCode, _ := strconv.Atoi(c.DefaultQuery("statusCode", "0"))
	adminID, _ := strconv.Atoi(c.DefaultQuery("adminId", "0"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	query := config.DB.Model(&models.ApiLog{})
	if method != "" {
		query = query.Where("method = ?", method)
	}
	if statusCode > 0 {
		query = query.Where("status_code = ?", statusCode)
	}
	if adminID > 0 {
		query = query.Where("admin_id = ?", adminID)
	}

	var total int64
	query.Count(&total)

	var logs []models.ApiLog
	if err := query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": gin.H{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// @Summary      Ringkasan statistik log API
// @Tags         Logs
// @Security     BearerAuth
// @Param        days  query     int  false  "Days back"  default(7)
// @Success      200   {object}  map[string]interface{}
// @Router       /logs/summary [get]
func (ctrl *LogsController) Summary(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days < 1 {
		days = 7
	}

	since := time.Now().AddDate(0, 0, -days)

	var totalCount int64
	config.DB.Model(&models.ApiLog{}).Where("created_at >= ?", since).Count(&totalCount)

	var byMethod []struct {
		Method string
		Count  int64
	}
	config.DB.Model(&models.ApiLog{}).
		Select("method, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("method").
		Order("count DESC").
		Scan(&byMethod)

	var byStatus []struct {
		StatusCode int
		Count      int64
	}
	config.DB.Model(&models.ApiLog{}).
		Select("status_code, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("status_code").
		Order("count DESC").
		Scan(&byStatus)

	var topPaths []struct {
		Path  string
		Count int64
	}
	config.DB.Model(&models.ApiLog{}).
		Select("path, COUNT(*) as count").
		Where("created_at >= ?", since).
		Group("path").
		Order("count DESC").
		Limit(10).
		Scan(&topPaths)

	c.JSON(http.StatusOK, gin.H{
		"period": gin.H{
			"days":  days,
			"since": since,
		},
		"totalRequests": totalCount,
		"byMethod":      byMethod,
		"byStatus":      byStatus,
		"topPaths":      topPaths,
	})
}
