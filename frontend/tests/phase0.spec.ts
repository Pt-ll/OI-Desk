import { test, expect } from '@playwright/test'
import { pathToFileURL } from 'node:url'
import { resolve } from 'node:path'

test('sandboxed plugin completes isolation and MessageChannel assertions', async ({ page }) => {
  const fixture = pathToFileURL(resolve('../phase0/plugin-host.html')).href
  await page.goto(fixture)
  await expect(page.locator('#status')).toHaveAttribute('data-success', 'true')
  await expect(page.locator('#status')).toHaveText('8/8 isolation and transport checks passed')
  await expect(page.locator('#plugin')).toHaveAttribute('sandbox', 'allow-scripts')
})
