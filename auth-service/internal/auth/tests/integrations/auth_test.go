package authIntegrationTest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	authDto "github.com/diki-haryadi/go-micro-template/internal/auth/dto"
	authFixture "github.com/diki-haryadi/go-micro-template/internal/auth/tests/fixtures"
	httpError "github.com/diki-haryadi/ztools/error/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	fixture *authFixture.IntegrationTestFixture
}

func (suite *testSuite) SetupSuite() {
	fixture, err := authFixture.NewIntegrationTestFixture()
	if err != nil {
		assert.Error(suite.T(), err)
	}

	suite.fixture = fixture
}

func (suite *testSuite) TearDownSuite() {
	suite.fixture.TearDown()
}

func (suite *testSuite) TestSuccessHttpSignUp() {
	userJSON := `{"username":"jhon@gmail.com","password":"password"}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", strings.NewReader(userJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusOK, response.Code)

	caDto := new(authDto.SignUpRequestDto)
	if assert.NoError(suite.T(), json.Unmarshal(response.Body.Bytes(), caDto)) {
		assert.Equal(suite.T(), "jhon@gmail.com", caDto.Username)
	}

}

func (suite *testSuite) TestUsernameValidationErrHttpSignUp() {
	userJSON := `{"username":"","password":"password"}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", strings.NewReader(userJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.SetupDefaultMiddlewares()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)

	httpErr := httpError.ParseExternalHttpErr(response.Result().Body)
	if assert.NotNil(suite.T(), httpErr) {
		assert.Equal(suite.T(), http.StatusBadRequest, httpErr.GetStatus())
		assert.Contains(suite.T(), httpErr.GetDetails(), "username")
	}
}

func (suite *testSuite) TestPasswordValidationErrHttpSignUp() {
	userJSON := `{"username":"username","password":""}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", strings.NewReader(userJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.SetupDefaultMiddlewares()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)

	httpErr := httpError.ParseExternalHttpErr(response.Result().Body)
	if assert.NotNil(suite.T(), httpErr) {
		assert.Equal(suite.T(), http.StatusBadRequest, httpErr.GetStatus())
		assert.Contains(suite.T(), httpErr.GetDetails(), "passworod")
	}
}

func (suite *testSuite) TestSuccessHttpSignIn() {
	userJSON := `{"username":"jhon@gmail.com","password":"password"}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signin", strings.NewReader(userJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusOK, response.Code)

	caDto := new(authDto.SignUpRequestDto)
	if assert.NoError(suite.T(), json.Unmarshal(response.Body.Bytes(), caDto)) {
		assert.Equal(suite.T(), "jhon@gmail.com", caDto.Username)
	}

}

func (suite *testSuite) TestUsernameValidationErrHttpSignIn() {
	userJSON := `{"username":"","password":"password"}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", strings.NewReader(userJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.SetupDefaultMiddlewares()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)

	httpErr := httpError.ParseExternalHttpErr(response.Result().Body)
	if assert.NotNil(suite.T(), httpErr) {
		assert.Equal(suite.T(), http.StatusBadRequest, httpErr.GetStatus())
		assert.Contains(suite.T(), httpErr.GetDetails(), "username")
	}
}

func (suite *testSuite) TestPasswordValidationErrHttpSignIn() {
	userJSON := `{"username":"jhon@gmail.com","password":""}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signin", strings.NewReader(userJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.SetupDefaultMiddlewares()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)

	httpErr := httpError.ParseExternalHttpErr(response.Result().Body)
	if assert.NotNil(suite.T(), httpErr) {
		assert.Equal(suite.T(), http.StatusBadRequest, httpErr.GetStatus())
		assert.Contains(suite.T(), httpErr.GetDetails(), "passworod")
	}
}

func TestRunSuite(t *testing.T) {
	tSuite := new(testSuite)
	suite.Run(t, tSuite)
}
