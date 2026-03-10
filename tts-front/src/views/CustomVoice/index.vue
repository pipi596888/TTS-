<template>
  <div class="cv-page">
    <div class="page-head">
      <div class="head-left">
        <div class="title">定制声音</div>
        <div class="sub">通过录音训练属于您的专属 AI 音色，打造极具辨识度的品牌声音</div>
      </div>
      <div class="head-right">
        <el-button type="primary" @click="openCreate">
          <el-icon><Plus /></el-icon>
          开始新克隆
        </el-button>
      </div>
    </div>

    <div class="steps">
      <div class="step-card">
        <div class="icon">
          <el-icon><Microphone /></el-icon>
        </div>
        <div class="step-title">1. 提供录音</div>
        <div class="step-desc">上传在线录音 5-10 分钟音频。音频质量越高，克隆效果越逼真。</div>
      </div>
      <div class="arrow">›</div>
      <div class="step-card">
        <div class="icon">
          <el-icon><Cpu /></el-icon>
        </div>
        <div class="step-title">2. 智能训练</div>
        <div class="step-desc">深度学习大模型自动提取声纹特征，通常在 15-30 分钟内完成。</div>
      </div>
      <div class="arrow">›</div>
      <div class="step-card">
        <div class="icon">
          <el-icon><Cloudy /></el-icon>
        </div>
        <div class="step-title">3. 模型部署</div>
        <div class="step-desc">一键部署至云端模型库，您可以立即在控制台或通过 API 调用。</div>
      </div>
    </div>

    <div class="list-head">
      <div class="lh-left">
        <div class="h">我的音色列表</div>
      </div>
      <div class="lh-right">
        <span class="auto">
          <span class="dot"></span>
          列表每 5 分钟自动更新状态
        </span>
        <el-button text :loading="loading" @click="refresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <div class="list-card" v-loading="loading">
      <div class="list-table">
        <div class="th">
          <div class="c name">音色名称</div>
          <div class="c time">创建时间</div>
          <div class="c prog">训练进度</div>
          <div class="c st">当前状态</div>
          <div class="c act">操作</div>
        </div>

        <div v-if="!loading && pagedItems.length === 0" class="empty">
          <el-empty description="暂无音色" :image-size="120" />
        </div>

        <div v-else class="tb">
          <div v-for="it in pagedItems" :key="it.id" class="tr" @click="openDetail(it)">
            <div class="c name">
              <div class="name-cell">
                <div class="avatar">
                  <el-icon><User /></el-icon>
                </div>
                <div class="nm">
                  <div class="n1">{{ it.name }}</div>
                  <div class="n2">音色: {{ it.tone || '-' }}</div>
                </div>
              </div>
            </div>
            <div class="c time">{{ formatDateTime(it.createdAt) }}</div>
            <div class="c prog" @click.stop>
              <div class="prog-row">
                <el-progress :percentage="progressPercent(it)" :stroke-width="10" />
                <div class="pct">{{ progressPercent(it) }}%</div>
              </div>
            </div>
            <div class="c st">
              <span class="status-pill" :class="statusClass(it)">{{ statusText(it) }}</span>
            </div>
            <div class="c act" @click.stop>
              <el-button v-if="canPreview(it)" link type="primary" @click="preview(it)">立即试听</el-button>
              <span v-else class="muted">试听暂不可用</span>
            </div>
          </div>
        </div>
      </div>

      <div class="pager" v-if="filteredItems.length > 0">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[4, 6, 8, 12]"
          layout="prev, pager, next"
          :total="filteredItems.length"
        />
      </div>
    </div>

    <div class="metrics">
      <div class="m-card">
        <div class="k">总计音色</div>
        <div class="v">{{ items.length }}</div>
      </div>
      <div class="m-card">
        <div class="k">训练中</div>
        <div class="v">{{ trainingCount }}</div>
      </div>
      <div class="m-card">
        <div class="k">剩余额度</div>
        <div class="v">{{ quotaText }}</div>
      </div>
      <div class="m-card">
        <div class="k">平均耗时</div>
        <div class="v">{{ avgTimeText }}</div>
      </div>
    </div>

    <el-drawer v-model="playerVisible" title="试听音色" :size="playerDrawerSize" @close="closePlayer">
      <div v-if="playingItem" class="player-meta">
        <div class="row"><span class="k">音色名称</span><span class="v">{{ playingItem.name }}</span></div>
        <div class="row"><span class="k">音色</span><span class="v">{{ playingItem.tone || '-' }}</span></div>
        <div class="row"><span class="k">状态</span><span class="v">{{ statusText(playingItem) }}</span></div>
      </div>
      <AudioPlayer v-if="playerUrl" :audio-url="playerUrl" :filename="`preview_${playingItem?.id || ''}.mp3`" />
      <el-empty v-else description="暂无可试听链接" :image-size="100" />
    </el-drawer>

    <el-drawer v-model="detailVisible" title="定制请求详情" :size="detailDrawerSize">
      <div v-if="current" class="detail">
        <div class="row"><span class="k">ID</span><span class="v mono">{{ current.id }}</span></div>
        <div class="row"><span class="k">名称</span><span class="v">{{ current.name }}</span></div>
        <div class="row"><span class="k">音色</span><span class="v">{{ current.tone || '-' }}</span></div>
        <div class="row"><span class="k">性别</span><span class="v">{{ current.gender || '-' }}</span></div>
        <div class="row">
          <span class="k">状态</span>
          <span class="v">
            <span class="status-pill" :class="statusClass(current)">{{ statusText(current) }}</span>
          </span>
        </div>
        <div class="row"><span class="k">创建时间</span><span class="v">{{ formatDateTime(current.createdAt) }}</span></div>

        <div class="block">
          <div class="block-title">样本文本</div>
          <div class="block-body pre">{{ current.sampleText || '-' }}</div>
        </div>

        <div class="block">
          <div class="block-title">样本链接</div>
          <div class="block-body">
            <div v-if="(current.sampleUrls || []).length === 0" class="muted">-</div>
            <div v-for="(u, i) in current.sampleUrls" :key="i" class="url">
              <el-link :href="u" target="_blank" type="primary">{{ u }}</el-link>
            </div>
          </div>
        </div>

        <div class="detail-actions" @click.stop>
          <el-button type="danger" plain @click="remove(current)">删除请求</el-button>
        </div>
      </div>
    </el-drawer>

    <el-dialog v-model="createVisible" title="开始新克隆" width="720px">
      <el-form label-width="96px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="例如：商务男声 / 客服助手" />
        </el-form-item>
        <el-form-item label="音色">
          <el-input v-model="form.tone" placeholder="例如：沉稳 / 轻柔 / 活泼" />
        </el-form-item>
        <el-form-item label="性别">
          <el-select v-model="form.gender" class="full">
            <el-option label="男" value="male" />
            <el-option label="女" value="female" />
            <el-option label="中性" value="neutral" />
          </el-select>
        </el-form-item>
        <el-form-item label="样本文本">
          <el-input v-model="form.sampleText" type="textarea" :rows="3" placeholder="可选：提供样本文本以提升训练效果" />
        </el-form-item>
        <el-form-item label="样本链接" required>
          <el-input v-model="urlsText" type="textarea" :rows="5" placeholder="每行一个音频链接（至少 1 条）" />
          <div class="sub-hint">建议提供清晰、无噪声、单人说话的录音链接。</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">提交</el-button>
      </template>
    </el-dialog>

    <div v-if="isAdmin" class="admin-area">
      <div class="admin-head">
        <div class="admin-title">管理员：全部请求</div>
        <el-button text :loading="adminLoading" @click="refreshAdmin">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
      <div class="admin-card">
        <el-table :data="adminItems" v-loading="adminLoading" style="width: 100%" row-key="id">
          <el-table-column prop="id" label="ID" width="90" />
          <el-table-column prop="userId" label="用户ID" width="110" />
          <el-table-column prop="name" label="名称" min-width="160" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="{ row }">
              <el-tag v-if="row.status === 'pending'" type="warning">待处理</el-tag>
              <el-tag v-else-if="row.status === 'success'" type="success">已通过</el-tag>
              <el-tag v-else type="danger">失败</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="时间" width="200">
            <template #default="{ row }">
              {{ formatDateTime(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="管理" width="240" fixed="right">
            <template #default="{ row }">
              <el-button size="small" type="success" plain :disabled="row.status !== 'pending'" @click="approve(row)">
                通过
              </el-button>
              <el-button size="small" type="danger" plain :disabled="row.status !== 'pending'" @click="openReject(row)">
                拒绝
              </el-button>
              <el-button size="small" text @click="openDetail(row, true)">详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <el-dialog v-model="rejectVisible" title="拒绝原因" width="520px">
      <el-input v-model="rejectText" type="textarea" :rows="5" placeholder="请输入拒绝原因（可选）" />
      <template #footer>
        <el-button @click="rejectVisible = false">取消</el-button>
        <el-button type="danger" :loading="rejecting" @click="reject">确认拒绝</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Cloudy, Cpu, Microphone, Plus, Refresh, User } from '@element-plus/icons-vue'
import { customVoiceApi, type CustomVoiceRequest } from '@/api/customVoice'
import { useUserStore } from '@/store/user'
import { useVoiceStore } from '@/store/voice'
import AudioPlayer from '@/components/AudioPlayer/index.vue'

const userStore = useUserStore()
const { userInfo } = storeToRefs(userStore)
const isAdmin = computed(() => (userInfo.value?.id || 0) === 1)

const voiceStore = useVoiceStore()

const loading = ref(false)
const items = ref<CustomVoiceRequest[]>([])

const page = ref(1)
const pageSize = ref(4)

const filteredItems = computed(() => {
  return [...items.value].sort((a, b) => {
    const ta = new Date(a.createdAt).getTime()
    const tb = new Date(b.createdAt).getTime()
    if (Number.isNaN(ta) && Number.isNaN(tb)) return 0
    if (Number.isNaN(ta)) return 1
    if (Number.isNaN(tb)) return -1
    return tb - ta
  })
})

const pagedItems = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredItems.value.slice(start, start + pageSize.value)
})

watch(pageSize, () => {
  page.value = 1
})

const trainingCount = computed(() => items.value.filter((i) => i.status === 'pending' || i.status === 'processing').length)
const quotaLimit = 15
const quotaText = computed(() => `${Math.min(items.value.length, quotaLimit)}/${quotaLimit}`)
const avgTimeText = computed(() => '18 分钟')

function formatDateTime(raw: string): string {
  const d = new Date(raw)
  if (Number.isNaN(d.getTime())) return raw || '-'
  const pad = (n: number) => `${n}`.padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function statusText(it: CustomVoiceRequest) {
  if (it.status === 'success') return '已可用'
  if (it.status === 'failed') return '失败'
  return '训练中'
}

function statusClass(it: CustomVoiceRequest) {
  if (it.status === 'success') return 'ok'
  if (it.status === 'failed') return 'bad'
  return 'run'
}

function progressPercent(it: CustomVoiceRequest): number {
  if (it.status === 'success') return 100
  if (it.status === 'failed') return 100
  const t = new Date(it.createdAt).getTime()
  if (!Number.isFinite(t)) return 12
  const mins = Math.max(0, (Date.now() - t) / 1000 / 60)
  return Math.max(5, Math.min(95, Math.round(mins * 3)))
}

const createVisible = ref(false)
const submitting = ref(false)
const form = ref<{ name: string; gender: string; tone: string; sampleText: string }>({
  name: '',
  gender: 'female',
  tone: '',
  sampleText: '',
})
const urlsText = ref('')

function openCreate() {
  form.value = { name: '', gender: 'female', tone: '', sampleText: '' }
  urlsText.value = ''
  createVisible.value = true
}

function toUrls(): string[] {
  return urlsText.value
    .split('\n')
    .map((s) => s.trim())
    .filter(Boolean)
}

async function submit() {
  const urls = toUrls()
  if (!form.value.name.trim()) {
    ElMessage.warning('请填写名称')
    return
  }
  if (urls.length === 0) {
    ElMessage.warning('请至少提供 1 条样本链接')
    return
  }

  submitting.value = true
  try {
    await customVoiceApi.createRequest({
      name: form.value.name,
      gender: form.value.gender,
      tone: form.value.tone,
      sampleText: form.value.sampleText,
      sampleUrls: urls,
    })
    ElMessage.success('已提交')
    createVisible.value = false
    await refresh()
    await voiceStore.fetchVoiceList()
    if (isAdmin.value) await refreshAdmin()
  } catch (err: any) {
    ElMessage.error(err?.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

// Detail drawer
const detailVisible = ref(false)
const current = ref<CustomVoiceRequest | null>(null)
const detailFromAdmin = ref(false)

function openDetail(row: CustomVoiceRequest, fromAdmin = false) {
  current.value = row
  detailFromAdmin.value = fromAdmin
  detailVisible.value = true
}

async function remove(row: CustomVoiceRequest) {
  try {
    await ElMessageBox.confirm(`确定删除请求「${row.name}」吗？`, '提示', { type: 'warning' })
    if (detailFromAdmin.value && isAdmin.value) {
      await customVoiceApi.adminDelete(row.id)
    } else {
      await customVoiceApi.delete(row.id)
    }
    ElMessage.success('已删除')
    await refresh()
    await voiceStore.fetchVoiceList()
    if (isAdmin.value) await refreshAdmin()
  } catch {
    // ignore
  }
}

async function refresh() {
  loading.value = true
  try {
    const res = await customVoiceApi.listMy()
    items.value = res.list || []
    await voiceStore.fetchVoiceList()
  } catch (err: any) {
    ElMessage.error(err?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

// Preview
const playerVisible = ref(false)
const playingItem = ref<CustomVoiceRequest | null>(null)
const playerUrl = ref<string>('')

const isMobile = ref(false)

function syncIsMobile() {
  isMobile.value = window.innerWidth < 980
}

const playerDrawerSize = computed(() => (isMobile.value ? '100%' : '520px'))
const detailDrawerSize = computed(() => (isMobile.value ? '100%' : '640px'))

function canPreview(it: CustomVoiceRequest) {
  if (it.status !== 'success') return false
  if (!it.resultVoiceId) return false
  const v = voiceStore.getVoiceById(it.resultVoiceId)
  return !!v?.previewUrl
}

async function preview(it: CustomVoiceRequest) {
  if (!it.resultVoiceId) return
  const v = voiceStore.getVoiceById(it.resultVoiceId)
  if (!v?.previewUrl) {
    ElMessage.warning('暂无可试听链接')
    return
  }
  playingItem.value = it
  playerUrl.value = v.previewUrl
  playerVisible.value = true
}

function closePlayer() {
  playerVisible.value = false
  playingItem.value = null
  playerUrl.value = ''
}

// Admin
const adminLoading = ref(false)
const adminItems = ref<CustomVoiceRequest[]>([])

async function refreshAdmin() {
  adminLoading.value = true
  try {
    const res = await customVoiceApi.adminList()
    adminItems.value = res.list || []
  } catch (err: any) {
    ElMessage.error(err?.message || '加载失败')
  } finally {
    adminLoading.value = false
  }
}

async function approve(row: CustomVoiceRequest) {
  try {
    await ElMessageBox.confirm(`确认通过「${row.name}」并生成音色吗？`, '提示', { type: 'warning' })
    await customVoiceApi.adminApprove(row.id)
    ElMessage.success('已通过')
    await voiceStore.fetchVoiceList()
    await refresh()
    await refreshAdmin()
  } catch {
    // ignore
  }
}

const rejectVisible = ref(false)
const rejecting = ref(false)
const rejectTargetId = ref<number>(0)
const rejectText = ref('')

function openReject(row: CustomVoiceRequest) {
  rejectTargetId.value = row.id
  rejectText.value = ''
  rejectVisible.value = true
}

async function reject() {
  if (!rejectTargetId.value) return
  rejecting.value = true
  try {
    await customVoiceApi.adminReject(rejectTargetId.value, rejectText.value)
    ElMessage.success('已拒绝')
    rejectVisible.value = false
    await voiceStore.fetchVoiceList()
    await refresh()
    await refreshAdmin()
  } catch (err: any) {
    ElMessage.error(err?.message || '操作失败')
  } finally {
    rejecting.value = false
  }
}

let autoTimer: number | null = null
function startAutoRefresh() {
  if (autoTimer !== null) return
  autoTimer = window.setInterval(() => {
    void refresh()
  }, 5 * 60 * 1000)
}

function stopAutoRefresh() {
  if (autoTimer !== null) {
    clearInterval(autoTimer)
    autoTimer = null
  }
}

onMounted(async () => {
  syncIsMobile()
  window.addEventListener('resize', syncIsMobile, { passive: true })

  await refresh()
  if (!userInfo.value) {
    await userStore.fetchUserInfo()
  }
  await voiceStore.fetchVoiceList()
  startAutoRefresh()
  if (isAdmin.value) {
    await refreshAdmin()
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', syncIsMobile)
  stopAutoRefresh()
  closePlayer()
})
</script>

<style scoped>
.cv-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f5f7fb;
  padding: 18px;
  overflow: auto;
}

.page-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 4px 2px 14px;
}

.head-left {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.title {
  font-size: 22px;
  font-weight: 900;
  color: #111827;
}

.sub {
  font-size: 13px;
  color: #7b7b7b;
  max-width: 720px;
}

.steps {
  display: grid;
  grid-template-columns: 1fr auto 1fr auto 1fr;
  gap: 14px;
  align-items: center;
  margin-bottom: 18px;
}

.step-card {
  background: #fff;
  border: 1px solid #e8eef7;
  border-radius: 14px;
  padding: 18px 18px;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  min-height: 132px;
}

.step-card .icon {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: rgba(47, 107, 255, 0.12);
  color: #2f6bff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 10px;
  font-size: 18px;
}

.step-title {
  font-weight: 900;
  text-align: center;
  color: #111827;
  margin-bottom: 8px;
}

.step-desc {
  font-size: 12px;
  color: #8c8c8c;
  text-align: center;
  line-height: 1.6;
}

.arrow {
  color: #cbd5e1;
  font-size: 28px;
  font-weight: 900;
  user-select: none;
}

.list-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2px 10px;
}

.h {
  font-size: 16px;
  font-weight: 900;
  color: #111827;
}

.lh-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.auto {
  font-size: 12px;
  color: #8c8c8c;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.auto .dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #9ca3af;
}

.list-card {
  background: #fff;
  border: 1px solid #e8eef7;
  border-radius: 14px;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.list-table {
  padding: 12px 14px 6px;
}

.th {
  display: grid;
  grid-template-columns: 2.1fr 1fr 1.2fr 0.8fr 0.7fr;
  gap: 10px;
  padding: 10px 10px;
  color: #64748b;
  font-size: 12px;
  font-weight: 800;
  border-bottom: 1px solid #eef0f5;
}

.tb .tr {
  display: grid;
  grid-template-columns: 2.1fr 1fr 1.2fr 0.8fr 0.7fr;
  gap: 10px;
  padding: 14px 10px;
  border-bottom: 1px solid #f1f5f9;
  cursor: pointer;
}

.tb .tr:last-child {
  border-bottom: none;
}

.name-cell {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.avatar {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: rgba(47, 107, 255, 0.12);
  color: #2f6bff;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nm {
  min-width: 0;
}

.n1 {
  font-weight: 900;
  color: #111827;
  line-height: 1.2;
}

.n2 {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 3px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 360px;
}

.prog-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.pct {
  font-size: 12px;
  font-weight: 800;
  color: #2f6bff;
}

.status-pill {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 999px;
  border: 1px solid transparent;
  font-weight: 800;
}

.status-pill.ok {
  color: #16a34a;
  background: rgba(34, 197, 94, 0.12);
  border-color: rgba(34, 197, 94, 0.2);
}

.status-pill.run {
  color: #2f6bff;
  background: rgba(47, 107, 255, 0.12);
  border-color: rgba(47, 107, 255, 0.2);
}

.status-pill.bad {
  color: #f56c6c;
  background: rgba(245, 108, 108, 0.12);
  border-color: rgba(245, 108, 108, 0.2);
}

.muted {
  color: #9ca3af;
  font-size: 12px;
}

.empty {
  padding: 24px 0 10px;
}

.pager {
  padding: 10px 14px 14px;
  display: flex;
  justify-content: center;
}

.metrics {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.m-card {
  background: #fff;
  border: 1px solid #e8eef7;
  border-radius: 14px;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  padding: 14px 16px;
}

.m-card .k {
  font-size: 12px;
  color: #6b7280;
  font-weight: 800;
  margin-bottom: 10px;
}

.m-card .v {
  font-size: 22px;
  font-weight: 900;
  color: #111827;
}

.player-meta {
  margin-bottom: 12px;
  padding: 10px 12px;
  background: #f5f7fa;
  border-radius: 12px;
}

.row {
  display: flex;
  gap: 12px;
  padding: 6px 0;
}

.k {
  width: 90px;
  color: #606266;
  font-size: 13px;
  flex-shrink: 0;
}

.v {
  flex: 1;
  color: #303133;
  font-size: 13px;
  word-break: break-all;
}

.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
}

.detail {
  padding: 6px 2px;
}

.block {
  margin-top: 12px;
  border: 1px solid #eef0f5;
  border-radius: 12px;
  padding: 12px;
  background: #fafafa;
}

.block-title {
  font-weight: 900;
  color: #111827;
  font-size: 13px;
  margin-bottom: 8px;
}

.block-body {
  color: #2b2b2b;
  line-height: 1.65;
  font-size: 13px;
}

.pre {
  white-space: pre-wrap;
}

.url {
  margin-bottom: 6px;
}

.detail-actions {
  margin-top: 14px;
}

.sub-hint {
  margin-top: 6px;
  font-size: 12px;
  color: #909399;
}

.full {
  width: 100%;
}

.admin-area {
  margin-top: 18px;
}

.admin-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2px 10px;
}

.admin-title {
  font-size: 14px;
  font-weight: 900;
  color: #111827;
}

.admin-card {
  background: #fff;
  border: 1px solid #e8eef7;
  border-radius: 14px;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  padding: 14px;
}

@media (max-width: 980px) {
  .cv-page {
    padding: 12px;
  }

  .page-head {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }

  .head-right {
    width: 100%;
    justify-content: flex-end;
  }

  .steps {
    grid-template-columns: 1fr;
  }

  .arrow {
    display: none;
  }

  .th,
  .tb .tr {
    grid-template-columns: 1.6fr 1fr;
    grid-auto-rows: auto;
  }

  .th .prog,
  .th .st,
  .th .act,
  .tb .tr .prog,
  .tb .tr .st,
  .tb .tr .act {
    display: none;
  }

  .metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>


