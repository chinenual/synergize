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
            (await envLoop2.getValue()).should.equal('');
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
            (await envLoop2.getValue()).should.equal(''); 
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
            await hooks.screenshotAndCompare(app, 'envs-'+which.toLowerCase()+'-deleteLoopPoint-alert');

            (await alertText.getText()).should.include('Cannot remove envelope point');
            await alertButton.click();
            await alertText.waitForDisplayed({inverse: true});
        });

        it(which.toLowerCase()+'-cannot delete sustain point if there are loop points', async () => {
            const envLoop4 = await app.client.$(cssQuoteId('#env'+which+'Loop[4]'));
            const alertButton = await app.client.$(cssQuoteId('#alertModal button'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));

            (await envLoop4.getValue()).should.equal('S');

            await app.client.pause(TYPING_PAUSE);
            await envLoop4.selectByVisibleText('');
            await app.client.pause(TYPING_PAUSE);
           
            await alertText.waitForDisplayed();
            
            await hooks.screenshotAndCompare(app, 'envs-'+which.toLowerCase()+'-deleteSustainPoint-alert');

            (await alertText.getText()).should.include('Cannot remove');
            await app.client.pause(TYPING_PAUSE);
            await alertButton.click();
            await alertText.waitForDisplayed({reverse: true});
        });
        
        it(which.toLowerCase()+'-remove loops', async () => {
            const envLoop1 = await app.client.$(cssQuoteId('#env'+which+'Loop[1]'));
            const envLoop2 = await app.client.$(cssQuoteId('#env'+which+'Loop[2]'));
            const envLoop3 = await app.client.$(cssQuoteId('#env'+which+'Loop[3]'));
            const envLoop4 = await app.client.$(cssQuoteId('#env'+which+'Loop[4]'));
            const accelLow = await app.client.$(cssQuoteId('#accel'+which+'Low'));
            const accelUp = await app.client.$(cssQuoteId('#accel'+which+'Up'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));

            await app.client.pause(TYPING_PAUSE);
            (await envLoop2.getValue()).should.equal('R');
            (await envLoop4.getValue()).should.equal('S');

            await app.client.pause(TYPING_PAUSE);
            await envLoop2.selectByVisibleText('');
            await app.client.pause(TYPING_PAUSE);
            await envLoop4.selectByVisibleText('');
            await app.client.pause(TYPING_PAUSE);
            
            (await envLoop1.getValue()).should.equal('');
            (await envLoop2.getValue()).should.equal('');
            (await envLoop3.getValue()).should.equal('');
            (await envLoop4.getValue()).should.equal('');            
            (await accelLow.isDisplayed()).should.equal(true); 
            (await accelUp.isDisplayed()).should.equal(true);
            (await alertText.isDisplayed()).should.equal(false); 
        });
        it(which.toLowerCase()+'-now can remove a point', async () => {
            const envLoop3 = await app.client.$(cssQuoteId('#env'+which+'Loop[3]'));
            const envLoop4 = await app.client.$(cssQuoteId('#env'+which+'Loop[4]'));
            const alertText = await app.client.$(cssQuoteId('#alertText'));
            const delenvpoint = await app.client.$(cssQuoteId('#del-'+which.toLowerCase()+'-env-point'));

            await app.client.pause(TYPING_PAUSE);
            (await envLoop3.isDisplayed()).should.equal(true);
            (await envLoop4.isDisplayed()).should.equal(true);

            await delenvpoint.click();
            await app.client.pause(TYPING_PAUSE);

            await envLoop4.waitForDisplayed({reverse: true});
            
            (await envLoop3.isDisplayed()).should.equal(true);
            (await envLoop4.isDisplayed()).should.equal(false);

            (await alertText.isDisplayed()).should.equal(false);            
        });

    }
    sharedTests('Freq');
    sharedTests('Amp');

    // test that all the spinner text conversions work at the right ranges
    describe('freq values', () => {
        it('type to envFreqLowVal[1] to and past -127', async () => {
            const ele = await app.client.$(cssQuoteId('#envFreqLowVal[1]'))
            const otherele = await app.client.$(cssQuoteId('#envFreqUpVal[1]'))
            await app.client.pause(TYPING_PAUSE)
            await ele.setValue('-126')
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('-126');

            // should not be able to go below min
            
            await app.client.pause(TYPING_PAUSE)
            await ele.click()
            await app.client.keys('ArrowDown');
            (await ele.getValue()).should.equal('-127');
            
            await app.client.pause(TYPING_PAUSE)
            await ele.click()
            await app.client.keys('ArrowDown');
            (await ele.getValue()).should.equal('-127');
        });
        
        it('type to envFreqUpVal[1] to and past 127', async () => {
            const ele = await app.client.$(cssQuoteId('#envFreqUpVal[1]'))
            const otherele = await app.client.$(cssQuoteId('#envFreqLowVal[1]'))
            await app.client.pause(TYPING_PAUSE)
            await ele.setValue('126')
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('126');

            // should not be able to go above max
            
            await app.client.pause(TYPING_PAUSE)
            await ele.click()
            await app.client.keys('ArrowUp');
            (await ele.getValue()).should.equal('127');
            
            await app.client.pause(TYPING_PAUSE)
            await ele.click()
            await app.client.keys('ArrowUp');
            (await ele.getValue()).should.equal('127');
        });
    });
    describe('freq times', () => {
        describe('type to envFreqLowTime[2] to and past 0', () => {
            it('set envFreqLowTime[2] to 1', async () => {
                const ele = await app.client.$(cssQuoteId('#envFreqLowTime[2]'));
                const otherele = await app.client.$(cssQuoteId('#envAmpUpTime[1]'));
                await app.client.pause(TYPING_PAUSE)
                await ele.setValue('1')
                await otherele.click(); // click in a different input to force onchange
                (await ele.getValue()).should.equal('1');
            });
            // should not be able to go below min
            it('set envFreqLowTime[2] down to 0', async () => {
                const ele = await app.client.$(cssQuoteId('#envFreqLowTime[2]'));
                (await ele.getValue()).should.equal('1'); 
                await ele.click();
                await app.client.keys('ArrowDown');
                await app.client.pause(TYPING_PAUSE);
                (await ele.getValue()).should.equal('0');               
            });
            it('set envFreqLowTime[2] down doesnt go past 0', async () => {
                const ele = await app.client.$(cssQuoteId('#envFreqLowTime[2]'));
                (await ele.getValue()).should.equal('0');               
                await ele.click();
                await app.client.keys('ArrowDown');
                await app.client.pause(TYPING_PAUSE);
                (await ele.getValue()).should.equal('0');               
            });
        });
    });
    describe('amp values', () => {
        it('type to envAmpLowVal[1] to and past 0', async () => {
            const ele = await app.client.$(cssQuoteId('#envAmpLowVal[1]'));
            const otherele = await app.client.$(cssQuoteId('#envAmpUpVal[1]'));
            await app.client.pause(TYPING_PAUSE);
            await ele.setValue('1');
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('1');

            // should not be able to go below min
            await ele.click();
            await app.client.keys('ArrowDown');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('0');         
            
            await ele.click();
            await app.client.keys('ArrowDown');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('0'); 
        });
        it('type to envAmpUpVal[1] to and past 72', async () => {
            const ele = await app.client.$(cssQuoteId('#envAmpUpVal[1]'));
            const otherele = await app.client.$(cssQuoteId('#envAmpLowVal[1]'));
            
            await app.client.pause(TYPING_PAUSE);
            await ele.setValue('71');
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('71');

            await app.client.pause(TYPING_PAUSE)

            await ele.click();
            await app.client.keys('ArrowUp');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('72');               

            await ele.click();
            await app.client.keys('ArrowUp');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('72');
        });
    });

    describe('amp times', () => {
        it('type to envAmpLowTime[1] to and past 0', async () => {
            const ele = await app.client.$(cssQuoteId('#envAmpLowTime[1]'));
            const otherele = await app.client.$(cssQuoteId('#envAmpUpTime[1]'));
            await app.client.pause(TYPING_PAUSE);
            await ele.setValue('1');
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('1');

            // should not be able to go below min
            await ele.click();
            await app.client.keys('ArrowDown');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('0');         
            
            await ele.click();
            await app.client.keys('ArrowDown');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('0');                     
        });
        
        it('type to envAmpUpTime[1] to and past 6576', async () => {
            const ele = await app.client.$(cssQuoteId('#envAmpUpTime[1]'));
            const otherele = await app.client.$(cssQuoteId('#envAmpLowTime[2]'));
            
            await app.client.pause(TYPING_PAUSE);
            await ele.setValue('5858');
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('5859');

            await app.client.pause(TYPING_PAUSE);

            await ele.click();
            await app.client.keys('ArrowUp');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('6576');               

            await ele.click();
            await app.client.keys('ArrowUp');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('6576'); 
        });

    });

    
    describe('accelerations', () => {
        it('type to accelAmpLow to and past 0', async () => {
            const ele = await app.client.$(cssQuoteId('#accelAmpLow'));
            const otherele = await app.client.$(cssQuoteId('#accelAmpUp'));
            await app.client.pause(TYPING_PAUSE);
            await ele.setValue('1')
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('1');

            // should not be able to go below min
            await ele.click();
            await app.client.keys('ArrowDown');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('0');         
            
            await ele.click();
            await app.client.keys('ArrowDown');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('0');
        });
        it('type to accelAmpUp to and past 126', async () => {
            const ele = await app.client.$(cssQuoteId('#accelAmpUp'));
            const otherele = await app.client.$(cssQuoteId('#accelAmpLow'));
            await app.client.pause(TYPING_PAUSE);
            await ele.setValue('126');
            await otherele.click(); // click in a different input to force onchange
            (await ele.getValue()).should.equal('126');

            // should not be able to go below min
            await ele.click();
            await app.client.keys('ArrowUp');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('127');         
            
            await ele.click();
            await app.client.keys('ArrowUp');
            await app.client.pause(TYPING_PAUSE);
            (await ele.getValue()).should.equal('127');
        });

    });

    describe('copy env', () => {
        it('check osc 1 initial conditions', async () => {
            const envOscSelect = await app.client.$(cssQuoteId('#envOscSelect'));
            const envAmpUpTime1 = await app.client.$(cssQuoteId('#envAmpUpTime[1]'));
            const accelAmpUp = await app.client.$(cssQuoteId('#accelAmpUp'));

            (await envOscSelect.getValue()).should.equal('1');
            (await envAmpUpTime1.getValue()).should.equal('6576');
            (await accelAmpUp.getValue()).should.equal('127');
        });

        it('switch to osc 2', async () => {
            const envOscSelect = await app.client.$(cssQuoteId('#envOscSelect'));
            const envAmpUpTime1 = await app.client.$(cssQuoteId('#envAmpUpTime[1]'));
            const accelAmpUp = await app.client.$(cssQuoteId('#accelAmpUp'));

            await app.client.pause(TYPING_PAUSE)
            await envOscSelect.selectByVisibleText('2');
            await app.client.pause(TYPING_PAUSE)
            await hooks.waitUntilValueExists('#tabTelltaleContent', 'osc:2');
            
            (await envOscSelect.getValue()).should.equal('2');
            (await envAmpUpTime1.getValue()).should.equal('0');
            (await accelAmpUp.getValue()).should.equal('30');
        });

        it('copy from 1', async () => {
            const envCopySelect = await app.client.$(cssQuoteId('#envCopySelect'));
            const envAmpUpTime1 = await app.client.$(cssQuoteId('#envAmpUpTime[1]'));
            const accelAmpUp = await app.client.$(cssQuoteId('#accelAmpUp'));

            await app.client.pause(TYPING_PAUSE)
            await envCopySelect.selectByVisibleText('1');
            await app.client.pause(TYPING_PAUSE)
            
            await hooks.waitUntilValueExists('#accelAmpUp', '127');
            
            (await envAmpUpTime1.getValue()).should.equal('6576');
            (await accelAmpUp.getValue()).should.equal('127');
        });

    });
    
    it('screenshot', async () => {
        await hooks.screenshotAndCompare(app, `INITVRAM-after-edit-envs`);
    });
});
