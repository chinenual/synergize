const hooks = require('./hooks');
const WINDOW_PAUSE = 1000;
const voiceGS7 = require('./page-objects/voiceG7S');

let app;

function cssQuoteId(id) {
  return id.replace(/\[/g,'\\[').replace(/\]/g, '\\]');
}

describe('Load G7S', () => {
  before(async () => {
    console.log("====== reuse the app");
    app = await hooks.getApp();
  });

  it('view G7S', async () => {
    await app.client
      .click('.file=G7S')
      .waitUntilTextExists('#name', 'G7S')

      .saveScreenshot('screenshots-G7S-voiceTab.png')
  });

  describe('check voice-tab', () => {
    for (k in voiceGS7) {
      let key = k; // without this let, the value is not consistnent inside the it()'s
      let value = voiceGS7[key];
      if (typeof value === 'string') {
        it(`${key} should be ${value}`, async () => {
          await app.client
            .getText('#' + cssQuoteId(key)).should.eventually.equal(value)
          });
        } else if (value["exist"] != undefined) {
          it(`${key} should ${value.exist ? '' :'not '}exist`, async () => {
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

