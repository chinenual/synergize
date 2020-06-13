const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;
const INIT_VOICING_TIMEOUT = 30000; // 30s to init voicemode and load the initial VRAM image

let app;

describe('Test Voicing Mode ON', () => {
  before(async () => {
    console.log("====== reuse the app");
    app = await hooks.getApp();
  });

  it('should enable voicing mode', async () => {
    await app.client
      .click('#voicingModeButton')

      .waitForVisible('#alertText', INIT_VOICING_TIMEOUT)
      .saveScreenshot('screenshots-voicingModeOn-alert.png')

      .getText('#alertText').should.eventually.include('enabled')

      .click('#alertModal button')
      .waitForVisible('#alertText', 1000, true) // wait to disappear

      .getAttribute("#voicingModeButton img", "src").should.eventually.include('static/images/red-button-on-full.png')

      .saveScreenshot('screenshots-voicingModeOn.png')
  });
});

