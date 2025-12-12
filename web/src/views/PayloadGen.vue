<template>
  <div class="payload-gen">
    <el-card>
      <template #header>
        <div class="header-actions">
          <span class="header-title">载荷生成</span>
          <el-button @click="$router.push('/')">返回主菜单</el-button>
        </div>
      </template>

      <el-form :model="form" label-width="120px" style="max-width: 800px">
        <el-form-item label="目标URL" required>
          <el-input v-model="form.url" placeholder="http://example.com" />
        </el-form-item>
        <el-form-item label="内存马类型">
          <el-radio-group v-model="shellType">
            <el-radio value="nextjs">Next.js</el-radio>
            <el-radio value="nodejs">Node.js</el-radio>
          </el-radio-group>
        </el-form-item>

        <template v-if="shellType === 'nextjs'">
          <el-form-item label="Shell路径" required>
            <el-input v-model="form.shell_path" placeholder="/api/shell" />
          </el-form-item>
        </template>

        <template v-if="shellType === 'nodejs'">
          <el-form-item label="密码" required>
            <el-input v-model="form.password" type="password" show-password />
          </el-form-item>
          <el-form-item label="编码类型">
            <el-select v-model="form.encode_type" style="width: 100%">
              <el-option label="Base64" value="base64" />
              <el-option label="XOR" value="xor" />
              <el-option label="AES" value="aes" />
            </el-select>
          </el-form-item>
          <el-form-item label="模板名称">
            <el-select v-model="form.template_name" style="width: 100%">
              <el-option
                v-for="template in templates"
                :key="template"
                :label="template"
                :value="template"
              />
            </el-select>
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" @click="handleInject" :loading="injecting">
            注入内存马
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { payloadApi } from '@/api/payload'
import { ElMessage } from 'element-plus'

const router = useRouter()
const shellType = ref<'nextjs' | 'nodejs'>('nodejs')
const templates = ref<string[]>([])
const injecting = ref(false)

const form = ref({
  url: '',
  password: '',
  encode_type: 'base64',
  template_name: '',
  shell_path: '/api/shell',
})

const fetchTemplates = async () => {
  try {
    templates.value = await payloadApi.getTemplates()
    if (templates.value.length > 0) {
      form.value.template_name = templates.value[0]
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取模板列表失败')
  }
}

const handleInject = async () => {
  if (!form.value.url) {
    ElMessage.warning('请输入目标URL')
    return
  }

  if (shellType.value === 'nextjs' && !form.value.shell_path) {
    ElMessage.warning('请输入Shell路径')
    return
  }

  if (shellType.value === 'nodejs' && !form.value.password) {
    ElMessage.warning('请输入密码')
    return
  }

  injecting.value = true
  try {
    const data: any = {
      url: form.value.url,
    }

    if (shellType.value === 'nextjs') {
      data.shell_path = form.value.shell_path
    } else {
      data.password = form.value.password
      data.encode_type = form.value.encode_type
      data.template_name = form.value.template_name
    }

    const result = await payloadApi.inject(data)
    ElMessage.success(result.message || '注入成功')
  } catch (error: any) {
    ElMessage.error(error.message || '注入失败')
  } finally {
    injecting.value = false
  }
}

onMounted(() => {
  fetchTemplates()
})
</script>

<style scoped>
.payload-gen {
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



