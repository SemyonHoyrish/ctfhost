/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

export default class User {
	id: number;
	nickname: string;
	email: string;
	isAdmin: boolean;

	constructor(id: number, nickname: string, email: string, isAdmin: boolean) {
		this.id = id;
		this.nickname = nickname;
		this.email = email;
		this.isAdmin = isAdmin;
	}
}
