import { expect, $ } from '@wdio/globals'
import mainPage from '../pages/MainPage.ts';

describe('Synergize Testing', () => {
    it('basics', async function () {
        // this seems to be the only way to get the window title via Appium/Mac:
        console.log('one');
        await expect($('//XCUIElementTypeWindow')).toHaveText('Synergize');

        console.log('two');
        let element = await mainPage.synergyName;
        let text = await element.getText();
        console.log("TEXT: " + text)
//        await expect(await mainPage.synergyName).toHaveText('not connected');
        console.log('three');
        await expect(mainPage.controlSurfaceName.toHaveText(''));
        console.log('four');
    })
})

