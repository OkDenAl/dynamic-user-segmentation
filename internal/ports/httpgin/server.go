package httpgin

import (
	operationPort "dynamic-user-segmentation/internal/ports/httpgin/operation"
	segmentPort "dynamic-user-segmentation/internal/ports/httpgin/segment"
	userSegmentPort "dynamic-user-segmentation/internal/ports/httpgin/user_segment"
	"dynamic-user-segmentation/internal/service"
	"dynamic-user-segmentation/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewServer(port string, services *service.Services, log logging.Logger) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	api := handler.Group("api/v1", LoggerMiddleware(log), gin.Recovery())
	{
		segmentPort.SetRouter(api, services.SegmentService)
		userSegmentPort.SetRouter(api, services.UserSegmentService)
		operationPort.SetRouter(api, services.OperationService)
	}
	return &http.Server{Addr: port, Handler: handler}
}
