package service

import (
	"context"
	"database/sql"
	"task-management/exception"
	"task-management/helper"
	"task-management/model/domain"
	"task-management/model/web"
	"task-management/repository"
	"github.com/go-playground/validator/v10"
)

type TaskServiceImpl struct {
	TaskRepository repository.TaskRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewTaskService(taskRepository repository.TaskRepository, db *sql.DB, validate *validator.Validate) TaskService {
	return &TaskServiceImpl{
		TaskRepository: taskRepository,
		DB:             db,
		Validate:       validate,
	}
}

// Insert implements [TaskService].
func (service *TaskServiceImpl) Insert(ctx context.Context, request web.TaskCreateRequest) web.TaskResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	task := domain.Task{
		Title: request.Title,
		Description: request.Description,
		Status: request.Status,
	}

	task = service.TaskRepository.Insert(ctx, tx, task)

	return helper.ToTaskResponse(task)
}

// Update implements [TaskService].
func (service *TaskServiceImpl) Update(ctx context.Context, request web.TaskUpdateRequest) web.TaskResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	task, err := service.TaskRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	task.Title = request.Title
	task.Description = request.Description
	task.Status = request.Status
	task = service.TaskRepository.Update(ctx, tx, task)

	return helper.ToTaskResponse(task)
}

// Delete implements [TaskService].
func (service *TaskServiceImpl) Delete(ctx context.Context, taskId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	task, err := service.TaskRepository.FindById(ctx, tx, taskId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.TaskRepository.Delete(ctx, tx, task)
}

// FindById implements [TaskService].
func (service *TaskServiceImpl) FindById(ctx context.Context, taskId int) web.TaskResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	task, err := service.TaskRepository.FindById(ctx, tx, taskId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToTaskResponse(task)
}

// FindAll implements [TaskService].
func (service *TaskServiceImpl) FindAll(ctx context.Context) []web.TaskResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	tasks := service.TaskRepository.FindAll(ctx, tx)

	return helper.ToTaskResponses(tasks)
}