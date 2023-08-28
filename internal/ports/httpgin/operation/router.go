package operation

import (
	"dynamic-user-segmentation/internal/service"
	"github.com/gin-gonic/gin"
)

func SetRouter(api *gin.RouterGroup, operationsService service.OperationService) {
	api.GET("operations/report", getReportLink(operationsService))
}
