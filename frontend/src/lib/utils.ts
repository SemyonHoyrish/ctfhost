/* Copyright (C) 2023  Semyon Hoyrish */
/* For more details see the "LICENSE" file */

import Task from './types/Task';

export type response_t = {
	task: {
		id: number;
		title: string;
		points: number;
		category: string;
		description: string;
		flag: string;
		visible: boolean;
	};
	solved: boolean;
};
export const TaskGetResponseToTask = (item: response_t): Task => {
	const t = new Task(item.task.id, item.task.title, item.task.points, item.task.category);
	t.descriptionItems = JSON.parse(item.task.description);
	t.visible = item.task.visible;
	t.flag = item.task.flag;
	t.solved = item.solved;
	return t;
};
export const TaskGetResponseToTaskArray = (data: Array<response_t>): Array<Task> => {
	const tasks: Array<Task> = new Array<Task>();
	data.forEach((item: response_t) => {
		tasks.push(TaskGetResponseToTask(item));
	});
	return tasks;
};
export const TaskToBackendTask = (task: Task) => {
	return {
		id: task.id,
		title: task.title,
		points: task.points,
		category: task.category,
		description: JSON.stringify(task.descriptionItems),
		flag: task.flag,
		visible: task.visible,
		removed: task.removed
	};
};
