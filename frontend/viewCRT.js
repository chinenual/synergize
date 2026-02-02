import {UIService} from '/bindings/github.com/chinenual/synergize';
import * as wails from '@wailsio/runtime';

import { index } from './index';
import { viewVCE_voice } from './viewVCE_voice';

export let viewCRT = {
  editMode: false,

  crt: {},
  crt_path: '',
  crt_name: '',

  makeSlotOnclick(slot) {
    return function() {
      index.viewVCESlot(slot);
    }
  },

  setCRT: function(path, name, data) {
    viewCRT.crt_path = path;
    viewCRT.crt_name = name;
    viewCRT.crt = data;
  },

  editCRT: function(name, path) {
    if (viewCRT.editMode) {
      index.confirmDialog(
          'Disabling edit mode will discard any pending edits. Are you sure?',
          function() {
            viewCRT.editMode = false;
            viewCRT.reinit();
          });
    } else {
      viewVCE_voice.connectSynergy(function() {
        viewCRT.editMode = true;
        viewCRT.reinit();
      });
    }
  },

  add: async function(slot) {
    let path = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Add Voice',
      'Filters': [
        {DisplayName: 'Voice', Pattern: 'vce'},
        {DisplayName: 'All Files', Pattern: '*'}
      ]
    });

    console.log('in fileDialog: ' + path);
    if (path != undefined) {
      let crt;
      try {
        crt = UIService.CrtEditAddVoice(viewCRT.crt, path[0]. slot);
      } catch (exc) {
        index.errorNotification(exc);
      }
      viewCRT.crt = crt;
      viewCRT.reinit();
      index.refreshConnectionStatus();
    }
  },

  clear: function(slot) {
    let ele = document.getElementById('crt_voicename_' + slot);
    ele.innerHTML = '';
    ele.onclick = function() {
      // nop
    };
    viewCRT.crt.Voices[slot - 1] = null;
  },

  loadCRT: async function () {
    if (viewCRT.crt_path != undefined) {
      viewVCE_voice.connectSynergy(function () {
        index.spinnerOn();
        try {
          UIService.CrtEditLoadCRT(viewCRT.crt);
          index.infoNotification('Successfully loaded CRT to Synergy');
        } catch (exc) {
          index.errorNotification(exc);
        }
        index.spinnerOff();
        index.refreshConnectionStatus();
      });
    }
  },
        
  saveCRT: async function(name, path_ignored) {
    let path = await wails.Dialogs.SaveFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Save CRT',
      'Filters': [
        {DisplayName: 'Cartridge', Pattern: 'CRT'},
        {DisplayName: 'All Files', Pattern: '*'}
      ]
    });
    console.log('in fileDialog: ' + path);

    if (path != undefined) {
      viewVCE_voice.connectSynergy(function () {
        index.spinnerOn();
        try {
          UIService.CrtEditSaveCRT(viewCRT.crt, path);
          index.infoNotification('Successfully saved CRT to ' + path);
        } catch (exc) {
          index.errorNotification(exc);
        }
        index.spinnerOff();
        index.refreshConnectionStatus();
      });
    }
  },

  viewLoadedCRT: function() {
    index.load('viewCRT.html', 'content', function() {
      viewCRT.reinit();
    });
  },

  reinit: function() {
    console.log('view CRT ' + viewCRT.crt_name)
    document.getElementById('crt_path').innerHTML = viewCRT.crt_name;
    // clear everything
    for (let i = 0; i < 24; i++) {
      let ele = document.getElementById('crt_voicename_' + (i + 1));
      ele.innerHTML = '';
      ele.onclick = function() {
        // nop
      };
    }
    // set active voices
    for (let i = 0; i < viewCRT.crt.Voices.length; i++) {
      let ele = document.getElementById('crt_voicename_' + (i + 1));
      if (viewCRT.crt.Voices[i] != null) {
        ele.innerHTML = viewCRT.crt.Voices[i].Head.VNAME;
        ele.onclick = viewCRT.makeSlotOnclick(i);
      }
    }

    if (viewCRT.editMode) {
      document.getElementById('editCRTButtonImg').src =
          `static/images/red-button-on-full.png`;
      document.querySelector('#saveCRTButtonDiv').style.display = 'block';
      document.querySelectorAll('.crtSlotAddButton').forEach(el => { el.style.display = 'block'; });
      document.querySelectorAll('.crtSlotClearButton').forEach(el => { el.style.display = 'block'; });
    } else {
      document.getElementById('editCRTButtonImg').src =
          `static/images/red-button-off-full.png`;
      document.querySelector('#saveCRTButtonDiv').style.display = 'none';
      document.querySelectorAll('.crtSlotAddButton').forEach(el => { el.style.display = 'none'; });
      document.querySelectorAll('.crtSlotClearButton').forEach(el => { el.style.display = 'none'; });
    }
  },

  init: function() {
    viewCRT.editMode = false;
    viewCRT.reinit();
  }
};

// make the letiable visible to HTML:
window.viewCRT = viewCRT;