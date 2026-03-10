<template>
  <div class="feedback-page">
    <div class="page-head">
      <div class="head-left">
        <div class="title">意见反馈</div>
        <div class="sub">提交问题与建议，管理员回复后会在这里展示</div>
      </div>
      <div class="head-right">
        <el-input
          v-model="keyword"
          clearable
          placeholder="搜索反馈内容..."
          class="search"
          @clear="page = 1"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="openCreate">
          <el-icon><Plus /></el-icon>
          提交反馈
        </el-button>
      </div>
    </div>

    <div class="tabs-row">
      <el-tabs v-model="statusTab" class="filter-tabs">
        <el-tab-pane label="全部" name="all" />
        <el-tab-pane label="处理中" name="open" />
        <el-tab-pane label="已回复" name="closed" />
      </el-tabs>

      <div class="tabs-right">
        <el-button :loading="loading" text @click="refresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <div class="grid-wrap" v-loading="loading">
      <div v-if="!loading && pagedItems.length === 0" class="empty">
        <el-empty description="暂无反馈" :image-size="120" />
      </div>

      <div v-else class="grid">
        <div v-for="it in pagedItems" :key="it.id" class="fb-card" @click="openDetail(it)">
          <div class="card-top">
            <div class="cover">
              <el-icon class="msg"><ChatDotRound /></el-icon>
            </div>
            <div class="top-right">
              <span class="status-pill" :class="statusClass(it)">{{ statusText(it) }}</span>
            </div>
          </div>

          <div class="card-main">
            <div class="row1">
              <el-tag size="small" effect="plain">{{ it.category || '其他' }}</el-tag>
              <span class="id">#{{ it.id }}</span>
              <span class="time">
                <el-icon><Clock /></el-icon>
                {{ formatDateTime(it.createdAt) }}
              </span>
            </div>

            <div class="content" :title="it.content">{{ it.content }}</div>

            <div v-if="it.reply" class="reply">
              <div class="reply-title">回复</div>
              <div class="reply-content" :title="it.reply">{{ it.reply }}</div>
            </div>
          </div>

          <div class="card-actions" @click.stop>
            <el-button class="act-btn" @click="copyText(String(it.id), '已复制反馈ID')">
              <el-icon><DocumentCopy /></el-icon>
              复制ID
            </el-button>
            <el-button class="act-btn" @click="copyText(it.content, '已复制反馈内容')">
              <el-icon><DocumentCopy /></el-icon>
              复制内容
            </el-button>
            <el-button class="act-btn primary" @click="openDetail(it)">
              <el-icon><View /></el-icon>
              查看详情
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <div class="pager" v-if="!loading && filteredItems.length > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :page-sizes="[6, 9, 12, 18]"
        layout="prev, pager, next, jumper"
        :total="filteredItems.length"
      />
    </div>

    <el-drawer v-model="detailVisible" title="反馈详情" :size="drawerSize">
      <div v-if="current" class="detail">
        <div class="detail-row">
          <span class="k">状态</span>
          <span class="v">
            <span class="status-pill" :class="statusClass(current)">{{ statusText(current) }}</span>
          </span>
        </div>
        <div class="detail-row">
          <span class="k">分类</span>
          <span class="v">
            <el-tag size="small" effect="plain">{{ current.category || '其他' }}</el-tag>
          </span>
        </div>
        <div class="detail-row">
          <span class="k">反馈ID</span>
          <span class="v mono">
            {{ current.id }}
            <el-button text size="small" :icon="DocumentCopy" @click="copyText(String(current.id), '已复制反馈ID')" />
          </span>
        </div>
        <div class="detail-row">
          <span class="k">提交时间</span>
          <span class="v">{{ formatDateTime(current.createdAt) }}</span>
        </div>
        <div class="detail-row">
          <span class="k">联系方式</span>
          <span class="v">{{ current.contact || '-' }}</span>
        </div>

        <div class="detail-block">
          <div class="block-title">内容</div>
          <div class="block-body">{{ current.content }}</div>
        </div>

        <div class="detail-block" v-if="current.reply">
          <div class="block-title">回复</div>
          <div class="block-body">{{ current.reply }}</div>
        </div>
      </div>
    </el-drawer>

    <el-dialog v-model="createVisible" title="提交反馈" width="640px">
      <el-form label-width="88px">
        <el-form-item label="分类">
          <el-select v-model="form.category" class="full" placeholder="可选">
            <el-option label="Bug" value="Bug" />
            <el-option label="功能建议" value="Feature" />
            <el-option label="体验优化" value="UX" />
            <el-option label="其他" value="Other" />
          </el-select>
        </el-form-item>
        <el-form-item label="内容" required>
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="6"
            placeholder="请描述你遇到的问题或建议（可包含复现步骤、期望效果等）"
          />
        </el-form-item>
        <el-form-item label="联系方式">
          <el-input v-model="form.contact" placeholder="邮箱/微信/手机号（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">提交</el-button>
      </template>
    </el-dialog>

    <div v-if="isAdmin" class="admin-area">
      <div class="admin-head">
        <div class="admin-title">管理员：全部反馈</div>
        <el-button text :loading="adminLoading" @click="refreshAdmin">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>

      <div class="admin-card" v-loading="adminLoading">
        <el-table
          :data="adminItems"
          style="width: 100%"
          row-key="id"
          :header-cell-style="{ background: '#f6f8fc', color: '#6b7280' }"
        >
          <el-table-column prop="id" label="ID" width="90" />
          <el-table-column prop="username" label="用户" width="140" />
          <el-table-column prop="category" label="分类" width="120" />
          <el-table-column prop="content" label="内容" min-width="240">
            <template #default="{ row }">
              <span class="admin-content" :title="row.content">{{ row.content }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="120">
            <template #default="{ row }">
              <span class="status-pill" :class="String(row.status) === 'open' ? 'open' : 'closed'">
                {{ String(row.status) === 'open' ? '处理中' : '已回复' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="140" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openReply(row)">回复</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <el-dialog v-model="replyVisible" title="管理员回复" width="640px">
      <div v-if="replyTarget" class="reply-meta">
        <div class="meta-row">
          <span class="mk">反馈ID</span><span class="mv mono">#{{ replyTarget.id }}</span>
        </div>
        <div class="meta-row">
          <span class="mk">用户</span><span class="mv">{{ replyTarget.username || replyTarget.userId }}</span>
        </div>
        <div class="meta-row">
          <span class="mk">内容</span><span class="mv">{{ replyTarget.content }}</span>
        </div>
      </div>
      <el-input v-model="replyText" type="textarea" :rows="6" placeholder="请输入回复内容..." />
      <template #footer>
        <el-button @click="replyVisible = false">取消</el-button>
        <el-button type="primary" :loading="replying" @click="submitReply">提交回复</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import { ChatDotRound, Clock, DocumentCopy, Plus, Refresh, Search, View } from '@element-plus/icons-vue'
import { feedbackApi, type FeedbackItem } from '@/api/feedback'
import { useUserStore } from '@/store/user'

const loading = ref(false)
const submitting = ref(false)
const items = ref<FeedbackItem[]>([])

const statusTab = ref<'all' | 'open' | 'closed'>('all')
const keyword = ref('')

const page = ref(1)
const pageSize = ref(9)

const isMobile = ref(false)

function syncIsMobile() {
  isMobile.value = window.innerWidth < 980
}

onMounted(() => {
  syncIsMobile()
  window.addEventListener('resize', syncIsMobile, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('resize', syncIsMobile)
})

const drawerSize = computed(() => (isMobile.value ? '100%' : '520px'))

const createVisible = ref(false)
const form = ref<{ category: string; content: string; contact: string }>({
  category: 'Feature',
  content: '',
  contact: '',
})

const detailVisible = ref(false)
const current = ref<FeedbackItem | null>(null)

const filteredItems = computed(() => {
  const base = [...items.value].sort((a, b) => {
    const ta = new Date(a.createdAt).getTime()
    const tb = new Date(b.createdAt).getTime()
    if (Number.isNaN(ta) && Number.isNaN(tb)) return 0
    if (Number.isNaN(ta)) return 1
    if (Number.isNaN(tb)) return -1
    return tb - ta
  })

  const st = statusTab.value
  const byStatus = st === 'all' ? base : base.filter((i) => String(i.status) === st)

  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return byStatus
  return byStatus.filter((i) => {
    const t = `${i.id} ${i.category} ${i.content} ${i.reply || ''}`.toLowerCase()
    return t.includes(kw)
  })
})

const pagedItems = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredItems.value.slice(start, start + pageSize.value)
})

watch([statusTab, keyword], () => {
  page.value = 1
})

function openCreate() {
  form.value = { category: 'Feature', content: '', contact: '' }
  createVisible.value = true
}

function openDetail(it: FeedbackItem) {
  current.value = it
  detailVisible.value = true
}

function formatDateTime(raw: string): string {
  const d = new Date(raw)
  if (Number.isNaN(d.getTime())) return raw || '-'
  const pad = (n: number) => `${n}`.padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function statusText(it: FeedbackItem) {
  return String(it.status) === 'open' ? '处理中' : '已回复'
}

function statusClass(it: FeedbackItem) {
  return String(it.status) === 'open' ? 'open' : 'closed'
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

async function copyText(text: string, okMsg: string) {
  try {
    await writeToClipboard(text)
    ElMessage.success(okMsg)
  } catch {
    ElMessage.error('复制失败')
  }
}

async function refresh() {
  loading.value = true
  try {
    const res = await feedbackApi.listMy()
    items.value = res.list || []
  } catch (err: any) {
    ElMessage.error(err?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function submit() {
  if (!form.value.content.trim()) {
    ElMessage.warning('请填写反馈内容')
    return
  }
  submitting.value = true
  try {
    await feedbackApi.create({
      category: form.value.category,
      content: form.value.content,
      contact: form.value.contact,
    })
    ElMessage.success('提交成功')
    createVisible.value = false
    await refresh()
  } catch (err: any) {
    ElMessage.error(err?.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

const userStore = useUserStore()
const { userInfo } = storeToRefs(userStore)
const isAdmin = computed(() => (userInfo.value?.id || 0) === 1)

const adminLoading = ref(false)
const adminItems = ref<FeedbackItem[]>([])
const replyVisible = ref(false)
const replying = ref(false)
const replyTarget = ref<FeedbackItem | null>(null)
const replyText = ref('')

async function refreshAdmin() {
  adminLoading.value = true
  try {
    const res = await feedbackApi.adminList()
    adminItems.value = res.list || []
  } catch (err: any) {
    ElMessage.error(err?.message || '加载失败')
  } finally {
    adminLoading.value = false
  }
}

function openReply(row: FeedbackItem) {
  replyTarget.value = row
  replyText.value = row.reply || ''
  replyVisible.value = true
}

async function submitReply() {
  if (!replyTarget.value) return
  if (!replyText.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  replying.value = true
  try {
    await feedbackApi.adminReply(replyTarget.value.id, { reply: replyText.value, status: 'closed' })
    ElMessage.success('已回复')
    replyVisible.value = false
    await refresh()
    await refreshAdmin()
  } catch (err: any) {
    ElMessage.error(err?.message || '操作失败')
  } finally {
    replying.value = false
  }
}

onMounted(async () => {
  if (!userInfo.value) {
    await userStore.fetchUserInfo()
  }
  await refresh()
  if (isAdmin.value) {
    await refreshAdmin()
  }
})
</script>

<style scoped>
.feedback-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f5f7fb;
  padding: 18px 18px 18px;
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
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 16px;
}

.fb-card {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e8eef7;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  cursor: pointer;
  transition: transform 0.12s ease, box-shadow 0.12s ease;
}

.fb-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 24px rgba(0, 0, 0, 0.06);
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

.msg {
  font-size: 18px;
}

.status-pill {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 999px;
  border: 1px solid transparent;
  user-select: none;
}

.status-pill.open {
  color: #409eff;
  background: rgba(64, 158, 255, 0.12);
  border-color: rgba(64, 158, 255, 0.2);
}

.status-pill.closed {
  color: #67c23a;
  background: rgba(103, 194, 58, 0.12);
  border-color: rgba(103, 194, 58, 0.2);
}

.card-main {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.row1 {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.id {
  font-size: 12px;
  color: #909399;
}

.time {
  font-size: 12px;
  color: #8c8c8c;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.content {
  font-weight: 800;
  color: #1f1f1f;
  font-size: 14px;
  line-height: 1.35;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 56px;
}

.reply {
  border-top: 1px dashed #eef0f5;
  padding-top: 10px;
}

.reply-title {
  font-size: 12px;
  color: #8c8c8c;
  margin-bottom: 6px;
}

.reply-content {
  font-size: 13px;
  color: #2b2b2b;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-actions {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 10px;
}

.act-btn {
  justify-content: center;
  height: 34px;
  border-radius: 10px;
}

.act-btn.primary {
  color: #fff;
  background: #2f6bff;
  border-color: #2f6bff;
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

.admin-area {
  margin-top: 16px;
}

.admin-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 2px 10px;
}

.admin-title {
  font-size: 14px;
  font-weight: 900;
  color: #111827;
}

.admin-card {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e8eef7;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.admin-content {
  display: inline-block;
  max-width: 520px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.reply-meta {
  background: #f5f7fa;
  border-radius: 12px;
  padding: 10px 12px;
  margin-bottom: 10px;
}

.meta-row {
  display: flex;
  gap: 10px;
  padding: 4px 0;
}

.mk {
  width: 70px;
  color: #606266;
  font-size: 13px;
  flex-shrink: 0;
}

.mv {
  flex: 1;
  color: #303133;
  font-size: 13px;
  word-break: break-all;
}

.detail {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.detail-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.k {
  font-size: 12px;
  color: #8c8c8c;
  flex-shrink: 0;
}

.v {
  font-size: 12px;
  color: #1f1f1f;
  text-align: right;
  word-break: break-all;
}

.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
}

.detail-block {
  margin-top: 10px;
  border: 1px solid #eef0f5;
  border-radius: 12px;
  padding: 12px;
  background: #fafafa;
}

.block-title {
  font-weight: 800;
  color: #1f1f1f;
  font-size: 13px;
  margin-bottom: 8px;
}

.block-body {
  color: #2b2b2b;
  line-height: 1.65;
  white-space: pre-wrap;
  font-size: 13px;
}

.full {
  width: 100%;
}

@media (max-width: 980px) {
  .feedback-page {
    padding: 12px;
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
    grid-template-columns: 1fr;
  }

  .card-actions {
    grid-template-columns: 1fr 1fr;
  }
}
</style>


