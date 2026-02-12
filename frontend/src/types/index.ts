export type UserRole = 'admin' | 'manager' | 'employee';
export type TaskStatus = 'todo' | 'in_progress' | 'done';
export type TaskPriority = 'low' | 'medium' | 'high' | 'critical';

export interface User {
  id: number;
  email: string;
  name: string;
  role: UserRole;
  created_at: string;
  updated_at: string;
}

export interface Project {
  id: number;
  name: string;
  description: string | null;
  owner_id: number;
  deadline: string | null;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: number;
  project_id: number;
  assignee_id: number | null;
  title: string;
  description: string | null;
  status: TaskStatus;
  priority: TaskPriority;
  deadline: string | null;
  created_at: string;
  updated_at: string;
}

export interface TimeEntry {
  id: number;
  task_id: number;
  user_id: number;
  minutes: number;
  description: string | null;
  date: string;
  created_at: string;
}

export interface Comment {
  id: number;
  task_id: number;
  user_id: number;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface Attachment {
  id: number;
  task_id: number;
  uploaded_by: number;
  filename: string;
  size: number;
  created_at: string;
}

export interface ProjectReport {
  project_id: number;
  project_name: string;
  total_tasks: number;
  completed_tasks: number;
  completion_pct: number;
  total_minutes: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
  role?: UserRole;
}

export interface LoginResponse {
  token: string;
}

export interface ApiError {
  error: string;
}
