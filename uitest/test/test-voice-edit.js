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
        console.log("====== reuse the app");
        app = await hooks.getApp();
    });

    it('voice tab should display', async () => {
        await app.client
            .click(`#vceTabs a[href='#vceVoiceTab']`)
            .getAttribute(`#vceTabs a[href='#vceVoiceTab']`, 'class').should.eventually.include('active')
            .waitForVisible('#voiceOscTable')
            .waitForVisible('#voiceParamTable')
    });

    // assumes we are on the INITVRAM voice
    it('sanity check initial state', async () => {
        await app.client
            .getValue('#patchType').should.eventually.equal('0')
            .getValue('#nOsc').should.eventually.equal('1')
            .isVisible(cssQuoteId('#MUTE[1]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#MUTE[2]')).should.eventually.equal(false)
    });

    // test that vibrato type changes when negative depth
    it('down-arrow to VIBDEP 0->-1', async () => {
        await app.client
            .getValue(cssQuoteId('#VIBDEP')).should.eventually.equal('0')
            .getText(cssQuoteId('#vibType')).should.eventually.equal('Sine')
            .click(cssQuoteId('#VIBDEP')).keys('ArrowDown')

            .waitForValue(cssQuoteId('#VIBDEP'), '-1')
            .pause(1000) /// hack: just waiting for text or "should eventually" as below doesnt seem to work :(
            .waitForText(cssQuoteId('#vibType'), 'Random')
            .getText(cssQuoteId('#vibType')).should.eventually.equal('Random')
    });

    // test increasing Osc count
    it('up-arrow to osc count 1->2', async () => {
        await app.client
            .click(cssQuoteId('#nOsc')).keys('ArrowUp')
            .click(cssQuoteId('#nOsc')).keys('ArrowUp')
            .click(cssQuoteId('#nOsc')).keys('ArrowUp')

            .waitForValue(cssQuoteId('#nOsc'), '4')
            .waitForVisible(cssQuoteId('#MUTE[4]'))
    });

    it('keys playable changes', async () => {
        await app.client
            .waitForValue(cssQuoteId('#nOsc'), '4')
            .getText(cssQuoteId('#keysPlayable')).should.eventually.equal('8')
    });

    it('4 rows in the osc table', async () => {
        await app.client
            .isVisible(cssQuoteId('#MUTE[1]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#MUTE[2]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#MUTE[3]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#MUTE[4]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#MUTE[5]')).should.eventually.equal(false)
    });

    it('patch routing matches patchtype 0', async () => {
        await app.client
            .getText(cssQuoteId('#patchFOInputDSR[1]')).should.eventually.equal('')
            .getText(cssQuoteId('#patchAdderInDSR[1]')).should.eventually.equal('1')
            .getText(cssQuoteId('#patchOutputDSR[1]')).should.eventually.equal('1')

            .getText(cssQuoteId('#patchFOInputDSR[2]')).should.eventually.equal('')
            .getText(cssQuoteId('#patchAdderInDSR[2]')).should.eventually.equal('1')
            .getText(cssQuoteId('#patchOutputDSR[2]')).should.eventually.equal('1')

            .getText(cssQuoteId('#patchFOInputDSR[2]')).should.eventually.equal('')
            .getText(cssQuoteId('#patchAdderInDSR[2]')).should.eventually.equal('1')
            .getText(cssQuoteId('#patchOutputDSR[2]')).should.eventually.equal('1')

            .getText(cssQuoteId('#patchFOInputDSR[2]')).should.eventually.equal('')
            .getText(cssQuoteId('#patchAdderInDSR[2]')).should.eventually.equal('1')
            .getText(cssQuoteId('#patchOutputDSR[2]')).should.eventually.equal('1')
    });

    /// change patch type to 1 - routing should change
    it('select patchtype 1', async () => {
        await app.client
            .selectByValue('#patchType', '1')
            .getValue(cssQuoteId('#patchType')).should.eventually.equal('1')
            .waitForValue(cssQuoteId('#patchType'), '1')
    });

    it('patch routing matches patchtype 1', async () => {
        await app.client
            .pause(1000) // hack to allow page to update - nothing to reliably wait for
            .getText(cssQuoteId('#patchFOInputDSR[1]')).should.eventually.equal('')
            .getText(cssQuoteId('#patchAdderInDSR[1]')).should.eventually.equal('')
            .getText(cssQuoteId('#patchOutputDSR[1]')).should.eventually.equal('2')

            .getText(cssQuoteId('#patchFOInputDSR[2]')).should.eventually.equal('2')
            .getText(cssQuoteId('#patchAdderInDSR[2]')).should.eventually.equal('1')
            .getText(cssQuoteId('#patchOutputDSR[2]')).should.eventually.equal('1')

            .getText(cssQuoteId('#patchFOInputDSR[3]')).should.eventually.equal('')
            .getText(cssQuoteId('#patchAdderInDSR[3]')).should.eventually.equal('')
            .getText(cssQuoteId('#patchOutputDSR[3]')).should.eventually.equal('2')

            .getText(cssQuoteId('#patchFOInputDSR[4]')).should.eventually.equal('2')
            .getText(cssQuoteId('#patchAdderInDSR[4]')).should.eventually.equal('1')
            .getText(cssQuoteId('#patchOutputDSR[4]')).should.eventually.equal('1')
    });

    // click MUTE and SOLO - should see class change (.on)
    it('MUTE and SOLO', async () => {
        await app.client
            .isExisting(cssQuoteId('#MUTE[1].on')).should.eventually.equal(false)
            .click(cssQuoteId('#MUTE[1]'))
            .isExisting(cssQuoteId('#MUTE[1].on')).should.eventually.equal(true)

        await app.client
            .isExisting(cssQuoteId('#SOLO[2].on')).should.eventually.equal(false)
            .click(cssQuoteId('#SOLO[2]'))
            .isExisting(cssQuoteId('#SOLO[2].on')).should.eventually.equal(true)

    });

    /// now test that the text-value conversions work for the OHARM and FDETUN spinners
    // OHARM defaults to 1, -1 should display as s1, 0 should display as 0
    // should be able to type those strings (s1 and ran2)
    // should be able to type something in between stepped values for FDETUN and get to the nearest value (247 should end up as 246)
    it('OHARM text conversions', async () => {
        await app.client
            .click(cssQuoteId('#OHARM[1]')).keys('ArrowDown')
            .waitForValue(cssQuoteId('#OHARM[1]'), '0')
            .click(cssQuoteId('#OHARM[1]')).keys('ArrowDown')
            .waitForValue(cssQuoteId('#OHARM[1]'), 's1')

            .clearElement(cssQuoteId('#OHARM[2]'))
            .setValue(cssQuoteId('#OHARM[2]'), 's3')
            .click(cssQuoteId('#OHARM[1]')) // click in a different input to force onchange (FIXME! Enter should be enough)
            .getValue(cssQuoteId('#OHARM[2]')).should.eventually.equal('s3')
            // and validate that the unnderlying value used by the spinner is in sync
            .click(cssQuoteId('#OHARM[2]')).keys('ArrowDown')
            .waitForValue(cssQuoteId('#OHARM[2]'), 's4')
    });


    it('FDETUN text conversions', async () => {

        await app.client
            .click(cssQuoteId('#FDETUN[1]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[1]'), '3')
            .click(cssQuoteId('#FDETUN[1]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[1]'), '6')

            .clearElement(cssQuoteId('#FDETUN[2]'))
            .setValue(cssQuoteId('#FDETUN[2]'), '247')
            .click(cssQuoteId('#FDETUN[1]')) // click in a different input to force onchange (FIXME! Enter should be enough)
            .getValue(cssQuoteId('#FDETUN[2]')).should.eventually.equal('246') // rounded to nearest value

            .click(cssQuoteId('#FDETUN[2]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[2]'), '252')
            .click(cssQuoteId('#FDETUN[2]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[2]'), 'ran1')

            .clearElement(cssQuoteId('#FDETUN[3]'))
            .setValue(cssQuoteId('#FDETUN[3]'), 'ran3')
            .click(cssQuoteId('#FDETUN[1]')) // click in a different input to force onchange (FIXME! Enter should be enough)
            .waitForValue(cssQuoteId('#FDETUN[3]'), 'ran3')
            .getValue(cssQuoteId('#FDETUN[3]')).should.eventually.equal('ran3')

            // and validate that the unnderlying value used by the spinner is in sync
            .click(cssQuoteId('#FDETUN[3]')).keys('ArrowUp')
            .waitForValue(cssQuoteId('#FDETUN[3]'), 'ran4')
    });

    it('screenshot', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `INITVRAM-after-edit-voice`) })
    });
    it('capture renderer logs', async () => {
        await app.client.getRenderProcessLogs().then(function (logs) {
            logs.forEach(function (log) {
                console.log("RENDERER: " + log.level + ": " + log.source + " : " + log.message);
            });
        });
    });

});
