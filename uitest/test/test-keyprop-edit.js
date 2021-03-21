const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...
const TYPING_PAUSE = 500; // slow down typing just a bit to reduce stress on Synergy for non-debounced typing to separate fields

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test keyprop page edits', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    it('voice tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceVoiceTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const table = await app.client.$('#voiceParamTable')
        await table.waitForDisplayed();
        (await table.isDisplayed()).should.equal(true);
    });

    it('click load G7S', async () => {
        
        const vname = await app.client.$('#VNAME')
        // need to clear this since previous test may also be using same voice
        await vname.clearValue()
        
        const link = await app.client.$('.file=G7S')
        await link.click()
                
        const confirmText = await app.client.$('#confirmText')

        await confirmText.waitForDisplayed();
        (await confirmText.getText()).should.include('pending edits');
        
        const confirmOk = await app.client.$('#confirmOKButton')
        await confirmOk.click()
        await confirmText.waitForDisplayed({reverse: true})  // wait to disappear
        
        app.client.waitUntil(
            () => vname.getText() == 'G7S',
            {
                timeout: LOAD_VCE_TIMEOUT
            }
        );       

        (await vname.getValue()).should.equal('G7S');

    });

    it('keyprop tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceKeyPropTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const table = await app.client.$('#keyPropTable')
        await table.waitForDisplayed();
        (await table.isDisplayed()).should.equal(true);
        
    });

    // element 1-4's initial value is 0 (min)
    // element 20-24's initial value is 32 (max)

    // Test up arrow, down arrow - both in and out of range
    it('up-arrow to element 1 - 0->1', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[1]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.click()
        await app.client.keys('ArrowUp');
        (await ele.getValue()).should.equal('1');
    });
    it('up-arrow to element 24 - 32->32 at limit', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[24]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.click()
        await app.client.keys('ArrowUp');
        (await ele.getValue()).should.equal('32');
    });
    it('down-arrow to element 23 - 32->31', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[24]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.click()
        await app.client.keys('ArrowDown');
        (await ele.getValue()).should.equal('31');
    });
    it('down-arrow to element 2 - 0->0 at limit', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[2]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.click()
        await app.client.keys('ArrowDown');
        (await ele.getValue()).should.equal('0');
    });


    // test typing value directly - both in and out of range
    it('type to element 10 via setvalue - 0->18', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[10]'))
        const otherele = await app.client.$(cssQuoteId('#keyprop[1]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.setValue('18')
        await otherele.click(); // click in a different input to force onchange
        (await ele.getValue()).should.equal('18');
    });
    it('type to element 11 via setvalue - 0->100 - above range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[11]'))
        const otherele = await app.client.$(cssQuoteId('#keyprop[1]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.setValue('100')
        await otherele.click(); // click in a different input to force onchange
        (await ele.getValue()).should.equal('32');
    });
    it('type to element 12 via setvalue - 9->-100 - below range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[12]'))
        const otherele = await app.client.$(cssQuoteId('#keyprop[1]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.setValue('-100')
        await otherele.click(); // click in a different input to force onchange
        (await ele.getValue()).should.equal('0');
    });

    // test spinner buttons
    it('button-down to element 3 0->-1 - below range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[3]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.moveTo()
        const arrow = await app.client.$(cssQuoteId('#keyprop[3] ~ span button.bootstrap-touchspin-down'));
        (await arrow.isDisplayed()).should.be.true;
        arrow.click();        
        (await ele.getValue()).should.equal('0');
    });
    it('button-down to element 14 11->10 - in range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[14]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.moveTo()
        const arrow = await app.client.$(cssQuoteId('#keyprop[14] ~ span button.bootstrap-touchspin-down'));
        (await arrow.isDisplayed()).should.be.true;
        arrow.click();        
        (await ele.getValue()).should.equal('10');
    });

    it('button-up to element 15 13->14 - in range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[15]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.moveTo()
        const arrow = await app.client.$(cssQuoteId('#keyprop[15] ~ span button.bootstrap-touchspin-up'));
        (await arrow.isDisplayed()).should.be.true;
        arrow.click();
        (await ele.getValue()).should.equal('14');
    });
    it('button-up to element 21 32->33 - above range', async () => {
        const ele = await app.client.$(cssQuoteId('#keyprop[21]'))
        await app.client.pause(TYPING_PAUSE)
        await ele.moveTo()
        const arrow = await app.client.$(cssQuoteId('#keyprop[21] ~ span button.bootstrap-touchspin-up'));
        (await arrow.isDisplayed()).should.be.true;
        arrow.click();
        (await ele.getValue()).should.equal('32');
    });

    it('screenshot', async () => {
        await hooks.screenshotAndCompare(app, `G7S-after-edit-keypropTab`)
    });

});
