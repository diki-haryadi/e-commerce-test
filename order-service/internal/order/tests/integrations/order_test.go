package artcileIntegrationTest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	orderDto "github.com/diki-haryadi/go-micro-template/internal/order/dto"
	orderFixture "github.com/diki-haryadi/go-micro-template/internal/order/tests/fixtures"
	httpError "github.com/diki-haryadi/ztools/error/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	fixture *orderFixture.IntegrationTestFixture
}

func (suite *testSuite) SetupSuite() {
	fixture, err := orderFixture.NewIntegrationTestFixture()
	if err != nil {
		assert.Error(suite.T(), err)
	}

	suite.fixture = fixture
}

func (suite *testSuite) TearDownSuite() {
	suite.fixture.TearDown()
}

func (suite *testSuite) TestSuccessHttpCheckoutOrder() {
	orderJSON := `{"user_id":"550e8400-e29b-41d4-a716-44665544000a","items":{
	"product_id":"550e8400-e29b-41d4-a716-446655440000",
	"warehouse_id":"550e8400-e29b-41d4-a716-446655440200",
	"quantity":2,
	}}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/orders/checkout", strings.NewReader(orderJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusOK, response.Code)

	caDto := new(orderDto.CheckoutRequestDto)
	if assert.NoError(suite.T(), json.Unmarshal(response.Body.Bytes(), caDto)) {
		assert.Contains(suite.T(), "-", caDto.UserID)
	}

}

func (suite *testSuite) TestUserIDValidationErrHttpCheckoutOrder() {
	orderJSON := `{"user_id":"","items":{
	"product_id":"550e8400-e29b-41d4-a716-446655440000",
	"warehouse_id":"550e8400-e29b-41d4-a716-446655440200",
	"quantity":2,
	}}`
	request := httptest.NewRequest(http.MethodPost, "/api/v1/orders/checkout", strings.NewReader(orderJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.SetupDefaultMiddlewares()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)

	httpErr := httpError.ParseExternalHttpErr(response.Result().Body)
	if assert.NotNil(suite.T(), httpErr) {
		assert.Equal(suite.T(), http.StatusBadRequest, httpErr.GetStatus())
		assert.Contains(suite.T(), httpErr.GetDetails(), "user_id")
	}

}

func TestRunSuite(t *testing.T) {
	tSuite := new(testSuite)
	suite.Run(t, tSuite)
}
