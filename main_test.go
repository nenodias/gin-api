package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
var ALUNO_MOCK_NOME = "Nome do Aluno Teste"
var ALUNO_MOCK_CPF = "12345678901"
var ALUNO_MOCK_RG = "123456789"

func SetupTestes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	return routes
}

func CriaAlunoMock() {
	aluno := models.Aluno{Nome: ALUNO_MOCK_NOME, CPF: ALUNO_MOCK_CPF, RG: ALUNO_MOCK_RG}
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

func TestBuscaAlunoPorId(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupTestes()
	r.GET("/alunos/:id", controllers.BuscaAlunoPorID)
	url := fmt.Sprintf("/alunos/%d", ID)
	req, _ := http.NewRequest("GET", url, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(
		resposta, req,
	)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
	var alunoMock models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &alunoMock)
	assert.Equal(t, ALUNO_MOCK_NOME, alunoMock.Nome, "Deveriam ser iguais")
	assert.Equal(t, ALUNO_MOCK_CPF, alunoMock.CPF, "Deveriam ser iguais")
	assert.Equal(t, ALUNO_MOCK_RG, alunoMock.RG, "Deveriam ser iguais")
}

func TestDeletaAlunoPorId(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	r := SetupTestes()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	url := fmt.Sprintf("/alunos/%d", ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	resposta := httptest.NewRecorder()
	r.ServeHTTP(
		resposta, req,
	)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
}

func TestAtualizaAlunoPorId(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()
	r := SetupTestes()
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	url := fmt.Sprintf("/alunos/%d", ID)
	aluno := models.Aluno{
		Nome: ALUNO_MOCK_NOME + " Novo",
		CPF:  "11122244411",
		RG:   "112244556",
	}
	body, _ := json.Marshal(aluno)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	resposta := httptest.NewRecorder()
	r.ServeHTTP(
		resposta, req,
	)
	var alunoMock models.Aluno
	json.Unmarshal(resposta.Body.Bytes(), &alunoMock)
	assert.Equal(t, http.StatusOK, resposta.Code, "Deveriam ser iguais")
	assert.Equal(t, aluno.Nome, alunoMock.Nome, "Deveriam ser iguais")
	assert.Equal(t, aluno.CPF, alunoMock.CPF, "Deveriam ser iguais")
	assert.Equal(t, aluno.RG, alunoMock.RG, "Deveriam ser iguais")
}
