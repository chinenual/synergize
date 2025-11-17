import {UIService} from '/bindings/github.com/chinenual/synergize';
import * as wails from '@wailsio/runtime';
import { index } from './index';

import { $ } from 'jquery';

export let dx2syn = {
  init : function() {
    // if (process.platform == 'darwin') {
    // 	// Macos dialog can select both folders and files
    // 	document.getElementById('dx2synFileMenuItem').hidden = true
    // 	document.getElementById('dx2synDirMenuItem').hidden = true
    // 	document.getElementById('dx2synEitherMenuItem').hidden = false
    // } else {
    // 	// windows and linux need a specific file and directory variant
    // 	document.getElementById('dx2synFileMenuItem').hidden = false
    // 	document.getElementById('dx2synDirMenuItem').hidden = false
    // 	document.getElementById('dx2synEitherMenuItem').hidden = true
    // }
  },

  convertEitherDialog: async function() {
    let path = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': true,
      'CanChooseFiles': true,
      'Title': 'Choose DX7 Sysex or folder containing DX7 Sysex\'s',
      'Filters': [
        {DisplayName: 'DX Sysex', Pattern: ['syx', 'sysx']},
        {DisplayName: 'All Files', Pattern: ['*']}
      ]
    });
    console.log('in convertFileDialog: ' + path);
    if (path != undefined) {
      dx2syn.runConvert(path[0]);
    }
  },

  convertFileDialog: async function() {
    let path = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Choose DX7 Sysex',
      'Filters': [
        {DisplayName: 'DX Sysex', Pattern: ['syx', 'sysx']},
        {DisplayName: 'All Files', Pattern: ['*']}
      ]
    });
    console.log('in convertFileDialog: ' + path);
    if (path != undefined) {
      dx2syn.runConvert(path[0]);
    }
  },

  convertFolderDialog: async function() {
    let path = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': true,
      'CanChooseFiles': false,
      'Title': 'Choose folder containing DX7 Sysex\'s'
    });
    console.log('in convertFolderDialog: ' + path);
    if (path != undefined) {
      dx2syn.runConvert(path[0]);
    }
  },

  runConvert: async function(path) {
    console.log('runConvert: ' + path)
    document.getElementById('subprocessTitle').innerHTML = 'Convert DX7 Sysex';
    document.getElementById('logOutput').innerHTML = '';
    document.getElementById('subprocessCloseButton')
        .setAttribute('disabled', 'disabled');
    document.getElementById('subprocessCancelButton')
        .removeAttribute('disabled');

    document.getElementById('subprocessCancelButton').onclick =
        async function() {
      console.log('SAW CANCEL');
      // Send message
      console.log('call dx2synCancel: ' + path);
      try {
        await UIService.Dx2synCancel();
      } catch (exc) {
        index.errorNotification(exc);
      }
      document.getElementById('subprocessCloseButton')
          .removeAttribute('disabled');
      document.getElementById('subprocessCancelButton')
          .setAttribute('disabled', 'disabled');
    };

    $('#subprocessModal').modal({
      backdrop: 'static'  // clicking outside the dialog doesnt close the dialog
    });
	  try {
		  await UIService.Dx2synStart(path);
	} catch (exc) {
          index.errorNotification(exc);		
	}
  },

  finishConvert: function(msg) {
    dx2syn.addProcessLog('\n' + msg);
    document.getElementById('subprocessCloseButton')
        .removeAttribute('disabled');
    document.getElementById('subprocessCancelButton')
        .setAttribute('disabled', 'disabled');
  },

  addProcessLog: function(msgs) {
    let html = msgs.replaceAll('\n', '<br/>\n');
    if (html[html.length - 1] != '\n') {
      html += '<br/>\n';
    }
    document.getElementById('logOutput').innerHTML =
        document.getElementById('logOutput').innerHTML + html;
    $('#subprocessModal').modal('handleUpdate')
  },
};
// for menus:
window.dx2syn = dx2syn;