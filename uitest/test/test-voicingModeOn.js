const hooks = require('./hooks');
const viewVCE = require('./test-viewVCE');
const WINDOW_PAUSE = 1000;
const INIT_VOICING_TIMEOUT = 30000; // 30s to init voicemode and load the initial VRAM image
const voiceINITVRAM = require('./page-objects/voice-INITVRAM');

let app;

describe('Test Voicing Mode ON', () => {
  afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this,app); });
  before(async () => {
    console.log("====== reuse the app");
    app = await hooks.getApp();
  });

  it('should enable voicing mode', async () => {
    await app.client
      .click('#voicingModeButton')

      .waitForVisible('#alertText', INIT_VOICING_TIMEOUT)
      .then(() => {return hooks.screenshotAndCompare(app, 'voicingModeOn-alert')})

      .getText('#alertText').should.eventually.include('enabled')

      .click('#alertModal button')
      .waitForVisible('#alertText', 1000, true) // wait to disappear

      .getAttribute("#voicingModeButton img", "src").should.eventually.include('static/images/red-button-on-full.png')

      .then(() => {return hooks.screenshotAndCompare(app, 'voicingModeOn')})
  });

  describe('initial VRAM image should be loaded', () => {
    viewVCE.testViewVCE([voiceINITVRAM], null, "voicemode");
  });
});

