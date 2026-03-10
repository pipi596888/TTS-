<template>
  <div class="gen-page">
    <div class="mode-tabs">
      <el-tabs v-model="mode" class="tabs">
        <el-tab-pane label="单角色模式" name="single" />
        <el-tab-pane label="多角色模式" name="multi" />
      </el-tabs>
    </div>

    <div class="layout">
      <aside class="left">
        <div class="card">
          <div class="card-title">
            <div class="title-left">
              <el-icon><Microphone /></el-icon>
              <span>音色选择</span>
            </div>
          </div>

          <div class="role-list">
            <div
              v-for="r in roles"
              :key="r.id"
              class="role-item"
              :class="{ active: r.id === activeRoleId }"
              @click="selectRole(r.id)"
            >
              <div class="role-main">
                <div class="role-name">{{ r.name }}</div>
                <div class="role-meta">{{ voiceMeta(r.voiceId) }}</div>
              </div>
              <div class="role-actions">
                <el-button size="small" @click.stop="openEditRole(r)">
                  设置
                </el-button>
                <el-button size="small" type="danger" plain @click.stop="deleteRole(r)">
                  删除
                </el-button>
              </div>
            </div>

            <el-empty v-if="roles.length === 0" description="暂无角色" :image-size="72" />

            <el-button class="add-role" text @click="openAddRole">
              <el-icon><Plus /></el-icon>
              添加角色
            </el-button>
          </div>
        </div>

        <div class="card">
          <div class="card-subtitle">语言增强</div>
          <el-select v-model="enhanceLang" class="full" size="default">
            <el-option label="中文" value="zh" />
            <el-option label="英文" value="en" />
            <el-option label="日文" value="ja" />
          </el-select>
        </div>

        <div class="card">
          <div class="card-subtitle">音频格式</div>
          <el-radio-group v-model="format" class="radio-row">
            <el-radio value="mp3">MP3</el-radio>
            <el-radio value="wav">WAV</el-radio>
            <el-radio value="flac">FLAC</el-radio>
            <el-radio value="pcm">PCM</el-radio>
          </el-radio-group>
        </div>

        <div class="card">
          <div class="card-subtitle">声道</div>
          <el-radio-group v-model="channel" class="radio-row">
            <el-radio value="mono">单声道</el-radio>
            <el-radio value="stereo">双声道</el-radio>
          </el-radio-group>
        </div>

        <div class="left-actions">
          <el-button class="full" @click="resetParams">
            <el-icon><Refresh /></el-icon>
            重置参数
          </el-button>
          <el-button class="full" type="primary" plain @click="goToWorks">
            前往作品
          </el-button>
        </div>
      </aside>

      <main class="main">
        <div class="card main-card">
          <div class="seg-header">
            <div class="seg-title">
              <span class="title">片段列表</span>
              <span class="count">（共 {{ segmentsOrdered.length }} 条）</span>
              <el-button class="add-at-top" text @click="addAtTop">
                <el-icon><Plus /></el-icon>
                在开头添加
              </el-button>
            </div>

            <div class="seg-tools">
              <el-button type="primary" plain @click="smartProcess = !smartProcess">
                <el-icon><MagicStick /></el-icon>
                智能处理文本
              </el-button>
              <el-input
                v-model="workTitle"
                class="title-input"
                clearable
                placeholder="作品名称（必填）"
                :maxlength="60"
                show-word-limit
              />
              <el-checkbox v-model="removeNarration">去掉旁白</el-checkbox>
              <el-button type="primary" :loading="isTaskRunning" :disabled="!canGenerateAll" @click="generateAll">
                全部生成
              </el-button>
              <span class="price-hint">音频生成：每字 2 积分</span>
            </div>
          </div>

          <div class="seg-list">
            <div v-if="segmentsOrdered.length === 0" class="empty">
              <el-empty description="暂无片段" :image-size="120">
                <el-button type="primary" @click="addAtTop">
                  <el-icon><Plus /></el-icon>
                  添加片段
                </el-button>
              </el-empty>
            </div>

            <div v-for="(seg, idx) in segmentsOrdered" :key="seg.id" class="seg-item">
              <div class="seg-row">
                <div class="seg-tags">
                  <el-tag size="small" type="info" effect="plain">{{ roleNameByVoice(seg.voiceId) }}</el-tag>
                  <el-tag v-if="seg.emotion" size="small" type="warning" effect="plain">{{ emotionLabel(seg.emotion) }}</el-tag>
                </div>

                <div class="seg-actions">
                  <el-button text @click="openEditText(seg)">
                    <el-icon><EditPen /></el-icon>
                    文本修改
                  </el-button>
                  <el-button text @click="openSegSettings(seg)">
                    <el-icon><Setting /></el-icon>
                    设置
                  </el-button>
                  <el-button text type="danger" @click="deleteSegment(seg.id)">
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-button>
                  <el-button type="primary" plain @click="generateOne(seg)">
                    生成
                  </el-button>
                </div>
              </div>

              <div class="seg-text" @dblclick="openEditText(seg)">
                {{ seg.text || '（空文本）' }}
              </div>

              <el-button class="add-below" text @click="addBelow(idx)">
                <el-icon><Plus /></el-icon>
                在此下方添加片段
              </el-button>
            </div>
          </div>
        </div>

        <transition name="slide-up">
          <div v-if="taskStatus?.audioUrl" class="card player-card">
            <AudioPlayer :audio-url="taskStatus.audioUrl" />
          </div>
        </transition>
      </main>
    </div>

    <el-dialog v-model="roleDialogVisible" :title="roleDialogTitle" width="460px">
      <el-form label-width="88px">
        <el-form-item label="角色名">
          <el-input v-model="roleForm.name" placeholder="例如：旁白 / 如冬" />
        </el-form-item>
        <el-form-item label="音色">
          <el-select v-model="roleForm.voiceId" class="full" filterable placeholder="选择音色">
            <el-option
              v-for="v in voiceList"
              :key="v.id"
              :label="`${v.name}（${v.gender} / ${v.tone}）`"
              :value="v.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmRole">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="textDialogVisible" title="文本修改" width="720px">
      <el-input v-model="textDraft" type="textarea" :rows="8" placeholder="请输入文本..." />
      <template #footer>
        <el-button @click="textDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmText">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="settingsDialogVisible" title="片段设置" width="520px">
      <el-form label-width="88px">
        <el-form-item label="角色">
          <el-select v-model="settingsDraft.voiceId" class="full">
            <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.voiceId" />
          </el-select>
        </el-form-item>
        <el-form-item label="情绪">
          <el-select v-model="settingsDraft.emotion" class="full" clearable placeholder="可选">
            <el-option label="中性" value="neutral" />
            <el-option label="开心" value="happy" />
            <el-option label="悲伤" value="sad" />
            <el-option label="愤怒" value="angry" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="settingsDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmSettings">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, EditPen, MagicStick, Microphone, Plus, Refresh, Setting } from '@element-plus/icons-vue'
import AudioPlayer from '@/components/AudioPlayer/index.vue'
import { useTTSStore } from '@/store/tts'
import { useVoiceStore } from '@/store/voice'
import type { Segment } from '@/types/tts'
import { ttsApi } from '@/api/tts'
import { worksApi } from '@/api/works'

type Mode = 'single' | 'multi'

type Role = {
  id: string
  name: string
  voiceId: number
}

const router = useRouter()
const route = useRoute()
const ttsStore = useTTSStore()
const voiceStore = useVoiceStore()

const { segments, format, channel, taskStatus, isTaskRunning } = storeToRefs(ttsStore)
const { voiceList, loading: voiceLoading } = storeToRefs(voiceStore)

const mode = ref<Mode>('multi')
const enhanceLang = ref<'zh' | 'en' | 'ja'>('zh')
const smartProcess = ref(true)
const removeNarration = ref(false)
const workTitle = ref('')

const roles = ref<Role[]>([])
const activeRoleId = ref<string>('')

const segmentsOrdered = computed(() => {
  return [...segments.value].sort((a, b) => a.order - b.order)
})

const effectiveSegments = computed(() => {
  const base = segmentsOrdered.value
  return removeNarration.value ? base.filter((s) => !isNarration(s)) : base
})

const canGenerateAll = computed(() => effectiveSegments.value.length > 0)

function getQueryString(v: unknown): string {
  if (typeof v === 'string') return v
  if (Array.isArray(v) && typeof v[0] === 'string') return v[0]
  return ''
}

const fromTaskId = computed(() => getQueryString(route.query.fromTaskId))

function roleNameByVoice(voiceId: number) {
  const r = roles.value.find((x) => x.voiceId === voiceId)
  return r?.name || voiceStore.getVoiceById(voiceId)?.name || `音色 ${voiceId}`
}

function voiceMeta(voiceId: number) {
  const v = voiceStore.getVoiceById(voiceId)
  if (!v) return ''
  return `音色：${v.gender} / ${v.tone}`
}

function emotionLabel(emotion: string) {
  const map: Record<string, string> = {
    neutral: '中性',
    happy: '高兴',
    sad: '悲伤',
    angry: '愤怒',
  }
  return map[emotion] || emotion
}

function selectRole(roleId: string) {
  activeRoleId.value = roleId
  const r = roles.value.find((x) => x.id === roleId)
  if (!r) return
  ttsStore.selectedVoiceId = r.voiceId

  if (mode.value === 'single') {
    ttsStore.setAllSegmentsVoice(r.voiceId)
  }
}

function activeVoiceId(): number {
  const r = roles.value.find((x) => x.id === activeRoleId.value)
  return r?.voiceId ?? ttsStore.selectedVoiceId
}

function ensureRolesFromVoices() {
  if (voiceList.value.length === 0) {
    roles.value = []
    activeRoleId.value = ''
    return
  }

  if (roles.value.length > 0) return

  const initial =
    mode.value === 'multi' ? voiceList.value.slice(0, Math.min(3, voiceList.value.length)) : voiceList.value.slice(0, 1)

  roles.value = initial.map((v) => ({
    id: `role_${v.id}`,
    name: v.name,
    voiceId: v.id,
  }))

  activeRoleId.value = roles.value[0]!.id
  ttsStore.selectedVoiceId = roles.value[0]!.voiceId
}

watch(
  () => mode.value,
  (m) => {
    if (m === 'single') {
      // 保留当前激活角色（或第一个），并同步到所有片段
      if (roles.value.length === 0) return
      const keep = roles.value.find((r) => r.id === activeRoleId.value) ?? roles.value[0]!
      roles.value = [keep]
      activeRoleId.value = keep.id
      ttsStore.selectedVoiceId = keep.voiceId
      ttsStore.setAllSegmentsVoice(keep.voiceId)
    } else {
      // 多角色：不强制回填片段 voiceId，仅保证有角色
      ensureRolesFromVoices()
    }
  }
)

onMounted(async () => {
  await voiceStore.fetchVoiceList()
  ensureRolesFromVoices()
  await loadFromWorkIfNeeded()
  if (voiceLoading.value) return
})

const roleDialogVisible = ref(false)
const roleDialogTitle = ref('添加角色')
const roleEditingId = ref<string>('')
const roleForm = ref<{ name: string; voiceId: number | null }>({ name: '', voiceId: null })

function openAddRole() {
  if (mode.value === 'single') {
    ElMessage.info('单角色模式下仅保留一个角色，可切换到多角色模式添加更多角色')
    return
  }
  roleEditingId.value = ''
  roleDialogTitle.value = '添加角色'
  const firstVoice = voiceList.value.find((v) => !roles.value.some((r) => r.voiceId === v.id)) ?? voiceList.value[0]
  roleForm.value = {
    name: firstVoice?.name ?? '',
    voiceId: firstVoice?.id ?? null,
  }
  roleDialogVisible.value = true
}

function openEditRole(r: Role) {
  roleEditingId.value = r.id
  roleDialogTitle.value = '角色设置'
  roleForm.value = { name: r.name, voiceId: r.voiceId }
  roleDialogVisible.value = true
}

function confirmRole() {
  const voiceId = roleForm.value.voiceId
  const name = roleForm.value.name.trim()
  if (!voiceId || !name) {
    ElMessage.warning('请填写角色名并选择音色')
    return
  }

  if (roleEditingId.value) {
    roles.value = roles.value.map((r) => (r.id === roleEditingId.value ? { ...r, name, voiceId } : r))
    if (activeRoleId.value === roleEditingId.value) {
      ttsStore.selectedVoiceId = voiceId
      if (mode.value === 'single') ttsStore.setAllSegmentsVoice(voiceId)
    }
  } else {
    const id = `role_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
    roles.value.push({ id, name, voiceId })
    if (!activeRoleId.value) selectRole(id)
  }

  roleDialogVisible.value = false
}

async function deleteRole(r: Role) {
  if (mode.value === 'single') return
  if (roles.value.length <= 1) {
    ElMessage.warning('至少保留一个角色')
    return
  }
  try {
    await ElMessageBox.confirm(`确定删除角色「${r.name}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }

  const fallback = roles.value.find((x) => x.id !== r.id) ?? roles.value[0]
  roles.value = roles.value.filter((x) => x.id !== r.id)
  if (activeRoleId.value === r.id && fallback) {
    selectRole(fallback.id)
  }
  // 将使用该 voiceId 的片段回退到 fallback voiceId
  if (fallback) {
    segments.value
      .filter((s) => s.voiceId === r.voiceId)
      .forEach((s) => ttsStore.updateSegment(s.id, { voiceId: fallback.voiceId }))
  }
}

function resetParams() {
  enhanceLang.value = 'zh'
  smartProcess.value = true
  removeNarration.value = false
  format.value = 'mp3'
  channel.value = 'mono'
  ttsStore.resetTask()
  ElMessage.success('已重置')
}

function goToWorks() {
  const taskId = ttsStore.currentTaskId || taskStatus.value?.taskId || fromTaskId.value
  router.push({
    path: '/works',
    query: taskId ? { playTaskId: taskId } : undefined,
  })
}

function normalizeLoadedVoiceId(voiceId: number): number {
  if (voiceStore.getVoiceById(voiceId)) return voiceId
  return voiceStore.getDefaultVoice()?.id ?? voiceList.value[0]?.id ?? voiceId
}

async function loadFromWorkIfNeeded() {
  if (!fromTaskId.value) return

  try {
    const detail = await ttsApi.getTaskDetail(fromTaskId.value)

    ttsStore.resetTask()

    const sorted = [...(detail.segments || [])].sort((a, b) => a.sort - b.sort)
    const loadedSegments: Segment[] = sorted.map((s, idx) => ({
      id: `seg_${fromTaskId.value}_${s.sort}_${idx}`,
      voiceId: normalizeLoadedVoiceId(s.voiceId),
      emotion: s.emotion,
      text: s.text ?? '',
      order: idx,
    }))

    segments.value = loadedSegments

    const f = detail.format as any
    if (f === 'mp3' || f === 'wav' || f === 'flac' || f === 'pcm') format.value = f
    const c = detail.channel as any
    if (c === 'mono' || c === 'stereo') channel.value = c

    const voiceIds = Array.from(new Set(loadedSegments.map((s) => s.voiceId).filter(Boolean)))
    const ids =
      voiceIds.length > 0
        ? voiceIds
        : [voiceStore.getDefaultVoice()?.id ?? voiceList.value[0]?.id ?? 0].filter(Boolean)

    roles.value = ids.map((voiceId, index) => ({
      id: `role_${voiceId}`,
      name: voiceStore.getVoiceById(voiceId)?.name ?? `角色 ${index + 1}`,
      voiceId,
    }))

    activeRoleId.value = roles.value[0]?.id ?? ''
    ttsStore.selectedVoiceId = roles.value[0]?.voiceId ?? 0
    mode.value = ids.length > 1 ? 'multi' : 'single'

    if (detail.title) {
      workTitle.value = `${detail.title}（编辑）`
    }

    ElMessage.success('已载入作品，可继续编辑并重新生成')
  } catch (err: any) {
    ElMessage.error(err?.message || '载入作品失败')
  }
}

async function ensureWorkTitle(): Promise<string> {
  const trimmed = workTitle.value.trim()
  if (trimmed) return trimmed

  try {
    const { value } = await ElMessageBox.prompt('请输入作品名称', '作品名称', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '例如：科技新闻播报-11月特刊',
      inputValue: workTitle.value,
      inputValidator: (v) => {
        const t = String(v || '').trim()
        if (!t) return '作品名称不能为空'
        if (t.length > 60) return '最多 60 个字符'
        return true
      },
    })
    workTitle.value = String(value || '').trim()
    return workTitle.value
  } catch {
    throw new Error('已取消')
  }
}

function addAtTop() {
  const id = ttsStore.insertSegmentAt(0, '', activeVoiceId())
  openEditTextById(id)
}

function addBelow(index: number) {
  const id = ttsStore.insertSegmentAt(index + 1, '', activeVoiceId())
  openEditTextById(id)
}

function deleteSegment(id: string) {
  ttsStore.removeSegment(id)
  ttsStore.resetTask()
}

const textDialogVisible = ref(false)
const editingTextSegId = ref<string>('')
const textDraft = ref('')

function openEditText(seg: Segment) {
  editingTextSegId.value = seg.id
  textDraft.value = seg.text
  textDialogVisible.value = true
}

function openEditTextById(id: string) {
  const seg = segmentsOrdered.value.find((s) => s.id === id)
  if (!seg) return
  openEditText(seg)
}

function confirmText() {
  const id = editingTextSegId.value
  if (!id) return
  ttsStore.updateSegment(id, { text: textDraft.value })
  textDialogVisible.value = false
}

const settingsDialogVisible = ref(false)
const editingSettingsSegId = ref<string>('')
const settingsDraft = ref<{ voiceId: number; emotion: string | '' }>({ voiceId: 0, emotion: '' })

function openSegSettings(seg: Segment) {
  editingSettingsSegId.value = seg.id
  settingsDraft.value = {
    voiceId: seg.voiceId || activeVoiceId(),
    emotion: (seg.emotion as any) || '',
  }
  settingsDialogVisible.value = true
}

function confirmSettings() {
  const id = editingSettingsSegId.value
  if (!id) return
  ttsStore.updateSegment(id, {
    voiceId: settingsDraft.value.voiceId,
    emotion: settingsDraft.value.emotion || undefined,
  })
  settingsDialogVisible.value = false
}

function normalizeText(text: string) {
  return text
    .replace(/\r\n/g, '\n')
    .replace(/[ \t]+\n/g, '\n')
    .replace(/\n{3,}/g, '\n\n')
    .trim()
}

function isNarration(seg: Segment) {
  const rn = roleNameByVoice(seg.voiceId)
  return rn.includes('旁白')
}

async function generateAll() {
  let title = ''
  try {
    title = await ensureWorkTitle()
  } catch (err: any) {
    if (err?.message !== '已取消') ElMessage.error(err?.message || '请填写作品名称')
    return
  }

  const base = segmentsOrdered.value
  const filtered = removeNarration.value ? base.filter((s) => !isNarration(s)) : base
  const override = smartProcess.value
    ? filtered.map((s) => ({ ...s, text: normalizeText(s.text) }))
    : filtered.map((s) => ({ ...s }))

  try {
    const taskId = await ttsStore.generateAudio(override)
    try {
      await worksApi.updateTitle(taskId, title)
    } catch (e: any) {
      console.warn('Failed to set work title:', e)
      ElMessage.warning('已开始生成，但作品名称写入失败，可在作品页重命名')
    }
  } catch (err: any) {
    ElMessage.error(err?.message || '生成失败')
  }
}

async function generateOne(seg: Segment) {
  const override: Segment = {
    ...seg,
    text: smartProcess.value ? normalizeText(seg.text) : seg.text,
  }
  try {
    await ttsStore.generateAudio([override])
  } catch (err: any) {
    ElMessage.error(err?.message || '生成失败')
  }
}
</script>

<style scoped>
.gen-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px;
  gap: 12px;
  overflow: hidden;
}

.mode-tabs {
  background: #fff;
  border-radius: 12px;
  padding: 0 12px;
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.06);
}

.tabs :deep(.el-tabs__header) {
  margin: 0;
}

.layout {
  flex: 1;
  min-height: 0;
  display: flex;
  gap: 16px;
  overflow: hidden;
}

.left {
  width: 320px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow: auto;
  padding-right: 4px;
}

.main {
  flex: 1;
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow: hidden;
}

.card {
  background: #fff;
  border-radius: 14px;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.06);
  padding: 14px;
}

.card-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 12px;
}

.title-left {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #1f1f1f;
}

.card-subtitle {
  font-size: 13px;
  color: #5a5a5a;
  font-weight: 600;
  margin-bottom: 10px;
}

.full {
  width: 100%;
}

.radio-row :deep(.el-radio) {
  margin-right: 14px;
}

.role-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.role-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.15s ease;
  background: #fafafa;
}

.role-item:hover {
  background: #f4f7ff;
}

.role-item.active {
  border-color: #409eff;
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.14), rgba(64, 158, 255, 0.06));
}

.role-main {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.role-name {
  font-weight: 700;
  color: #1f1f1f;
  font-size: 14px;
  line-height: 1.2;
}

.role-meta {
  font-size: 12px;
  color: #7b7b7b;
  line-height: 1.2;
}

.role-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.add-role {
  justify-content: flex-start;
  padding-left: 2px;
}

.left-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.main-card {
  flex: 1;
  min-height: 0;
  padding: 0;
  overflow: hidden;
}

.seg-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid #f0f0f0;
}

.seg-title {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.seg-title .title {
  font-weight: 700;
  color: #1f1f1f;
}

.seg-title .count {
  color: #8c8c8c;
  font-size: 13px;
}

.add-at-top {
  margin-left: 6px;
}

.seg-tools {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.title-input {
  width: 260px;
}

.price-hint {
  font-size: 12px;
  color: #8c8c8c;
  white-space: nowrap;
}

.seg-list {
  padding: 14px 16px;
  overflow: auto;
  height: 100%;
}

.seg-item {
  border: 1px solid #eef0f5;
  border-radius: 14px;
  padding: 12px;
  background: #fff;
  margin-bottom: 12px;
}

.seg-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.seg-tags {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.seg-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.seg-text {
  margin-top: 10px;
  padding: 10px 12px;
  background: #fafafa;
  border-radius: 12px;
  color: #2a2a2a;
  line-height: 1.7;
  white-space: pre-wrap;
}

.add-below {
  margin-top: 6px;
  padding-left: 2px;
  justify-content: flex-start;
}

.player-card {
  padding: 14px 16px;
}

.empty {
  padding: 32px 0;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.25s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

@media (max-width: 980px) {
  .gen-page {
    padding: 12px;
  }

  .layout {
    flex-direction: column;
  }

  .left {
    width: 100%;
    max-height: 420px;
    padding-right: 0;
  }

  .seg-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .seg-tools {
    width: 100%;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .title-input {
    width: 100%;
  }
}
</style>


