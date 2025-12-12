<template>
  <div class="dashboard">
    <el-card class="header-card">
      <div class="header-content">
        <h1>NodeJsshell Manager</h1>
        <div>
          <el-button type="primary" @click="showAddDialog = true">
            <el-icon><Plus /></el-icon>
            添加Shell
          </el-button>
          <el-button @click="$router.push('/proxy')" style="margin-left: 10px">
            <el-icon><Setting /></el-icon>
            代理管理
          </el-button>
        </div>
      </div>
    </el-card>

    <div class="stats-grid">
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">总Shell数</div>
        </div>
      </el-card>
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-value online">{{ stats.online }}</div>
          <div class="stat-label">在线Shell</div>
        </div>
      </el-card>
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-value">{{ stats.offline }}</div>
          <div class="stat-label">离线Shell</div>
        </div>
      </el-card>
    </div>

    <el-card class="shell-manager-card">
      <template #header>
        <div class="card-header">
          <span>Shell管理</span>
        </div>
      </template>
      <el-table 
        :data="shellStore.shells" 
        style="width: 100%" 
        class="modern-shell-table"
        stripe
        :header-cell-style="{ background: '#f8f9fa', color: '#1d2129', fontWeight: 600, fontSize: '14px' }"
        max-height="calc(100vh - 280px)"
      >
        <el-table-column prop="id" label="ID" width="70" align="center">
          <template #default="{ row }">
            <span style="color: #909399; font-weight: 500;">#{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" width="160" show-overflow-tooltip>
          <template #default="{ row }">
            <div style="display: flex; align-items: center; gap: 8px;">
              <el-icon style="color: #409eff;"><Connection /></el-icon>
              <span style="font-weight: 500; color: #1d2129;">{{ row.name || '未命名' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="URL" min-width="250" show-overflow-tooltip>
          <template #default="{ row }">
            <span style="color: #606266; font-family: 'Courier New', monospace; font-size: 13px;">{{ row.url }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="方法" width="90" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="row.method === 'POST' ? 'success' : 'info'" effect="plain">
              {{ row.method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="130" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="row.status === 'online' ? 'success' : 'info'" 
              effect="dark"
              size="small"
              class="status-tag"
            >
              <el-icon v-if="row.status === 'online'" class="status-icon">
                <CircleCheck />
              </el-icon>
              <el-icon v-else class="status-icon">
                <CircleClose />
              </el-icon>
              <span class="status-text">{{ row.status }}</span>
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_active" label="最后活跃" width="160">
          <template #default="{ row }">
            <span style="color: #909399; font-size: 13px;">{{ row.last_active || '从未' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="380" fixed="right">
          <template #default="{ row }">
            <div style="display: flex; gap: 8px; flex-wrap: wrap;">
              <el-button type="primary" size="small" @click="enterShell(row.id)" :icon="Promotion">
                进入
              </el-button>
              <el-button type="warning" size="small" @click="editShell(row)" :icon="Edit">
                编辑
              </el-button>
              <el-button type="info" size="small" @click="testShell(row.id)" :icon="VideoPlay">
                测试
              </el-button>
              <el-button type="danger" size="small" @click="deleteShell(row.id)" :icon="Delete">
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加Shell对话框 -->
    <el-dialog v-model="showAddDialog" title="添加Shell" width="600px" class="modern-dialog">
      <el-form :model="newShell" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="newShell.name" />
        </el-form-item>
        <el-form-item label="URL" required>
          <el-input v-model="newShell.url" />
        </el-form-item>
        <el-form-item label="密码" required>
          <el-input v-model="newShell.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="编码类型">
          <el-select v-model="newShell.encode_type" style="width: 100%">
            <el-option label="无编码" value="" />
            <el-option label="Base64" value="base64" />
          </el-select>
        </el-form-item>
        <el-form-item label="请求方法">
          <el-select v-model="newShell.method" style="width: 100%">
            <el-option label="POST" value="POST" />
            <el-option label="GET" value="GET" />
          </el-select>
        </el-form-item>
        <el-form-item label="分组">
          <el-input v-model="newShell.group" />
        </el-form-item>
        <el-form-item label="代理">
          <el-select v-model="newShell.proxy_id" clearable placeholder="选择代理（可选）" style="width: 100%">
            <el-option
              v-for="proxy in proxies"
              :key="proxy.id"
              :label="`${proxy.name} (${proxy.host}:${proxy.port})`"
              :value="proxy.id"
              :disabled="!proxy.enabled"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAdd">确定</el-button>
      </template>
    </el-dialog>

    <!-- 编辑Shell对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑Shell" width="600px" class="modern-dialog">
      <el-form :model="editingShell" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="editingShell.name" />
        </el-form-item>
        <el-form-item label="URL" required>
          <el-input v-model="editingShell.url" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="editingShell.password" type="password" show-password placeholder="留空不修改" />
        </el-form-item>
        <el-form-item label="编码类型">
          <el-select v-model="editingShell.encode_type" style="width: 100%">
            <el-option label="无编码" value="" />
            <el-option label="Base64" value="base64" />
          </el-select>
        </el-form-item>
        <el-form-item label="请求方法">
          <el-select v-model="editingShell.method" style="width: 100%">
            <el-option label="POST" value="POST" />
            <el-option label="GET" value="GET" />
          </el-select>
        </el-form-item>
        <el-form-item label="分组">
          <el-input v-model="editingShell.group" />
        </el-form-item>
        <el-form-item label="代理">
          <el-select v-model="editingShell.proxy_id" clearable placeholder="选择代理（可选）" style="width: 100%">
            <el-option
              v-for="proxy in proxies"
              :key="proxy.id"
              :label="`${proxy.name} (${proxy.host}:${proxy.port})`"
              :value="proxy.id"
              :disabled="!proxy.enabled"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useShellStore } from '@/stores/shellStore'
import { useProxyStore } from '@/stores/proxyStore'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Setting, Connection, CircleCheck, CircleClose, Promotion, Edit as EditIcon, VideoPlay, Delete as DeleteIcon } from '@element-plus/icons-vue'
import request from '@/api/request'
import type { Shell, Proxy } from '@/types'

const router = useRouter()
const shellStore = useShellStore()
const proxyStore = useProxyStore()

const showAddDialog = ref(false)
const showEditDialog = ref(false)
const proxies = ref<Proxy[]>([])

const newShell = ref<Partial<Shell>>({
  name: '',
  url: '',
  password: '',
  encode_type: '',
  method: 'POST',
  group: '',
  proxy_id: undefined,
})

const editingShell = ref<Partial<Shell>>({})

const stats = computed(() => {
  const total = shellStore.shells.length
  const online = shellStore.shells.filter(s => s.status === 'online').length
  const offline = total - online
  return { total, online, offline }
})

const enterShell = (id: number) => {
  router.push(`/shell/${id}/terminal`)
}

const editShell = (shell: Shell) => {
  editingShell.value = { ...shell }
  showEditDialog.value = true
}

const testShell = async (id: number) => {
  try {
    const result = await shellStore.testShell(id)
    if (result.success) {
      ElMessage.success(`测试成功，延迟: ${result.latency}ms`)
    } else {
      ElMessage.error('测试失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '测试失败')
  }
}

const deleteShell = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个Shell吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await shellStore.deleteShell(id)
    ElMessage.success('删除成功')
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const fetchProxies = async () => {
  try {
    proxies.value = await request.get('/proxies')
  } catch (error: any) {
    console.error('获取代理列表失败:', error)
  }
}

const handleAdd = async () => {
  try {
    const shellData = { ...newShell.value }
    if (shellData.proxy_id === undefined || shellData.proxy_id === null) {
      delete shellData.proxy_id
    }
    await shellStore.createShell(shellData)
    ElMessage.success('添加成功')
    showAddDialog.value = false
    Object.assign(newShell.value, {
      name: '',
      url: '',
      password: '',
      encode_type: '',
      method: 'POST',
      group: '',
      proxy_id: undefined,
    })
  } catch (error: any) {
    ElMessage.error(error.message || '添加失败')
  }
}

const handleUpdate = async () => {
  try {
    if (!editingShell.value.id) return
    const shellData = { ...editingShell.value }
    if (shellData.proxy_id === undefined || shellData.proxy_id === null) {
      shellData.proxy_id = 0 as any
    }
    await shellStore.updateShell(editingShell.value.id, shellData)
    ElMessage.success('更新成功')
    showEditDialog.value = false
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  }
}

onMounted(() => {
  shellStore.fetchShells()
  fetchProxies()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8ecf1 100%);
  min-height: 100vh;
  height: 100vh;
  overflow: auto;
  display: flex;
  flex-direction: column;
}

.header-card {
  margin-bottom: 20px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h1 {
  margin: 0;
  color: #1d2129;
  font-weight: 600;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-content {
  padding: 16px 20px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #1d2129;
  margin-bottom: 6px;
}

.stat-value.online {
  color: #67c23a;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  font-weight: 500;
}

.shell-manager-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.shell-manager-card :deep(.el-card__body) {
  flex: 1;
  overflow: auto;
  padding: 20px;
}

.card-header {
  font-weight: 600;
  font-size: 16px;
  color: #1d2129;
}

.modern-shell-table {
  border-radius: 8px;
  overflow: hidden;
}

.modern-shell-table :deep(.el-table__row) {
  transition: all 0.2s ease;
}

.modern-shell-table :deep(.el-table__row:hover) {
  background-color: #f8f9fa;
  transform: translateX(2px);
}

.modern-shell-table :deep(.el-table td),
.modern-shell-table :deep(.el-table th) {
  border-bottom: 1px solid #f0f0f0;
}

.modern-shell-table :deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
  background: #fafbfc;
}

:deep(.status-tag) {
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  flex-wrap: nowrap !important;
  white-space: nowrap !important;
  gap: 4px !important;
  flex-direction: row !important;
}

:deep(.status-icon) {
  flex-shrink: 0 !important;
  display: inline-flex !important;
  align-items: center !important;
  font-size: 14px !important;
  line-height: 1 !important;
  margin: 0 !important;
}

:deep(.status-text) {
  white-space: nowrap !important;
  display: inline-block !important;
  line-height: 1 !important;
  margin: 0 !important;
}

.modern-dialog :deep(.el-dialog) {
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
}
</style>

