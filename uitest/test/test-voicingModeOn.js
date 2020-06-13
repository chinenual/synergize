const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;

let app;

describe('Test Voicing Mode ON', () => {
  before(async () => {
    console.log("====== reuse the app");
    app = await hooks.getApp();
  });

  it('should enable voicing mode', async () => {
    await app.client
      .click('#voicingModeButton')

      .waitForVisible('#alertText')
      .saveScreenshot('screenshots-voicingModeOn-alert.png')

      .getText('#alertText').should.eventually.include('enabled')

      .click('#alertModal button')
      .waitForVisible('#alertText', {reverse: true})

      .getAttribute("#voicingModeButton img", "src").should.eventually.include('static/images/red-button-on-full.png')

      .saveScreenshot('screenshots-voicingModeOn.png')
  });
});

