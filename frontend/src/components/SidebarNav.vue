<script lang="ts" setup>
import { navItems, useShellStore, type RouteId } from '../stores/shell'
const shell = useShellStore()
function navigate(route: RouteId) { shell.navigate(route) }
</script>
<template>
  <aside class="sidebar" aria-label="主导航">
    <div class="brand"><div class="brand-mark">OI</div><div><strong>OI-Desk</strong><span>Training workspace</span></div></div>
    <nav class="nav-list">
      <button v-for="item in navItems" :key="item.id" class="nav-item" :class="{ active: shell.activeRoute === item.id }" :aria-current="shell.activeRoute === item.id ? 'page' : undefined" :data-testid="`nav-${item.id}`" @click="navigate(item.id)">
        <span class="nav-icon" aria-hidden="true">{{ item.icon }}</span><span class="nav-copy"><strong>{{ item.label }}</strong><small>{{ item.description }}</small></span>
      </button>
    </nav>
    <div class="sidebar-footer"><span class="status-dot" aria-hidden="true"></span><span>Local-first mode</span></div>
  </aside>
</template>
