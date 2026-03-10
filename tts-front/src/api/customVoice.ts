import { request } from '@/utils/request'

export type CustomVoiceStatus = 'pending' | 'success' | 'failed' | string

export interface CustomVoiceRequest {
  id: number
  userId: number
  name: string
  tone: string
  gender: string
  sampleText: string
  sampleUrls: string[]
  status: CustomVoiceStatus
  resultVoiceId?: number
  errorMsg?: string
  createdAt: string
  updatedAt: string
}

export interface CustomVoiceListResponse {
  list: CustomVoiceRequest[]
  total: number
}

export interface CreateCustomVoiceParams {
  name: string
  tone?: string
  gender?: string
  sampleText?: string
  sampleUrls: string[]
}

export const customVoiceApi = {
  createRequest(params: CreateCustomVoiceParams) {
    return request<{ id: number }>({
      url: '/voice/custom/request',
      method: 'POST',
      data: params,
    })
  },

  listMy() {
    return request<CustomVoiceListResponse>({
      url: '/voice/custom/list',
      method: 'GET',
    })
  },

  delete(id: number) {
    return request<void>({
      url: `/voice/custom/${id}`,
      method: 'DELETE',
    })
  },

  adminList() {
    return request<CustomVoiceListResponse>({
      url: '/admin/voice/custom/list',
      method: 'GET',
    })
  },

  adminApprove(id: number) {
    return request<{ voice: any }>({
      url: `/admin/voice/custom/${id}/approve`,
      method: 'PUT',
    })
  },

  adminReject(id: number, errorMsg: string) {
    return request<void>({
      url: `/admin/voice/custom/${id}/reject`,
      method: 'PUT',
      data: { errorMsg },
    })
  },

  adminDelete(id: number) {
    return request<void>({
      url: `/admin/voice/custom/${id}`,
      method: 'DELETE',
    })
  },
}


