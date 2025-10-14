package analytics

import (
	"github.com/K1la/sales-tracker/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"net/http"
)

func (h *Handler) GetAnalytics(c *ginext.Context) {
	var query dto.AnalyticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.log.Error().Err(err).Msg("Error parsing analytics query")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.GetAnalytics(c.Request.Context(), query)
	if err != nil {
		h.log.Error().Err(err).Msg("Error getting analytics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
