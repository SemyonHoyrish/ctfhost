/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

package main

import "encoding/json"

const READ_BUFFER = 1024 * 16

type Response struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`
}

func CreateJsonResponse(success bool, data string) []byte {
	bytes, err := json.Marshal(Response{
		success,
		data,
	})

	if err == nil {
		return bytes
	}
	return make([]byte, 0)
}
