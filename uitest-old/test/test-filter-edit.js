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
        const tab = await app.client.$(`#vceTabs a[href='#vceFiltersTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const ele = await app.client.$('#filterSelect')
        await ele.waitForDisplayed()
    });

    it('select Af', async () => {
        const ele = await app.client.$('#filterSelect')
        await ele.selectByVisibleText('Af')
        const table = await app.client.$('#filterTable')
        await table.waitForDisplayed()
    });


    // test that all the spinner text conversions work at the right ranges
    describe('filter values', () => {
        it('type to flt[1] to and past -64', async () => {
            const ele = await app.client.$(cssQuoteId('#flt[1]'))
            const otherele = await app.client.$(cssQuoteId('#flt[2]'))
            await app.client.pause(TYPING_PAUSE)
            await ele.setValue('-63')
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('-63');

            await app.client.pause(TYPING_PAUSE)
            await ele.click()
            await app.client.keys('ArrowDown');
            (await ele.getValue()).should.equal('-64');
            
            await app.client.pause(TYPING_PAUSE)
            await ele.click()
            await app.client.keys('ArrowDown');
            (await ele.getValue()).should.equal('-64')            ;
        });
        it('type to flt[2] to and past 63', async () => {
            const ele = await app.client.$(cssQuoteId('#flt[2]'))
            const otherele = await app.client.$(cssQuoteId('#flt[1]'))
            await app.client.pause(TYPING_PAUSE)
            await ele.setValue('62')
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('62');

            await app.client.pause(TYPING_PAUSE)
            await ele.click()
            await app.client.keys('ArrowUp');
            (await ele.getValue()).should.equal('63');
            
            await app.client.pause(TYPING_PAUSE)
            await ele.click()
            await app.client.keys('ArrowUp');
            (await ele.getValue()).should.equal('63');
        });
    });

    describe('copy filter', () => {
        it('check Af initial conditions', async () => {
            const filterSelect = await app.client.$('#filterSelect')
            await filterSelect.selectByVisibleText('Af')
            const filterTable = await app.client.$('#filterTable')
            await filterTable.waitForDisplayed()

            const flt1 = await app.client.$(cssQuoteId('#flt[1]'))
            const flt2 = await app.client.$(cssQuoteId('#flt[2]'));
            (await flt1.getValue()).should.equal('-64');
            (await flt2.getValue()).should.equal('63');
        });

        it('switch to Bf 2', async () => {
            const filterSelect = await app.client.$('#filterSelect')
            await filterSelect.selectByVisibleText('Bf 2')
            const filterTable = await app.client.$('#filterTable')
            await filterTable.waitForDisplayed()

            const flt1 = await app.client.$(cssQuoteId('#flt[1]'))
            const flt2 = await app.client.$(cssQuoteId('#flt[2]'));
            (await flt1.getValue()).should.equal('0');
            (await flt2.getValue()).should.equal('0');
        });

        it('copy from 1', async () => {
            const filterSelect = await app.client.$('#filterSelect')
            await filterSelect.selectByVisibleText('Af')
            const filterTable = await app.client.$('#filterTable')
            await filterTable.waitForDisplayed()

            const flt1 = await app.client.$(cssQuoteId('#flt[1]'))
            const flt2 = await app.client.$(cssQuoteId('#flt[2]'));
            (await flt1.getValue()).should.equal('-64');
            (await flt2.getValue()).should.equal('63');
        });

    });

    it('screenshot', async () => {
        await hooks.screenshotAndCompare(app, `INITVRAM-after-edit-filters`) 
    });
});
