import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './tests',
  timeout: 10_000,
  webServer: { command: 'vite preview --host 127.0.0.1 --port 4173', url: 'http://127.0.0.1:4173', reuseExistingServer: true },
  use: { browserName: 'chromium', headless: true },
})
