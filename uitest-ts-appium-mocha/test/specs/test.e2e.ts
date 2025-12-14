import { expect, driver, $ } from '@wdio/globals'

describe('Synergize Testing', () => {
    it('basics', async function () {
        await expect(driver.getTitle().toHaveText('Synergize'));
        await expect($('~synergyName').toHaveText('not connected'));
        await expect($('~controlSurfaceName').toHaveText(''));
    })
})

