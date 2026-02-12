import type { ProjectReport } from '@/types';
import { api } from '@/api/client';

export function getProjectReport(projectId: number) {
  return api.get<ProjectReport>(`/projects/${projectId}/report`);
}
