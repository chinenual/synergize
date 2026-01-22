// const {lookupService} = require("dns");
// const {env} = require("process");
// const {DH_CHECK_P_NOT_PRIME} = require("constants");
import {UIService} from '/bindings/github.com/chinenual/synergize';
import {Chart} from 'chart.js';
import {_} from 'lodash';

import {index} from './index';
import {viewVCE} from './viewVCE';
import {viewVCE_voice} from './viewVCE_voice';

let dragOldValue = {x: undefined, y: undefined};

export let viewVCE_envs = {

  chart: null,

  // amp values for the currently displayed envelope - computed lazily
  floatAmpVal: null,

  clearFloatAmpVal: function() {
    this.floatAmpVal = null;
  },

  init: function(incrementalUpdate) {
    // console.log('--- start viewVCE_envs init');
    viewVCE_envs.unsetFloatVals();

    if (viewVCE_envs.deb_onchange == null) {
      viewVCE_envs.deb_onchange = index.debounceFirstArg(
          viewVCE_envs.raw_onchange, index.DEBOUNCE_WAIT);
    }
    if (viewVCE_envs.deb_onchangeGain == null) {
      viewVCE_envs.deb_onchangeGain = index.debounceFirstArg(
          viewVCE_envs.raw_onchangeGain, index.DEBOUNCE_WAIT);
    }
    if (viewVCE_envs.deb_onchangeEnvAccel == null) {
      viewVCE_envs.deb_onchangeEnvAccel =
          _.debounce(viewVCE_envs.raw_onchangeEnvAccel, index.DEBOUNCE_WAIT);
    }
    if (viewVCE_envs.deb_copyFrom == null) {
      viewVCE_envs.deb_copyFrom =
          _.debounce(viewVCE_envs.raw_copyFrom, index.DEBOUNCE_WAIT);
    }

    let selectEle = document.getElementById('envOscSelect');
    // remove old options:
    while (selectEle.firstChild) {
      selectEle.removeChild(selectEle.firstChild);
    }

    for (let i = 0; i <= viewVCE.vce.Head.VOITAB; i++) {
      let option = document.createElement('option');
      option.value = '' + (i + 1);
      option.innerHTML = '' + (i + 1);
      selectEle.appendChild(option);
    }
    document.querySelector('#envCopySelectDiv').style.display = 'none';

    viewVCE_envs.envChartUpdate(1, -1, true)
    // console.log('--- finish viewVCE_envs init');
  },


  // Freq values:
  //   as displayed: -127 .. 127 TODO
  //   byte range:   0x00 .. 0xff  -- -127 .. 127
  scaleFreqEnvValue: function(v) {
    // See OSCDSP.Z80 DISVAL: DVAL10:
    // bytes are store as unsigned values in the file, but displayed as
    // "signed" int8's -- do the 2's complement change here
    if (v > 127) {
      v = v - 256;
    }
    return v;
  },

  unscaleFreqEnvValue: function(v) {
    // bytes are store as unsigned values in the file, but displayed as
    // "signed" int8's -- do the 2's complement change here
    if (v < 0) {
      v = 256 + v;
    }
    return v;
  },

  FreqEnvValueToText(v) {
    return '' + viewVCE_envs.scaleFreqEnvValue(v);
  },

  TextToFreqEnvValue(v) {
    if (v == null || v === '') {
      return 0;
    }
    let val = parseInt(v, 10);
    return '' + viewVCE_envs.unscaleFreqEnvValue(val);
  },

  // Amp values:
  //   as displayed: 0 .. 72
  //   byte range:   0x37 .. 0x7f (55 .. 127)
  scaleAmpEnvValue: function(v) {
    // See OSCDSP.Z80 DISVAL: DVAL30:
    // if (last) return 0;
    return Math.max(0, v - 55);
  },

  unscaleAmpEnvValue: function(v) {
    return v + 55;
  },

  AmpEnvValueToText(v) {
    // console.log("AmpEnvValueToText '" + v + "' -->
    // "+viewVCE_envs.scaleAmpEnvValue(v));
    return '' + viewVCE_envs.scaleAmpEnvValue(v);
  },

  TextToAmpEnvValue(v) {
    if (v == null || v === '') {
      // console.log("TextToAmpEnvValue '" + v + "' --> 55");
      return 55;
    }
    let val = parseInt(v, 10);
    // console.log("TextToAmpEnvValue '" + v + "' ->
    // "+viewVCE_envs.unscaleAmpEnvValue(val));
    return '' + viewVCE_envs.unscaleAmpEnvValue(val);
  },

  // NOTE: the ftab based scaling functions as done in SYNHCS are not
  // reversable (for the freq case, several values of x map to the same y,
  // so reversing y can never map to some x's.  In SYNHCS, this didnt matter
  // since the mapping was one-way (the raw x values go to the synergy, the
  // y values were only used to show the values to the user. For us, we need
  // to convert the "user y" values to the "x" values to send to the
  // synergy.)
  //
  // Problem is the exponential nature of the scaling would lead to HUGE
  // numbers.  (my guess is that few real patches use these large time
  // values.) In any case, I've chosen to just create a table based mapping
  // approach rather than do all the math that SYNHCS does. This allows me
  // to substitute some unique values at the end of the range to keep things
  // unique, but not allow them to get outrageously large.
  //
  // For most values, it's exactly the same as SYNHCS, but for those extra
  // values of x, there are new y's so the editor can do its job

  freqTimeScale: [
    0,     1,     2,     3,     4,     5,     6,     7,     8,     9,     10,
    11,    12,    13,    14,    15,    25,    28,    32,    36,    40,    45,
    51,    57,    64,    72,    81,    91,    102,   115,   129,   145,   163,
    183,   205,   230,   258,   290,   326,   366,   411,   461,   517,   581,
    652,   732,   822,   922,   1035,  1162,  1304,  1464,  1644,  1845,  2071,
    2325,  2609,  2929,  3288,  3691,  4143,  4650,  5219,  5859,  6576,  7382,
    8286,  9300,  10439, 11718, 13153, 14764, 16572, 18600, 20078, 23436, 26306,
    29528, 29529, 29530, 29531, 29532, 29533, 29534, 29535
  ],

  ampTimeScale: [
    0,    1,    2,    3,    4,    5,    6,    7,    8,    9,    10,
    11,   12,   13,   14,   15,   16,   17,   18,   19,   20,   21,
    22,   23,   24,   25,   26,   27,   28,   29,   30,   31,   32,
    33,   34,   35,   36,   37,   38,   39,   40,   45,   51,   57,
    64,   72,   81,   91,   102,  115,  129,  145,  163,  183,  205,
    230,  258,  290,  326,  366,  411,  461,  517,  581,  652,  732,
    822,  922,  1035, 1162, 1304, 1464, 1644, 1845, 2071, 2325, 2609,
    2929, 3288, 3691, 4143, 4650, 5219, 5859, 6576
  ],


  // Freq Time values:
  //   as displayed: 0 .. 29528
  //   byte range:   0x0 .. 0x54 (0 .. 84)
  scaleFreqTimeValue: function(v) {
    // See OSCDSP.Z80 DISVAL for the original ftab-baased scaling which is
    // roughly:
    //	if (v <= 15) return v;
    //  return viewVCE_envs.scaleViaRtab((2 * v) - 14);
    if (v < 0) {
      return 0;
    } else if (v >= viewVCE_envs.freqTimeScale.length) {
      return viewVCE_envs.freqTimeScale[viewVCE_envs.freqTimeScale.length - 1];
    }
    return viewVCE_envs.freqTimeScale[v];
  },

  unscaleFreqTimeValue: function(v) {
    // fixme: linear search is brute force - but the list is short -
    // performance is "ok" as is...
    for (let i = 0; i < viewVCE_envs.freqTimeScale.length; i++) {
      if (viewVCE_envs.freqTimeScale[i] >= v) {
        return i;
      }
    }
    // shouldnt happen!
    return viewVCE_envs.freqTimeScale.length - 1;
  },

  FreqTimeValueToText(v) {
    return '' + viewVCE_envs.scaleFreqTimeValue(v);
  },

  TextToFreqTimeValue(v) {
    if (v == null || v === '') {
      return 0;
    }
    let val = parseInt(v, 10);
    return '' + viewVCE_envs.unscaleFreqTimeValue(val);
  },

  // Freq Time values:
  //   as displayed: 0 .. 6576
  //   byte range:   0x0 .. 0x54 (0 .. 84)
  scaleAmpTimeValue: function(v) {
    // See OSCDSP.Z80 DISVAL: DVAL20: which is:
    // if (v < 39) return v;
    // return viewVCE_envs.scaleViaRtab((v * 2) - 54);
    // console.log("scale amp time value: " + v + ", -> " +
    // viewVCE_envs.ampTimeScale[v])
    if (v < 0) {
      return 0;
    } else if (v >= viewVCE_envs.ampTimeScale.length) {
      return viewVCE_envs.ampTimeScale[viewVCE_envs.ampTimeScale.length - 1];
    }
    return viewVCE_envs.ampTimeScale[v];
  },

  unscaleAmpTimeValue: function(v) {
    // fixme: linear search is brute force - but the list is short -
    // performance is "ok" as is...
    for (let i = 0; i < viewVCE_envs.ampTimeScale.length; i++) {
      if (viewVCE_envs.ampTimeScale[i] >= v) {
        // console.log("unscale amp time value: " + v + ", -> " + i + " (" +
        // viewVCE_envs.ampTimeScale[i])
        return i;
      }
    }
    // shouldnt happen!
    // console.log("unscale amp time value fall through " + v + " " + typeof
    // (v))
    return viewVCE_envs.ampTimeScale.length - 1;
  },

  AmpTimeValueToText(v) {
    return '' + viewVCE_envs.scaleAmpTimeValue(v);
  },

  TextToAmpTimeValue(v) {
    if (v == null || v === '') {
      return 0;
    }
    let val = parseInt(v, 10);
    return '' + viewVCE_envs.unscaleAmpTimeValue(val);
  },

  testConversionFunctions: function() {
    let ok = true;
    let result = [];
    for (let i = 0; i <= 255; i++) {
      let scaled = viewVCE_envs.scaleFreqEnvValue(i);
      let unscaled = viewVCE_envs.unscaleFreqEnvValue(scaled);
      if (i != unscaled) {
        ok = false;
        let err = 'ERROR: freqEnvValue ' + i + ' totext: ' + scaled +
            ' reversed to ' + unscaled;
        result.push(err);
        console.log(err);
      }
    }
    for (let i = 55; i <= 127; i++) {
      let scaled = viewVCE_envs.scaleAmpEnvValue(i);
      let unscaled = viewVCE_envs.unscaleAmpEnvValue(scaled);
      if (i != unscaled) {
        ok = false;
        let err = 'ERROR: ampEnvValue ' + i + ' totext: ' + scaled +
            ' reversed to ' + unscaled;
        result.push(err);
        console.log(err);
      }
    }
    for (let i = 0; i <= 79; i++) {
      let scaled = viewVCE_envs.scaleFreqTimeValue(i);
      let unscaled = viewVCE_envs.unscaleFreqTimeValue(scaled);
      if (i != unscaled) {
        ok = false;
        let err = 'ERROR: ampTimeValue ' + i + ' totext: ' + scaled +
            ' reversed to ' + unscaled;
        result.push(err);
        console.log(err);
      }
    }
    for (let i = 0; i <= 79; i++) {
      let scaled = viewVCE_envs.scaleAmpTimeValue(i);
      let unscaled = viewVCE_envs.unscaleAmpTimeValue(scaled);
      if (i != unscaled) {
        ok = false;
        let err = 'ERROR: ampTimeValue ' + i + ' totext: ' + scaled +
            ' reversed to ' + unscaled;
        result.push(err);
        console.log(err);
      }
    }
    // Spot check some values to ensure the forumlae are computing same
    // values as SYNHCS did (except for the upper range of freq time which
    // we delibrartely change to make the function reversable)
    let expects = [
      {
        arr: [
          [0, 0], [10, 10], [15, 15], [16, 25], [54, 2071], [75, 23436],
          [76, 26306], [77, 29528], [84, 29535], [85, 29535]
        ],
        name: 'freqTimeValue',
        func: viewVCE_envs.scaleFreqTimeValue,
      },
      {
        arr: [
          [0, 0], [20, 20], [40, 40], [41, 45], [54, 205], [75, 2325],
          [76, 2609], [83, 5859], [84, 6576], [85, 6576]
        ],
        name: 'ampTimeValue',
        func: viewVCE_envs.scaleAmpTimeValue,
      },
      {
        arr: [[-61, -61], [-15, -15], [0, 0], [63, 63]],
        name: 'freqValue',
        func: viewVCE_envs.scaleFreqEnvValue,
      },
      {
        arr: [[55, 0], [56, 1], [126, 71], [127, 72]],
        name: 'ampValue',
        func: viewVCE_envs.scaleAmpEnvValue,
      },
      {
        arr: [['', 55]],
        name: 'empty string amp val',
        func: viewVCE_envs.TextToAmpEnvValue,
      },
      {
        arr: [['', 0]],
        name: 'empty string amp time',
        func: viewVCE_envs.TextToAmpTimeValue,
      },
      {
        arr: [['', 0]],
        name: 'empty string freq time',
        func: viewVCE_envs.TextToFreqTimeValue,
      }
    ];

    for (let j = 0; j < expects.length; j++) {
      let expect = expects[j];
      for (let i = 0; i < expect.arr.length; i++) {
        let scaled = expect.func(expect.arr[i][0]);
        if (scaled != expect.arr[i][1]) {
          ok = false;
          let err = 'ERROR: ' + expect.name + '(' + expect.arr[i][0] +
              ') == ' + scaled + ', expected ' + expect.arr[i][1];
          result.push(err);
          console.log(err);
        }
      }
    }
    console.log(
        'viewVCE_envs.testConversionFunctions: ' + (ok ? 'PASS' : 'FAIL'));
    return ok ? null : result;
  },

  unsetFloatVals: function() {
    console.log('unsetFloatVals');
    viewVCE_envs.floatAmpVal = null;
  },

  initFloatVals: function() {
    // floatAmpVal.low/up : float versions of the displayed amp values
    // floatAmpVal.referenceLow/ReferenceUp: snapshot of the current values,
    // but scaled to 100% per envelope.  This allows subsequent gain changes
    // to retain the original env shape, even at extreme gain amounts.

    if (viewVCE_envs.floatAmpVal === null) {
      console.log('initFloatVals');
      viewVCE_envs.floatAmpVal = []
      console.log('initFloatVals init: ', viewVCE_envs.floatAmpVal)
      for (let osc = 0; osc <= viewVCE.vce.Head.VOITAB; osc++) {
        let currentOscGain = viewVCE_envs.raw_computeOscGain(osc);
        if (currentOscGain[0] <= 0.0) {
          // avoid divide by zero!  when original gain was zero, may as well
          // just set the new gain as requested
          currentOscGain[0] = 1.0;
        }
        if (currentOscGain[1] <= 0.0) {
          // avoid divide by zero!  when original gain was zero, may as well
          // just set the new gain as requested
          currentOscGain[1] = 1.0;
        }

        viewVCE_envs.floatAmpVal.push({
          low: [],
          up: [],
          referenceLow: [],
          referenceUp: [],
          origOscGain: currentOscGain
        })
        console.log('initFloatVals top: ' + osc, viewVCE_envs.floatAmpVal)

        for (let eleIndex = 0;
             eleIndex < viewVCE.vce.Envelopes[osc].AmpEnvelope.NPOINTS;
             eleIndex++) {
          let v = viewVCE_envs.scaleAmpEnvValue(
              viewVCE.vce.Envelopes[osc].AmpEnvelope.Table[(eleIndex * 4) + 0]);
          viewVCE_envs.floatAmpVal[osc].low.push(v)
          v = 100.0 / currentOscGain[0] * v;
          viewVCE_envs.floatAmpVal[osc].referenceLow.push(v)
          v = viewVCE_envs.scaleAmpEnvValue(
              viewVCE.vce.Envelopes[osc].AmpEnvelope.Table[(eleIndex * 4) + 1]);
          viewVCE_envs.floatAmpVal[osc].up.push(v)
          v = 100.0 / currentOscGain[1] * v;
          viewVCE_envs.floatAmpVal[osc].referenceUp.push(v)
        }
        console.log('initFloatVals eles: ' + osc, viewVCE_envs.floatAmpVal)
      }
    }
  },

  raw_computeOscGain: function(osc /* zero-based*/) {
    // initial computation of each env's gain from the byte values used to
    // initialize the floatVal's
    // -- all subsequent gain calculations are based on the float vals
    let maxLow = 0;
    let maxUp = 0;
    for (let eleIndex = 0;
         eleIndex < viewVCE.vce.Envelopes[osc].AmpEnvelope.NPOINTS;
         eleIndex++) {
      // low
      let v = viewVCE_envs.scaleAmpEnvValue(
          viewVCE.vce.Envelopes[osc].AmpEnvelope.Table[(eleIndex * 4) + 0]);
      maxLow = Math.max(maxLow, v)
      // up
      v = viewVCE_envs.scaleAmpEnvValue(
          viewVCE.vce.Envelopes[osc].AmpEnvelope.Table[(eleIndex * 4) + 1]);
      maxUp = Math.max(maxUp, v)
    }
    let result = [
      100.0 * maxLow / 72.0,
      100.0 * maxUp / 72.0
    ];  // 72 == MAX allowed Amp Val
    console.log(
        'raw_computeOscGain ' + osc + ' ' + maxLow + ' ' + maxUp + ' -> ' +
        result)
    return result;
  },

  computeOscGain: function(osc /* zero-based*/, lowupboth /*0,1 or 2*/) {
    // lazy init:
    viewVCE_envs.initFloatVals();

    // compute the request gain computation based on the floating point
    // values in the floatAmpVal arrays

    let max = 0;
    for (let eleIndex = 0;
         eleIndex < viewVCE.vce.Envelopes[osc].AmpEnvelope.NPOINTS;
         eleIndex++) {
      if (lowupboth == 0 || lowupboth == 2) {
        // low
        let v = viewVCE_envs.floatAmpVal[osc].low[eleIndex];
        max = Math.max(max, v)
      }
      if (lowupboth == 1 || lowupboth == 2) {
        // up
        let v = viewVCE_envs.floatAmpVal[osc].up[eleIndex];
        max = Math.max(max, v)
      }
    }
    let result = Math.round(100.0 * max / 72.0);  // 72 == MAX allowed Amp Val
    console.log(
        'computeOscGain ' + osc + ' ' + lowupboth + ' ' + max + ' -> ' + result)
    return result;
  },

  setOscGain: function(
      osc /* zero-based */, gain /*0..100*/) {  // both low and up change
    // need to recompute the individual gains in terms of the overall gain
    viewVCE_envs.initFloatVals();

    let gainLow = viewVCE_envs.floatAmpVal[osc].origOscGain[0];
    let gainUp = viewVCE_envs.floatAmpVal[osc].origOscGain[1];

    let origOscGain = Math.max(gainLow, gainUp);
    if (origOscGain <= 0.0) {
      // avoid divide by zero!  when original gain was zero, may as well
      // just set the new gain as requested
      origOscGain = 1.0;
    }

    // proportional change for each point
    let proportion = gain / origOscGain;
    console.log('setOscGain ' + osc + ' ' + gain + ' ' + proportion);
    viewVCE_envs.setGain(osc, gainLow * proportion, 0);
    viewVCE_envs.setGain(osc, gainUp * proportion, 1);
  },

  setGain: async function(
      osc /* zero-based */, gain /*0..100*/, lowup /*0 or 1*/) {
    // this can be called for an OSC whose controls are not actually visible
    // (it's called from the main voice tab).  So if visible, update the
    // <inputs> and send updates to the csurface. If not visible, update the
    // data structures and send changes to the Synergy, but dont update any
    // visible controls.
    //
    // Gain is set by simply multiplying the requested gain against the
    // "reference" gains in the floatAmpVal arrays. (reference gains are
    // scaled to "100%").

    console.log('setGain: ' + osc + ' ' + gain + ' ' + lowup);

    let envOscSelectEle = document.getElementById('envOscSelect');
    let visibleOsc =
        parseInt(envOscSelectEle.value, 10);  // one-based osc index
    visibleOsc--;                             // convert to zero-base

    viewVCE_envs.initFloatVals();

    let referenceFloatVals = lowup == 0 ?
        viewVCE_envs.floatAmpVal[osc].referenceLow :
        viewVCE_envs.floatAmpVal[osc].referenceUp;
    let floatVals = lowup == 0 ? viewVCE_envs.floatAmpVal[osc].low :
                                 viewVCE_envs.floatAmpVal[osc].up;

    for (let eleIndex = 0;
         eleIndex < viewVCE.vce.Envelopes[osc].AmpEnvelope.NPOINTS;
         eleIndex++) {
      let stub = lowup == 0 ? 'envAmpLowVal' : 'envAmpUpVal';
      // reference val is for the "100%" gain case - doesnt change when we
      // change gain
      let refval = referenceFloatVals[eleIndex];
      let floatNewVal = refval * gain / 100.0;
      floatNewVal = Math.min(72, Math.max(0, floatNewVal));
      let newval = Math.round(floatNewVal);
      // non-reference entries in the floatvals need to be kept up to date
      // with the gain change
      floatVals[eleIndex] = floatNewVal;

      console.log(
          '  setGain (point) ' + gain + ' ' + stub + '[' + eleIndex + '] ' +
          refval + ' ' + floatVals[eleIndex] + ' ' + newval);

      if (visibleOsc == osc) {
        let input = document.getElementById(`${stub}[${eleIndex + 1}]`)
        input.value = '' + newval;
        // set value in the input, set in the viewVCE.vce.Envelopes and send
        // to csurface and synergy: we are already debounced, so call raw to
        // avoid a delay in updates:
        viewVCE_envs.raw_onchange(input, floatNewVal);
      } else {
        // no UI field to change, so need to call the backend to send values
        // to Synergy here rather than relying on onchange
        viewVCE.vce.Envelopes[osc].AmpEnvelope.Table[(eleIndex * 4) + lowup] =
            viewVCE_envs.unscaleAmpEnvValue(newval);
        try {
          await UIService.SetEnvEle(
              lowup ? 'setEnvAmpLowVal' : 'setEnvAmpUpVal'.at,
              osc + 1 /* one-based */, eleIndex + 1 /*one-based*/,
              viewVCE.vce.Envelopes[osc]
                  .AmpEnvelope.Table[(eleIndex * 4) + lowup]);
        } catch (exc) {
          // failed - dont change the value
          index.errorNotification(exc);
          return false;
        }
      }

      if (osc == visibleOsc) {
        document.querySelector('#gainAmpLow').value =
            viewVCE_envs.computeOscGain(osc, 0);
        document.querySelector('#gainAmpUp').value =
            viewVCE_envs.computeOscGain(osc, 1);
      }
      // update the Voice tab
      document.getElementById(`OscGain[${osc + 1}]`).value =
          viewVCE_envs.computeOscGain(osc, 2);
      console.log(
          'setGain: update voice tab: ' + osc + ' ' + visibleOsc + ' ' +
          document.getElementById(`OscGain[${osc + 1}]`).value);
    }
    return true;
  },

  onchangeGain: function(ele) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    if (viewVCE_envs
            .supressOnchange) { /*console.log("viewVCE_envs.suppressOnChange");*/
      return;
    }
    viewVCE_envs.deb_onchangeGain(ele);
  },

  deb_onchangeGain: null,  // initialized during init()

  raw_onchangeGain: function(ele) {  // low or up gain change
    let gain = parseInt(ele.value, 10);

    console.log('onchangeGain: ' + ele.id + ' ' + gain);
    let envOscSelectEle = document.getElementById('envOscSelect');
    let osc = parseInt(envOscSelectEle.value, 10);  // one-based osc index
    osc--;                                          // convert to zero-base

    if (ele.id.match(/Low/)) {
      viewVCE_envs.setGain(osc, gain, 0);
      // recompute the "original gains" so that osc-level gain can
      // create proportional changes to reflect this latest change
      viewVCE_envs.floatAmpVal[osc].origOscGain[0] = gain;
    } else {
      viewVCE_envs.setGain(osc, gain, 1);
      // recompute the "original gains" so that osc-level gain can
      // create proportional changes to reflect this latest change
      viewVCE_envs.floatAmpVal[osc].origOscGain[1] = gain;
    }

    viewVCE_voice.sendToCSurface(null, ele.id, gain);
  },

  supressOnChange: false,

  onchangeLoop: async function(ele) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    if (viewVCE_envs.supressOnchange) {
      console.log('viewVCE_envs.suppressOnChange');
      return;
    }

    let eleIndex;
    let envOscSelectEle = document.getElementById('envOscSelect');
    let osc = parseInt(envOscSelectEle.value, 10);  // one-based osc index
    // let envEnvSelectEle = document.getElementById("envEnvSelect");
    // let selectedEnv = parseInt(envEnvSelectEle.value, 10);
    let eleValue = ele.value;

    let pattern = /([A-Za-z]+)\[(\d+)\]/;
    let ret = ele.id.match(pattern);
    if (ret) {
      eleIndex = parseInt(ret[2])
    }

    let env;
    let envid;
    if (ele.id.includes('Freq')) {
      env = viewVCE.vce.Envelopes[osc - 1].FreqEnvelope;
      envid = 'Freq';
    } else {
      env = viewVCE.vce.Envelopes[osc - 1].AmpEnvelope;
      envid = 'Amp';
    }

    let accelLow = 30;  // defaults
    let accelUp = 30;   // defaults

    if (env.ENVTYPE == 1) {
      accelLow = env.SUSTAINPT;
      accelUp = env.LOOPPT;

      // temporarily replace the accellerations with point indexes to
      // make all the validations below less messy than they might be
      env.SUSTAINPT = 0;
      env.LOOPPT = 0;
    }

    // type = 1  : no loop (and LOOPPT and SUSTAINPT are accelleration
    // rates not point positions) type = 2  : S only type = 3  : L and
    // S - L must be before S type = 4  : R and S - R must be before S
    // WARNING: when type1, the LOOPPT and SUSTAINPT values are
    // _acceleration_ rates, not point positions. What a pain.
    console.log('loop change ' + eleIndex + ' ' + eleValue + ' ' + envid);
    console.dir(env);
    if (eleValue == '') {
      // always safe to remove a loop point
      if (env.LOOPPT == eleIndex) {
        env.LOOPPT = 0;
        env.ENVTYPE = 2;  // sustain only
      }
      // can't remove SUSTAIN POINT if there's a loop point
      if (env.SUSTAINPT == eleIndex) {
        if (env.ENVTYPE === 3 || env.ENVTYPE == 4) {
          index.errorNotification(
              'Cannot remove SUSTAIN point if there are LOOP or REPEAT points')
          ele.value = 'S';
          if (env.ENVTYPE == 1) {
            // restore the accellerations
            env.SUSTAINPT = accelLow;
            env.LOOPPT = accelUp;
          }
          return;
        }
        env.SUSTAINPT = 0;
      }
      if (env.LOOPPT === 0 && env.SUSTAINPT === 0) {
        env.ENVTYPE = 1;  // no loops
      } else if (env.LOOPPT === 0) {
        env.ENVTYPE = 2;  // sustain-only
      }
    } else if (eleValue === 'S') {
      if (env.LOOPPT > eleIndex) {
        index.errorNotification(
            'Cannot set SUSTAIN point before LOOP or REPEAT point');
        ele.value = '';
        if (env.ENVTYPE == 1) {
          // restore the accellerations
          env.SUSTAINPT = accelLow;
          env.LOOPPT = accelUp;
        }
        return;
      } else if (env.LOOPPT === eleIndex) {
        // replacing a L or R with an S:
        env.ENVTYPE = 2;
        env.LOOPPT = 0;
        env.SUSTAINPT = eleIndex;
      } else if (env.SUSTAINPT == 0) {
        env.ENVTYPE = 2;
        env.SUSTAINPT = eleIndex;
      } else {
        // S is after an L/R
        // env type remains unchanged
        env.SUSTAINPT = eleIndex;
      }
    } else {  // 'R' or 'L'
      if (env.SUSTAINPT <= eleIndex) {
        index.errorNotification('SUSTAIN point must be after LOOP or REPEAT');
        ele.value = '';
        if (env.ENVTYPE == 1) {
          // restore the accellerations
          env.SUSTAINPT = accelLow;
          env.LOOPPT = accelUp;
        }
        return;
      } else {
        // we use the most recent change to set the env type
        env.ENVTYPE = eleValue === 'L' ? 3 : 4;
        env.LOOPPT = eleIndex;
      }
    }

    if (env.ENVTYPE == 1) {
      // restore the accellerations
      env.SUSTAINPT = accelLow;
      env.LOOPPT = accelUp;
    }

    // validation is OK, but now need to clean up any selects (i.e. if
    // loop point moved, need to set the previous location to '')
    // easiest thing to do is just brute force reset each ele to
    // reflect the value in the envelope
    for (let p = 0; p < 16; p++) {
      //$(`#env${envid}Loop\\[${p + 1}\\] option[value='']`)
      //    .prop('selected', true);
      //$(`#env${envid}Loop\\[${p + 1}\\] option[value='L']`)
      //    .prop('selected', false);
      //$(`#env${envid}Loop\\[${p + 1}\\] option[value='R']`)
      //    .prop('selected', false);
      //$(`#env${envid}Loop\\[${p + 1}\\] option[value='S']`)
      //    .prop('selected', false);
      document.querySelector(`#env${envid}Loop\\[${p + 1}\\]`).value = '';
    }
    if (env.ENVTYPE != 1 && env.SUSTAINPT > 0) {
      //$(`#env${envid}Loop\\[${env.SUSTAINPT}\\] option[value='S']`)
      //    .prop('selected', true);
      document.querySelector(`#env${envid}Loop\\[${env.SUSTAINPT}\\]`).value =
          'S';
    }
    if (env.ENVTYPE != 1 && env.LOOPPT > 0) {
      let v = env.ENVTYPE == 3 ? 'L' : 'R'
      //$(`#env${envid}Loop\\[${env.LOOPPT}\\] option[value='${v}']`)
      //    .prop('selected', true);
      document.querySelector(`#env${envid}Loop\\[${env.LOOPPT}\\]`).value = v;
    }

    // only show accelleration values if type1 envelope
    if (env.ENVTYPE === 1) {
      document.querySelectorAll(`.type1accel div.${envid}`).forEach(el => {
        el.style.display = 'block';
      });
      document.querySelector(`#accel${envid}Low`).value = env.SUSTAINPT;
      document.querySelector(`#accel${envid}Up`).value = env.LOOPPT;
    } else {
      document.querySelectorAll(`.type1accel div.${envid}`).forEach(el => {
        el.style.display = 'none';
      });
    }

    console.log('resulting env: ' + envid);
    console.dir(env);
    try {
      await UIService.SetLoopPoint(
          osc, envid, env.ENVTYPE, env.SUSTAINPT, env.LOOPPT);
    } catch (exc) {
      // failed - dont change the value
      index.errorNotification(exc);
      return false;
    }

    return true;
  },

  copyFrom: function(fromOsc) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    if (viewVCE_envs
            .supressOnchange) { /*console.log("viewVCE_envs.suppressOnChange");*/
      return;
    }
    viewVCE_envs.deb_copyFrom(fromOsc);
  },

  deb_copyFrom: null,

  raw_copyFrom: async function(fromOsc) {
    let oscSelectEle = document.getElementById('envOscSelect');
    let toOsc = oscSelectEle.options[oscSelectEle.selectedIndex].value;
    toOsc = parseInt(toOsc, 10);

    // "copy" means copy all the osc-specific stuff related to the
    // envelopes - but not the patch and detuning fields.

    // abuse JSON to do a deep copy:
    let newEnvelopes =
        JSON.parse(JSON.stringify(viewVCE.vce.Envelopes[fromOsc - 1]))
    // retain the stuff we don't want copied:
    newEnvelopes.FreqEnvelope.OPTCH =
        viewVCE.vce.Envelopes[toOsc - 1].FreqEnvelope.OPTCH;
    newEnvelopes.FreqEnvelope.OHARM =
        viewVCE.vce.Envelopes[toOsc - 1].FreqEnvelope.OHARM;
    newEnvelopes.FreqEnvelope.FDETUN =
        viewVCE.vce.Envelopes[toOsc - 1].FreqEnvelope.FDETUN;

    index.spinnerOn();
    try {
      await UIService.SetEnvelopes(toOsc, newEnvelopes);
      index.spinnerOff();
    } catch (exc) {
      index.spinnerOff();
      // failed - dont change the value
      index.errorNotification(exc);
      return false;
    }
    viewVCE.vce.Envelopes[toOsc - 1] = newEnvelopes

    // update the floatVal so gains work
    viewVCE_envs.unsetFloatVals();
    document.querySelector('#gainAmpLow').value =
        viewVCE_envs.computeOscGain(toOsc - 1, 0);
    document.querySelector('#gainAmpUp').value =
        viewVCE_envs.computeOscGain(toOsc - 1, 1);
    // update the Voice tab
    document.getElementById(`OscGain[${toOsc}]`).value =
        viewVCE_envs.computeOscGain(toOsc - 1, 2)
    console.log(
        'onchange: setGain: update voice tab: ' + toOsc + ' ' +
        document.getElementById(`OscGain[${toOsc}]`).value)

    viewVCE_envs.envChartUpdate(toOsc, -1, true);

    return true;
  },

  onchangeEnvAccel: function(ele) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    if (viewVCE_envs
            .supressOnchange) { /*console.log("viewVCE_envs.suppressOnChange");*/
      return;
    }
    viewVCE_envs.deb_onchangeEnvAccel(ele);
  },

  deb_onchangeEnvAccel: null,  // initialized during init()

  raw_onchangeEnvAccel: async function(ele) {
    if (viewVCE.supressOnchange) { /*console.log("raw
                                      viewVCE.suppressOnChange");*/
      return;
    }
    if (viewVCE_envs.supressOnchange) { /*console.log("raw
                                           viewVCE_envs.suppressOnChange");*/
      return;
    }

    // type1 accelerations are really just the SUSTAIN and LOOP
    // points.  We use the same backend function as the loop change
    // event


    let envOscSelectEle = document.getElementById('envOscSelect');
    let osc = parseInt(envOscSelectEle.value, 10);  // one-based osc index
    // let envEnvSelectEle = document.getElementById("envEnvSelect");
    // let selectedEnv = parseInt(envEnvSelectEle.value, 10);
    let eleValue = index.checkInputElementValue(ele);
    if (eleValue == undefined) {
      return;
    }

    // console.log("changed: " + ele.id + " val: " + ele.value);

    let env;
    let envid;
    // id is accelFreqUp , accelAmpLow etc.
    if (ele.id.includes('Freq')) {
      env = viewVCE.vce.Envelopes[osc - 1].FreqEnvelope;
      envid = 'Freq';
    } else {
      env = viewVCE.vce.Envelopes[osc - 1].AmpEnvelope;
      envid = 'Amp';
    }
    if (ele.id.includes('Low')) {
      env.SUSTAINPT = parseInt(eleValue, 10);
    } else {
      env.LOOPPT = parseInt(eleValue, 10);
    }

    console.log('ACCEL change ' + ele.id + ' ' + eleValue + ' ' + envid);

    try {
      await UIService.SetLoopPoint(
          osc, envid, env.ENVTYPE, env.SUSTAINPT, env.LOOPPT);
    } catch (exc) {
      // failed - dont change the value
      index.errorNotification(exc);
      return false;
    }
    // currently there's no visual rendering of the loop, so no
    // need to refresh the chart
    //				viewVCE_envs.envChartUpdate(osc,
    // selectedEnv, false);
    viewVCE_voice.sendToCSurface(ele, ele.id, parseInt(eleValue, 10));
    return true;
  },

  onchange: async function(ele, extraarg) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    if (viewVCE_envs
            .supressOnchange) { /*console.log("viewVCE_envs.suppressOnChange");*/
      return;
    }
    viewVCE_envs.deb_onchange(ele, extraarg);
  },

  deb_onchange: null,  // initialized during init()

  raw_onchange: async function(ele, extraarg) {
    if (viewVCE.supressOnchange) { /*console.log("raw
                                      viewVCE.suppressOnChange");*/
      return;
    }
    if (viewVCE_envs.supressOnchange) { /*console.log("raw
                                           viewVCE_envs.suppressOnChange");*/
      return;
    }

    let envOscSelectEle = document.getElementById('envOscSelect');
    let osc = parseInt(envOscSelectEle.value, 10);  // one-based osc index
    let envEnvSelectEle = document.getElementById('envEnvSelect');
    let selectedEnv = parseInt(envEnvSelectEle.value, 10);

    // Don't call checkInoutElementValue() - it assumes that there is
    // no scaling and would apply the "byte" min/max to the "text"
    // scaled value
    //	  let value = index.checkInputElementValue(ele);
    let value = parseInt(ele.value, 10);

    if (value == undefined) {
      return;
    }
    // console.log("in onchange - value: " + value + " " +
    // typeof(value))

    let pattern = /([A-Za-z]+)\[(\d+)\]/;
    let funcName;
    let eleIndex;
    let ret = ele.id.match(pattern);
    let bytevalue;
    if (ret) {
      let fieldType = ret[1];
      funcName = 'set' + fieldType.charAt(0).toUpperCase() + fieldType.slice(1);
      eleIndex = parseInt(ret[2])
      // now scale the value to the byte value the synergy wants to
      // see:
      switch (fieldType) {
        case 'envFreqLowVal':
          bytevalue = viewVCE_envs.unscaleFreqEnvValue(value);
          viewVCE.vce.Envelopes[osc - 1]
              .FreqEnvelope.Table[((eleIndex - 1) * 4) + 0] = bytevalue;
          break;
        case 'envFreqUpVal':
          bytevalue = viewVCE_envs.unscaleFreqEnvValue(value);
          viewVCE.vce.Envelopes[osc - 1]
              .FreqEnvelope.Table[((eleIndex - 1) * 4) + 1] = bytevalue;
          break;
        case 'envFreqLowTime':
          bytevalue = viewVCE_envs.unscaleFreqTimeValue(value);
          viewVCE.vce.Envelopes[osc - 1]
              .FreqEnvelope.Table[((eleIndex - 1) * 4) + 2] = bytevalue;
          break;
        case 'envFreqUpTime':
          bytevalue = viewVCE_envs.unscaleFreqTimeValue(value);
          viewVCE.vce.Envelopes[osc - 1]
              .FreqEnvelope.Table[((eleIndex - 1) * 4) + 3] = bytevalue;
          break;
        case 'envAmpLowVal':
          bytevalue = viewVCE_envs.unscaleAmpEnvValue(value);
          viewVCE.vce.Envelopes[osc - 1]
              .AmpEnvelope.Table[((eleIndex - 1) * 4) + 0] = bytevalue;
          if (extraarg == undefined) {
            // update the floatVal so gains work
            viewVCE_envs.unsetFloatVals();
            document.querySelector('#gainAmpLow').value =
                viewVCE_envs.computeOscGain(osc - 1, 0);
            // update the Voice tab
            document.getElementById(`OscGain[${osc}]`).value =
                viewVCE_envs.computeOscGain(osc - 1, 2)
            console.log(
                'onchange: setGain: update voice tab: ' + osc + 1 + ' ' +
                document.getElementById(`OscGain[${osc}]`).value)
          }
          break;
        case 'envAmpUpVal':
          bytevalue = viewVCE_envs.unscaleAmpEnvValue(value);
          viewVCE.vce.Envelopes[osc - 1]
              .AmpEnvelope.Table[((eleIndex - 1) * 4) + 1] = bytevalue;
          if (extraarg == undefined) {
            // update the floatVal so gains work
            viewVCE_envs.unsetFloatVals();
            document.querySelector('#gainAmpUp').value =
                viewVCE_envs.computeOscGain(osc - 1, 1);
            document.getElementById(`OscGain[${osc}]`).value =
                viewVCE_envs.computeOscGain(osc - 1, 2)
            console.log(
                'onchange: setGain: update voice tab: ' + osc + 1 + ' ' +
                document.getElementById(`OscGain[${osc}]`).value)
          }
          break;
        case 'envAmpLowTime':
          bytevalue = viewVCE_envs.unscaleAmpTimeValue(value);
          viewVCE.vce.Envelopes[osc - 1]
              .AmpEnvelope.Table[((eleIndex - 1) * 4) + 2] = bytevalue;
          break;
        case 'envAmpUpTime':
          bytevalue = viewVCE_envs.unscaleAmpTimeValue(value);
          viewVCE.vce.Envelopes[osc - 1]
              .AmpEnvelope.Table[((eleIndex - 1) * 4) + 3] = bytevalue;
          break;
      }
    }
    console.log(
        'env ele change ' + ele.id + ' rawval: ' + ele.value + ' -> ' + value +
        ' -> ' + bytevalue);
    // console.log("in onchange - bytevalue: " + bytevalue + " " +
    // typeof(bytevalue))

    try {
      await UIService.SetEnvEle(funcName, osc, eleIndex, bytevalue);
    } catch (exc) {
      // failed - dont change the value
      index.errorNotification(exc);
      return false;
    }
    viewVCE_envs.envChartUpdate(osc, selectedEnv, false);
    viewVCE_voice.sendToCSurface(ele, ele.id, bytevalue);
    return true;
  },

  changeEnvPoints: async function(whichEnv, increment) {
    let envOscSelectEle = document.getElementById('envOscSelect');
    let osc = parseInt(envOscSelectEle.value, 10);  // one-based osc index
    let envEnvSelectEle = document.getElementById('envEnvSelect');
    let selectedEnv = parseInt(envEnvSelectEle.value, 10);
    let envs = viewVCE.vce.Envelopes[osc - 1];

    console.log('#points changed: ' + whichEnv + ' increment: ' + increment);

    // reset the floating point backing arrays: FIXME: this may mean
    // that someone who is in midst of twiddling gain, then
    // adds/deletes a point, then expects gain to keep working may get
    // a suprised distortion in the envelope shape.
    viewVCE_envs.unsetFloatVals();

    let changed = false;
    if (whichEnv === 'freq') {
      let newlen = envs.FreqEnvelope.NPOINTS + increment;
      if (newlen >= 1 && newlen <= 16) {
        if (increment === -1 && envs.FreqEnvelope.ENVTYPE != 1 &&
            (envs.FreqEnvelope.LOOPPT > newlen ||
             envs.FreqEnvelope.SUSTAINPT > newlen)) {
          index.errorNotification(
              'Cannot remove envelope point with SUSTAIN/LOOP marker.  Remove the SUSTAIN/LOOP marker before trying to remove points.');
          return;
        }
        envs.FreqEnvelope.NPOINTS = newlen;
        changed = true;
      }
    } else {
      let newlen = envs.AmpEnvelope.NPOINTS + increment;
      if (newlen >= 1 && newlen <= 16) {
        if (increment === -1 && envs.AmpEnvelope.ENVTYPE != 1 &&
            (envs.AmpEnvelope.LOOPPT > newlen ||
             envs.AmpEnvelope.SUSTAINPT > newlen)) {
          index.errorNotification(
              'Cannot remove envelope point with SUSTAIN/LOOP marker.  Remove the SUSTAIN/LOOP marker before trying to remove points.');
          return;
        }
        envs.AmpEnvelope.NPOINTS = newlen;
        changed = true;
      }
    }
    if (changed) {
      try {
        await UIService.SetOscEnvLengths(
            osc, envs.FreqEnvelope.NPOINTS, envs.AmpEnvelope.NPOINTS);
      } catch (exc) {
        // failed - dont change the value
        index.errorNotification(exc);
        return false;
      }
      viewVCE_envs.envChartUpdate(osc, selectedEnv, true);
    }
  },

  uncompressEnvelopes: function() {
    // the first time we evaluate this viewVCE.vce, the envelopes may
    // be compressed.  To make it easier to add/remove filters in the
    // editor, we rewrite the envelopes arrays such each has the max
    // amount of elements and each are initialized as SYNHCS does.
    if (viewVCE.vce.Extra['uncompressedEnvelopes'] != undefined) {
      // no need to do it again
      return;
    }
    // only need to worry about the number of oscillators in the
    // Envelopes table; any addition osc's added will automatically
    // fill in "full length" envelopes (but use the length of the
    // array not the current value of VOITAB lowered the number of
    // osc's)
    for (let i = 0; i < viewVCE.vce.Envelopes.length; i++) {
      const FULL_LENGTH = 16 * 4;  // 16 rows, each with 4 values
      for (let j = viewVCE.vce.Envelopes[i].FreqEnvelope.Table.length;
           j < FULL_LENGTH; j++) {
        viewVCE.vce.Envelopes[i].FreqEnvelope.Table.push(0);
      }
      for (let j = viewVCE.vce.Envelopes[i].AmpEnvelope.Table.length;
           j < FULL_LENGTH; j++) {
        viewVCE.vce.Envelopes[i].AmpEnvelope.Table.push(0);
      }
    }
    viewVCE.vce.Extra.uncompressedEnvelopes = true;
  },

  changeTimeScale: function(val) {
    this.chart.options.scales.xAxes[0].type = val;
    this.chart.update();
  },

  changeFreqScale: function(val) {
    this.chart.options.scales.yAxes[0].type = val;
    this.chart.update();
  },

  changeTimeZoom: function(val) {
    let div = document.getElementById('envZoomDiv');
    div.style.width = val;
    // this.chart.update();
  },

  // XREF: needs to match the dataset order inside envChangeUpdate()
  valFieldNameByDatasetIdx:
      ['envFreqLowVal', 'envFreqUpVal', 'envAmpLowVal', 'envAmpUpVal'],

  timeFieldNameByDatasetIdx:
      ['envFreqLowTime', 'envFreqUpTime', 'envAmpLowTime', 'envAmpUpTime'],


  envChartUpdate: function(oscNum, envNum, animate) {
    viewVCE_envs.supressOnchange = true;

    viewVCE_envs.uncompressEnvelopes();

    let envCopySelectEle = document.getElementById('envCopySelect');
    // remove old options:
    while (envCopySelectEle.firstChild) {
      envCopySelectEle.removeChild(envCopySelectEle.firstChild);
    }
    // hide the copy selector for All or cases where there are no
    // filters, or when we're not in voicing mode
    document.querySelector('#envCopySelectDiv').style.display = 'none';
    if (viewVCE_voice.voicingMode) {
      document.querySelector('#envCopySelectDiv').style.display = 'block';
      // populate options in the select with only "other" osc (i.e.
      // "this" osc should be not shown or at least unselectable)

      // first element is empty to avoid confusing the user if they
      // havent selected something:
      let option = document.createElement('option');
      option.value = -1;
      option.innerHTML = '';
      envCopySelectEle.appendChild(option);

      for (let i = 0; i <= viewVCE.vce.Head.VOITAB; i++) {
        if ((i + 1) != oscNum) {
          let option = document.createElement('option');
          option.value = i + 1;
          option.innerHTML = i + 1;
          envCopySelectEle.appendChild(option);
        }
      }
    }

    let oscIndex = oscNum - 1;
    let envelopes = viewVCE.vce.Envelopes[oscIndex];

    let pointStyleMetadata = [
      // order needs to match the dataset array
      {color: 2, loopPt: -1, repeatPt: -1, sustainPt: -1},
      {color: 3, loopPt: -1, repeatPt: -1, sustainPt: -1},
      {color: 0, loopPt: -1, repeatPt: -1, sustainPt: -1},
      {color: 1, loopPt: -1, repeatPt: -1, sustainPt: -1},
    ];

    function annotatePointStyle(ctx) {
      let styleMeta = filteredPointStyleMetadata[ctx.datasetIndex]
          // console.log("annotate point ctx ",ctx)
          let img = new Image(14, 14);
      if (styleMeta.sustainPt == ctx.dataIndex) {
        img.src = `static/images/loopS-${styleMeta.color}.png`;
        return img;
      } else if (styleMeta.loopPt == ctx.dataIndex) {
        img.src = `static/images/loopL-${styleMeta.color}.png`;
        return img;
      } else if (styleMeta.repeatPt == ctx.dataIndex) {
        img.src = `static/images/loopR-${styleMeta.color}.png`;
        return img;
      }
      return 'circle'
    }

    // XREF: order needs to match the valFieldNameByDatasetIdx and
    // timeFieldNameByDatasetIdx above
    let freqLowIdx = 0;
    let freqUpIdx = 1;
    let ampLowIdx = 2;
    let ampUpIdx = 3;
    let datasets = [
      {
        label: 'Freq Low',
        xAxisID: 'time-axis',
        yAxisID: 'freq-axis',
        fill: false,
        lineTension: 0,
        pointRadius: 2,
        pointHitRadius: 5,
        pointStyle: annotatePointStyle,
        showLine: true,
        borderWidth: 3,
        backgroundColor: viewVCE.chartColors[2],
        borderColor: viewVCE.chartColors[2],
        data: []
      },
      {
        label: 'Freq Up',
        xAxisID: 'time-axis',
        yAxisID: 'freq-axis',
        fill: false,
        lineTension: 0,
        pointRadius: 2,
        pointHitRadius: 5,
        pointStyle: annotatePointStyle,
        showLine: true,
        borderWidth: 3,
        backgroundColor: viewVCE.chartColors[3],
        borderColor: viewVCE.chartColors[3],
        data: []
      },
      {
        label: 'Amp Low',
        xAxisID: 'time-axis',
        yAxisID: 'amp-axis',
        fill: false,
        lineTension: 0,
        pointRadius: 2,
        pointHitRadius: 5,
        pointStyle: annotatePointStyle,
        showLine: true,
        borderWidth: 3,
        backgroundColor: viewVCE.chartColors[0],
        borderColor: viewVCE.chartColors[0],
        data: []
      },
      {
        label: 'Amp Up',
        xAxisID: 'time-axis',
        yAxisID: 'amp-axis',
        fill: false,
        lineTension: 0,
        pointRadius: 2,
        pointHitRadius: 5,
        pointStyle: annotatePointStyle,
        showLine: true,
        borderWidth: 3,
        backgroundColor: viewVCE.chartColors[1],
        borderColor: viewVCE.chartColors[1],
        data: []
      },
    ];

    viewVCE_voice.sendToCSurface(
        null, `num-freq-env-points`, envelopes.FreqEnvelope.NPOINTS);
    viewVCE_voice.sendToCSurface(
        null, `num-amp-env-points`, envelopes.AmpEnvelope.NPOINTS);

    // clear old values:
    document.querySelector('#envTable td.val input').value = '';
    document.querySelector('#envTable td.total div').innerHTML = ('');
    // clear the loop points
    document.querySelector(`#envTable select option[value='']`).selected = true;
    document.querySelector(`#envTable select option[value='L']`).selected =
        false;
    document.querySelector(`#envTable select option[value='S']`).selected =
        false;
    document.querySelector(`#envTable select option[value='R']`).selected =
        false;

    // fill in freq env data:

    // scaling algorithms derived from DISVAL: in OSCDSP.Z80

    let totalTimeLow = 0;
    let totalTimeUp = 0;
    // let lastFreqLow = 0;
    // let lastFreqUp = 0;
    // let lastAmpLow = 0;
    // let lastAmpUp = 0;

    for (let i = 0; i < 16; i++) {
      // completely hide the rows for rows not used by either envelope
      let tr = document.querySelector(`#envTableTr\\[${(i + 1)}\\]`);
      console.log('tr', i, tr)
      if (i <
          Math.max(
              envelopes.FreqEnvelope.NPOINTS, envelopes.AmpEnvelope.NPOINTS)) {
        tr.style.display = 'block';
      }
      else {
        tr.style.display = 'none';
      }
    }
    if (viewVCE_voice.voicingMode) {
      document.querySelectorAll('.listplusminus div').forEach(el => {
        el.style.display = 'block';
      });
    } else {
      document.querySelectorAll('.listplusminus div').forEach(el => {
        el.style.display = 'none';
      });
    }
    viewVCE_voice.sendToCSurface(
        null, `num-freq-env-points`, envelopes.FreqEnvelope.NPOINTS);
    viewVCE_voice.sendToCSurface(
        null, `num-amp-env-points`, envelopes.AmpEnvelope.NPOINTS);

    // only show accelleration values if type1 envelope
    if (envelopes.FreqEnvelope.ENVTYPE === 1) {
      document.querySelectorAll('.type1accel div.Freq').forEach(el => {
        el.style.display = 'block';
      });
      document.querySelector('#accelFreqLow').value =
          envelopes.FreqEnvelope.SUSTAINPT;
      document.querySelector('#accelFreqUp').value =
          envelopes.FreqEnvelope.LOOPPT;
      viewVCE_voice.sendToCSurface(null, `freq-env-accel-visible`, 1);
      if (animate) {
        viewVCE_voice.sendToCSurface(
            null, `accelFreqLow`, envelopes.FreqEnvelope.SUSTAINPT);
        viewVCE_voice.sendToCSurface(
            null, `accelFreqUp`, envelopes.FreqEnvelope.LOOPPT);
      }
    } else {
      document.querySelectorAll('.type1accel div.Freq').forEach(el => {
        el.style.display = 'none';
      });
      viewVCE_voice.sendToCSurface(null, `freq-env-accel-visible`, 0);
      if (animate) {
        viewVCE_voice.sendToCSurface(null, `accelFreqLow`, 0);
        viewVCE_voice.sendToCSurface(null, `acceFreqUp`, 0);
      }
    }
    // only show accelleration values if type1 envelope
    if (envelopes.AmpEnvelope.ENVTYPE === 1) {
      document.querySelectorAll('.type1accel div.Amp').forEach(el => {
        el.style.display = 'block';
      });
      document.querySelector('#accelAmpLow').value =
          envelopes.AmpEnvelope.SUSTAINPT;
      document.querySelector('#accelAmpUp').value -
          envelopes.AmpEnvelope.LOOPPT;
      viewVCE_voice.sendToCSurface(null, `amp-env-accel-visible`, 1);
      if (animate) {
        viewVCE_voice.sendToCSurface(
            null, `accelAmpLow`, envelopes.AmpEnvelope.SUSTAINPT);
        viewVCE_voice.sendToCSurface(
            null, `accelAmpUp`, envelopes.AmpEnvelope.LOOPPT);
      }
    } else {
      document.querySelectorAll('.type1accel div.Amp').forEach(el => {
        el.style.display = 'none';
      });
      viewVCE_voice.sendToCSurface(null, `amp-env-accel-visible`, 0);
      if (animate) {
        viewVCE_voice.sendToCSurface(null, `accelAmpLow`, 0);
        viewVCE_voice.sendToCSurface(null, `accelAmpUp`, 0);
      }
    }
    document.querySelector('#gainAmpLow').value =
        viewVCE_envs.computeOscGain(oscIndex, 0);
    document.querySelector('#gainAmpUp').value =
        viewVCE_envs.computeOscGain(oscIndex, 1);
    if (animate) {
      console.log(
          'GAIN TO CS: ' + document.querySelector('#gainAmpLow').value +
          ' and ' + document.querySelector('#gainAmpUp').value)
      viewVCE_voice.sendToCSurface(
          null, `gainAmpLow`, document.querySelector('#gainAmpLow').value);
      viewVCE_voice.sendToCSurface(
          null, `gainAmpUp`, document.querySelector('#gainAmpUp').value);
    }

    for (let i = envelopes.FreqEnvelope.NPOINTS; i < 16; i++) {
      document.querySelector(`#envFreqLoop\\[${i + 1}\\]`).style.display =
          'none';
      document.querySelector(`#envFreqLowVal\\[${i + 1}\\]`).style.display =
          'none';
      document.querySelector(`#envFreqUpVal\\[${i + 1}\\]`).style.display =
          'none';
      if (i != 0) {
        // no xxxTime[1] on the display
        document.querySelector(`#envFreqLowTime\\[${i + 1}\\]`).style.display =
            'none';
        document.querySelector(`#envFreqUpTime\\[${i + 1}\\]`).style.display =
            'none';
      }
      if (animate) {
        viewVCE_voice.sendToCSurface(null, `envFreqLowVal[${i + 1}]`, 0);
        viewVCE_voice.sendToCSurface(null, `envFreqUpVal[${i + 1}]`, 0);
        if (i != 0) {
          // no xxxTime[1] on the display
          viewVCE_voice.sendToCSurface(null, `envFreqLowTime[${i + 1}]`, 0);
          viewVCE_voice.sendToCSurface(null, `envFreqUpTime[${i + 1}]`, 0);
        }
      }
    }
    for (let i = 0; i < envelopes.FreqEnvelope.NPOINTS; i++) {
      document.querySelector(`#envFreqLoop\\[${i + 1}\\]`).style.display =
          'block';
      document.querySelector(`#envFreqLowVal\\[${i + 1}\\]`).style.display =
          'block';
      document.querySelector(`#envFreqUpVal\\[${i + 1}\\]`).style.display =
          'block';
      if (i != 0) {
        // no xxxTime[1] on the display
        document.querySelector(`#envFreqLowTime\\[${i + 1}\\]`).style.display =
            'block';
        document.querySelector(`#envFreqUpTime\\[${i + 1}\\]`).style.display =
            'block';
      }
      // table is logically in groups of 4
      let freqLow = viewVCE_envs.scaleFreqEnvValue(
          envelopes.FreqEnvelope.Table[i * 4 + 0]);
      let freqUp = viewVCE_envs.scaleFreqEnvValue(
          envelopes.FreqEnvelope.Table[i * 4 + 1]);
      let timeLow = viewVCE_envs.scaleFreqTimeValue(
          envelopes.FreqEnvelope.Table[i * 4 + 2], i == 0);
      let timeUp = viewVCE_envs.scaleFreqTimeValue(
          envelopes.FreqEnvelope.Table[i * 4 + 3], i == 0);

      if (animate) {
        viewVCE_voice.sendToCSurface(
            null, `envFreqLowVal[${i + 1}]`,
            envelopes.FreqEnvelope.Table[i * 4 + 0]);
        viewVCE_voice.sendToCSurface(
            null, `envFreqUpVal[${i + 1}]`,
            envelopes.FreqEnvelope.Table[i * 4 + 1]);
        viewVCE_voice.sendToCSurface(
            null, `envFreqLowTime[${i + 1}]`,
            envelopes.FreqEnvelope.Table[i * 4 + 2]);
        viewVCE_voice.sendToCSurface(
            null, `envFreqUpTime[${i + 1}]`,
            envelopes.FreqEnvelope.Table[i * 4 + 3]);
      }

      if (i == 0) {
        // first row's time values are fixed at zero (the entries in
        // the table are used for wave/kprop markers)
        timeLow = 0;
        timeUp = 0;
      }
      // lastFreqLow = freqLow;
      // lastFreqUp = freqUp;
      totalTimeLow += timeLow;
      totalTimeUp += timeUp;

      datasets[freqLowIdx].data.push({x: totalTimeLow, y: freqLow});
      datasets[freqUpIdx].data.push({x: totalTimeUp, y: freqUp});

      document.getElementById(`envFreqLowVal[${i + 1}]`).value = freqLow;
      document.getElementById(`envFreqUpVal[${i + 1}]`).value = freqUp;
      if (i !== 0) {
        document.getElementById(`envFreqLowTime[${i + 1}]`).value = timeLow;
        document.getElementById(`envFreqUpTime[${i + 1}]`).value = timeUp;
      }
      document.getElementById(`envFreqTotLowTime[${i + 1}]`).innerHTML =
          totalTimeLow;
      document.getElementById(`envFreqTotUpTime[${i + 1}]`).innerHTML =
          totalTimeUp;

      if (envelopes.FreqEnvelope.ENVTYPE != 1) {
        if (envelopes.FreqEnvelope.SUSTAINPT == (i + 1)) {
          //$(`#envFreqLoop\\[${i + 1}\\] option[value='S']`)
          //    .prop('selected', true);
          document.querySelector(`#envFreqLoop\\[${i + 1}\\]`).value = 'S';
          pointStyleMetadata[freqLowIdx].sustainPt = i;
          pointStyleMetadata[freqUpIdx].sustainPt = i;
        }
        if (envelopes.FreqEnvelope.LOOPPT == (i + 1)) {
          let v = envelopes.FreqEnvelope.ENVTYPE == 3 ? 'L' : 'R'
          //$(`#envFreqLoop\\[${i + 1}\\] option[value='${v}']`)
          //    .prop('selected', true);
          document.querySelector(`#envFreqLoop\\[${i + 1}\\]`).value = v;
          if (v === 'L') {
            pointStyleMetadata[freqLowIdx].loopPt = i;
            pointStyleMetadata[freqUpIdx].loopPt = i;
          } else {
            pointStyleMetadata[freqLowIdx].repeatPt = i;
            pointStyleMetadata[freqUpIdx].repeatPt = i;
          }
        }
      }
    }
    // let maxTotalTime = Math.max(totalTimeLow, totalTimeUp);

    totalTimeLow = 0;
    totalTimeUp = 0;


    for (let i = envelopes.FreqEnvelope.NPOINTS; i < 16; i++) {
      // hide unused rows

      document.querySelector(`#envAmpLoop\\[${i + 1}\\]`).style.display =
          'none';
      document.querySelector(`#envAmpLowVal\\[${i + 1}\\]`).style.display =
          'none';
      document.querySelector(`#envAmpUpVal\\[${i + 1}\\]`).style.display =
          'none';
      document.querySelector(`#envAmpLowTime\\[${i + 1}\\]`).style.display =
          'none';
      document.querySelector(`#envAmpUpTime\\[${i + 1}\\]`).style.display =
          'none';

      if (animate) {
        viewVCE_voice.sendToCSurface(null, `envAmpLowVal[${i + 1}]`, 0);
        viewVCE_voice.sendToCSurface(null, `envAmpUpVal[${i + 1}]`, 0);
        viewVCE_voice.sendToCSurface(null, `envAmpLowTime[${i + 1}]`, 0);
        viewVCE_voice.sendToCSurface(null, `envAmpUpTime[${i + 1}]`, 0);
      }
    }

    // Amp envelopes have an implicit start point at time zero, value
    // zero
    datasets[ampLowIdx].data.push({x: 0, y: 0});
    datasets[ampUpIdx].data.push({x: 0, y: 0});

    for (let i = 0; i < envelopes.AmpEnvelope.NPOINTS; i++) {
      document.querySelector(`#envAmpLoop\\[${i + 1}\\]`).style.display =
          'block';
      document.querySelector(`#envAmpLowVal\\[${i + 1}\\]`).style.display =
          'block';
      document.querySelector(`#envAmpUpVal\\[${i + 1}\\]`).style.display =
          'block';
      document.querySelector(`#envAmpLowTime\\[${i + 1}\\]`).style.display =
          'block';
      document.querySelector(`#envAmpUpTime\\[${i + 1}\\]`).style.display =
          'block';

      // table is logically in groups of 4.
      // "j" accounts for the difference in column index due to the
      // row-spanning separators (only in i==0):
      // let j = (i == 0) ? 11 : 9;

      //	    console.dir(tr);
      //	    console.dir(tr.find('td:eq(' +(j+0)+ ')'));
      let isLast = (i + 1) >= envelopes.AmpEnvelope.NPOINTS;
      let ampLow =
          viewVCE_envs.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i * 4 + 0]);
      let ampUp =
          viewVCE_envs.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i * 4 + 1]);
      let timeLow = viewVCE_envs.scaleAmpTimeValue(
          envelopes.AmpEnvelope.Table[i * 4 + 2]);
      let timeUp = viewVCE_envs.scaleAmpTimeValue(
          envelopes.AmpEnvelope.Table[i * 4 + 3]);

      if (animate) {
        viewVCE_voice.sendToCSurface(
            null, `envAmpLowVal[${i + 1}]`,
            envelopes.AmpEnvelope.Table[i * 4 + 0]);
        viewVCE_voice.sendToCSurface(
            null, `envAmpUpVal[${i + 1}]`,
            envelopes.AmpEnvelope.Table[i * 4 + 1]);
        viewVCE_voice.sendToCSurface(
            null, `envAmpLowTime[${i + 1}]`,
            envelopes.AmpEnvelope.Table[i * 4 + 2]);
        viewVCE_voice.sendToCSurface(
            null, `envAmpUpTime[${i + 1}]`,
            envelopes.AmpEnvelope.Table[i * 4 + 3]);
      }

      // lastAmpLow = ampLow;
      // lastAmpUp = ampUp;
      totalTimeLow += timeLow;
      totalTimeUp += timeUp;

      datasets[ampLowIdx].data.push({x: totalTimeLow, y: ampLow});
      datasets[ampUpIdx].data.push({x: totalTimeUp, y: ampUp});

      //	    console.dir(datasets[ampLowIdx]);

      document.getElementById(`envAmpLowVal[${i + 1}]`).value = ampLow;
      document.getElementById(`envAmpUpVal[${i + 1}]`).value = ampUp;
      document.getElementById(`envAmpLowTime[${i + 1}]`).value = timeLow;
      document.getElementById(`envAmpUpTime[${i + 1}]`).value = timeUp;
      document.getElementById(`envAmpTotLowTime[${i + 1}]`).innerHTML =
          totalTimeLow;
      document.getElementById(`envAmpTotUpTime[${i + 1}]`).innerHTML =
          totalTimeUp;

      if (isLast) {
        document.getElementById(`envAmpLowVal[${i + 1}]`).disabled = true;
        document.getElementById(`envAmpUpVal[${i + 1}]`).disabled = true;
      } else if (viewVCE_voice.voicingMode) {
        document.getElementById(`envAmpLowVal[${i + 1}]`).disabled = false;
        document.getElementById(`envAmpUpVal[${i + 1}]`).disabled = false;
      }

      if (envelopes.AmpEnvelope.ENVTYPE != 1) {
        if (envelopes.AmpEnvelope.SUSTAINPT == (i + 1)) {
          //$(`#envAmpLoop\\[${i + 1}\\] option[value='S']`)
          //  .prop('selected', true);
          document.querySelector(`#envAmpLoop\\[${i + 1}\\]`).value = 'S';
          // we draw an extra point for amp curve - so the index of
          // the loop point is i+1
          pointStyleMetadata[ampLowIdx].sustainPt = i + 1;
          pointStyleMetadata[ampUpIdx].sustainPt = i + 1;
        }
        if (envelopes.AmpEnvelope.LOOPPT == (i + 1)) {
          let v = envelopes.AmpEnvelope.ENVTYPE == 3 ? 'L' : 'R'
          //$(`#envAmpLoop\\[${i + 1}\\] option[value='${v}']`)
          //    .prop('selected', true);
          document.querySelector(`#envAmpLoop\\[${i + 1}\\]`).value = v;
          if (v === 'L') {
            // we draw an extra point for amp curve - so the index of
            // the loop point is i+1
            pointStyleMetadata[ampLowIdx].loopPt = i + 1;
            pointStyleMetadata[ampUpIdx].loopPt = i + 1;
          } else {
            pointStyleMetadata[ampLowIdx].repeatPt = i + 1;
            pointStyleMetadata[ampUpIdx].repeatPt = i + 1;
          }
        }
      }
    }

    // maxTotalTime = Math.max(maxTotalTime, totalTimeLow);
    // maxTotalTime = Math.max(maxTotalTime, totalTimeUp);

    /*
need to confirm the actual behavior of the envelopes to determine if
this visualization makes sense datasets[freqLowIdx].data.push({x:
maxTotalTime, y: lastFreqLow}); datasets[freqUpIdx].data.push( {x:
maxTotalTime,  y: lastFreqUp}); datasets[ampLowIdx].data.push({x:
maxTotalTime, y: lastAmpLow}); datasets[ampUpIdx].data.push( {x:
maxTotalTime,  y: lastAmpUp});
    */

    //	console.dir(datasets);

    let animation_duration = animate ? 1000 : 0;

    let filteredPointStyleMetadata = [];
    let filteredDatasets = [];
    if (envNum < 0) {
      // all of them:
      filteredDatasets = datasets;
      filteredPointStyleMetadata = pointStyleMetadata;
    } else {
      filteredDatasets.push(datasets[envNum])
      filteredPointStyleMetadata.push(pointStyleMetadata[envNum])
    }
    let ctx = document.getElementById('envChart').getContext('2d');
    if (viewVCE_envs.chart != null) {
      viewVCE_envs.chart.destroy();
    }
    let timeAxisType = document.getElementById('timeScale').value;
    let freqAxisType = document.getElementById('freqScale').value;

    viewVCE_envs.chart = new Chart(ctx, {

      type: 'scatter',
      data: {
        //		labels: ['','','','','','','','','','',
        //'','','','','','','','','','',
        //'','','','','','','','','','','',''],
        datasets: filteredDatasets
      },

      // Configuration options go here
      options: {
        animation: {duration: animation_duration},
        tooltips: {
          mode: 'index',
        },
        scales: {
          xAxes: [{
            position: 'bottom',
            id: 'time-axis',
            type: timeAxisType,
            gridLines: {
              color: '#666',
              display: true,
              drawBorder: true,
              drawOnChartArea: true
            },
            scaleLabel: {display: true, labelString: 'Time (ms)'},

            // after a pan, the graph sometimes shows poorly formatted
            // values for min and max (apparently bypasses the tick
            // callback?)
            // Workaround by just not displaying them
            afterTickToLabelConversion: function(scaleInstance) {
              // set the first and last tick to null so it does not
              // display note, ticks[0] is the last tick and
              // ticks[length - 1] is the first
              scaleInstance.ticks[0] = null;
              scaleInstance.ticks[scaleInstance.ticks.length - 1] = null;

              // need to do the same thing for this similiar array
              // which is used internally
              // scaleInstance.ticksAsNumbers[0] = null;
              // scaleInstance.ticksAsNumbers[scaleInstance.ticksAsNumbers.length
              // - 1] = null;
            },

            ticks: {
              precision: 2,
              callback: function(value, index, values) {
                if (value >= 1.0) {
                  return value.toFixed(0);
                } else {
                  let v = value.toFixed(2);
                  if (v.endsWith('.00')) {
                    return value.toFixed(0);
                  } else {
                    return v;
                  }
                }
              },
              color: '#666',
              display: true
            }
          }],
          yAxes: [
            {
              position: 'left',
              id: 'freq-axis',
              type: freqAxisType,
              gridLines: {
                color: '#666',
                display: true,
                drawBorder: true,
                drawOnChartArea: false
              },
              scaleLabel: {display: true, labelString: 'Frequency (Hz)'},

              // after a pan, the graph sometimes shows poorly
              // formatted values for min and max (apparently bypasses
              // the tick callback?)
              // Workaround by just not displaying them
              afterTickToLabelConversion: function(scaleInstance) {
                // set the first and last tick to null so it does not
                // display note, ticks[0] is the last tick and
                // ticks[length - 1] is the first
                scaleInstance.ticks[0] = null;
                scaleInstance.ticks[scaleInstance.ticks.length - 1] = null;

                // need to do the same thing for this similiar array
                // which is used internally
                // scaleInstance.ticksAsNumbers[0] = null;
                // scaleInstance.ticksAsNumbers[scaleInstance.ticksAsNumbers.length
                // - 1] = null;
              },

              ticks: {
                precision: 2,
                callback: function(value, index, values) {
                  // don't use scientific notation
                  if (value >= 1.0) {
                    return value.toFixed(0);
                  } else {
                    let v = value.toFixed(2);
                    if (v.endsWith('.00')) {
                      return value.toFixed(0);
                    } else {
                      return v;
                    }
                  }
                },
                color: '#eee',
                display: true
              }
            },
            {
              position: 'right',
              id: 'amp-axis',
              type: 'linear',
              gridLines: {
                color: '#666',
                display: true,
                drawBorder: true,
                drawOnChartArea: false
              },
              scaleLabel: {display: true, labelString: 'Amplitude dB'},

              // after a pan, the graph sometimes shows poorly
              // formatted values for min and max (apparently bypasses
              // the tick callback?)
              // Workaround by just not displaying them
              afterTickToLabelConversion: function(scaleInstance) {
                // set the first and last tick to null so it does not
                // display note, ticks[0] is the last tick and
                // ticks[length - 1] is the first
                scaleInstance.ticks[0] = null;
                scaleInstance.ticks[scaleInstance.ticks.length - 1] = null;

                // need to do the same thing for this similiar array
                // which is used internally
                // scaleInstance.ticksAsNumbers[0] = null;
                // scaleInstance.ticksAsNumbers[scaleInstance.ticksAsNumbers.length
                // - 1] = null;
              },

              ticks: {color: '#eee', display: true}
            }
          ],
        },
        responsive: true,
        maintainAspectRatio: false,

        dragData: viewVCE_voice.voicingMode,
        dragDataRound: 0,
        dragX: true,

        dragOptions: {showTooltip: true},

        onDragStart: function(e, element) {
          viewVCE_envs.dragging = true;
          // console.log('onDragStart: ', envNum, e, element)
          //  constrain amp curve:  first point fixed at 0,0
          //
          //  HACK: also if only a single point in the dataset (common
          //  for freq envs), don't allow drag.
          // This works around an as yet undiagnosed bug that cases
          // vlaues to go to floating point numbers smaller than zero
          // and cofuse the auto-scaling.
          if (viewVCE_envs.chart.data.datasets[element._datasetIndex]
                      .data.length == 1 ||
              ((envNum < 0 && element._datasetIndex >= 2) || (envNum >= 2)) &&
                  element._index === 0) {
            // can't move the first amp point
            dragOldValue.x =
                viewVCE_envs.chart.data.datasets[element._datasetIndex]
                    .data[element._index]
                    .x;
            dragOldValue.y =
                viewVCE_envs.chart.data.datasets[element._datasetIndex]
                    .data[element._index]
                    .y;
            // console.log("ondragStart: freeze 0th
            // amp",element._datasetIndex,element._index,dragOldValue)
            viewVCE_envs.chart.update(0);
          }
        },
        onDrag: function(e, datasetIndex, index, value) {
          if (viewVCE_envs.chart.data.datasets[datasetIndex].data.length == 1 ||
              ((envNum < 0 && datasetIndex >= 2) || (envNum >= 2)) &&
                  index === 0) {
            // can't move the first amp point
            viewVCE_envs.chart.data.datasets[datasetIndex].data[index].x =
                dragOldValue.x;
            viewVCE_envs.chart.data.datasets[datasetIndex].data[index].y =
                dragOldValue.y;
            // console.log("ondrag: freeze 0th
            // amp",datasetIndex,index,dragOldValue)
            viewVCE_envs.chart.update(0);
            return
          }
          e.target.style.cursor = 'grabbing'
          // time must stay between neighboring points:
          let min = index > 0 ?
              viewVCE_envs.chart.data.datasets[datasetIndex].data[index - 1].x :
              0;
          // if the last point, use the scale max
          let max =
              (index ===
               (viewVCE_envs.chart.data.datasets[datasetIndex].data.length -
                1)) ?
              (viewVCE_envs.chart.scales['time-axis'].max + 1) :
              viewVCE_envs.chart.data.datasets[datasetIndex].data[index + 1].x;
          // if this is a freq env, then the 0th point's x value is
          // fixed at 0
          if (((envNum < 0 && datasetIndex < 2) || (envNum < 2)) &&
              index === 0) {
            min = -1;
            max = 2;  // the clamping expression below subtracts or
                      // adds 1
          }
          // console.log('ondrag: ', envNum, datasetIndex, index,
          // value, min, max)
          if (value.x >= max) {
            value.x = max - 1;
            // console.log('onDrag: CLAMP ', datasetIndex, index,
            // value)
            viewVCE_envs.chart.update();
          } else if (value.x <= min) {
            value.x = min + 1;
            // console.log('onDrag: CLAMP ', datasetIndex, index,
            // value)
            viewVCE_envs.chart.update();
          }
          viewVCE_envs.updateEnvFromGraphChange(
              datasetIndex, index, value, false)
        },
        onDragEnd: function(e, datasetIndex, index, value) {
          viewVCE_envs.dragging = false;
          e.target.style.cursor = 'default'
          // console.log('onDragEnd: ', datasetIndex, index, value)

          viewVCE_envs.updateEnvFromGraphChange(
              datasetIndex, index, value, true)
        },

        hover: {
          mode: 'index',
          intersect: true,
          onHover: function(e) {
            const point = this.getElementAtEvent(e)
            if (viewVCE_voice.voicingMode && point.length &&
                !(((envNum < 0 && point[0]._datasetIndex >= 2) ||
                   (envNum >= 2)) &&
                  point[0]._index === 0)) {
              e.target.style.cursor = 'grab';
            }
            else {
              e.target.style.cursor = 'default';
            }
          }
        },
        plugins: {
          // zoom plugin is only used by the env graphs
          zoom: {
            zoom: {enabled: false},
            pan: {
              enabled: true,
              mode: function({chart}) {
                if (viewVCE_envs.dragging) {
                  return '';
                }
                return 'xy';
              },
            }
          }
        }
      }
    });
    document.getElementById('tabTelltaleContent').value = `osc:${oscNum}`;
    viewVCE_envs.supressOnchange = false;
  }

  ,

  updateEnvFromGraphChange: function(datasetIndex, index, value, fireOnChange) {
    // now reverse engineer the changed env values and update the
    // corresponding point value or time if x has changed, then the TIME
    // value for both the point and the preceding point need to change
    // (since the env values are the delta-t from the previous point, not
    // the absolute t of the point) if y has changed, only its value needs
    // to be updated.
    let newV = value.y
    let newT
    let nextNewT = undefined
    let fieldIndex = index + 1;  // fields are 1-based

    if (datasetIndex >= 2) {
      // amp.  the first point in the env corresponds to the second point
      // on the graph
      fieldIndex = fieldIndex - 1;
      if (index === 1) {
        newT = value.x
      } else {
        newT =
            (value.x -
             viewVCE_envs.chart.data.datasets[datasetIndex].data[index - 1].x);
      }
    } else {
      // freq.  the first point in the env corresponds to the first point
      // on the graph
      if (index === 0) {
        newT = undefined;
      } else {
        newT =
            (value.x -
             viewVCE_envs.chart.data.datasets[datasetIndex].data[index - 1].x);
      }
    }
    // last point case is common to both types of env
    // if last point, there's no nextT
    if (index !=
        viewVCE_envs.chart.data.datasets[datasetIndex].data.length - 1) {
      nextNewT =
          (viewVCE_envs.chart.data.datasets[datasetIndex].data[index + 1].x -
           value.x);
    }
    console.log('UPDATE VALUES', fieldIndex, newV, newT, nextNewT)

    function setValueAndFireOnchange(id, val) {
      let ele = document.getElementById(id);
      ele.value = val;
      // don't run onchange during the drag - since we redraw the graph
      // after sending data to the Synergy
      //(and that aborts the drag)
      if (fireOnChange) {
        // just call the function directly; faking the event in the
        // browser is error prone
        viewVCE_envs.onchange(ele);
      }
    }

    setValueAndFireOnchange(
        `${viewVCE_envs.valFieldNameByDatasetIdx[datasetIndex]}[${fieldIndex}]`,
        newV);
    if (newT != undefined) {
      setValueAndFireOnchange(
          `${viewVCE_envs.timeFieldNameByDatasetIdx[datasetIndex]}[${
              fieldIndex}]`,
          newT);
    }
    if (nextNewT != undefined) {
      setValueAndFireOnchange(
          `${viewVCE_envs.timeFieldNameByDatasetIdx[datasetIndex]}[${
              fieldIndex + 1}]`,
          nextNewT);
    }
  }
};
