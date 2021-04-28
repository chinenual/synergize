const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...
const TYPING_PAUSE = 500; // slow down typing just a bit to reduce stress on Synergy for non-debounced typing to separate fields
const REDRAW_ENV_PAUSE = 2000;

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test gain', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    it('voice tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceVoiceTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const paramTable = await app.client.$('#voiceParamTable')
        const oscTable = await app.client.$('#voiceOscTable') 
        await oscTable.waitForDisplayed()
        await paramTable.waitForDisplayed()
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
        
        const confirmOk = await app.client.$('#confirmOKButton');
        await confirmOk.click();
        await confirmText.waitForDisplayed({reverse: true})  // wait to disappear
        
        await hooks.waitUntilValueExists('#VNAME', 'G7S');
        (await vname.getValue()).should.equal('G7S');

    });
    
    it('voice tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceVoiceTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const paramTable = await app.client.$('#voiceParamTable')
        const oscTable = await app.client.$('#voiceOscTable') 
        await oscTable.waitForDisplayed()
        await paramTable.waitForDisplayed()
    });
    
    it('check initial osc gains', async() => {
        const oscGain1 = await app.client.$(cssQuoteId('#oscGain[1]'));
        const oscGain2 = await app.client.$(cssQuoteId('#oscGain[2]'));
        const oscGain3 = await app.client.$(cssQuoteId('#oscGain[3]'));
        const oscGain4 = await app.client.$(cssQuoteId('#oscGain[4]'));

        (await oscGain1.getValue()).should.equal('97');
        (await oscGain2.getValue()).should.equal('83');
        (await oscGain3.getValue()).should.equal('100');
        (await oscGain4.getValue()).should.equal('100');
    });


    it('env tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceEnvsTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active')
        const table = await app.client.$('#envTable')
        await table.waitForDisplayed();
        (await table.isDisplayed()).should.equal(true);
    });

    it('change env gains', async() => {
        const lowGain = await app.client.$(cssQuoteId('#gainAmpLow'));
        const upGain = await app.client.$(cssQuoteId('#gainAmpUp'));
        const lowVal1 = await app.client.$(cssQuoteId('#envAmpLowVal[1]'));
        const lowVal2 = await app.client.$(cssQuoteId('#envAmpLowVal[2]'));
        const upVal1 = await app.client.$(cssQuoteId('#envAmpUpVal[1]'));
        const upVal2 = await app.client.$(cssQuoteId('#envAmpUpVal[2]'));

        // check osc1 initial state:
        (await lowGain.getValue()).should.equal('97');
        (await upGain.getValue()).should.equal('97');
        (await lowVal1.getValue()).should.equal('70');
        (await lowVal2.getValue()).should.equal('69'); 
        (await upVal1.getValue()).should.equal('70');
        (await upVal2.getValue()).should.equal('69'); 

        // set Low to 0, check values: 0,0 (Up unchanged)
        await app.client.pause(TYPING_PAUSE);
        await lowGain.setValue('0');        
        await upGain.click(); // click off the element to force the onchange event
        await app.client.pause(REDRAW_ENV_PAUSE);
        await hooks.waitUntilValueExists(cssQuoteId('#envAmpLowVal[1]'), '0');
        
        (await lowGain.getValue()).should.equal('0');
        (await upGain.getValue()).should.equal('97');
        (await lowVal1.getValue()).should.equal('0');
        (await lowVal2.getValue()).should.equal('0'); 
        (await upVal1.getValue()).should.equal('70');
        (await upVal2.getValue()).should.equal('69'); 

        // set Low to 50, check values: 36,35 (Up unchanged)
        await app.client.pause(TYPING_PAUSE);
        await lowGain.setValue('50'); 
        await upGain.click(); // click off the element to force the onchange event
        await app.client.pause(REDRAW_ENV_PAUSE);
        await hooks.waitUntilValueExists(cssQuoteId('#envAmpLowVal[1]'), '36');
        
        (await lowGain.getValue()).should.equal('50');
        (await upGain.getValue()).should.equal('97');
        (await lowVal1.getValue()).should.equal('36');
        (await lowVal2.getValue()).should.equal('35'); 
        (await upVal1.getValue()).should.equal('70');
        (await upVal2.getValue()).should.equal('69');
        
        // set Up to 75,  check values: 54, 53 (low unchanged)
        await app.client.pause(TYPING_PAUSE);
        await upGain.setValue('75'); 
        await lowGain.click(); // click off the element to force the onchange event
        await app.client.pause(REDRAW_ENV_PAUSE);
        await hooks.waitUntilValueExists(cssQuoteId('#envAmpUpVal[1]'), '54');
        
        (await lowGain.getValue()).should.equal('50');
        (await upGain.getValue()).should.equal('75');
        (await lowVal1.getValue()).should.equal('36');
        (await lowVal2.getValue()).should.equal('35'); 
        (await upVal1.getValue()).should.equal('54');
        (await upVal2.getValue()).should.equal('53');
    });


    it('voice tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceVoiceTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const paramTable = await app.client.$('#voiceParamTable')
        const oscTable = await app.client.$('#voiceOscTable') 
        await oscTable.waitForDisplayed()
        await paramTable.waitForDisplayed()
    });

    it('change osc gains', async() => {
        // osc1 should now say 75
        // set it to 90
        const oscGain1 = await app.client.$(cssQuoteId('#oscGain[1]'));
        const oscGain2 = await app.client.$(cssQuoteId('#oscGain[2]'));
        const oscGain3 = await app.client.$(cssQuoteId('#oscGain[3]'));
        const oscGain4 = await app.client.$(cssQuoteId('#oscGain[4]'));

        (await oscGain1.getValue()).should.equal('75');
        (await oscGain2.getValue()).should.equal('83');
        (await oscGain3.getValue()).should.equal('100');
        (await oscGain4.getValue()).should.equal('100');

        await app.client.pause(TYPING_PAUSE);
        await oscGain1.setValue('90');
        await oscGain2.click(); // click off the element to force the onchange event
        await app.client.pause(REDRAW_ENV_PAUSE);
        await hooks.waitUntilValueExists(cssQuoteId('#oscGain[1]'), '90');
        (await oscGain1.getValue()).should.equal('90');        
    });

    
    it('env tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceEnvsTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const table = await app.client.$('#envTable')
        await table.waitForDisplayed();
        (await table.isDisplayed()).should.equal(true);
    });

    it('check modified gains', async() => {
        // osc1: should be 60,90 (low: 43,43, up: 65,64)

        const lowGain = await app.client.$(cssQuoteId('#gainAmpLow'));
        const upGain = await app.client.$(cssQuoteId('#gainAmpUp'));
        const lowVal1 = await app.client.$(cssQuoteId('#envAmpLowVal[1]'));
        const lowVal2 = await app.client.$(cssQuoteId('#envAmpLowVal[2]'));
        const upVal1 = await app.client.$(cssQuoteId('#envAmpUpVal[1]'));
        const upVal2 = await app.client.$(cssQuoteId('#envAmpUpVal[2]'));

        (await lowGain.getValue()).should.equal('60');
        (await upGain.getValue()).should.equal('90');
        (await lowVal1.getValue()).should.equal('43');
        (await lowVal2.getValue()).should.equal('43'); 
        (await upVal1.getValue()).should.equal('65');
        (await upVal2.getValue()).should.equal('64'); 
    });

    // set osc gain to zero 
    it('voice tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceVoiceTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const paramTable = await app.client.$('#voiceParamTable')
        const oscTable = await app.client.$('#voiceOscTable') 
        await oscTable.waitForDisplayed()
        await paramTable.waitForDisplayed()
    });

    it('change osc gains', async() => {
        // osc1 should now say 90
        // set it to 0
        const oscGain1 = await app.client.$(cssQuoteId('#oscGain[1]'));
        const oscGain2 = await app.client.$(cssQuoteId('#oscGain[2]'));
        const oscGain3 = await app.client.$(cssQuoteId('#oscGain[3]'));
        const oscGain4 = await app.client.$(cssQuoteId('#oscGain[4]'));

        (await oscGain1.getValue()).should.equal('90');
        (await oscGain2.getValue()).should.equal('83');
        (await oscGain3.getValue()).should.equal('100');
        (await oscGain4.getValue()).should.equal('100');

        await app.client.pause(TYPING_PAUSE);
        await oscGain1.setValue('0');
        await oscGain2.click(); // click off the element to force the onchange event
        await app.client.pause(REDRAW_ENV_PAUSE);
        await hooks.waitUntilValueExists(cssQuoteId('#oscGain[1]'), '0');
        (await oscGain1.getValue()).should.equal('0');        
    });

    
    it('env tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceEnvsTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const table = await app.client.$('#envTable')
        await table.waitForDisplayed();
        (await table.isDisplayed()).should.equal(true);
    });

    it('check modified gains', async() => {
        // osc1: should be 0,0

        const lowGain = await app.client.$(cssQuoteId('#gainAmpLow'));
        const upGain = await app.client.$(cssQuoteId('#gainAmpUp'));
        const lowVal1 = await app.client.$(cssQuoteId('#envAmpLowVal[1]'));
        const lowVal2 = await app.client.$(cssQuoteId('#envAmpLowVal[2]'));
        const upVal1 = await app.client.$(cssQuoteId('#envAmpUpVal[1]'));
        const upVal2 = await app.client.$(cssQuoteId('#envAmpUpVal[2]'));

        (await lowGain.getValue()).should.equal('0');
        (await upGain.getValue()).should.equal('0');
        (await lowVal1.getValue()).should.equal('0');
        (await lowVal2.getValue()).should.equal('0'); 
        (await upVal1.getValue()).should.equal('0');
        (await upVal2.getValue()).should.equal('0'); 
    });

    // now change it back via osc gain
    it('voice tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceVoiceTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const paramTable = await app.client.$('#voiceParamTable')
        const oscTable = await app.client.$('#voiceOscTable') 
        await oscTable.waitForDisplayed()
        await paramTable.waitForDisplayed()
    });

    it('change osc gains', async() => {
        // osc1 should now say 0
        // set it to 90
        const oscGain1 = await app.client.$(cssQuoteId('#oscGain[1]'));
        const oscGain2 = await app.client.$(cssQuoteId('#oscGain[2]'));
        const oscGain3 = await app.client.$(cssQuoteId('#oscGain[3]'));
        const oscGain4 = await app.client.$(cssQuoteId('#oscGain[4]'));

        (await oscGain1.getValue()).should.equal('0');
        (await oscGain2.getValue()).should.equal('83');
        (await oscGain3.getValue()).should.equal('100');
        (await oscGain4.getValue()).should.equal('100');

        await app.client.pause(TYPING_PAUSE);
        await oscGain1.setValue('90');
        await oscGain2.click(); // click off the element to force the onchange event
        await app.client.pause(REDRAW_ENV_PAUSE);
        await hooks.waitUntilValueExists(cssQuoteId('#oscGain[1]'), '90');
        (await oscGain1.getValue()).should.equal('90');        
    });

    
    it('env tab should display', async () => {
        const tab = await app.client.$(`#vceTabs a[href='#vceEnvsTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const table = await app.client.$('#envTable')
        await table.waitForDisplayed();
        (await table.isDisplayed()).should.equal(true);
    });

    it('check modified gains', async() => {
        // osc1: should be 60,90 (low: 43,43, up: 65,64)

        const lowGain = await app.client.$(cssQuoteId('#gainAmpLow'));
        const upGain = await app.client.$(cssQuoteId('#gainAmpUp'));
        const lowVal1 = await app.client.$(cssQuoteId('#envAmpLowVal[1]'));
        const lowVal2 = await app.client.$(cssQuoteId('#envAmpLowVal[2]'));
        const upVal1 = await app.client.$(cssQuoteId('#envAmpUpVal[1]'));
        const upVal2 = await app.client.$(cssQuoteId('#envAmpUpVal[2]'));

        (await lowGain.getValue()).should.equal('60');
        (await upGain.getValue()).should.equal('90');
        (await lowVal1.getValue()).should.equal('43');
        (await lowVal2.getValue()).should.equal('43'); 
        (await upVal1.getValue()).should.equal('65');
        (await upVal2.getValue()).should.equal('64'); 
    });

});
