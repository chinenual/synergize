const hooks = require('./hooks');
const config = require('../config').get(process.env.NODE_ENV);
const viewVCE = require('./test-viewVCE');


const WINDOW_PAUSE = 1000;

const voiceG7S = require('./page-objects/voice-G7S');
const voiceCATHERG = require('./page-objects/voice-CATHERG');
const voiceGUITAR2A = require('./page-objects/voice-GUITAR2A');


//const SearchPage = require('./page-objects/search.page');

var app;

describe('Setup', () => {
  before(async () => {
    console.log("====== remove test preferences.json file if exists");

    var fs = require('fs');
    var path = './preferences.json';
    if (fs.existsSync(path)) {
      fs.unlinkSync(path);
    }
    console.log("====== start up the app");
    app = await hooks.startApp();
  });

  it('opens a window', async () => {
    await app.client
      // HACK:  astilectron seems to not have executed onWait if we start too fast - 
      // we often get a segfault when trying to display the about window - prob because 
      // the about_w variable hasnt been initialized yet?  So wait a bit...
      .pause(5000)
      .waitUntilWindowLoaded()
      .getWindowCount()
      .should.eventually.be.above(0)
      .saveScreenshot('screenshots-mainWindow-startup.png')
      .getTitle().should.eventually.equal('Synergize')
  });

});

require('./test-about');

require('./test-prefs');
describe('Test READ-ONLY views', () => {
  viewVCE.testViewVCE([voiceG7S, voiceCATHERG, voiceGUITAR2A], viewVCE.loadVCEViaLeftPanel, "readonly");
});
describe('Test Voicing Mode views', () => {
  require('./test-voicingModeOn');
  viewVCE.testViewVCE([voiceG7S, voiceCATHERG, voiceGUITAR2A], viewVCE.loadVCEViaLeftPanel, "voicemode");
  require('./test-voicingModeOff');
});



describe('Tear Down', () => {
  after(async () => {
    console.log("====== tear down the app");
    await hooks.stopApp(app);
  });
  it('last gasp', async () => {
    await app.client
      .getTitle().should.eventually.equal('Synergize')
  });
});

