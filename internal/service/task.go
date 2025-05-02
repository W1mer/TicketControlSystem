package service

import (
	"github.com/W1mer/TicketControlSystem/internal/repository"
)

// TaskService - сервис для работы с задачами
type TaskService struct {
	repo *repository.TaskRepository
}

// NewTaskService - конструктор сервиса
func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask - создание задачи
func (s *TaskService) CreateTask(title, status string) (*repository.Task, error) {
	task := &repository.Task{
		Title:  title,
		Status: status,
	}
	if err := s.repo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

// GetTask - получение задачи
func (s *TaskService) GetTask(id int) (*repository.Task, error) {
	return s.repo.GetByID(id)
}

// UpdateTask - обновление задачи
func (s *TaskService) UpdateTask(id int, title, status string) error {
	task := &repository.Task{
		ID:     id,
		Title:  title,
		Status: status,
	}
	return s.repo.Update(task)
}

// DeleteTask - удаление задачи
func (s *TaskService) DeleteTask(id int) error {
	return s.repo.Delete(id)
}
