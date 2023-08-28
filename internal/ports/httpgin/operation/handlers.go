package operation

import (
	"dynamic-user-segmentation/internal/ports/httpgin/responses"
	"dynamic-user-segmentation/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getReportLink(operationsService service.OperationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := reportFileLinkRequest{}
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		link, err := operationsService.MakeReportLink(c, req.Month, req.Year)
		if err != nil {
			switch err {
			default:
				c.JSON(http.StatusInternalServerError, responses.Error(err))
				return
			}
		}
		c.JSON(http.StatusOK, responses.Link(link))
	}
}
