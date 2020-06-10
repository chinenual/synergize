const hooks = require('./hooks');
const config = require('../config').get(process.env.NODE_ENV);
//const SearchPage = require('./page-objects/search.page');

// switchWindow seems unreliable - ~20% of the time I call it to check the About window, it crashes electron.  I can't find a way to "wait" for it to be visible (assuming the problem is a timing error and switchToWindow doesnt know what to do with a hidden window) 
var mainWindow, aboutWindow, prefsWindow;

describe('Sample Test', () => {
  let app;

  before(async () => {
    app = await hooks.startApp();
  });

  after(async () => {
    await hooks.stopApp(app);
  });

  it('opens a window', async () => {
    await app.client
      .waitUntilWindowLoaded()
      .getWindowCount()
      .should.eventually.be.above(0)
      .saveScreenshot('screenshots-mainWindow-startup.png')
      .getTitle()
      .should.eventually.equal('Synergize')

    await app.client.windowHandles()
      .then((handles) => {
        console.log("***** window handles: " + JSON.stringify(handles));
        mainWindow = handles.value[0];
        aboutWindow = handles.value[1];
        prefsWindow = handles.value[2];
      });

      var n = await app.client.window(mainWindow).getTitle();
      console.log(" title[mainWindow]: " + n)
      var n = await app.client.window(aboutWindow).getTitle();
      console.log(" title[aboutWindow]: " + n)
      var n = await app.client.window(prefsWindow).getTitle();
      console.log(" title[prefsWindow]: " + n)

      await app.client.window(mainWindow)

  });

  it('click Help/About', async () => {
    await app.client
      .click('#helpButton')
      .waitForVisible("#aboutMenuItem")
      .click('#aboutMenuItem')

      .window(aboutWindow)
      //      .switchWindow('About Synergize')
      .saveScreenshot('screenshots-aboutWindow.png')
      .getTitle()
      .should.eventually.equal('About Synergize')
  });

  it('show main window', async () => {
    await app.client
      .switchWindow('Synergize')
      .getTitle()
      .should.eventually.equal('Synergize')
  });

  //  it('should get a url', async() => {
  //    await app.client.url(config.url)
  //      .getTitle()
  //      .should.eventually.include('DuckDuckGo');
  //  });
  //
  //  it('should search', async() => {
  //    const input = 'this is a test';
  //    await app.client.url(config.url)
  //      .setValue(SearchPage.searchField, input)
  //      .getValue(SearchPage.searchField)
  //      .should.eventually.equal(input)
  //      .click(SearchPage.searchButton)
  //      .element(SearchPage.searchResult)
  //      .should.eventually.exist;
  //  });

});
