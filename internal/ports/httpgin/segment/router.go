package segment

import (
	"dynamic-user-segmentation/internal/service"
	"github.com/gin-gonic/gin"
)

func SetRouter(api *gin.RouterGroup, segmentService service.SegmentService) {
	api.POST("/segment/create", createSegment(segmentService))
	api.DELETE("/segment/delete", deleteSegment(segmentService))
}
