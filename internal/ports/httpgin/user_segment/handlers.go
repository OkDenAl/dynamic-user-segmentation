package user_segment

import (
	"dynamic-user-segmentation/internal/ports/httpgin/responses"
	"dynamic-user-segmentation/internal/repository/dberrors"
	"dynamic-user-segmentation/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func makeOperationWithUsersSegment(userSegmentService service.UserSegmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := userSegmentOperationRequest{}
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}

		err = userSegmentService.AddSegmentsToUser(c, req.UserId, req.SegmentsToAdd, req.ExpiresAt)
		if err != nil {
			switch err {
			case dberrors.ErrAlreadyExists:
				fallthrough
			case dberrors.ErrUnknownData:
				fallthrough
			case service.ErrInvalidUserId:
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

func getAllSegmentsOfUser(userSegmentService service.UserSegmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := strconv.Atoi(c.Param("user_id"))
		segments, err := userSegmentService.GetAllUserSegments(c, int64(userId))
		if err != nil {
			switch err {
			case service.ErrInvalidUserId:
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
