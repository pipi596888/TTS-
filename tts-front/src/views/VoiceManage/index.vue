<template>
  <div class="voice-page">
    <div class="page-head">
      <div class="head-left">
        <div class="title">音色管理</div>
        <div class="sub">管理可用音色与默认音色</div>
      </div>
      <div class="head-right">
        <el-input v-model="keyword" clearable placeholder="搜索音色名称..." class="search" @clear="page = 1">
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="openCreate">
          <el-icon><Plus /></el-icon>
          新增音色
        </el-button>
      </div>
    </div>

    <div class="tabs-row">
      <el-tabs v-model="tab" class="filter-tabs">
        <el-tab-pane label="全部音色" name="all" />
        <el-tab-pane label="默认音色" name="default" />
        <el-tab-pane label="男声" name="male" />
        <el-tab-pane label="女声" name="female" />
      </el-tabs>
      <div class="tabs-right">
        <el-button :loading="loading" text @click="voiceStore.fetchVoiceList()">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <div class="grid-wrap" v-loading="loading">
      <div v-if="!loading && pagedVoices.length === 0" class="empty">
        <el-empty description="暂无音色" :image-size="120" />
      </div>

      <div v-else class="grid">
        <div v-for="v in pagedVoices" :key="v.id" class="voice-card">
          <div class="card-top">
            <div class="cover" :class="{ def: v.isDefault }">
              <el-icon class="mic"><Microphone /></el-icon>
            </div>
            <div class="status">
              <span v-if="v.isDefault" class="pill ok">默认</span>
              <span v-else class="pill" :class="genderClass(v.gender)">{{ genderText(v.gender) }}</span>
            </div>
          </div>

          <div class="card-main">
            <div class="name" :title="v.name">{{ v.name }}</div>
            <div class="meta">音色: {{ v.tone || '-' }}</div>
            <div class="meta-row">
              <span class="meta-item">
                <el-icon><Key /></el-icon>
                ID {{ v.id }}
              </span>
              <span class="meta-item">
                <el-icon><Link /></el-icon>
                {{ v.previewUrl ? '可试听' : '无试听' }}
              </span>
            </div>
          </div>

          <div class="card-actions">
            <el-button class="act-btn" :disabled="!v.previewUrl" @click="preview(v)">
              <el-icon><VideoPlay /></el-icon>
              试听
            </el-button>
            <el-button class="act-btn" :disabled="v.isDefault" @click="setDefault(v)">
              <el-icon><Check /></el-icon>
              设为默认
            </el-button>
            <el-button class="act-btn" @click="copyId(v.id)">
              <el-icon><DocumentCopy /></el-icon>
              复制ID
            </el-button>
            <el-button class="act-btn danger" @click="handleDelete(v.id)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <div class="pager" v-if="!loading && filteredVoices.length > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :page-sizes="[6, 9, 12, 18]"
        layout="prev, pager, next, jumper"
        :total="filteredVoices.length"
      />
    </div>

    <div v-if="playerVisible" class="player-bar">
      <audio ref="audioEl" :src="playerUrl" preload="metadata" />

      <div class="pb-left">
        <div class="pb-icon">
          <el-icon><Microphone /></el-icon>
        </div>
        <div class="pb-info">
          <div class="pb-title">正在试听：{{ playingName }}</div>
          <div class="pb-sub">音色: {{ playingTone }}</div>
        </div>
      </div>

      <div class="pb-center">
        <div class="pb-controls">
          <el-button circle type="primary" @click="togglePlay">
            <el-icon><component :is="isPlaying ? VideoPause : VideoPlay" /></el-icon>
          </el-button>
        </div>

        <div class="pb-progress">
          <span class="pb-time">{{ formatTime(currentTime) }}</span>
          <el-slider
            v-model="currentTime"
            class="pb-slider"
            :min="0"
            :max="duration"
            :show-tooltip="false"
            @change="seek"
          />
          <span class="pb-time">{{ formatTime(duration) }}</span>
        </div>
      </div>

      <div class="pb-right">
        <el-icon class="vol"><component :is="'VolumeUp'" /></el-icon>
        <el-slider
          v-model="volume"
          class="vol-slider"
          :min="0"
          :max="1"
          :step="0.05"
          @change="setVolume"
        />
        <el-button text @click="closePlayer">关闭</el-button>
      </div>
    </div>

    <el-dialog v-model="createVisible" title="新增音色" width="420px">
      <el-form :model="createForm" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="createForm.name" placeholder="如：旁白 / 女声1" />
        </el-form-item>
        <el-form-item label="音色" required>
          <el-input v-model="createForm.tone" placeholder="如：轻柔 / 沉稳 / 活泼" />
        </el-form-item>
        <el-form-item label="性别" required>
          <el-select v-model="createForm.gender" placeholder="请选择">
            <el-option label="男" value="male" />
            <el-option label="女" value="female" />
            <el-option label="中性" value="neutral" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="submitCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Check,
  Delete,
  DocumentCopy,
  Key,
  Link,
  Microphone,
  Plus,
  Refresh,
  Search,
  VideoPause,
  VideoPlay,
} from '@element-plus/icons-vue'
import { storeToRefs } from 'pinia'
import { useVoiceStore } from '@/store/voice'
import type { Voice } from '@/types/voice'

const voiceStore = useVoiceStore()
const { voiceList, loading } = storeToRefs(voiceStore)

const tab = ref<'all' | 'default' | 'male' | 'female'>('all')
const keyword = ref('')

const page = ref(1)
const pageSize = ref(9)

const filteredVoices = computed(() => {
  const base = [...voiceList.value]
  const t = tab.value
  const byTab =
    t === 'all'
      ? base
      : t === 'default'
        ? base.filter((v) => v.isDefault)
        : base.filter((v) => normalizeGender(v.gender) === t)

  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return byTab
  return byTab.filter((v) => `${v.name} ${v.tone} ${v.id}`.toLowerCase().includes(kw))
})

const pagedVoices = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredVoices.value.slice(start, start + pageSize.value)
})

watch([tab, keyword], () => {
  page.value = 1
})

onMounted(() => {
  voiceStore.fetchVoiceList()
})

function normalizeGender(g: string) {
  const s = `${g || ''}`.toLowerCase()
  if (s.includes('男') || s === 'male' || s === 'm') return 'male'
  if (s.includes('女') || s === 'female' || s === 'f') return 'female'
  return 'neutral'
}

function genderText(g: string) {
  const s = normalizeGender(g)
  if (s === 'male') return '男声'
  if (s === 'female') return '女声'
  return '中性'
}

function genderClass(g: string) {
  const s = normalizeGender(g)
  if (s === 'male') return 'male'
  if (s === 'female') return 'female'
  return 'neutral'
}

async function writeToClipboard(text: string) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text)
    return
  }
  const el = document.createElement('textarea')
  el.value = text
  el.style.position = 'fixed'
  el.style.left = '-9999px'
  el.style.top = '-9999px'
  document.body.appendChild(el)
  el.focus()
  el.select()
  const ok = document.execCommand('copy')
  document.body.removeChild(el)
  if (!ok) throw new Error('copy_failed')
}

async function copyId(id: number) {
  try {
    await writeToClipboard(String(id))
    ElMessage.success('已复制音色ID')
  } catch {
    ElMessage.error('复制失败')
  }
}

async function setDefault(v: Voice) {
  try {
    await voiceStore.setDefaultVoice(v.id)
    ElMessage.success('默认音色已更新')
  } catch {
    ElMessage.error('设置默认音色失败')
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除这个音色吗？', '警告', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
    await voiceStore.deleteVoice(id)
    ElMessage.success('删除成功')
    if (playingVoiceId.value === id) closePlayer()
  } catch {
    // cancelled
  }
}

// Preview player
const audioEl = ref<HTMLAudioElement | null>(null)
const playerVisible = ref(false)
const isPlaying = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const volume = ref(1)
const playerUrl = ref<string>('')
const playingVoiceId = ref<number | null>(null)
const playingName = ref<string>('-')
const playingTone = ref<string>('-')

watch(volume, (v) => {
  if (audioEl.value) audioEl.value.volume = v
})

function bindAudioEvents(el: HTMLAudioElement) {
  el.addEventListener('timeupdate', () => {
    currentTime.value = el.currentTime || 0
  })
  el.addEventListener('loadedmetadata', () => {
    duration.value = Number.isFinite(el.duration) ? el.duration : 0
  })
  el.addEventListener('ended', () => {
    isPlaying.value = false
  })
  el.addEventListener('pause', () => {
    isPlaying.value = false
  })
  el.addEventListener('play', () => {
    isPlaying.value = true
  })
}

watch(audioEl, (el) => {
  if (el) {
    bindAudioEvents(el)
    el.volume = volume.value
  }
})

function preview(v: Voice) {
  if (!v.previewUrl) return
  playerVisible.value = true
  playerUrl.value = v.previewUrl
  playingVoiceId.value = v.id
  playingName.value = v.name || '-'
  playingTone.value = v.tone || '-'
  currentTime.value = 0
  duration.value = 0
  const el = audioEl.value
  if (!el) return
  el.currentTime = 0
  void el.play().catch((err) => {
    ElMessage.error(err?.message || '无法播放试听音频')
  })
}

function togglePlay() {
  const el = audioEl.value
  if (!el) return
  if (isPlaying.value) el.pause()
  else void el.play()
}

function seek(t: number) {
  const el = audioEl.value
  if (!el) return
  el.currentTime = t
}

function setVolume(v: number) {
  if (audioEl.value) audioEl.value.volume = v
}

function closePlayer() {
  const el = audioEl.value
  if (el) el.pause()
  playerVisible.value = false
  isPlaying.value = false
  playingVoiceId.value = null
  playerUrl.value = ''
  currentTime.value = 0
  duration.value = 0
}

function formatTime(seconds: number): string {
  const s = Math.max(0, Math.floor(seconds || 0))
  const m = Math.floor(s / 60)
  const r = s % 60
  return `${m.toString().padStart(2, '0')}:${r.toString().padStart(2, '0')}`
}

// Create dialog
const createVisible = ref(false)
const creating = ref(false)
const createForm = ref<{ name: string; tone: string; gender: string }>({
  name: '',
  tone: '',
  gender: 'female',
})

function openCreate() {
  createForm.value = { name: '', tone: '', gender: 'female' }
  createVisible.value = true
}

async function submitCreate() {
  const name = createForm.value.name.trim()
  const tone = createForm.value.tone.trim()
  const gender = createForm.value.gender
  if (!name || !tone || !gender) {
    ElMessage.warning('请填写完整信息')
    return
  }
  creating.value = true
  try {
    await voiceStore.createVoice({ name, tone, gender })
    ElMessage.success('创建成功')
    createVisible.value = false
  } catch (err: any) {
    ElMessage.error(err?.message || '创建失败')
  } finally {
    creating.value = false
  }
}

onUnmounted(() => {
  closePlayer()
})
</script>

<style scoped>
.voice-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f5f7fb;
  padding: 18px 18px 92px;
  overflow: hidden;
}

.page-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 4px 2px 12px;
}

.head-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.title {
  font-size: 22px;
  font-weight: 800;
  color: #1f1f1f;
}

.sub {
  font-size: 13px;
  color: #7b7b7b;
}

.head-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search {
  width: 320px;
}

.tabs-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 0 2px;
}

.tabs-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.grid-wrap {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 14px 2px 0;
}

.grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.voice-card {
  width: 320px;
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e8eef7;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.cover {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(64, 158, 255, 0.12);
  color: #409eff;
}

.cover.def {
  background: rgba(103, 194, 58, 0.12);
  color: #67c23a;
}

.mic {
  font-size: 18px;
}

.pill {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 999px;
  border: 1px solid transparent;
  user-select: none;
  color: #409eff;
  background: rgba(64, 158, 255, 0.12);
  border-color: rgba(64, 158, 255, 0.2);
}

.pill.ok {
  color: #67c23a;
  background: rgba(103, 194, 58, 0.12);
  border-color: rgba(103, 194, 58, 0.2);
}

.pill.male {
  color: #409eff;
  background: rgba(64, 158, 255, 0.12);
  border-color: rgba(64, 158, 255, 0.2);
}

.pill.female {
  color: #a855f7;
  background: rgba(168, 85, 247, 0.12);
  border-color: rgba(168, 85, 247, 0.2);
}

.pill.neutral {
  color: #909399;
  background: rgba(144, 147, 153, 0.12);
  border-color: rgba(144, 147, 153, 0.2);
}

.card-main {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.name {
  font-weight: 800;
  color: #1f1f1f;
  font-size: 14px;
  line-height: 1.25;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 36px;
}

.meta {
  font-size: 12px;
  color: #6b7280;
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 14px;
  font-size: 12px;
  color: #8c8c8c;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.card-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.act-btn {
  justify-content: center;
  height: 34px;
  border-radius: 10px;
}

.act-btn.danger {
  color: #f56c6c;
  border-color: rgba(245, 108, 108, 0.35);
  background: rgba(245, 108, 108, 0.06);
}

.empty {
  padding-top: 40px;
  display: flex;
  justify-content: center;
}

.pager {
  display: flex;
  justify-content: center;
  padding: 14px 0 0;
}

.player-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  height: 74px;
  background: #fff;
  border-top: 1px solid #eaeef5;
  box-shadow: 0 -6px 20px rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 18px;
  gap: 14px;
}

.pb-left {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 260px;
}

.pb-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: rgba(64, 158, 255, 0.12);
  color: #409eff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pb-info {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.pb-title {
  font-weight: 800;
  color: #1f1f1f;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 260px;
}

.pb-sub {
  font-size: 12px;
  color: #8c8c8c;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 260px;
}

.pb-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.pb-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pb-progress {
  width: 100%;
  max-width: 520px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.pb-time {
  font-size: 12px;
  color: #8c8c8c;
}

.pb-slider {
  flex: 1;
}

.pb-right {
  min-width: 260px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
}

.vol {
  color: #8c8c8c;
}

.vol-slider {
  width: 120px;
}

@media (max-width: 980px) {
  .voice-page {
    padding: 12px 12px 92px;
  }

  .page-head {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }

  .head-right {
    width: 100%;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .search {
    width: 100%;
  }

  .tabs-row {
    flex-direction: column;
    align-items: stretch;
  }

  .tabs-right {
    justify-content: flex-end;
  }

  .grid {
    justify-content: center;
  }

  .voice-card {
    width: 100%;
    max-width: 520px;
  }

  .pb-right {
    display: none;
  }
  .pb-left {
    min-width: 0;
  }
  .pb-title,
  .pb-sub {
    max-width: 100%;
  }
}
</style>


