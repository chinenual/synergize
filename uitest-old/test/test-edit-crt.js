const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...
const TYPING_PAUSE = 500; // slow down typing just a bit to reduce stress on Synergy for non-debounced typing to separate fields

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test CRT edits', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    describe('Load INTERNAL CRT', () => {
        it('click load INTERNAL', async () => {
            const link = await app.client.$('.file=INTERNAL')
            await link.click()
            
            const crt_path = await app.client.$('#crt_path')
            await app.client.waitUntilTextExists('#crt_path', 'INTERNAL');

            (await crt_path.getText()).should.equal('INTERNAL');

            const clearButton = await app.client.$('#crt_tables > div:nth-child(1) > table > tbody > tr:nth-child(2) > td > span.crtSlotClearButton');
            (await clearButton.isDisplayed()).should.equal(false)
        });
        it('sanity check slots', async () => {
            const crt_voicename_1 = await app.client.$('#crt_voicename_1');
            const crt_voicename_2 = await app.client.$('#crt_voicename_2');
            const crt_voicename_3 = await app.client.$('#crt_voicename_3');
            (await crt_voicename_1.getText()).should.equal('G7S');
            (await crt_voicename_2.getText()).should.equal('HORNSXX');
            (await crt_voicename_3.getText()).should.equal('RRHODES');
        });
        it('toggle edit mode', async () => {
            const edit = await app.client.$('#editCRTButton')
            await edit.click()

            const clearButton = await app.client.$('#crt_tables > div:nth-child(1) > table > tbody > tr:nth-child(2) > td > span.crtSlotClearButton')
            await clearButton.waitForDisplayed();
            (await clearButton.isDisplayed()).should.equal(true)
        });
        it('sanity check slots', async () => {
            const crt_voicename_1 = await app.client.$('#crt_voicename_1');
            const crt_voicename_2 = await app.client.$('#crt_voicename_2');
            const crt_voicename_3 = await app.client.$('#crt_voicename_3');
            (await crt_voicename_1.getText()).should.equal('G7S');
            (await crt_voicename_2.getText()).should.equal('HORNSXX');
            (await crt_voicename_3.getText()).should.equal('RRHODES');
        });
    });
    describe('Edit slots', () => {
        it('clear slot 2', async () => {
            const clearButton = await app.client.$('#crt_tables > div:nth-child(1) > table > tbody > tr:nth-child(2) > td > span.crtSlotClearButton')
            await clearButton.click()
            const crt_voicename_2 = await app.client.$('#crt_voicename_2');
            (await crt_voicename_2.getText()).should.equal('')
        });
        /*
         * Can't call native file dialogs from Spectron.  FIXME: add mocks
         
         it('alter slot 3', async () => {
         await app.client
         .click('#crt_tables > div:nth-child(1) > table > tbody > tr:nth-child(3) > td > span.crtSlotAddButton')
         .waitUntilTextExists("#crt_voicename_3", 'GUITAR2A', LOAD_VCE_TIMEOUT)
         .getText('#crt_voicename_3').should.eventually.equal('GUITAR2A')
         });
        */
    });

    it('screenshot', async () => {
        await hooks.screenshotAndCompare(app, `INITVRAM-after-edit-crt`) 
    });
});
