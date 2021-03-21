const hooks = require('./hooks');
const viewVCE = require('./test-viewVCE');
const WINDOW_PAUSE = 1000;
const INIT_VOICING_TIMEOUT = 30000; // 30s to init voicemode and load the initial VRAM image
const voiceINITVRAM = require('./page-objects/voice-INITVRAM');

let app;

describe('Test Voicing Mode ON', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    it('should enable voicing mode', async () => {
        const voicingModeButton = await app.client.$('#voicingModeButton')
        await voicingModeButton.click()

        const alertText = await app.client.$('#alertText')
        const button = await app.client.$('#alertModal button')
        const img = await app.client.$('#voicingModeButton img')

        await alertText.waitForDisplayed({timeout: INIT_VOICING_TIMEOUT})
        await hooks.screenshotAndCompare(app, 'voicingModeOn-alert');

        (await alertText.getText()).should.include('enabled');

        await button.click()

        await alertText.waitForDisplayed({timeout: 1000, reverse: true});  // wait to disappear

        (await img.getAttribute("src")).should.include('static/images/red-button-on-full.png')

        await hooks.screenshotAndCompare(app, 'voicingModeOn')
    });

});

