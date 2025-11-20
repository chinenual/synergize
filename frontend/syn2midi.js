import {UIService} from '/bindings/github.com/chinenual/synergize';
import * as wails from '@wailsio/runtime';

import {index} from './index';
import $ from 'jquery';
Object.assign(window, { $: $, jQuery: $ });

// Import all of Bootstrap's JS
import * as bootstrap from 'bootstrap';

export let syn2midi = {

  init: function() {},

  lastTempo: 120.0,
  lastRaw: false,
  lastMaxClock: 2 * 60,

  synpathDialog: async function(ele) {
    let path = await wails.Dialogs.OpenFile({
      'CanChooseDirectories': false,
      'CanChooseFiles': true,
      'Title': 'Choose SYN file to convert to MIDI',
      'Filters': [
        {DisplayName: 'State', Pattern: 'syn'},
        {DisplayName: 'All Files', Pattern: '*'}
      ]
    });

    // console.log("in folderDialog: " +
    // folder);
    if (path != undefined && ele != undefined && ele != null) {
      ele.value = path;

      syn2midi.getSynSequencerState(path[0]);
    }
    return path;
  },

  modeChange: function() {
    let raw = !document.getElementById('syn2midiMode').checked;
    console.log('mode change ' + raw);
    if (raw) {
      document.querySelector('#syn2midiMaxClockDiv').style.display = 'none';
      document.querySelector('#syn2midiTrackModeDiv').style.display = 'none';
    } else {
      document.querySelector('#syn2midiMaxClockDiv').style.display = 'block';
      document.querySelector('#syn2midiTrackModeDiv').style.display = 'block';
    }
  },

  convertDialog: function() {
    document.getElementById('syn2midiTempo').value = syn2midi.lastTempo;
    document.getElementById('syn2midiMode').checked = !syn2midi.lastRaw;
    document.getElementById('syn2midiMaxClock').value = syn2midi.lastMaxClock;
    syn2midi.modeChange();
    document.getElementById('syn2midiConvertButton').onclick = function() {
      if (document.getElementById('syn2midiPath').value != undefined &&
          document.getElementById('syn2midiTempo').value != undefined) {
        let buttons = [
          parseInt(document.getElementById('syn2midiTrack1').value, 10),
          parseInt(document.getElementById('syn2midiTrack2').value, 10),
          parseInt(document.getElementById('syn2midiTrack3').value, 10),
          parseInt(document.getElementById('syn2midiTrack4').value, 10)
        ]
        syn2midi.lastTempo = document.getElementById('syn2midiTempo').value
        syn2midi.runConvert(
            document.getElementById('syn2midiPath').value,
            parseInt(document.getElementById('syn2midiTempo').value, 10),
            !document.getElementById('syn2midiMode').checked,
            parseInt(document.getElementById('syn2midiMaxClock').value, 10),
            buttons);
      }
    };
    let modal = new bootstrap.Modal(document.getElementById('syn2midiModal'), {backdrop: 'static'});
    console.log("modal", modal);
    modal.show();
    
  },

  getSynSequencerState: async function(path) {
    console.log('call getSynSequencerState: ' + path);
    let trackButtons;
    try {
      trackButtons = UIService.GetSynSequencerState(path);
    } catch (exc) {
      index.errorNotification(exc);
    }
    document.getElementById('syn2midiTrack1').value = '' + trackButtons[0];
    document.getElementById('syn2midiTrack2').value = '' + trackButtons[1];
    document.getElementById('syn2midiTrack3').value = '' + trackButtons[2];
    document.getElementById('syn2midiTrack4').value = '' + trackButtons[3];
  },

  runConvert: async function(path, tempo, raw, maxClockSeconds, buttons) {
    console.log('runConvert: ' + path + ' ' + tempo);
    try {
      UIService.Syn2midi(
          path, parseFloat(tempo), raw, maxClockSeconds, buttons);
      index.infoNotification(
          'Successfully converted Synergy sequence data to ' + path + '.mid');

    } catch (exc) {
      index.errorNotification(exc);
    }
  },
};
window.syn2midi = syn2midi;
