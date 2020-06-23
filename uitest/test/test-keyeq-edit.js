const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test keyeq page edits', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    it('voice tab should display', async () => {
        await app.client
            .click(`#vceTabs a[href='#vceVoiceTab']`)
            .getAttribute(`#vceTabs a[href='#vceVoiceTab']`, 'class').should.eventually.include('active')
            .waitForVisible('#voiceParamTable')
            .waitUntil(() => {
                return app.client.$('#voiceParamTable').isVisible()
            })
            .isVisible('#voiceParamTable').should.eventually.equal(true)
    });

    it('click load G7S', async () => {
        await app.client
        // need to clear this since previous test may also be using same voice
            .clearElement("#VNAME")
            .click('.file=G7S')

            .waitForVisible('#confirmText')
            .getText('#confirmText').should.eventually.include('pending edits')
            .click('#confirmOKButton')
            .waitForVisible('#confirmText', 1000, true) // wait to disappear

            .waitUntilTextExists("#vce_name", 'G7S', LOAD_VCE_TIMEOUT)

            .getValue('#VNAME').should.eventually.equal('G7S')
    });
    it('keyeq tab should display', async () => {
        await app.client
            .click(`#vceTabs a[href='#vceKeyEqTab']`)
            .getAttribute(`#vceTabs a[href='#vceKeyEqTab']`, 'class').should.eventually.include('active')
            .waitForVisible('#keyEqTable')
    });

    // element 1-5's initial value is -24 (min)
    // element 16-17's initial value is 7 (max) // DK MAN file says -24..6, but SYNHCS supports -24..7

    // test typing value directly - both in and out of range
    it('type to element 9 via setvalue - 0->-24', async () => {
        await app.client
            .clearElement(cssQuoteId('#keyeq[9]'))
            .setValue(cssQuoteId('#keyeq[9]'), '-24')
            .click(cssQuoteId('#keyeq[1]')) // click in a different input to force onchange
            .getValue(cssQuoteId('#keyeq[9]')).should.eventually.equal('-24')
    });
    it('type to element 10 via setvalue - 1->-24', async () => {
        await app.client
            .clearElement(cssQuoteId('#keyeq[10]'))
            .setValue(cssQuoteId('#keyeq[10]'), '-24')
            .click(cssQuoteId('#keyeq[1]')) // click in a different input to force onchange
            .getValue(cssQuoteId('#keyeq[10]')).should.eventually.equal('-24')
    });
    it('type to element 11 via setvalue - 0->100 - above range', async () => {
        await app.client
            .clearElement(cssQuoteId('#keyeq[11]'))
            .setValue(cssQuoteId('#keyeq[11]'), '100')
            .click(cssQuoteId('#keyeq[1]')) // click in a different input to force onchange
            .getValue(cssQuoteId('#keyeq[11]')).should.eventually.equal('7')
    });
    it('type to element 12 via setvalue - 9->-100 - below range', async () => {
        await app.client
            .clearElement(cssQuoteId('#keyeq[12]'))
            .setValue(cssQuoteId('#keyeq[12]'), '-100')
            .click(cssQuoteId('#keyeq[1]')) // click in a different input to force onchange
            .getValue(cssQuoteId('#keyeq[12]')).should.eventually.equal('-24')
    });


    // Test up arrow, down arrow - both in and out of range
    it('up-arrow to element 1 - -4->-3', async () => {
        await app.client
            .click(cssQuoteId('#keyeq[1]')).keys('ArrowUp')
            .getValue(cssQuoteId('#keyeq[1]')).should.eventually.equal('-3')
    });
    it('up-arrow to element 16 - 7->8 at limit', async () => {
        await app.client
            .click(cssQuoteId('#keyeq[16]')).keys('ArrowUp')
            .getValue(cssQuoteId('#keyeq[16]')).should.eventually.equal('7')
    });
    it('down-arrow to element 19 - 6->5', async () => {
        await app.client
            .click(cssQuoteId('#keyeq[19]')).keys('ArrowDown')
            .getValue(cssQuoteId('#keyeq[19]')).should.eventually.equal('5')
    });
    it('down-arrow to element 10 - -24->-24 at limit', async () => {
        await app.client
            .click(cssQuoteId('#keyeq[10]')).keys('ArrowDown')
            .getValue(cssQuoteId('#keyeq[10]')).should.eventually.equal('-24')
    });

    // test spinner buttons
    it('button-down to element 9 -24->-25 - below range', async () => {
        await app.client
            .moveToObject(cssQuoteId('#keyeq[9]'))
            .isVisible(cssQuoteId('#keyeq[9] ~ span button.bootstrap-touchspin-down')).should.eventually.be.true
            .click(cssQuoteId('#keyeq[9] ~ span button.bootstrap-touchspin-down'))
            .getValue(cssQuoteId('#keyeq[9]')).should.eventually.equal('-24')
    });
    it('button-down to element 22 2->1 - in range', async () => {
        await app.client
            .moveToObject(cssQuoteId('#keyeq[22]'))
            .isVisible(cssQuoteId('#keyeq[22] ~ span button.bootstrap-touchspin-down')).should.eventually.be.true
            .click(cssQuoteId('#keyeq[22] ~ span button.bootstrap-touchspin-down'))
            .getValue(cssQuoteId('#keyeq[22]')).should.eventually.equal('1')
    });

    it('button-up to element 23 2->3 - in range', async () => {
        await app.client
            .moveToObject(cssQuoteId('#keyeq[23]'))
            .isVisible(cssQuoteId('#keyeq[23] ~ span button.bootstrap-touchspin-up')).should.eventually.be.true
            .click(cssQuoteId('#keyeq[23] ~ span button.bootstrap-touchspin-up'))
            .getValue(cssQuoteId('#keyeq[23]')).should.eventually.equal('3')
    });
    it('button-up to element 17 7->8 - above range', async () => {
        await app.client
            .moveToObject(cssQuoteId('#keyeq[17]'))
            .isVisible(cssQuoteId('#keyeq[17] ~ span button.bootstrap-touchspin-up')).should.eventually.be.true
            .click(cssQuoteId('#keyeq[17] ~ span button.bootstrap-touchspin-up'))
            .getValue(cssQuoteId('#keyeq[17]')).should.eventually.equal('7')
    });

    it('screenshot', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `G7S-after-edit-keyeqTab`) })
    });

});
