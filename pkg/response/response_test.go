package response

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	ErrorResponse(c.Writer, "error test", 400)

	assert.Equal(t, rr.Code, 400)
}

func TestSuccessResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)

	SuccessResponse(c.Writer, nil, 200)

	assert.Equal(t, rr.Code, 200)
}

func TestErrorResponseData(t *testing.T) {
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)

	ErrorResponseData(c.Writer, nil, "error test", 500)

	assert.Equal(t, rr.Code, 500)
}
