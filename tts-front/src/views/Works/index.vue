<template>
  <div class="works-page">
    <div class="page-head">
      <div class="head-left">
        <div class="title">我的作品</div>
        <div class="sub">系统会自动保存您生成的所有音频项目</div>
      </div>
      <div class="head-right">
        <el-input v-model="keyword" clearable placeholder="搜索作品名称或任务ID..." class="search" @clear="page = 1">
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新建项目
        </el-button>
      </div>
    </div>

    <div class="tabs-row">
      <el-tabs v-model="statusFilter" class="filter-tabs">
        <el-tab-pane label="全部作品" name="all" />
        <el-tab-pane label="正在处理" name="processing" />
        <el-tab-pane label="已完成" name="completed" />
        <el-tab-pane label="失败" name="failed" />
      </el-tabs>

      <div class="tabs-right">
        <el-tag v-if="hasRunning" type="warning" effect="plain">自动刷新中</el-tag>
        <el-button :loading="loading" text @click="refreshList(false)">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <div class="grid-wrap" v-loading="loading">
      <div v-if="!loading && pagedWorks.length === 0" class="empty">
        <el-empty description="暂无作品" :image-size="120" />
      </div>

      <div v-else class="grid">
        <div v-for="w in pagedWorks" :key="w.taskId" class="work-card">
          <div class="card-top">
            <div class="cover" :class="coverClass(w)">
              <el-icon class="note"><Headset /></el-icon>
            </div>
            <div class="status">
              <span class="status-pill" :class="statusClass(w)">{{ statusText(w) }}</span>
            </div>
          </div>

          <div class="card-main">
            <div class="name" :title="workTitle(w)">{{ workTitle(w) }}</div>
            <div class="meta">格式: {{ (w.format || '-').toUpperCase() }}</div>
            <div class="meta-row">
              <span class="meta-item">
                <el-icon><Calendar /></el-icon>
                {{ formatDate(w.createdAt) }}
              </span>
              <span class="meta-item mono" :title="w.taskId">
                <el-icon><DocumentCopy /></el-icon>
                {{ shortTaskId(w.taskId) }}
                <el-button class="copy-btn" text size="small" :icon="DocumentCopy" @click="copyText(w.taskId)" />
              </span>
            </div>
            <div v-if="isFailed(w) && w.errorMsg" class="err" :title="w.errorMsg">失败原因: {{ w.errorMsg }}</div>
          </div>

          <div class="processing-box" v-if="isProcessing(w)">
            <div class="processing-tip">生成中，请稍候... ({{ progressPercent(w) }}%)</div>
            <el-progress :percentage="progressPercent(w)" :stroke-width="10" />
          </div>

          <div class="card-actions" v-else-if="isCompleted(w)">
            <el-button class="act-btn" @click="playWork(w)">
              <el-icon><VideoPlay /></el-icon>
              播放
            </el-button>
            <el-button class="act-btn" @click="handleDownload(w)">
              <el-icon><Download /></el-icon>
              下载
            </el-button>
            <el-button class="act-btn" @click="handleShare(w)">
              <el-icon><Share /></el-icon>
              分享
            </el-button>
            <el-button class="act-btn" @click="handleEdit(w)">
              <el-icon><EditPen /></el-icon>
              编辑
            </el-button>
            <el-button class="act-btn" @click="handleRename(w)">
              <el-icon><EditPen /></el-icon>
              重命名
            </el-button>
            <el-button class="act-btn danger" @click="handleDelete(w)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>

          <div class="card-actions" v-else>
            <el-button class="act-btn" @click="handleEdit(w)">
              <el-icon><EditPen /></el-icon>
              编辑重试
            </el-button>
            <el-button class="act-btn" @click="handleRename(w)">
              <el-icon><EditPen /></el-icon>
              重命名
            </el-button>
            <el-button class="act-btn danger" @click="handleDelete(w)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <div class="pager" v-if="!loading && filteredWorks.length > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :page-sizes="[6, 9, 12, 18]"
        layout="prev, pager, next, jumper"
        :total="filteredWorks.length"
      />
    </div>

    <div v-if="playerVisible" class="player-bar">
      <audio ref="audioEl" :src="playerUrl" preload="metadata" />

      <div class="pb-left">
        <div class="pb-icon">
          <el-icon><Headset /></el-icon>
        </div>
        <div class="pb-info">
          <div class="pb-title">正在播放: {{ playingTitle }}</div>
          <div class="pb-sub">{{ playingSub }}</div>
        </div>
      </div>

      <div class="pb-center">
        <div class="pb-controls">
          <el-button circle text :disabled="!hasPrev" @click="playPrev">
            <el-icon><ArrowLeft /></el-icon>
          </el-button>
          <el-button circle type="primary" @click="togglePlay">
            <el-icon><component :is="isPlaying ? VideoPause : VideoPlay" /></el-icon>
          </el-button>
          <el-button circle text :disabled="!hasNext" @click="playNext">
            <el-icon><ArrowRight /></el-icon>
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
        <el-slider v-model="volume" class="vol-slider" :min="0" :max="1" :step="0.05" @change="setVolume" />
        <el-button text @click="closePlayer">关闭</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft,
  ArrowRight,
  Calendar,
  Delete,
  DocumentCopy,
  Download,
  EditPen,
  Headset,
  Plus,
  Refresh,
  Search,
  Share,
  VideoPause,
  VideoPlay,
} from '@element-plus/icons-vue'
import { worksApi, type Work } from '@/api/works'
import { downloadAudio } from '@/utils/audio'

type WorkView = Work

const router = useRouter()
const route = useRoute()

const worksList = ref<WorkView[]>([])
const loading = ref(false)

const statusFilter = ref<'all' | 'processing' | 'completed' | 'failed'>('all')
const keyword = ref('')

const page = ref(1)
const pageSize = ref(9)

const autoRefreshTimer = ref<number | null>(null)

const hasRunning = computed(() =>
  worksList.value.some((w) => w.status === 'pending' || w.status === 'processing')
)

const filteredWorks = computed(() => {
  const list = [...worksList.value].sort((a, b) => {
    const ta = new Date(a.createdAt).getTime()
    const tb = new Date(b.createdAt).getTime()
    if (Number.isNaN(ta) && Number.isNaN(tb)) return 0
    if (Number.isNaN(ta)) return 1
    if (Number.isNaN(tb)) return -1
    return tb - ta
  })

  const f = statusFilter.value
  const byStatus =
    f === 'all'
      ? list
      : f === 'processing'
        ? list.filter((w) => isProcessing(w))
        : f === 'completed'
          ? list.filter((w) => isCompleted(w))
          : list.filter((w) => isFailed(w))

  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return byStatus
  return byStatus.filter((w) => {
    const t = `${workTitle(w)} ${w.taskId || ''}`.toLowerCase()
    return t.includes(kw)
  })
})

const pagedWorks = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredWorks.value.slice(start, start + pageSize.value)
})

watch([statusFilter, keyword], () => {
  page.value = 1
})

const pendingAutoPlay = ref(true)

async function refreshList(silent = false) {
  if (!silent) loading.value = true
  try {
    const res = await worksApi.list()
    worksList.value = res.list || []
    syncAutoRefresh()
    autoPlayFromQuery()
  } catch (e) {
    console.error(e)
  } finally {
    if (!silent) loading.value = false
  }
}

function handleCreate() {
  router.push('/generate')
}

function progressPercent(row: Work): number {
  const p = Number(row.progress)
  if (Number.isNaN(p)) return 0
  if (row.status === 'success') return 100
  if (row.status === 'failed') return Math.max(0, Math.min(100, p))
  return Math.max(0, Math.min(99, p))
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

async function copyText(text: string) {
  try {
    await writeToClipboard(text)
    ElMessage.success('已复制')
  } catch {
    ElMessage.error('复制失败')
  }
}

function formatDateTime(raw: string): string {
  const d = new Date(raw)
  if (Number.isNaN(d.getTime())) return raw || '-'
  const pad = (n: number) => `${n}`.padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function formatDate(raw: string): string {
  const d = new Date(raw)
  if (Number.isNaN(d.getTime())) return raw || '-'
  const pad = (n: number) => `${n}`.padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

function downloadName(row: Work): string {
  const ext = (row.format || 'mp3').toLowerCase()
  return `tts_${row.taskId}.${ext}`
}

async function handleDownload(row: Work) {
  if (!row.audioUrl) return
  try {
    await downloadAudio(row.audioUrl, downloadName(row))
    ElMessage.success('下载成功')
  } catch {
    ElMessage.error('下载失败')
  }
}

async function handleShare(row: Work) {
  if (!row.audioUrl) return
  try {
    await writeToClipboard(row.audioUrl)
    ElMessage.success('已复制分享链接')
  } catch {
    ElMessage.error('复制失败')
  }
}

function handleEdit(row: WorkView) {
  router.push({ path: '/generate', query: { fromTaskId: row.taskId } })
}

async function handleRename(row: WorkView) {
  try {
    const res = await ElMessageBox.prompt('请输入新的作品名称', '重命名', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputValue: workTitle(row),
    })
    const next = `${res.value || ''}`.trim()
    if (!next) return
    await worksApi.updateTitle(row.taskId, next)
    row.title = next
    ElMessage.success('已重命名')
  } catch (err: any) {
    if (err !== 'cancel') ElMessage.error(err?.message || '重命名失败')
  }
}

async function handleDelete(row: Work) {
  try {
    await ElMessageBox.confirm('确定要删除这个作品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await worksApi.delete(row.taskId)
    ElMessage.success('删除成功')
    if (playingTaskId.value === row.taskId) closePlayer()
    worksList.value = worksList.value.filter((w) => w.taskId !== row.taskId)
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '删除失败')
    }
  }
}

function syncAutoRefresh() {
  if (!hasRunning.value) {
    stopAutoRefresh()
    return
  }
  if (autoRefreshTimer.value !== null) return
  autoRefreshTimer.value = window.setInterval(() => {
    void refreshList(true)
  }, 3000)
}

function stopAutoRefresh() {
  if (autoRefreshTimer.value !== null) {
    clearInterval(autoRefreshTimer.value)
    autoRefreshTimer.value = null
  }
}

function isCompleted(w: WorkView) {
  return w.status === 'success' && !!w.audioUrl
}

function isProcessing(w: WorkView) {
  return w.status === 'pending' || w.status === 'processing'
}

function isFailed(w: WorkView) {
  return w.status === 'failed'
}

function statusText(w: WorkView) {
  if (isCompleted(w)) return '已完成'
  if (isProcessing(w)) return '处理中'
  if (isFailed(w)) return '失败'
  return '未知'
}

function statusClass(w: WorkView) {
  if (isCompleted(w)) return 'ok'
  if (isProcessing(w)) return 'running'
  if (isFailed(w)) return 'bad'
  return 'unknown'
}

function coverClass(w: WorkView) {
  if (isCompleted(w)) return 'c1'
  if (isProcessing(w)) return 'c3'
  if (isFailed(w)) return 'c4'
  return 'c4'
}

function shortTaskId(id: string): string {
  if (!id) return ''
  if (id.length <= 14) return id
  return `${id.slice(0, 8)}...${id.slice(-4)}`
}

function workTitle(w: WorkView): string {
  return w.title?.trim() || `音频项目 ${shortTaskId(w.taskId)}`
}

function autoPlayFromQuery() {
  if (!pendingAutoPlay.value) return
  const play = `${route.query.playTaskId || ''}`.trim()
  if (!play) return
  const target = worksList.value.find((x) => x.taskId === play)
  if (target && isCompleted(target)) {
    playWork(target)
    pendingAutoPlay.value = false
  }
}

// Bottom player bar
const audioEl = ref<HTMLAudioElement | null>(null)
const playerVisible = ref(false)
const isPlaying = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const volume = ref(1)
const playerUrl = ref<string>('')
const playingTaskId = ref<string | null>(null)

const playableList = computed(() => filteredWorks.value.filter((w) => isCompleted(w)))
const playingIndex = computed(() => playableList.value.findIndex((w) => w.taskId === playingTaskId.value))
const hasPrev = computed(() => playingIndex.value > 0)
const hasNext = computed(() => playingIndex.value >= 0 && playingIndex.value < playableList.value.length - 1)
const playingWork = computed(() => playableList.value[playingIndex.value] || null)
const playingTitle = computed(() => (playingWork.value ? workTitle(playingWork.value) : '-'))
const playingSub = computed(() =>
  playingWork.value ? `${formatDateTime(playingWork.value.createdAt)} · ${(playingWork.value.format || 'mp3').toUpperCase()}` : ''
)

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

function playWork(w: WorkView) {
  if (!w.audioUrl || !isCompleted(w)) return
  playerVisible.value = true
  playingTaskId.value = w.taskId
  playerUrl.value = w.audioUrl
  currentTime.value = 0
  duration.value = 0
  const el = audioEl.value
  if (!el) return
  el.currentTime = 0
  void el.play().catch((err) => {
    ElMessage.error(err?.message || '无法播放音频')
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
  playingTaskId.value = null
  playerUrl.value = ''
  currentTime.value = 0
  duration.value = 0
}

function playPrev() {
  if (!hasPrev.value) return
  const target = playableList.value[playingIndex.value - 1]
  if (target) playWork(target)
}

function playNext() {
  if (!hasNext.value) return
  const target = playableList.value[playingIndex.value + 1]
  if (target) playWork(target)
}

function formatTime(seconds: number): string {
  const s = Math.max(0, Math.floor(seconds || 0))
  const m = Math.floor(s / 60)
  const r = s % 60
  return `${m.toString().padStart(2, '0')}:${r.toString().padStart(2, '0')}`
}

onMounted(() => {
  refreshList()
})

onUnmounted(() => {
  stopAutoRefresh()
  closePlayer()
})
</script>

<style scoped>
.works-page {
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

.work-card {
  width: 340px;
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
}

.cover.c1 {
  background: rgba(64, 158, 255, 0.12);
  color: #409eff;
}
.cover.c3 {
  background: rgba(103, 194, 58, 0.12);
  color: #67c23a;
}
.cover.c4 {
  background: rgba(245, 108, 108, 0.12);
  color: #f56c6c;
}

.note {
  font-size: 18px;
}

.status-pill {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 999px;
  border: 1px solid transparent;
  user-select: none;
}
.status-pill.ok {
  color: #67c23a;
  background: rgba(103, 194, 58, 0.12);
  border-color: rgba(103, 194, 58, 0.2);
}
.status-pill.running {
  color: #409eff;
  background: rgba(64, 158, 255, 0.12);
  border-color: rgba(64, 158, 255, 0.2);
}
.status-pill.bad {
  color: #f56c6c;
  background: rgba(245, 108, 108, 0.12);
  border-color: rgba(245, 108, 108, 0.2);
}

.card-main {
  display: flex;
  flex-direction: column;
  gap: 8px;
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
  flex-wrap: wrap;
  font-size: 12px;
  color: #8c8c8c;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
}

.copy-btn {
  margin-left: 4px;
}

.err {
  font-size: 12px;
  color: #f56c6c;
  background: rgba(245, 108, 108, 0.06);
  border: 1px solid rgba(245, 108, 108, 0.14);
  border-radius: 10px;
  padding: 8px 10px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.processing-box {
  padding-top: 6px;
}

.processing-tip {
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
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
  .works-page {
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

  .work-card {
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
