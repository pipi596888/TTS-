import { request } from '@/utils/request'

export interface SystemStats {
  users: number
  voices: number
  tasks: number
  tasksPending: number
  tasksProcessing: number
  tasksSuccess: number
  tasksFailed: number
  feedbackOpen: number
  customVoicePending: number
}

export const systemApi = {
  stats() {
    return request<SystemStats>({
      url: '/admin/system/stats',
      method: 'GET',
    })
  },
}


