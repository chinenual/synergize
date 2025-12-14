import { expect, $ } from '@wdio/globals'

describe('Synergize Testing', () => {
    it('basics', async function () {
        await expect($('~synergyName').toHaveText('not connected'));
        await expect($('~controlSurfaceName').toHaveText(''));
    })
})

