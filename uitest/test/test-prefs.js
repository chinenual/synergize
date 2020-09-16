const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;

let app;

describe('Check initial preferences', () => {
  afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this, app); });
  before(async () => {
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
      .then(() => { return hooks.screenshotAndCompare(app, 'prefsWindow') })

      .getTitle().should.eventually.equal('Synergize Preferences')

        .$('#libraryPath').setValue('../data/testfiles')
        .$('#oscAutoConfig').click() // set to off
        .$('#vstAutoConfig').click() // set to off

        .click('button[type=submit]')
        .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable.

      .switchWindow('Synergize')
      .getTitle().should.eventually.equal('Synergize')

      .then(() => { return hooks.screenshotAndCompare(app, 'mainAfterSetLibraryPath') })

      .getText('#path').should.eventually.equal('testfiles')
  });


  it('show main window', async () => {
    await app.client
      .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable. 

      .switchWindow('Synergize')
      .getTitle().should.eventually.equal('Synergize')
  });

});

