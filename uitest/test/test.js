const hooks = require('./hooks');
const config = require('../config').get(process.env.NODE_ENV);
const viewVCE = require('./test-viewVCE');

var chai = require('chai')
chai.config.includeStack = true


const WINDOW_PAUSE = 1000;

const voiceG7S = require('./page-objects/voice-G7S');
const voiceCATHERG = require('./page-objects/voice-CATHERG');
const voiceGUITAR2A = require('./page-objects/voice-GUITAR2A');


//const SearchPage = require('./page-objects/search.page');

var app;

describe('Setup', () => {
    afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this,app); });

    before(async () => {
        console.log("====== remove test preferences.json file if exists");

        var fs = require('fs');
        var path = './preferences.json';
        if (fs.existsSync(path)) {
            fs.unlinkSync(path);
        }
        console.log("====== start up the app");
        app = await hooks.startApp();

        app.client.setTimeout({'implicit': 0})
    });

    it('opens a window', async () => {
        await app.client
        // HACK:  astilectron seems to not have executed onWait if we start too fast - 
        // we often get a segfault when trying to display the about window - prob because 
        // the about_w variable hasnt been initialized yet?  So wait a bit...
            .pause(5000)
        await app.client
            .waitUntilWindowLoaded()
        await app.client
            .getWindowCount()
            .should.eventually.be.above(0)

        // no screenshot: default lib directory is my $HOME - which can change from minute to minute
        // .then(() => {return hooks.screenshotAndCompare(app, 'mainWindow-startup')})

        await app.client
            .getTitle().should.eventually.equal('Synergize')

    });

});


/*
describe('Render unit tests', () => {

    it('voice page text conversions', async () => {
        await app.webContents
            .executeJavaScript('viewVCE_voice.testConversionFunctions()').should.eventually.be.true
    });
    it('env page text conversions', async () => {
        await app.webContents
            .executeJavaScript('viewVCE_envs.testConversionFunctions()').should.eventually.be.true
    });
});

require('./test-about');
*/

require('./test-prefs');

/*
  require('./test-edit-crt');

  describe('Test Voicing Mode views', () => {
  afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this,app); });

  require('./test-voicingModeOn');

  require('./test-voice-edit');
  require('./test-envs-edit');
  require('./test-filter-edit.js');
  require('./test-keyeq-edit');
  require('./test-keyprop-edit');

  viewVCE.testViewVCE([voiceG7S, voiceCATHERG, voiceGUITAR2A], viewVCE.loadVCEViaLeftPanelVoicingMode, "voicemode");

  require('./test-voicingModeOff');
  });
*/

describe('Test READ-ONLY views', () => {
    afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this,app); });
    
    viewVCE.testViewVCE([voiceG7S, voiceCATHERG, voiceGUITAR2A], viewVCE.loadVCEViaLeftPanel, "readonlyVCE");
    viewVCE.testViewVCE([voiceG7S, voiceCATHERG, voiceGUITAR2A], viewVCE.loadVCEViaINTERNALCRT, "readonlyCRT");
});

describe('Tear Down', () => {
    afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this,app); });
    after(async () => {
        console.log("====== tear down the app");
        await hooks.stopApp(app);
    });
    it('last gasp', async () => {
        await app.client
            .getTitle().should.eventually.equal('Synergize')
    });
});

