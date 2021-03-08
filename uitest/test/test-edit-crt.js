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
            await app.client
                .click('.file=INTERNAL')
                .waitUntilTextExists("#crt_path", 'INTERNAL', LOAD_VCE_TIMEOUT)
                .getText('#crt_path').should.eventually.equal('INTERNAL')
                .isVisible('#crt_tables > div:nth-child(1) > table > tbody > tr:nth-child(2) > td > span.crtSlotClearButton').should.eventually.equal(false)
        });
        it('sanity check slots', async () => {
            await app.client
                .getText('#crt_voicename_1').should.eventually.equal('G7S')
                .getText('#crt_voicename_2').should.eventually.equal('HORNSXX')
                .getText('#crt_voicename_3').should.eventually.equal('RRHODES')
        });
        it('toggle edit mode', async () => {
            await app.client
                .click('#editCRTButton')
                .waitForVisible('#crt_tables > div:nth-child(1) > table > tbody > tr:nth-child(2) > td > span.crtSlotClearButton')
                .isVisible('#crt_tables > div:nth-child(1) > table > tbody > tr:nth-child(2) > td > span.crtSlotClearButton').should.eventually.equal(true)
        });
        it('sanity check slots', async () => {
            await app.client
                .getText('#crt_voicename_1').should.eventually.equal('G7S')
                .getText('#crt_voicename_2').should.eventually.equal('HORNSXX')
                .getText('#crt_voicename_3').should.eventually.equal('RRHODES')
        });
    });
    describe('Edit slots', () => {
        it('clear slot 2', async () => {
            await app.client
                .click('#crt_tables > div:nth-child(1) > table > tbody > tr:nth-child(2) > td > span.crtSlotClearButton')
                .waitUntilTextExists("#crt_voicename_2", '', LOAD_VCE_TIMEOUT)
                .getText('#crt_voicename_2').should.eventually.equal('')
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
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `INITVRAM-after-edit-crt`) })
    });
});
