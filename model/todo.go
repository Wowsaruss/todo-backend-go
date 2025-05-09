package model

import (
	"log"
	"database/sql"
	"todo-backend-go/db"
)

// ErrNoRows is returned when a requested record is not found
var ErrNoRows = sql.ErrNoRows

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Order     int    `json:"order"`
	UserId    *int   `json:"user_id"`
}

func GetAllTodos() ([]Todo, error) {
	log.Printf("Executing GetAllTodos query")
	rows, err := db.DB.Query(`SELECT id, title, completed, "order", user_id FROM todos`)
	if err != nil {
		log.Printf("Error querying todos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Order, &todo.UserId); err != nil {
			log.Printf("Error scanning todo row: %v", err)
			return nil, err
		}
		todos = append(todos, todo)
	}
	
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating todo rows: %v", err)
		return nil, err
	}
	
	log.Printf("Successfully retrieved %d todos", len(todos))
	return todos, nil
}

func CreateTodo(todo *Todo) error {
	log.Printf("Creating new todo with title: %s", todo.Title)
	query := `INSERT INTO todos (title, completed, "order", user_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.DB.QueryRow(query, todo.Title, todo.Completed, todo.Order, todo.UserId).Scan(&todo.ID)
	if err != nil {
		log.Printf("Error creating todo: %v", err)
		return err
	}
	log.Printf("Successfully created todo with ID: %d", todo.ID)
	return nil
}

func UpdateTodo(id int) (*Todo, error) {
	log.Printf("Updating todo with ID: %d", id)
	var todo Todo
	query := `
		UPDATE todos 
		SET completed = NOT completed 
		WHERE id = $1 
		RETURNING id, title, completed, "order", user_id`
	
	err := db.DB.QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.Order, &todo.UserId)
	if err == sql.ErrNoRows {
		log.Printf("No todo found with ID: %d", id)
		return nil, nil
	}
	if err != nil {
		log.Printf("Error updating todo: %v", err)
		return nil, err
	}
	log.Printf("Successfully updated todo with ID: %d", id)
	return &todo, nil
}

func DeleteTodo(id int) error {
	log.Printf("Deleting todo with ID: %d", id)
	query := "DELETE FROM todos WHERE id = $1"
	result, err := db.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting todo: %v", err)
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		log.Printf("No todo found with ID: %d", id)
		return sql.ErrNoRows
	}
	log.Printf("Successfully deleted todo with ID: %d", id)
	return nil
}
