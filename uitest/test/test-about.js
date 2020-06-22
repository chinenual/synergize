const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;

let app;

describe('About window', () => {
    afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this,app); });
    before(async () => {
        console.log("====== reuse the app");
        app = await hooks.getApp();
      });
    
    it('click Help/About', async () => {
      await app.client
        .waitUntilWindowLoaded()
        .click('#helpButton')
        .waitForVisible("#aboutMenuItem")
        .click('#aboutMenuItem')
  
        .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable. 
  
        .switchWindow('About Synergize')
        .then(() => {return hooks.screenshotAndCompare(app, 'aboutWindow')})

        .getTitle().should.eventually.equal('About Synergize')
    });
  
    /***
     * Can't seem to really test this via webdriver.  client.close() closes the window, but then
     * webdriver loses its handle to it - so it thinks it's not around any more. (that or the hook 
     * that causes the close to be interprted as a hide is bypassed.  in any case, dont try to test this)
     *
    it('close and reopen About', async () => {
      await app.client
      .close()
  
      .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable. 
  
      .switchWindow('Synergize')
      .getTitle()
      .should.eventually.equal('Synergize')
      .click('#helpButton')
      .waitForVisible("#aboutMenuItem")
      .click('#aboutMenuItem')
  
      .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable. 
  
      .switchWindow('About Synergize')
      .saveScreenshot('screenshots-aboutWindow.png')
      .getTitle()
      .should.eventually.equal('About Synergize')
  });
  
  *
  ***/
  
    it('show main window', async () => {
      await app.client
        .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable. 
  
        .switchWindow('Synergize')
        .getTitle().should.eventually.equal('Synergize')
    });
  
  });
  
  