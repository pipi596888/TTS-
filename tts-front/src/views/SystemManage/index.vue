<template>
  <div class="admin-page">
    <div class="layout">
      <aside class="left">
        <div class="side-card">
          <div class="side-title">管理中心</div>
          <el-menu class="side-menu" :default-active="active" @select="handleSelect">
            <el-menu-item index="users">
              <el-icon><User /></el-icon>
              <span>用户管理</span>
            </el-menu-item>
            <el-menu-item index="roles">
              <el-icon><Avatar /></el-icon>
              <span>角色权限</span>
            </el-menu-item>
            <el-menu-item index="logs">
              <el-icon><Document /></el-icon>
              <span>系统日志</span>
            </el-menu-item>
            <el-menu-item index="monitor">
              <el-icon><DataAnalysis /></el-icon>
              <span>资源监控</span>
            </el-menu-item>
          </el-menu>
        </div>

        <div class="metric-card">
          <div class="metric-title">总用户数</div>
          <div class="metric-row">
            <div class="metric-value">{{ stats?.users ?? 0 }}</div>
            <div class="metric-trend up">+12%</div>
          </div>
        </div>

        <div class="metric-card">
          <div class="metric-title">活跃用户</div>
          <div class="metric-row">
            <div class="metric-value">{{ activeUsers }}</div>
            <div class="metric-trend up">+5%</div>
          </div>
        </div>
      </aside>

      <main class="right">
        <div v-if="!isAdmin" class="non-admin">
          <el-alert
            title="当前账号没有管理员权限（默认 userId=1 为管理员）"
            type="warning"
            show-icon
            :closable="false"
          />
        </div>

        <template v-else>
          <!-- 用户管理 -->
          <section v-show="active === 'users'">
            <div class="page-head">
              <div class="head-left">
                <div class="title">用户管理</div>
                <div class="sub">管理平台用户信息、角色分配及账户状态</div>
              </div>
              <div class="head-right">
                <el-input v-model="userKeyword" class="search" clearable placeholder="搜索功能或用户..." @clear="userPage = 1">
                  <template #prefix>
                    <el-icon><Search /></el-icon>
                  </template>
                </el-input>
                <el-button class="ghost" @click="exportUsers">
                  <el-icon><Download /></el-icon>
                  导出
                </el-button>
                <el-button type="primary" @click="openCreateUser">
                  <el-icon><Plus /></el-icon>
                  新增用户
                </el-button>
              </div>
            </div>

            <div class="panel">
              <div class="panel-top">
                <div class="bulk">
                  <el-checkbox v-model="bulkEnabled" label="批量操作" />
                </div>
                <div class="count">共 {{ usersTotal }} 条记录</div>
              </div>

              <el-table
                :data="users"
                v-loading="usersLoading"
                style="width: 100%"
                row-key="id"
                :header-cell-style="{ background: '#f6f8fc', color: '#6b7280' }"
              >
                <el-table-column type="selection" width="44" />
                <el-table-column prop="uid" label="用户 ID" width="140">
                  <template #default="{ row }">
                    <span class="uid">{{ row.uid }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="username" label="用户名" min-width="160">
                  <template #default="{ row }">
                    <div class="user-cell">
                      <div class="avatar">{{ row.username.slice(0, 1) }}</div>
                      <div class="uname">{{ row.username }}</div>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="role" label="角色" width="140">
                  <template #default="{ row }">
                    <span class="role-pill" :class="roleClass(row.role)">{{ roleLabel(row.role) }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="status" label="状态" width="140">
                  <template #default="{ row }">
                    <div class="status-cell">
                      <span class="dot" :class="row.status === 'active' ? 'ok' : 'off'"></span>
                      <span :class="row.status === 'active' ? 'st-ok' : 'st-off'">{{ statusLabel(row.status) }}</span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="160" fixed="right">
                  <template #default="{ row }">
                    <div class="ops">
                      <el-button link type="primary" :disabled="row.id === 1" @click="openEditUser(row)">
                        <el-icon><EditPen /></el-icon>
                        编辑
                      </el-button>
                      <el-button link type="danger" :disabled="row.id === 1" @click="deleteUser(row)">
                        <el-icon><Delete /></el-icon>
                        删除
                      </el-button>
                    </div>
                  </template>
                </el-table-column>
              </el-table>

              <div class="panel-bottom">
                <el-pagination
                  v-model:current-page="userPage"
                  v-model:page-size="userPageSize"
                  :page-sizes="[5, 10, 20]"
                  layout="prev, pager, next"
                  :total="usersTotal"
                />
              </div>
            </div>

            <div class="panel logs-panel">
              <div class="logs-head">
                <div class="logs-title">最近系统操作</div>
                <el-button link type="primary" @click="active = 'logs'">查看全部日志</el-button>
              </div>
              <div class="logs-body">
                <div v-if="recentLogs.length === 0" class="muted">暂无记录</div>
                <div v-for="op in recentLogs" :key="op.id" class="log-item">
                  <div class="log-main">
                    <span class="bullet"></span>
                    <span class="log-text">{{ logText(op) }}</span>
                  </div>
                  <div class="log-sub">
                    <span>{{ timeAgo(op.createdAt) }}</span>
                    <span class="sep">·</span>
                    <span class="muted">{{ op.ip }}</span>
                  </div>
                </div>
              </div>
            </div>
          </section>

          <!-- 角色权限 -->
          <section v-show="active === 'roles'" class="panel">
            <div class="panel-top">
              <div class="bulk"><span class="logs-title2">角色权限</span></div>
              <div class="count">共 {{ roles.length }} 个角色</div>
            </div>
            <el-table
              :data="roles"
              v-loading="rolesLoading"
              style="width: 100%"
              row-key="key"
              :header-cell-style="{ background: '#f6f8fc', color: '#6b7280' }"
            >
              <el-table-column prop="name" label="角色" width="160" />
              <el-table-column prop="description" label="说明" min-width="220" />
              <el-table-column label="权限" min-width="260">
                <template #default="{ row }">
                  <div class="perm-tags">
                    <el-tag v-for="p in row.permissions" :key="p" size="small" effect="plain" type="info">{{ p }}</el-tag>
                    <span v-if="!row.permissions || row.permissions.length === 0" class="muted">无</span>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </section>

          <!-- 系统日志 -->
          <section v-show="active === 'logs'" class="panel">
            <div class="panel-top logs-top">
              <div class="bulk"><span class="logs-title2">系统日志</span></div>
              <div class="count">共 {{ logsTotal }} 条</div>
            </div>
            <div class="logs-body full" v-loading="logsLoading">
              <div v-if="logs.length === 0" class="muted">暂无记录</div>
              <div v-for="op in logs" :key="op.id" class="log-item">
                <div class="log-main">
                  <span class="bullet"></span>
                  <span class="log-text">{{ logText(op) }}</span>
                </div>
                <div class="log-sub">
                  <span>{{ timeAgo(op.createdAt) }}</span>
                  <span class="sep">·</span>
                  <span class="muted">{{ op.ip }}</span>
                </div>
              </div>
            </div>
            <div class="panel-bottom">
              <el-pagination
                v-model:current-page="logsPage"
                v-model:page-size="logsPageSize"
                :page-sizes="[10, 20, 50]"
                layout="prev, pager, next"
                :total="logsTotal"
              />
            </div>
          </section>

          <!-- 资源监控 -->
          <section v-show="active === 'monitor'" class="panel">
            <div class="panel-top">
              <div class="bulk"><span class="logs-title2">资源监控</span></div>
              <div class="count">数据概览</div>
            </div>
            <div class="monitor-grid" v-loading="statsLoading">
              <div class="mon-card">
                <div class="k">任务总数</div>
                <div class="v">{{ stats?.tasks ?? '-' }}</div>
              </div>
              <div class="mon-card">
                <div class="k">处理中</div>
                <div class="v">{{ stats?.tasksProcessing ?? '-' }}</div>
              </div>
              <div class="mon-card">
                <div class="k">未处理反馈</div>
                <div class="v">{{ stats?.feedbackOpen ?? '-' }}</div>
              </div>
              <div class="mon-card">
                <div class="k">待审定制音色</div>
                <div class="v">{{ stats?.customVoicePending ?? '-' }}</div>
              </div>
            </div>
          </section>
        </template>
      </main>
    </div>

    <el-dialog v-model="userDialogVisible" :title="editingUser ? '编辑用户' : '新增用户'" width="520px">
      <el-form :model="userForm" label-width="90px">
        <el-form-item label="用户名" required>
          <el-input v-model="userForm.username" placeholder="如：张三" />
        </el-form-item>
        <el-form-item v-if="!editingUser" label="密码" required>
          <el-input v-model="userForm.password" type="password" show-password placeholder="设置登录密码" />
        </el-form-item>
        <el-form-item v-else label="重置密码">
          <el-input v-model="userForm.password" type="password" show-password placeholder="不修改请留空" />
        </el-form-item>
        <el-form-item label="角色" required>
          <el-select v-model="userForm.role" placeholder="请选择">
            <el-option label="管理员" value="admin" />
            <el-option label="工程师" value="engineer" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" required>
          <el-select v-model="userForm.status" placeholder="请选择">
            <el-option label="正常" value="active" />
            <el-option label="已禁用" value="disabled" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="userDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitUser">{{ editingUser ? '保存' : '创建' }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { saveAs } from 'file-saver'
import {
  Avatar,
  DataAnalysis,
  Delete,
  Document,
  Download,
  EditPen,
  Plus,
  Search,
  User,
} from '@element-plus/icons-vue'
import { systemApi, type SystemStats } from '@/api/system'
import { adminApi, type AdminLog, type AdminRole, type AdminUser, type AdminUserRoleKey, type AdminUserStatusKey } from '@/api/admin'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
const { userInfo } = storeToRefs(userStore)
const isAdmin = computed(() => (userInfo.value?.id || 0) === 1)

const active = ref<'users' | 'roles' | 'logs' | 'monitor'>('users')
function handleSelect(key: string) {
  active.value = key as any
}

const statsLoading = ref(false)
const stats = ref<SystemStats | null>(null)
const activeUsers = computed(() => {
  const u = stats.value?.users ?? 0
  const open = stats.value?.feedbackOpen ?? 0
  return Math.max(0, u - open)
})

async function loadStats() {
  statsLoading.value = true
  try {
    stats.value = await systemApi.stats()
  } catch (err: any) {
    ElMessage.error(err?.message || '加载失败')
  } finally {
    statsLoading.value = false
  }
}

const bulkEnabled = ref(false)
const userKeyword = ref('')
const userPage = ref(1)
const userPageSize = ref(5)
const usersLoading = ref(false)
const users = ref<AdminUser[]>([])
const usersTotal = ref(0)

function roleLabel(role: AdminUserRoleKey) {
  if (role === 'admin') return '管理员'
  if (role === 'engineer') return '工程师'
  return '普通用户'
}

function roleClass(role: AdminUserRoleKey) {
  if (role === 'admin') return 'admin'
  if (role === 'engineer') return 'eng'
  return 'user'
}

function statusLabel(status: AdminUserStatusKey) {
  return status === 'active' ? '正常' : '已禁用'
}

async function loadUsers() {
  usersLoading.value = true
  try {
    const res = await adminApi.listUsers({
      keyword: userKeyword.value.trim() || undefined,
      page: userPage.value,
      pageSize: userPageSize.value,
    })
    users.value = res.list || []
    usersTotal.value = res.total || 0
  } catch (err: any) {
    ElMessage.error(err?.message || '加载用户失败')
  } finally {
    usersLoading.value = false
  }
}

let kwTimer: number | null = null
watch(userKeyword, () => {
  userPage.value = 1
  if (kwTimer) window.clearTimeout(kwTimer)
  kwTimer = window.setTimeout(() => {
    void loadUsers()
  }, 250)
})
watch([userPage, userPageSize], () => {
  void loadUsers()
})

async function exportUsers() {
  const keyword = userKeyword.value.trim() || undefined
  const all: AdminUser[] = []
  const pageSize = 100
  for (let page = 1; page <= 10; page++) {
    const res = await adminApi.listUsers({ keyword, page, pageSize })
    all.push(...(res.list || []))
    if (all.length >= (res.total || 0)) break
  }

  const rows = all
  const csv = [
    ['uid', 'username', 'role', 'status'].join(','),
    ...rows.map((r) =>
      [r.uid, r.username, roleLabel(r.role), statusLabel(r.status)].map((x) => `"${String(x).replace(/"/g, '""')}"`).join(',')
    ),
  ].join('\n')
  saveAs(new Blob([csv], { type: 'text/csv;charset=utf-8' }), `users_${Date.now()}.csv`)
  await adminApi.appendLog('导出了用户列表数据')
  await loadLogs(true)
  ElMessage.success('已导出')
}

// Logs (real backend)
const logsLoading = ref(false)
const logs = ref<AdminLog[]>([])
const logsTotal = ref(0)
const logsPage = ref(1)
const logsPageSize = ref(10)
const recentLogs = computed(() => logs.value.slice(0, 2))

function logText(op: AdminLog) {
  const actor = op.actorUsername ? `管理员 ${op.actorUsername}` : '管理员'
  return `${actor} ${op.action}`
}

function timeAgo(createdAt: string) {
  const d = new Date(String(createdAt || '').replace(' ', 'T'))
  const t = d.getTime()
  if (!t) return createdAt || '-'

  const diff = Date.now() - t
  if (diff < 60 * 1000) return '刚刚'
  const min = Math.floor(diff / (60 * 1000))
  if (min < 60) return `${min} 分钟前`
  const hour = Math.floor(min / 60)
  if (hour < 24) return `${hour} 小时前`
  const day = Math.floor(hour / 24)
  return `${day} 天前`
}

async function loadLogs(resetPage = false) {
  if (resetPage) logsPage.value = 1
  logsLoading.value = true
  try {
    const res = await adminApi.listLogs({
      page: logsPage.value,
      pageSize: logsPageSize.value,
    })
    logs.value = res.list || []
    logsTotal.value = res.total || 0
  } catch (err: any) {
    ElMessage.error(err?.message || '加载日志失败')
  } finally {
    logsLoading.value = false
  }
}
watch([logsPage, logsPageSize], () => {
  void loadLogs()
})

// Roles (real backend)
const rolesLoading = ref(false)
const roles = ref<AdminRole[]>([])
async function loadRoles() {
  rolesLoading.value = true
  try {
    const res = await adminApi.listRoles()
    roles.value = res.list || []
  } catch (err: any) {
    ElMessage.error(err?.message || '加载角色失败')
  } finally {
    rolesLoading.value = false
  }
}

// Create / Edit user dialog (local only)
const userDialogVisible = ref(false)
const editingUser = ref<AdminUser | null>(null)
const userForm = ref<{ username: string; password: string; role: AdminUserRoleKey; status: AdminUserStatusKey }>({
  username: '',
  password: '',
  role: 'user',
  status: 'active',
})

function openCreateUser() {
  editingUser.value = null
  userForm.value = { username: '', password: '', role: 'user', status: 'active' }
  userDialogVisible.value = true
}

function openEditUser(u: AdminUser) {
  editingUser.value = u
  userForm.value = { username: u.username, password: '', role: u.role, status: u.status }
  userDialogVisible.value = true
}

async function submitUser() {
  const username = userForm.value.username.trim()
  const password = userForm.value.password
  if (!username) {
    ElMessage.warning('请输入用户名')
    return
  }
  if (!editingUser.value && !password) {
    ElMessage.warning('请设置密码')
    return
  }

  try {
    if (editingUser.value) {
      await adminApi.updateUser(editingUser.value.id, {
        username,
        role: userForm.value.role,
        status: userForm.value.status,
        password: password || undefined,
      })
      ElMessage.success('已保存')
    } else {
      await adminApi.createUser({
        username,
        password,
        role: userForm.value.role,
        status: userForm.value.status,
      })
      ElMessage.success('已创建')
    }
    userDialogVisible.value = false
    await loadUsers()
    await loadLogs(true)
    await loadStats()
  } catch (err: any) {
    ElMessage.error(err?.message || '操作失败')
  }
}

async function deleteUser(u: AdminUser) {
  try {
    await ElMessageBox.confirm(`确定要删除用户 ${u.uid} 吗？`, '提示', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await adminApi.deleteUser(u.id)
    ElMessage.success('已删除')
    await loadUsers()
    await loadLogs(true)
    await loadStats()
  } catch {
    // cancelled
  }
}

onMounted(async () => {
  if (!userInfo.value) {
    await userStore.fetchUserInfo()
  }
  if (isAdmin.value) {
    await loadStats()
    await loadUsers()
    await loadLogs(true)
    await loadRoles()
  }
})
</script>

<style scoped>
.admin-page {
  height: 100%;
  background: #f5f7fb;
  overflow: hidden;
}

.layout {
  height: 100%;
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 18px;
  padding: 18px;
  overflow: hidden;
}

.left {
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-height: 0;
}

.side-card {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e8eef7;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  padding: 14px;
}

.side-title {
  font-weight: 800;
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 10px;
}

.side-menu {
  border-right: none;
}

.side-menu :deep(.el-menu-item) {
  height: 42px;
  border-radius: 10px;
  margin: 6px 0;
}

.side-menu :deep(.el-menu-item.is-active) {
  background: rgba(47, 107, 255, 0.1);
  color: #2f6bff;
}

.metric-card {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e8eef7;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  padding: 14px;
}

.metric-title {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 10px;
  font-weight: 700;
}

.metric-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.metric-value {
  font-size: 26px;
  font-weight: 900;
  color: #111827;
}

.metric-trend {
  font-size: 12px;
  font-weight: 800;
  padding: 2px 8px;
  border-radius: 999px;
}

.metric-trend.up {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.12);
}

.right {
  min-height: 0;
  overflow: auto;
  padding-right: 2px;
}

.non-admin {
  padding: 8px 2px 0;
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
  font-size: 20px;
  font-weight: 900;
  color: #111827;
}

.sub {
  font-size: 13px;
  color: #7b7b7b;
}

.head-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.search {
  width: 320px;
}

.ghost {
  background: #fff;
  border: 1px solid #e8eef7;
}

.panel {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #e8eef7;
  box-shadow: 0 1px 10px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.panel-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-bottom: 1px solid #eef0f5;
  background: #fbfcff;
}

.bulk {
  font-weight: 800;
  color: #374151;
  font-size: 13px;
}

.count {
  font-size: 12px;
  color: #9ca3af;
}

.uid {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
  color: #64748b;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: rgba(47, 107, 255, 0.12);
  color: #2f6bff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 900;
  font-size: 12px;
}

.uname {
  font-weight: 800;
  color: #111827;
}

.role-pill {
  display: inline-flex;
  align-items: center;
  height: 22px;
  padding: 0 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 800;
  border: 1px solid transparent;
}

.role-pill.admin {
  color: #a855f7;
  background: rgba(168, 85, 247, 0.12);
  border-color: rgba(168, 85, 247, 0.22);
}

.role-pill.eng {
  color: #2f6bff;
  background: rgba(47, 107, 255, 0.12);
  border-color: rgba(47, 107, 255, 0.22);
}

.role-pill.user {
  color: #6b7280;
  background: rgba(107, 114, 128, 0.12);
  border-color: rgba(107, 114, 128, 0.22);
}

.status-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 800;
}

.dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
}

.dot.ok {
  background: #22c55e;
}

.dot.off {
  background: #9ca3af;
}

.st-ok {
  color: #16a34a;
}

.st-off {
  color: #6b7280;
}

.ops {
  display: flex;
  gap: 10px;
}

.perm-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
}

.panel-bottom {
  padding: 12px 14px;
  display: flex;
  justify-content: center;
}

.logs-panel {
  margin-top: 16px;
  padding: 14px;
}

.logs-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.logs-title {
  font-weight: 900;
  color: #111827;
}

.logs-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.logs-body.full {
  padding: 14px;
}

.log-item {
  padding-left: 8px;
}

.log-main {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #111827;
  font-weight: 700;
}

.bullet {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #2f6bff;
}

.log-sub {
  margin-left: 14px;
  margin-top: 4px;
  font-size: 12px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  gap: 6px;
}

.sep {
  opacity: 0.7;
}

.muted {
  color: #9ca3af;
}

.placeholder {
  padding: 26px;
}

.ph-title {
  font-size: 16px;
  font-weight: 900;
  color: #111827;
  margin-bottom: 8px;
}

.ph-sub {
  color: #6b7280;
  font-size: 13px;
}

.logs-title2 {
  font-weight: 900;
  color: #111827;
}

.monitor-grid {
  padding: 14px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.mon-card {
  border: 1px solid #eef0f5;
  border-radius: 14px;
  padding: 14px;
  background: #fff;
}

.mon-card .k {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 8px;
  font-weight: 800;
}

.mon-card .v {
  font-size: 22px;
  font-weight: 900;
  color: #111827;
}

@media (max-width: 980px) {
  .admin-page {
    padding: 12px;
  }

  .layout {
    grid-template-columns: 1fr;
  }

  .left {
    width: 100%;
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

  .monitor-grid {
    grid-template-columns: 1fr;
  }
}
</style>

