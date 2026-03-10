import { request } from '@/utils/request'
import type { TTSGenerateParams, TTSGenerateResponse, TTSTaskResponse } from '@/types/tts'

export interface TTSTaskDetailResponse {
  taskId: string
  title?: string
  status: string
  progress: number
  audioUrl?: string
  error?: string
  format: string
  channel: string
  segments: {
    voiceId: number
    emotion?: string
    text: string
    sort: number
  }[]
  createdAt?: string
  updatedAt?: string
}

export const ttsApi = {
  generate(params: TTSGenerateParams) {
    return request<TTSGenerateResponse>({
      url: '/tts/generate',
      method: 'POST',
      data: params,
    })
  },

  getTask(taskId: string) {
    return request<TTSTaskResponse>({
      url: `/tts/task/${taskId}`,
      method: 'GET',
    })
  },

  getTaskDetail(taskId: string) {
    return request<TTSTaskDetailResponse>({
      url: `/tts/task/${taskId}/detail`,
      method: 'GET',
    })
  },
}
