const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...
const TYPING_PAUSE = 500; // slow down typing just a bit to reduce stress on Synergy for non-debounced typing to separate fields

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test voice page edits', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    it('voice tab should display', async () => {
        await app.client
            .click(`#vceTabs a[href='#vceVoiceTab']`)
            .getAttribute(`#vceTabs a[href='#vceVoiceTab']`, 'class').should.eventually.include('active')
            .waitForVisible('#voiceOscTable')
            .waitForVisible('#voiceParamTable')
            .pause(2000)//HACK
    });

    // assumes we are on the INITVRAM voice
    it('sanity check initial state', async () => {
        await app.client
            .pause(TYPING_PAUSE)
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
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#VIBDEP')).keys('ArrowDown')

            .waitForValue(cssQuoteId('#VIBDEP'), '-1')
            .pause(1000) /// hack: just waiting for text or "should eventually" as below doesnt seem to work :(
            .waitForText(cssQuoteId('#vibType'), 'Random')
            .getText(cssQuoteId('#vibType')).should.eventually.equal('Random')
    });

    // test increasing Osc count
    it('up-arrow to osc count 1->4', async () => {
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
            .getValue(cssQuoteId('#patchFOInputDSR[1]')).should.eventually.equal('')
            .getValue(cssQuoteId('#patchAdderInDSR[1]')).should.eventually.equal('1')
            .getValue(cssQuoteId('#patchOutputDSR[1]')).should.eventually.equal('1')

            .getValue(cssQuoteId('#patchFOInputDSR[2]')).should.eventually.equal('')
            .getValue(cssQuoteId('#patchAdderInDSR[2]')).should.eventually.equal('1')
            .getValue(cssQuoteId('#patchOutputDSR[2]')).should.eventually.equal('1')

            .getValue(cssQuoteId('#patchFOInputDSR[2]')).should.eventually.equal('')
            .getValue(cssQuoteId('#patchAdderInDSR[2]')).should.eventually.equal('1')
            .getValue(cssQuoteId('#patchOutputDSR[2]')).should.eventually.equal('1')

            .getValue(cssQuoteId('#patchFOInputDSR[2]')).should.eventually.equal('')
            .getValue(cssQuoteId('#patchAdderInDSR[2]')).should.eventually.equal('1')
            .getValue(cssQuoteId('#patchOutputDSR[2]')).should.eventually.equal('1')
    });

    /// change patch type to 1 - routing should change
    it('select patchtype 1', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue('#patchType', '1')
            .getValue(cssQuoteId('#patchType')).should.eventually.equal('1')
            .waitForValue(cssQuoteId('#patchType'), '1')
    });

    it('patch routing matches patchtype 1', async () => {
        await app.client
            .pause(1000) // hack to allow page to update - nothing to reliably wait for
            .getValue(cssQuoteId('#patchFOInputDSR[1]')).should.eventually.equal('')
            .getValue(cssQuoteId('#patchAdderInDSR[1]')).should.eventually.equal('')
            .getValue(cssQuoteId('#patchOutputDSR[1]')).should.eventually.equal('2')

            .getValue(cssQuoteId('#patchFOInputDSR[2]')).should.eventually.equal('2')
            .getValue(cssQuoteId('#patchAdderInDSR[2]')).should.eventually.equal('1')
            .getValue(cssQuoteId('#patchOutputDSR[2]')).should.eventually.equal('1')

            .getValue(cssQuoteId('#patchFOInputDSR[3]')).should.eventually.equal('')
            .getValue(cssQuoteId('#patchAdderInDSR[3]')).should.eventually.equal('')
            .getValue(cssQuoteId('#patchOutputDSR[3]')).should.eventually.equal('2')

            .getValue(cssQuoteId('#patchFOInputDSR[4]')).should.eventually.equal('2')
            .getValue(cssQuoteId('#patchAdderInDSR[4]')).should.eventually.equal('1')
            .getValue(cssQuoteId('#patchOutputDSR[4]')).should.eventually.equal('1')
    });

    // click MUTE and SOLO - should see class change (.on)
    it('MUTE and SOLO', async () => {
        await app.client
            .isExisting(cssQuoteId('#MUTE[1].on')).should.eventually.equal(false)
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#MUTE[1]'))
            .isExisting(cssQuoteId('#MUTE[1].on')).should.eventually.equal(true)

        await app.client
            .isExisting(cssQuoteId('#SOLO[2].on')).should.eventually.equal(false)
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#SOLO[2]'))
            .isExisting(cssQuoteId('#SOLO[2].on')).should.eventually.equal(true)

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
            .click(cssQuoteId('#OHARM[3]')).keys('ArrowDown')
            .waitForValue(cssQuoteId('#OHARM[3]'), '31')
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#OHARM[3]')).keys('ArrowUp')
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
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#wkWAVE[1]'), 'Tri')
            .getValue(cssQuoteId('#wkWAVE[1]')).should.eventually.equal('Tri')
            .waitForValue(cssQuoteId('#wkWAVE[1]'), 'Tri')
    });
    it('Keyprop select', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .click(cssQuoteId('#wkKEYPROP[1]'))
            .getValue(cssQuoteId('#wkKEYPROP[1]')).should.eventually.equal('true')
            .waitForValue(cssQuoteId('#wkKEYPROP[1]'), 'true')
    });
    it('Filter select', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByVisibleText(cssQuoteId('#FILTER[1]'), 'Af')
            .getValue(cssQuoteId('#FILTER[1]')).should.eventually.equal('-1')
            .waitForValue(cssQuoteId('#FILTER[1]'), '-1')

            .pause(TYPING_PAUSE)
            .selectByVisibleText(cssQuoteId('#FILTER[2]'), 'Bf')
            .getValue(cssQuoteId('#FILTER[2]')).should.eventually.equal('2')
            .waitForValue(cssQuoteId('#FILTER[2]'), '2')
    });

    it('screenshot', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `INITVRAM-after-edit-voice`) })
    });

});
