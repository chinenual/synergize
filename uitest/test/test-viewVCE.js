const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...

let app;

function padName(str) {
    while (str.length < 8) {
        str = str + " ";
    }
    return str;
}

module.exports = {

    loadVCEViaLeftPanel(name) {
        describe('Load ' + name + '.VCE from left panel', () => {
            it('click load ' + name, async () => {
                // NOTE: using webdriverio selector extension ("=FOO" means look for text "FOO")
                //    https://webdriver.io/docs/selectors/#element-with-certain-text
                const link = await app.client.$('.file=' + name)
                await link.click()
                
                const vce_name = await app.client.$('#vce_name')
                app.client.waitUntil(
                    () => vce_name.getText() == name,
                    {
                        timeout: LOAD_VCE_TIMEOUT
                    }
                );

//                await app.client.pause(2000) // HACK

                const vname = await app.client.$('#VNAME')
                await vname.getValue().should.eventually.equal(name)
            });
        });
    },

    loadVCEViaLeftPanelVoicingMode(name) {
        describe('Load ' + name + '.VCE from left panel', () => {
            // need to take extra care since previous test may have been editing G7S -- which could confuse the loader into
            // thinking it's ok to go ahead -- but in fact while the old edited values are still populated.
            // this is not necessary when not in voicing mode since the tests cycle through voices with different names
            it('voice tab should display', async () => {
                await app.client
                    .click(`#vceTabs a[href='#vceVoiceTab']`)
                    .getAttribute(`#vceTabs a[href='#vceVoiceTab']`, 'class').should.eventually.include('active')
                    .waitForVisible('#voiceParamTable')
                    .waitUntil(() => {
                        return app.client.$('#voiceParamTable').isVisible()
                    })
                    .isVisible('#voiceParamTable').should.eventually.equal(true)
            });
            it('click load ' + name, async () => {
                await app.client
                // need to clear this since previous test may also be using same voice
                    .clearElement("#VNAME")
                    .click('.file=' + name)

                    .waitForVisible('#confirmText')
                    .getText('#confirmText').should.eventually.include('pending edits')
                    .click('#confirmOKButton')
                    .waitForVisible('#confirmText', 1000, true) // wait to disappear

                    .waitForValue("#VNAME", LOAD_VCE_TIMEOUT)
                    .getValue('#VNAME').should.eventually.equal(name)
            });
        });
    },

    loadVCEViaINTERNALCRT(name) {
        describe('Load ' + name + ' from Internal CRT', () => {
            it('click load ' + name, async () => {
                const link = await app.client.$('.file=INTERNAL')
                await link.click()

                const crt_path = await app.client.$('#crt_path')
                app.client.waitUntil(
                    () => crt_path.getText() == 'INTERNAL',
                    {
                        timeout: LOAD_VCE_TIMEOUT
                    }
                );

                await crt_path.getText().should.eventually.equal('INTERNAL')

                const v_link = await app.client.$(`//*[@id='content']//span[text()='${padName(name)}']`)
                await v_link.click()

                const vce_name = await app.client.$('#vce_name')
                app.client.waitUntil(
                    () => vce_name.getText() == name,
                    {
                        timeout: LOAD_VCE_TIMEOUT
                    }
                );

                await app.client.pause(2000) // HACK

                const vname = await app.client.$('#VNAME')
                await vname.getValue().should.eventually.equal(name)
                
                const vce_crt_name = await app.client.$('#vce_crt_name')
                await vce_crt_name.getText().should.eventually.equal('INTERNAL')
                
                await vce_name.getText().should.eventually.equal(name)
            });
        });
    },

    testViewVCE(arrayOfVoices, voiceLoader, context) {


        function cssQuoteId(id) {
            return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
        }

        describe('Test VCE loading - ' + context, () => {
            before(async () => {
                app = await hooks.getApp();
            });

            arrayOfVoices.forEach(function (vv, vidx) {
                let v = vv;

                if (voiceLoader != null) {
                    voiceLoader(v.name);
                }
/*
                describe('Check fields for ' + v.name, () => {
                    describe('check voice-tab', () => {
                        it('grab screenshot', async () => {
                            await app.client
                                .then(() => { return hooks.screenshotAndCompare(app, `${v.name}-${context}-voiceTab`) })
                        });
                        for (k in v.voicetab) {
                            let key = k; // without this let, the value is not consistnent inside the it()'s
                            let value = v.voicetab[key];
                            if (typeof value === 'string') {
                                it(`${key} should be ${value}`, async () => {
                                    await app.client
                                        .getText('#' + cssQuoteId(key)).should.eventually.equal(value)
                                });
                            } else if (value["exist"] != undefined) {
                                it(`${key} should ${value.exist ? '' : 'not '}exist`, async () => {
                                    await app.client
                                        .isExisting('#' + cssQuoteId(key)).should.eventually.equal(value.exist)
                                });

                            } else if (value["visible"] != undefined) {
                                it(`${key} should ${value.exist ? '' : 'not '}be visible`, async () => {
                                    await app.client
                                        .isVisible('#' + cssQuoteId(key)).should.eventually.equal(value.visible)
                                });

                            } else if (value["value"] != undefined) {
                                it(`${key} should be ${value.value}`, async () => {
                                    await app.client
                                        .getValue('#' + cssQuoteId(key)).should.eventually.equal(value.value)
                                });

                            } else if (value["selected"] != undefined) {
                                it(`${key} should be ${value.selected}`, async () => {
                                    await app.client
                                        .$('#' + cssQuoteId(key)).isSelected().should.eventually.equal(value.selected)
                                });
                            }
                        }
                    });
                    describe('check envelopes-tab', () => {
                        it('env tab should display', async () => {
                            await app.client
                                .click(`#vceTabs a[href='#vceEnvsTab']`)
                                .getAttribute(`#vceTabs a[href='#vceEnvsTab']`, 'class').should.eventually.include('active')
                                .waitForVisible('#vceEnvsTab')
                                .then(() => { return hooks.screenshotAndCompare(app, `${v.name}-${context}-envTab`) })
                        });
                        it('default osc should be 1 and env should be All', async () => {
                            await app.client
                                .getValue('#envOscSelect').should.eventually.equal('1')
                                .getValue('#envEnvSelect').should.eventually.equal('-1')
                        });
                        describe('check fields for each osc', () => {
                            // for each filter, spot check some fields
                            v.envelopestab.selections.forEach(function (osc, oidx) {
                                describe('check fields for osc ' + osc.select.text, () => {
                                    it('select osc ' + osc.select.text, async () => {
                                        await app.client
                                            .selectByVisibleText('#envOscSelect', osc.select.text)
                                            .getValue(cssQuoteId('#envOscSelect')).should.eventually.equal(osc.select.value)
                                            .waitForVisible('#envTable')
                                            .getValue('#tabTelltaleContent').should.eventually.equal(`osc:${osc.select.text}`)
                                            .then(() => { return hooks.screenshotAndCompare(app, `${v.name}-${context}-envTab-${osc.select.text}`) })
                                    });
                                    // spot check some elements
                                    for (k in osc) {
                                        let key = k; // without this let, the value is not consistnent inside the it()
                                        if (key != 'select') {

                                            let value = osc[key];
                                            if (typeof value === 'string') {
                                                it(`${key} should be ${value}`, async () => {
                                                    await app.client
                                                        .isVisible('#' + cssQuoteId(key)).should.eventually.equal(true)
                                                        .getText('#' + cssQuoteId(key)).should.eventually.equal(value)
                                                });
                                            } else if (value["exist"] != undefined) {
                                                it(`${key} should ${value.exist ? '' : 'not '}exist`, async () => {
                                                    await app.client
                                                        .isExisting('#' + cssQuoteId(key)).should.eventually.equal(value.exist)
                                                });

                                            } else if (value["visible"] != undefined) {
                                                it(`${key} should ${value.exist ? '' : 'not '}be visible`, async () => {
                                                    await app.client
                                                        .isVisible('#' + cssQuoteId(key)).should.eventually.equal(value.visible)
                                                });

                                            } else if (value["value"] != undefined) {
                                                it(`${key} should be ${value.value}`, async () => {
                                                    await app.client
                                                        .isVisible('#' + cssQuoteId(key)).should.eventually.equal(true)
                                                        .getValue('#' + cssQuoteId(key)).should.eventually.equal(value.value)
                                                });

                                            } else if (value["selected"] != undefined) {
                                                it(`${key} should be ${value.selected}`, async () => {
                                                    await app.client
                                                        .isVisible('#' + cssQuoteId(key)).should.eventually.equal(true)
                                                        .$('#' + cssQuoteId(key)).isSelected().should.eventually.equal(value.selected)
                                                });
                                            }
                                        }
                                    }
                                });
                            });
                        });

                    });
                    describe('check filters-tab', () => {
                        it('filters tab should display', async () => {
                            await app.client
                                .click(`#vceTabs a[href='#vceFiltersTab']`)
                                .getAttribute(`#vceTabs a[href='#vceFiltersTab']`, 'class').should.eventually.include('active')
                                .waitForVisible('#vceFiltersTab')
                                .then(() => { return hooks.screenshotAndCompare(app, `${v.name}-${context}-filtersTab`) })
                        });
                        it('default filter should be ' + v.filterstab.select.text, async () => {
                            await app.client
                                .getValue(cssQuoteId(v.filterstab.select.selector)).should.eventually.equal(v.filterstab.select.value)
                            //.getText(cssQuoteId(v.filterstab.select.selector)).should.eventually.equal(v.filterstab.select.text)
                        });
                        // for each filter, spot check some fields
                        v.filterstab.selections.forEach(function (flt, fidx) {
                            it('check filter ' + fidx, async () => {
                                await app.client
                                    .selectByVisibleText(v.filterstab.select.selector, flt.select.text)
                                    .getValue(cssQuoteId(v.filterstab.select.selector)).should.eventually.equal(flt.select.value)
                                    .waitForVisible('#filterTable')
                                    .then(() => { return hooks.screenshotAndCompare(app, `${v.name}-${context}-filtersTab-${flt.select.text}`) })
                            });
                            // spot check some elements
                            for (k in flt) {
                                let key = k; // without this let, the value is not consistnent inside the it()'s
                                if (key != 'select') {
                                    let value = flt[key];
                                    it(`${key} should be ${value.value}`, async () => {
                                        await app.client
                                            .getValue('#' + cssQuoteId(key)).should.eventually.equal(value.value)
                                    });
                                }
                            }
                        });
                    });
                    describe('check keyeq-tab', () => {
                        it('keyeq tab should display', async () => {
                            await app.client
                                .click(`#vceTabs a[href='#vceKeyEqTab']`)
                                .getAttribute(`#vceTabs a[href='#vceKeyEqTab']`, 'class').should.eventually.include('active')
                                .waitForVisible('#keyEqTable')
                                .then(() => { return hooks.screenshotAndCompare(app, `${v.name}-${context}-keyeqTab`) })
                        });
                        for (k in v.keyeqtab) {
                            let key = k; // without this let, the value is not consistnent inside the it()'s
                            let value = v.keyeqtab[key];
                            it(`${key} should be ${value.value}`, async () => {
                                await app.client
                                    .getValue('#' + cssQuoteId(key)).should.eventually.equal(value.value)
                            });
                        }
                    });
                    describe('check keyprop-tab', () => {
                        it('keyprop tab should display', async () => {
                            await app.client
                                .click(`#vceTabs a[href='#vceKeyPropTab']`)
                                .getAttribute(`#vceTabs a[href='#vceKeyPropTab']`, 'class').should.eventually.include('active')
                                .waitForVisible('#keyPropTable')
                                .then(() => { return hooks.screenshotAndCompare(app, `${v.name}-${context}-keypropTab`) })
                        });
                        for (k in v.keyproptab) {
                            let key = k; // without this let, the value is not consistnent inside the it()'s
                            let value = v.keyproptab[key];
                            it(`${key} should be ${value.value}`, async () => {
                                await app.client
                                    .getValue('#' + cssQuoteId(key)).should.eventually.equal(value.value)
                            });
                        }
                    });

                });
*/
            });

        });

    }
}
