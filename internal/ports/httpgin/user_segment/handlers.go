package user_segment

import (
	"dynamic-user-segmentation/internal/ports/httpgin/responses"
	"dynamic-user-segmentation/internal/service/user_segment"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func MakeOperationWithUsersSegment(userSegmentService user_segment.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := userSegmentOperationRequest{}
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}

		err = userSegmentService.AddSegmentsToUser(c, req.UserId, req.SegmentsToAdd)
		if err != nil {
			switch err {
			case user_segment.ErrInvalidUserId:
				c.JSON(http.StatusBadRequest, responses.Error(err))
				return
			default:
				c.JSON(http.StatusInternalServerError, responses.Error(err))
				return
			}
		}

		err = userSegmentService.DeleteSegmentsFromUser(c, req.UserId, req.SegmentsToDelete)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Error(err))
			return
		}
		c.JSON(http.StatusOK, responses.Ok())
	}
}

func GetAllSegmentsOfUser(userSegmentService user_segment.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := strconv.Atoi(c.Param("user_id"))
		segments, err := userSegmentService.GetAllUserSegments(c, int64(userId))
		if err != nil {
			switch err {
			case user_segment.ErrInvalidUserId:
				c.JSON(http.StatusBadRequest, responses.Error(err))
				return
			default:
				c.JSON(http.StatusInternalServerError, responses.Error(err))
				return
			}
		}
		c.JSON(http.StatusOK, responses.Data(segments))
	}
}
