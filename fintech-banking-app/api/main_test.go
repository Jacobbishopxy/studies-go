package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// `cd api/`
// `go test -v`
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
