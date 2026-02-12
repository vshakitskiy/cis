import type { Attachment } from '@/types';
import { api } from '@/api/client';

function base(projectId: number, taskId: number) {
  return `/projects/${projectId}/tasks/${taskId}/attachments`;
}

export function listAttachments(projectId: number, taskId: number) {
  return api.get<Attachment[]>(base(projectId, taskId));
}

export function uploadAttachment(projectId: number, taskId: number, file: File) {
  return api.upload<Attachment>(base(projectId, taskId), file);
}

export function deleteAttachment(projectId: number, taskId: number, id: number) {
  return api.delete<void>(`${base(projectId, taskId)}/${id}`);
}

export function downloadAttachment(projectId: number, taskId: number, id: number) {
  return api.download(`${base(projectId, taskId)}/${id}/download`);
}
