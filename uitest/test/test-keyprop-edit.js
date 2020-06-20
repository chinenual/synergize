const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test keyprop page edits', () => {
    before(async () => {
        console.log("====== reuse the app");
        app = await hooks.getApp();
    });

    it('click load G7S', async () => {
        await app.client
            .click('.file=G7S')
            .waitUntilTextExists("#name", 'G7S', LOAD_VCE_TIMEOUT)
            .getText('#name').should.eventually.equal('G7S')
    });
    it('keyprop tab should display', async () => {
        await app.client
            .click(`#vceTabs a[href='#vceKeyPropTab']`)
            .getAttribute(`#vceTabs a[href='#vceKeyPropTab']`, 'class').should.eventually.include('active')
            .waitForVisible('#keyPropTable')
    });

    // element 1-4's initial value is 0 (min)
    // element 20-24's initial value is 32 (max)

    // Test up arrow, down arrow - both in and out of range
    it('up-arrow to element 1 - 0->1', async () => {
        await app.client
            .click(cssQuoteId('#keyprop[1]')).keys('ArrowUp')
            .getValue(cssQuoteId('#keyprop[1]')).should.eventually.equal('1')
    });
    it('up-arrow to element 24 - 32->32 at limit', async () => {
        await app.client
            .click(cssQuoteId('#keyprop[24]')).keys('ArrowUp')
            .getValue(cssQuoteId('#keyprop[24]')).should.eventually.equal('32')
    });
    it('down-arrow to element 23 - 32->31', async () => {
        await app.client
            .click(cssQuoteId('#keyprop[24]')).keys('ArrowDown')
            .getValue(cssQuoteId('#keyprop[24]')).should.eventually.equal('31')
    });
    it('down-arrow to element 2 - 0->0 at limit', async () => {
        await app.client
            .click(cssQuoteId('#keyprop[2]')).keys('ArrowDown')
            .getValue(cssQuoteId('#keyprop[2]')).should.eventually.equal('0')
    });


    // test typing value directly - both in and out of range
    it('type to element 10 via setvalue - 0->18', async () => {
        await app.client
            .clearElement(cssQuoteId('#keyprop[10]'))
            .setValue(cssQuoteId('#keyprop[10]'), '18')
            .click(cssQuoteId('#keyprop[1]')) // click in a different input to force onchange
            .getValue(cssQuoteId('#keyprop[10]')).should.eventually.equal('18')
    });
    it('type to element 11 via setvalue - 0->100 - above range', async () => {
        await app.client
            .clearElement(cssQuoteId('#keyprop[11]'))
            .setValue(cssQuoteId('#keyprop[11]'), '100')
            .click(cssQuoteId('#keyprop[1]')) // click in a different input to force onchange
            .getValue(cssQuoteId('#keyprop[11]')).should.eventually.equal('32')
    });
    it('type to element 12 via setvalue - 9->-100 - below range', async () => {
        await app.client
            .clearElement(cssQuoteId('#keyprop[12]'))
            .setValue(cssQuoteId('#keyprop[12]'), '-100')
            .click(cssQuoteId('#keyprop[1]')) // click in a different input to force onchange
            .getValue(cssQuoteId('#keyprop[12]')).should.eventually.equal('0')
    });

    // test spinner buttons
    it('button-down to element 3 0->-1 - below range', async () => {
        await app.client
            .moveToObject(cssQuoteId('#keyprop[3]'))
            .isVisible(cssQuoteId('#keyprop[3] ~ span button.bootstrap-touchspin-down')).should.eventually.be.true
            .click(cssQuoteId('#keyprop[3] ~ span button.bootstrap-touchspin-down'))
            .getValue(cssQuoteId('#keyprop[3]')).should.eventually.equal('0')
    });
    it('button-down to element 14 11->10 - in range', async () => {
        await app.client
            .moveToObject(cssQuoteId('#keyprop[14]'))
            .isVisible(cssQuoteId('#keyprop[14] ~ span button.bootstrap-touchspin-down')).should.eventually.be.true
            .click(cssQuoteId('#keyprop[14] ~ span button.bootstrap-touchspin-down'))
            .getValue(cssQuoteId('#keyprop[14]')).should.eventually.equal('10')
    });

    it('button-up to element 15 13->14 - in range', async () => {
        await app.client
            .moveToObject(cssQuoteId('#keyprop[15]'))
            .isVisible(cssQuoteId('#keyprop[15] ~ span button.bootstrap-touchspin-up')).should.eventually.be.true
            .click(cssQuoteId('#keyprop[15] ~ span button.bootstrap-touchspin-up'))
            .getValue(cssQuoteId('#keyprop[15]')).should.eventually.equal('14')
    });
    it('button-up to element 21 32->33 - above range', async () => {
        await app.client
            .moveToObject(cssQuoteId('#keyprop[21]'))
            .isVisible(cssQuoteId('#keyprop[21] ~ span button.bootstrap-touchspin-up')).should.eventually.be.true
            .click(cssQuoteId('#keyprop[21] ~ span button.bootstrap-touchspin-up'))
            .getValue(cssQuoteId('#keyprop[21]')).should.eventually.equal('32')
    });

    it('screenshot', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `G7S-after-edit-keypropTab`) })
    });
    it('capture renderer logs', async () => {
        await app.client.getRenderProcessLogs().then(function (logs) {
            logs.forEach(function (log) {
                console.log("RENDERER: " + log.level + ": " + log.source + " : " + log.message);
            });
        });
    });

});
