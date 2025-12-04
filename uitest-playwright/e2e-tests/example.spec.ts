import { test, expect } from '@playwright/test';
import testPrefs  from './test-prefs';

// Annotate entire file as serial.
// Playwright docs discourage this since its bad practice when testing most web apps. 
// But we're specifically testing a stateful instance of Synergize, so we want to run
// tests in a given order - and not in parallel.
test.describe.configure({ mode: 'serial' });

test.beforeEach(async ({ page }) => {
  await page.goto('');
});

test('has title', async ({ page }) => {
  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/Synergize/);
});

test('initial state unconnected', async ({ page }) => {
  await expect(page.locator('css=#synergyName')).toContainText('not connected');
});

test('Conversion function unit tests', async ({ page }) => {
 let r = await page.evaluate('viewVCE_voice.testConversionFunctions()');
 expect(r, 'result of viewVCE_voice.testConversionFunctions()').toBe(true);

// r = await page.evaluate('viewVCE_envs.testConversionFunctions()');
// expect(r, 'result of viewVCE_envs.testConversionFunctions()').toBe(true);
});

test.describe(testPrefs);

// require('./test-edit-crt');

// describe('Test Voicing Mode views', () => {
//     afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this,app); });
    
//     require('./test-voicingModeOn');

//     describe('initial VRAM image should be loaded', () => {
//         viewVCE.testViewVCE([voiceINITVRAM], null, "voicemode");
//     });
    
//     require('./test-voice-edit');
//     require('./test-envs-edit');
//     require('./test-filter-edit.js');
//     require('./test-keyeq-edit');
//     require('./test-keyprop-edit');
//     require('./test-gain');

//     viewVCE.testViewVCE([voiceG7S, voiceCATHERG, voiceGUITAR2A], viewVCE.loadVCEViaLeftPanelVoicingMode, "voicemode");
    
//     require('./test-voicingModeOff');
    
// });

// describe('Test READ-ONLY views', () => {
//     afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this,app); });
    
//     viewVCE.testViewVCE([voiceG7S, voiceCATHERG, voiceGUITAR2A], viewVCE.loadVCEViaLeftPanel, "readonlyVCE");
//     viewVCE.testViewVCE([voiceG7S, voiceCATHERG, voiceGUITAR2A], viewVCE.loadVCEViaINTERNALCRT, "readonlyCRT");
// });


// // at end since we can't close the window - can just let it get closed implicity by the tear down
// require('./test-about');

// describe('Tear Down', () => {
//     afterEach("screenshot on failure", function () { hooks.screenshotIfFailed(this,app); });
//     after(async () => {
//         console.log("====== tear down the app");
//         await hooks.stopApp(app);
//     });
//     it('last gasp', async () => {
//         (await app.client.getTitle()).should.equal('Synergize')
//     });
// });

