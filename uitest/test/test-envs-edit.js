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
        const tab = await app.client.$(`#vceTabs a[href='#vceEnvsTab']`)
        await tab.click();
        (await tab.getAttribute('class')).should.include('active');
        const ele = await app.client.$('#envTable')
        await ele.waitForDisplayed()
    });

    // assumes we are on the INITVRAM voice
    it('sanity check initial state', async () => {
        const accelFreqLow = await app.client.$(cssQuoteId('#accelFreqLow'));
        const accelFreqUp = await app.client.$(cssQuoteId('#accelFreqUp'));
        const accelAmpLow = await app.client.$(cssQuoteId('#accelAmpLow'));
        const accelAmpUp = await app.client.$(cssQuoteId('#accelAmpUp'));

        (await accelFreqLow.isDisplayed()).should.equal(true);
        (await accelFreqUp.isDisplayed()).should.equal(true);
        (await accelAmpLow.isDisplayed()).should.equal(true);
        (await accelAmpUp.isDisplayed()).should.equal(true);      
    });

    function sharedTests(which) {
        
        // add some rows so we have enough variations to play with:
        it(which.toLowerCase() + '-adds and removes points', async () => {
            const envLoop2 = await app.client.$(cssQuoteId('#env'+which+'Loop[2]'));
            const envLoop3 = await app.client.$(cssQuoteId('#env'+which+'Loop[3]'));
            const envLoop4 = await app.client.$(cssQuoteId('#env'+which+'Loop[4]'));
            const envLoop5 = await app.client.$(cssQuoteId('#env'+which+'Loop[5]'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));
            const addenvpoint = await app.client.$(cssQuoteId('#add-'+which.toLowerCase()+'-env-point'));
            const delenvpoint = await app.client.$(cssQuoteId('#del-'+which.toLowerCase()+'-env-point'));

            (await envLoop2.isDisplayed()).should.equal(false);
            
            await addenvpoint.click();
            await app.client.pause(2000);//HACK
            await envLoop2.waitForDisplayed();
            
            await addenvpoint.click();
            await app.client.pause(2000);//HACK
            await envLoop3.waitForDisplayed();
            
            await addenvpoint.click();
            await envLoop4.waitForDisplayed();
            
            await addenvpoint.click();
            await envLoop5.waitForDisplayed();
            
            (await envLoop2.isDisplayed()).should.equal(true);
            (await envLoop3.isDisplayed()).should.equal(true);
            (await envLoop4.isDisplayed()).should.equal(true);
            (await envLoop5.isDisplayed()).should.equal(true);
            (await alertText.isDisplayed()).should.equal(false);
            
            await delenvpoint.click();
            await envLoop5.waitForDisplayed({reverse: true});
            
            (await envLoop5.isDisplayed()).should.equal(false);
            (await alertText.isDisplayed()).should.equal(false);
        });

        it(which.toLowerCase()+'-SUSTAIN at 1', async () => {
            const envLoop1 = await app.client.$(cssQuoteId('#env'+which+'Loop[1]'));
            const accelLow = await app.client.$(cssQuoteId('#accel'+which+'Low'));
            const accelUp = await app.client.$(cssQuoteId('#accel'+which+'Up'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));
            
            await envLoop1.selectByVisibleText('S->');
            await app.client.pause(TYPING_PAUSE);
            (await envLoop1.getValue()).should.equal('S');
            (await accelLow.isDisplayed()).should.equal(false);
            (await accelUp.isDisplayed()).should.equal(false);
            (await alertText.isDisplayed()).should.equal(false);
        });

        
        it(which.toLowerCase()+'-LOOP at 2 should fail', async () => {
            const envLoop2 = await app.client.$(cssQuoteId('#env'+which+'Loop[2]'));
            const alertButton = await app.client.$(cssQuoteId('#alertModal button'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));
            
            await app.client.pause(TYPING_PAUSE);
            await envLoop2.selectByVisibleText('L->');
            await app.client.pause(TYPING_PAUSE);

            await alertText.waitForDisplayed();

            await hooks.screenshotAndCompare(app, 'envs-'+which.toLowerCase()+'-LLoopAfterSustain-alert');

            (await alertText.getText()).should.include('SUSTAIN point must be after')
            
            await alertButton.click();
            
            await alertText.waitForDisplayed({inverse: true});
            
        });

        it(which.toLowerCase()+'-REPEAT at 2 should fail', async () => {
            const envLoop2 = await app.client.$(cssQuoteId('#env'+which+'Loop[2]'));
            const alertButton = await app.client.$(cssQuoteId('#alertModal button'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));
            
            await app.client.pause(TYPING_PAUSE);
            await envLoop2.selectByVisibleText('R->');
            await app.client.pause(TYPING_PAUSE);
            
            await alertText.waitForDisplayed();

            await hooks.screenshotAndCompare(app, 'envs-'+which.toLowerCase()+'-RLoopAfterSustain-alert');

            (await alertText.getText()).should.include('SUSTAIN point must be after')
            
            await alertButton.click();
            
            await alertText.waitForDisplayed({inverse: true});
            
        });

        it(which.toLowerCase()+'-should be able to move SUSTAIN', async () => {
            const envLoop1 = await app.client.$(cssQuoteId('#env'+which+'Loop[1]'));
            const envLoop4 = await app.client.$(cssQuoteId('#env'+which+'Loop[4]'));
            const accelLow = await app.client.$(cssQuoteId('#accel'+which+'Low'));
            const accelUp = await app.client.$(cssQuoteId('#accel'+which+'Up'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));

            await envLoop4.selectByVisibleText('S->');
            await app.client.pause(TYPING_PAUSE);

            await hooks.waitUntilValueExists(cssQuoteId('#env'+which+'Loop[4]'), 'S');
            (await envLoop1.getValue()).should.equal('');
            (await envLoop4.getValue()).should.equal('S');
            (await accelLow.isDisplayed()).should.equal(false); 
            (await accelUp.isDisplayed()).should.equal(false);
            (await alertText.isDisplayed()).should.equal(false);
        });
        it(which.toLowerCase()+'-now should be able to add a LOOP or REPEAT', async () => {
            const envLoop1 = await app.client.$(cssQuoteId('#env'+which+'Loop[1]'));
            const accelLow = await app.client.$(cssQuoteId('#accel'+which+'Low'));
            const accelUp = await app.client.$(cssQuoteId('#accel'+which+'Up'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));

            await envLoop1.selectByVisibleText('L->');
            await app.client.pause(TYPING_PAUSE);

            await hooks.waitUntilValueExists(cssQuoteId('#env'+which+'Loop[1]'), 'L');
            (await envLoop1.getValue()).should.equal('L');
            (await accelLow.isDisplayed()).should.equal(false); 
            (await accelUp.isDisplayed()).should.equal(false);
            (await alertText.isDisplayed()).should.equal(false);
        });
        it(which.toLowerCase()+'-now should be able to move the loop', async () => {
            const envLoop1 = await app.client.$(cssQuoteId('#env'+which+'Loop[1]'));
            const envLoop2 = await app.client.$(cssQuoteId('#env'+which+'Loop[2]'));
            const accelLow = await app.client.$(cssQuoteId('#accel'+which+'Low'));
            const accelUp = await app.client.$(cssQuoteId('#accel'+which+'Up'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));

            (await envLoop1.getValue()).should.equal('L');
            await envLoop2.selectByVisibleText('R->');
            await app.client.pause(TYPING_PAUSE);

            await hooks.waitUntilValueExists(cssQuoteId('#env'+which+'Loop[2]'), 'R');
            (await envLoop1.getValue()).should.equal('');
            (await envLoop2.getValue()).should.equal('R');
            (await accelLow.isDisplayed()).should.equal(false); 
            (await accelUp.isDisplayed()).should.equal(false);
            (await alertText.isDisplayed()).should.equal(false);            
	    });

        it(which.toLowerCase()+'-should disallow removing row if it contains a loop point', async () => {
            const alertButton = await app.client.$(cssQuoteId('#alertModal button'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));
            const delenvpoint = await app.client.$(cssQuoteId('#del-'+which.toLowerCase()+'-env-point'));

            await delenvpoint.click();
            await alertText.waitForDisplayed();
            await hooks.screenshotAndCompare(app, 'envs-freq-deleteLoopPoint-alert');

            (await alertText.getText()).should.include('Cannot remove envelope point');
            await alertButton.click();
            await alertText.waitForDisplayed({inverse: true});
        });
    }
    sharedTests('Freq');
    sharedTests('Amp');
   

/*********

    it('freq-cannot delete sustain point if there are loop points', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envFreqLoop[4]'), '')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-freq-deleteSustainPoint-alert') })

            .getText('#alertText').should.eventually.include('Cannot remove')

            .pause(TYPING_PAUSE)
            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear

    });

    it('amp-cannot delete sustain point if there are loop points', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envAmpLoop[4]'), '')

            .waitForVisible('#alertText')
            .then(() => { return hooks.screenshotAndCompare(app, 'envs-amp-deleteSustainPoint-alert') })

            .getText('#alertText').should.eventually.include('Cannot remove')

            .pause(TYPING_PAUSE)
            .click('#alertModal button')
            .waitForVisible('#alertText', 1000, true) // wait to disappear

    });

    it('freq-remove loops', async () => {
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

    it('amp-remove loops', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envAmpLoop[2]'), '')
            .pause(TYPING_PAUSE)
            .selectByValue(cssQuoteId('#envAmpLoop[4]'), '')
            .getValue(cssQuoteId('#envAmpLoop[1]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envAmpLoop[2]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envAmpLoop[3]')).should.eventually.equal('')
            .getValue(cssQuoteId('#envAmpLoop[4]')).should.eventually.equal('')

            .isVisible(cssQuoteId('#accelAmpLow')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#accelAmpUp')).should.eventually.equal(true)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)

    });

    it('freq-now can remove a point', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .click('#del-freq-env-point')
            .waitForVisible(cssQuoteId('#envFreqLoop[5]'), 1000, true) // wait to disappear
            .isVisible(cssQuoteId('#envFreqLoop[4]')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    it('amp-now can remove a point', async () => {
        await app.client
            .pause(TYPING_PAUSE)
            .click('#del-amp-env-point')
            .waitForVisible(cssQuoteId('#envAmpLoop[5]'), 1000, true) // wait to disappear
            .isVisible(cssQuoteId('#envAmpLoop[4]')).should.eventually.equal(false)
            .isVisible(cssQuoteId('#alertText')).should.eventually.equal(false)
    });

    // test that all the spinner text conversions work at the right ranges
    describe('freq values', () => {
        it('type to envFreqLowVal[1] to and past -127', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#envFreqLowVal[1]'))
                .setValue(cssQuoteId('#envFreqLowVal[1]'), '-126')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqUpVal[1]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#envFreqLowVal[1]')).should.eventually.equal('-126')
            // should not be able to go below min
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqLowVal[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#envFreqLowVal[1]')).should.eventually.equal('-127')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqLowVal[1]')).keys('ArrowDown')
                .getValue(cssQuoteId('#envFreqLowVal[1]')).should.eventually.equal('-127')
        });
        it('type to envFreqUpVal[1] to and past 127', async () => {
            await app.client
                .pause(TYPING_PAUSE)
                .clearElement(cssQuoteId('#envFreqUpVal[1]'))
                .setValue(cssQuoteId('#envFreqUpVal[1]'), '126')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqLowVal[1]')) // click in a different input to force onchange
                .getValue(cssQuoteId('#envFreqUpVal[1]')).should.eventually.equal('126')
            // should not be able to go above max
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqUpVal[1]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envFreqUpVal[1]')).should.eventually.equal('127')
                .pause(TYPING_PAUSE)
                .click(cssQuoteId('#envFreqUpVal[1]')).keys('ArrowUp')
                .getValue(cssQuoteId('#envFreqUpVal[1]')).should.eventually.equal('127')

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
*********/
    
    it('screenshot', async () => {
        await hooks.screenshotAndCompare(app, `INITVRAM-after-edit-envs`);
    });
});
