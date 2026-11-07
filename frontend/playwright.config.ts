import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './tests',
  timeout: 10_000,
  use: { browserName: 'chromium', headless: true },
})
