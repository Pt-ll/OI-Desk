<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { Greet, ReportFrontendProbe } from '../wailsjs/go/main/App'
import ErrorNotice from './components/ErrorNotice.vue'
import SidebarNav from './components/SidebarNav.vue'
import { normalizeError } from './types/errors'
import { useShellStore } from './stores/shell'

const shell = useShellStore()
const status = ref('正在验证 Wails 绑定…')
const probeFrame = ref<HTMLIFrameElement | null>(null)

async function onProbeResult(event: MessageEvent) {
  if (event.source !== probeFrame.value?.contentWindow || event.data?.type !== 'phase0-result') return
  shell.setDiagnostics(shell.bindingOK, event.data.success === true)
  status.value = event.data.detail
  try { await ReportFrontendProbe(shell.bindingOK, shell.pluginOK) } catch (error) { shell.setError(normalizeError(error, '无法回传 Phase 0 诊断结果')) }
}

onMounted(async () => {
  await shell.initialize()
  window.addEventListener('message', onProbeResult)
  try {
    shell.setDiagnostics((await Greet('Phase 0')) === "Hello Phase 0, It's show time!", false)
    status.value = shell.bindingOK ? 'Wails Go ↔ Vue 通信通过，正在验证隔离…' : '绑定结果不匹配'
  } catch (error) { shell.setError(normalizeError(error, '无法连接 Go 后端')); status.value = '绑定失败' }
})
onUnmounted(() => window.removeEventListener('message', onProbeResult))
</script>

<template>
  <div class="shell">
    <SidebarNav />
    <main class="content">
      <header class="topbar">
        <div><span class="eyebrow">OI-Desk / workspace</span><h1 data-testid="route-title">{{ shell.activeNav.label }}</h1></div>
        <div class="topbar-actions"><span class="sync-status"><span class="status-dot" aria-hidden="true"></span>本地数据</span><button class="theme-toggle" type="button" data-testid="theme-toggle" @click="shell.setTheme(shell.theme === 'dark' ? 'light' : 'dark')">{{ shell.theme === 'dark' ? '☼ 浅色' : '◐ 深色' }}</button></div>
      </header>
      <ErrorNotice v-if="shell.error" :error="shell.error" @dismiss="shell.clearError" />
      <section class="hero-card"><div><span class="eyebrow">{{ shell.activeNav.description }}</span><h2>把训练数据留在自己的工作台</h2><p>基础导航和主题系统已就绪，后续模块会在各自阶段逐步接入。</p></div><div class="hero-badge"><strong>Phase 1</strong><span>foundation</span></div></section>
      <section class="workspace-grid">
        <article v-for="item in [{ title: '本地优先', text: '数据默认保存在本机，服务边界保持清晰。', icon: '◉' }, { title: '可扩展', text: 'OJ、插件和 AI 按独立阶段接入。', icon: '◇' }, { title: '可验证', text: '每个阶段都有独立测试和验收条件。', icon: '✓' }]" :key="item.title" class="info-card"><span class="card-icon">{{ item.icon }}</span><h3>{{ item.title }}</h3><p>{{ item.text }}</p></article>
      </section>
      <footer class="diagnostics" data-testid="diagnostics">{{ status }} · Go：{{ shell.bindingOK ? '通过' : '等待' }} · iframe：{{ shell.pluginOK ? '通过' : '等待' }}</footer>
      <iframe ref="probeFrame" class="probe-frame" title="Phase 0 plugin isolation probe" src="/phase0/plugin-host.html" aria-hidden="true" />
    </main>
  </div>
</template>
