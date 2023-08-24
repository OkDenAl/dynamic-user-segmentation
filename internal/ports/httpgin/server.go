package httpgin

import (
	segmentPort "dynamic-user-segmentation/internal/ports/httpgin/segment"
	"dynamic-user-segmentation/internal/service/segment"
	"dynamic-user-segmentation/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewServer(port string, segmentService segment.Service, log logging.Logger) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	api := handler.Group("api/v1", LoggerMiddleware(log), gin.Recovery())
	{
		segmentPort.SetRouter(api, segmentService)
	}
	return &http.Server{Addr: port, Handler: handler}
}
