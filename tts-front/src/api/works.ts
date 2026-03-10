import { request } from '@/utils/request'

export interface Work {
  taskId: string
  title?: string
  status: string
  progress: number
  format: string
  audioUrl?: string
  errorMsg?: string
  createdAt: string
}

export interface WorksResponse {
  list: Work[]
  total: number
}

export const worksApi = {
  list() {
    return request<WorksResponse>({
      url: '/works/list',
      method: 'GET',
    })
  },
  updateTitle(taskId: string, title: string) {
    return request<void>({
      url: `/works/${taskId}/title`,
      method: 'PUT',
      data: { title },
    })
  },
  delete(taskId: string) {
    return request({
      url: `/works/${taskId}`,
      method: 'DELETE',
    })
  },
}

