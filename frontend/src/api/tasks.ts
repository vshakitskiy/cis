import type { Task, TaskStatus, TaskPriority } from '@/types';
import { api } from '@/api/client';

export function listTasks(projectId: number) {
  return api.get<Task[]>(`/projects/${projectId}/tasks`);
}

export function getTask(projectId: number, taskId: number) {
  return api.get<Task>(`/projects/${projectId}/tasks/${taskId}`);
}

export function createTask(
  projectId: number,
  data: {
    title: string;
    description?: string;
    assignee_id?: number;
    priority?: TaskPriority;
    deadline?: string;
  },
) {
  return api.post<Task>(`/projects/${projectId}/tasks`, data);
}

export function updateTask(
  projectId: number,
  taskId: number,
  data: {
    title: string;
    description?: string;
    assignee_id?: number | null;
    status: TaskStatus;
    priority: TaskPriority;
    deadline?: string;
  },
) {
  return api.put<Task>(`/projects/${projectId}/tasks/${taskId}`, data);
}

export function deleteTask(projectId: number, taskId: number) {
  return api.delete<void>(`/projects/${projectId}/tasks/${taskId}`);
}
