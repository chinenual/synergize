const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...

let app;

module.exports = {

  loadVCEViaLeftPanel(name) {
    describe('Load ' + name + '.VCE from left panel', () => {
      it('click load ' + name, async () => {
        await app.client
          .click('.file=' + name)
          .waitUntilTextExists("#name", name, LOAD_VCE_TIMEOUT)
          .getText('#name').should.eventually.equal(name)
      });
    });
  },

  loadVCEViaINTERNALCRT(name) {
    describe('Load ' + name + ' from Internal CRT', () => {
      it('click load ' + name, async () => {
        await app.client
          .click('.file=INTERNAL')
          .waitUntilTextExists("#crt_path", 'INTERNAL', LOAD_VCE_TIMEOUT)
          .getText('#crt_path').should.eventually.equal('INTERNAL')
          .click(`//*[text()='${name}']`)
          .waitUntilTextExists("#name", name, LOAD_VCE_TIMEOUT)
          .getText('#name').should.eventually.equal(name)
      });
    });
  },

  testViewVCE(arrayOfVoices, voiceLoader, context) {


    function cssQuoteId(id) {
      return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
    }

    describe('Test VCE loading - ' + context, () => {
      before(async () => {
        console.log("====== reuse the app");
        app = await hooks.getApp();
      });

      arrayOfVoices.forEach(function (vv, vidx) {
        let v = vv;

        if (voiceLoader != null) {
          voiceLoader(v.name);
        }
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
            // TODO
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
                  .selectByVisibleText(v.filterstab.select.selector,flt.select.text)
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

      });

    });

  }
}