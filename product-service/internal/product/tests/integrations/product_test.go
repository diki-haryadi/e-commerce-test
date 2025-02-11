package artcileIntegrationTest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	articleDto "github.com/diki-haryadi/go-micro-template/internal/product/dto"
	articleFixture "github.com/diki-haryadi/go-micro-template/internal/product/tests/fixtures"
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

func (suite *testSuite) TestSuccessHttpGetProduct() {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/product", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusOK, response.Code)

	caDto := new(articleDto.GetProductsResponseDto)
	if assert.NoError(suite.T(), json.Unmarshal(response.Body.Bytes(), caDto)) {
		assert.Equal(suite.T(), "Gaming Laptop X1", caDto.Products[0].Name)
		assert.Equal(suite.T(), "High-performance gaming laptop with RTX 3080", caDto.Products[0].Description)
	}
}

func (suite *testSuite) TestSuccessHttpGetProductByID() {
	request := httptest.NewRequest(http.MethodGet, "/api/v1/product/550e8400-e29b-41d4-a716-446655440000", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	response := httptest.NewRecorder()
	suite.fixture.InfraContainer.EchoHttpServer.GetEchoInstance().ServeHTTP(response, request)

	assert.Equal(suite.T(), http.StatusOK, response.Code)

	caDto := new(articleDto.ProductResponse)
	if assert.NoError(suite.T(), json.Unmarshal(response.Body.Bytes(), caDto)) {
		assert.Equal(suite.T(), "Gaming Laptop X1", caDto.Name)
		assert.Equal(suite.T(), "High-performance gaming laptop with RTX 3080", caDto.Description)
	}
}

func TestRunSuite(t *testing.T) {
	tSuite := new(testSuite)
	suite.Run(t, tSuite)
}
