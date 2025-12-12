<template>
  <div class="proxy-manager">
    <el-card>
      <template #header>
        <div class="header-actions">
          <span class="header-title">代理管理</span>
          <div>
            <el-button type="primary" @click="showAddDialog = true">添加代理</el-button>
            <el-button @click="$router.push('/')">返回主菜单</el-button>
          </div>
        </div>
      </template>

      <el-table :data="proxies" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="host" label="主机" width="150" />
        <el-table-column prop="port" label="端口" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="enabled" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="warning" size="small" @click="editProxy(row)">编辑</el-button>
            <el-button type="info" size="small" @click="testProxy(row.id)">测试</el-button>
            <el-button type="danger" size="small" @click="deleteProxy(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showAddDialog" title="添加代理" width="600px">
      <el-form :model="newProxy" label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="newProxy.name" />
        </el-form-item>
        <el-form-item label="类型" required>
          <el-select v-model="newProxy.type" style="width: 100%">
            <el-option label="HTTP" value="http" />
            <el-option label="SOCKS5" value="socks5" />
          </el-select>
        </el-form-item>
        <el-form-item label="主机" required>
          <el-input v-model="newProxy.host" />
        </el-form-item>
        <el-form-item label="端口" required>
          <el-input-number v-model="newProxy.port" :min="1" :max="65535" style="width: 100%" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="newProxy.username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="newProxy.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="newProxy.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAdd">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showEditDialog" title="编辑代理" width="600px">
      <el-form :model="editingProxy" label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="editingProxy.name" />
        </el-form-item>
        <el-form-item label="类型" required>
          <el-select v-model="editingProxy.type" style="width: 100%">
            <el-option label="HTTP" value="http" />
            <el-option label="SOCKS5" value="socks5" />
          </el-select>
        </el-form-item>
        <el-form-item label="主机" required>
          <el-input v-model="editingProxy.host" />
        </el-form-item>
        <el-form-item label="端口" required>
          <el-input-number v-model="editingProxy.port" :min="1" :max="65535" style="width: 100%" />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="editingProxy.username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="editingProxy.password" type="password" show-password placeholder="留空不修改" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="editingProxy.enabled" />
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
import request from '@/api/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Proxy } from '@/types'

const router = useRouter()
const proxies = ref<Proxy[]>([])

const showAddDialog = ref(false)
const showEditDialog = ref(false)

const newProxy = ref<Partial<Proxy>>({
  name: '',
  type: 'http',
  host: '',
  port: 8080,
  username: '',
  password: '',
  enabled: true,
})

const editingProxy = ref<Partial<Proxy>>({})

const fetchProxies = async () => {
  try {
    proxies.value = await request.get('/proxies')
  } catch (error: any) {
    ElMessage.error(error.message || '获取代理列表失败')
  }
}

const editProxy = (proxy: Proxy) => {
  editingProxy.value = { ...proxy }
  showEditDialog.value = true
}

const testProxy = async (id: number) => {
  try {
    const result = await request.post(`/proxies/${id}/test`)
    if (result.success) {
      ElMessage.success('代理测试成功')
    } else {
      ElMessage.error('代理测试失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '代理测试失败')
  }
}

const deleteProxy = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个代理吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await request.delete(`/proxies/${id}`)
    ElMessage.success('删除成功')
    await fetchProxies()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleAdd = async () => {
  try {
    await request.post('/proxies', newProxy.value)
    ElMessage.success('添加成功')
    showAddDialog.value = false
    await fetchProxies()
    Object.assign(newProxy.value, {
      name: '',
      type: 'http',
      host: '',
      port: 8080,
      username: '',
      password: '',
      enabled: true,
    })
  } catch (error: any) {
    ElMessage.error(error.message || '添加失败')
  }
}

const handleUpdate = async () => {
  try {
    if (!editingProxy.value.id) return
    await request.put(`/proxies/${editingProxy.value.id}`, editingProxy.value)
    ElMessage.success('更新成功')
    showEditDialog.value = false
    await fetchProxies()
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  }
}

onMounted(() => {
  fetchProxies()
})
</script>

<style scoped>
.proxy-manager {
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



