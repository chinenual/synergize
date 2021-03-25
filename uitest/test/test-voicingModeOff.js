const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;

let app;

describe('Test Voicing Mode OFF', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    it('should disable voicing mode', async () => {
        const voicingModeButton = await app.client.$('#voicingModeButton')
        await voicingModeButton.click()

        const confirmText = await app.client.$('#confirmText')

        await confirmText.waitForDisplayed();
        (await confirmText.getText()).should.include('pending edits');
        
        const confirmOk = await app.client.$('#confirmOKButton')
        await confirmOk.click()

        await confirmText.waitForDisplayed({reverse: true})  // wait to disappear
        
        const alertText = await app.client.$('#alertText')
        const img = await app.client.$('#voicingModeButton img')

        await alertText.waitForDisplayed()
        await hooks.screenshotAndCompare(app, 'voicingModeOff-alert');

        (await alertText.getText()).should.include('disabled')

        const alertClose = await app.client.$('#alertModal button')
        await alertClose.click()

        await alertText.waitForDisplayed({reverse: true});  // wait to disappear

        (await img.getAttribute("src")).should.include('static/images/red-button-off-full.png');

        await hooks.screenshotAndCompare(app, 'voicingModeOff')

    });
});

