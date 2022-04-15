package seeders

import (
	"database/sql"
	"fmt"
)

type createTodoData struct {
	ID      int64
	Title   string
	Comment string
	UserID  int64
}

func CreateTodoData(con *sql.DB) error {
	insertDataList := [...]*createTodoData{
		{
			ID:      1,
			Title:   "todo1",
			Comment: "todo1のコメント",
			UserID:  1,
		},
		{
			ID:      2,
			Title:   "todo2",
			Comment: "todo2のコメント",
			UserID:  1,
		},
		{
			ID:      3,
			Title:   "todo3",
			Comment: "todo3のコメント",
			UserID:  2,
		},
	}

	for _, insertData := range insertDataList {
		ins, err = con.Prepare("INSERT INTO todos (id, title, comment, user_id) VALUES ($1,$2,$3,$4)")
		if err != nil {
			return err
		}
		_, err = ins.Exec(insertData.ID, insertData.Title, insertData.Comment, insertData.UserID)
		if err != nil {
			return err
		}
	}
	get, getErr := con.Prepare("SELECT setval('todos_id_seq', (SELECT MAX(id) FROM todos));")
	if getErr != nil {
		return getErr
	}
	_, err = get.Exec()
	if err != nil {
		return err
	}
	fmt.Println("create data at todos table")
	return nil
}
