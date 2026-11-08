import { defineStore } from 'pinia'
import type { AppError } from '../types/errors'

export type Theme = 'dark' | 'light'
export type RouteId = 'dashboard' | 'oj' | 'algorithms' | 'plugins' | 'ai' | 'settings'
export interface NavItem { id: RouteId; label: string; description: string; icon: string }
export const navItems: NavItem[] = [
  { id: 'dashboard', label: 'Dashboard', description: '训练概览与快捷入口', icon: '⌂' },
  { id: 'oj', label: 'OJ', description: '账号、题目与提交记录', icon: '◈' },
  { id: 'algorithms', label: 'Algorithms', description: '算法分类与掌握状态', icon: '⌘' },
  { id: 'plugins', label: 'Plugins', description: '插件安装与权限', icon: '⊞' },
  { id: 'ai', label: 'AI', description: 'Provider、模型与用量', icon: '✦' },
  { id: 'settings', label: 'Settings', description: '应用设置与数据管理', icon: '⚙' },
]
const THEME_KEY = 'oidesk.theme'
export const useShellStore = defineStore('shell', {
  state: () => ({ theme: 'dark' as Theme, activeRoute: 'dashboard' as RouteId, error: null as AppError | null, bindingOK: false, pluginOK: false }),
  getters: { activeNav: (state) => navItems.find((item) => item.id === state.activeRoute) ?? navItems[0] },
  actions: {
    initialize() { if (typeof window !== 'undefined') { const saved = window.localStorage.getItem(THEME_KEY); if (saved === 'dark' || saved === 'light') this.theme = saved }; this.applyTheme() },
    applyTheme() { if (typeof document !== 'undefined') document.documentElement.dataset.theme = this.theme },
    setTheme(theme: Theme) { this.theme = theme; this.applyTheme(); if (typeof window !== 'undefined') window.localStorage.setItem(THEME_KEY, theme) },
    navigate(route: RouteId) { this.activeRoute = route; this.error = null },
    setError(error: AppError) { this.error = error },
    clearError() { this.error = null },
    setDiagnostics(bindingOK: boolean, pluginOK: boolean) { this.bindingOK = bindingOK; this.pluginOK = pluginOK },
  },
})
