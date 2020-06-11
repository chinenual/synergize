const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;

let app;

describe('Check initial preferences', () => {
    before(async () => {
        console.log("====== reuse the app");
        app = await hooks.getApp();
      });
    
      it('click Help/Preferences', async () => {
        await app.client
          .waitUntilWindowLoaded()
          .click('#helpButton')
          .waitForVisible("#preferencesMenuItem")
          .click('#preferencesMenuItem')
    
          .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable. 
    
          .switchWindow('Synergize Preferences')
          .saveScreenshot('screenshots-prefsWindow.png')
          .getTitle().should.eventually.equal('Synergize Preferences')

          .$('#libraryPath').setValue('../data/testfiles')
          .click('button[type=submit]')

          .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable. 
  
          .switchWindow('Synergize')
          .getTitle().should.eventually.equal('Synergize')

          .saveScreenshot('screenshots-mainAfterSetLibraryPath.png')

          .getText('#path').should.eventually.equal('testfiles')
        });
    
    
    it('show main window', async () => {
      await app.client
        .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable. 
  
        .switchWindow('Synergize')
        .getTitle().should.eventually.equal('Synergize')
    });
  
  });
  
  