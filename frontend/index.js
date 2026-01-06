// const { dialog } = require('electron').remote;
//
// let shell = require('electron').shell
import './scss/styles.scss'

import {UIService} from '/bindings/github.com/chinenual/synergize';
import * as wails from '@wailsio/runtime';
import {_} from 'lodash';

import {dx2syn} from './dx2syn';
// import Modal from 'modal-vanilla';
import {closeModal, openModal} from './modal';
import {syn2midi} from './syn2midi';
import {viewCRT} from './viewCRT';
import {viewVCE} from './viewVCE';
import {viewVCE_envs} from './viewVCE_envs';
import {viewVCE_voice} from './viewVCE_voice';

// hook the console log to log both to the browser console, but also to the
// backend
const orig_consoleLog = console.log;
const orig_consoleInfo = console.info;
const orig_consoleWarn = console.warn;
const orig_consoleError = console.error;
function redirectConsole() {
  console.log = function(...args) {
    orig_consoleLog(...args);
    UIService.LogInfo('LOG', args);
  };
  console.info = function(...args) {
    orig_consoleInfo(...args);
    UIService.LogInfo('INFO', args);
  };
  console.warn = function(...args) {
    orig_consoleWarn(...args);
    UIService.LogWarn('WARN', args);
  };
  console.error = function(...args) {
    orig_consoleError(...args);
    UIService.LogError('ERROR', args);
  };
}
redirectConsole();

export let index = {
  DEBOUNCE_WAIT_SHORT: 50,
  DEBOUNCE_WAIT: 250,

  init: function () {
    console.log('TOP OF INIT');
    dx2syn.init();
    syn2midi.init();

    // init menus to default state
    index.updateConnectionStatus('', '')

    // Explore default path
    index.explore();

    index.runUnitTests();
  },

  runUnitTests: function() {
    let result = viewVCE_voice.testConversionFunctions();
    if (result != null) {
      index.errorNotification(
          'viewVCE_voice.testConversionFunctions failed ' + result);
    }
    result = viewVCE_envs.testConversionFunctions();
    if (result != null) {
      index.errorNotification(
          'viewVCE_envs.testConversionFunctions failed ' + result);
    }
  },

  browserOpenURL: function(url) {
    wails.Browser.OpenURL(url);
  },

  checkInputElementValue: function(ele) {
    if (!ele.value.match(/^-?\d+$/)) {
      return undefined;
    }
    let result = parseInt(ele.value, 10);
    if (ele.hasAttribute('min')) {
      let min = parseInt(ele.getAttribute('min'), 10);
      if (result < min) result = min;
    }
    if (ele.hasAttribute('max')) {
      let max = parseInt(ele.getAttribute('max'), 10);
      if (result > max) result = max;
    }
    return result;
  },

  spinnerOn: function() {
    document.getElementById('spinner').style.display = 'block';
  },
  spinnerOff: function() {
    document.getElementById('spinner').style.display = 'none';
  },

  confirmDialog: function(message, successCallback) {
    if (typeof message != 'string') {
      message = JSON.stringify(message);
    }
    console.log('CONFIRM DIALOG: ' + message)
    document.getElementById('confirmTitle').innerHTML = 'Confirm';
    document.getElementById('confirmText').innerHTML = message;
    document.getElementById('confirmOKButton').onclick = successCallback;
    //    let modal = new Modal({
    //        el: document.getElementById('confirmModal'),
    //       backdrop: 'static'
    //      });
    //    console.log('modal', modal);
    //    modal.show();
    openModal('confirmModal');
  },

  errorNotification: function(message) {
    // if (typeof message != 'string') {
    //   message = JSON.stringify(message);
    // }
    console.log('ERROR NOTIFICATION: ' + message)
    document.getElementById('alertTitle').innerHTML = 'Error';
    document.getElementById('alertText').innerHTML = message;
    //    let modal = new Modal({
    //      el: document.getElementById('alertModal'),
    //      backdrop: 'static'
    //    });
    //    console.log('modal', modal);
    //    modal.show();
    openModal('alertModal');
  },
  infoNotification: function(message) {
    if (typeof message != 'string') {
      message = JSON.stringify(message);
    }
    console.log('INFO NOTIFICATION: ' + message)
    document.getElementById('alertTitle').innerHTML = 'Info';
    document.getElementById('alertText').innerHTML = message;
    // Make alert messages hide themselves after 3s - no need to click
    //    let modal = new Modal({
    //      el: document.getElementById('alertModal'),
    //      backdrop: 'static'
    //    });
    setTimeout(function() {
      //      modal.hide();
      closeModal(document.getElementById('alertModal'));
    }, 3000);
    //   modal.show();
    openModal('alertModal');
  },

  chooseZeroconfService: function(
      prompt1, choices1, prompt2, choices2, onCancel, onOK, onRescan) {
    console.log(
        'chooseZeroconfService: ' + prompt1 + ' ' + JSON.stringify(choices1) +
        ' ' + prompt2 + ' ' + JSON.stringify(choices2));
    if (prompt1 != null) {
      document.getElementById('chooseZeroconf1Prompt').innerHTML = prompt1;
      document.querySelector('#zeroconf1Div').style.display = 'block';
    } else {
      document.getElementById('chooseZeroconf1Prompt').innerHTML = '';
      document.getElementById('chooseZeroconf1Items').innerHTML = '';
      document.querySelector('#zeroconf1Div').style.display = 'none';
    }
    if (prompt2 != null) {
      document.getElementById('chooseZeroconf2Prompt').innerHTML = prompt2;
      document.querySelector('#zeroconf2Div').style.display = 'block';
    } else {
      document.getElementById('chooseZeroconf2Prompt').innerHTML = '';
      document.getElementById('chooseZeroconf2Items').innerHTML = '';
      document.querySelector('#zeroconf2Div').style.display = 'none';
    }
    if (prompt1 != null && choices1 != null) {
      let html = '';
      for (let i = 0; i < choices1.length; i++) {
        let addr = ''
        if (choices1[i].Port != 0) {
          addr = ` (${choices1[i].HostName}:${choices1[i].Port})`
        }
        html = html +
            `
		    <div class="form-check">
                <input class="form-check-input" type="radio" name="chooseZeroconf1Radios"
		aria-label="chooseZeroconf1Radio${i}" id="chooseZeroconf1Radio${
                   i}" 
                value="${i}" ${i == 0 ? 'checked' : ''}>
                <label class="form-check-label" for="chooseZeroconf1Radio${i}">
                   ${choices1[i].InstanceName}${addr}
                </label>
			</div>`
        console.log('html now ' + html);
      }
      document.getElementById('chooseZeroconf1Items').innerHTML = html;
      console.log(
          'innerHTML now ' +
          document.getElementById('chooseZeroconf1Items').innerHTML);
    }

    if (prompt2 != null && choices2 != null) {
      let html = '';
      for (let i = 0; i < choices2.length; i++) {
        let addr = ''
        if (choices2[i].Port != 0) {
          addr = ` (${choices2[i].HostName}:${choices2[i].Port})`
        }
        html = html +
            `
		    <div class="form-check">
                <input class="form-check-input" type="radio" name="chooseZeroconf2Radios" 
		 aria-label="chooseZeroconf2Radio${
                   i}" id="chooseZeroconf2Radio${i}" 
                 value="${i}" ${i == 0 ? 'checked' : ''}>
                <label class="form-check-label" for="chooseZeroconf2Radio${i}">
                   ${choices2[i].InstanceName}${addr}
                </label>
			</div>`
        console.log('html now ' + html);
      }
      document.getElementById('chooseZeroconf2Items').innerHTML = html;
      console.log(
          'innerHTML now ' +
          document.getElementById('chooseZeroconf2Items').innerHTML);
    }

    document.getElementById('chooseZeroconfCancelButton').onclick = function() {
      console.log('Cancelled');
      onCancel();
    };
    document.getElementById('chooseZeroconfOKButton').onclick = function() {
      let idx1 = null
      let selected1 = null
      let idx2 = null
      let selected2 = null
      if (choices1 != null && prompt1 != null) {
        idx1 = parseInt(
            document.querySelector('#chooseZeroconf1Items input:checked').value,
            10);
        selected1 = choices1[idx1];
      }
      if (choices2 != null && prompt2 != null) {
        idx2 = parseInt(
            document.querySelector('#chooseZeroconf2Items input:checked').value,
            10);
        selected2 = choices2[idx2];
      }
      console.log('Selected ' + idx1 + ' ' + idx2);
      onOK(selected1, selected2);
    };
    document.getElementById('chooseZeroconfRescanButton').onclick = function() {
      console.log('Rescan');
      onRescan();
    };

    //    let modal = new Modal({
    //      el: document.getElementById('chooseZeroconfModal'),
    //      backdrop: 'static'
    //    });
    //    console.log('modal', modal);
    //    modal.show();
    openModal('chooseZeroconfModal');
  },

  saveSYNDialog: async function() {
    let path = await wails.Dialogs.SaveFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Save state to SYN file',
      'Filters': [
        {DisplayName: 'State', Pattern: 'syn'},
        {DisplayName: 'All Files', Pattern: '*'}
      ]
    });
    console.log('in fileDialog: ' + path);

    if (path != undefined) {
      viewVCE_voice.connectSynergy(async function() {
        // Send message
        index.spinnerOn();
        try {
          await UIService.SaveSYN(path);
          index.infoNotification('Successfully saved Synergy state to ' + path);
        } catch (exc) {
          index.errorNotification(exc);
        }
        index.refreshConnectionStatus();
      });
    }
  },
  loadSYNDialog: async function() {
    let path = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Load state from SYN file',
      'Filters': [
        {DisplayName: 'State', Pattern: 'syn'},
        {DisplayName: 'All Files', Pattern: '*'}
      ]
    });
    console.log('in fileDialog: ' + path);
    if (path != undefined) {
      index.loadSYN(path[0], path[0]);
    }
  },
  loadCRTDialog: async function() {
    let path = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Load CRT Cartridge file',
      'Filters': [
        {DisplayName: 'Cartridge', Pattern: 'crt'},
        {DisplayName: 'All Files', Pattern: '*'}
      ]
    });
    console.log('in fileDialog: ' + path);
    if (path != undefined) {
      index.viewCRT(path[0], path[0]);
    }
  },
  loadVCEDialog: async function() {
    let path = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Load VCE Voice file',
      'Filters': [
        {DisplayName: 'Voice', Pattern: 'vce'},
        {DisplayName: 'All Files', Pattern: '*'}
      ]
    });
    console.log('in fileDialog: ' + path);
    if (path != undefined) {
      index.viewVCE(path[0], path[0]);
    }
  },
  saveVCEDialog: async function() {
    let path = await wails.Dialogs.SaveFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Save VCE Voice file',
      'Filters': [
        {DisplayName: 'Voice', Pattern: 'vce'},
        {DisplayName: 'All Files', Pattern: '*'}
      ]
    });
    console.log('in saveVCEDialog: ' + path);
    if (path != undefined) {
      index.spinnerOn();
      try {
        await UIService.SaveVCE(path);
        index.spinnerOff();
        index.infoNotification(
            'Successfully saved Synergy voice file to ' + path);
      } catch (exc) {
        index.spinnerOff();
        index.errorNotification(exc);
      }

      index.refreshConnectionStatus();
    }
  },
  loadSYN: function(name, path) {
    viewVCE_voice.connectSynergy(function() {
      index.confirmDialog('Load Synergy state file ' + path, async function() {
        // Send message
        index.spinnerOn();
        try {
          await UIService.LoadSYN(path);
          index.spinnerOff();
          index.infoNotification('Successfully loaded ' + name + ' to Synergy')
        } catch (exc) {
          index.spinnerOff();
          index.errorNotification(exc);
        }
        index.refreshConnectionStatus();
      })
    });
  },

  viewCRT: async function(name, path) {
    if (viewVCE_voice.voicingMode) {
      index.errorNotification('Can\'t load a CRT file while in Voicing mode');
      return;
    }
    try {
      let c = await UIService.ReadCRT(path)
      viewCRT.setCRT(path, name, c);
      index.load('viewCRT.html', 'content', function() {
        viewCRT.init();
      });
      index.refreshConnectionStatus();

    } catch (exc) {
      index.errorNotification(exc);
    }
  },
  viewVCE: function(name, path) {
    console.log('index.viewVCE ' + name + ' ' + path)
    if (viewVCE_voice.voicingMode) {
      index.confirmDialog(
          'Loading voice file will overwrite any pending edits - continue?',
          function() {
            index.raw_viewVCE(name, path);
          });
    }
    else {
      index.raw_viewVCE(name, path);
    }
  },
  raw_viewVCE: async function(name, path) {
    console.log('index.raw_viewVCE ' + name + ' ' + path)
    // Send message
    index.spinnerOn();
    try {
      let v;
      if (viewVCE_voice.voicingMode) {
        v = await UIService.LoadVceVoicingMode(path);
      } else {
        v = await UIService.ReadVCE(path);
      }
      index.spinnerOff();
      viewVCE.setVCE(v);

      viewCRT.setCRT(null, null);
      index.load('viewVCE.html', 'content', function() {
        viewVCE.init();
      });
      index.refreshConnectionStatus();
    } catch (exc) {
      index.spinnerOff();
      index.errorNotification(exc);
      return
    }
  },
  viewVCESlot: function(slot) {
    viewVCE.setVCE(viewCRT.crt.Voices[slot]);

    console.log('view voice slot ' + slot + ' : ' + viewVCE.vce);
    index.load('viewVCE.html', 'content', function() {
      viewVCE.init();
    });
  },
  addFolder: function(name, path) {
    let div = document.createElement('div');
    div.className = 'dir';
    div.onclick = function() {
      index.explore(path)
    };
    if (name == '..') name = '&lt;Parent&gt;';
    div.innerHTML = `<i class="fa fa-folder"></i><span>` + name + `</span>`;
    document.getElementById('dirs').appendChild(div)
  },
  addSYNFile: function(name, path) {
    let div = document.createElement('div');
    div.className = 'file';
    div.onclick = function() {
      index.loadSYN(name, path)
    };
    div.innerHTML = `<i class="fa fa-file"></i><span>` + name + `</span>`;
    document.getElementById('SYNfiles').appendChild(div)
  },
  addCRTFile: function(name, path) {
    let div = document.createElement('div');
    div.className = 'file';
    div.onclick = function() {
      index.viewCRT(name, path)
    };
    div.innerHTML = `<i class="fa fa-file"></i><span>` + name + `</span>`;
    document.getElementById('CRTfiles').appendChild(div)
  },
  addVCEFile: function(name, path) {
    let div = document.createElement('div');
    div.className = 'file';
    div.onclick = function() {
      index.viewVCE(name, path)
    };
    div.innerHTML = `<i class="fa fa-file"></i><span>` + name + `</span>`;
    document.getElementById('VCEfiles').appendChild(div)
  },
  explore: async function(path) {
    console.log('explore:', path);
    if (path == undefined) path = '';
    try {
      let exploration = await UIService.Explore(path);
      console.log('exploration:', exploration);

      // Process path
      document.getElementById('path').innerHTML = exploration.path;

      // Process dirs
      document.getElementById('dirs').innerHTML = ''
      for (let i = 0; i < exploration.dirs.length; i++) {
        index.addFolder(exploration.dirs[i].name, exploration.dirs[i].path);
      }

      document.getElementById('CRTfiles').innerHTML = ''
      if (exploration.CRTfiles.length > 0) {
        let div = document.createElement('div')
        div.innerHTML =
            '<div class=\'horizSeparator\'></div><b>Cartridge Files (.CRT)</b>';
        document.getElementById('CRTfiles').appendChild(div);

        for (let i = 0; i < exploration.CRTfiles.length; i++) {
          index.addCRTFile(
              exploration.CRTfiles[i].name, exploration.CRTfiles[i].path);
        }
      }

      document.getElementById('SYNfiles').innerHTML = ''
      if (exploration.SYNfiles.length > 0) {
        let div = document.createElement('div')
        div.innerHTML =
            '<div class=\'horizSeparator\'></div><b>Synergy State (.SYN)</b>';
        document.getElementById('SYNfiles').appendChild(div);
        for (let i = 0; i < exploration.SYNfiles.length; i++) {
          index.addSYNFile(
              exploration.SYNfiles[i].name, exploration.SYNfiles[i].path);
        }
      }

      document.getElementById('VCEfiles').innerHTML = ''
      if (exploration.VCEfiles.length > 0) {
        let div = document.createElement('div')
        div.innerHTML =
            '<div class=\'horizSeparator\'></div><b>Voice Files (.VCE)</b>';
        document.getElementById('VCEfiles').appendChild(div);
        for (let i = 0; i < exploration.VCEfiles.length; i++) {
          index.addVCEFile(
              exploration.VCEfiles[i].name, exploration.VCEfiles[i].path);
        }
      }
    } catch (err) {
      index.errorNotification(err);
      console.log('explore threw err: ', err);
    }
  },
  disconnectSynergy: function() {
    if (viewVCE_voice.voicingMode) {
      index.confirmDialog(
          'Disconnecting the Synergy will will discard any pending edits. Are you sure?',
          function() {
            viewVCE_voice.raw_voicingModeOff(true);
          });
    } else {
      index.raw_disconnectSynergy();
    }
  },

  raw_disconnectSynergy: async function() {
    index.spinnerOn();
    try {
      let status = await UIService.DisconnectSynergy();
      index.updateConnectionStatus(
          status.SynergyName, status.ControlSurfaceName);
      index.infoNotification('Disconnected Synergy');
      index.spinnerOff();
    } catch (exc) {
      index.spinnerOff();
      index.errorNotification(exc);
    }
  },

  disconnectControlSurface: async function() {
    index.spinnerOn();
    try {
      let status = await UIService.DisconnectControlSurface();
      index.spinnerOff();
      index.updateConnectionStatus(
          status.SynergyName, status.ControlSurfaceName);
      index.infoNotification('Disconnected Control Surface');
      document.querySelector('#disableControlSurfaceMenuItem')
          .classList.add('disabled');
    } catch (exc) {
      index.spinnerOff();
      index.errorNotification(exc);
    }
  },

  disableVRAM: function() {
    viewVCE_voice.connectSynergy(async function() {
      index.spinnerOff();
      try {
        await UIService.DisableVRAM();
        index.spinnerOff();
        index.infoNotification('Successfully disabled Synergy\'s VRAM')
      } catch (exc) {
        index.spinnerOff();
        index.errorNotification(exc);
      }
      index.refreshConnectionStatus();
    });
  },

  refreshConnectionStatus: async function() {
    console.log('refreshing connection status');
    try {
      let status = await UIService.GetConnectionStatus();
      index.updateConnectionStatus(
          status.SynergyName, status.ControlSurfaceName);
    } catch (exc) {
      index.errorNotification(exc);
    }
  },

  synergyName: null,
  controlSurfaceName: null,

  updateConnectionStatus: function(synergyName, csName) {
    index.synergyName =
        (synergyName == null || synergyName === '') ? null : synergyName;
    index.contronSurfaceName =
        (csName == null || csName === '') ? null : csName;

    console.log('update status: ' + synergyName + ' ' + csName);
    document.getElementById('synergyName').innerHTML = synergyName;
    document.getElementById('controlSurfaceName').innerHTML = csName;
    if (synergyName === null || synergyName === '') {
      document.getElementById('synergyName').innerHTML = 'not connected';
      document.querySelector('#disconnectSynergyMenuItem')
          .classList.add('disabled');
      document.querySelector('#connectSynergyMenuItem')
          .classList.remove('disabled');
      document.getElementById('connectButtonImg').src =
          `static/images/grey-button-off-full.png`;
    } else {
      document.querySelector('#disconnectSynergyMenuItem')
          .classList.remove('disabled');
      document.querySelector('#connectSynergyMenuItem')
          .classList.add('disabled');
      document.getElementById('connectButtonImg').src =
          `static/images/grey-button-on-full.png`;
    }
    if (csName === null || csName === '') {
      document.querySelector('#controlSurfaceStatus').style.display = 'none';
      document.querySelector('#disconnectControlSurfaceMenuItem')
          .classList.add('disabled');
    } else {
      document.querySelector('#controlSurfaceStatus').style.display = 'block';
      document.querySelector('#disconnectControlSurfaceMenuItem')
          .classList.remove('disabled');
    }
  },

  checkVersion:
      async function(synergyWasDisconnected, controlSurfaceWasDisconnected) {
        console.log(
            'checkVersion ' + synergyWasDisconnected + ' ' +
            controlSurfaceWasDisconnected);
        try {
          await UIService.CheckVersion(
              synergyWasDisconnected, controlSurfaceWasDisconnected)
        } catch (exc) {
          index.errorNotification(exc);
        }
      },

  //   fileDialog: function() {
  //     let files = dialog.showOpenDialogSync({
  //       // electron bug? filter files cause the dialog to look wonky
  //       filters: [
  //         {name: 'Voice', extensions: ['vce']},
  //         {name: 'Cartridge', extensions: ['crt']},
  //         {name: 'State', extensions: ['syn']},
  //         {name: 'All Files', extensions: ['*']}
  //       ],
  //       properties: ['openFile']
  //     });
  //     console.log('in fileDialog: ' + files);
  //     return files;
  //   },

  runCOMTST: async function() {
    viewVCE_voice.connectSynergy(async function() {
      index.spinnerOn();
      try {
        let status = await UIService.RunCOMTST();
        index.spinnerOff();
        console.log('runCOMTST returned: ' + status);
        index.infoNotification(status);
      } catch (exc) {
        index.spinnerOff();
        index.errorNotification(exc);
      }
    });
    index.refreshConnectionStatus();
  },

  load: async function(url, eleId, callback) {
    console.log(
        'load ' + url + ' into ' + eleId + ' ' +
        document.querySelector(('#' + eleId)));
    try {
      const r = await fetch(url);
      if (!r.ok) {
        throw new Error(`URL: ${url}: Response status: ${r.status}`);
      }
      const body = await r.text();
      document.querySelector('#' + eleId).innerHTML = body;
      if (callback != undefined) {
        let element = document.getElementById(eleId);
        callback(element);
      }
    } catch (exc) {
      console.log('ERROR: ', exc);
    }

    /*
    timing bug - onreadystatechange fires before the DOM is ready to query

    element = document.getElementById(eleId);
    req = new XMLHttpRequest();

    req.onreadystatechange = function () {
            if (this.readyState == 4 && this.status == 200) {
                    element.innerHTML = req.responseText;
                    if (callback != undefined) {
                            callback(element);
                    }
            }
    };

    req.open("GET", url, false);
    req.send(null);
    */
  },
  dropdownMenu: function(contentId) {
    // console.log("toggle display on " + contentId);
    document.getElementById(contentId).style.display = 'block';
  },

  viewDiag: function() {
    index.load('diag.html', 'content');
  },
  showAbout: async function() {
    try {
      await UIService.ShowAbout();
    } catch (err) {
      index.errorNotification(err);
      console.log('show about threw err: ', err);
    }
  },
  showPreferences: async function() {
    try {
      await UIService.ShowPreferences();
    } catch (err) {
      index.errorNotification(err);
      console.log('show preferences threw err: ', err);
    }
  },

  // debounce a function separately for each "first" argument - we use
  // this with first argument being the input ele being debounced -
  // this allows each input to be independently debounced even if all
  // using the same onchange function Adapted from:
  // https://github.com/lodash/lodash/issues/2403 and
  // https://stackoverflow.com/a/28795512
  debounceFirstArg: function(func, wait = 0, options = {}) {
    let mem = _.memoize(function() {
      return _.debounce(func, wait, options)
    });
    return function() {
      mem.apply(this, arguments).apply(this, arguments)
    }
  }
};

/***
listen: function () {
console.log("index listening...")
astilectron.onMessage(function (message) {
        switch (message.name) {
                case "explore":
                        index.explore(message.payload);
                        return { payload: "ok" };
                case "updateConnectionStatus":
                        index.updateConnectionStatus(message.payload.SynergyName,
message.payload.ControlSurfaceName); return { payload: "ok" }; case
"fileDialog": f = index.fileDialog(message.payload); return { payload: f };
                        break;
                case "viewVCE":
                        console.log("viewVCE: " +
JSON.stringify(message.payload)); vce = message.payload; index.load("view.html",
"content", function () { viewVCE.init();
                                });

                        return { payload: "ok" };
                        break;
                case "runDiag":
                        index.viewDiag();
                        return { payload: "ok" };
                        break;
                case "updateFromCSurface":
                        valueString =
viewVCE_voice.updateFromCSurface(message.payload) return { payload: valueString
};

                case "dx2synAddProcessLog":
                        console.log("dx2synAddProcessLog  - " + message.payload)
                        //????dx2syn.addProcessLog(message.payload);
                        return { payload: "ok" };

                case "dx2synFinish":
                        console.log("dx2synFinish  - " + message.payload)
                        //?????dx2syn.finishConvert(message.payload);
                        return { payload: "ok" };

        }
});
}
**/


// make the letiable visible to HTML:
window.index = index;

function inDropbtn(ele) {
  if (ele == null) {
    // console.log("ele is null");
    return false
  } else if (ele.classList.contains('dropbtn')) {
    // console.log("ele has dropbtn " + JSON.stringify(ele));
    return true;
  }
  // console.log("ele doesnt have dropbtn - try parent " + JSON.stringify(ele) +
  // " " + JSON.stringify(ele.parentElement));
  return inDropbtn(ele.parentElement);
}

/* close dropdowns if user clicks outside the menu */
window.onclick =
    function(event) {
  if (!inDropbtn(event.target)) {
    let dropdowns = document.getElementsByClassName('dropdown-content');
    let i;
    for (i = 0; i < dropdowns.length; i++) {
      let openDropdown = dropdowns[i];
      if (openDropdown.style.display === 'block') {
        // console.log("toggle display off " + openDropdown.id);
        openDropdown.style.display = 'none';
      }
    }
  }
}


    wails.Events.On('explore', (path) => {
      console.log('explore event: ', path);
      index.explore(path[0])
    });
