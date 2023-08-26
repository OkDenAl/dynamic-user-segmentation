package segment

import (
	"dynamic-user-segmentation/internal/ports/httpgin/responses"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/service/segment"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateSegment(segmentService segment.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := segmentOperationRequest{}
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		err = segmentService.CreateSegment(c, req.Name, req.PercentOfUsers)
		if err != nil {
			switch err {
			case segment.ErrInvalidName:
				fallthrough
			case repository.ErrAlreadyExists:
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

func DeleteSegment(segmentService segment.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := segmentOperationRequest{}
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error(err))
			return
		}
		err = segmentService.DeleteSegment(c, req.Name)
		if err != nil {
			switch err {
			case segment.ErrInvalidName:
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
