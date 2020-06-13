const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;

let app;

describe('Test Voicing Mode OFF', () => {
  before(async () => {
    console.log("====== reuse the app");
    app = await hooks.getApp();
  });

  it('should disable voicing mode', async () => {
    await app.client
      .click('#voicingModeButton')

      .waitForVisible('#alertText')
      .then(() => {return hooks.screenshotAndCompare(app, 'voicingModeOff-alert')})

      .getText('#alertText').should.eventually.include('disabled')

      .click('#alertModal button')
      .waitForVisible('#alertText', 1000, true) // wait to disappear

      .then(() => {return hooks.screenshotAndCompare(app, 'voicingModeOff')})
      .getAttribute("#voicingModeButton img", "src").should.eventually.include('static/images/red-button-off-full.png')

      .then(() => {return hooks.screenshotAndCompare(app, 'voicingModeOff')})
  });
});

