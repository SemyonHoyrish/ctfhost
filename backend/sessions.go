/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

package main

import (
	"crypto/md5"
	"fmt"
	"time"
)

type Session struct {
	Id        int
	UserId    int
	Token     string
	ExpiresOn int64
	Cancelled bool
}

const SESSION_EXPIRES_IN = 60 * 60 * 24

func newSession(userId int) (string, error) {
	time := time.Now().Unix()
	token := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(userId)+fmt.Sprint(time))))
	expires := time + SESSION_EXPIRES_IN

	DB_conn, err := DB_connect()
	if err != nil {
		return "", err
	}
	defer DB_conn.Close()

	_, err = DB_conn.Query("INSERT INTO `sessions` (id, user_id, token, expires_on, cancelled) VALUES (NULL, ?, ?, ?, 0)", userId, token, expires)

	if err != nil {
		return "", err
	}

	return token, nil
}

func checkSession(token string) (bool, int) {
	DB_conn, err := DB_connect()
	if err != nil {
		return false, -1
	}
	defer DB_conn.Close()

	row := DB_conn.QueryRow("SELECT * FROM `sessions` WHERE token = ? AND cancelled = 0", token)
	var db_session Session
	err = row.Scan(&db_session.Id, &db_session.UserId, &db_session.Token, &db_session.ExpiresOn, &db_session.Cancelled)

	if err != nil {
		return false, -1
	}

	time := time.Now().Unix()

	if time > db_session.ExpiresOn {
		return false, -1
	}

	return true, db_session.UserId
}

func cancelSession(token string) error {
	DB_conn, err := DB_connect()
	if err != nil {
		return err
	}
	defer DB_conn.Close()

	row := DB_conn.QueryRow("SELECT * FROM `sessions` WHERE cancelled = 0 AND token = ?", token)
	var db_session Session
	err = row.Scan(&db_session.Id, &db_session.UserId, &db_session.Token, &db_session.ExpiresOn, &db_session.Cancelled)

	if err != nil {
		return err
	}

	_, err = DB_conn.Query("UPDATE `sessions` SET cancelled = 1 WHERE cancelled = 0 AND token = ?", token)

	if err != nil {
		return err
	}

	return nil
}

func cancelAllSessions(userId int) error {
	DB_conn, err := DB_connect()
	if err != nil {
		return err
	}
	defer DB_conn.Close()

	_, err = DB_conn.Query("SELECT * FROM `sessions` WHERE cancelled = 0 AND user_id = ?", userId)
	if err != nil {
		return err
	}

	_, err = DB_conn.Query("UPDATE `sessions` SET cancelled = 1 WHERE cancelled = 0 AND user_id = ?", userId)

	if err != nil {
		return err
	}

	return nil
}
