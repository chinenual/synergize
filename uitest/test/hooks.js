const Application = require('spectron').Application;
const chai = require('chai');
const chaiAsPromised = require('chai-as-promised');
const electron = require('electron');
const {exec} = require("child_process");

global.before(() => {
  chai.should();
  chai.use(chaiAsPromised);
});

module.exports = {
  async startMainApp() {
     exec("../output/darwin-amd64/Synergize.app/Contents/MacOS/Synergize -UITEST 55555", (error, stdout, stderr) => {
    if (error) {
        console.log(`error: ${error.message}`);
        return;
    }
    if (stderr) {
        console.log(`stderr: ${stderr}`);
        return;
    }
    console.log(`stdout: ${stdout}`);

  });
  },

  async startApp() {
    module.exports.startMainApp();
   
    const rendererApp = await new Application({
	
	path: '../output/darwin-amd64/Synergize.app/Contents/MacOS/vendor/electron-darwin-amd64/Synergize.app/Contents/MacOS/Synergize',
	args: ['../output/darwin-amd64/Synergize.app/Contents/MacOS/vendor/astilectron/main.js', '127.0.0.1:55555', 'true' ],

	
        chromeDriverLogPath: './chromedriver.log',
        webdriverLogPath: './webdriver.log'

    }).start();
    chaiAsPromised.transferPromiseness = rendererApp.transferPromiseness;
    return rendererApp;
  },

  async stopApp(app) {
    if (app && app.isRunning()) {
      await app.stop();
    }
  }
};
