import type { Comment } from '@/types';
import { api } from '@/api/client';

function base(projectId: number, taskId: number) {
  return `/projects/${projectId}/tasks/${taskId}/comments`;
}

export function listComments(projectId: number, taskId: number) {
  return api.get<Comment[]>(base(projectId, taskId));
}

export function createComment(projectId: number, taskId: number, data: { content: string }) {
  return api.post<Comment>(base(projectId, taskId), data);
}

export function updateComment(projectId: number, taskId: number, id: number, data: { content: string }) {
  return api.put<Comment>(`${base(projectId, taskId)}/${id}`, data);
}

export function deleteComment(projectId: number, taskId: number, id: number) {
  return api.delete<void>(`${base(projectId, taskId)}/${id}`);
}
