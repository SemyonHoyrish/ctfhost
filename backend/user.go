/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Id           int
	Nickname     string
	Email        string
	PasswordHash string
	IsAdmin      bool
	Removed      bool
}

type UsersRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func getPasswordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	s := fmt.Sprintf("%x", sum)
	return s
}

func getUserById(userId int) (User, error) {
	DB_conn, err := DB_connect()
	if err != nil {
		return User{}, err
	}

	row := DB_conn.QueryRow("SELECT * FROM `users` WHERE `id` = ?", userId)
	var db_user User
	err = row.Scan(&db_user.Id, &db_user.Nickname, &db_user.Email, &db_user.PasswordHash, &db_user.IsAdmin, &db_user.Removed)
	if err != nil {
		return User{}, err
	}
	return db_user, nil
}
func getUserByEmail(userEmail string) (User, error) {
	DB_conn, err := DB_connect()
	if err != nil {
		return User{}, err
	}

	row := DB_conn.QueryRow("SELECT * FROM `users` WHERE `email` = ?", userEmail)
	var db_user User
	err = row.Scan(&db_user.Id, &db_user.Nickname, &db_user.Email, &db_user.PasswordHash, &db_user.IsAdmin, &db_user.Removed)
	if err != nil {
		return User{}, err
	}
	return db_user, nil
}
func getUserByNickname(userNickname string) (User, error) {
	DB_conn, err := DB_connect()
	if err != nil {
		return User{}, err
	}

	row := DB_conn.QueryRow("SELECT * FROM `users` WHERE `nickname` = ?", userNickname)
	var db_user User
	err = row.Scan(&db_user.Id, &db_user.Nickname, &db_user.Email, &db_user.PasswordHash, &db_user.IsAdmin, &db_user.Removed)
	if err != nil {
		return User{}, err
	}
	return db_user, nil
}

func user_register(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UsersRequest
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

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanRegister {
		w.Write(CreateJsonResponse(false, "According to current settings new users cannot be registered"))
		return
	}

	row := DB_conn.QueryRow("SELECT * FROM `users` WHERE `nickname` = ? OR `email` = ?", req.Nickname, req.Email)
	var db_user User
	err = row.Scan(&db_user.Id, &db_user.Nickname, &db_user.Email, &db_user.PasswordHash, &db_user.IsAdmin, &db_user.Removed)
	if err == nil {
		w.Write(CreateJsonResponse(false, "User with given nickname or email already exists"))
		return
	}

	// @ToDo: validate nickname email
	_, err = DB_conn.Query("INSERT INTO `users` (id, nickname, email, password_hash, removed) VALUES (NULL, ?, ?, ?, 0)", req.Nickname, req.Email, getPasswordHash(req.Password))

	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while adding to the DB"))
		return
	}

	row = DB_conn.QueryRow("SELECT * FROM `users` WHERE `nickname` = ? AND `email` = ?", req.Nickname, req.Email)
	err = row.Scan(&db_user.Id, &db_user.Nickname, &db_user.Email, &db_user.PasswordHash, &db_user.IsAdmin, &db_user.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	session, err := newSession(db_user.Id)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Cannot create session"))
		return
	}

	w.Write(CreateJsonResponse(true, session))
}

func user_login(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UsersRequest
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

	row := DB_conn.QueryRow("SELECT * FROM `users` WHERE `email` = ?", req.Email)
	var db_user User
	err = row.Scan(&db_user.Id, &db_user.Nickname, &db_user.Email, &db_user.PasswordHash, &db_user.IsAdmin, &db_user.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "User not found in the DB"))
		return
	}

	if getPasswordHash(req.Password) != db_user.PasswordHash {
		w.Write(CreateJsonResponse(false, "Wrong password"))
		return
	}

	session, err := newSession(db_user.Id)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Cannot create session"))
		return
	}

	w.Write(CreateJsonResponse(true, session))
}

func user_logout(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UserEditRequest
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

	logged, _ := checkSession(req.Token)
	if !logged {
		w.Write(CreateJsonResponse(false, "You are not logged in"))
		return
	}

	err = cancelSession(req.Token)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	w.Write(CreateJsonResponse(true, ""))
}

func user_logout_all(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UserEditRequest
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

	err = cancelAllSessions(userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	w.Write(CreateJsonResponse(true, ""))
}

func user_get_sessions_count(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UserEditRequest
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

	time := time.Now().Unix()

	row := DB_conn.QueryRow("SELECT COUNT(id) FROM `sessions` WHERE user_id = ? AND cancelled = 0 AND expires_on > ?", userId, time)
	var count int
	err = row.Scan(&count)
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB error"))
		return
	}

	w.Write(CreateJsonResponse(true, fmt.Sprint(count)))
}

type UsersRemoveRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func user_remove(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UsersRemoveRequest
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

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanManageAccount {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot manage accounts"))
		return
	}

	var row *sql.Row
	email := req.Email

	if req.Email == "" {
		row = DB_conn.QueryRow("SELECT * FROM `users` WHERE `id` = ? AND removed = 0", userId)
		u, err := getUserById(userId)
		if err != nil {
			w.Write(CreateJsonResponse(false, "ERROR"))
			return
		}

		email = u.Email
	} else {
		row = DB_conn.QueryRow("SELECT * FROM `users` WHERE `email` = ? AND removed = 0", req.Email)
	}

	var db_user User
	err = row.Scan(&db_user.Id, &db_user.Nickname, &db_user.Email, &db_user.PasswordHash, &db_user.IsAdmin, &db_user.Removed)
	if err != nil {
		w.Write(CreateJsonResponse(false, "User not found in the DB"))
		return
	}

	if userId != db_user.Id {
		row := DB_conn.QueryRow("SELECT * FROM `users` WHERE `email` = ? AND removed = 0", email)
		var db_req_user User
		err = row.Scan(&db_req_user.Id, &db_req_user.Nickname, &db_req_user.Email, &db_req_user.PasswordHash, &db_req_user.IsAdmin, &db_req_user.Removed)
		if err != nil {
			w.Write(CreateJsonResponse(false, "User not found in the DB"))
			return
		}
		if !db_req_user.IsAdmin {
			w.Write(CreateJsonResponse(false, "You cannot delete other users"))
			return
		}
	}

	func() { // remove team
		var row *sql.Row
		tId := getUserTeamId(userId)
		row = DB_conn.QueryRow("SELECT * FROM `teams` WHERE `id` = ? AND removed = 0", tId)

		var db_team Team
		err = row.Scan(&db_team.Id, &db_team.Name, &db_team.CaptainId, &db_team.Score, &db_team.Removed)
		if err != nil {
			return
		}

		if userId != db_team.CaptainId {
			return
		}

		_, err = DB_conn.Query("UPDATE `teams` SET removed = 1 WHERE id = ? AND removed = 0", db_team.Id)
		if err != nil {
			return
		}

		_, err = DB_conn.Query("UPDATE `team_members` SET removed = 1 WHERE team_id = ? AND removed = 0", db_team.Id)
		if err != nil {
			return
		}
	}()

	_, err = DB_conn.Query("UPDATE `users` SET removed = 1 WHERE email = ? AND removed = 0", email)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while editing DB (1)"))
		return
	}

	_, err = DB_conn.Query("UPDATE `sessions` SET cancelled = 1 WHERE user_id = ? AND cancelled = 0", db_user.Id)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while editing DB (2)"))
		return
	}

	w.Write(CreateJsonResponse(true, ""))
}

type UserEditRequest struct {
	Token string       `json:"token"`
	User  UsersRequest `json:"user"`
}

func user_edit(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UserEditRequest
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

	settings, err := getSettings()
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while checking settings"))
		return
	}
	if !settings.UserCanManageAccount {
		w.Write(CreateJsonResponse(false, "According to current settings users cannot manage accounts"))
		return
	}

	// @ToDo: admin can change other users info

	user, err := getUserById(userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	newUserData := req.User
	if newUserData.Nickname == "" {
		newUserData.Nickname = user.Nickname
	}
	if newUserData.Email == "" {
		newUserData.Email = user.Email
	}

	_, err = DB_conn.Query("UPDATE `users` SET nickname = ?, email = ? WHERE id = ?", newUserData.Nickname, newUserData.Email, userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Error while editing DB (1)"))
		return
	}

	if newUserData.Password != "" && getPasswordHash(newUserData.Password) != user.PasswordHash {
		_, err = DB_conn.Query("UPDATE `users` SET password_hash = ? WHERE id = ?", getPasswordHash(newUserData.Password), userId)
	}

	w.Write(CreateJsonResponse(true, ""))
}

type UserTeamInvite struct {
	InviteId int    `json:"inviteId"`
	TeamId   int    `json:"teamId"`
	TeamName string `json:"teamName"`
}

func getUserTeamInvites(userId int) []UserTeamInvite {
	DB_conn, err := DB_connect()
	if err != nil {
		return make([]UserTeamInvite, 0)
	}
	defer DB_conn.Close()

	rows, err := DB_conn.Query("SELECT * FROM `team_invites` WHERE `user_id` = ? AND accepted = 0 AND declined = 0 AND removed = 0", userId)
	if err != nil {
		return make([]UserTeamInvite, 0)
	}

	var invites []UserTeamInvite = make([]UserTeamInvite, 0)
	for rows.Next() {
		var inv DB_TeamInvite
		err := rows.Scan(&inv.Id, &inv.TeamId, &inv.UserId, &inv.Accepted, &inv.Declined, &inv.Removed)
		if err != nil {
			return make([]UserTeamInvite, 0)
		}

		t, err := getTeamById(inv.TeamId)
		if err != nil {
			return make([]UserTeamInvite, 0)
		}

		invites = append(invites, UserTeamInvite{
			InviteId: inv.Id,
			TeamId:   t.Id,
			TeamName: t.Name,
		})
	}

	return invites
}

type UserManageInvitesRequest struct {
	Token    string `json:"token"`
	InviteId int    `json:"inviteId"`
}

func user_get_invites(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UserManageInvitesRequest
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

	invites := getUserTeamInvites(userId)

	bytes, err = json.Marshal(invites)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Json error"))
		return
	}

	w.Write(CreateJsonResponse(true, string(bytes)))
}

func user_accept_invite(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UserManageInvitesRequest
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

	tId := getUserTeamId(userId)
	if tId != -1 {
		w.Write(CreateJsonResponse(false, "You already in a team"))
		return
	}

	if settings.MaxMembersPerTeam != -1 && len(getTeamMembers(tId)) >= settings.MaxMembersPerTeam {
		w.Write(CreateJsonResponse(false, "Team has reched limit of memers"))
		return
	}

	invites := getUserTeamInvites(userId)

	for _, i := range invites {
		if i.InviteId == req.InviteId {
			_, err = DB_conn.Query("UPDATE `team_invites` SET accepted = 1 WHERE id = ?", i.InviteId)
			if err != nil {
				w.Write(CreateJsonResponse(false, "Error while editing DB (1)"))
				return
			}
			_, err = DB_conn.Query("INSERT INTO `team_members` (id, team_id, user_id, removed) VALUES (NULL, ?, ?, 0)", i.TeamId, userId)
			if err != nil {
				w.Write(CreateJsonResponse(false, "Error while editing DB (2)"))
				return
			}
			break
		}
	}

	w.Write(CreateJsonResponse(true, "Accepted"))
}

func user_decline_invite(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UserManageInvitesRequest
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

	invites := getUserTeamInvites(userId)

	for _, i := range invites {
		if i.InviteId == req.InviteId {
			_, err = DB_conn.Query("UPDATE `team_invites` SET declined = 1 WHERE id = ?", i.InviteId)
			if err != nil {
				w.Write(CreateJsonResponse(false, "Error while editing DB"))
				return
			}
			break
		}
	}

	w.Write(CreateJsonResponse(true, "Declined"))
}

func user_leave_team(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UserEditRequest
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
		w.Write(CreateJsonResponse(false, "You are not in any team"))
		return
	}

	_, err = DB_conn.Query("UPDATE `team_members` SET removed = 1 WHERE removed = 0 AND team_id = ? AND user_id = ?", teamId, userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "DB error"))
		return
	}

	w.Write(CreateJsonResponse(true, "Declined"))
}

type UserGetRequest struct {
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
}

type ResponseUser struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
}

func user_get(w http.ResponseWriter, r *http.Request) {
	var bytes []byte = make([]byte, READ_BUFFER)
	n, err := r.Body.Read(bytes)

	if err != nil && err.Error() != "EOF" {
		w.Write(CreateJsonResponse(false, "Error while reading request's body"))
		return
	}

	var req UserGetRequest
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

	user, err := getUserById(userId)
	if err != nil {
		w.Write(CreateJsonResponse(false, "ERROR"))
		return
	}

	var res ResponseUser

	if req.Nickname != "" {
		requestedUser, err := getUserByNickname(req.Nickname)
		if err != nil {
			w.Write(CreateJsonResponse(false, "ERROR"))
			return
		}

		if user.IsAdmin || user.Id == requestedUser.Id {
			res = ResponseUser{
				Id:       requestedUser.Id,
				Nickname: requestedUser.Nickname,
				Email:    requestedUser.Email,
				IsAdmin:  requestedUser.IsAdmin,
			}

		} else {
			res = ResponseUser{
				Id:       requestedUser.Id,
				Nickname: requestedUser.Nickname,
			}
		}
	} else {
		res = ResponseUser{
			Id:       user.Id,
			Nickname: user.Nickname,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
		}
	}

	bytes, err = json.Marshal(res)
	if err != nil {
		w.Write(CreateJsonResponse(false, "Json error"))
		return
	}

	w.Write(CreateJsonResponse(true, string(bytes)))
}
