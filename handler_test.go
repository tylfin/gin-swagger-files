package handler_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gotest.tools/assert"

	swaggerFiles "github.com/tylfin/gin-swagger-files"
)

func performRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	return w
}

func TestWrapHandler(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)
}

func TestWrapCustomHandler(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/*any", ginSwagger.CustomWrapHandler(&ginSwagger.Config{}, swaggerFiles.Handler))

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)

	w3 := performRequest("GET", "/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)

	w4 := performRequest("GET", "/notfound", router)
	assert.Equal(t, 404, w4.Code)
}

func TestDisablingWrapHandler(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	disablingKey := "SWAGGER_DISABLE"

	router.GET("/simple/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, disablingKey))

	w1 := performRequest("GET", "/simple/index.html", router)
	assert.Equal(t, 200, w1.Code)

	w3 := performRequest("GET", "/simple/favicon-16x16.png", router)
	assert.Equal(t, 200, w3.Code)

	w4 := performRequest("GET", "/simple/notfound", router)
	assert.Equal(t, 404, w4.Code)

	os.Setenv(disablingKey, "true")

	router.GET("/disabling/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, disablingKey))

	w11 := performRequest("GET", "/disabling/index.html", router)
	assert.Equal(t, 404, w11.Code)

	w33 := performRequest("GET", "/disabling/favicon-16x16.png", router)
	assert.Equal(t, 404, w33.Code)

	w44 := performRequest("GET", "/disabling/notfound", router)
	assert.Equal(t, 404, w44.Code)
}

func TestDisablingCustomWrapHandler(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	disablingKey := "SWAGGER_DISABLE2"

	router.GET("/simple/*any", ginSwagger.DisablingCustomWrapHandler(&ginSwagger.Config{}, swaggerFiles.Handler, disablingKey))

	w1 := performRequest("GET", "/simple/index.html", router)
	assert.Equal(t, 200, w1.Code)

	os.Setenv(disablingKey, "true")

	router.GET("/disabling/*any", ginSwagger.DisablingCustomWrapHandler(&ginSwagger.Config{}, swaggerFiles.Handler, disablingKey))

	w11 := performRequest("GET", "/disabling/index.html", router)
	assert.Equal(t, 404, w11.Code)
}

func TestWithGzipMiddleware(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(gzip.Gzip(gzip.BestSpeed))

	router.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	w1 := performRequest("GET", "/index.html", router)
	assert.Equal(t, 200, w1.Code)
	assert.Equal(t, w1.Header()["Content-Type"][0], "text/html; charset=utf-8")

	w2 := performRequest("GET", "/swagger-ui.css", router)
	assert.Equal(t, 200, w2.Code)
	assert.Equal(t, w2.Header()["Content-Type"][0], "text/css; charset=utf-8")

	w3 := performRequest("GET", "/swagger-ui-bundle.js", router)
	assert.Equal(t, 200, w3.Code)
	assert.Equal(t, w3.Header()["Content-Type"][0], "application/javascript")
}
