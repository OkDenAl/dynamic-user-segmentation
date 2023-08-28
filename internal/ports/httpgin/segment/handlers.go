package segment

import (
	"dynamic-user-segmentation/internal/ports/httpgin/responses"
	"dynamic-user-segmentation/internal/repository/dberrors"
	"dynamic-user-segmentation/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func createSegment(segmentService service.SegmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := segmentCreatingRequest{}
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		err = segmentService.CreateSegment(c, req.Name, req.PercentOfUsers)
		if err != nil {
			switch err {
			case service.ErrInvalidPercentData:
				fallthrough
			case service.ErrInvalidSegment:
				fallthrough
			case dberrors.ErrAlreadyExists:
				c.JSON(http.StatusBadRequest, responses.Error(err))
				return
			default:
				c.JSON(http.StatusInternalServerError, responses.Error(err))
				return
			}
		}
		c.JSON(http.StatusCreated, responses.Ok())
	}
}

func deleteSegment(segmentService service.SegmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := segmentDeletingRequest{}
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		err = segmentService.DeleteSegment(c, req.Name)
		if err != nil {
			switch err {
			case service.ErrInvalidSegment:
				c.JSON(http.StatusBadRequest, responses.Error(err))
				return
			default:
				c.JSON(http.StatusInternalServerError, responses.Error(err))
				return
			}
		}
		c.JSON(http.StatusOK, responses.Ok())
	}
}
