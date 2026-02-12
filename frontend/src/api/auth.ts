import type { LoginRequest, LoginResponse, RegisterRequest, User } from '@/types';
import { api } from '@/api/client';

export function login(data: LoginRequest) {
  return api.post<LoginResponse>('/auth/login', data);
}

export function register(data: RegisterRequest) {
  return api.post<User>('/auth/register', data);
}
