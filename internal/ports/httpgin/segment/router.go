package segment

import (
	"dynamic-user-segmentation/internal/service/segment"
	"github.com/gin-gonic/gin"
)

func SetRouter(api *gin.RouterGroup, segmentService segment.Service) {
	//api.GET("/balance/:user_id", GetBalance(balanceService))
	api.POST("/segment/create", CreateSegment(segmentService))
	api.DELETE("/segment/delete", DeleteSegment(segmentService))
}
