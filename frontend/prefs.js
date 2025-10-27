import {UIService} from '/bindings/github.com/chinenual/synergize';
import * as wails from '@wailsio/runtime';

export let prefs = {
  init: async function() {
    let preferences;
    let os
    try {
      let result = await UIService.GetPreferences();
      os = result[0];
      preferences = result[1];
    } catch (err) {
      console.log('getpref err', err);
      index.errorNotification(err);
      return
    }

    document.getElementById('useSerial').checked =
        preferences.UseSerial ? 'checked' : '';
    document.getElementById('serialPort').value = preferences.SerialPort;
    document.getElementById('serialBaud').value = preferences.SerialBaud;
    document.getElementById('flowControl').checked =
        preferences.SerialFlowControl;

    document.getElementById('libraryPath').value = preferences.LibraryPath;

    document.getElementById('useOsc').checked =
        preferences.UseOsc ? 'checked' : '';
    document.getElementById('oscAutoConfig').checked =
        preferences.OscAutoConfig ? 'checked' : '';
    document.getElementById('oscPort').value = preferences.OscPort;
    document.getElementById('oscCSurfaceAddress').value =
        preferences.OscCSurfaceAddress;
    document.getElementById('oscCSurfacePort').value =
        preferences.OscCSurfacePort;
    prefs.toggleOsc();

    if (os === 'darwin') {
      /*
        * as nice as this would be, macos hides the /dev directory
        * from the UI dialogs.  Best to just let folk type it in...
        *

      // add an onclick handler to popup a file dialog
      var ele = document.getElementById("serialPort");
      ele.onclick = function() {
      prefs.serialPortDialog(this,this.value);
      }

      *
      */
    } else {
      // on windows, use a straight text box
    }
  },

  toggleOsc: function() {
    var useSerialChecked = document.getElementById('useSerial').checked;
    var useOscChecked = document.getElementById('useOsc').checked;
    var autoChecked = document.getElementById('oscAutoConfig').checked;


    document.getElementById('serialPort').disabled = (!useSerialChecked);
    document.getElementById('serialBaud').disabled = (!useSerialChecked);
    document.getElementById('flowControl').disabled = (!useSerialChecked);

    document.getElementById('oscPort').disabled = (!useOscChecked);
    document.getElementById('oscAutoConfig').disabled = (!useOscChecked);
    document.getElementById('oscCSurfaceAddress').disabled =
        (!useOscChecked) || autoChecked;
    document.getElementById('oscCSurfacePort').disabled =
        (!useOscChecked) || autoChecked;
  },



  libraryFolderDialog: async function(ele, defaultValue) {
    let folder = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': true,
      'CanChooseFiles': false,
      'Title': 'Choose Voice Library Path',
      'Directory': defaultValue
    });
    console.log('folder', folder);
    if (folder != undefined && ele != undefined && ele != null) {
      ele.value = folder;
    }
    return folder;
  },

  cancelAndClose: async function() {
    try {
      UIService.CancelPreferences()
    } catch (err) {
      console.log('CANCEL pref err', err);
      index.errorNotification(err);
    }
  },

  saveAndClose: async function() {
      preferences = {
        'UseSerial': document.getElementById('useSerial').checked,
        'SerialPort': document.getElementById('serialPort').value,
        'SerialBaud': parseInt(document.getElementById('serialBaud').value, 10),
        'SerialFlowControl': document.getElementById('flowControl').checked,
        'LibraryPath': document.getElementById('libraryPath').value,
        'UseOsc': document.getElementById('useOsc').checked,
        'OscAutoConfig': document.getElementById('oscAutoConfig').checked,
        'OscPort': parseInt(document.getElementById('oscPort').value, 10),
        'OscCSurfaceAddress':
            document.getElementById('oscCSurfaceAddress').value,
        'OscCSurfacePort':
            parseInt(document.getElementById('oscCSurfacePort').value, 10)
      };
    try {
      console.log('before save');
      await UIService.SavePreferences(preferences);
      console.log('after save');
    } catch (err) {
      index.errorNotification(err);
    }
  }
};
window.prefs = prefs;