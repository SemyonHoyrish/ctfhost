/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

var DB_NAME string
var DB_USER string
var DB_PASSWORD string
var ADMIN_PASSWORD string

func DB_connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@/"+DB_NAME)
	if err != nil {
		return nil, err
	}

	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)

	// If client make to many request, we cannot open more db connections,
	// so for now just increase their number and make their lifetime less.

	db.SetConnMaxLifetime(time.Second * 10)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	return db, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		DB_NAME = "ctfhost"
		DB_USER = "root"
		DB_PASSWORD = "pass"
		ADMIN_PASSWORD = "pass"
	} else {
		DB_NAME = os.Getenv("DB_NAME")
		DB_USER = os.Getenv("DB_USER")
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
		ADMIN_PASSWORD = os.Getenv("ADMIN_PASSWORD")
	}

	{
		db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@/")
		if err != nil {
			panic(err)
		}

		_, err = db.Query("CREATE DATABASE `" + DB_NAME + "` CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")
		if err == nil {
			fmt.Println("CREATED DATABASE " + DB_NAME)
		}

		db.Close()
	}

	db, err := DB_connect()
	if err != nil {
		panic(err)
	}

	_, err = db.Query("CREATE TABLE `ctfhost`.`settings` (`id` INT NOT NULL AUTO_INCREMENT , `user_can_register` BOOLEAN NOT NULL , `user_can_manage_team` BOOLEAN NOT NULL , `user_can_manage_account` BOOLEAN NOT NULL , `user_can_view_tasks` BOOLEAN NOT NULL , `max_members_per_team` INT NOT NULL , `main_page_content` MEDIUMTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL , `changed_by` INT NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")
	if err == nil {
		fmt.Println("Created table 'settings'")
		_, err = db.Query("INSERT INTO `settings` (id, user_can_register, user_can_manage_team, user_can_manage_account, user_can_view_tasks, max_members_per_team, main_page_content, changed_by) VALUES (NULL, 1, 1, 1, 1, -1, '<h1>Welcome to the CTF!</h1>', 0)")
		if err == nil {
			fmt.Println("Inserted default settings in 'settings'")
		}
	}
	_, err = db.Query("CREATE TABLE `ctfhost`.`users` (`id` INT NOT NULL AUTO_INCREMENT , `nickname` VARCHAR(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL , `email` VARCHAR(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL , `password_hash` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL , `is_admin` BOOLEAN NOT NULL DEFAULT FALSE , `removed` BOOLEAN NOT NULL DEFAULT FALSE , PRIMARY KEY (`id`), UNIQUE (`nickname`), UNIQUE (`email`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")
	if err == nil {
		fmt.Println("Created table 'users'")
		_, err = db.Query("INSERT INTO `users` (id, nickname, email, password_hash, is_admin, removed) VALUES (NULL, 'admin', 'admin@localhost', ?, 1, 0)", getPasswordHash(ADMIN_PASSWORD))
		if err == nil {
			fmt.Println("Created admin account")
		}
	}
	_, err = db.Query("CREATE TABLE `ctfhost`.`sessions` (`id` INT NOT NULL AUTO_INCREMENT , `user_id` INT NOT NULL , `token` VARCHAR(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL , `expires_on` BIGINT NOT NULL , `cancelled` BOOLEAN NOT NULL DEFAULT FALSE , PRIMARY KEY (`id`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")
	if err == nil {
		fmt.Println("Created table 'sessions'")
	}
	_, err = db.Query("CREATE TABLE `ctfhost`.`teams` (`id` INT NOT NULL AUTO_INCREMENT , `name` VARCHAR(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL , `captain_id` INT NOT NULL , `score` INT NOT NULL DEFAULT 0 , `removed` BOOLEAN NOT NULL DEFAULT FALSE , PRIMARY KEY (`id`) ) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")
	if err == nil {
		fmt.Println("Created table 'teams'")
	}
	_, err = db.Query("CREATE TABLE `ctfhost`.`team_members` (`id` INT NOT NULL AUTO_INCREMENT , `team_id` INT NOT NULL , `user_id` INT NOT NULL , `removed` BOOLEAN NOT NULL DEFAULT FALSE , PRIMARY KEY (`id`)) ENGINE = InnoDB;")
	if err == nil {
		fmt.Println("Created table 'team_members'")
	}
	_, err = db.Query("CREATE TABLE `ctfhost`.`team_invites` (`id` INT NOT NULL AUTO_INCREMENT , `team_id` INT NOT NULL , `user_id` INT NOT NULL , `accepted` BOOLEAN NOT NULL DEFAULT FALSE , `declined` BOOLEAN NOT NULL DEFAULT FALSE , `removed` BOOLEAN NOT NULL DEFAULT FALSE , PRIMARY KEY (`id`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")
	if err == nil {
		fmt.Println("Created table 'team_invites'")
	}
	_, err = db.Query("CREATE TABLE `ctfhost`.`tasks` (`id` INT NOT NULL AUTO_INCREMENT , `title` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL , `points` INT NOT NULL , `category` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL , `description` JSON NOT NULL , `flag` VARCHAR(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL , `visible` tinyint(1) NOT NULL DEFAULT '1', `removed` tinyint(1) NOT NULL DEFAULT '0', PRIMARY KEY (`id`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")
	if err == nil {
		fmt.Println("Created table 'tasks'")
	}
	_, err = db.Query("CREATE TABLE `ctfhost`.`solved_tasks` (`id` INT NOT NULL AUTO_INCREMENT , `task_id` INT NOT NULL , `team_id` INT NOT NULL , `solved_at` BIGINT NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")
	if err == nil {
		fmt.Println("Created table 'solved_tasks'")
	}

	db.Close()

	http.HandleFunc("/settings/get", settings_get)
	http.HandleFunc("/settings/set", settings_set)
	http.HandleFunc("/settings/get_main_page_content", settings_get_main_page_content)

	http.HandleFunc("/user/register", user_register)
	http.HandleFunc("/user/login", user_login)
	http.HandleFunc("/user/remove", user_remove)
	http.HandleFunc("/user/edit", user_edit)
	http.HandleFunc("/user/get", user_get)
	http.HandleFunc("/user/get_invites", user_get_invites)
	http.HandleFunc("/user/accept_invite", user_accept_invite)
	http.HandleFunc("/user/decline_invite", user_decline_invite)
	http.HandleFunc("/user/leave_team", user_leave_team)
	http.HandleFunc("/user/logout", user_logout)
	http.HandleFunc("/user/logout_all", user_logout_all)
	http.HandleFunc("/user/get_sessions_count", user_get_sessions_count)

	http.HandleFunc("/team/create", team_create)
	http.HandleFunc("/team/edit_name", team_edit_name)
	http.HandleFunc("/team/remove", team_remove)
	http.HandleFunc("/team/submit_flag", team_submit_flag)
	http.HandleFunc("/team/get", team_get)
	http.HandleFunc("/team/get_members", team_get_members)
	http.HandleFunc("/team/get_invites", team_get_invites)
	http.HandleFunc("/team/send_invite", team_send_invite)
	http.HandleFunc("/team/remove_invite", team_remove_invite)
	http.HandleFunc("/team/remove_member", team_remove_member)

	http.HandleFunc("/team/get_leaderboard", team_get_leaderboard)

	http.HandleFunc("/tasks/new", tasks_new)
	http.HandleFunc("/tasks/edit", tasks_edit)
	http.HandleFunc("/tasks/remove", tasks_remove)
	http.HandleFunc("/tasks/get_all", tasks_get_all)
	http.HandleFunc("/tasks/get", tasks_get)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("ERROR: '%s'", err.Error())
	}
}
