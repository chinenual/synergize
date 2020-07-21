const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...
const TYPING_PAUSE = 500; // slow down typing just a bit to reduce stress on Synergy for non-debounced typing to separate fields

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test filter page edits', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    it('filter tab should display', async () => {
        await app.client
            .click(`#vceTabs a[href='#vceFiltersTab']`)
            .getAttribute(`#vceTabs a[href='#vceFiltersTab']`, 'class').should.eventually.include('active')
            .waitForVisible('#filterSelect')
    });

    it('select Af', async () => {
        await app.client
            .selectByVisibleText(cssQuoteId('#filterSelect'), 'Af')
            .waitForVisible('#filterTable')

    });


    // test that all the spinner text conversions work at the right ranges
    describe('filter values', () => {
        it('type to flt[1] to and past -64', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#flt[1]'))
                .setValue(cssQuoteId('#flt[1]'), '-63')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#flt[2]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#flt[1]')).should.eventually.equal('-63')
                // should not be able to go below min
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#flt[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#flt[1]')).should.eventually.equal('-64')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#flt[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#flt[1]')).should.eventually.equal('-64')
        });
        it('type to flt[2] to and past 63', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#flt[2]'))
                .setValue(cssQuoteId('#flt[2]'), '1')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#flt[1]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#flt[2]')).should.eventually.equal('62')
                // should not be able to go above max
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#flt[2]')).keys('ArrowUp')
                .getValue(cssQuoteId('#flt[2]')).should.eventually.equal('63')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#flt[2]')).keys('ArrowUp')
                .getValue(cssQuoteId('#flt[2]')).should.eventually.equal('63')

        });
    });

    it('screenshot', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `INITVRAM-after-edit-filters`) })
    });
});
