package sqlstore

import (
	"fmt"
	"main/internal/app/model"
	"strings"
)

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) Create(task *model.Task) (int, error) {
	//for i, attachment := range task.Attachments {
	//	task.Attachments[i] = "'" + attachment + "'"
	//}
	attachments := fmt.Sprintf("{%v}", strings.Join(task.Attachments, ","))
	err := r.store.db.QueryRow("INSERT INTO public.app_tasks (header, body, deadline_date, author, attachments, groupname) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		task.Header, task.Body, task.DeadlineDate.Time, task.Author, attachments, task.Groupname).Scan(&task.ID)
	if err != nil {
		return 0, err
	}
	return task.ID, nil
}

func (r *TaskRepository) GetAll(groupname int) ([]model.Task, error) {
	var (
		tasks        []model.Task
		attachmentsB []byte
		attachments  []string
	)
	rows, err := r.store.db.Query("SELECT id, header, body, create_date, deadline_date, author, attachments FROM public.app_tasks WHERE groupname = $1", groupname)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // не забываем закрыть rows

	for rows.Next() {
		var task model.Task
		err = rows.Scan(&task.ID, &task.Header, &task.Body, &task.CreateDate, &task.DeadlineDate, &task.Author, &attachmentsB)
		if err != nil {
			return nil, err
		}
		attachmentsBStr := string(attachmentsB)
		attachmentsBStr = strings.ReplaceAll(attachmentsBStr, "{", "")
		attachmentsBStr = strings.ReplaceAll(attachmentsBStr, "}", "")
		attachmentsBStr = strings.ReplaceAll(attachmentsBStr, "\"\"", "")
		attachments = strings.Split(attachmentsBStr, ",")
		task.Attachments = attachments
		tasks = append(tasks, task)
	}

	// конвертируем []byte в []string
	// разделяем байты по запятой

	return tasks, nil
}

func (r *TaskRepository) Delete(id, groupname int) error {
	_, err := r.store.db.Exec("DELETE FROM public.app_tasks WHERE id = $1 AND groupname = $2", id, groupname)
	if err != nil {
		return err
	}
	return nil
}
