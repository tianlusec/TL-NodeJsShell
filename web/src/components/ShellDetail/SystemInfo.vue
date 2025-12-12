<template>
  <div class="system-info">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>系统信息</span>
          <div>
            <el-tag v-if="isCached" type="info" size="small" style="margin-right: 8px">缓存</el-tag>
            <el-button size="small" @click="fetchInfo(true)">刷新</el-button>
          </div>
        </div>
      </template>

      <div v-if="info" class="info-content">
        <div class="info-item">
          <span class="label">操作系统:</span>
          <span class="value">{{ info.platform }} {{ info.release }}</span>
        </div>
        <div class="info-item">
          <span class="label">架构:</span>
          <span class="value">{{ info.arch }}</span>
        </div>
        <div class="info-item">
          <span class="label">主机名:</span>
          <span class="value">{{ info.hostname }}</span>
        </div>
        <div class="info-item">
          <span class="label">类型:</span>
          <span class="value">{{ info.type }}</span>
        </div>
        <div class="info-item">
          <span class="label">运行时间:</span>
          <span class="value">{{ formatUptime(info.uptime) }}</span>
        </div>
        <div class="info-item">
          <span class="label">CPU核心数:</span>
          <span class="value">{{ info.cpus }}</span>
        </div>
        <div class="info-item">
          <span class="label">内存使用:</span>
          <div class="memory-info">
            <el-progress
              :percentage="memoryPercentage"
              :status="memoryPercentage > 90 ? 'exception' : 'success'"
            />
            <div class="memory-details">
              <span>已用: {{ formatBytes(info.totalmem - info.freemem) }}</span>
              <span>总计: {{ formatBytes(info.totalmem) }}</span>
              <span>可用: {{ formatBytes(info.freemem) }}</span>
            </div>
          </div>
        </div>
        <div class="info-item">
          <span class="label">用户信息:</span>
          <div class="user-info">
            <div>用户名: {{ info.userInfo?.username || 'N/A' }}</div>
            <div>UID: {{ info.userInfo?.uid || 'N/A' }}</div>
            <div>GID: {{ info.userInfo?.gid || 'N/A' }}</div>
            <div>Home: {{ info.userInfo?.homedir || 'N/A' }}</div>
            <div>Shell: {{ info.userInfo?.shell || 'N/A' }}</div>
          </div>
        </div>
        
        <div v-if="info.networkInterfaces && info.networkInterfaces.length > 0" class="info-item">
          <span class="label">网络接口:</span>
          <div class="network-interfaces">
            <div v-for="iface in info.networkInterfaces" :key="iface.interface" class="interface-item">
              <div class="interface-name">{{ iface.interface }}</div>
              <div class="interface-details">
                <span>IP: {{ iface.address }}</span>
                <span>类型: {{ iface.family }}</span>
                <span v-if="iface.netmask">掩码: {{ iface.netmask }}</span>
                <span v-if="iface.mac">MAC: {{ iface.mac }}</span>
                <el-tag v-if="iface.internal" type="info" size="small">内部</el-tag>
              </div>
            </div>
          </div>
        </div>
        
        <div v-if="info.envVars" class="info-item">
          <span class="label">环境变量:</span>
          <el-collapse>
            <el-collapse-item title="查看环境变量" :name="1">
              <div class="env-vars">
                <div v-for="(value, key) in info.envVars" :key="key" class="env-item">
                  <span class="env-key">{{ key }}:</span>
                  <span class="env-value">{{ value }}</span>
                </div>
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
        
        <div v-if="info.hosts" class="info-item">
          <span class="label">Hosts 文件:</span>
          <el-collapse>
            <el-collapse-item title="查看 Hosts 内容" :name="2">
              <pre class="hosts-content">{{ info.hosts }}</pre>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
      <div v-else class="no-info">
        暂无系统信息
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { shellApi } from '@/api/shell'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  shellId: number
}>()

const loading = ref(false)
const info = ref<any>(null)
const isCached = ref(false)

const memoryPercentage = computed(() => {
  if (!info.value || !info.value.totalmem) return 0
  return Math.round(((info.value.totalmem - info.value.freemem) / info.value.totalmem) * 100)
})

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatUptime = (seconds: number) => {
  if (!seconds) return '0秒'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  let result = ''
  if (days > 0) result += `${days}天 `
  if (hours > 0) result += `${hours}小时 `
  if (minutes > 0) result += `${minutes}分钟 `
  result += `${secs}秒`
  return result.trim()
}

const fetchInfo = async (forceRefresh = false) => {
  loading.value = true
  try {
    const params = forceRefresh ? { refresh: 'true' } : {}
    const response = await shellApi.getInfo(props.shellId, params)
    info.value = response.system_info || response
    isCached.value = response.cached === true
  } catch (error: any) {
    ElMessage.error(error.message || '获取系统信息失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchInfo()
})

defineExpose({
  fetchInfo,
})
</script>

<style scoped>
.system-info {
  height: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.system-info :deep(.el-card) {
  border: none;
  box-shadow: none;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.system-info :deep(.el-card__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
  background: #f8f9fa;
  border-radius: 12px 12px 0 0;
  flex-shrink: 0;
}

.system-info :deep(.el-card__body) {
  padding: 20px;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s ease;
}

.info-item:hover {
  background: #f0f2f5;
  border-color: #d4d7de;
}

.label {
  font-weight: 600;
  color: #1d2129;
  font-size: 14px;
  letter-spacing: 0.2px;
}

.value {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
}

.memory-info {
  width: 100%;
}

.memory-details {
  display: flex;
  gap: 16px;
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 13px;
  color: #606266;
}

.network-interfaces {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.interface-item {
  padding: 12px;
  background: #ffffff;
  border-radius: 6px;
  border-left: 3px solid #409eff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.interface-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transform: translateX(2px);
}

.interface-name {
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
  font-size: 14px;
}

.interface-details {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 12px;
  color: #606266;
}

.env-vars {
  max-height: 300px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
}

.env-item {
  display: flex;
  gap: 8px;
  font-size: 12px;
  padding: 4px;
}

.env-key {
  font-weight: bold;
  color: #303133;
  min-width: 150px;
}

.env-value {
  color: #606266;
  word-break: break-all;
}

.hosts-content {
  max-height: 300px;
  overflow-y: auto;
  padding: 12px;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
  font-family: 'Courier New', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  color: #303133;
}

.system-info :deep(.el-card__body)::-webkit-scrollbar {
  width: 6px;
}

.system-info :deep(.el-card__body)::-webkit-scrollbar-track {
  background: #f5f7fa;
  border-radius: 3px;
}

.system-info :deep(.el-card__body)::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 3px;
}

.system-info :deep(.el-card__body)::-webkit-scrollbar-thumb:hover {
  background: #a8abb2;
}

.no-info {
  text-align: center;
  color: #909399;
  padding: 40px;
}
</style>

