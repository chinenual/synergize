const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;
const voiceGS7 = require('./page-objects/voice-G7S');
const voiceCATHERG = require('./page-objects/voice-CATHERG');
const voiceGUITAR2A = require('./page-objects/voice-GUITAR2A');

let app;

function cssQuoteId(id) {
  return id.replace(/\[/g, '\\[').replace(/\]/g, '\\]');
}

describe('Test read-only VCE loading', () => {
  before(async () => {
    console.log("====== reuse the app");
    app = await hooks.getApp();
  });

  [voiceGS7, voiceCATHERG, voiceGUITAR2A].forEach(function (v, vidx) {

    describe('Load ' + v.name, () => {

      it('view ' + v.name, async () => {
        await app.client
          .click('.file=' + v.name)
          .waitUntilTextExists('#name', v.name)

          .saveScreenshot('screenshots-+v.name+-voiceTab.png')
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

