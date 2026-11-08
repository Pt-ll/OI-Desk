import { test, expect } from '@playwright/test'
test('shell navigation and theme toggle are deterministic', async ({ page }) => {
  await page.goto('http://127.0.0.1:4173/')
  await expect(page.getByTestId('route-title')).toHaveText('Dashboard')
  await page.getByTestId('nav-settings').click()
  await expect(page.getByTestId('route-title')).toHaveText('Settings')
  await page.getByTestId('theme-toggle').click()
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'light')
})
