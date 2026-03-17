package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Todo update model params

type UpdateTodoParams struct {
	Id        uuid.UUID
	Title     *string
	Completed *bool
}

func makeTodoSetQuery(params UpdateTodoParams) updateSetQuery {
	query := make([]string, 0)
	var result updateSetQuery
	result.ArgId = 1

	if params.Title != nil {
		query = append(query, fmt.Sprintf("title=$%d", result.ArgId))
		result.Args = append(result.Args, *params.Title)
		result.ArgId++
	}
	if params.Completed != nil {
		query = append(query, fmt.Sprintf("completed=$%d", result.ArgId))
		result.Args = append(result.Args, *params.Completed)
		result.ArgId++
	}

	query = append(query, "updated_at=CURRENT_TIMESTAMP")
	result.Sql = strings.Join(query, ", ")
	result.Args = append(result.Args, params.Id)

	return result
}

type updateSetQuery struct {
	Sql   string
	Args  []interface{}
	ArgId int
}

// List update model params

type UpdateListParams struct {
	Id          uuid.UUID
	Title       *string
	Description *string
	UserId      uuid.UUID
}

func makeListSetQuery(params UpdateListParams) updateSetQuery {
	query := make([]string, 0)
	result := updateSetQuery{
		Args:  make([]interface{}, 0),
		ArgId: 1,
	}

	if params.Title != nil {
		query = append(query, fmt.Sprintf("title = $%d", result.ArgId))
		result.Args = append(result.Args, &params.Title)
		result.ArgId++
	}
	if params.Description != nil {
		query = append(query, fmt.Sprintf("description = $%d", result.ArgId))
		result.Args = append(result.Args, &params.Description)
		result.ArgId++
	}

	query = append(query, "updated_at = CURRENT_TIMESTAMP")
	result.Sql = strings.Join(query, ", ")
	result.Args = append(result.Args, params.Id, params.UserId)

	return result
}
