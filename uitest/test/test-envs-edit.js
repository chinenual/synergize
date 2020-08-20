const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...
const TYPING_PAUSE = 500; // slow down typing just a bit to reduce stress on Synergy for non-debounced typing to separate fields

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test envs page edits', () => {
    before(async () => {
        app = await hooks.getApp();
    });

    it('envs tab should display', async () => {
        await app.client
            .click(`#vceTabs a[href='#vceEnvsTab']`)
            .getAttribute(`#vceTabs a[href='#vceEnvsTab']`, 'class').should.eventually.include('active')
            .waitForVisible('#envTable')
            .pause(2000)//HACK
    });

    // assumes we are on the INITVRAM voice
    it('sanity check initial state', async () => {
        await app.client
            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#accelAmpLow')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#accelAmpUp')).should.eventually.equal(true)
    });

    // add some rows so we have enough variations to play with:
    it('adds and removes points', async () => {
        await app.client
            .isVisible(cssQuoteId('#envFreqLoop[2]')).should.eventually.equal(false)
            .pause(TYPING_PAUSE)
            .click('#add-freq-env-point')
            .waitForVisible(cssQuoteId('#envFreqLoop[2]'))
            .pause(TYPING_PAUSE)
            .click('#add-freq-env-point')
            .waitForVisible(cssQuoteId('#envFreqLoop[3]'))
            .pause(TYPING_PAUSE)
            .click('#add-freq-env-point')
            .waitForVisible(cssQuoteId('#envFreqLoop[4]'))
            .pause(TYPING_PAUSE)
            .click('#add-freq-env-point')
            .waitForVisible(cssQuoteId('#envFreqLoop[5]'))

            .pause(TYPING_PAUSE)
            .click('#add-amp-env-point')
            .waitForVisible(cssQuoteId('#envAmpLoop[2]'))

            .isVisible(cssQuoteId('#envFreqLoop[2]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#envFreqLoop[3]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#envFreqLoop[4]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#envFreqLoop[5]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#envAmpLoop[2]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)

        await app.client
            .pause(TYPING_PAUSE)
            .click('#del-freq-env-point')
            .waitForVisible(cssQuoteId('#envFreqLoop[5]'), 1000, true) // wait to disappear
            .isVisible(cssQuoteId('#envFreqLoop[5]')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('SUSTAIN at 1', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envFreqLoop[1]'), 'S')
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('S')
            .waitForValue(cssQuoteId('#envFreqLoop[1]'), 'S')
            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('LOOP and REPEAT at 2 should fail', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envFreqLoop[2]'), 'L')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-LoopAfterSustain-alert') })

            .getText('#alertText').should.eventually.include('SUSTAIN point must be after')

            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear

            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envFreqLoop[2]'), 'R')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-LoopAfterSustain-alert') })

            .getText('#alertText').should.eventually.include('SUSTAIN point must be after')

            .pause(TYPING_PAUSE)
            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear

    });

    it('should be able to move SUSTAIN', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('S')
            .selectByValue(cssQuoteId('#envFreqLoop[4]'), 'S')
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envFreqLoop[4]')).should.eventually.equal('S')
            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(false)

            .waitForVisible('#accelFreqUp', 1000, true) // wait to disappear
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('now should be able to add a LOOP or REPEAT', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envFreqLoop[1]'), 'L')
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('L')
            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('now should be able to move the loop', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('L')
            .selectByValue(cssQuoteId('#envFreqLoop[2]'), 'R')
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envFreqLoop[2]')).should.eventually.equal('R')
            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('should disallow removing row if it contains a loop point', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .click('#del-freq-env-point')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-deleteLoopPoint-alert') })

            .getText('#alertText').should.eventually.include('Cannot remove envelope point')

            .pause(TYPING_PAUSE)
            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear
    });

    it('cannot delete sustain point if there are loop points', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envFreqLoop[4]'), '')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-deleteSustainPoint-alert') })

            .getText('#alertText').should.eventually.include('Cannot remove')

            .pause(TYPING_PAUSE)
            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear

    });

    it('remove loops', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envFreqLoop[2]'), '')
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envFreqLoop[4]'), '')
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envFreqLoop[2]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envFreqLoop[3]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envFreqLoop[4]')).should.eventually.equal('')

            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)

    });

    it('now can remove a point', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .click('#del-freq-env-point')
            .waitForVisible(cssQuoteId('#envFreqLoop[5]'), 1000, true) // wait to disappear
            .isVisible(cssQuoteId('#envFreqLoop[4]')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    // test that all the spinner text conversions work at the right ranges
    describe('freq values', () => {
        it('type to envFreqLowVal[1] to and past -61', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#envFreqLowVal[1]'))
                .setValue(cssQuoteId('#envFreqLowVal[1]'), '-60')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqUpVal[1]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#envFreqLowVal[1]')).should.eventually.equal('-60')
                // should not be able to go below min
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqLowVal[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#envFreqLowVal[1]')).should.eventually.equal('-61')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqLowVal[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#envFreqLowVal[1]')).should.eventually.equal('-61')
        });
        it('type to envFreqUpVal[1] to and past 63', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#envFreqUpVal[1]'))
                .setValue(cssQuoteId('#envFreqUpVal[1]'), '62')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqLowVal[1]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#envFreqUpVal[1]')).should.eventually.equal('62')
                // should not be able to go above max
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqUpVal[1]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envFreqUpVal[1]')).should.eventually.equal('63')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqUpVal[1]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envFreqUpVal[1]')).should.eventually.equal('63')

        });
    });
    describe('freq times', () => {
        describe('type to envFreqLowTime[2] to and past 0', () => {
            it('set envFreqLowTime[2] to 1', async () => {
                await app.client
                    .pause(TYPING_PAUSE)
                    .clearElement(cssQuoteId('#envFreqLowTime[2]'))
                    .setValue(cssQuoteId('#envFreqLowTime[2]'), '1')
                    .pause(TYPING_PAUSE)
                    .click(cssQuoteId('#envAmpUpTime[1]')) // click in a different input to force onchange
                    .getValue(cssQuoteId('#envFreqLowTime[2]')).should.eventually.equal('1')
            });
            // should not be able to go below min
            it('set envFreqLowTime[2] down to 0', async () => {
                await app.client
                    .pause(TYPING_PAUSE)
                    .click(cssQuoteId('#envFreqLowTime[2]')).keys('ArrowDown')
                    .getValue(cssQuoteId('#envFreqLowTime[2]')).should.eventually.equal('0')
            });
            it('set envFreqLowTime[2] down doesnt go past 0', async () => {
                await app.client
                    .pause(TYPING_PAUSE)
                    .click(cssQuoteId('#envFreqLowTime[2]')).keys('ArrowDown')
                    .pause(TYPING_PAUSE)
                    .getValue(cssQuoteId('#envFreqLowTime[2]')).should.eventually.equal('0')
            });
        });
        it('type to envAmpUpTime[2] to and past 6576', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#envAmpUpTime[2]'))
                .setValue(cssQuoteId('#envAmpUpTime[2]'), '5858')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqLowTime[2]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#envAmpUpTime[2]')).should.eventually.equal('5859')
                // should not be able to go above max
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpUpTime[2]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envAmpUpTime[2]')).should.eventually.equal('6576')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpUpTime[2]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envAmpUpTime[2]')).should.eventually.equal('6576')

        });

    });
    describe('amp values', () => {
        it('type to envAmpLowVal[1] to and past 0', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#envAmpLowVal[1]'))
                .setValue(cssQuoteId('#envAmpLowVal[1]'), '1')
                .click(cssQuoteId('#envAmpUpVal[1]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#envAmpLowVal[1]')).should.eventually.equal('1')
                // should not be able to go below min
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpLowVal[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#envAmpLowVal[1]')).should.eventually.equal('0')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpLowVal[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#envAmpLowVal[1]')).should.eventually.equal('0')
        });
        it('type to envAmpUpVal[1] to and past 72', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#envAmpUpVal[1]'))
                .setValue(cssQuoteId('#envAmpUpVal[1]'), '71')
                .click(cssQuoteId('#envAmpLowVal[1]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#envAmpUpVal[1]')).should.eventually.equal('71')
                // should not be able to go above max
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpUpVal[1]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envAmpUpVal[1]')).should.eventually.equal('72')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpUpVal[1]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envAmpUpVal[1]')).should.eventually.equal('72')

        });
    });
    describe('amp times', () => {
        it('type to envAmpLowTime[1] to and past 0', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#envAmpLowTime[1]'))
                .setValue(cssQuoteId('#envAmpLowTime[1]'), '1')
                .click(cssQuoteId('#envAmpUpTime[1]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#envAmpLowTime[1]')).should.eventually.equal('1')
                // should not be able to go below min
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpLowTime[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#envAmpLowTime[1]')).should.eventually.equal('0')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpLowTime[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#envAmpLowTime[1]')).should.eventually.equal('0')
        });
        it('type to envAmpUpTime[1] to and past 6576', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#envAmpUpTime[1]'))
                .setValue(cssQuoteId('#envAmpUpTime[1]'), '5858') // expect the conversion to round to the right value
                .click(cssQuoteId('#envAmpLowTime[1]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#envAmpUpTime[1]')).should.eventually.equal('5859')
                // should not be able to go above max
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpUpTime[1]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envAmpUpTime[1]')).should.eventually.equal('6576')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envAmpUpTime[1]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envAmpUpTime[1]')).should.eventually.equal('6576')

        });
    });

    describe('accelerations', () => {
        it('type to accelAmpLow to and past 0', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#accelAmpLow'))
                .setValue(cssQuoteId('#accelAmpLow'), '1')
                .click(cssQuoteId('#accelAmpUp')) // click in a different input to force onchange
                .getValue(cssQuoteId('#accelAmpLow')).should.eventually.equal('1')
                // should not be able to go below min
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#accelAmpLow')).keys('ArrowDown')
                .getValue(cssQuoteId('#accelAmpLow')).should.eventually.equal('0')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#accelAmpLow')).keys('ArrowDown')
                .getValue(cssQuoteId('#accelAmpLow')).should.eventually.equal('0')
        });
        it('type to accelAmpUp to and past 126', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#accelAmpUp'))
                .setValue(cssQuoteId('#accelAmpUp'), '126') // expect the conversion to round to the right value
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#accelAmpLow')) // click in a different input to force onchange
                .getValue(cssQuoteId('#accelAmpUp')).should.eventually.equal('126')
                // should not be able to go above max
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#accelAmpUp')).keys('ArrowUp')
                .getValue(cssQuoteId('#accelAmpUp')).should.eventually.equal('127')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#accelAmpUp')).keys('ArrowUp')
                .getValue(cssQuoteId('#accelAmpUp')).should.eventually.equal('127')

        });

    });

    describe('copy env', () => {
        it('check osc 1 initial conditions', async () => {

            await app.client
                .getValue('#envOscSelect').should.eventually.equal('1')
                .getValue(cssQuoteId('#envAmpUpTime[1]')).should.eventually.equal('6576')
                .getValue(cssQuoteId('#accelAmpUp')).should.eventually.equal('127')
        });

        it('switch to osc 2', async () => {
            await app.client
                .selectByVisibleText('#envOscSelect', '2')
                .pause(TYPING_PAUSE)
                .getValue('#tabTelltaleContent').should.eventually.equal(`osc:2`)
                .getValue(cssQuoteId('#envAmpUpTime[1]')).should.eventually.equal('0')
                .getValue(cssQuoteId('#accelAmpUp')).should.eventually.equal('30')
        });

        it('copy from 1', async () => {
            await app.client
                .selectByVisibleText('#envCopySelect', '1')
                .pause(TYPING_PAUSE)
                .getValue(cssQuoteId('#envAmpUpTime[1]')).should.eventually.equal('6576')
                .getValue(cssQuoteId('#accelAmpUp')).should.eventually.equal('127')
        });

    });

    it('screenshot', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `INITVRAM-after-edit-envs`) })
    });
});
