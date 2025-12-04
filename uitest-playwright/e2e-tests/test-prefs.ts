import { test, expect } from '@playwright/test';
import testPrefs  from './test-prefs';

export default function createTests() {

    test.describe.configure({ mode: 'serial' });

    test.beforeEach(async ({ page }) => {
        await page.goto('');
    });

    test.describe('Check initial preferences', () => {

        test('click Help/Preferences', async ({ page }) => {
            const context = page.context();
            await page.locator("css=#helpButton").click()

            const prefsPagePromise = context.waitForEvent('page');

            await page.locator('css=#preferencesMenuItem').click()

            const prefsPage = await context.newPage();
            await prefsPage.goto('/prefs.html')

            await expect(prefsPage).toHaveTitle(/Synergize Preferences/);

            await prefsPage.locator('#libraryPath').fill('./data/testfiles')
            await prefsPage.locator('#saveButton').click()
            
            await page.goto('/')
            
            await expect(page).toHaveTitle(/Synergize/);
            await expect(page.locator('#path')).toHaveText('testfiles')
        });


//        test('show main window', async ({ page }) => {
//
//            await page.goto('/')
//            await expect(page).toHaveTitle(/Synergize/);
//        });

    });

}