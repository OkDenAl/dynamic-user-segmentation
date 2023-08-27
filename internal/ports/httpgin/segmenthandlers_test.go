package httpgin

import (
	"bytes"
	"dynamic-user-segmentation/internal/mocks/service/segmentserv_mocks"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/service/segment"
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

type succesResponse struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

type SegmentApiTestSuite struct {
	suite.Suite
	segmentService *segmentserv_mocks.Service
	client         *http.Client
	baseURL        string
}

func (suite *SegmentApiTestSuite) SetupTest() {
	suite.segmentService = &segmentserv_mocks.Service{}
	server := NewServer(":18080", suite.segmentService, nil, nil, logging.NewForMocks())
	testServer := httptest.NewServer(server.Handler)
	suite.client = testServer.Client()
	suite.baseURL = testServer.URL
}

func (suite *SegmentApiTestSuite) TestCreateSegment_OK() {
	suite.segmentService.On("CreateSegment", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string"),
		mock.AnythingOfType("float64")).Return(nil)
	body := map[string]any{
		"name":             "test",
		"percent_of_users": 0,
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/segment/create", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusCreated)
	var response succesResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Message, "success")
}

func (suite *SegmentApiTestSuite) TestCreateSegment_ErrAlreadyExists() {
	errCreateSuite(suite, repository.ErrAlreadyExists, http.StatusBadRequest)
}

func (suite *SegmentApiTestSuite) TestCreateSegment_ErrInvalidName() {
	errCreateSuite(suite, segment.ErrInvalidSegment, http.StatusBadRequest)
}

func (suite *SegmentApiTestSuite) TestCreateSegment_ErrInvalidPercentData() {
	errCreateSuite(suite, segment.ErrInvalidPercentData, http.StatusBadRequest)
}
func (suite *SegmentApiTestSuite) TestCreateSegment_UnexpectedError() {
	errCreateSuite(suite, errors.New("some error"), http.StatusInternalServerError)
}

func (suite *SegmentApiTestSuite) TestDeleteSegment_OK() {
	suite.segmentService.On("DeleteSegment", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string")).
		Return(nil)
	body := map[string]any{
		"name": "test",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/segment/delete", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var response succesResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Message, "success")
}

func (suite *SegmentApiTestSuite) TestDeleteSegment_ErrInvalidName() {
	errDeleteSuite(suite, segment.ErrInvalidSegment, http.StatusBadRequest)
}

func (suite *SegmentApiTestSuite) TestDeleteSegment_UnexpectedError() {
	errDeleteSuite(suite, errors.New("some error"), http.StatusInternalServerError)
}

func errCreateSuite(suite *SegmentApiTestSuite, err error, status int) {
	suite.segmentService.On("CreateSegment", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string"),
		mock.AnythingOfType("float64")).Return(err)
	body := map[string]any{}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/segment/create", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, status)
	var response errorResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Error, err.Error())
}

func errDeleteSuite(suite *SegmentApiTestSuite, err error, status int) {
	suite.segmentService.On("DeleteSegment", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("string")).
		Return(err)
	body := map[string]any{
		"name": "test",
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodDelete, suite.baseURL+"/api/v1/segment/delete", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, status)
	var response errorResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Error, err.Error())
}

func TestAdsApi(t *testing.T) {
	suite.Run(t, new(SegmentApiTestSuite))
}
