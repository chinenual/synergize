const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;
const LOAD_VCE_TIMEOUT = 20000; // loading in voicing mode can take a while...

module.exports = {

  testViewVCE(arrayOfVoices) {

    let app;

    function cssQuoteId(id) {
      return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
    }

    describe('Test VCE loading', () => {
      before(async () => {
        console.log("====== reuse the app");
        app = await hooks.getApp();
      });

      arrayOfVoices.forEach(function (v, vidx) {

        describe('Load ' + v.name, () => {

          it('load ' + v.name, async () => {
            await app.client
              .click('.file=' + v.name)
              .waitForText('#name', v.name, LOAD_VCE_TIMEOUT)
            await app.client
              .saveScreenshot(`screenshots-${v.name}-voiceTab.png`)
          });

          describe('check voice-tab', () => {
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