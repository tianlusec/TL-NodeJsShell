<template>
  <div class="shell-manager">
    <el-card>
      <template #header>
        <div class="header-actions">
          <span class="header-title">Shell管理</span>
          <div>
            <el-button type="primary" @click="showAddDialog = true">添加Shell</el-button>
            <el-button @click="$router.push('/')">返回主菜单</el-button>
          </div>
        </div>
      </template>

      <el-table :data="shellStore.shells" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="url" label="URL" width="300" show-overflow-tooltip />
        <el-table-column prop="method" label="方法" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'online' ? 'success' : 'info'">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_active" label="最后活跃" width="180" />
        <el-table-column label="操作" width="340" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="enterShell(row.id)">
              进入
            </el-button>
            <el-button type="warning" size="small" @click="editShell(row)">
              编辑
            </el-button>
            <el-button type="info" size="small" @click="testShell(row.id)">
              测试
            </el-button>
            <el-button type="danger" size="small" @click="deleteShell(row.id)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showAddDialog" title="添加Shell" width="600px">
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

    <el-dialog v-model="showEditDialog" title="编辑Shell" width="600px">
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useShellStore } from '@/stores/shellStore'
import { useProxyStore } from '@/stores/proxyStore'
import request from '@/api/request'
import { ElMessage, ElMessageBox } from 'element-plus'
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
.shell-manager {
  padding: 20px;
  background: #f5f5f5;
  min-height: 100vh;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 18px;
  font-weight: bold;
}
</style>

