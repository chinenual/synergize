const hooks = require('./hooks');
const { DownloadItem } = require('electron');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...

let app;

function cssQuoteId(id) {
    return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test envs page edits', () => {
    before(async () => {
        console.log("====== reuse the app");
        app = await hooks.getApp();
    });

    it('envs tab should display', async () => {
        await app.client
            .click(`#vceTabs a[href='#vceEnvsTab']`)
            .getAttribute(`#vceTabs a[href='#vceEnvsTab']`, 'class').should.eventually.include('active')
            .waitForVisible('#envTable')
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
            .click('#addFreqPoint')
            .click('#addFreqPoint')
            .click('#addFreqPoint')
            .click('#addFreqPoint')

            .waitForVisible(cssQuoteId('#envFreqLoop[5]'))

            .isVisible(cssQuoteId('#envFreqLoop[2]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#envFreqLoop[3]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#envFreqLoop[4]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#envFreqLoop[5]')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)

        await app.client
            .click('#delFreqPoint')
            .waitForVisible(cssQuoteId('#envFreqLoop[5]'), 1000, true) // wait to disappear
            .isVisible(cssQuoteId('#envFreqLoop[5]')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('SUSTAIN at 1', async () => {
        await app.client
            .selectByValue(cssQuoteId('#envFreqLoop[1]'), 'S')
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('S')
            .waitForValue(cssQuoteId('#envFreqLoop[1]'), 'S')
            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('LOOP and REPEAT at 2 should fail', async () => {
        await app.client
            .selectByValue(cssQuoteId('#envFreqLoop[2]'), 'L')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-LoopAfterSustain-alert') })

            .getText('#alertText').should.eventually.include('SUSTAIN point must be after')

            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear

            .selectByValue(cssQuoteId('#envFreqLoop[2]'), 'R')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-LoopAfterSustain-alert') })

            .getText('#alertText').should.eventually.include('SUSTAIN point must be after')

            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear

    });

    it('should be able to move SUSTAIN', async () => {
        await app.client
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('S')
            .selectByValue(cssQuoteId('#envFreqLoop[4]'), 'S')
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envFreqLoop[4]')).should.eventually.equal('S')
            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(false)

            .waitForVisible('#accelFreqUp', 1000, true) // wait to disappear
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('SCREENSHOT', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `TEMP4.failed.`) })

    });

    it('now should be able to add a LOOP or REPEAT', async () => {
        await app.client
            .selectByValue(cssQuoteId('#envFreqLoop[1]'), 'L')
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('L')
            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('SCREENSHOT', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `TEMP3.failed.`) })

    });
    it('now should be able to move the loop', async () => {
        await app.client
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('L')
            .selectByValue(cssQuoteId('#envFreqLoop[2]'), 'R')
            .getValue(cssQuoteId('#envFreqLoop[1]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envFreqLoop[2]')).should.eventually.equal('R')
            .isVisible(cssQuoteId('#accelFreqLow')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#accelFreqUp')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });
    it('SCREENSHOT', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `TEMP2.failed.`) })

    });

    it('should disallow removing row if it contains a loop point', async () => {
        await app.client
            .click('#delFreqPoint')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-deleteLoopPoint-alert') })

            .getText('#alertText').should.eventually.include('Cannot remove envelope point')

            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear
    });
    it('SCREENSHOT', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `TEMP1.failed.`) })

    });

    it('cannot delete sustain point if there are loop points', async () => {
        await app.client
            .selectByValue(cssQuoteId('#envFreqLoop[4]'), '')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-deleteSustainPoint-alert') })

            .getText('#alertText').should.eventually.include('Cannot remove')

            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear

    });

    it('remove loops', async () => {
        await app.client
            .selectByValue(cssQuoteId('#envFreqLoop[2]'), '')
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
            .click('#delFreqPoint')
            .waitForVisible(cssQuoteId('#envFreqLoop[5]'), 1000, true) // wait to disappear
            .isVisible(cssQuoteId('#envFreqLoop[4]')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('screenshot', async () => {
        await app.client
            .then(() => { return hooks.screenshotAndCompare(app, `INITVRAM-after-edit-envs`) })
    });
    it('capture renderer logs', async () => {
        await app.client.getRenderProcessLogs().then(function (logs) {
            logs.forEach(function (log) {
                console.log("RENDERER: " + log.level + ": " + log.source + " : " + log.message);
            });
        });
    });
});
