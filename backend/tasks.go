/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Points      int    `json:"points"`
	Category    string `json:"category"`
	Description string `json:"description"` // json
	Flag        string `json:"flag"`
	// Type = "default" | "docker"
	Visible bool `json:"visible"`
	Removed bool `json:"removed"`
}

type TasksRequest struct {
	Token string
	Task  Task
}

func tasks_new(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TasksRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	DB_conn, err := DB_connect()
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB connection error"))
		return
	}
	defer DB_conn.Close()

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "Wrong auth token"))
		return
	}

	user, err := getUserById(userId)

	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	if !user.IsAdmin {
		w.Write(CreateJsonResponse(false, "Only admin can manage tasks"))
		return
	}

	task := req.Task

	_, err = DB_conn.Query("INSERT INTO `tasks` (id, title, points, category, description, flag, visible, removed) VALUES (NULL, ?, ?, ?, ?, ?, ?, 0)", task.Title, task.Points, task.Category, task.Description, task.Flag, task.Visible)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while adding to the DB"))
		return
	}

	w.Write(CreateJsonResponse(true, ""))
}

func tasks_edit(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TasksRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	DB_conn, err := DB_connect()
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB connection error"))
		return
	}
	defer DB_conn.Close()

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "Wrong auth token"))
		return
	}

	user, err := getUserById(userId)

	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	if !user.IsAdmin {
		w.Write(CreateJsonResponse(false, "Only admin can manage tasks"))
		return
	}

	task := req.Task

	row := DB_conn.QueryRow("SELECT * FROM `tasks` WHERE `id` = ?", task.Id)
	var db_task Task
	err = row.Scan(&db_task.Id, &db_task.Title, &db_task.Points, &db_task.Category, &db_task.Description, &db_task.Flag, &db_task.Visible, &db_task.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Task with give id not found in the DB"))
		return
	}

	flag := db_task.Flag
	if task.Flag != "" {
		flag = task.Flag
	}

	_, err = DB_conn.Query("UPDATE `tasks` SET title = ?, points = ?, category = ?, description = ?, flag = ?, visible = ? WHERE `id` = ?", task.Title, task.Points, task.Category, task.Description, flag, task.Visible, task.Id)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while editting the DB"))
		return
	}

	w.Write(CreateJsonResponse(true, ""))
}

func tasks_remove(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TasksRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	DB_conn, err := DB_connect()
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB connection error"))
		return
	}
	defer DB_conn.Close()

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "Wrong auth token"))
		return
	}

	user, err := getUserById(userId)

	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	if !user.IsAdmin {
		w.Write(CreateJsonResponse(false, "Only admin can manage tasks"))
		return
	}

	task := req.Task

	row := DB_conn.QueryRow("SELECT * FROM `tasks` WHERE `id` = ?", task.Id)
	var db_task Task
	err = row.Scan(&db_task.Id, &db_task.Title, &db_task.Points, &db_task.Category, &db_task.Description, &db_task.Flag, &db_task.Visible, &db_task.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Task with give id not found in the DB"))
		return
	}

	_, err = DB_conn.Query("UPDATE `tasks` SET removed = 1 WHERE `id` = ?", task.Id)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while adding to the DB"))
		return
	}

	w.Write(CreateJsonResponse(true, ""))
}

type TasksGetRequest struct {
	Token  string `json:"token"`
	TaskId int    `json:"taskId"`
}

type TaskGetResponse struct {
	Task   Task `json:"task"`
	Solved bool `json:"solved"`
}

type SolvedTask struct {
	Id       int
	TaskId   int
	TeamId   int
	SolvedAt int64
}

func isTaskSolvedByTeam(taskId int, teamId int) bool {
	DB_conn, err := DB_connect()
	if err != nil {
		return false
	}
	defer DB_conn.Close()

	row := DB_conn.QueryRow("SELECT * FROM `solved_tasks` WHERE task_id = ? AND team_id = ?", taskId, teamId)
	var st SolvedTask
	err = row.Scan(&st.Id, &st.TaskId, &st.TeamId, &st.SolvedAt)
	if err != nil {
		return false
	}
	return true
}

func tasks_get_all(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TasksGetRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	DB_conn, err := DB_connect()
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB connection error"))
		return
	}
	defer DB_conn.Close()

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "You are not logged in"))
		return
	}

	user, err := getUserById(userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanViewTasks && !user.IsAdmin {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot view tasks"))
		return
	}

	teamId := getUserTeamId(userId)
	if teamId == -1 {
		w.Write(CreateJsonResponse(false, "You are not a member of any team"))
		return
	}

	var rows *sql.Rows

	if user.IsAdmin {
		rows, err = DB_conn.Query("SELECT * FROM `tasks` WHERE removed = 0")

	} else {
		rows, err = DB_conn.Query("SELECT * FROM `tasks` WHERE visible = 1 AND removed = 0")

	}

	if err != nil {
		w.Write(CreateJsonResponse(false, "DB error"))
		return
	}

	var tasks []Task
	for rows.Next() {
		var task Task
		rows.Scan(&task.Id, &task.Title, &task.Points, &task.Category, &task.Description, &task.Flag, &task.Visible, &task.Removed)
		if !user.IsAdmin {
			task.Flag = ""
		}
		tasks = append(tasks, task)
	}

	var res []TaskGetResponse
	for _, t := range tasks {
		res = append(res, TaskGetResponse{
			Task:   t,
			Solved: isTaskSolvedByTeam(t.Id, teamId),
		})
	}

	bytes, err = json.Marshal(res)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Json error"))
		return
	}

	w.Write(CreateJsonResponse(true, string(bytes)))
}

func tasks_get(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TasksGetRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	DB_conn, err := DB_connect()
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB connection error"))
		return
	}
	defer DB_conn.Close()

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "You are not logged in"))
		return
	}

	user, err := getUserById(userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanViewTasks && !user.IsAdmin {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot view tasks"))
		return
	}

	teamId := getUserTeamId(userId)
	if teamId == -1 {
		w.Write(CreateJsonResponse(false, "You are not a member of any team"))
		return
	}

	var row *sql.Row
	if user.IsAdmin {
		row = DB_conn.QueryRow("SELECT * FROM `tasks` WHERE removed = 0 AND id = ?", req.TaskId)
	} else {
		row = DB_conn.QueryRow("SELECT * FROM `tasks` WHERE visible = 1 AND removed = 0 AND id = ?", req.TaskId)
	}

	var task Task
	err = row.Scan(&task.Id, &task.Title, &task.Points, &task.Category, &task.Description, &task.Flag, &task.Visible, &task.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB error"))
		return
	}
	if !user.IsAdmin {
		task.Flag = ""
	}

	bytes, err = json.Marshal(TaskGetResponse{
		Task:   task,
		Solved: isTaskSolvedByTeam(task.Id, teamId),
	})
	if err != nil {
		w.Write(CreateJsonResponse(false, "Json error"))
		return
	}

	w.Write(CreateJsonResponse(true, string(bytes)))
}
