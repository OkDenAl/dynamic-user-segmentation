package user_segment

import (
	"dynamic-user-segmentation/internal/service"
	"github.com/gin-gonic/gin"
)

func SetRouter(api *gin.RouterGroup, userSegmentService service.UserSegmentService) {
	api.GET("/user_segment/:user_id", getAllSegmentsOfUser(userSegmentService))
	api.POST("/user_segment/operation", makeOperationWithUsersSegment(userSegmentService))
}
