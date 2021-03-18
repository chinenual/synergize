const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...
const TYPING_PAUSE = 500; // slow down typing just a bit to reduce stress on Synergy for non-debounced typing to separate fields

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test keyeq page edits', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    it('voice tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceVoiceTab']`)
        await tab.click()
        await tab.getAttribute('class').should.eventually.include('active')
        const table = await app.client.$('#voiceParamTable')
        await table.waitForDisplayed()
        await table.isDisplayed().should.eventually.equal(true)
    });

    it('click load G7S', async () => {
        
        const vname = await app.client.$('#VNAME')
        // need to clear this since previous test may also be using same voice
        await vname.clearValue()
        
        const link = await app.client.$('.file=G7S')
        await link.click()
                
        const confirmText = await app.client.$('#confirmText')

        await confirmText.waitForDisplayed()
        await confirmText.getText().should.eventually.include('pending edits')
        
        const confirmOk = await app.client.$('#confirmOKButton')
        await confirmOk.click()
        await confirmText.waitForDisplayed({reverse: true})  // wait to disappear
        
        app.client.waitUntil(
            () => vname.getText() == 'G7S',
            {
                timeout: LOAD_VCE_TIMEOUT
            }
        );       

        await vname.getValue().should.eventually.equal('G7S')

    });

    it('keyeq tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceKeyEqTab']`)
        await tab.click()
        await tab.getAttribute('class').should.eventually.include('active')
        const table = await app.client.$('#keyEqTable')
        await table.waitForDisplayed()
        await table.isDisplayed().should.eventually.equal(true)
    });

    // element 1-5's initial value is -24 (min)
    // element 16-17's initial value is 7 (max) // DK MAN file says -24..6, but SYNHCS supports -24..7

    // test typing value directly - both in and out of range
    it('type to element 9 via setvalue - 0->-24', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[9]'))
        const otherele = await app.client.$(cssQuoteId('#keyeq[1]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.setValue('-24')
        await otherele.click() // click in a different input to force onchange
        await ele.getValue().should.eventually.equal('-24')
    });
    it('type to element 10 via setvalue - 1->-24', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[10]'))
        const otherele = await app.client.$(cssQuoteId('#keyeq[1]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.setValue('-24')
        await otherele.click() // click in a different input to force onchange
        await ele.getValue().should.eventually.equal('-24')
    });
    it('type to element 11 via setvalue - 0->100 - above range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[11]'))
        const otherele = await app.client.$(cssQuoteId('#keyeq[1]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.setValue('100')
        await otherele.click() // click in a different input to force onchange
        await ele.getValue().should.eventually.equal('7')
    });
    it('type to element 12 via setvalue - 9->-100 - below range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[12]'))
        const otherele = await app.client.$(cssQuoteId('#keyeq[1]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.setValue('-100')
        await otherele.click() // click in a different input to force onchange
        await ele.getValue().should.eventually.equal('-24')
    });


    // Test up arrow, down arrow - both in and out of range
    it('up-arrow to element 1 - -4->-3', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[1]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.click()
        await app.client.keys('ArrowUp')
        await ele.getValue().should.eventually.equal('-3')
    });
    it('up-arrow to element 16 - 7->8 at limit', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[16]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.click()
        await app.client.keys('ArrowUp')
        await ele.getValue().should.eventually.equal('7')
    });
    it('down-arrow to element 19 - 6->5', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[19]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.click()
        await app.client.keys('ArrowDown')
        await ele.getValue().should.eventually.equal('5')
    });
    it('down-arrow to element 10 - -24->-24 at limit', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[10]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.click()
        await app.client.keys('ArrowDown')
        await ele.getValue().should.eventually.equal('-24')
    });

    // test spinner buttons
    it('button-down to element 9 -24->-25 - below range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[9]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.moveTo()
        const arrow = await app.client.$(cssQuoteId('#keyeq[9] ~ span button.bootstrap-touchspin-down'))
        await arrow.isDisplayed().should.eventually.be.true
        arrow.click()        
        await ele.getValue().should.eventually.equal('-24')
    });
    it('button-down to element 22 2->1 - in range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[22]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.moveTo()
        const arrow = await app.client.$(cssQuoteId('#keyeq[22] ~ span button.bootstrap-touchspin-down'))
        await arrow.isDisplayed().should.eventually.be.true
        arrow.click()        
        await ele.getValue().should.eventually.equal('1')        
    });

    it('button-up to element 23 2->3 - in range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[23]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.moveTo()
        const arrow = await app.client.$(cssQuoteId('#keyeq[23] ~ span button.bootstrap-touchspin-up'))
        await arrow.isDisplayed().should.eventually.be.true
        arrow.click()        
        await ele.getValue().should.eventually.equal('3') 
    });
    it('button-up to element 17 7->8 - above range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyeq[17]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.moveTo()
        const arrow = await app.client.$(cssQuoteId('#keyeq[17] ~ span button.bootstrap-touchspin-up'))
        await arrow.isDisplayed().should.eventually.be.true
        arrow.click()        
        await ele.getValue().should.eventually.equal('7') 
    });

    it('screenshot', async () => {
        await hooks.screenshotAndCompare(app, `G7S-after-edit-keyeqTab`)
    });

});
