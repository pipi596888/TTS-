import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Voice } from '@/types/voice'
import { voiceApi } from '@/api/voice'

export const useVoiceStore = defineStore('voice', () => {
  const voiceList = ref<Voice[]>([])
  const loading = ref(false)
  const error = ref<string>('')

  async function fetchVoiceList() {
    loading.value = true
    error.value = ''
    try {
      const res = await voiceApi.list()
      voiceList.value = res?.list || []
    } catch (err: any) {
      console.error('Failed to fetch voice list:', err)
      error.value = err.message || '获取音色列表失败'
    } finally {
      loading.value = false
    }
  }

  async function createVoice(params: { name: string; tone: string; gender: string }) {
    const res = await voiceApi.create(params)
    voiceList.value.push(res)
    return res
  }

  async function deleteVoice(id: number) {
    await voiceApi.delete(id)
    voiceList.value = voiceList.value.filter((v) => v.id !== id)
  }

  async function setDefaultVoice(id: number) {
    await voiceApi.setDefault(id)
    voiceList.value = voiceList.value.map((v) => ({
      ...v,
      isDefault: v.id === id,
    }))
  }

  function getDefaultVoice(): Voice | undefined {
    return voiceList.value.find((v) => v.isDefault)
  }

  function getVoiceById(id: number): Voice | undefined {
    return voiceList.value.find((v) => v.id === id)
  }

  return {
    voiceList,
    loading,
    error,
    fetchVoiceList,
    createVoice,
    deleteVoice,
    setDefaultVoice,
    getDefaultVoice,
    getVoiceById,
  }
})

