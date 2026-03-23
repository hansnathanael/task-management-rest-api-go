package service

import (
	"context"
	"task-management/model/web"
)

type TaskService interface {
	Insert(ctx context.Context, request web.TaskCreateRequest) web.TaskResponse
	Update(ctx context.Context, request web.TaskUpdateRequest) web.TaskResponse
	Delete(ctx context.Context, taskId int)
	FindById(ctx context.Context, taskId int) web.TaskResponse
	FindAll(ctx context.Context) []web.TaskResponse
}