import { test, expect } from '@playwright/test';
import testPrefs  from './test-prefs';

export default function createTests() {

    test.describe.configure({ mode: 'serial' });

    test.beforeEach(async ({ page }) => {
        await page.goto('');
    });

    test.describe('Check initial preferences', () => {

        test('click Help/Preferences', async ({ page }) => {
            await page.locator("css=#helpButton").click()

            //const item = await page.$('#preferencesMenuItem')
            //await item.waitForDisplayed()

            await page.locator('css=#preferencesMenuItem').click()

            await page.goto('/prefs.html')
            await expect(page).toHaveTitle(/Synergize Preferences/);

            await page.locator('#libraryPath').fill('./data/testfiles')

            await page.locator('#saveButton').click()


            await page.goto('/')
            await expect(page).toHaveTitle(/Synergize/);

            await expect(page.locator('#path')).toHaveText('testfiles')
        });


        test('show main window', async ({ page }) => {

            await page.goto('/')
            await expect(page).toHaveTitle(/Synergize/);
        });

    });

}