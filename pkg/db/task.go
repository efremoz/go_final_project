package db

import (
	"database/sql"
	"fmt"
	"go_final_project/pkg/constants"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)`

	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = result.LastInsertId()
	}

	return id, err
}

func Tasks(search string, limit int) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	if isDateQuery(search) {
		date := convertDateFormat(search)
		query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE date = ?
		ORDER BY date
		LIMIT ?`
		rows, err = db.Query(query, date, limit)
	} else if search != "" {
		pattern := "%" + search + "%"
		query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ?
		ORDER BY date
		LIMIT ?`
		rows, err = db.Query(query, pattern, pattern, limit)
	} else {
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
