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
      .saveScreenshot('screenshots-voicingModeOn-alert.png')

      .getText('#alertText').should.eventually.include('disabled')

      .click('#alertModal button')
      .waitForVisible('#alertText', {reverse: true})

      //.saveScreenshot('screenshots-voicingModeOffXXX.png')
      //.getAttribute("#voicingModeButton img", "src").should.eventually.include('static/images/red-button-off-full.png')

      .saveScreenshot('screenshots-voicingModeOff.png')
  });
});

