package repository

import (
	"database/sql"
	"time"
)

// Task - структура задачи
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskRepository - репозиторий для работы с задачами
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository - конструктор репозитория
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create - создание задачи
func (r *TaskRepository) Create(task *Task) error {
	query := `INSERT INTO tasks (title, status) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.QueryRow(query, task.Title, task.Status).Scan(&task.ID, &task.CreatedAt)
}

// GetByID - получение задачи по ID
func (r *TaskRepository) GetByID(id int) (*Task, error) {
	task := &Task{}
	query := `SELECT id, title, status, created_at FROM tasks WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Update - обновление задачи
func (r *TaskRepository) Update(task *Task) error {
	query := `UPDATE tasks SET title = $1, status = $2 WHERE id = $3`
	_, err := r.db.Exec(query, task.Title, task.Status, task.ID)
	return err
}

// Delete - удаление задачи
func (r *TaskRepository) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
