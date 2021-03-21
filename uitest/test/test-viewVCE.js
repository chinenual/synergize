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

                const vname = await app.client.$('#VNAME');
                await (vname.getValue()).should.equal(name);
            });
        });
    },

    loadVCEViaLeftPanelVoicingMode(name) {
        describe('Load ' + name + '.VCE from left panel', () => {
            // need to take extra care since previous test may have been editing G7S -- which could confuse the loader into
            // thinking it's ok to go ahead -- but in fact while the old edited values are still populated.
            // this is not necessary when not in voicing mode since the tests cycle through voices with different names
            it('voice tab should display', async () => {
                const tab = await app.client.$(`#vceTabs a[href='#vceVoiceTab']`)
                await tab.click();
                (await tab.getAttribute('class')).should.include('active')
                const table = await app.client.$('#voiceParamTable')
                await table.waitForDisplayed();
                (await table.isDisplayed('#voiceParamTable')).should.equal(true)
            });
            it('click load ' + name, async () => {
                const vname = await app.client.$('#VNAME')
                // need to clear this since previous test may also be using same voice
                await vname.clearValue()
                const link = await app.client.$('.file=' + name)
                await link.click()
                const confirmText = await app.client.$('#confirmText')

                await confirmText.waitForDisplayed();
                (await confirmText.getText()).should.include('pending edits')

                const confirmOk = await app.client.$('#confirmOkButton')
                await confirmOk.click()
                await confirmText.waitForDisplayed({reverse: true})
                
                app.client.waitUntil(
                    () => vname.getValue() == name,
                    {
                        timeout: LOAD_VCE_TIMEOUT
                    }
                );

                (await vname.getValue()).should.equal(name)
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

                (await crt_path.getText()).should.equal('INTERNAL')

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

                const vname = await app.client.$('#VNAME');
                (await vname.getValue()).should.equal(name)
                
                const vce_crt_name = await app.client.$('#vce_crt_name');
                (await vce_crt_name()).should.equal('INTERNAL');;
                
                (await vce_name.getText()).should.equal(name)
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
                describe('Check fields for ' + v.name, () => {
                    describe('check voice-tab', () => {
                        it('grab screenshot', () => {
                                hooks.screenshotAndCompare(app, `${v.name}-${context}-voiceTab`)
                        });
                        for (k in v.voicetab) {
                            let key = k; // without this let, the value is not consistnent inside the it()'s
                            let value = v.voicetab[key];
                            if (typeof value === 'string') {
                                it(`${key} should be ${value}`, async () => {
                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                    (await ele.getText()).should.equal(value)
                                });
                            } else if (value["exist"] != undefined) {
                                it(`${key} should ${value.exist ? '' : 'not '}exist`, async () => {
                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                    (await ele.isExisting()).should.equal(value.exist);
                                });
                                
                            } else if (value["visible"] != undefined) {
                                it(`${key} should ${value.exist ? '' : 'not '}be visible`, async () => {
                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                    (await ele.isDisplayed()).should.equal(value.visible);
                                });
                                
                            } else if (value["value"] != undefined) {
                                it(`${key} should be ${value.value}`, async () => {
                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                    (await ele.getValue()).should.equal(value.value);
                                });
                                
                            } else if (value["selected"] != undefined) {
                                it(`${key} should be ${value.selected}`, async () => {
                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                    (await ele.isSelected()).should.equal(value.selected);
                                });
                            }
                        }
                    });
                    describe('check envelopes-tab', () => {
                        it('env tab should display', async () => {
                            const tab = await app.client.$(`#vceTabs a[href='#vceEnvsTab']`)
                            await tab.click();
                            (await tab.getAttribute('class')).should.include('active');
                            await hooks.screenshotAndCompare(app, `${v.name}-${context}-envTab`) 

                        });
                        it('default osc should be 1 and env should be All', async () => {
                            const oscEle = await app.client.$('#envOscSelect')
                            const envEle = await app.client.$('#envEnvSelect');
                            (await oscEle.getValue()).should.equal('1');
                            (await envEle.getValue()).should.equal('-1');
                        });
                        describe('check fields for each osc', () => {
                            // for each filter, spot check some fields
                            v.envelopestab.selections.forEach(function (osc, oidx) {
                                describe('check fields for osc ' + osc.select.text, () => {
                                    it('select osc ' + osc.select.text, async () => {
                                        const oscEle = await app.client.$('#envOscSelect')
                                        const envEle = await app.client.$('#envEnvSelect')
                                        
                                        await oscEle.selectByVisibleText(osc.select.text);
                                        (await oscEle.getValue()).should.equal(osc.select.value);
                                
                                        const table = await app.client.$('#envTable')
                                        await table.waitForDisplayed()

                                        const telltale = await app.client.$('#tabTelltaleContent');
                                        (await telltale.getValue()).should.equal(`osc:${osc.select.text}`);
                                        
                                        await hooks.screenshotAndCompare(app, `${v.name}-${context}-filtersTab-{flt.select.text}`)
                                    });
                                    // spot check some elements
                                    for (k in osc) {
                                        let key = k; // without this let, the value is not consistnent inside the it()
                                        if (key != 'select') {

                                            let value = osc[key];
                                            if (typeof value === 'string') {
                                                it(`${key} should be ${value}`, async () => {
                                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                                    (await ele.isDisplayed()).should.equal(true);
                                                    (await ele.isText()).should.equal(value);
                                                });
                                            } else if (value["exist"] != undefined) {
                                                it(`${key} should ${value.exist ? '' : 'not '}exist`, async () => {
                                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                                    (await ele.isExisting()).should.equal(value.exist);
                                                });

                                            } else if (value["visible"] != undefined) {
                                                it(`${key} should ${value.visible ? '' : 'not '}be visible`, async () => {
                                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                                    (await ele.isDisplayed()).should.equal(value.visible);
                                                });

                                            } else if (value["value"] != undefined) {
                                                it(`${key} should be ${value.value}`, async () => {
                                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                                    (await ele.isDisplayed()).should.equal(true);
                                                    (await ele.getValue()).should.equal(value.value);
                                                });

                                            } else if (value["selected"] != undefined) {
                                                it(`${key} should be ${value.selected}`, async () => {
                                                    const ele = await app.client.$('#' + cssQuoteId(key));
                                                    (await ele.isDisplayed()).should.equal(true);
                                                    (await ele.isSelected()).should.equal(value.selected);
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
                            const tab = await app.client.$(`#vceTabs a[href='#vceFiltersTab']`)
                            await tab.click();
                            (await tab.getAttribute('class')).should.include('active');
                            await hooks.screenshotAndCompare(app, `${v.name}-${context}-filtersTab`) 
                        });
                        it('default filter should be ' + v.filterstab.select.text, async () => {
                            const ele = await app.client.$(cssQuoteId(v.filterstab.select.selector));
                            (await ele.getValue()).should.equal(v.filterstab.select.value);
                        });
                        // for each filter, spot check some fields
                        v.filterstab.selections.forEach(function (flt, fidx) {
                            it('check filter ' + fidx, async () => {
                                const ele = await app.client.$(cssQuoteId(v.filterstab.select.selector));
                                await ele.selectByVisibleText(flt.select.text);
                                (await ele.getValue()).should.equal(flt.select.value);
                                
                                const table = await app.client.$('#filterTable')
                                await table.waitForDisplayed()
                                await hooks.screenshotAndCompare(app, `${v.name}-${context}-filtersTab-{flt.select.text}`)
                            });
                            // spot check some elements
                            for (k in flt) {
                                let key = k; // without this let, the value is not consistnent inside the it()'s
                                if (key != 'select') {
                                    let value = flt[key];
                                    it(`${key} should be ${value.value}`, async () => {
                                        const ele = await app.client.$('#' + cssQuoteId(key));
                                        (await ele.getValue()).should.equal(value.value);
                                    });
                                }
                            }
                        });
                    });
                    describe('check keyeq-tab', () => {
                        it('keyeq tab should display', async () => {
                            const tab = await app.client.$(`#vceTabs a[href='#vceKeyEqTab']`)
                            await tab.click();
                            (await tab.getAttribute('class')).should.include('active');
                            const table = await app.client.$('#keyEqTable')
                            await table.waitForDisplayed()
                            await hooks.screenshotAndCompare(app, `${v.name}-${context}-keyeqTab`) 
                        });
                        for (k in v.keyeqtab) {
                            let key = k; // without this let, the value is not consistnent inside the it()'s
                            let value = v.keyeqtab[key];
                            it(`${key} should be ${value.value}`, async () => {
                                const ele = await app.client.$('#' + cssQuoteId(key));
                                (await ele.getValue()).should.equal(value.value);
                            });
                        }
                    });

                    describe('check keyprop-tab', () => {
                        it('keyprop tab should display', async () => {
                            const tab = await app.client.$(`#vceTabs a[href='#vceKeyPropTab']`)
                            await tab.click();
                            (await tab.getAttribute('class')).should.include('active');
                            const table = await app.client.$('#keyPropTable')
                            await table.waitForDisplayed()
                            await hooks.screenshotAndCompare(app, `${v.name}-${context}-keypropTab`) 
                        });
                        for (k in v.keyproptab) {
                            let key = k; // without this let, the value is not consistnent inside the it()'s
                            let value = v.keyproptab[key];
                            it(`${key} should be ${value.value}`, async () => {
                                const ele = await app.client.$('#' + cssQuoteId(key));
                                (await ele.getValue()).should.equal(value.value);
                            });
                        }
                    });


                });
            });

        });

    }
}
