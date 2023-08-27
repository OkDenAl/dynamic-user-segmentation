package segment

import (
	"dynamic-user-segmentation/internal/service/segment"
	"github.com/gin-gonic/gin"
)

func SetRouter(api *gin.RouterGroup, segmentService segment.Service) {
	api.POST("/segment/create", createSegment(segmentService))
	api.DELETE("/segment/delete", deleteSegment(segmentService))
}
