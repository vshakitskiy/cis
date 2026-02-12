import type { TimeEntry } from '@/types';
import { api } from '@/api/client';

function base(projectId: number, taskId: number) {
  return `/projects/${projectId}/tasks/${taskId}/time-entries`;
}

export function listTimeEntries(projectId: number, taskId: number) {
  return api.get<TimeEntry[]>(base(projectId, taskId));
}

export function createTimeEntry(
  projectId: number,
  taskId: number,
  data: { minutes: number; description?: string; date: string },
) {
  return api.post<TimeEntry>(base(projectId, taskId), data);
}

export function deleteTimeEntry(projectId: number, taskId: number, id: number) {
  return api.delete<void>(`${base(projectId, taskId)}/${id}`);
}
