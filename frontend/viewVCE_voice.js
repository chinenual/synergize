import {UIService} from '/bindings/github.com/chinenual/synergize';
import {TouchSpin} from '@touchspin/core';
import VanillaRenderer from '@touchspin/renderer-vanilla';
import {Chart} from 'chart.js';
import _ from 'lodash';
import * as nomnoml from 'nomnoml';

import {index} from './index';
import {viewCRT} from './viewCRT';
import {viewVCE} from './viewVCE';
import {viewVCE_envs} from './viewVCE_envs';
import {viewVCE_filters} from './viewVCE_filters';

export let viewVCE_voice = {
  voicingMode: false,
  csEnabled: false,

  timbreProportionCurve: function(center, sensitivity) {
    let result = [];
    if (sensitivity == 0) {
      for (let v = 0; v < 32; v++) {
        result[v] = center;
      }
      return result;
    }
    // center = 0..32
    // sensitivity = 1..31
    for (let v = 0; v < 32; v++) {
      // this appears to be what the z80 code is doing for timbre PROPC:
      let p = (center * 2) - 15 + ((v / 10) * sensitivity) - (2 * sensitivity);
      if (p > 31) p = 31;
      if (p < 0) p = 0;
      result[v] = p;
    }
    return result;
  },

  ampProportionCurve: function(center, sensitivity) {
    let result = [];
    if (sensitivity == 0) {
      for (let v = 0; v < 32; v++) {
        result[v] = center;
      }
      return result;
    }
    // center = 0..32
    // sensitivity = 1..31
    for (let v = 0; v < 32; v++) {
      // this appears to be what the z80 code is doing for timbre PROPC:
      let p =
          ((((v / 10) * sensitivity) / 2.0 - sensitivity) / 2.0) + (center - 24)
      if (p > 6) p = 6;
      if (p < -24) p = -24;
      result[v] = p + 25;
    }
    return result;
  },

  SOLO: [],
  MUTE: [],

  toggleOsc: async function(ele) {
    console.log('toggle ' + ele.id);
    let oscPattern = /([A-Z]+)\[(\d+)\]/;
    let ret = ele.id.match(oscPattern);
    if (ret) {
      let param = ret[1];
      let osc = parseInt(ret[2], 10); /* 1-based */

      let state = !ele.classList.contains('on');
      viewVCE_voice[param][osc - 1] = state;
      ele.classList.toggle('on');

      try {
        await UIService.SetOscSolo(viewVCE_voice.MUTE, viewVCE_voice.SOLO);
        // return array is ignored
      } catch (exc) {
        // failed - dont change the boolean
        index.errorNotification(exc);
        return false;
      }
      console.log('send to csurface ' + ele.id + ' ' + state)
      viewVCE_voice.sendToCSurface(ele, ele.id, state ? 1 : 0);
    }
  },

  filterChanged: function(ele) {
    let id = ele.id;
    console.log('filterChanged: ' + id + ' val: ' + ele.value);

    let filterPattern = /FILTER\[(\d+)\]/;
    let ret;
    let osc;
    if ((ret = id.match(filterPattern))) {
      osc = parseInt(ret[1])
    } else {
      console.log('ERROR: filterCHanged called with bad ele ' + ele);
      osc = 0;
    }
    let filterValue = parseInt(ele.value, 10);
    viewVCE.vce.Head.FILTER[osc - 1] = filterValue;
    // don't wait for this - we want to get feedback to the control surface asap
    // async () => {
    viewVCE_filters.init(true);
    //}
  },

  OHARMToText: function(str) {
    let newStr;
    let val = parseInt(str, 10);
    if (val == -12) {
      newStr = 'dc';
    } else if (val < 0) {
      // bug#79:  -1 (ff)  -> s1, -2 (fe) -> s2
      // correct: -1 (ff) -> s11, -2 (fe) -> s10
      newStr = 's' + (12 + val);
    } else {
      newStr = str;
    }
    // console.log("   OHARMToText(" + str + ") returns " + newStr);
    return newStr;
  },

  TextToOHARM: function(str) {
    let newStr;
    let ret;
    if (str === 'dc') {
      newStr = '-12';
    } else if ((ret = str.match(/s(\d+)/))) {
      // bug#79:  s1 -> -1 (ff), s2 -> -2 (fe)
      // correct: s1 -> -11 (f5), s2 -> -10 (f6)
      //          s11 -> -1 (ff),  s10 -> -2 (fe)
      let val = parseInt(ret[1], 10);
      newStr = '' + (val - 12);
    } else if ((ret = str.match(/\d+/))) {
      newStr = str;
    } else {
      /// error!
      console.log('ERROR: TextToOHARM cant decode ' + str);
      newStr = str;
    }
    // console.log("   TextToOHARM(" + str + ") returns " + newStr);
    return newStr;
  },

  FDETUNToText: function(str) {
    /* THIS IS UGLY.  The Synergy uses a non-linear mapping of this byte to a
     * 1/30Hz increment, with 5 positive values reserved for "random" settings.
     * Ick ick ick.
     *
     * See D.DTN routine in VOIDSP.Z80 in the SYNHCS sourcecode.
     */
    let newStr;
    let val = parseInt(str, 10);
    if (val > 58) {
      // CASE A
      newStr = 'ran' + (val - 58);
    } else if (val >= -32 && val <= 32) {
      // CASE B
      val = val * 3;
      newStr = '' + val;
    } else if (val > 0) {
      // CASE C
      val = ((val * 2) - 32) * 3;
      newStr = '' + val;
    } else {  // negative
      // CASE D
      val = ((val * 2) + 32) * 3;
      newStr = '' + val;
    }
    // console.log("   FDETUNToText(" + str + ") returns " + newStr);
    return newStr;
  },

  TextToFDETUN: function(str) {
    // See FDETUNToText.  This "reverses" that attrocity
    let newStr;
    let ret;
    if ((ret = str.match(/ran(\d+)/))) {
      // CASE A
      let val = parseInt(ret[1], 10);
      val += 58;
      newStr = '' + val;
    } else {
      let val = parseInt(str, 10);
      if (val >= (-32 * 3) && val <= (32 * 3)) {
        // CASE B
        val /= 3;
      } else if (val > 0) {
        // CASE C
        val = ((val / 3) + 32) / 2;
      } else {
        // CASE D
        val = ((val / 3) - 32) / 2;
      }
      newStr = '' + Math.round(val);
    }
    // console.log("   TextToFDETUN(" + str + ") returns " + newStr);
    return newStr;
  },

  NullablePatchRegisterToText: function(str) {
    if (str.trim() === '') {
      return '';
    }
    let val = parseInt(str, 10);
    if (val === 0) {
      return '';
    }
    return '' + val;
  },

  TextToNullablePatchRegister: function(str) {
    if (str.trim() === '') {
      return '0';
    }
    let val = parseInt(str, 10);
    if (val === 0) {
      return '0'
    }
    return '' + val;
  },

  testConversionFunctions: function() {
    let ok = true;
    let result = [];
    for (let i = 0; i < 3; i++) {
      let str = viewVCE_voice.NullablePatchRegisterToText('' + i);
      let reverseStr = viewVCE_voice.TextToNullablePatchRegister(str);
      if (('' + i) != reverseStr) {
        ok = false;
        let err = 'ERROR: PatchReg ' + i + ' totext: ' + str + ' reversed to ' +
            reverseStr;
        result.push(err);
        console.log(err);
      }
    }
    for (let i = -12; i <= 30; i++) {
      let str = viewVCE_voice.OHARMToText('' + i);
      let reverseStr = viewVCE_voice.TextToOHARM(str);
      if (('' + i) != reverseStr) {
        ok = false;
        let err = 'ERROR: OHARM ' + i + ' totext: ' + str + ' reversed to ' +
            reverseStr;
        result.push(err);
        console.log(err);
      }
    }
    for (let i = -63; i <= 63; i++) {
      let str = viewVCE_voice.FDETUNToText('' + i);
      let reverseStr = viewVCE_voice.TextToFDETUN(str);
      if (('' + i) != reverseStr) {
        ok = false;
        let err = 'ERROR: FDETUN ' + i + ' totext: ' + str + ' reversed to ' +
            reverseStr;
        result.push(err);
        console.log(err);
      }
    }
    // if (val >= (-32 * 3) && val <= (32 * 3)) {
    //	// CASE B
    // } else if (val > 0) {
    //	// CASE C
    // } else {
    //	// CASE D
    // }
    //  test that the FDETUN rounding does the right thing:
    let v = viewVCE_voice.FDETUNToText(viewVCE_voice.TextToFDETUN('20'));
    if ('21' != v) {
      ok = false;
      let err = 'ERROR: FDETUN CASE B ' + 20 + ' rounded to ' + v +
          ' - expected ' + 21;
      result.push(err);
      console.log(err);
    }
    v = viewVCE_voice.FDETUNToText(viewVCE_voice.TextToFDETUN('247'));
    if ('246' != v) {
      ok = false;
      let err = 'ERROR: FDETUN CASE C ' + 247 + ' rounded to ' + v +
          ' - expected ' + 246;
      result.push(err);
      console.log(err);
    }
    v = viewVCE_voice.FDETUNToText(viewVCE_voice.TextToFDETUN('-205'));
    if ('-204' != v) {
      ok = false;
      let err = 'ERROR: FDETUN CASE D ' + -205 + ' rounded to ' + v +
          ' - expected ' + -204;
      result.push(err);
      console.log(err);
    }

    console.log(
        'viewVCE_voice.testConversionFunctions: ' + (ok ? 'PASS' : 'FAIL'));
    return ok ? null : result;
  },

  onchangeDSR: async function(param, osc /*1-based*/, value) {
    osc = parseInt(osc, 10);

    console.log('onchangeDSR: ' + param + '[' + osc + '] == ' + value);

    // displayed values are 1-based, bit values in the patch byte are 0-based
    let patchFOInputDSR =
        document.getElementById(`patchFOInputDSR[${osc}]`).value;
    let patchAdderInDSR =
        document.getElementById(`patchAdderInDSR[${osc}]`).value;
    let patchOutputDSR =
        document.getElementById(`patchOutputDSR[${osc}]`).value;

    // XREF: patch byte encode/decode
    let patchInhibitAddr = patchAdderInDSR == '' ? true : false;
    let patchInhibitF0 = patchFOInputDSR == '' ? true : false;
    let patchByte = 0;
    patchByte |= ((parseInt(patchFOInputDSR, 10) - 1) & 0x03);
    patchByte |= (((parseInt(patchAdderInDSR, 10) - 1) << 3) & 0x18);
    patchByte |= (((parseInt(patchOutputDSR, 10) - 1) << 6) & 0xc0);
    if (patchInhibitAddr) {
      patchByte |= 0x20;
    }
    if (patchInhibitF0) {
      patchByte |= 0x04;
    }

    console.log(
        osc + ' old patch byte: ' +
        viewVCE.vce.Envelopes[osc - 1].FreqEnvelope.OPTCH + '\n' +
        ' new patch byte: ' + patchByte + '\n' +
        ' patchInhibitAddr : ' + patchInhibitAddr + '\n' +
        ' patchInhibitF0   : ' + patchInhibitF0 + '\n' +
        ' patchOutputDSR   : ' + patchOutputDSR + '\n' +
        ' patchAdderInDSR  : ' + patchAdderInDSR + '\n' +
        ' patchFOInputDSR  : ' + patchFOInputDSR + '\n');

    try {
      await UIService.SetVoiceOscDataByte(osc, patchByte);
    } catch (exc) {
      // failed - dont change the boolean
      index.errorNotification(exc);
      return false;
    }
    viewVCE.vce.Envelopes[osc - 1].FreqEnvelope.OPTCH = patchByte;
    viewVCE_voice
        .patchTable();  // in case the patch diagram changes due to the edit
    viewVCE_voice.voicingModeVisuals();
    return true;
  },

  onchange: function(ele, updater, valueConverter) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    viewVCE_voice.deb_onchange(ele, updater, valueConverter);
  },

  deb_onchange: null,

  raw_onchange: async function(ele, updater, valueConverter) {
    if (viewVCE
            .supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/
      return;
    }
    let id = ele.id;
    console.log('changed: ' + id + ', new value: ' + ele.value);

    if (valueConverter == undefined) {
      // identity function if none passed
      valueConverter = function(v) {
        return v
      };
    }

    let value
    let param
    let args
    let osc
    let gainPattern = /OscGain\[(\d+)\]/;
    let filterPattern = /FILTER\[(\d+)\]/;
    let dsrPattern = /([0-9A-Za-z]+DSR)\[(\d+)\]/;
    let waveKeyPattern = /wk([A-Z]+)\[(\d+)\]/;
    let oscPattern = /([A-Z]+)\[(\d+)\]/;
    let headPattern = /([A-Z]+)/;
    let funcname;
    let ret;
    if (id === 'VNAME') {
      param = 'VNAME'
      funcname = 'setVNAME';
      args = ele.value;
    } else if ((ret = id.match(gainPattern))) {
      osc = parseInt(ret[1])
      value = parseInt(ele.value, 10);
      viewVCE_envs.setOscGain(osc - 1, value);
      viewVCE_voice.sendToCSurface(ele, ele.id, value);
      return

    } else if ((ret = id.match(dsrPattern))) {
      param = ret[1]
      osc = ret[2]
      value = parseInt(valueConverter(ele.value), 10);
      return viewVCE_voice.onchangeDSR(param, osc, value);
    } else if ((ret = id.match(filterPattern))) {
      param = 'FILTER'
      funcname = 'setOscFILTER';
      osc = parseInt(ret[1])
      // the synergy wants to see -1 for Af, <osc#> for Bf and 0 for no filter
      value = parseInt(valueConverter(ele.value), 10);
      args = [osc, value];
    } else if ((ret = id.match(waveKeyPattern))) {
      param = ret[1]
      osc = parseInt(ret[2], 10)
      if (param == 'WAVE') {
        funcname = 'setOscWAVE'
        value = ele.value == 'Sin' ? 0 : 1;
        args = [osc, value];
        viewVCE.vce.Envelopes[osc - 1].FreqEnvelope.Table[3] |= value;
      }
      else {
        funcname = 'setOscKEYPROP'
        value = ele.checked ? 1 : 0;
        args = [osc, value];
        viewVCE.vce.Envelopes[osc - 1].FreqEnvelope.Table[3] |=
            (value ? 0x10 : 0);
      }
    } else if ((ret = id.match(oscPattern))) {
      param = ret[1];
      osc = parseInt(ret[2], 10)
      value = parseInt(valueConverter(ele.value), 10)
      args = [osc, value]

      // console.log("changed: " + id + " param: " + param + " osc: " + osc);
      viewVCE.vce.Envelopes[osc - 1].FreqEnvelope[param] = value;
      funcname = 'setVoiceByte'

    } else if ((ret = id.match(headPattern))) {
      param = id;
      value = parseInt(valueConverter(ele.value), 10)
      args = [value]
      viewVCE.vce.Head[param] = valueConverter(ele.value);
      funcname = 'setVoiceByte'
    }
    // console.dir(viewVCE.viewVCE.vce);
    if (param != null) {
      try {
        switch (funcname) {
          case 'setVNAME':
            await UIService.SetVNAME(args);
            break;
          case 'setOscFILTER':
            await UIService.SetOscFILTER(args);
            break;
          case 'setOscWAVE':
            await UIService.SetOscWAVE(args);
            break;
          case 'setOscKEYPROP':
            await UIService.SetOscKEYPROP(args);
            break;
          case 'setVoiceByte':
            await UIService.SetVoiceByte(param, args);
            break;
        }
        viewVCE_voice.sendToCSurface(ele, ele.id, value);
        if (updater != undefined) {
          console.log('updater: ' + updater);
          updater(ele);
        }
      } catch (exc) {
        // failed - dont change the boolean
        index.errorNotification(exc);
        return false;
      }
    }
    return true;
  },

  patchTable: function() {
    // FIXME: Dynamic update of option not yet working:
    // The _first_ load of the page appear that the DOM is not ready for the
    // updated options, despite being called from the jquery load() "complete"
    // callback.   All subsequent loads(), however work.
    // Fix this once the DOM lifecycle issue is sorted out.
    //
    //$("#patchType").empty().append(viewVCE_voice.patchTypeOptions);
    document.querySelector('#patchType').value = viewVCE.vce.Extra.PatchType;

    let tbody = document.getElementById('patchTbody');
    // remove old rows:
    while (tbody.firstChild) {
      tbody.removeChild(tbody.firstChild);
    }

    let outRegisters = [[], [], [], []];
    let freqDAG = '';

    viewVCE_voice.sendToCSurface(null, `num-osc`, viewVCE.vce.Head.VOITAB + 1);
    /*
    for (let osc = viewVCE.vce.Head.VOITAB + 1; osc < 16; osc++) {
            // midi initialation for unused osc's
            viewVCE_voice.sendToCSurface(null, `OHARM[${osc + 1}]`, 0);
            viewVCE_voice.sendToCSurface(null, `FDETUN[${osc + 1}]`, 0);
            viewVCE_voice.sendToCSurface(null, `MUTE[${osc + 1}]`, 0);
            viewVCE_voice.sendToCSurface(null, `SOLO[${osc + 1}]`, 0);
            viewVCE_voice.sendToCSurface(null, `wkWAVE[${osc + 1}]`, 0);
            viewVCE_voice.sendToCSurface(null, `wkKEYPROP[${osc + 1}]`, 0);
            viewVCE_voice.sendToCSurface(null, `FILTER[${osc + 1}]`, 0);
            viewVCE_voice.sendToCSurface(null, `osc-enabled[${osc + 1}]`, 0);
    }*/

    let debug_patchBytes = '';
    // populate new ones:
    for (let osc = 0; osc <= viewVCE.vce.Head.VOITAB; osc++) {
      viewVCE_voice.sendToCSurface(null, `osc-enabled[${osc + 1}]`, 1);

      let tr = document.createElement('tr');
      let td = document.createElement('td');

      //--- OSC
      td.innerHTML = osc + 1;  // Osc
      // Mute
      let span;
      span = document.createElement('span');
      span.innerHTML =
          `&nbsp;&nbsp;<span onclick="viewVCE_voice.toggleOsc(this)" class="vceEditToggleText" aria-label="MUTE[${
              osc + 1}]" id="MUTE[${osc + 1}]">M</span>`;
      td.append(span);
      viewVCE_voice.sendToCSurface(null, `MUTE[${osc + 1}]`, 0)

      // Solo
      span = document.createElement('span');
      span.innerHTML =
          `&nbsp;<span onclick="viewVCE_voice.toggleOsc(this)" class="vceEditToggleText" aria-label="SOLO[${
              osc + 1}]" id="SOLO[${osc + 1}]">S</span>`;
      td.append(span);
      viewVCE_voice.sendToCSurface(null, `SOLO[${osc + 1}]`, 0)

      tr.appendChild(td);

      // Gain
      let gain = viewVCE_envs.computeOscGain(osc, 2);
      td = document.createElement('td');
      td.innerHTML =
          `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinPLAIN" aria-label="OscGain[${
              osc + 1}]" id="OscGain[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this,undefined,undefined)" value="${
              gain}"
			min="0" max="100"
			disabled/></div>`;
      tr.appendChild(td);
      viewVCE_voice.sendToCSurface(null, `OscGain[${osc + 1}]`, gain)

      // XREF: patch byte encode/decode
      // FIXME: assumes envelopes are sorted in oscillator order
      let patchByte = viewVCE.vce.Envelopes[osc].FreqEnvelope.OPTCH;
      let patchInhibitAddr = (patchByte & 0x20) != 0;
      let patchInhibitF0 = (patchByte & 0x04) != 0;
      let patchOutputDSR = ((patchByte & 0xc0) >> 6);
      let patchAdderInDSR = ((patchByte & 0x18) >> 3);
      let patchFOInputDSR = (patchByte & 0x03);

      //			console.log(osc + " patch byte: " + patchByte +
      //"\n" + 				" patchInhibitAddr : " +
      // patchInhibitAddr + "\n" + 				" patchInhibitF0
      // : " + patchInhibitF0 + "\n" + 				" patchOutputDSR
      // : " + patchOutputDSR + "\n" + 				"
      // patchAdderInDSR  : " + patchAdderInDSR + "\n"
      //+ 				" patchFOInputDSR  : " + patchFOInputDSR
      //+ "\n");

      debug_patchBytes += ', ' + patchByte;

      // compute the DAG based on current register usage:
      if (!patchInhibitF0) {
        let modulatingOscs = outRegisters[patchFOInputDSR];
        for (let i = 0; i < modulatingOscs.length; i++) {
          freqDAG += `[${modulatingOscs[i] + 1}]-[${osc + 1}]\n`;
        }
      } else {
        freqDAG += `[${osc + 1}]\n`;
      }
      if (patchInhibitAddr) {
        // no longer summing, this output starts a new set of addrs:
        outRegisters[patchOutputDSR] = [];
      }
      outRegisters[patchOutputDSR].push(osc);

      //--- Patch F
      td = document.createElement('td');
      let reg = 0;
      if (!patchInhibitF0) {
        reg = patchFOInputDSR + 1;
      } else {
        reg = 0;
      }
      if (osc == 0) {
        // the first osc's Freq DSR can't be altered - render as a disabled
        // input control so we can get its value in the onchange function
        // without any special casing
        td.innerHTML =
            `<div class="spinwrapper"><input type="text" class="vceNum vceEditDisabled" aria-label="patchFOInputDSR[${
                osc + 1}]" id="patchFOInputDSR[${osc + 1}]" 
				value="${
                viewVCE_voice.NullablePatchRegisterToText('' + reg)}" 
				disabled/></div>`;
      } else {
        td.innerHTML =
            `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinNullablePatchReg" aria-label="patchFOInputDSR[${
                osc + 1}]" id="patchFOInputDSR[${osc + 1}]" 
				onchange="viewVCE_voice.onchange(this,undefined,viewVCE_voice.TextToNullablePatchRegister)" value="${
                viewVCE_voice.NullablePatchRegisterToText('' + reg)}" 
				min="0" max="2"
				disabled/></div>`;
      }
      tr.appendChild(td);

      //--- Patch A
      td = document.createElement('td');
      if (!patchInhibitAddr) {
        reg = patchAdderInDSR + 1;
      } else {
        reg = 0;
      }
      td.innerHTML =
          `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinNullablePatchReg" aria-label="patchAdderInDSR[${
              osc + 1}]" id="patchAdderInDSR[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this,undefined,viewVCE_voice.TextToNullablePatchRegister)" value="${
              viewVCE_voice.NullablePatchRegisterToText('' + reg)}" 
			min="0" max="2"
			disabled/></div>`;
      tr.appendChild(td);

      //--- Patch O
      td = document.createElement('td');
      td.innerHTML =
          `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinPlain" aria-label="patchOutputDSR[${
              osc + 1}]" id="patchOutputDSR[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this,undefined,undefined)" value="${
              patchOutputDSR + 1}" 
			min="1" max="2"
			disabled/></div>`;
      tr.appendChild(td);

      //--- Hrm
      // HACK: we wrap these elements in <div> of a fixed size to keep the
      // up/down buttons positioned properly. Someone more skilled in the ways
      // of CSS would surely have a cleaner solution.
      td = document.createElement('td');
      td.innerHTML =
          `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinOHARM" aria-label="OHARM[${
              osc + 1}]" id="OHARM[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this,undefined,viewVCE_voice.TextToOHARM)" value="${
              viewVCE_voice.OHARMToText(
                  viewVCE.vce.Envelopes[osc].FreqEnvelope.OHARM)}" 
			min="-12" max="30"
			disabled/></div>`;
      tr.appendChild(td);
      viewVCE_voice.sendToCSurface(
          null, `OHARM[${osc + 1}]`,
          viewVCE.vce.Envelopes[osc].FreqEnvelope.OHARM)

      //--- Detn
      td = document.createElement('td');
      td.innerHTML =
          `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinFDETUN" aria-label="FDETUN[${
              osc + 1}]" id="FDETUN[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this,undefined,viewVCE_voice.TextToFDETUN)" value="${
              viewVCE_voice.FDETUNToText(
                  viewVCE.vce.Envelopes[osc].FreqEnvelope.FDETUN)}" 
			min="-63" max="63"
			disabled/></div>`;
      tr.appendChild(td);
      viewVCE_voice.sendToCSurface(
          null, `FDETUN[${osc + 1}]`,
          viewVCE.vce.Envelopes[osc].FreqEnvelope.FDETUN)

      let waveByte = viewVCE.vce.Envelopes[osc].FreqEnvelope.Table[3];
      let wave = ((waveByte & 0x1) == 0) ? 'Sin' : 'Tri';
      let keyprop = ((waveByte & 0x10) == 0) ? false : true;

      //--- Wave
      td = document.createElement('td');
      td.innerHTML = wave;
      td.innerHTML = `<select class="vceEdit" aria-label="wkWAVE[${
          osc + 1}]" id="wkWAVE[${osc + 1}]" value="${wave}" 
			onchange="viewVCE_voice.onchange(this)" disabled/>
			<option ${
          wave == 'Sin' ? 'selected' : ''} value="Sin">Sin</option>
			<option ${
          wave == 'Tri' ? 'selected' : ''} value="Tri">Tri</option>
			</select>
			`;
      tr.appendChild(td);
      viewVCE_voice.sendToCSurface(
          null, `wkWAVE[${osc + 1}]`, wave == 'Sin' ? 0 : 1);

      //--- Key
      td = document.createElement('td');
      // can't use disabled attr - bootstrap styling hides it - use javascript
      // hack to make it readonnly
      td.innerHTML = `<input type="checkbox" aria-label="wkKEYPROP[${
          osc + 1}]" id="wkKEYPROP[${osc + 1}]" value="true" 
			${keyprop ? ' checked ' : ''} 
			onchange="viewVCE_voice.voicingMode ? viewVCE_voice.onchange(this) : (this.checked=!this.checked)"/>`;
      tr.appendChild(td);
      viewVCE_voice.sendToCSurface(
          null, `wkKEYPROP[${osc + 1}]`, keyprop ? 1 : 0);

      //--- Flt
      td = document.createElement('td');
      let filter = viewVCE.vce.Head.FILTER[osc];
      td.innerHTML = (filter == 0) ? '' :
          (filter > 0)             ? ('Bf ' + filter) :
                                     ('Af ' + -filter);
      td = document.createElement('td');
      td.innerHTML = wave;
      td.innerHTML = `<select class="vceEdit" aria-label="FILTER[${
          osc + 1}]" id="FILTER[${osc + 1}]" value="${filter}" 
					onchange="viewVCE_voice.onchange(this,viewVCE_voice.filterChanged)" disabled/>
					<option ${
          filter == 0 ? 'selected' : ''} value="0"></option>
					<option ${
          filter < 0 ? 'selected' : ''} value="-1">Af</option>
					<option ${
          filter > 0 ? 'selected' : ''} value="${osc + 1}">Bf</option>
					</select>
					`;
      tr.appendChild(td);
      viewVCE_voice.sendToCSurface(null, `FILTER[${osc + 1}]`, filter);

      tbody.appendChild(tr);
    }
    // final row is the plus/minus buttons
    {
      let temp = document.createElement('template');
      temp.innerHTML =
          `<tr class="listplusminus" aria-label="oscPlusMinus" id="oscPlusMinus" style="display:none;">
							    <td colspan="9">
			    					<div style="margin-top: 5px; float: left;">
				    					<input aria-label="del-osc" id="del-osc" type='button' value='-'
										    onclick='viewVCE_voice.setNumOscillators(viewVCE.vce.Head.VOITAB)' />
									    <input aria-label="add-osc" id="add-osc" type='button' value='+'
									   	    onclick='viewVCE_voice.setNumOscillators(viewVCE.vce.Head.VOITAB+2)' />
								    </div>
								</td>
							</tr>`;
      tbody.appendChild(temp.content.firstChild);
    }

    console.log('patchBytes: ' + debug_patchBytes);
    console.log('freqDAG: ' + freqDAG);
    // Generate the patch diagram:
    let patchDiagramCanvas = document.getElementById('patchDiagram');
    // nomnoml is confused by leading spaces on directives lines, so...:
    let patchDiagramSource = `
#ranker: longest-path
#spacing: 12
#padding: 3
#fontSize: 10
#fill: #333
#lineWidth:1
#stroke: #fff
#background: #252525
#bendSize: 1
${freqDAG}
`;
    // console.log("nomnoml src: " + patchDiagramSource);
    nomnoml.draw(patchDiagramCanvas, patchDiagramSource);
  },

  changePatchType: async function(newIndex) {
    let patchBytes = [];
    index.spinnerOn();
    try {
      patchBytes = await UIService.SetPatchType(parseInt(newIndex, 10));
      index.spinnerOff();
    } catch (exc) {
      index.spinnerOff();
      index.errorNotification(exc);
      return
    }
    for (let i = 0; i < viewVCE.vce.Envelopes.length; i++) {
      viewVCE.vce.Envelopes[i].FreqEnvelope.OPTCH = patchBytes[i];
    }
    viewVCE.vce.Extra.PatchType = parseInt(newIndex, 0);
    viewVCE.init();
    index.refreshConnectionStatus();
  },

  setNumOscillators: function(newNum) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    viewVCE_voice.deb_setNumOscillators(newNum);
  },

  deb_setNumOscillators: null,

  raw_setNumOscillators: async function(newNum) {
    if (viewVCE.supressOnchange) {
      /*console.log("raw viewVCE.suppressOnChange");*/
      return;
    }
    console.log('setNumOscillators: ' + newNum);
    if (newNum < 1 || newNum > 16) {
      return;
    }

    // reset the floating point backing arrays: FIXME: this may mean that
    // someone who is in midst of twiddling gain, then adds/deletes a
    // oscillator, then expects gain to keep working may get a suprised
    // distortion in the envelope shape.
    viewVCE_envs.clearFloatAmpVal();

    try {
      let r = await UIService.SetNumOscillators(
          parseInt(newNum, 10),
          parseInt(document.getElementById('patchType').value, 10))
      /// now the tricky part - update the in memory version of vce to reflect
      /// what just happened:
      viewVCE.vce.Head.VOITAB = newNum - 1
      let oldLength = viewVCE.vce.Envelopes.length

      if (viewVCE.vce.Head.VOITAB <= 0) {
        document.querySelector('#del-osc').classList.add('disabled');
      }
      else {
        document.querySelector('#del-osc').classList.remove('disabled');
      }
      if (viewVCE.vce.Head.VOITAB >= 15) {
        document.querySelector('#add-osc').classList.add('disabled');
      } else {
        document.querySelector('#add-osc').classList.remove('disabled');
      }
      if (newNum <= oldLength) {
        // nothing to do - just ignored the extra envelopes
      } else {
        for (let i = oldLength; i < newNum; i++) {
          // copy the envelope template into the vce:
          // abuse JSON to do a deep copy:
          viewVCE.vce.Envelopes[i] =
              JSON.parse(JSON.stringify(r.EnvelopeTemplate));
          // overwrite the default patch type
          console.log('before copy', viewVCE.vce.Envelopes);
          viewVCE.vce.Envelopes[i].OPTCH = r.PatchBytes[i];
          console.log(i, 'copy env - now ', viewVCE.vce.Envelopes[i]);
          console.log('AFTER copy', viewVCE.vce.Envelopes);
        }
      }
      viewVCE.init();
    } catch (exc) {
      // failed - dont change the boolean
      index.errorNotification(exc);
    }
    index.refreshConnectionStatus();
  },

  _withZeroconf: async function(
      prompt1, prompt2,
      zeroconfSelector /* "getSynergy" or "getSynergyAndControlSurface" */,
      actionAfterSelect, recurseAction, successCallback) {
    // fetch the config from the server (server will return already selected
    // config, or "not enabled" or a list of selections) if server returns
    // "already configured" or "not enabled"
    //      send request to start voicemode to server
    // if server returns a list of options
    //     popup a selection dialog
    //        Cancel event aborts attempt to start voicemode
    //        OK event sends the selection and request to start voicemode to
    //        server Rescan event sends request to rescan to the server -- then
    //        recursively calls config

    console.log(
        'top withZeroconf ' + zeroconfSelector + ' ' + actionAfterSelect);

    let r;
    try {
      switch (zeroconfSelector) {
        case 'getSynergy':
          r = await UIService.GetSynergy();
          break;
        case 'getSynergyAndControlSurface':
          r = await UIService.GetSynergyAndControlSurface();
          break;
        default:
          index.errorNotification(
              'invalid arg to _withZeroconf - ' + zeroconfSelector);
          break;
      }
    } catch (exc) {
      console.log(
          '_withZeroconfig - sendMessage failed: ' + JSON.stringify(exc))
      index.errorNotification(exc);
    }
    console.log(zeroconfSelector + ' returned ' + JSON.stringify(r));
    if (((!r[0].HasDevice) || r[0].AlreadyConfigured) &&
        ((!r[1].HasDevice) || r[1].AlreadyConfigured)) {
      console.log('call actionAfterSelect');
      actionAfterSelect(null, null, successCallback);
    } else {
      // zeroconf found more than one option - show dialog
      console.log('show menu');
      if ((!r[0].HasDevice) || r[0].AlreadyConfigured) {
        prompt1 = null;
      }
      if ((!r[1].HasDevice) || r[1].AlreadyConfigured) {
        prompt2 = null;
      }
      index.chooseZeroconfService(
          prompt1, r[0].Choices, prompt2, r[1].Choices,
          function() {
            console.log('cancelled');
            // cancelled - do nothing
          },
          function(choice1, choice2) {
            console.log(
                'user chose ' + JSON.stringify(choice1) + ' ' +
                JSON.stringify(choice2));
            // user selected one of the options
            actionAfterSelect(choice1, choice2, successCallback);
          },
          async function() {
            console.log('rescan');
            // user asked for a rescan
            index.spinnerOn();
            try {
              await UIService.RescanZeroconf();
              console.log('rescan done');
              // recurse
              index.spinnerOff();
              recurseAction();
            } catch (exc) {
              // failed - abort
              index.spinnerOff();
              index.errorNotification(exc);
            }
          });
    }
  },

  connectSynergy: function(successCallback) {
    let wasDisconnectedSynergy = index.synergyName === null;
    return viewVCE_voice._withZeroconf(
        'Choose Synergy', null, 'getSynergy', viewVCE_voice.raw_connectSynergy,
        viewVCE_voice.connectSynergy, function() {
          if (wasDisconnectedSynergy) {
            index.checkVersion(wasDisconnectedSynergy, false /*don't care*/);
          }
          successCallback();
        });
  },

  voicingModeOn: function() {
    let wasDisconnectedSynergy = index.synergyName === null;
    let wasDisconnectedCs = index.controlSurfaceName === null;
    return viewVCE_voice._withZeroconf(
        'Choose Synergy', 'Choose Control Surface',
        'getSynergyAndControlSurface', viewVCE_voice.raw_voicingModeOn,
        viewVCE_voice.voicingModeOn, function() {
          if (wasDisconnectedSynergy || wasDisconnectedCs) {
            index.checkVersion(wasDisconnectedSynergy, wasDisconnectedCs);
          }
        });
  },

  raw_connectSynergy: async function(zeroconfChoice, ignored, callback) {
    index.spinnerOn();
    let r;
    try {
      r = await UIService.ConnectSynergy(zeroconfChoice);
      index.spinnerOff();
    } catch (exc) {
      index.spinnerOff();
      index.errorNotification(exc);
      return
    }
    index.updateConnectionStatus(
        r.Status.SynergyName, r.Status.ControlSurfaceName);
    if (!r.AlreadyConnected) {
      index.infoNotification(
          'Successfully connected to Synergy: ' + r.Status.SynergyName);
    }
    callback();
    return
  },

  toggleVoicingMode: function(mode) {
    if (!mode) {
      index.confirmDialog(
          'Disabling Voicing Mode will discard any pending edits. Are you sure?',
          function() {
            viewVCE_voice.raw_voicingModeOff(false);
          });
    } else {
      viewVCE_voice.voicingModeOn();
    }
  },

  raw_voicingModeOff: async function(disconnect) {
    console.log(`VoicingMode off`);

    index.spinnerOn();
    try {
      let r = await UIService.ToggleVoicingMode(
          false /*mode*/, disconnect /*disconnect*/, null /*vce*/,
          null /*zeroconfSynergy*/, null /*zeroconfCs*/);
      index.spinnerOff();
      viewVCE_voice.voicingMode = false;
      if (r != null) {
        // FIXME: why are we ignoring the values returned by the service??

        // if we just disabled voicing, clear the VCE view
        document.getElementById('content').innerHTML = '';
        document.querySelector('#disableControlSurfaceMenuItem')
            .classList.add('disabled');
        viewVCE_voice.csEnabled = false;
      }
      viewVCE.setVCE(null);
      viewVCE_voice.voicingModeVisuals();
      let msg = 'Voicing mode disabled.';
      if (disconnect) {
        msg = msg + ' Synergy Disconnected.';
      }
      index.infoNotification(msg);
    } catch (exc) {
      index.spinnerOff();
      console.log(exc);
      index.errorNotification(exc);
    }
    index.refreshConnectionStatus();
  },

  raw_voicingModeOn: async function(
      synergyZeroconfChoice, csZeroconfChoice, callback) {
    console.log(`VoicingMode on`);

    index.spinnerOn();
    try {
      let r = await UIService.ToggleVoicingMode(
          true /*mode*/, false /*disconnect*/, viewVCE.vcr /*vce*/,
          synergyZeroconfChoice /*zeroconfSynergy*/,
          csZeroconfChoice /*zeroconfCs*/);
      index.spinnerOff();
      console.log('toggleVoiceMode returned: ', r);
      //  Check error

      viewVCE_voice.voicingMode = true;
      let csMessage = '';
      if (r != null) {
        let loaded_vce = r[0];
        let csEnabled = r[1];
        let csName = r[2];
        // let synergyName_ignored = r[3];

        viewVCE.setVCE(loaded_vce);
        viewVCE_voice.csEnabled = csEnabled;

        if (viewVCE_voice.csEnabled) {
          csMessage = `.<br>Control Surface is enabled: ${csName}.`;
        } else {
          csMessage = `.<br>Control Surface is not enabled.`;
        }
        viewCRT.setCRT(null, null);
        index.load('viewVCE.html', 'content', function(ele) {
          viewVCE.init();
        });
      }
      index.infoNotification(`Voicing mode ${
          viewVCE_voice.voicingMode ? 'enabled' : 'disabled'}.${csMessage}`);

      index.refreshConnectionStatus();

      // reset the SOLO arrays
      for (let osc = 0; osc < 16; osc++) {
        viewVCE_voice.MUTE[osc] = false;
        viewVCE_voice.SOLO[osc] = false;
        document.querySelectorAll('.vceEditToggle').forEach(el => {
          el.classList.remove('on');
        });
        viewVCE_voice.sendToCSurface(null, `MUTE[${osc + 1}]`, 0)
        viewVCE_voice.sendToCSurface(null, `MUTE[${osc + 1}]`, 0)
      }
      viewVCE_voice.voicingModeVisuals();
      callback();

    } catch (exc) {
      index.spinnerOff();
      index.errorNotification(exc);
      viewVCE_voice.csEnabled = false;
    }
  },

  touchspin_init: function(els, callback_before, callback_after) {
    els.forEach(el => {
      let spinner;
      if (callback_before === undefined) {
        spinner = TouchSpin(el, {
          renderer: VanillaRenderer,
          verticalbuttons: true,
          verticalup: '\u25b4',     //'\u25b2',
          verticaldown: '\u25be',   //'\u25bc',
          buttonup_txt: '\u25b4',   //'\u25b2',
          buttondown_txt: '\u25be'  //'\u25bc',
        });
      } else {
        spinner = TouchSpin(el, {
          renderer: VanillaRenderer,
          verticalbuttons: true,
          verticalup: '\u25b4',      //'\u25b2',
          verticaldown: '\u25be',    //'\u25bc',
          buttonup_txt: '\u25b4',    //'\u25b2',
          buttondown_txt: '\u25be',  //'\u25bc',
          callback_before_calculation: function(value) {
            return callback_before(value);
          },
          callback_after_calculation: function(value) {
            return callback_after(value);
          }
        });
      }
      // spinner.addEventListener('keydown', (e) => {
      //   switch (e.key) {
      //     case 'ArrowUp':
      //       e.preventDefault();
      //       spinner.upOnce();
      //       break;
      //     case 'ArrowDown':
      //       e.preventDefault();
      //       spinner.downOnce();
      //       break;
      //   }
      // });
      //   el.addEventListener('change:start', (event) => {
      //    event.preventDefault();
      //  })
    });
  },

  voicingModeVisuals: function() {
    let mode = viewVCE_voice.voicingMode;
    // mode = true; // For debugging and CSS tweaking: force edit controls to be
    // visible

    // XREF: converter mappings : these are duplicated in updateFromCSurface
    if (mode) {
      // CSS for styling the buttons when disabled is HARD.  So avoid it.
      viewVCE_voice.touchspin_init(
          document.querySelectorAll('.vceNum.spinNullablePatchReg'),
          viewVCE_voice.TextToNullablePatchRegister,
          viewVCE_voice.NullablePatchRegisterToText);

      viewVCE_voice.touchspin_init(
          document.querySelectorAll('.vceNum.spinOHARM'),
          viewVCE_voice.TextToOHARM, viewVCE_voice.OHARMToText);
      viewVCE_voice.touchspin_init(
          document.querySelectorAll('.vceNum.spinFDETUN'),
          viewVCE_voice.TextToFDETUN, viewVCE_voice.FDETUNToText);
      viewVCE_voice.touchspin_init(
          document.querySelectorAll('.vceNum.spinAmpEnv'),
          viewVCE_envs.TextToAmpEnvValue, viewVCE_envs.AmpEnvValueToText);
      viewVCE_voice.touchspin_init(
          document.querySelectorAll('.vceNum.spinAmpTime'),
          viewVCE_envs.TextToAmpTimeValue, viewVCE_envs.AmpTimeValueToText);
      viewVCE_voice.touchspin_init(
          document.querySelectorAll('.vceNum.spinFreqTime'),
          viewVCE_envs.TextToFreqTimeValue, viewVCE_envs.FreqTimeValueToText);
      // plain number variant:
      viewVCE_voice.touchspin_init(
          document.querySelectorAll('.vceNum.spinPLAIN'), undefined, undefined);
      // make any plain-text spans align:
      document.querySelectorAll('.spinNOSPIN')
          .forEach(el => {el.classList.add('spinNOSPIN-Enabled')});
    }

    // Load/Save menu items get disabled/enabled:
    if (mode) {
      document.querySelector('#disableVRAMMenuItem').classList.add('disabled');
      document.querySelector('#loadCRTMenuItem').classList.add('disabled');
      document.querySelector('#saveVCEMenuItem').classList.remove('disabled');
      document.querySelector('#oscPlusMinus').style.display = 'block';
    } else {
      document.querySelector('#disableVRAMMenuItem')
          .classList.remove('disabled');
      document.querySelector('#loadCRTMenuItem').classList.remove('disabled');
      document.querySelector('#saveVCEMenuItem').classList.add('disabled');
      document.querySelector('#oscPlusMinus').style.display = 'none';
    }
    if (viewVCE.vce && viewVCE.vce.Head.VOITAB <= 0) {
      document.querySelector('#del-osc').classList.add('disabled');
    } else {
      document.querySelector('#del-osc').classList.remove('disabled');
    }
    if (viewVCE.vce && viewVCE.vce.Head.VOITAB >= 15) {
      document.querySelector('#add-osc').classList.add('disabled');
    } else {
      document.querySelector('#add-osc').classList.remove('disabled');
    }

    document.querySelectorAll('.vceEdit').forEach(el => {
      el.disabled = !mode;
    });
    if (mode) {
      document.querySelectorAll('.vceEditToggleText').forEach(el => {
        el.style.display = 'block';
      });
    } else {
      document.querySelectorAll('.vceEditToggleText').forEach(el => {
        el.style.display = 'none';
      });
    }
    document.getElementById('voiceModeButtonImg').src =
        `static/images/red-button-${
            viewVCE_voice.voicingMode ? 'on' : 'off'}-full.png`;
  },

  chart: null,

  updateChart: function() {
    let ampData = viewVCE_voice.ampProportionCurve(
        viewVCE.vce.Head.VACENT, viewVCE.vce.Head.VASENS);
    let timbreData = viewVCE_voice.timbreProportionCurve(
        viewVCE.vce.Head.VTCENT, viewVCE.vce.Head.VTSENS);

    viewVCE_voice.chart.data.datasets[0].data = ampData;
    viewVCE_voice.chart.data.datasets[1].data = timbreData;
    viewVCE_voice.chart.update();
  },

  updateVibType: function() {
    document.getElementById('vibType').innerHTML =
        (viewVCE.vce.Head.VIBDEP >= 0) ? 'Sine' : 'Random';
  },

  patchTypeOptions: null,  // HTML fragment used to populate the <options> for
                           // the patchType select

  getPatchTypeNames: async function() {
    let nameArray;
    try {
      nameArray = await UIService.GetPatchTypeNames()
    } catch (exc) {
      index.errorNotification(exc);
    }
    let html = ''
    for (let i = 0; i < nameArray.length; i++) {
      html += `<option value="${i}">${nameArray[i]}</option>\n`;
    }
    viewVCE_voice.patchTypeOptions = html;
  },

  init: function(incrementalUpdate) {
    console.log('--- start viewVCE_voice init');

    if (viewVCE_voice.deb_onchange == null) {
      viewVCE_voice.deb_onchange = index.debounceFirstArg(
          viewVCE_voice.raw_onchange, index.DEBOUNCE_WAIT);
    }
    if (viewVCE_voice.deb_setNumOscillators == null) {
      viewVCE_voice.deb_setNumOscillators =
          _.debounce(viewVCE_voice.raw_setNumOscillators, index.DEBOUNCE_WAIT);
    }
    if (viewVCE_voice.patchTypeNames == null) {
      viewVCE_voice.getPatchTypeNames();
    }

    document.querySelector('#vceTabs a[href="#vceVoiceTab"]')
        .addEventListener('shown.bs.tab', () => {
          viewVCE_voice.sendToCSurface(null, 'voice-tab', 1);
        });
    document.querySelector('#vceTabs a[href="#vceEnvsTab"]')
        .addEventListener('shown.bs.tab', () => {
          viewVCE_voice.sendToCSurface(null, 'freq-envelopes-tab', 1);
        });
    document.querySelector('#vceTabs a[href="#vceFiltersTab"]')
        .addEventListener('shown.bs.tab', () => {
          viewVCE_voice.sendToCSurface(null, 'filters-tab', 1);
        });
    document.querySelector('#vceTabs a[href="#vceKeyEqTab"]')
        .addEventListener('shown.bs.tab', () => {
          viewVCE_voice.sendToCSurface(null, 'keyeq-tab', 1);
        });
    document.querySelector('#vceTabs a[href="#vceKeyPropTab"]')
        .addEventListener('shown.bs.tab', () => {
          viewVCE_voice.sendToCSurface(null, 'keyprop-tab', 1);
        });

    viewVCE_voice.patchTable();

    console.log(
        'view VCE, CRT:' + viewCRT.crt_name +
        ', VCE: ' + viewVCE.vce.Head.VNAME)

    if (viewCRT.crt_name == null) {
      document.getElementById('backToCRT').hidden = true;
    }
    else {
      document.getElementById('backToCRT').hidden = false;
    }

    document.getElementById('nOsc').innerHTML = viewVCE.vce.Head.VOITAB + 1;
    document.getElementById('keysPlayable').innerHTML =
        Math.floor(32 / (viewVCE.vce.Head.VOITAB + 1));
    viewVCE_voice.updateVibType();
    document.getElementById('VIBRAT').value = viewVCE.vce.Head.VIBRAT;
    viewVCE_voice.sendToCSurface(
        document.getElementById('VIBRAT'), 'VIBRAT', viewVCE.vce.Head.VIBRAT)
    document.getElementById('VIBDEL').value = viewVCE.vce.Head.VIBDEL;
    viewVCE_voice.sendToCSurface(
        document.getElementById('VIBDEL'), 'VIBDEL', viewVCE.vce.Head.VIBDEL)
    document.getElementById('VIBDEP').value = viewVCE.vce.Head.VIBDEP;
    viewVCE_voice.sendToCSurface(
        document.getElementById('VIBDEP'), 'VIBDEP', viewVCE.vce.Head.VIBDEP)
    document.getElementById('APVIB').value = viewVCE.vce.Head.APVIB;
    viewVCE_voice.sendToCSurface(
        document.getElementById('APVIB'), 'APVIB', viewVCE.vce.Head.APVIB)

    document.getElementById('VTRANS').value = viewVCE.vce.Head.VTRANS;
    viewVCE_voice.sendToCSurface(
        document.getElementById('VTRANS'), 'VTRANS', viewVCE.vce.Head.VTRANS)
    document.getElementById('VACENT').value = viewVCE.vce.Head.VACENT;
    viewVCE_voice.sendToCSurface(
        document.getElementById('VACENT'), 'VACENT', viewVCE.vce.Head.VACENT)
    document.getElementById('VASENS').value = viewVCE.vce.Head.VASENS;
    viewVCE_voice.sendToCSurface(
        document.getElementById('VASENS'), 'VASENS', viewVCE.vce.Head.VASENS)
    document.getElementById('VTCENT').value = viewVCE.vce.Head.VTCENT;
    viewVCE_voice.sendToCSurface(
        document.getElementById('VTCENT'), 'VTCENT', viewVCE.vce.Head.VTCENT)
    document.getElementById('VTSENS').value = viewVCE.vce.Head.VTSENS;
    viewVCE_voice.sendToCSurface(
        document.getElementById('VTSENS'), 'VTSENS', viewVCE.vce.Head.VTSENS)

    let count = 0;
    for (let i = 0; i < viewVCE.vce.Head.FILTER.length; i++) {
      if (viewVCE.vce.Head.FILTER[i] != 0) {
        count++;
      }
    }
    document.getElementById('nFilter').innerHTML = count;

    Chart.defaults.global.defaultFontColor = 'white';
    Chart.defaults.global.defaultFontSize = 14;

    let ampData = viewVCE_voice.ampProportionCurve(
        viewVCE.vce.Head.VACENT, viewVCE.vce.Head.VASENS);
    let timbreData = viewVCE_voice.timbreProportionCurve(
        viewVCE.vce.Head.VTCENT, viewVCE.vce.Head.VTSENS);

    let ctx = document.getElementById('velocityChart').getContext('2d');
    if (viewVCE_voice.chart != null) {
      // kill off the old chart so we dont get conflicts
      viewVCE_voice.chart.destroy();
    }
    viewVCE_voice.chart = new Chart(ctx, {

      type: 'line',
      data: {
        labels: [
          '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '',
          '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''
        ],
        // labels: ['','','',''],
        datasets: [
          {
            fill: false,
            lineTension: 0,
            pointRadius: 0,
            pointHitRadius: 5,
            label: 'Amplitude',
            backgroundColor: viewVCE.chartColors[0],
            borderColor: viewVCE.chartColors[0],
            data: ampData
          },
          {
            fill: false,
            lineTension: 0,
            pointRadius: 0,
            pointHitRadius: 5,
            label: 'Timbre',
            backgroundColor: viewVCE.chartColors[1],
            borderColor: viewVCE.chartColors[1],
            data: timbreData
          }
        ]
      },

      // Configuration options go here
      options: {
        animation: {duration: 0},
        scales: {
          xAxes: [{
            gridlines: {
              color: '#666',
              display: true,
              drawBorder: true,
              drawOnChartArea: false
            },
            scaleLabel: {display: true, labelString: 'Velocity'},
            ticks: {display: true}
          }],
          yAxes: [{
            gridlines: {
              color: '#666',
              display: true,
              drawBorder: true,
              drawOnChartArea: false
            },
            scaleLabel: {display: true, labelString: 'Proportion'},
            ticks: {
              display: true,
              max: 33,
              min: 0,
              callback: function(dataLabel, index) {
                return '';
              }
            }
          }],
        },
        responsive: false,
        maintainAspectRatio: false
      }
    });

    document.getElementById('vce_crt_name').innerHTML = viewCRT.crt_name;
    // do this last to help the uitest to not start testing too soon
//    let name = ''
//    for (let i = 0; i < viewVCE.vce.Head.VNAME.length; i++) {
//      name = name + String.fromCharCode(viewVCE.vce.Head.VNAME[i]);
//      console.log("code", viewVCE.vce.Head.VNAME[i], "name: ", name);
//    }
    let name = viewVCE.vce.Head.VNAME;
    console.log("name: ", name);
    name = name.replace(/ +$/g, ''); // trim trailing spaces for editing
    console.log("name: ", name);
    document.getElementById('vce_name').innerHTML = name;
    document.getElementById('VNAME').value = name;
    console.log('--- finish viewVCE_voice init');
  },

  updateFromCSurface: function(payload) {
    if (!viewVCE_voice.voicingMode) {
      // ignore unless we're voicing
      return
    }

    // special handling for switching tabs -- fieldnames end with "-tab":
    if (payload.Field.search('-tab') >= 0) {
      if (payload.Field === 'voice-tab' ||
          payload.Field === 'voice-freqs-tab' ||
          payload.Field === 'osc-gain-tab') {
        // let el = document.querySelector('#vceTabs a[href="#vceVoiceTab"]');
        // let tab = bootstrap.Tab.getInstance(el);
        // tab.show();
        let tab = new bootstrap.Tab(
            document.getSelection('#vceTabs a[href="#vceVoiceTab"]'));
        console.log('tab', tab);
        tab.show();
      } else if (payload.Field === 'freq-envelopes-tab') {
        let tab = new bootstrap.Tab(
            document.getSelection('#vceTabs a[href="#vceEnvsTab"]'));
        console.log('tab', tab);
        tab.show();
      } else if (payload.Field === 'amp-envelopes-tab') {
        let tab = new bootstrap.Tab(
            document.getSelection('#vceTabs a[href="#vceEnvsTab"]'));
        console.log('tab', tab);
        tab.show();
      } else if (payload.Field === 'filters-tab') {
        let tab = new bootstrap.Tab(
            document.getSelection('#vceTabs a[href="#vceFiltersTab"]'));
        console.log('tab', tab);
        tab.show();
      } else if (payload.Field === 'keyeq-tab') {
        let tab = new bootstrap.Tab(
            document.getSelection('#vceTabs a[href="#vceKeyEqTab"]'));
        console.log('tab', tab);
        tab.show();
      } else if (payload.Field === 'keyprop-tab') {
        let tab = new bootstrap.Tab(
            document.getSelection('#vceTabs a[href="#vceKeyPropTab"]'));
        console.log('tab', tab);
        tab.show();
      }
      return;
    }

    let ele = document.getElementById(payload.Field)

    // special handling for the +/- buttons
    if (payload.Field.search('add-') == 0 ||
        payload.Field.search('del-') == 0) {
      ele.onclick();
      valueString = payload.Field.search('add-') == 0 ? 'Add' : 'Delete';
      return valueString;  // don't trigger non-existent onchange()
    }

    if (ele === undefined || ele === null) {
      console.log('updateFromCSurface ' + payload.Field + ' element not found');
      return
    }
    let value = payload.Value


    // when using MIDI
    // value comes in unscaled (it's a 0-based MIDI value).
    // Use the min value on the input control to correct for an offset and then
    // use the text conversion function attached the the touchspin (if any) to
    // turn that into a string)
    //
    // when using OSC, the value is the direct Synergy byte value - no offset
    /*
    if (ele.hasAttribute("min")) {
            let min = parseInt(ele.getAttribute("min"), 10);
            value = value + min;
    }
    */

    // XREF: converter mappings: it would be nicer to directly query the input
    // element to determine what sort of touchspin callbacks are associate, if
    // any.  But its not obvous how to do that, so I duplicate some logic here
    let converter = function(value) {
      return value;
    };
    if (ele.classList.contains('spinFreqTime')) {
      converter = viewVCE_envs.FreqTimeValueToText;
    } else if (ele.classList.contains('spinNullablePatchReg')) {
      converter = viewVCE_voice.NullablePatchRegisterToText;
    } else if (ele.classList.contains('spinOHARM')) {
      converter = viewVCE_voice.OHARMToText;
    } else if (ele.classList.contains('spinFDETUN')) {
      converter = viewVCE_voice.FDETUNToText;
    } else if (ele.classList.contains('spinAmpEnv')) {
      converter = viewVCE_envs.AmpEnvValueToText;
    } else if (ele.classList.contains('spinAmpTime')) {
      converter = viewVCE_envs.AmpTimeValueToText;
    }

    let valueString = converter('' + value)

    // console.log("  updateFromCSurface " + payload.Field + "was " +
    // ele.value);
    if (ele.disabled) {
      // console.log("   disabled!");
      return;
    }
    if (ele.display == 'none') {
      // console.log("   hidden/disabled!");
      return;
    }
    if (ele.nodeName == 'SELECT') {
      // cycle through the options in each click
      let options = ele.options
      //			console.log("cycle SELECT: currently " +
      // options.selectedIndex + " len: " + options.length);
      let i = options.selectedIndex + 1
      if (i >= options.length) {
        i = 0
      }
      options.selectedIndex = i;
      //			console.log("cycle SELECT: now " +
      // options.selectedIndex);
      valueString = options[i].text;
    } else if (ele.nodeName == 'SPAN') {
      // SOLO/MUTE buttons
      ele.onclick();

      valueString = ele.classList.contains('on') ? 'ON' : 'OFF'
      return valueString;  // don't trigger non-existent onchange()
    } else if (ele.type == 'checkbox') {
      ele.checked = value == 0 ? '' : 'checked';
    } else if (ele.type == 'text') {
      ele.value = valueString;
    }
    // console.log("  updateFromCSurface " + payload.Field + "NOW " +
    // ele.value);
    ele.onchange();
    return valueString;
  },

  sendToCSurface: async function(ele, field, value) {
    if (!viewVCE_voice.csEnabled) {
      return
    }
    // pass ele==null to force the value to just be sent without scaling

    // when using MIDI
    // value comes in scaled to "synergy byte" value - but MIDI values are
    // always 0 based. when using OSC, no need to scale the offset Use the min
    // value on the input control to correct for an offset
    /*
    if (ele != null && ele.hasAttribute("min")) {
            let min = parseInt(ele.getAttribute("min"), 10);
            value = value - min;
    }
    */

    try {
      await UIService.SendToCSurface(field, parseInt(value, 10));
    } catch (exc) {
      index.errorNotification(exc);
    }
  }
};
window.viewVCE_voice = viewVCE_voice;
