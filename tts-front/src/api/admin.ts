import { request } from '@/utils/request'

export type AdminUserRoleKey = 'admin' | 'engineer' | 'user'
export type AdminUserStatusKey = 'active' | 'disabled'

export interface AdminUser {
  id: number
  uid: string
  username: string
  email?: string
  role: AdminUserRoleKey
  status: AdminUserStatusKey
  createdAt?: string
  updatedAt?: string
}

export interface AdminUsersResponse {
  list: AdminUser[]
  total: number
}

export interface AdminRole {
  key: AdminUserRoleKey
  name: string
  description: string
  permissions: string[]
}

export interface AdminRolesResponse {
  list: AdminRole[]
}

export interface AdminLog {
  id: number
  actorUserId: number
  actorUsername: string
  action: string
  ip: string
  createdAt: string
}

export interface AdminLogsResponse {
  list: AdminLog[]
  total: number
}

export const adminApi = {
  listUsers(params: { keyword?: string; page?: number; pageSize?: number }) {
    return request<AdminUsersResponse>({
      url: '/admin/users',
      method: 'GET',
      params,
    })
  },

  createUser(params: { username: string; password: string; email?: string; role: AdminUserRoleKey; status: AdminUserStatusKey }) {
    return request<{ id: number }>({
      url: '/admin/users',
      method: 'POST',
      data: params,
    })
  },

  updateUser(id: number, params: { username?: string; password?: string; email?: string; role?: AdminUserRoleKey; status?: AdminUserStatusKey }) {
    return request<{ ok: boolean }>({
      url: `/admin/users/${id}`,
      method: 'PUT',
      data: params,
    })
  },

  deleteUser(id: number) {
    return request<{ ok: boolean }>({
      url: `/admin/users/${id}`,
      method: 'DELETE',
    })
  },

  listRoles() {
    return request<AdminRolesResponse>({
      url: '/admin/roles',
      method: 'GET',
    })
  },

  listLogs(params: { keyword?: string; page?: number; pageSize?: number }) {
    return request<AdminLogsResponse>({
      url: '/admin/logs',
      method: 'GET',
      params,
    })
  },

  appendLog(action: string) {
    return request<{ ok: boolean }>({
      url: '/admin/logs',
      method: 'POST',
      data: { action },
    })
  },
}
