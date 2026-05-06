package db

import (
	"database/sql"
	"fmt"
	"time"

	"go_final_project/pkg/constants"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func UpdateDate(id string, newDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	result, err := db.Exec(query, newDate, id)
	if err != nil {
		return fmt.Errorf("issue date update error: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting the number of updated rows: %w", err)
	}

	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("issue deletion error: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting the number of updated rows: %w", err)
	}

	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil

}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func GetTask(id string) (*Task, error) {
	var task Task
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?`

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("issue update error: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}

func Tasks(search string, limit int) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	switch {
	case isDateQuery(search):
		date := convertDateFormat(search)
		query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE date = ?
		ORDER BY date
		LIMIT ?`
		rows, err = db.Query(query, date, limit)
	case search != "":
		pattern := "%" + search + "%"
		query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ?
		ORDER BY date
		LIMIT ?`
		rows, err = db.Query(query, pattern, pattern, limit)
	default:
		query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?`
		rows, err = db.Query(query, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting the task list: %w", err)
	}

	defer rows.Close()
	tasks := make([]*Task, 0)

	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("issue scanning error: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in processing the results: %w", err)
	}

	return tasks, nil
}

func isDateQuery(s string) bool {
	_, err := time.Parse(constants.InputDateFormat, s)
	return err == nil
}

func convertDateFormat(dateStr string) string {
	t, _ := time.Parse(constants.InputDateFormat, dateStr)
	return t.Format(constants.DateFormat)
}
