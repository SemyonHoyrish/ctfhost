/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

import type IDescriptionItem from './IDescriptionItem';

export default class Task {
	id: number;
	title: string;
	points: number;
	category: string;
	flag = '';
	solved? = false;
	visible = true;
	removed = false;
	descriptionItems: Array<IDescriptionItem> = new Array<IDescriptionItem>();

	constructor(id: number, title: string, points: number, category: string) {
		this.id = id;
		this.title = title;
		this.points = points;
		this.category = category;
	}

	getDescriptionHTML(): string {
		let result = '';

		for (const item of this.descriptionItems) {
			result += item.getHTML();
		}

		return result;
	}
}
