<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { Greet, ReportFrontendProbe } from '../wailsjs/go/main/App'
const status = ref('正在验证 Wails 绑定…')
const bindingOK = ref(false)
const pluginOK = ref(false)
const probeFrame = ref<HTMLIFrameElement>()
onMounted(async () => {
  window.addEventListener('message', async (event: MessageEvent) => {
    if (event.source !== probeFrame.value?.contentWindow || event.data?.type !== 'phase0-result') return
    pluginOK.value = event.data.success === true
    status.value = event.data.detail
    await ReportFrontendProbe(bindingOK.value, pluginOK.value)
  })
  try {
    bindingOK.value = (await Greet('Phase 0')) === "Hello Phase 0, It's show time!"
    status.value = bindingOK.value ? 'Wails Go ↔ Vue 通信通过，正在验证隔离…' : '绑定结果不匹配'
  } catch (error) { status.value = `绑定失败：${String(error)}` }
})
</script>
<template>
  <main>
    <h1>OI-Desk · Phase 0</h1>
    <p>技术验证工作台 · 不包含业务功能</p>
    <p>{{ status }}</p>
    <p>Go ↔ Vue：{{ bindingOK ? '通过' : '等待' }} · iframe：{{ pluginOK ? '通过' : '等待' }}</p>
    <iframe v-if="bindingOK" ref="probeFrame" title="Phase 0 plugin isolation probe" src="/phase0/plugin-host.html" />
  </main>
</template>
<style scoped>
main { padding: 2rem; }
iframe { width: 90%; height: 240px; background: white; border: 1px solid #667; }
</style>
