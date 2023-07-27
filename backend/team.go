/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Team struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	CaptainId int    `json:"captainId"`
	Score     int    `json:"score"`
	Removed   bool   `json:"removed"`
}

type TeamEditRequest struct {
	TeamName    string `json:"teamName"`
	NewTeamName string `json:"newTeamName"`
	Token       string `json:"token"`
}

func team_create(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TeamEditRequest
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
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanManageTeam {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot manage teams"))
		return
	}

	row := DB_conn.QueryRow("SELECT * FROM `teams` WHERE `name` = ? AND removed = 0", req.NewTeamName)
	var db_team Team
	err = row.Scan(&db_team.Id, &db_team.Name, &db_team.CaptainId, &db_team.Score, &db_team.Removed)
	if err == nil {
		w.Write(CreateJsonResponse(false, "Team with given name already exists"))
		return
	}

	teamId := getUserTeamId(userId)
	if teamId != -1 {
		w.Write(CreateJsonResponse(false, "You already in a team"))
		return
	}

	_, err = DB_conn.Query("INSERT INTO `teams` (id, name, captain_id, score) VALUES (NULL, ?, ?, 0)", req.NewTeamName, userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while changing DB data (1)"))
		return
	}

	row = DB_conn.QueryRow("SELECT * FROM `teams` WHERE `name` = ? AND removed = 0", req.NewTeamName)
	err = row.Scan(&db_team.Id, &db_team.Name, &db_team.CaptainId, &db_team.Score, &db_team.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	_, err = DB_conn.Query("INSERT INTO `team_members` (id, team_id, user_id, removed) VALUES (NULL, ?, ?, 0)", db_team.Id, userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while changing DB data (2)"))
		return
	}

	w.Write(CreateJsonResponse(true, "Team created successfully"))
}

func team_edit_name(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TeamEditRequest
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
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanManageTeam {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot manage teams"))
		return
	}

	{
		row := DB_conn.QueryRow("SELECT * FROM `teams` WHERE `name` = ? AND removed = 0", req.NewTeamName)
		var db_team Team
		err = row.Scan(&db_team.Id, &db_team.Name, &db_team.CaptainId, &db_team.Score, &db_team.Removed)
		if err == nil {
			w.Write(CreateJsonResponse(false, "Team with given name already exists"))
			return
		}
	}

	row := DB_conn.QueryRow("SELECT * FROM `teams` WHERE `name` = ? AND removed = 0", req.TeamName)
	var db_team Team
	err = row.Scan(&db_team.Id, &db_team.Name, &db_team.CaptainId, &db_team.Score, &db_team.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Team do not exists"))
		return
	}

	if userId != db_team.CaptainId {
		w.Write(CreateJsonResponse(false, "Only captain can change team's name"))
		return
	}

	_, err = DB_conn.Query("UPDATE `teams` SET name = ? WHERE id = ?", req.NewTeamName, db_team.Id)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while changing DB data"))
		return
	}

	w.Write(CreateJsonResponse(true, "Team's name changed successfully"))
}

func team_remove(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TeamEditRequest
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
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanManageTeam {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot manage teams"))
		return
	}

	var row *sql.Row
	if req.TeamName == "" {
		tId := getUserTeamId(userId)
		row = DB_conn.QueryRow("SELECT * FROM `teams` WHERE `id` = ? AND removed = 0", tId)
	} else {
		row = DB_conn.QueryRow("SELECT * FROM `teams` WHERE `name` = ? AND removed = 0", req.TeamName)

	}
	var db_team Team
	err = row.Scan(&db_team.Id, &db_team.Name, &db_team.CaptainId, &db_team.Score, &db_team.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Team do not exists"))
		return
	}

	if userId != db_team.CaptainId {
		w.Write(CreateJsonResponse(false, "Only captain can remove team"))
		return
	}

	_, err = DB_conn.Query("UPDATE `teams` SET removed = 1 WHERE id = ? AND removed = 0", db_team.Id)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while changing DB data (1)"))
		return
	}

	_, err = DB_conn.Query("UPDATE `team_members` SET removed = 1 WHERE team_id = ? AND removed = 0", db_team.Id)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while changing DB data (2)"))
		return
	}

	w.Write(CreateJsonResponse(true, "Team removed successfully"))
}

type TeamMember struct {
	Id      int
	TeamId  int
	UserId  int
	Removed bool
}

func getUserTeamId(userId int) int {
	DB_conn, err := DB_connect()
	if err != nil {
		return -1 // @Maybe: return other value on error??
	}
	defer DB_conn.Close()

	row := DB_conn.QueryRow("SELECT * FROM `team_members` WHERE user_id = ? AND removed = 0 ORDER BY id DESC", userId)
	var teamMember TeamMember
	err = row.Scan(&teamMember.Id, &teamMember.TeamId, &teamMember.UserId, &teamMember.Removed)
	if err != nil {
		return -1
	}
	return teamMember.TeamId
}

func getTeamById(teamId int) (Team, error) {
	DB_conn, err := DB_connect()
	if err != nil {
		return Team{}, err
	}
	defer DB_conn.Close()

	row := DB_conn.QueryRow("SELECT * FROM `teams` WHERE id = ? AND removed = 0", teamId)
	var team Team
	err = row.Scan(&team.Id, &team.Name, &team.CaptainId, &team.Score, &team.Removed)
	if err != nil {
		return Team{}, err
	}
	return team, nil
}

func getTeamByName(teamName string) (Team, error) {
	DB_conn, err := DB_connect()
	if err != nil {
		return Team{}, err
	}
	defer DB_conn.Close()

	row := DB_conn.QueryRow("SELECT * FROM `teams` WHERE name = ? AND removed = 0", teamName)
	var team Team
	err = row.Scan(&team.Id, &team.Name, &team.CaptainId, &team.Score, &team.Removed)
	if err != nil {
		return Team{}, err
	}
	return team, nil
}

func getTeamMembers(teamId int) []ResponseUser {
	members := make([]TeamMember, 0)

	DB_conn, err := DB_connect()
	if err != nil {
		return make([]ResponseUser, 0)
	}
	defer DB_conn.Close()

	rows, err := DB_conn.Query("SELECT * FROM `team_members` WHERE team_id = ? AND removed = 0", teamId)
	if err != nil {
		return make([]ResponseUser, 0)
	}
	for rows.Next() {
		var member TeamMember
		err := rows.Scan(&member.TeamId, &member.TeamId, &member.UserId, &member.Removed)
		if err == nil {
			members = append(members, member)
		}
	}

	users := make([]ResponseUser, len(members))
	for i, m := range members {
		u, err := getUserById(m.UserId)
		if err == nil {
			users[i] = ResponseUser{
				Id:       u.Id,
				Nickname: u.Nickname,
				Email:    u.Email,
				IsAdmin:  u.IsAdmin,
			}
		} else {
			users[i] = ResponseUser{}
		}
	}

	return users
}

func getTeamScore(teamId int) (int, error) {
	DB_conn, err := DB_connect()
	if err != nil {
		return 0, err
	}
	defer DB_conn.Close()

	rows, err := DB_conn.Query("SELECT task_id FROM `solved_tasks` WHERE team_id = ?", teamId)
	if err != nil {
		return 0, err
	}
	result := 0
	for rows.Next() {
		var taskId int
		err = rows.Scan(&taskId)
		if err != nil {
			return 0, err
		}
		row := DB_conn.QueryRow("SELECT points FROM `tasks` WHERE id = ? AND visible = 1 AND removed = 0", taskId)
		var points int
		err = row.Scan(&points)
		if err == nil {
			result += points
		}
	}

	{ // @Maybe remove, but now it's possible that team solved task, but second DB query return error and in `teams` table there is no change
		_, err := DB_conn.Query("UPDATE `teams` SET score = ? WHERE id = ? AND removed = 0", result, teamId)
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}

type SendInviteRequest struct {
	Token        string `json:"token"`
	UserNickname string `json:"userNickname"`
	UserEmail    string `jsion:"userEmail"`
}

func team_send_invite(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req SendInviteRequest
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
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanManageTeam {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot manage teams"))
		return
	}

	teamId := getUserTeamId(userId)
	if teamId == -1 {
		w.Write(CreateJsonResponse(false, "You are not member of any team"))
		return
	}

	if settings.MaxMembersPerTeam != -1 &&
		len(getTeamMembers(teamId))+len(getTeamUserInvites(teamId)) >= settings.MaxMembersPerTeam {
		w.Write(CreateJsonResponse(false, "You have reached limit of memers in the team (including invites)"))
		return
	}

	db_team, err := getTeamById(teamId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Team do not exists"))
		return
	}

	if userId != db_team.CaptainId {
		w.Write(CreateJsonResponse(false, "Only captain can invite users to the team"))
		return
	}

	var invitedUser User
	if req.UserEmail != "" {
		invitedUser, err = getUserByEmail(req.UserEmail)
	} else if req.UserNickname != "" {
		invitedUser, err = getUserByNickname(req.UserNickname)
	} else {
		w.Write(CreateJsonResponse(false, "You should specify nickname or email of user you want to invite"))
		return
	}
	if err != nil {
		w.Write(CreateJsonResponse(false, "Requested user not found"))
		return
	}

	for _, i := range getUserTeamInvites(invitedUser.Id) {
		if i.TeamId == db_team.Id {
			w.Write(CreateJsonResponse(false, "You already invited this user"))
			return
		}
	}

	_, err = DB_conn.Query("INSERT INTO `team_invites` (id, team_id, user_id, accepted, declined, removed) VALUES (NULL, ?, ?, 0, 0, 0)", teamId, invitedUser.Id)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while editing DB"))
		return
	}

	w.Write(CreateJsonResponse(true, "User invited"))
}

type TeamUserInvite struct {
	InviteId     int    `json:"inviteId"`
	UserId       int    `json:"userId"`
	UserNickname string `json:"userNickname"`
}

type DB_TeamInvite struct {
	Id       int
	UserId   int
	TeamId   int
	Accepted bool
	Declined bool
	Removed  bool
}

func getTeamUserInvites(teamId int) []TeamUserInvite {
	DB_conn, err := DB_connect()
	if err != nil {
		return make([]TeamUserInvite, 0)
	}
	defer DB_conn.Close()

	rows, err := DB_conn.Query("SELECT * FROM `team_invites` WHERE `team_id` = ? AND accepted = 0 AND declined = 0 AND removed = 0", teamId)
	if err != nil {
		return make([]TeamUserInvite, 0)
	}

	var invites []TeamUserInvite = make([]TeamUserInvite, 0)
	for rows.Next() {
		var inv DB_TeamInvite
		err := rows.Scan(&inv.Id, &inv.TeamId, &inv.UserId, &inv.Accepted, &inv.Declined, &inv.Removed)
		if err != nil {
			return make([]TeamUserInvite, 0)
		}

		t, err := getUserById(inv.UserId)
		if err != nil {
			return make([]TeamUserInvite, 0)
		}

		invites = append(invites, TeamUserInvite{
			InviteId:     inv.Id,
			UserId:       t.Id,
			UserNickname: t.Nickname,
		})
	}

	return invites
}

func team_get_invites(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req SendInviteRequest
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

	teamId := getUserTeamId(userId)
	if teamId == -1 {
		w.Write(CreateJsonResponse(false, "You are not member of any team"))
		return
	}

	invites := getTeamUserInvites(teamId)

	bytes, err = json.Marshal(invites)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Json error"))
		return
	}

	w.Write(CreateJsonResponse(true, string(bytes)))
}

type TeamRemoveInviteRequest struct {
	Token    string `json:"token"`
	InviteId int    `json:"inviteId"`
}

func team_remove_invite(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TeamRemoveInviteRequest
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
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	teamId := getUserTeamId(userId)
	if teamId == -1 {
		w.Write(CreateJsonResponse(false, "You are not member of any team"))
		return
	}

	db_team, err := getTeamById(teamId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Team do not exists"))
		return
	}

	if userId != db_team.CaptainId {
		w.Write(CreateJsonResponse(false, "Only captain can invite users to the team"))
		return
	}

	invites := getTeamUserInvites(teamId)

	for _, i := range invites {
		if i.InviteId == req.InviteId {
			_, err = DB_conn.Query("UPDATE `team_invites` SET removed = 1 WHERE id = ?", i.InviteId)
			if err != nil {
				w.Write(CreateJsonResponse(false, "Error while editing DB"))
				return
			}
			break
		}
	}

	w.Write(CreateJsonResponse(true, "User invite removed"))
}

type TeamRemoveMemberRequest struct {
	Token  string `json:"token"`
	UserId int    `json:"userId"`
}

func team_remove_member(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TeamRemoveMemberRequest
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
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanManageTeam {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot manage teams"))
		return
	}

	teamId := getUserTeamId(userId)
	if teamId == -1 {
		w.Write(CreateJsonResponse(false, "You are not member of any team"))
		return
	}

	db_team, err := getTeamById(teamId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Team do not exists"))
		return
	}

	if userId != db_team.CaptainId {
		w.Write(CreateJsonResponse(false, "Only captain can invite users to the team"))
		return
	}

	members := getTeamMembers(teamId)

	found := false
	for _, m := range members {
		if m.Id == req.UserId {
			_, err := DB_conn.Query("UPDATE `team_members` SET removed = 1 WHERE removed = 0 AND team_id = ? AND user_id = ?", teamId, req.UserId)
			if err != nil {
				w.Write(CreateJsonResponse(false, "DB error"))
				return
			}
			found = true
			break
		}
	}
	if !found {
		w.Write(CreateJsonResponse(false, "It is not member of your team"))
		return
	}

	w.Write(CreateJsonResponse(true, "User removed"))
}

type SubmitFlagRequest struct {
	TaskId int    `json:"taskId"`
	Flag   string `json:"flag"`
	Token  string `json:"token"`
}

func team_submit_flag(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req SubmitFlagRequest
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
		w.Write(CreateJsonResponse(false, "You are not logged int"))
		return
	}

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanViewTasks {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot view tasks"))
		return
	}

	teamId := getUserTeamId(userId)
	if teamId == -1 {
		w.Write(CreateJsonResponse(false, "You are not member of any team"))
		return
	}

	row := DB_conn.QueryRow("SELECT * FROM `tasks` WHERE `id` = ? AND visible = 1 AND removed = 0", req.TaskId)
	var db_task Task
	err = row.Scan(&db_task.Id, &db_task.Title, &db_task.Points, &db_task.Category, &db_task.Description, &db_task.Flag, &db_task.Visible, &db_task.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Task with give id not found in the DB"))
		return
	}

	if req.Flag != db_task.Flag {
		w.Write(CreateJsonResponse(false, "Wrong flag"))
		return
	}

	time := time.Now().Unix()
	_, err = DB_conn.Query("INSERT INTO `solved_tasks` (id, task_id, team_id, solved_at) VALUES (NULL, ?, ?, ?)", req.TaskId, teamId, time)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while changing DB data (1)"))
		return
	}
	_, err = DB_conn.Query("UPDATE `teams` SET score = score + ? WHERE id = ?", db_task.Points, teamId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while changing DB data (2)"))
		return
	}

	w.Write(CreateJsonResponse(true, "Flag accepted! Good job!"))
}

type TeamGetRequest struct {
	Token    string `json:"token"`
	TeamName string `json:"teamName"`
}

func team_get(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TeamGetRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "You are not logged in"))
		return
	}

	var team Team

	if req.TeamName == "" {
		user, err := getUserById(userId)
		if err != nil {
			w.Write(CreateJsonResponse(false, "ERROR"))
			return
		}

		teamId := getUserTeamId(user.Id)
		if teamId == -1 {
			w.Write(CreateJsonResponse(false, "You are not member of any team"))
			return
		}

		team, err = getTeamById(teamId)
		if err != nil {
			w.Write(CreateJsonResponse(false, "ERROR"))
			return
		}
	} else {
		team, err = getTeamByName(req.TeamName)
		if err != nil {
			w.Write(CreateJsonResponse(false, "ERROR"))
			return
		}
	}

	score, err := getTeamScore(team.Id)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}
	team.Score = score

	bytes, err = json.Marshal(team)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Json error"))
		return
	}

	w.Write(CreateJsonResponse(true, string(bytes)))
}

func team_get_members(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req TeamGetRequest
	err = json.Unmarshal(bytes[:n], &req)

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while parsing request object"))
		return
	}

	logged, userId := checkSession(req.Token)

	if !logged {
		w.Write(CreateJsonResponse(false, "You are not logged in"))
		return
	}

	var team Team

	if req.TeamName == "" {
		user, err := getUserById(userId)
		if err != nil {
			w.Write(CreateJsonResponse(false, "ERROR"))
			return
		}

		teamId := getUserTeamId(user.Id)
		if teamId == -1 {
			w.Write(CreateJsonResponse(false, "You are not member of any team"))
			return
		}

		team, err = getTeamById(teamId)
		if err != nil {
			w.Write(CreateJsonResponse(false, "ERROR"))
			return
		}
	} else {
		team, err = getTeamByName(req.TeamName)
		if err != nil {
			w.Write(CreateJsonResponse(false, "ERROR"))
			return
		}
	}

	bytes, err = json.Marshal(getTeamMembers(team.Id))
	if err != nil {
		w.Write(CreateJsonResponse(false, "Json error"))
		return
	}

	w.Write(CreateJsonResponse(true, string(bytes)))
}

// @Maybe we can move it somewhere else
func team_get_leaderboard(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	_, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	DB_conn, err := DB_connect()
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB connection error"))
		return
	}
	defer DB_conn.Close()

	rows, err := DB_conn.Query("SELECT * FROM `teams` WHERE removed = 0 ORDER BY score DESC")
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB error"))
		return
	}

	var teams []Team
	for rows.Next() {
		var team Team
		err := rows.Scan(&team.Id, &team.Name, &team.CaptainId, &team.Score, &team.Removed)
		if err != nil {
			w.Write(CreateJsonResponse(false, "DB error (2)"))
			return
		}
		teams = append(teams, team)
	}

	bytes, err = json.Marshal(teams)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Json error"))
		return
	}

	w.Write(CreateJsonResponse(true, string(bytes)))
}
