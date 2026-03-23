package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"task-management/app"
	"task-management/controller"
	"task-management/helper"
	"task-management/middleware"
	"task-management/model/domain"
	"task-management/repository"
	"task-management/service"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func newDBTest() *sql.DB {
	db, err := sql.Open("mysql", "root:P@ssw0rd@tcp(localhost:3306)/task_management_test")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setUpRouter(db *sql.DB) http.Handler {
	validate := validator.New()

	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository, db, validate)
	taskController := controller.NewTaskController(taskService)

	router := app.NewRouter(taskController)

	return middleware.NewAuthMiddleware(router)
}

func truncateTask(db *sql.DB) {
	db.Exec("truncate task")
}

func TestInsertTaskSuccess(t *testing.T) {
	db := newDBTest()
	truncateTask(db)
	router := setUpRouter(db)

	requestBody := strings.NewReader(`
	{
    	"title": "Self Learning - Go",
    	"description": "Learn how to make Go program with RESTful API based",
    	"status": "Done"
	}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/tasks", requestBody)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("X-API-Key", "SecretKey")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "Self Learning - Go", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, "Learn how to make Go program with RESTful API based", responseBody["data"].(map[string]interface{})["description"])
	assert.Equal(t, "Done", responseBody["data"].(map[string]interface{})["status"])
}

func TestInsertTaskFailed(t *testing.T) {
	db := newDBTest()
	truncateTask(db)
	router := setUpRouter(db)

	requestBody := strings.NewReader(`{"title": "", "description": "", "status": ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/tasks", requestBody)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("X-API-Key", "SecretKey")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "Bad Request", responseBody["status"])
}

func TestUpdateTaskSuccess(t *testing.T) {
	db := newDBTest()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	task := taskRepository.Insert(context.Background(), tx, domain.Task{
		Title: "Self Learning - Go",
		Description: "Learn how to make Go program with RESTful API based",
		Status: "Done",
	})
	tx.Commit()

	router := setUpRouter(db)

	requestBody := strings.NewReader(`
	{
    	"title": "Self Learning - Golang",
    	"description": "Learn how to make Golang program with RESTful API based",
    	"status": "Inprogress"
	}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/tasks/" + strconv.Itoa(task.Id), requestBody)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("X-API-Key", "SecretKey")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, task.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Self Learning - Golang", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, "Learn how to make Golang program with RESTful API based", responseBody["data"].(map[string]interface{})["description"])
	assert.Equal(t, "Inprogress", responseBody["data"].(map[string]interface{})["status"])
}

func TestUpdateTaskFailed(t *testing.T) {
	db := newDBTest()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	task := taskRepository.Insert(context.Background(), tx, domain.Task{
		Title: "Self Learning - Go",
		Description: "Learn how to make Go program with RESTful API based",
		Status: "Done",
	})
	tx.Commit()

	router := setUpRouter(db)

	requestBody := strings.NewReader(`
	{
    	"title": "",
    	"description": "",
    	"status": ""
	}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/tasks/" + strconv.Itoa(task.Id), requestBody)
	request.Header.Add("content-type", "application/json")
	request.Header.Add("X-API-Key", "SecretKey")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "Bad Request", responseBody["status"])
}

func TestDeleteTaskSuccess(t *testing.T) {
	db := newDBTest()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	task := taskRepository.Insert(context.Background(), tx, domain.Task{
		Title: "Self Learning - Go",
		Description: "Learn how to make Go program with RESTful API based",
		Status: "Done",
	})
	tx.Commit()

	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/tasks/" + strconv.Itoa(task.Id), nil)
	request.Header.Add("X-API-Key", "SecretKey")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestDeleteTaskFailed(t *testing.T) {
	db := newDBTest()
	truncateTask(db)
	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/tasks/1", nil)
	request.Header.Add("X-API-Key", "SecretKey")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "Not Found", responseBody["status"])
}

func TestGetTaskSuccess(t *testing.T) {
	db := newDBTest()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	task := taskRepository.Insert(context.Background(), tx, domain.Task{
		Title: "Self Learning - Go",
		Description: "Learn how to make Go program with RESTful API based",
		Status: "Done",
	})
	tx.Commit()

	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/tasks/" + strconv.Itoa(task.Id), nil)
	request.Header.Add("X-API-Key", "SecretKey")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, task.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, task.Title, responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, task.Description, responseBody["data"].(map[string]interface{})["description"])
	assert.Equal(t, task.Status, responseBody["data"].(map[string]interface{})["status"])
}

func TestGetTaskFailed(t *testing.T) {
	db := newDBTest()
	truncateTask(db)
	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/tasks/1", nil)
	request.Header.Add("X-API-Key", "SecretKey")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "Not Found", responseBody["status"])
}

func TestGetListTasks(t *testing.T) {
	db := newDBTest()
	truncateTask(db)

	tx, _ := db.Begin()
	taskRepository := repository.NewTaskRepository()
	task1 := taskRepository.Insert(context.Background(), tx, domain.Task{
		Title: "Self Learning - Go",
		Description: "Learn how to make Go program with RESTful API based",
		Status: "Done",
	})

	task2 := taskRepository.Insert(context.Background(), tx, domain.Task{
		Title: "Self Learning - Golang",
		Description: "Learn how to make Golang program with RESTful API based",
		Status: "Inprogress",
	})
	tx.Commit()

	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/tasks", nil)
	request.Header.Add("X-API-Key", "SecretKey")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	tasks := responseBody["data"].([]interface{})
	taskResponse1 := tasks[0].(map[string]interface{})
	taskResponse2 := tasks[1].(map[string]interface{})

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	assert.Equal(t, task1.Id, int(taskResponse1["id"].(float64)))
	assert.Equal(t, task1.Title, taskResponse1["title"])
	assert.Equal(t, task1.Description, taskResponse1["description"])
	assert.Equal(t, task1.Status, taskResponse1["status"])

	assert.Equal(t, task2.Id, int(taskResponse2["id"].(float64)))
	assert.Equal(t, task2.Title, taskResponse2["title"])
	assert.Equal(t, task2.Description, taskResponse2["description"])
	assert.Equal(t, task2.Status, taskResponse2["status"])
}

func TestUnauthorized(t *testing.T) {
	db := newDBTest()
	truncateTask(db)
	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/tasks", nil)
	request.Header.Add("X-API-Key", "Secret")

	recorder := httptest.NewRecorder()
	
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	assert.Equal(t, "Unauthorized", responseBody["status"])
}