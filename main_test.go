package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nenodias/gin-api/controllers"
	"github.com/stretchr/testify/assert"
)

func SetupTestes() *gin.Engine {
	routes := gin.Default()
	return routes
}

func TestVerificaStatusCodeSaudacao(t *testing.T) {
	r := SetupTestes()
	r.GET("/:nome", controllers.Saudacao)
	req, _ := http.NewRequest("GET", "/Neno", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(
		resposta, req,
	)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
	mockDaResposta := `{"API diz:":"E ai Neno, tudo beleza?"}`
	respostaBody, _ := ioutil.ReadAll(resposta.Body)
	assert.Equal(t, mockDaResposta, string(respostaBody), "Deveriam ser iguais")
}
