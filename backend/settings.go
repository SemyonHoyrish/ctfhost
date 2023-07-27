/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

package main

import (
	"encoding/json"
	"net/http"
)

type Settings struct {
	Id                   int
	UserCanRegister      bool
	UserCanManageTeam    bool
	UserCanManageAccount bool
	UserCanViewTasks     bool
	MaxMembersPerTeam    int
	MainPageContent      string
	ChangedBy            int
}

func getSettings() (Settings, error) {
	DB_conn, err := DB_connect()
	if err != nil {
		return Settings{}, err
	}
	defer DB_conn.Close()

	row := DB_conn.QueryRow("SELECT * from `settings` ORDER BY id DESC LIMIT 1")
	var settings Settings
	err = row.Scan(&settings.Id, &settings.UserCanRegister, &settings.UserCanManageTeam, &settings.UserCanManageAccount, &settings.UserCanViewTasks, &settings.MaxMembersPerTeam, &settings.MainPageContent, &settings.ChangedBy)
	if err != nil {
		return Settings{}, err
	}

	return settings, nil
}

func setSettings(settings Settings, userId int) error {
	DB_conn, err := DB_connect()
	if err != nil {
		return err
	}
	defer DB_conn.Close()

	_, err = DB_conn.Query("INSERT INTO `settings` (id, user_can_register, user_can_manage_team, user_can_manage_account, user_can_view_tasks, max_members_per_team, main_page_content, changed_by) VALUES (NULL, ?, ?, ?, ?, ?, ?, ?)", settings.UserCanRegister, settings.UserCanManageTeam, settings.UserCanManageAccount, settings.UserCanViewTasks, settings.MaxMembersPerTeam, settings.MainPageContent, userId)

	return err
}

type SettingsRequest struct {
	Token   string   `json:"token"`
	Setting Settings `json:"settings"`
}

func settings_get(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req SettingsRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	user, err := getUserById(userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	if !user.IsAdmin {
		w.Write(CreateJsonResponse(false, "Only admin can request this"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while gettings settings"))
		return
	}

	bytes, err = json.Marshal(settings)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Json error"))
		return
	}

	w.Write(CreateJsonResponse(true, string(bytes)))
}

func settings_set(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req SettingsRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	user, err := getUserById(userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	if !user.IsAdmin {
		w.Write(CreateJsonResponse(false, "Only admin can do this"))
		return
	}

	err = setSettings(req.Setting, userId)

	if err != nil {
		w.Write(CreateJsonResponse(false, "DB error"))
		return
	}

	w.Write(CreateJsonResponse(true, "Settings changed"))
}

func settings_get_main_page_content(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req SettingsRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	_, err = getUserById(userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while gettings settings"))
		return
	}

	w.Write(CreateJsonResponse(true, settings.MainPageContent))
}
