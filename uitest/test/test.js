const hooks = require('./hooks');
const config = require('../config').get(process.env.NODE_ENV);
const WINDOW_PAUSE = 1000;

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
      .waitUntilWindowLoaded()
      .getWindowCount()
      .should.eventually.be.above(0)
      .saveScreenshot('screenshots-mainWindow-startup.png')
      .getTitle().should.eventually.equal('Synergize')
  });

});


require('./test-about');
require('./test-prefs');
require('./test-viewVCE');



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

