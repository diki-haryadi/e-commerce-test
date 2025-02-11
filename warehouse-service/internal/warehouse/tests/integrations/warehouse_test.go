package artcileIntegrationTest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	articleDto "github.com/diki-haryadi/go-micro-template/internal/warehouse/dto"
	articleFixture "github.com/diki-haryadi/go-micro-template/internal/warehouse/tests/fixtures"
	httpError "github.com/diki-haryadi/ztools/error/http"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	fixture *articleFixture.IntegrationTestFixture
}

func (suite *testSuite) SetupSuite() {
	fixture, err := articleFixture.NewIntegrationTestFixture()
	if err != nil {
		assert.Error(suite.T(), err)
	}

	suite.fixture = fixture
}

func (suite *testSuite) TearDownSuite() {
	suite.fixture.TearDown()
}

func (suite *testSuite) TestSuccessHttpCreateWarehouse() {
	warehouseJSON := `{"shop_id":"550e8400-e29b-41d4-a716-446655440100","name":"Central Warehouse"}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/warehouse", strings.NewReader(warehouseJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusOK, response.Code)

	caDto := new(articleDto.CreateWarehouseRequestDto)
	if assert.NoError(suite.T(), json.Unmarshal(response.Body.Bytes(), caDto)) {
		assert.Equal(suite.T(), "Central Warehouse", caDto.Name)
	}

}

func (suite *testSuite) TestShopIDValidationErrHttpCreateWarehouse() {
	warehouseJSON := `{"shop_id":"","name":"Central Warehouse"}`

	request := httptest.NewRequest(http.MethodPost, "/api/v1/warehouse", strings.NewReader(warehouseJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.SetupDefaultMiddlewares()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusBadRequest, response.Code)

	httpErr := httpError.ParseExternalHttpErr(response.Result().Body)
	if assert.NotNil(suite.T(), httpErr) {
		assert.Equal(suite.T(), http.StatusBadRequest, httpErr.GetStatus())
		assert.Contains(suite.T(), httpErr.GetDetails(), "shop_id")
	}

}

func TestRunSuite(t *testing.T) {
	tSuite := new(testSuite)
	suite.Run(t, tSuite)
}
