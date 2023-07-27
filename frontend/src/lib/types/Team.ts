/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

export default class Team {
	id: number;
	name: string;
	captainId: number;
	score: number;
	removed: boolean;

	constructor(id: number, name: string, captainId: number, score: number, removed: boolean) {
		this.id = id;
		this.name = name;
		this.captainId = captainId;
		this.score = score;
		this.removed = removed;
	}
}
