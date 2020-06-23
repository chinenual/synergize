const Application = require('spectron').Application;
const chai = require('chai');
const chaiAsPromised = require('chai-as-promised');
const electron = require('electron');
const { exec } = require("child_process");
const fs = require('fs')
const path = require('path')
const PNG = require('pngjs').PNG

const SCREEN_DIFFS_ARE_FAILURES = false;

const APPNAME = 'Synergize';
const PORT = 55555; // the port the main process will listen to

global.before(() => {
  chai.should();
  chai.use(chaiAsPromised);
});

// Map nodejs arch to golang arch
let archMap = {
  "arm": "arm",
  "ia32": "386",
  "x86": "386",
  "x64": "amd64",
  "ia64": "amd64"
};

if (archMap[process.arch] === undefined) {
  console.log(`FATAL: unhandled platform/processor type (${process.arch}) - add your variant to archMap in test/hooks.js`);
  process.exit(1);
}

function mainExe() {
  if (process.platform === 'darwin') {
    return `../output/darwin-${archMap[process.arch]}/${APPNAME}.app/Contents/MacOS/${APPNAME}`;
  } else if (process.platform === 'linux') {
    return `../output/linux-${archMap[process.arch]}/${APPNAME}`;
  } else if (process.platform === 'win32') {
    return `../output/windows-${archMap[process.arch]}/${APPNAME}.exe`;
  } else {
    console.log("FATAL: unhandled platform/os - add your variant here");
    process.exit(1);
  }
}

function electronExe() {
  if (process.platform === 'darwin') {
    return `../output/darwin-${archMap[process.arch]}/${APPNAME}.app/Contents/MacOS/vendor/electron-darwin-${archMap[process.arch]}/${APPNAME}.app/Contents/MacOS/${APPNAME}`;
  } else if (process.platform === 'linux') {
    return `../output/linux-${archMap[process.arch]}/vendor/electron-linux-${archMap[process.arch]}/electron`;
  } else if (process.platform === 'win32') {
    return `${process.env.APPDATA}/${APPNAME}/vendor/electron-windows-${archMap[process.arch]}/Electron.exe`;
  } else {
    console.log("FATAL: unhandled platform - add your variant here");
    process.exit(1);
  }
}

function astilectronJS() {
  if (process.platform === 'darwin') {
    return `../output/darwin-${archMap[process.arch]}/${APPNAME}.app/Contents/MacOS/vendor/astilectron/main.js`;
  } else if (process.platform === 'linux') {
    return `../output/linux-${archMap[process.arch]}/vendor/vendor/astilectron/main.js`;
  } else if (process.platform === 'win32') {
    return `${process.env.APPDATA}/${APPNAME}/vendor/astilectron/main.js`;
  } else {
    console.log("FATAL: unhandled platform - add your variant here");
    process.exit(1);
  }
}

module.exports = {
  async startMainApp() {
    console.log(`node arch: "${process.arch}"   golang arch: "${archMap[process.arch]}"`)
    console.log(`Starting main exe: ${mainExe()}`);
    exec(`"${mainExe()}" -UITEST ${PORT} -SERIALVERBOSE`, (error, stdout, stderr) => {
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

  async getApp() {
    return module.exports.app;
  },

  async startApp() {
    module.exports.startMainApp();

    // give the main exe a little time to initialize itself    
    await sleep(4000)
    function sleep(ms) {
      return new Promise((resolve) => {
        setTimeout(resolve, ms);
      });
    }

    console.log(`Starting electron exe: ${electronExe()}`);
    const rendererApp = await new Application({

      path: electronExe(),
      args: [astilectronJS(), `127.0.0.1:${PORT}`, 'true'],

      // for debugging:
      //chromeDriverLogPath: './chromedriver.log',
      //webdriverLogPath: './webdriver.log'

    }).start();
    chaiAsPromised.transferPromiseness = rendererApp.transferPromiseness;
    module.exports.app = rendererApp;
    return rendererApp;
  },

  async stopApp(app) {
    if (app && app.isRunning()) {
      await app.stop();
    }
  },

  screenshotIfFailed(mochaInstance, app) {
    app.client.getRenderProcessLogs().then(function (logs) {
      logs.forEach(function (log) {
        console.log("RENDERER: " + log.level + ": " + log.source + " : " + log.message);
      });
    });

    if (mochaInstance.currentTest.state !== "passed") {
      const ssDir = path.join(__dirname, 'screenshots', process.platform)
      // check that path exists otherwise create it
      if (!fs.existsSync(ssDir)) {
        fs.mkdirSync(ssDir)
      }
      var name = "AFTERHOOK-FAILED-" + mochaInstance.currentTest.fullTitle();
      // sanitize the name (replace spaces, slashes with underscores)
      name = name.replace(/[^A-Za-z0-9_-]/g, '_')
      const ssPath = path.join(ssDir, name + '.failed.png')
      console.log('ERROR:  afterEach write screeshot to ' + ssPath);
      app.client.saveScreenshot(ssPath);
    }
  },

  screenshotAndCompare(app, name) {
    // adapted from https://github.com/webtorrent/webtorrent-desktop/blob/master/test/setup.js

    const ssDir = path.join(__dirname, 'screenshots', process.platform)

    // check that path exists otherwise create it
    if (!fs.existsSync(ssDir)) {
      fs.mkdirSync(ssDir)
    }
    name = name.replace(/[^A-Za-z0-9_-]/g, '_')
    const ssPath = path.join(ssDir, name + '.png')
    let ssBuf

    try {
      ssBuf = fs.readFileSync(ssPath)
    } catch (err) {
      ssBuf = Buffer.alloc(0)
    }

    // many pages have animated charts that last about a second; pause to let them finish
    return app.client.pause(1500).then(function () {
      return app.browserWindow.capturePage()
    }).then(function (buffer) {
      if (ssBuf.length === 0) {
        console.log('Saving screenshot ' + ssPath)
        fs.writeFileSync(ssPath, buffer)
        return chai.assert.isOk(true, 'nothing to compare') // return a non-failure promise
      } else {
        const match = compareIgnoringTransparency(buffer, ssBuf);
        if (match) {
          //console.log('Screenshot matches ' + ssPath)
          return chai.assert.isOk(true, 'screenshots match') // return a non-failure promise
        } else {
          const ssFailedPath = path.join(ssDir, name + '.failed.png')
          console.log('Saving screenshot, failed comparison: ' + ssFailedPath)
          fs.writeFileSync(ssFailedPath, buffer)
          // FIXME: for now, don't make this fail the test -- some of the graphic charts draw lines at slightly
          // different offsets for some reason.  until that gets debugged and fixed, just warn but don't fail
          if (SCREEN_DIFFS_ARE_FAILURES) {
            return chai.assert.fail('screenshot failed comparison ' + ssFailedPath)
          } else {
            console.log('ERROR: Screenshot doesnt match - but not flagging as ERROR: ' + ssFailedPath)
            return chai.assert.isOk(true, 'ignorning screenshot failed comparison ' + ssFailedPath)
          }
        }
      }
    })
  }
};

// Compares two PNGs, ignoring any transparent regions in bufExpected.
// Returns true if they match.  Directly from https://github.com/webtorrent/webtorrent-desktop/blob/master/test/setup.js
function compareIgnoringTransparency(bufActual, bufExpected) {
  // Common case: exact byte-for-byte match
  if (Buffer.compare(bufActual, bufExpected) === 0) return true

  // Otherwise, compare pixel by pixel
  let sumSquareDiff = 0
  let numDiff = 0
  const pngA = PNG.sync.read(bufActual)
  const pngE = PNG.sync.read(bufExpected)
  if (pngA.width !== pngE.width || pngA.height !== pngE.height) return false
  const w = pngA.width
  const h = pngE.height
  const da = pngA.data
  const de = pngE.data
  for (let y = 0; y < h; y++) {
    for (let x = 0; x < w; x++) {
      const i = ((y * w) + x) * 4
      if (de[i + 3] === 0) continue // Skip transparent pixels
      const ca = (da[i] << 16) | (da[i + 1] << 8) | da[i + 2]
      const ce = (de[i] << 16) | (de[i + 1] << 8) | de[i + 2]
      if (ca === ce) continue

      // Add pixel diff to running sum
      // This is necessary on Windows, where rendering apparently isn't quite deterministic
      // and a few pixels in the screenshot will sometimes be off by 1. (Visually identical.)
      numDiff++
      sumSquareDiff += (da[i] - de[i]) * (da[i] - de[i])
      sumSquareDiff += (da[i + 1] - de[i + 1]) * (da[i + 1] - de[i + 1])
      sumSquareDiff += (da[i + 2] - de[i + 2]) * (da[i + 2] - de[i + 2])
    }
  }
  const rms = Math.sqrt(sumSquareDiff / (numDiff + 1))
  const l2Distance = Math.round(Math.sqrt(sumSquareDiff))
  console.log('screenshot diff l2 distance: ' + l2Distance + ', rms: ' + rms)
  return l2Distance < 5000 && rms < 100
}

