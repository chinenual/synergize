
export default function createTests() {

    test.describe.configure({ mode: 'serial' });

    test.beforeEach(async ({ page }) => {
        await page.goto('');
    });

    test.describe('Check initial preferences', () => {

        test('click Help/Preferences', async () => {
            await page.locator("css=#helpButton").click()

            //const item = await page.$('#preferencesMenuItem')
            //await item.waitForDisplayed()

            await page.locator('css=#preferencesMenuItem').click()

            await page.switchWindow('Synergize Preferences');
            await hooks.screenshotAndCompare(app, 'prefsWindow');

            (await page.getTitle()).should.equal('Synergize Preferences')

            const txt = await page.$('#libraryPath')
            await txt.setValue('../data/testfiles')

            const submit = await page.$('button[type=submit]')
            await submit.click()


            await page.switchWindow('Synergize');
            (await page.getTitle()).should.equal('Synergize')

            const txt2 = await page.$('#path');

            (await txt2.getText()).should.equal('testfiles')
        });


        test('show main window', async () => {

            await page.switchWindow('Synergize');
            (await page.getTitle()).should.equal('Synergize')
        });

    });

}