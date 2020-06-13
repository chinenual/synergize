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
            it('grab screenshot', async() => {
              await app.client
              .saveScreenshot(`screenshots-${v.name}-${context}-voiceTab.png`)
            });
            for (k in v) {
              let key = k; // without this let, the value is not consistnent inside the it()'s
              let value = v[key];
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
        });

      });

    });

  }
}