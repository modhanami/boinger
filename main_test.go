package main

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}
