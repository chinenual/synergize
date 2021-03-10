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

        await confirmText.waitForDisplayed()
        await confirmText.getText().should.eventually.include('pending edits')
        
        const confirmOk = await app.client.$('#confirmOKButton')
console.log("OFF 4")
        await confirmOk.click()
console.log("OFF 3")

        await confirmText.waitForDisplayed({reverse: true})  // wait to disappear
        
        const alertText = await app.client.$('#alertText')
        const img = await app.client.$('#voicingModeButton img')

        await alertText.waitForDisplayed()
        await hooks.screenshotAndCompare(app, 'voicingModeOff-alert')

        await alertText.getText().should.eventually.include('disabled')

        const alertClose = await app.client.$('#alertModal button')
console.log("OFF 2")
        await alertClose.click()
console.log("OFF 1")

        await alertText.waitForDisplayed({reverse: true})  // wait to disappear

        await img.getAttribute("src").should.eventually.include('static/images/red-button-off-full.png')

        await hooks.screenshotAndCompare(app, 'voicingModeOff')

    });
});

