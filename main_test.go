package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nenodias/gin-api/controllers"
	"github.com/nenodias/gin-api/database"
	"github.com/nenodias/gin-api/models"
	"github.com/stretchr/testify/assert"
)

var ID int

func SetupTestes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	return routes
}

func CriaAlunoMock() {
	aluno := models.Aluno{Nome: "Nome do Aluno Teste", CPF: "12345678901", RG: "123456789"}
	database.DB.Create(&aluno)
	ID = int(aluno.ID)
}

func DeletaAlunoMock() {
	var aluno models.Aluno
	database.DB.Delete(&aluno, ID)
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

func TestExibeTodosOsAlunos(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupTestes()
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	req, _ := http.NewRequest("GET", "/alunos", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(
		resposta, req,
	)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
}

func TestBuscaPorCPF(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupTestes()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	req, _ := http.NewRequest("GET", "/alunos/cpf/12345678901", nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(
		resposta, req,
	)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
}
