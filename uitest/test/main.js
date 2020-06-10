const hooks = require('./hooks');
const config = require('../config').get(process.env.NODE_ENV);
//const SearchPage = require('./page-objects/search.page');

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
  });

  it('click Help/About', async () => {
    await app.client
      .click('#helpButton')
      .click('#aboutMenuItem')
      .switchWindow('About Synergize')
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
