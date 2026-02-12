import type { Project } from '@/types';
import { api } from '@/api/client';

export function listProjects() {
  return api.get<Project[]>('/projects');
}

export function getProject(id: number) {
  return api.get<Project>(`/projects/${id}`);
}

export function createProject(data: { name: string; description?: string; deadline?: string }) {
  return api.post<Project>('/projects', data);
}

export function updateProject(id: number, data: { name: string; description?: string; deadline?: string }) {
  return api.put<Project>(`/projects/${id}`, data);
}

export function deleteProject(id: number) {
  return api.delete<void>(`/projects/${id}`);
}
