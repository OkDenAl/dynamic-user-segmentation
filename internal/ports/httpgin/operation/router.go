package operation

import (
	"dynamic-user-segmentation/internal/service/operation"
	"github.com/gin-gonic/gin"
)

func SetRouter(api *gin.RouterGroup, operationsService operation.Service) {
	api.GET("operations/report", getReportLink(operationsService))
}
