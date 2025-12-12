<template>
  <div class="shell-detail">
    <div class="header-bar">
      <div class="header-left">
        <el-button @click="$router.push('/')" plain>
          <el-icon><ArrowLeft /></el-icon>
          返回主菜单
        </el-button>
      </div>
      <div class="header-title">
        <h2>{{ shell?.name || `Shell #${shellId}` }}</h2>
        <el-tag 
          :type="shell?.status === 'online' ? 'success' : 'info'" 
          effect="dark"
          size="large"
          style="font-weight: 500;"
        >
          {{ shell?.status || 'unknown' }}
        </el-tag>
      </div>
    </div>

    <div class="content-wrapper">
      <!-- 顶部导航卡片 -->
      <div class="nav-cards">
        <div 
          class="nav-card" 
          :class="{ active: activeTab === 'system' }"
          @click="activeTab = 'system'"
        >
          <el-icon class="nav-icon"><Monitor /></el-icon>
          <span class="nav-title">系统信息</span>
        </div>
        <div 
          class="nav-card" 
          :class="{ active: activeTab === 'terminal' }"
          @click="activeTab = 'terminal'"
        >
          <el-icon class="nav-icon"><Operation /></el-icon>
          <span class="nav-title">虚拟终端</span>
        </div>
        <div 
          class="nav-card" 
          :class="{ active: activeTab === 'files' }"
          @click="activeTab = 'files'"
        >
          <el-icon class="nav-icon"><FolderOpened /></el-icon>
          <span class="nav-title">文件管理</span>
        </div>
      </div>

      <!-- 内容区域 -->
      <div class="content-area">
        <!-- 系统信息 -->
        <div v-show="activeTab === 'system'" class="tab-content">
          <el-card class="content-card">
            <SystemInfo :shell-id="shellId" ref="systemInfoRef" />
          </el-card>
        </div>

        <!-- 虚拟终端 -->
        <div v-show="activeTab === 'terminal'" class="tab-content">
          <el-card class="content-card">
            <div class="terminal-wrapper">
              <div ref="terminalRef" class="terminal-container"></div>
            </div>
          </el-card>
        </div>

        <!-- 文件管理 -->
        <div v-show="activeTab === 'files'" class="tab-content">
          <el-card class="content-card">
            <div class="file-manager">
              <div class="file-toolbar">
                <el-button type="primary" @click="refreshFiles" :icon="Refresh">
                  刷新
                </el-button>
                <el-button @click="goUp" :icon="ArrowUp">
                  上级目录
                </el-button>
                <el-button type="success" @click="showUploadDialog" :icon="UploadFilled">
                  上传
                </el-button>
                <el-input
                  v-model="currentPath"
                  style="flex: 1; max-width: 400px;"
                  @keyup.enter="changePath"
                  placeholder="输入路径后按 Enter 跳转"
                >
                  <template #prefix>
                    <el-icon><Folder /></el-icon>
                  </template>
                  <template #append>
                    <el-button @click="changePath">跳转</el-button>
                  </template>
                </el-input>
              </div>
              <el-table
                :data="files"
                style="width: 100%"
                @row-dblclick="handleFileDoubleClick"
                @row-contextmenu="handleContextMenu"
                v-loading="filesLoading"
                :header-cell-style="{ background: '#f8f9fa', color: '#1d2129', fontWeight: 600 }"
                :row-style="{ cursor: 'pointer' }"
                stripe
                class="modern-table"
                height="calc(100vh - 380px)"
              >
                <el-table-column prop="name" label="名称" min-width="200">
                  <template #default="{ row }">
                    <div style="display: flex; align-items: center; gap: 8px;">
                      <el-icon v-if="row.type === 'd'" style="color: #409eff;">
                        <FolderOpened />
                      </el-icon>
                      <el-icon v-else style="color: #909399;">
                        <Document />
                      </el-icon>
                      <span>{{ row.name }}</span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="size" label="大小" width="120" align="right">
                  <template #default="{ row }">
                    <span style="color: #909399; font-family: 'Courier New', monospace;">
                      {{ row.size }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column prop="time" label="时间" width="180">
                  <template #default="{ row }">
                    <span style="color: #606266;">{{ row.time }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="类型" width="90" align="center">
                  <template #default="{ row }">
                    <el-tag 
                      :type="row.type === 'd' ? 'primary' : 'info'"
                      size="small"
                      effect="plain"
                    >
                      {{ row.type === 'd' ? '目录' : '文件' }}
                    </el-tag>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>
        </div>
      </div>
    </div>

    <el-dialog 
      v-model="fileDialogVisible" 
      title="文件内容" 
      width="85%" 
      top="3vh"
      class="modern-dialog"
      destroy-on-close
    >
      <div style="margin-bottom: 12px; font-size: 13px; color: #909399;">
        <el-icon><Document /></el-icon>
        <span style="margin-left: 6px;">{{ editingFilePath }}</span>
      </div>
      <el-input
        v-model="fileContent"
        type="textarea"
        :rows="24"
        readonly
        v-if="!fileEditing"
        style="font-family: 'Courier New', 'Consolas', monospace; font-size: 13px;"
      />
      <el-input
        v-model="fileContent"
        type="textarea"
        :rows="24"
        v-else
        style="font-family: 'Courier New', 'Consolas', monospace; font-size: 13px;"
      />
      <template #footer>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <el-button @click="fileDialogVisible = false" plain>关闭</el-button>
          <div style="display: flex; gap: 10px;">
            <el-button v-if="!fileEditing" type="primary" @click="fileEditing = true">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button v-if="fileEditing" type="success" @click="saveFile">
              <el-icon><Check /></el-icon>
              保存
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>

    <el-dialog 
      v-model="uploadDialogVisible" 
      title="上传文件" 
      width="580px"
      class="modern-dialog"
      destroy-on-close
    >
      <el-form :model="uploadForm" label-width="100px" label-position="top">
        <el-form-item label="文件路径">
          <el-input 
            v-model="uploadForm.path" 
            placeholder="输入文件路径，例如：/path/to/file.txt"
            size="large"
          >
            <template #prefix>
              <el-icon><Folder /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="选择文件">
          <el-upload
            :auto-upload="false"
            :on-change="handleFileSelect"
            :show-file-list="false"
            drag
            class="modern-upload"
          >
            <el-icon class="el-icon--upload" style="font-size: 48px; color: #409eff; margin-bottom: 12px;">
              <UploadFilled />
            </el-icon>
            <div class="el-upload__text">
              <div style="font-size: 16px; color: #303133; margin-bottom: 4px;">将文件拖到此处</div>
              <div style="font-size: 14px; color: #909399;">或<em style="color: #409eff; font-style: normal;">点击选择文件</em></div>
            </div>
          </el-upload>
        </el-form-item>
        <el-form-item v-if="uploadForm.fileName">
          <el-alert
            :title="`已选择文件: ${uploadForm.fileName}`"
            type="info"
            :closable="false"
            show-icon
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 12px;">
          <el-button @click="uploadDialogVisible = false" :disabled="uploading">取消</el-button>
          <el-button type="primary" @click="uploadFile" :loading="uploading" :disabled="!uploadForm.file">
            <el-icon v-if="!uploading"><UploadFilled /></el-icon>
            上传
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 上传进度对话框 -->
    <el-dialog 
      v-model="uploadProgress.visible" 
      title="上传进度" 
      width="480px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="!uploading"
      class="modern-dialog progress-dialog"
      align-center
    >
      <div style="padding: 24px 0;">
        <div style="margin-bottom: 20px; font-size: 15px; color: #1d2129; font-weight: 500; text-align: center;">
          <el-icon v-if="uploadProgress.total > 0" style="margin-right: 8px; color: #409eff;">
            <UploadFilled />
          </el-icon>
          <span v-if="uploadProgress.total > 0">正在上传: {{ uploadProgress.current }}/{{ uploadProgress.total }} 分片</span>
          <span v-else>准备上传...</span>
        </div>
        <el-progress 
          :percentage="uploadProgress.percentage" 
          :status="uploadProgress.status"
          :stroke-width="26"
          striped
          striped-flow
        />
        <div style="margin-top: 16px; font-size: 16px; color: #409eff; text-align: center; font-weight: 600;">
          {{ uploadProgress.percentage }}%
        </div>
      </div>
      <template #footer>
        <el-button 
          type="primary" 
          @click="uploadProgress.visible = false"
          :disabled="uploading"
          size="large"
          style="width: 100%;"
        >
          {{ uploading ? '上传中...' : '关闭' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 下载进度对话框 -->
    <el-dialog 
      v-model="downloadProgress.visible" 
      title="下载进度" 
      width="480px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
      class="modern-dialog progress-dialog"
      align-center
    >
      <div style="padding: 24px 0;">
        <div style="margin-bottom: 20px; font-size: 15px; color: #1d2129; font-weight: 500; text-align: center;">
          <el-icon v-if="downloadProgress.total > 0" style="margin-right: 8px; color: #67c23a;">
            <Download />
          </el-icon>
          <span v-if="downloadProgress.total > 0">正在下载: {{ downloadProgress.current }}/{{ downloadProgress.total }} 分片</span>
          <span v-else>准备下载...</span>
        </div>
        <el-progress 
          :percentage="downloadProgress.percentage" 
          :status="downloadProgress.status"
          :stroke-width="26"
          striped
          striped-flow
        />
        <div style="margin-top: 16px; font-size: 16px; color: #67c23a; text-align: center; font-weight: 600;">
          {{ downloadProgress.percentage }}%
        </div>
      </div>
      <template #footer>
        <el-button 
          type="success" 
          @click="downloadProgress.visible = false"
          size="large"
          style="width: 100%;"
        >
          关闭
        </el-button>
      </template>
    </el-dialog>

    <!-- 右键菜单 -->
    <Teleport to="body">
      <div
        v-if="contextMenu.visible"
        class="context-menu"
        :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
        @click.stop
      >
      <div
        v-if="contextMenu.item && contextMenu.item.type === 'd'"
        class="context-menu-item"
        @click="handleMenuEnter"
      >
        <el-icon><FolderOpened /></el-icon>
        <span>进入</span>
      </div>
      <div
        v-if="contextMenu.item && contextMenu.item.type !== 'd'"
        class="context-menu-item"
        @click="handleMenuView"
      >
        <el-icon><View /></el-icon>
        <span>查看</span>
      </div>
      <div
        v-if="contextMenu.item && contextMenu.item.type !== 'd'"
        class="context-menu-item"
        @click="handleMenuDownload"
      >
        <el-icon><Download /></el-icon>
        <span>下载</span>
      </div>
      <div
        v-if="contextMenu.item && contextMenu.item.type !== 'd'"
        class="context-menu-item"
        @click="handleMenuEdit"
      >
        <el-icon><Edit /></el-icon>
        <span>编辑</span>
      </div>
      <div
        v-if="contextMenu.item"
        class="context-menu-item context-menu-item-danger"
        @click="handleMenuDelete"
      >
        <el-icon><Delete /></el-icon>
        <span>删除</span>
      </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import { shellApi } from '@/api/shell'
import { fileApi } from '@/api/file'
import { useTerminalStore } from '@/stores/terminalStore'
import SystemInfo from '@/components/ShellDetail/SystemInfo.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled, FolderOpened, View, Download, Edit, Delete, Document, Operation, Check, Refresh, ArrowUp, Folder, ArrowLeft, Monitor } from '@element-plus/icons-vue'
import type { Shell, FileItem } from '@/types'

const route = useRoute()
const router = useRouter()
const shellId = parseInt(route.params.shellId as string)
const activeTab = ref<'system' | 'terminal' | 'files'>('terminal')

const shell = ref<Shell | null>(null)
const executing = ref(false)
const terminalRef = ref<HTMLElement>()
const terminal = ref<Terminal | null>(null)
const fitAddon = ref<FitAddon | null>(null)
const terminalStore = useTerminalStore()
const currentLine = ref('')
const prompt = ref('$ ')

const files = ref<FileItem[]>([])
const currentPath = ref('')
const filesLoading = ref(false)
const fileDialogVisible = ref(false)
const fileContent = ref('')
const editingFilePath = ref('')
const fileEditing = ref(false)
const systemInfoRef = ref<InstanceType<typeof SystemInfo>>()

const uploadDialogVisible = ref(false)
const uploading = ref(false)
const uploadForm = ref({
  path: '',
  content: '',
  file: null as File | null,
  fileName: ''
})
const uploadProgress = ref({
  current: 0,
  total: 0,
  percentage: 0,
  status: 'success' as 'success' | 'exception' | 'warning',
  visible: false
})
const downloadProgress = ref({
  current: 0,
  total: 0,
  percentage: 0,
  status: 'success' as 'success' | 'exception' | 'warning',
  visible: false
})
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  item: null as FileItem | null
})

const loadShell = async () => {
  try {
    shell.value = await shellApi.get(shellId)
  } catch (error: any) {
    ElMessage.error(error.message || '加载Shell信息失败')
  }
}

const initTerminal = async () => {
  if (!terminalRef.value) return
  await nextTick()

  const term = new Terminal({
    cursorBlink: true,
    theme: {
      background: '#1e1e1e',
      foreground: '#d4d4d4',
    },
    fontSize: 14,
    fontFamily: 'Consolas, "Courier New", monospace',
  })

  const addon = new FitAddon()
  term.loadAddon(addon)
  term.open(terminalRef.value)
  addon.fit()

  terminal.value = term
  fitAddon.value = addon

  term.writeln('NodeShell Terminal')
  term.writeln('直接在终端中输入命令，按 Enter 执行')
  term.writeln('')
  term.write(prompt.value)

  // 监听用户输入
  term.onData(async (data: string) => {
    if (executing.value) return

    const charCode = data.charCodeAt(0)

    // Enter键 - 执行命令
    if (charCode === 13) {
      term.write('\r\n')
      if (currentLine.value.trim()) {
        await executeCommand(currentLine.value.trim())
        currentLine.value = ''
        term.write(prompt.value)
      } else {
        term.write(prompt.value)
      }
      return
    }

    // Backspace键
    if (charCode === 127) {
      if (currentLine.value.length > 0) {
        currentLine.value = currentLine.value.slice(0, -1)
        term.write('\b \b')
      }
      return
    }

    // Ctrl+C - 取消当前命令
    if (charCode === 3) {
      term.write('^C\r\n')
      currentLine.value = ''
      term.write(prompt.value)
      return
    }

    // 普通字符输入
    if (charCode >= 32) {
      currentLine.value += data
      term.write(data)
    }
  })

  term.focus()

  window.addEventListener('resize', () => {
    addon.fit()
  })
}

const executeCommand = async (cmd: string) => {
  if (!cmd || executing.value) return

  terminalStore.addToHistory(cmd)
  executing.value = true

  try {
    const response = await shellApi.execute(shellId, cmd)
    if (terminal.value) {
      if (response.data) {
        const lines = response.data.split('\n')
        for (const line of lines) {
          if (line === '' && lines.indexOf(line) === lines.length - 1) {
            continue
          }
          terminal.value.writeln(line || '')
        }
      } else {
        terminal.value.writeln('(无输出)')
      }
    }
  } catch (error: any) {
    if (terminal.value) {
      terminal.value.writeln(`错误: ${error.message || '执行失败'}`)
    }
  } finally {
    executing.value = false
    if (terminal.value) {
      terminal.value.focus()
    }
  }
}

const fetchFiles = async () => {
  console.log('[fetchFiles] 开始获取文件列表, shellId:', shellId, 'currentPath:', currentPath.value)
  filesLoading.value = true
  try {
    // 如果currentPath为空，传递"."给后端以获取当前工作目录
    const pathToFetch = currentPath.value || '.'
    console.log('[fetchFiles] 请求路径:', pathToFetch)
    const result = await fileApi.list(shellId, pathToFetch)
    console.log('[fetchFiles] 获取结果:', result)
    if (result && result.files) {
      files.value = result.files
      if (result.path) {
        currentPath.value = result.path
      }
      console.log('[fetchFiles] 文件列表更新, 文件数量:', result.files.length)
    }
  } catch (error: any) {
    console.error('[fetchFiles] 获取文件列表失败:', error)
    ElMessage.error(error.message || '获取文件列表失败')
  } finally {
    filesLoading.value = false
  }
}

const refreshFiles = () => {
  fetchFiles()
}

const goUp = () => {
  if (!currentPath.value || currentPath.value === '/' || currentPath.value === '.') {
    // 如果已经是根目录或当前目录，无法再向上
    ElMessage.info('已经在根目录')
    return
  }
  const parts = currentPath.value.split('/').filter(p => p && p !== '.')
  if (parts.length > 0) {
    parts.pop()
    currentPath.value = parts.length > 0 ? '/' + parts.join('/') : '/'
  } else {
    currentPath.value = '/'
  }
  fetchFiles()
}

const changePath = () => {
  fetchFiles()
}

const enterDirectory = (path: string) => {
  console.log('[enterDirectory] 点击进入目录:', path)
  if (!path) {
    ElMessage.warning('路径无效')
    return
  }
  currentPath.value = path
  fetchFiles()
}

const handleFileDoubleClick = (row: FileItem) => {
  if (row.type === 'd') {
    enterDirectory(row.path)
  } else {
    readFile(row.path)
  }
}

let closeMenuHandler: ((e: MouseEvent) => void) | null = null

const handleContextMenu = (row: FileItem, column: any, event: MouseEvent) => {
  event.preventDefault()
  event.stopPropagation()
  
  // 先关闭之前的菜单并清理旧的监听器
  if (closeMenuHandler) {
    document.removeEventListener('click', closeMenuHandler)
    document.removeEventListener('contextmenu', closeMenuHandler)
    closeMenuHandler = null
  }
  contextMenu.value.visible = false
  
  // 计算菜单位置，确保不会超出屏幕边界
  const menuWidth = 150 // 菜单宽度
  const estimatedMenuItemHeight = 36 // 每个菜单项的高度
  const menuPadding = 8 // 菜单内边距
  // 根据文件类型估算菜单项数量
  const itemCount = row.type === 'd' ? 2 : 5 // 目录2项（进入、删除），文件5项（查看、下载、编辑、删除）
  const estimatedMenuHeight = itemCount * estimatedMenuItemHeight + menuPadding * 2
  const padding = 10 // 距离屏幕边缘的最小距离
  
  let x = event.clientX
  let y = event.clientY
  
  // 检查右边界
  if (x + menuWidth + padding > window.innerWidth) {
    x = window.innerWidth - menuWidth - padding
    if (x < padding) x = padding
  }
  
  // 检查下边界
  if (y + estimatedMenuHeight + padding > window.innerHeight) {
    // 在鼠标上方显示
    y = event.clientY - estimatedMenuHeight
    // 如果还是超出上边界，则放在屏幕底部
    if (y < padding) {
      y = window.innerHeight - estimatedMenuHeight - padding
      if (y < padding) y = padding
    }
  }
  
  // 检查左边界
  if (x < padding) {
    x = padding
  }
  
  // 检查上边界
  if (y < padding) {
    y = padding
  }
  
  // 使用 nextTick 确保菜单先渲染
  nextTick(() => {
    contextMenu.value = {
      visible: true,
      x: x,
      y: y,
      item: row
    }
    
    // 在下一个tick中获取真实菜单高度并再次调整
    nextTick(() => {
      const menuEl = document.querySelector('.context-menu') as HTMLElement
      if (menuEl) {
        const realHeight = menuEl.offsetHeight
        const realWidth = menuEl.offsetWidth
        
        // 重新检查下边界
        if (contextMenu.value.y + realHeight + padding > window.innerHeight) {
          contextMenu.value.y = event.clientY - realHeight
          if (contextMenu.value.y < padding) {
            contextMenu.value.y = window.innerHeight - realHeight - padding
          }
        }
        
        // 重新检查右边界
        if (contextMenu.value.x + realWidth + padding > window.innerWidth) {
          contextMenu.value.x = window.innerWidth - realWidth - padding
          if (contextMenu.value.x < padding) {
            contextMenu.value.x = padding
          }
        }
      }
    })
  })
  
  // 点击其他地方关闭菜单
  closeMenuHandler = (e: MouseEvent) => {
    const target = e.target as HTMLElement
    const menuEl = document.querySelector('.context-menu')
    
    // 如果点击的不是菜单本身
    if (menuEl && !menuEl.contains(target)) {
      contextMenu.value.visible = false
      if (closeMenuHandler) {
        document.removeEventListener('click', closeMenuHandler)
        document.removeEventListener('contextmenu', closeMenuHandler)
        closeMenuHandler = null
      }
    }
  }
  
  // 延迟添加监听器，避免立即触发
  setTimeout(() => {
    if (closeMenuHandler) {
      document.addEventListener('click', closeMenuHandler, true)
      document.addEventListener('contextmenu', closeMenuHandler, true)
    }
  }, 100)
}

const closeContextMenu = () => {
  contextMenu.value.visible = false
  if (closeMenuHandler) {
    document.removeEventListener('click', closeMenuHandler)
    document.removeEventListener('contextmenu', closeMenuHandler)
    closeMenuHandler = null
  }
}

const handleMenuEnter = () => {
  if (contextMenu.value.item) {
    enterDirectory(contextMenu.value.item.path)
    closeContextMenu()
  }
}

const handleMenuView = () => {
  if (contextMenu.value.item) {
    readFile(contextMenu.value.item.path)
    closeContextMenu()
  }
}

const handleMenuDownload = () => {
  if (contextMenu.value.item) {
    downloadFile(contextMenu.value.item.path)
    closeContextMenu()
  }
}

const handleMenuEdit = () => {
  if (contextMenu.value.item) {
    readFile(contextMenu.value.item.path)
    // 读取文件后自动进入编辑模式
    setTimeout(() => {
      fileEditing.value = true
    }, 300)
    closeContextMenu()
  }
}

const handleMenuDelete = () => {
  if (contextMenu.value.item) {
    deleteFile(contextMenu.value.item.path)
    closeContextMenu()
  }
}

const readFile = async (path: string) => {
  console.log('[readFile] 点击查看文件:', path)
  if (!path) {
    ElMessage.warning('路径无效')
    return
  }
  try {
    const result = await fileApi.read(shellId, path)
    fileContent.value = result.content
    editingFilePath.value = path
    fileEditing.value = false
    fileDialogVisible.value = true
  } catch (error: any) {
    console.error('[readFile] 读取文件失败:', error)
    ElMessage.error(error.message || '读取文件失败')
  }
}

const saveFile = async () => {
  try {
    await fileApi.update(shellId, editingFilePath.value, fileContent.value)
    ElMessage.success('保存成功')
    fileEditing.value = false
    fileDialogVisible.value = false
  } catch (error: any) {
    ElMessage.error(error.message || '保存文件失败')
  }
}

const downloadFile = async (path: string) => {
  console.log('[downloadFile] 点击下载文件:', path)
  if (!path) {
    ElMessage.warning('路径无效')
    return
  }
  
  downloadProgress.value = {
    current: 0,
    total: 0,
    percentage: 0,
    status: 'success',
    visible: true
  }
  
  try {
    const CHUNK_SIZE = 300 * 1024 // 300KB
    
    // 先尝试获取第一个分片以确定文件大小和总块数
    const firstChunk = await fileApi.download(shellId, path, 0, CHUNK_SIZE)
    
    // 如果返回的是 JSON（分片响应），说明是大文件，需要分片下载
    if (firstChunk && typeof firstChunk === 'object' && 'total_chunks' in firstChunk) {
      const totalChunks = firstChunk.total_chunks
      const fileSize = firstChunk.file_size
      
      downloadProgress.value.total = totalChunks
      downloadProgress.value.current = 1
      downloadProgress.value.percentage = Math.round((1 / totalChunks) * 100)
      
      // 收集所有分片
      const chunks: string[] = []
      chunks[0] = firstChunk.data
      
      // 下载剩余分片
      for (let i = 1; i < totalChunks; i++) {
        const chunkResp = await fileApi.download(shellId, path, i, CHUNK_SIZE)
        chunks[i] = chunkResp.data
        
        // 更新进度
        downloadProgress.value.current = i + 1
        downloadProgress.value.percentage = Math.round(((i + 1) / totalChunks) * 100)
      }
      
      // 合并所有分片
      let binary = ''
      for (let i = 0; i < chunks.length; i++) {
        binary += atob(chunks[i])
      }
      
      // 转换为 Blob
      const uint8Array = new Uint8Array(binary.length)
      for (let i = 0; i < binary.length; i++) {
        uint8Array[i] = binary.charCodeAt(i)
      }
      const blob = new Blob([uint8Array])
      
      // 下载文件
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = path.split('/').pop() || 'download'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(url)
      
      downloadProgress.value.status = 'success'
      ElMessage.success('下载成功')
    } else {
      // 小文件：直接下载（兼容原有方式）
      downloadProgress.value.total = 1
      downloadProgress.value.current = 1
      downloadProgress.value.percentage = 100
      
      const blob = firstChunk as Blob
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = path.split('/').pop() || 'download'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      window.URL.revokeObjectURL(url)
      
      downloadProgress.value.status = 'success'
      ElMessage.success('下载成功')
    }
  } catch (error: any) {
    console.error('[downloadFile] 下载文件失败:', error)
    downloadProgress.value.status = 'exception'
    ElMessage.error(error.message || '下载失败')
  } finally {
    // 延迟关闭进度对话框，让用户看到完成状态
    setTimeout(() => {
      // 不自动关闭，让用户手动关闭
    }, 1000)
  }
}

const deleteFile = async (path: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个文件/目录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await fileApi.delete(shellId, path)
    ElMessage.success('删除成功')
    fetchFiles()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const showUploadDialog = () => {
  uploadForm.value.path = currentPath.value || '.'
  // 如果currentPath是目录，则在路径末尾添加文件名提示
  if (uploadForm.value.path && !uploadForm.value.path.includes('.')) {
    uploadForm.value.path = uploadForm.value.path.endsWith('/') 
      ? uploadForm.value.path + 'filename.txt'
      : uploadForm.value.path + '/filename.txt'
  }
  uploadForm.value.content = ''
  uploadForm.value.file = null
  uploadForm.value.fileName = ''
  uploadProgress.value = {
    current: 0,
    total: 0,
    percentage: 0,
    status: 'success',
    visible: false
  }
  uploadDialogVisible.value = true
}

const handleFileSelect = (file: any) => {
  // 统一走二进制读取，避免文本模式破坏二进制内容
  uploadForm.value.file = file.raw
  uploadForm.value.fileName = file.name
  uploadForm.value.content = '' // 清空文本内容，改用文件分片上传

  // 自动填充文件名到路径
  if (!uploadForm.value.path || uploadForm.value.path.endsWith('/') || uploadForm.value.path.endsWith('filename.txt')) {
    const basePath = currentPath.value || '.'
    uploadForm.value.path = basePath.endsWith('/') 
      ? basePath + file.name
      : basePath + '/' + file.name
  }
}

const uploadFile = async () => {
  if (!uploadForm.value.path) {
    ElMessage.warning('请输入文件路径')
    return
  }
  
  if (!uploadForm.value.file) {
    ElMessage.warning('请选择文件')
    return
  }
  
  uploading.value = true
  uploadProgress.value = {
    current: 0,
    total: 0,
    percentage: 0,
    status: 'success',
    visible: true
  }
  
  try {
    const CHUNK_SIZE = 70 * 1024 // 70KB
    const file = uploadForm.value.file as File
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE)
    
    uploadProgress.value.total = totalChunks
    
    // 读取文件并分片上传
    for (let i = 0; i < totalChunks; i++) {
      const start = i * CHUNK_SIZE
      const end = Math.min(start + CHUNK_SIZE, file.size)
      const chunk = file.slice(start, end)
      
      // 将 chunk 读取为 ArrayBuffer，然后转换为 Base64
      const arrayBuffer = await chunk.arrayBuffer()
      const uint8Array = new Uint8Array(arrayBuffer)
      let binary = ''
      for (let j = 0; j < uint8Array.length; j++) {
        binary += String.fromCharCode(uint8Array[j])
      }
      const base64 = btoa(binary)
      
      await fileApi.upload(shellId, uploadForm.value.path, base64, i, totalChunks)
      
      // 更新进度
      uploadProgress.value.current = i + 1
      uploadProgress.value.percentage = Math.round(((i + 1) / totalChunks) * 100)
    }
    
    uploadProgress.value.status = 'success'
    ElMessage.success('上传成功')
    uploadDialogVisible.value = false
    fetchFiles()
  } catch (error: any) {
    uploadProgress.value.status = 'exception'
    ElMessage.error(error.message || '上传失败')
  } finally {
    uploading.value = false
    // 延迟关闭进度对话框
    setTimeout(() => {
      uploadProgress.value.visible = false
      uploadProgress.value = {
        current: 0,
        total: 0,
        percentage: 0,
        status: 'success',
        visible: false
      }
    }, 2000)
  }
}

onMounted(async () => {
  await loadShell()
  // 根据路由参数决定默认显示哪个标签
  const tab = (route.params.tab as string) || 'terminal'
  if (tab === 'system' || tab === 'terminal' || tab === 'files') {
    activeTab.value = tab as 'system' | 'terminal' | 'files'
  }
  
  // 初始化终端（如果需要）
  if (activeTab.value === 'terminal') {
    await initTerminal()
  }
  
  // 加载文件列表（如果需要）
  if (activeTab.value === 'files') {
    currentPath.value = ''
    await fetchFiles()
  }
})

// 监听标签切换
watch(activeTab, async (newTab) => {
  if (newTab === 'terminal' && !terminal.value) {
    await initTerminal()
  }
  if (newTab === 'files') {
    currentPath.value = ''
    await fetchFiles()
  }
  router.replace(`/shell/${shellId}/${newTab}`)
})

onUnmounted(() => {
  if (terminal.value) {
    terminal.value.dispose()
  }
  window.removeEventListener('resize', () => {})
  // 清理右键菜单监听器
  if (closeMenuHandler) {
    document.removeEventListener('click', closeMenuHandler)
    document.removeEventListener('contextmenu', closeMenuHandler)
    closeMenuHandler = null
  }
})
</script>

<style scoped>
.shell-detail {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8ecf1 100%);
}

.header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 30px;
  background: #ffffff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  z-index: 100;
}

.header-left {
  display: flex;
  gap: 12px;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-title h2 {
  margin: 0;
  color: #1d2129;
  font-weight: 600;
  font-size: 20px;
  letter-spacing: -0.3px;
}

.content-wrapper {
  display: flex;
  flex-direction: column;
  flex: 1;
  gap: 16px;
  padding: 16px;
  overflow: auto;
  background: transparent;
}

/* 顶部导航卡片 */
.nav-cards {
  display: flex;
  gap: 12px;
  margin-bottom: 0;
}

.nav-card {
  flex: 1;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 12px 20px;
  background: #ffffff;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.04);
}

.nav-card:hover {
  border-color: #409eff;
  background: #f0f7ff;
  transform: translateY(-1px);
  box-shadow: 0 3px 10px rgba(64, 158, 255, 0.12);
}

.nav-card.active {
  border-color: #409eff;
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  color: #ffffff;
  box-shadow: 0 3px 12px rgba(64, 158, 255, 0.25);
}

.nav-card.active .nav-icon {
  color: #ffffff;
}

.nav-card.active .nav-title {
  color: #ffffff;
}

.nav-icon {
  font-size: 18px;
  color: #409eff;
  transition: all 0.3s ease;
}

.nav-card:not(.active) .nav-icon {
  color: #606266;
}

.nav-title {
  font-size: 13px;
  font-weight: 600;
  color: #1d2129;
  transition: all 0.3s ease;
}

/* 内容区域 */
.content-area {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.tab-content {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.content-card {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  border: 1px solid #e4e7ed;
}

.content-card :deep(.el-card__body) {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
}


.terminal-wrapper {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  background: #1e1e1e;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.3);
}

.terminal-container {
  flex: 1;
  min-height: 0;
  background: #1e1e1e;
  padding: 12px;
}


.context-menu {
  position: fixed;
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
  z-index: 9999;
  min-width: 160px;
  max-width: 200px;
  padding: 6px 0;
  user-select: none;
  backdrop-filter: blur(10px);
}

.context-menu-item {
  padding: 10px 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: #1d2129;
  transition: all 0.2s ease;
  font-weight: 400;
}

.context-menu-item:hover {
  background: linear-gradient(90deg, #f0f7ff 0%, #e8f4ff 100%);
  color: #409eff;
}

.context-menu-item-danger {
  color: #f56c6c;
}

.context-menu-item-danger:hover {
  background: linear-gradient(90deg, #fef0f0 0%, #fde8e8 100%);
  color: #f56c6c;
}

.context-menu-item .el-icon {
  font-size: 16px;
  width: 18px;
  text-align: center;
}

.file-manager {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  gap: 16px;
}

.file-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}
</style>

