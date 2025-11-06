import * as index from './index';
import {$} from './jquery-3.4.1.min';
import * as viewVCE_voice from './viewVCE_voice';
import * as wails from '@wailsio/runtime';

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
        {DisplayName: 'Voice', Pattern: ['vce']},
        {DisplayName: 'All Files', Pattern: ['*']}
      ]
    });

    console.log('in fileDialog: ' + path);
    if (path != undefined) {
      let message = {
        'name': 'crtEditAddVoice',
        'payload': {'VcePath': path[0], 'Slot': slot, 'Crt': viewCRT.crt}
      };
      console.dir(message.payload);
      astilectron.sendMessage(message, function(message) {
        // Check error
        if (message.name === 'error') {
          index.errorNotification(message.payload);
          return
        }
        console.dir(message.payload);
        viewCRT.crt = message.payload;
        viewCRT.reinit();
        index.refreshConnectionStatus();
      });
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

  loadCRT: function() {
    if (viewCRT.crt_path != undefined) {
      viewVCE_voice.connectSynergy(function() {
        let message = {'name': 'crtEditLoadCRT', 'payload': {Crt: viewCRT.crt}};
        // Send message
        index.spinnerOn();
        astilectron.sendMessage(message, function(message) {
          index.spinnerOff();
          // Check error
          if (message.name === 'error') {
            index.errorNotification(message.payload);
          } else {
            index.infoNotification('Successfully loaded CRT to Synergy');
          }
          index.refreshConnectionStatus();
        });
      });
    }
  },

  saveCRT: async function(name, path_ignored) {
   let path = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Save CRT',
      'Filters': [
        {DisplayName: 'Cartridge', Pattern: ['CRT']},
        {DisplayName: 'All Files', Pattern: ['*']}
      ]
    });
    console.log('in fileDialog: ' + path);

    if (path != undefined) {
      viewVCE_voice.connectSynergy(function() {
        let message = {
          'name': 'crtEditSaveCRT',
          'payload': {Path: path, Crt: viewCRT.crt}
        };
        // Send message
        index.spinnerOn();
        astilectron.sendMessage(message, function(message) {
          index.spinnerOff();
          // Check error
          if (message.name === 'error') {
            index.errorNotification(message.payload);
          } else {
            index.infoNotification('Successfully saved CRT to ' + path);
          }
          index.refreshConnectionStatus();
        });
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
      $('#saveCRTButtonDiv').show();
      $('.crtSlotAddButton').show();
      $('.crtSlotClearButton').show();
    } else {
      document.getElementById('editCRTButtonImg').src =
          `static/images/red-button-off-full.png`;
      $('#saveCRTButtonDiv').hide();
      $('.crtSlotAddButton').hide();
      $('.crtSlotClearButton').hide();
    }
  },

  init: function() {
    viewCRT.editMode = false;
    viewCRT.reinit();
  }
};
