import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Segment, AudioFormat, AudioChannel, TTSTaskResponse } from '@/types/tts'
import { ttsApi } from '@/api/tts'

export const useTTSStore = defineStore('tts', () => {
  const segments = ref<Segment[]>([])
  const selectedVoiceId = ref<number>(0)
  const format = ref<AudioFormat>('mp3')
  const channel = ref<AudioChannel>('mono')
  const currentTaskId = ref<string>('')
  const taskStatus = ref<TTSTaskResponse | null>(null)
  const isGenerating = ref(false)
  const pollingTimer = ref<number | null>(null)

  const totalCharacters = computed(() => {
    return segments.value.reduce((sum, seg) => sum + seg.text.length, 0)
  })

  const isPolling = computed(() => pollingTimer.value !== null)

  const isTaskRunning = computed(() => {
    const s = taskStatus.value?.status
    return s === 'pending' || s === 'processing'
  })

  function insertSegmentAt(index: number, text: string, voiceId?: number): string {
    const id = `seg_${Date.now()}_${Math.random().toString(36).slice(2, 11)}`
    const insertIndex = Math.max(0, Math.min(index, segments.value.length))
    segments.value.splice(insertIndex, 0, {
      id,
      voiceId: voiceId ?? selectedVoiceId.value,
      text,
      order: insertIndex,
    })
    segments.value = segments.value.map((s, idx) => ({ ...s, order: idx }))
    return id
  }

  function addSegment(text: string, voiceId?: number): string {
    return insertSegmentAt(segments.value.length, text, voiceId)
  }

  function updateSegment(id: string, updates: Partial<Segment>) {
    const index = segments.value.findIndex((s) => s.id === id)
    if (index === -1) return

    const existing = segments.value[index]!
    segments.value[index] = {
      id: existing.id,
      voiceId: updates.voiceId ?? existing.voiceId,
      emotion: updates.emotion ?? existing.emotion,
      text: updates.text ?? existing.text,
      order: existing.order,
    }
  }

  function removeSegment(id: string) {
    segments.value = segments.value
      .filter((s) => s.id !== id)
      .map((s, idx) => ({
        ...s,
        order: idx,
      }))
  }

  function reorderSegments(newOrder: Segment[]) {
    segments.value = newOrder.map((seg, index) => ({
      ...seg,
      order: index,
    }))
  }

  function clearSegments() {
    segments.value = []
  }

  function setAllSegmentsVoice(voiceId: number) {
    segments.value = segments.value.map((s) => ({ ...s, voiceId }))
  }

  async function generateAudio(overrideSegments?: Segment[]): Promise<string> {
    if (isTaskRunning.value) {
      throw new Error('当前已有任务在生成中，请等待完成或重置任务')
    }

    const source = overrideSegments ?? segments.value
    const validSegments = source.filter((s) => s.text.trim())
    if (validSegments.length === 0) {
      throw new Error('请输入文本内容')
    }

    isGenerating.value = true
    try {
      const res = await ttsApi.generate({
        segments: validSegments.map((s) => ({
          voiceId: s.voiceId,
          emotion: s.emotion || 'neutral',
          text: s.text,
        })),
        format: format.value,
        channel: channel.value,
      })

      currentTaskId.value = res.taskId
      // 先给一个可见的初始状态，避免 UI 等待轮询才变化
      taskStatus.value = {
        taskId: res.taskId,
        status: 'pending',
        progress: 0,
      }
      startPolling(res.taskId)
      return res.taskId
    } finally {
      isGenerating.value = false
    }
  }

  function startPolling(taskId: string) {
    stopPolling()
    const pollOnce = async () => {
      try {
        const res = await ttsApi.getTask(taskId)
        taskStatus.value = res
        if (res.status === 'success' || res.status === 'failed') {
          stopPolling()
        }
      } catch (error) {
        console.error('Polling error:', error)
      }
    }

    void pollOnce()

    pollingTimer.value = window.setInterval(() => {
      void pollOnce()
    }, 2000)
  }

  function stopPolling() {
    if (pollingTimer.value) {
      clearInterval(pollingTimer.value)
      pollingTimer.value = null
    }
  }

  function resetTask() {
    currentTaskId.value = ''
    taskStatus.value = null
    stopPolling()
  }

  return {
    segments,
    selectedVoiceId,
    format,
    channel,
    currentTaskId,
    taskStatus,
    isGenerating,
    isPolling,
    isTaskRunning,
    totalCharacters,
    insertSegmentAt,
    addSegment,
    updateSegment,
    removeSegment,
    reorderSegments,
    clearSegments,
    setAllSegmentsVoice,
    generateAudio,
    startPolling,
    stopPolling,
    resetTask,
  }
})

