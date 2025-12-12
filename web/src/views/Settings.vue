<template>
  <div class="settings">
    <el-container>
      <el-header>
        <div class="header-content">
          <div class="header-left">
            <el-button @click="$router.push('/')">返回主菜单</el-button>
          </div>
          <h2>设置</h2>
        </div>
      </el-header>
      <el-main>
        <el-card>
          <el-form :model="form" label-width="150px" style="max-width: 600px">
            <el-form-item label="主题">
              <el-select v-model="form.theme" @change="handleThemeChange" style="width: 200px">
                <el-option label="浅色" value="light" />
                <el-option label="深色" value="dark" />
              </el-select>
            </el-form-item>
            <el-form-item label="语言">
              <el-select v-model="form.language" @change="handleLanguageChange" style="width: 200px">
                <el-option label="简体中文" value="zh-CN" />
                <el-option label="English" value="en" />
              </el-select>
            </el-form-item>
            <el-form-item label="字体大小">
              <el-input-number v-model="form.fontSize" :min="12" :max="20" style="width: 200px" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useConfigStore } from '@/stores/configStore'
import { ElMessage } from 'element-plus'

const router = useRouter()
const configStore = useConfigStore()

const form = ref({
  theme: 'light',
  language: 'zh-CN',
  fontSize: 14,
})

const handleThemeChange = (value: string) => {
  configStore.setTheme(value)
  ElMessage.success('主题已更新')
}

const handleLanguageChange = (value: string) => {
  configStore.setLanguage(value)
  ElMessage.success('语言已更新')
}

onMounted(() => {
  configStore.init()
  form.value.theme = configStore.theme
  form.value.language = configStore.language
  form.value.fontSize = configStore.fontSize
})
</script>

<style scoped>
.settings {
  min-height: 100vh;
  background: #f5f5f5;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-left {
  display: flex;
  gap: 10px;
}

h2 {
  margin: 0;
  color: #333;
}
</style>

