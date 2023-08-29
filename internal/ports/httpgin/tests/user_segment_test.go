package tests

import (
	"bytes"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/mocks/service/usersegmentserv"
	"dynamic-user-segmentation/internal/ports/httpgin"
	"dynamic-user-segmentation/internal/repository/dberrors"
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

type UserSegmentApiTestSuite struct {
	suite.Suite
	userSegmentService *usersegmentserv_mocks.UserSegmentService
	client             *http.Client
	baseURL            string
}

func (suite *UserSegmentApiTestSuite) SetupTest() {
	suite.userSegmentService = &usersegmentserv_mocks.UserSegmentService{}
	services := &service.Services{
		OperationService:   nil,
		SegmentService:     nil,
		UserSegmentService: suite.userSegmentService,
	}
	server := httpgin.NewServer(":18080", services, logging.NewForMocks())
	testServer := httptest.NewServer(server.Handler)
	suite.client = testServer.Client()
	suite.baseURL = testServer.URL
}

func (suite *UserSegmentApiTestSuite) TestMakeOperationWithUsersSegment_OK() {
	suite.userSegmentService.On("AddSegmentsToUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string"), mock.AnythingOfType("entity.TTL")).Return(nil)
	suite.userSegmentService.On("DeleteSegmentsFromUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).Return(nil)
	body := map[string]any{}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/user_segment/operation", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var response succesResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Message, "success")
}

func (suite *UserSegmentApiTestSuite) TestMakeOperationWithUsersSegment_ErrAlreadyExists() {
	errAddSegment(suite, dberrors.ErrAlreadyExists, http.StatusBadRequest)
}
func (suite *UserSegmentApiTestSuite) TestMakeOperationWithUsersSegment_ErrUnknownData() {
	errAddSegment(suite, dberrors.ErrUnknownData, http.StatusBadRequest)
}
func (suite *UserSegmentApiTestSuite) TestMakeOperationWithUsersSegment_ErrInvalidUserId() {
	errAddSegment(suite, service.ErrInvalidUserId, http.StatusBadRequest)
}
func (suite *UserSegmentApiTestSuite) TestMakeOperationWithUsersSegment_UnexpectedError() {
	errAddSegment(suite, errors.New("some error"), http.StatusInternalServerError)
}
func (suite *UserSegmentApiTestSuite) TestMakeOperationWithUsersSegment_DeleteSegmentsUnexpectedError() {
	suite.userSegmentService.On("AddSegmentsToUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string"), mock.AnythingOfType("entity.TTL")).Return(nil)
	suite.userSegmentService.On("DeleteSegmentsFromUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).Return(errors.New("some error"))
	body := map[string]any{}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/user_segment/operation", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusInternalServerError)
	var response errorResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Error, errors.New("some error").Error())
}

func (suite *UserSegmentApiTestSuite) TestGetAllSegmentsOfUser_OK() {
	suite.userSegmentService.On("GetAllUserSegments", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return([]entity.Segment{{Name: "test1"}, {Name: "test2"}}, nil)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/user_segment/1", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusOK)
	var response dataResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Error, "")
	suite.Equal(response.Data, []interface{}{map[string]interface{}{"name": "test1"}, map[string]interface{}{"name": "test2"}})
}

func (suite *UserSegmentApiTestSuite) TestGetAllSegmentsOfUser_ErrInvalidUserId() {
	suite.userSegmentService.On("GetAllUserSegments", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64")).
		Return(nil, service.ErrInvalidUserId)
	req, _ := http.NewRequest(http.MethodGet, suite.baseURL+"/api/v1/user_segment/-1", nil)
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, http.StatusBadRequest)
	var response errorResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Error, service.ErrInvalidUserId.Error())
}

func errAddSegment(suite *UserSegmentApiTestSuite, err error, status int) {
	suite.userSegmentService.On("AddSegmentsToUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string"), mock.AnythingOfType("entity.TTL")).Return(err)
	suite.userSegmentService.On("DeleteSegmentsFromUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("int64"),
		mock.AnythingOfType("string")).Return(nil)
	body := map[string]any{}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, suite.baseURL+"/api/v1/user_segment/operation", bytes.NewReader(data))
	resp, _ := suite.client.Do(req)
	suite.Equal(resp.StatusCode, status)
	var response errorResponse
	respBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &response)
	suite.Equal(response.Error, err.Error())
}
func TestUserSegmentApi(t *testing.T) {
	suite.Run(t, new(UserSegmentApiTestSuite))
}
