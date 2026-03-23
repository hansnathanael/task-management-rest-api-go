package repository

import (
	"context"
	"database/sql"
	"errors"
	"task-management/helper"
	"task-management/model/domain"
)

type TaskRepositoryImpl struct {
}

func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImpl{}
}

// Insert implements [TaskRepository].
func (repository *TaskRepositoryImpl) Insert(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task {
	sql := "insert into task(title, description, status) values (?, ?, ?)"
	result, err := tx.ExecContext(ctx, sql, task.Title, task.Description, task.Status)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	task.Id = int(id)

	return task
}

// Update implements [TaskRepository].
func (repository *TaskRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, task domain.Task) domain.Task {
	sql := "update task set title = ?, description = ?, status = ? where id = ?"
	_, err := tx.ExecContext(ctx, sql, task.Title, task.Description, task.Status, task.Id)
	helper.PanicIfError(err)

	return task
}

// Delete implements [TaskRepository].
func (repository *TaskRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, task domain.Task) {
	sql := "delete from task where id = ?"
	_, err := tx.ExecContext(ctx, sql, task.Id)
	helper.PanicIfError(err)
}

// FindById implements [TaskRepository].
func (repository *TaskRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, taskId int) (domain.Task, error) {
	sql := "select id, title, description, status from task where id = ?"
	rows, err := tx.QueryContext(ctx, sql, taskId)
	helper.PanicIfError(err)
	defer rows.Close()

	task := domain.Task{}

	if rows.Next() {
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status)
		helper.PanicIfError(err)

		return task, nil
	} else {
		return task, errors.New("Task Not Found")
	}
}

// FindAll implements [TaskRepository].
func (repository *TaskRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Task {
	sql := "select id, title, description, status from task"
	rows, err := tx.QueryContext(ctx, sql)
	helper.PanicIfError(err)
	defer rows.Close()

	var tasks []domain.Task

	for rows.Next() {
		task := domain.Task{}
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status)
		helper.PanicIfError(err)
		tasks = append(tasks, task)
	}

	return tasks
}