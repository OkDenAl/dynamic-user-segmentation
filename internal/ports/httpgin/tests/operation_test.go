package tests

import (
	"bytes"
	"dynamic-user-segmentation/internal/mocks/service/operationserv_mocks"
	"dynamic-user-segmentation/internal/ports/httpgin"
	"dynamic-user-segmentation/internal/service"
	"dynamic-user-segmentation/pkg/logging"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type OperationApiTestSuite struct {
	suite.Suite
	operationService *operationserv_mocks.OperationService
	client           *http.Client
	baseURL          string
}

func (suite *OperationApiTestSuite) SetupTest() {
	suite.operationService = &operationserv_mocks.OperationService{}
	services := &service.Services{
		OperationService:   suite.operationService,
		SegmentService:     nil,
		UserSegmentService: nil,
	}
	server := httpgin.NewServer(":18080", services, logging.NewForMocks())
	testServer := httptest.NewServer(server.Handler)
	suite.client = testServer.Client()
	suite.baseURL = testServer.URL
}

func (suite *OperationApiTestSuite) TestGetReportLink_OK() {
	suite.operationService.On("MakeReportLink", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int"),
		mock.AnythingOfType("int")).
		Return("link", nil)
	body := map[string]any{
		"month": 8,
		"year":  2023,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/operations/report", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var response linkResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Link, "link")
}

func (suite *OperationApiTestSuite) TestGetReportLink_InternalServerError() {
	suite.operationService.On("MakeReportLink", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int"),
		mock.AnythingOfType("int")).
		Return("", errors.New("some error"))
	body := map[string]any{}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/operations/report", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
	var response errorResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Error, errors.New("some error").Error())
}

func TestOperationApi(t *testing.T) {
	suite.Run(t, new(OperationApiTestSuite))
}
