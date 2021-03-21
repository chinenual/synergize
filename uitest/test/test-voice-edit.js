const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...
const TYPING_PAUSE = 500; // slow down typing just a bit to reduce stress on Synergy for non-debounced typing to separate fields
const ADD_OSC_PAUSE = 2000;
const PATCH_PAUSE = 2000;

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test voice page edits', () => {
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

    // assumes we are on the INITVRAM voice
    it('sanity check initial state', async () => {
        const patchType = await app.client.$('#patchType')
        const nOsc = await app.client.$('#nOsc')
        const mute1 = await app.client.$(cssQuoteId('#MUTE[1]'))
        const mute2 = await app.client.$(cssQuoteId('#MUTE[2]'))

        await app.client.pause(TYPING_PAUSE);

        (await patchType.getValue()).should.equal('0');
        (await nOsc.getText()).should.equal('1');;
        (await mute1.isDisplayed()).should.equal(true);
        (await mute2.isDisplayed()).should.equal(false);
    });

    // test that vibrato type changes when negative depth
    it('down-arrow to VIBDEP 0->-1', async () => {
        const vibdep = await app.client.$('#VIBDEP');
        const vibtype = await app.client.$('#vibType');
        (await vibdep.getValue()).should.equal('0');
        (await vibtype.getText()).should.equal('Sine');

        await app.client.pause(TYPING_PAUSE)

        await vibdep.click()
        await app.client.keys('ArrowDown');
        
        await app.client.pause(TYPING_PAUSE)

        (await vibdep.getValue()).should.equal('-1');
        (await vibtype.getText()).should.equal('Random');
    });

    // test increasing Osc count
    it('up-arrow to osc count 1->4', async () => {
        const add = await app.client.$('#add-osc')
        const del = await app.client.$('#del-osc')
        const nOsc = await app.client.$('#nOsc')
        const mute4 = await app.client.$(cssQuoteId('#MUTE[4]'))
        const mute5 = await app.client.$(cssQuoteId('#MUTE[5]'))
 const start=Date.now()
        await add.click()
        console.log('NOSC: ' + await nOsc.getText() + ' - click took ' +(Date.now()-start));
        app.client.waitUntil(
            () => nOsc.getText() == '2',
            {
                timeout: LOAD_VCE_TIMEOUT
            }
        ); 
        console.log('0NOSC: ' + await nOsc.getText()+ ' - waitUntil took ' +(Date.now()-start));

        (await nOsc.getText()).should.equal('2');
        console.log('1NOSC: ' + await nOsc.getText()+ ' - eventually took ' +(Date.now()-start));
        
        await add.click()
        app.client.waitUntil(
            () => nOsc.getText() == '3',
            {
                timeout: LOAD_VCE_TIMEOUT
            }
        );       
        (await nOsc.getText()).should.equal('3');

        await add.click()
        app.client.waitUntil(
            () => nOsc.getText() == '4',
            {
                timeout: LOAD_VCE_TIMEOUT
            }
        );       

        await add.click()
        app.client.waitUntil(
            () => nOsc.getText() == '5',
            {
                timeout: LOAD_VCE_TIMEOUT
            }
        );       

        (await mute4.isDisplayed()).should.equal(true);
        (await mute5.isDisplayed()).should.equal(true);

        await del.click()
        app.client.waitUntil(
            () => nOsc.getText() == '4',
            {
                timeout: LOAD_VCE_TIMEOUT
            }
        );       
        
        (await mute4.isDisplayed()).should.equal(true);
        (await mute5.isDisplayed()).should.equal(false)        ;
 console.log('NOSC: ' + await nOsc.getText())
        await hooks.screenshotAndCompare(app, `DEBUG`)
    });

    it('keys playable changes', async () => {
        const nOsc = await app.client.$('#nOsc')
        const keysPlayable = await app.client.$('#keysPlayable');
        (await nOsc.getText()).should.equal('4');
        (await keysPlayable.getText()).should.equal('8');
    });

    it('4 rows in the osc table', async () => {
        const mute1 = await app.client.$(cssQuoteId('#MUTE[1]'));
        const mute2 = await app.client.$(cssQuoteId('#MUTE[2]'));
        const mute3 = await app.client.$(cssQuoteId('#MUTE[3]'));
        const mute4 = await app.client.$(cssQuoteId('#MUTE[4]'));
        const mute5 = await app.client.$(cssQuoteId('#MUTE[5]'));
        (await mute1.isDisplayed()).should.equal(true);
        (await mute2.isDisplayed()).should.equal(true);
        (await mute3.isDisplayed()).should.equal(true);
        (await mute4.isDisplayed()).should.equal(true);
        (await mute5.isDisplayed()).should.equal(false);
    });

    it('patch routing matches patchtype 0', async () => {
        {
            const ele = await app.client.$(cssQuoteId('#patchFOInputDSR[1]'));
            (await ele.getValue()).should.equal('');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchAdderInDSR[1]'));
            (await ele.getValue()).should.equal('1');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchOutputDSR[1]'));
            (await ele.getValue()).should.equal('1');
        }

        {
            const ele = await app.client.$(cssQuoteId('#patchFOInputDSR[2]'));
            (await ele.getValue()).should.equal('');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchAdderInDSR[2]'));
            (await ele.getValue()).should.equal('1');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchOutputDSR[2]'));
            (await ele.getValue()).should.equal('1');
        }

        {
            const ele = await app.client.$(cssQuoteId('#patchFOInputDSR[2]'));
            (await ele.getValue()).should.equal('');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchAdderInDSR[2]'));
            (await ele.getValue()).should.equal('1');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchOutputDSR[2]'));
            (await ele.getValue()).should.equal('1');
        }

        {
            const ele = await app.client.$(cssQuoteId('#patchFOInputDSR[2]'));
            (await ele.getValue()).should.equal('');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchAdderInDSR[2]'));
            (await ele.getValue()).should.equal('1');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchOutputDSR[2]'));
            (await ele.getValue()).should.equal('1');
        }
    });

    /// change patch type to 1 - routing should change
    it('select patchtype 1', async () => {
        const patchType = await app.client.$('#patchType')
        await patchType.selectByIndex(1);
        (await patchType.getValue()).should.equal('1');
    });

    it('patch routing matches patchtype 1', async () => {

        {
            const ele = await app.client.$(cssQuoteId('#patchFOInputDSR[1]'));
            (await ele.getValue()).should.equal('');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchAdderInDSR[1]'));
            (await ele.getValue()).should.equal('');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchOutputDSR[1]'));
            (await ele.getValue()).should.equal('2');
        }

        {
            const ele = await app.client.$(cssQuoteId('#patchFOInputDSR[2]'));
            (await ele.getValue()).should.equal('2');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchAdderInDSR[2]'));
            (await ele.getValue()).should.equal('1')
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchOutputDSR[2]'));
            (await ele.getValue()).should.equal('1');
        }

        {
            const ele = await app.client.$(cssQuoteId('#patchFOInputDSR[3]'));
            (await ele.getValue()).should.equal('');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchAdderInDSR[3]'))
            (await ele.getValue()).should.equal('')
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchOutputDSR[3]'));
            (await ele.getValue()).should.equal('2');
        }

        {
            const ele = await app.client.$(cssQuoteId('#patchFOInputDSR[4]'));
            (await ele.getValue()).should.equal('2')
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchAdderInDSR[4]'));
            (await ele.getValue()).should.equal('1');
        }
        {
            const ele = await app.client.$(cssQuoteId('#patchOutputDSR[4]'));
            (await ele.getValue()).should.equal('1');
        }
    });

    // click MUTE and SOLO - should see class change (.on)
    it('MUTE and SOLO', async () => {
        const mute1 = await app.client.$(cssQuoteId('#MUTE[1]'));
        const solo2 = await app.client.$(cssQuoteId('#SOLO[2]'));
        const mute1on = await app.client.$(cssQuoteId('#MUTE[1].on'));
        const solo2on = await app.client.$(cssQuoteId('#SOLO[2].on'));

        (await mute1on.isExisting()).should.equal(false);
        await mute1.click()
        await mute1on.isExisting().should.eventually.equal(true)
        
        await solo2on.isExisting().should.eventually.equal(false)
        await solo2.click();
        (await solo2on.isExisting()).should.equal(true);
    });

    /// now test that the text-value conversions work for the OHARM and FDETUN spinners
    // OHARM defaults to 1, -1 should display as s1, 0 should display as 0
    // should be able to type those strings (s1 and ran2)
    // should be able to type something in between stepped values for FDETUN and get to the nearest value (247 should end up as 246)
    it('OHARM text conversions', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#OHARM[1]')).keys('ArrowDown')
            .waitForValue(cssQuoteId('#OHARM[1]'), '0')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#OHARM[1]')).keys('ArrowDown')
            .waitForValue(cssQuoteId('#OHARM[1]'), 's1')

            .clearElement(cssQuoteId('#OHARM[2]'))
            .setValue(cssQuoteId('#OHARM[2]'), 's3')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#OHARM[1]')) // click in a different input to force onchange (FIXME! Enter should be enough)
            .getValue(cssQuoteId('#OHARM[2]')).should.eventually.equal('s3')
        // and validate that the unnderlying value used by the spinner is in sync
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#OHARM[2]')).keys('ArrowDown')
            .waitForValue(cssQuoteId('#OHARM[2]'), 's4')

            .pause(TYPING_PAUSE)
            .clearElement(cssQuoteId('#OHARM[3]'))
            .setValue(cssQuoteId('#OHARM[3]'), 'dc')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#OHARM[1]')) // click in a different input to force onchange (FIXME! Enter should be enough)
            .waitForValue(cssQuoteId('#OHARM[3]'), 'dc')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#OHARM[3]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#OHARM[3]'), 's11')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#OHARM[3]')).keys('ArrowDown')
            .waitForValue(cssQuoteId('#OHARM[3]'), 'dc')

    });


    it('FDETUN text conversions', async () => {

        await app.client
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#FDETUN[1]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[1]'), '3')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#FDETUN[1]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[1]'), '6')

            .pause(TYPING_PAUSE)
            .clearElement(cssQuoteId('#FDETUN[2]'))
            .setValue(cssQuoteId('#FDETUN[2]'), '247')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#FDETUN[1]')) // click in a different input to force onchange (FIXME! Enter should be enough)
            .getValue(cssQuoteId('#FDETUN[2]')).should.eventually.equal('246') // rounded to nearest value

            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#FDETUN[2]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[2]'), '252')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#FDETUN[2]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[2]'), 'ran1')

            .pause(TYPING_PAUSE)
            .clearElement(cssQuoteId('#FDETUN[3]'))
            .setValue(cssQuoteId('#FDETUN[3]'), 'ran3')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#FDETUN[1]')) // click in a different input to force onchange (FIXME! Enter should be enough)
            .waitForValue(cssQuoteId('#FDETUN[3]'), 'ran3')
            .getValue(cssQuoteId('#FDETUN[3]')).should.eventually.equal('ran3')

        // and validate that the unnderlying value used by the spinner is in sync
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#FDETUN[3]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[3]'), 'ran4')
    });

    it('Wave select', async () => {
        const ele = await app.client.$(cssQuoteId('#wkWAVE[1]'))
        await ele.selectByVisibleText('Tri')
        await ele.getValue().should.eventually.equal('Tri')        
    });
    it('Keyprop select', async () => {
        const ele = await app.client.$(cssQuoteId('#wkKEYPROP[1]'))
        await ele.click()
        await ele.isSelected().should.eventually.equal(true)
    });
    it('Filter select', async () => {
        {
            const ele = await app.client.$(cssQuoteId('#FILTER[1]'))
            await ele.selectByVisibleText('Af')
            await ele.getValue().should.eventually.equal('-1') 
        }
        {
            const ele = await app.client.$(cssQuoteId('#FILTER[2]'))
            await ele.selectByVisibleText('Bf')
            await ele.getValue().should.eventually.equal('2') 
        }
    });
    describe('patch edits', () => {

        it('addr', async () => {
            const ele = await app.client.$(cssQuoteId('#patchAdderInDSR[1]'))
            await ele.should.eventually.equal('')
            await ele.click()
            await app.client.keys('ArrowUp')
            await ele.should.eventually.equal('1')
        });
        it('freq', async () => {
            const ele = await app.client.$(cssQuoteId('#patchFOInputDSR[3]'))
            await ele.should.eventually.equal('')
            await ele.click()
            await app.client.keys('ArrowUp')
            await ele.should.eventually.equal('1')
            await ele.click()
            await app.client.keys('ArrowUp')
            await ele.should.eventually.equal('2')
        });
        it('out', async () => {
            const ele = await app.client.$(cssQuoteId('#patchOutputDSR[4]'))
            await ele.should.eventually.equal('1')
            await ele.click()
            await app.client.keys('ArrowUp')
            await ele.should.eventually.equal('2')
        });
    });
    it('screenshot', async () => {
        await hooks.screenshotAndCompare(app, `INITVRAM-after-edit-voice`)
    });

});
