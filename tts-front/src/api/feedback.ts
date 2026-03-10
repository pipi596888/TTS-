import { request } from '@/utils/request'

export type FeedbackStatus = 'open' | 'closed'

export interface FeedbackItem {
  id: number
  userId: number
  username?: string
  category: string
  content: string
  contact: string
  status: FeedbackStatus | string
  reply?: string
  createdAt: string
  updatedAt: string
}

export interface FeedbackListResponse {
  list: FeedbackItem[]
  total: number
}

export interface CreateFeedbackParams {
  category?: string
  content: string
  contact?: string
}

export const feedbackApi = {
  create(params: CreateFeedbackParams) {
    return request<{ id: number }>({
      url: '/feedback',
      method: 'POST',
      data: params,
    })
  },

  listMy() {
    return request<FeedbackListResponse>({
      url: '/feedback/my',
      method: 'GET',
    })
  },

  adminList() {
    return request<FeedbackListResponse>({
      url: '/admin/feedback/list',
      method: 'GET',
    })
  },

  adminReply(id: number, params: { reply: string; status?: string }) {
    return request<void>({
      url: `/admin/feedback/${id}/reply`,
      method: 'PUT',
      data: params,
    })
  },
}


