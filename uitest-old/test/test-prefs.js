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
        const button = await app.client.$('#helpButton')
        await button.click()

        const item = await app.client.$('#preferencesMenuItem')
        await item.waitForDisplayed()

        await item.click()

//        await app.client
//            .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable. 

        await app.client.switchWindow('Synergize Preferences');
        await hooks.screenshotAndCompare(app, 'prefsWindow');

        (await app.client.getTitle()).should.equal('Synergize Preferences')

        const txt = await app.client.$('#libraryPath')
        await txt.setValue('../data/testfiles')

        const submit = await app.client.$('button[type=submit]')
        await submit.click()

//        await app.client
//            .pause(WINDOW_PAUSE) // HACK: but without this switching windows is unreliable.

        await app.client.switchWindow('Synergize');
        (await app.client.getTitle()).should.equal('Synergize')

        const txt2 = await app.client.$('#path');

        (await txt2.getText()).should.equal('testfiles')
    });


    it('show main window', async () => {
//        await app.client
//            .pause(WINDOW_PAUSE); // HACK: but without this switching windows is unreliable. 

        await app.client.switchWindow('Synergize');
        (await app.client.getTitle()).should.equal('Synergize')
    });

});

