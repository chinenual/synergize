const { lookupService } = require("dns");
const { env } = require("process");
const { DH_CHECK_P_NOT_PRIME } = require("constants");

var dragOldValue = {x: undefined, y: undefined};

let viewVCE_envs = {

	chart: null,

	init: function (incrementalUpdate) {
		//console.log('--- start viewVCE_envs init');

		if (viewVCE_envs.deb_onchange == null) {
			viewVCE_envs.deb_onchange = index.debounceFirstArg(viewVCE_envs.raw_onchange, DEBOUNCE_WAIT);
		}
		if (viewVCE_envs.deb_onchangeEnvAccel == null) {
			viewVCE_envs.deb_onchangeEnvAccel = _.debounce(viewVCE_envs.raw_onchangeEnvAccel, DEBOUNCE_WAIT);
		}
		if (viewVCE_envs.deb_copyFrom == null) {
			viewVCE_envs.deb_copyFrom = _.debounce(viewVCE_envs.raw_copyFrom, DEBOUNCE_WAIT);
		}

		var selectEle = document.getElementById("envOscSelect");
		// remove old options:
		while (selectEle.firstChild) {
			selectEle.removeChild(selectEle.firstChild);
		}

		for (i = 0; i <= vce.Head.VOITAB; i++) {
			var option = document.createElement("option");
			option.value = "" + (i + 1);
			option.innerHTML = "" + (i + 1);
			selectEle.appendChild(option);
		}
		$('#envCopySelectDiv').hide();

		viewVCE_envs.envChartUpdate(1, -1, true)
		//console.log('--- finish viewVCE_envs init');
	},


	// Freq values:
	//   as displayed: -61 .. 63
	//   byte range:   0xc3 .. 0x3f (-61 .. 63)
	scaleFreqEnvValue: function (v) {
		// See OSCDSP.Z80 DISVAL: DVAL10:
		return v;
	},
	unscaleFreqEnvValue: function (v) {
		return v;
	},
	FreqEnvValueToText(v) {
		return v;
	},
	TextToFreqEnvValue(v) {
		return v;
	},

	// Amp values:
	//   as displayed: 0 .. 72
	//   byte range:   0x37 .. 0x7f (55 .. 127)
	scaleAmpEnvValue: function (v) {
		// See OSCDSP.Z80 DISVAL: DVAL30:
		//if (last) return 0;
		return Math.max(0, v - 55);
	},
	unscaleAmpEnvValue: function (v) {
		return v + 55;
	},
	AmpEnvValueToText(v) {
		//console.log("AmpEnvValueToText '" + v + "' --> "+viewVCE_envs.scaleAmpEnvValue(v));
		return '' + viewVCE_envs.scaleAmpEnvValue(v);
	},
	TextToAmpEnvValue(v) {
		if (v == null || v === '') {
			//console.log("TextToAmpEnvValue '" + v + "' --> 55");
			return 55;
		}
		var val = parseInt(v, 10);
		//console.log("TextToAmpEnvValue '" + v + "' -> "+viewVCE_envs.unscaleAmpEnvValue(val));
		return '' + viewVCE_envs.unscaleAmpEnvValue(val);
	},

	// NOTE: the ftab based scaling functions as done in SYNHCS are not reversable (for the freq case, 
	// several values of x map to the same y, so reversing y can never map to some x's.  In SYNHCS, 
	// this didnt matter since the mapping was one-way (the raw x values go to the synergy, the y values 
	// were only used to show the values to the user. For us, we need to convert the "user y" values to 
	// the "x" values to send to the synergy.)
	//
	// Problem is the exponential nature of the scaling would lead to HUGE numbers.  (my guess is that 
	// few real patches use these large time values.)  In any case, I've chosen to just create a table 
	// based mapping approach rather than do all the math that SYNHCS does. This allows me to substitute 
	// some unique values at the end of the range to keep things unique, but not allow them to
	// get outrageously large.
	//
	// For most values, it's exactly the same as SYNHCS, but for those extra values of x, there are 
	// new y's so the editor can do its job

	freqTimeScale: [0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		25, 28, 32, 36, 40, 45, 51, 57,
		64, 72, 81, 91, 102, 115, 129, 145,
		163, 183, 205, 230, 258, 290, 326, 366,
		411, 461, 517, 581, 652, 732, 822, 922,
		1035, 1162, 1304, 1464, 1644, 1845, 2071, 2325,
		2609, 2929, 3288, 3691, 4143, 4650, 5219, 5859,
		6576, 7382, 8286, 9300, 10439, 11718, 13153, 14764,
		16572, 18600, 20078, 23436, 26306, 29528, 29529, 29530,
		29531, 29532, 29533, 29534, 29535],
	ampTimeScale: [0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39,
		40, 45, 51, 57, 64, 72, 81, 91,
		102, 115, 129, 145, 163, 183, 205, 230,
		258, 290, 326, 366, 411, 461, 517, 581,
		652, 732, 822, 922, 1035, 1162, 1304, 1464,
		1644, 1845, 2071, 2325, 2609, 2929, 3288, 3691,
		4143, 4650, 5219, 5859, 6576],


	// Freq Time values:
	//   as displayed: 0 .. 29528
	//   byte range:   0x0 .. 0x54 (0 .. 84)
	scaleFreqTimeValue: function (v) {
		// See OSCDSP.Z80 DISVAL for the original ftab-baased scaling which is roughly:
		//	if (v <= 15) return v;
		//  return viewVCE_envs.scaleViaRtab((2 * v) - 14);
		if (v < 0) {
			return 0;
		} else if (v >= viewVCE_envs.freqTimeScale.length) {
			return viewVCE_envs.freqTimeScale[viewVCE_envs.freqTimeScale.length - 1];
		}
		return viewVCE_envs.freqTimeScale[v];
	},
	unscaleFreqTimeValue: function (v) {
		// fixme: linear search is brute force - but the list is short - performance is "ok" as is...
		for (var i = 0; i < viewVCE_envs.freqTimeScale.length; i++) {
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
		var val = parseInt(v, 10);
		return '' + viewVCE_envs.unscaleFreqTimeValue(val);
	},


	// Freq Time values:
	//   as displayed: 0 .. 6576
	//   byte range:   0x0 .. 0x54 (0 .. 84)
	scaleAmpTimeValue: function (v) {
		// See OSCDSP.Z80 DISVAL: DVAL20: which is:
		// if (v < 39) return v;
		// return viewVCE_envs.scaleViaRtab((v * 2) - 54);
		//console.log("scale amp time value: " + v + ", -> " + viewVCE_envs.ampTimeScale[v])
		if (v < 0) {
			return 0;
		} else if (v >= viewVCE_envs.ampTimeScale.length) {
			return viewVCE_envs.ampTimeScale[viewVCE_envs.ampTimeScale.length - 1];
		}
		return viewVCE_envs.ampTimeScale[v];

	},
	unscaleAmpTimeValue: function (v) {
		// fixme: linear search is brute force - but the list is short - performance is "ok" as is...
		for (var i = 0; i < viewVCE_envs.ampTimeScale.length; i++) {
			if (viewVCE_envs.ampTimeScale[i] >= v) {
				//console.log("unscale amp time value: " + v + ", -> " + i + " (" + viewVCE_envs.ampTimeScale[i])
				return i;
			}
		}
		// shouldnt happen!
		//console.log("unscale amp time value fall through " + v + " " + typeof (v))
		return viewVCE_envs.ampTimeScale.length - 1;
	},
	AmpTimeValueToText(v) {
		return '' + viewVCE_envs.scaleAmpTimeValue(v);
	},
	TextToAmpTimeValue(v) {
		if (v == null || v === '') {
			return 0;
		}
		var val = parseInt(v, 10);
		return '' + viewVCE_envs.unscaleAmpTimeValue(val);
	},

	testConversionFunctions: function () {
		var ok = true;
		for (var i = -61; i <= 63; i++) {
			var scaled = viewVCE_envs.scaleFreqEnvValue(i);
			var unscaled = viewVCE_envs.unscaleFreqEnvValue(scaled);
			if (i != unscaled) {
				ok = false;
				console.log("ERROR: freqEnvValue " + i + " totext: " + scaled + " reversed to " + unscaled)
			}
		}
		for (var i = 55; i <= 127; i++) {
			var scaled = viewVCE_envs.scaleAmpEnvValue(i);
			var unscaled = viewVCE_envs.unscaleAmpEnvValue(scaled);
			if (i != unscaled) {
				ok = false;
				console.log("ERROR: ampEnvValue " + i + " totext: " + scaled + " reversed to " + unscaled)
			}
		}
		for (var i = 0; i <= 79; i++) {
			var scaled = viewVCE_envs.scaleFreqTimeValue(i);
			var unscaled = viewVCE_envs.unscaleFreqTimeValue(scaled);
			if (i != unscaled) {
				ok = false;
				console.log("ERROR: ampTimeValue " + i + " totext: " + scaled + " reversed to " + unscaled)
			}
		}
		for (var i = 0; i <= 79; i++) {
			var scaled = viewVCE_envs.scaleAmpTimeValue(i);
			var unscaled = viewVCE_envs.unscaleAmpTimeValue(scaled);
			if (i != unscaled) {
				ok = false;
				console.log("ERROR: ampTimeValue " + i + " totext: " + scaled + " reversed to " + unscaled)
			}
		}
		// Spot check some values to ensure the forumlae are computing same values as SYNHCS did (except 
		// for the upper range of freq time which we delibrartely change to make the function reversable)
		var expects = [
			{
				arr: [[0, 0], [10, 10], [15, 15], [16, 25], [54, 2071], [75, 23436], [76, 26306], [77, 29528], [84, 29535], [85, 29535]],
				name: "freqTimeValue",
				func: viewVCE_envs.scaleFreqTimeValue,
			},
			{
				arr: [[0, 0], [20, 20], [40, 40], [41, 45], [54, 205], [75, 2325], [76, 2609], [83, 5859], [84, 6576], [85, 6576]],
				name: "ampTimeValue",
				func: viewVCE_envs.scaleAmpTimeValue,
			},
			{
				arr: [[-61, -61], [-15, -15], [0, 0], [63, 63]],
				name: "freqValue",
				func: viewVCE_envs.scaleFreqEnvValue,
			},
			{
				arr: [[55, 0], [56, 1], [126, 71], [127, 72]],
				name: "ampValue",
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


		for (var j = 0; j < expects.length; j++) {
			var expect = expects[j];
			for (var i = 0; i < expect.arr.length; i++) {
				var scaled = expect.func(expect.arr[i][0]);
				if (scaled != expect.arr[i][1]) {
					ok = false;
					console.log("ERROR: " + expect.name + "(" + expect.arr[i][0] + ") == " + scaled + ", expected " + expect.arr[i][1]);
				}
			}
		}
		console.log("viewVCE_envs.testConversionFunctions: " + (ok ? "PASS" : "FAIL"));
		return ok;
	},

	supressOnChange: false,

	onchangeLoop: function (ele) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		if (viewVCE_envs.supressOnchange) { console.log("viewVCE_envs.suppressOnChange"); return; }

		var eleIndex;
		var envOscSelectEle = document.getElementById("envOscSelect");
		var osc = parseInt(envOscSelectEle.value, 10); // one-based osc index
		var envEnvSelectEle = document.getElementById("envEnvSelect");
		var selectedEnv = parseInt(envEnvSelectEle.value, 10);
		var eleValue = ele.value;

		var pattern = /([A-Za-z]+)\[(\d+)\]/;
		if (ret = ele.id.match(pattern)) {
			eleIndex = parseInt(ret[2])
		}

		var env;
		var envid;
		if (ele.id.includes('Freq')) {
			env = vce.Envelopes[osc - 1].FreqEnvelope;
			envid = "Freq";
		} else {
			env = vce.Envelopes[osc - 1].AmpEnvelope;
			envid = "Amp";
		}

		var accelLow = 30; // defaults
		var accelUp = 30; // defaults

		if (env.ENVTYPE == 1) {
			accelLow = env.SUSTAINPT;
			accelUp = env.LOOPPT;

			// temporarily replace the accellerations with point indexes to make all the validations below less messy than they might be
			env.SUSTAINPT = 0;
			env.LOOPPT = 0;
		}

		// type = 1  : no loop (and LOOPPT and SUSTAINPT are accelleration rates not point positions)
		// type = 2  : S only
		// type = 3  : L and S - L must be before S
		// type = 4  : R and S - R must be before S
		// WARNING: when type1, the LOOPPT and SUSTAINPT values are _acceleration_ rates, not point positions. What a pain.
		console.log("loop change " + eleIndex + " " + eleValue + " " + envid);
		console.dir(env);
		if (eleValue == '') {
			// always safe to remove a loop point
			if (env.LOOPPT == eleIndex) {
				env.LOOPPT = 0;
				env.ENVTYPE = 2; // sustain only
			}
			// can't remove SUSTAIN POINT if there's a loop point
			if (env.SUSTAINPT == eleIndex) {
				if (env.ENVTYPE === 3 || env.ENVTYPE == 4) {
					index.errorNotification('Cannot remove SUSTAIN point if there are LOOP or REPEAT points')
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
				env.ENVTYPE = 1; // no loops
			} else if (env.LOOPPT === 0) {
				env.ENVTYPE = 2; // sustain-only 
			}
		} else if (eleValue === 'S') {
			if (env.LOOPPT > eleIndex) {
				index.errorNotification('Cannot set SUSTAIN point before LOOP or REPEAT point');
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
		} else { // 'R' or 'L'
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

		// validation is OK, but now need to clean up any selects (i.e. if loop point moved, need to set the previous location to '')
		// easiest thing to do is just brute force reset each ele to reflect the value in the envelope
		for (var p = 0; p < 16; p++) {
			$(`#env${envid}Loop\\[${p + 1}\\] option[value='']`).prop('selected', true);
			$(`#env${envid}Loop\\[${p + 1}\\] option[value='L']`).prop('selected', false);
			$(`#env${envid}Loop\\[${p + 1}\\] option[value='R']`).prop('selected', false);
			$(`#env${envid}Loop\\[${p + 1}\\] option[value='S']`).prop('selected', false);
		}
		if (env.ENVTYPE != 1 && env.SUSTAINPT > 0) {
			$(`#env${envid}Loop\\[${env.SUSTAINPT}\\] option[value='S']`).prop('selected', true);
		}
		if (env.ENVTYPE != 1 && env.LOOPPT > 0) {
			var v = env.ENVTYPE == 3 ? 'L' : 'R'
			$(`#env${envid}Loop\\[${env.LOOPPT}\\] option[value='${v}']`).prop('selected', true);
		}

		// only show accelleration values if type1 envelope
		if (env.ENVTYPE === 1) {
			$(`.type1accel div.${envid}`).show();
			$(`#accel${envid}Low`).val(env.SUSTAINPT);
			$(`#accel${envid}Up`).val(env.LOOPPT);
		} else {
			$(`.type1accel div.${envid}`).hide();
		}

		console.log("resulting env: " + envid);
		console.dir(env);
		let message = {
			"name": "setLoopPoint",
			"payload": {
				"Osc": osc,
				"Env": envid,
				"EnvType": env.ENVTYPE,
				"SustainPt": env.SUSTAINPT,
				"LoopPt": env.LOOPPT,
			}
		};
		astilectron.sendMessage(message, function (message) {
			//console.log("setLoopPoint returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the value
				index.errorNotification(message.payload);
				return false;
			} else {
				// currently there's no visual rendering of the loop, so no need to refresh the chart
				//				viewVCE_envs.envChartUpdate(osc, selectedEnv, false);
			}
		});
		return true;

	},

	copyFrom: function (fromOsc) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		if (viewVCE_envs.supressOnchange) { /*console.log("viewVCE_envs.suppressOnChange");*/ return; }
		viewVCE_envs.deb_copyFrom(fromOsc);
	},

	deb_copyFrom: null,

	raw_copyFrom: function (fromOsc) {
		var oscSelectEle = document.getElementById("envOscSelect");
		var toOsc = oscSelectEle.options[oscSelectEle.selectedIndex].value;
		toOsc = parseInt(toOsc, 10);

		// "copy" means copy all the osc-specific stuff related to the envelopes - but not the patch and detuning fields.

		// abuse JSON to do a deep copy:
		newEnvelopes = JSON.parse(JSON.stringify(vce.Envelopes[fromOsc - 1]))
		// retain the stuff we don't want copied: 
		newEnvelopes.FreqEnvelope.OPTCH = vce.Envelopes[toOsc - 1].FreqEnvelope.OPTCH;
		newEnvelopes.FreqEnvelope.OHARM = vce.Envelopes[toOsc - 1].FreqEnvelope.OHARM;
		newEnvelopes.FreqEnvelope.FDETUN = vce.Envelopes[toOsc - 1].FreqEnvelope.FDETUN;

		let message = {
			"name": "setEnvelopes",
			"payload": {
				"Osc": toOsc,
				"Envelopes": newEnvelopes
			}
		};
		index.spinnerOn();
		astilectron.sendMessage(message, function (message) {
			index.spinnerOff();
			//console.log("setEnvelopes returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the value
				index.errorNotification(message.payload);
				return false;
			} else {
				vce.Envelopes[toOsc - 1] = newEnvelopes
				viewVCE_envs.envChartUpdate(toOsc, -1, true);
			}
		});
		return true;
	},

	onchangeEnvAccel: function (ele) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		if (viewVCE_envs.supressOnchange) { /*console.log("viewVCE_envs.suppressOnChange");*/ return; }
		viewVCE_envs.deb_onchangeEnvAccel(ele);
	},

	deb_onchangeEnvAccel: null,

	raw_onchangeEnvAccel: function (ele) {
		if (viewVCE.supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/ return; }
		if (viewVCE_envs.supressOnchange) { /*console.log("raw viewVCE_envs.suppressOnChange");*/ return; }

		// type1 accelerations are really just the SUSTAIN and LOOP points.  We use the same backend function as the loop change event


		var envOscSelectEle = document.getElementById("envOscSelect");
		var osc = parseInt(envOscSelectEle.value, 10); // one-based osc index
		var envEnvSelectEle = document.getElementById("envEnvSelect");
		var selectedEnv = parseInt(envEnvSelectEle.value, 10);
		var eleValue = index.checkInputElementValue(ele);
		if (eleValue == undefined) {
			return;
		}

		//console.log("changed: " + ele.id + " val: " + ele.value);

		var env;
		var envid;
		// id is accelFreqUp , accelAmpLow etc.
		if (ele.id.includes('Freq')) {
			env = vce.Envelopes[osc - 1].FreqEnvelope;
			envid = "Freq";
		} else {
			env = vce.Envelopes[osc - 1].AmpEnvelope;
			envid = "Amp";
		}
		if (ele.id.includes('Low')) {
			env.SUSTAINPT = parseInt(eleValue, 10);
		} else {
			env.LOOPPT = parseInt(eleValue, 10);
		}

		console.log("ACCEL change " + ele.id + " " + eleValue + " " + envid);


		let message = {
			"name": "setLoopPoint",
			"payload": {
				"Osc": osc,
				"Env": envid,
				"EnvType": env.ENVTYPE,
				"SustainPt": env.SUSTAINPT,
				"LoopPt": env.LOOPPT,
			}
		};
		astilectron.sendMessage(message, function (message) {
			//console.log("setLoopPoint returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the value
				index.errorNotification(message.payload);
				return false;
			} else {
				// currently there's no visual rendering of the loop, so no need to refresh the chart
				//				viewVCE_envs.envChartUpdate(osc, selectedEnv, false);
				viewVCE_voice.sendToCSurface(ele, ele.id, parseInt(eleValue, 10));
			}
		});
		return true;

	},

	onchange: function (ele) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		if (viewVCE_envs.supressOnchange) { /*console.log("viewVCE_envs.suppressOnChange");*/return; }
		viewVCE_envs.deb_onchange(ele);
	},

	deb_onchange: null, // initialized during init()

	raw_onchange: function (ele) {
		if (viewVCE.supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/ return; }
		if (viewVCE_envs.supressOnchange) { /*console.log("raw viewVCE_envs.suppressOnChange");*/ return; }

		var eleIndex;
		var envOscSelectEle = document.getElementById("envOscSelect");
		var osc = parseInt(envOscSelectEle.value, 10); // one-based osc index
		var envEnvSelectEle = document.getElementById("envEnvSelect");
		var selectedEnv = parseInt(envEnvSelectEle.value, 10);

		// Don't call checkInoutElementValue() - it assumes that there is no scaling 
		// and would apply the "byte" min/max to the "text" scaled value
		//	  var value = index.checkInputElementValue(ele);
		var value = parseInt(ele.value, 10);;
		if (value == undefined) {
			return;
		}
		//console.log("in onchange - value: " + value + " " + typeof(value))

		var pattern = /([A-Za-z]+)\[(\d+)\]/;
		var funcName;
		var eleIndex;
		if (ret = ele.id.match(pattern)) {
			fieldType = ret[1];
			funcName = 'set' + fieldType.charAt(0).toUpperCase() + fieldType.slice(1);
			eleIndex = parseInt(ret[2])
			var bytevalue;
			// now scale the value to the byte value the synergy wants to see:
			switch (fieldType) {
				case "envFreqLowVal":
					bytevalue = viewVCE_envs.unscaleFreqEnvValue(value);
					vce.Envelopes[osc - 1].FreqEnvelope.Table[((eleIndex - 1) * 4) + 0] = bytevalue;
					break;
				case "envFreqUpVal":
					bytevalue = viewVCE_envs.unscaleFreqEnvValue(value);
					vce.Envelopes[osc - 1].FreqEnvelope.Table[((eleIndex - 1) * 4) + 1] = bytevalue;
					break;
				case "envFreqLowTime":
					bytevalue = viewVCE_envs.unscaleFreqTimeValue(value);
					vce.Envelopes[osc - 1].FreqEnvelope.Table[((eleIndex - 1) * 4) + 2] = bytevalue;
					break;
				case "envFreqUpTime":
					bytevalue = viewVCE_envs.unscaleFreqTimeValue(value);
					vce.Envelopes[osc - 1].FreqEnvelope.Table[((eleIndex - 1) * 4) + 3] = bytevalue;
					break;
				case "envAmpLowVal":
					bytevalue = viewVCE_envs.unscaleAmpEnvValue(value);
					vce.Envelopes[osc - 1].AmpEnvelope.Table[((eleIndex - 1) * 4) + 0] = bytevalue;
					break;
				case "envAmpUpVal":
					bytevalue = viewVCE_envs.unscaleAmpEnvValue(value);
					vce.Envelopes[osc - 1].AmpEnvelope.Table[((eleIndex - 1) * 4) + 1] = bytevalue;
					break;
				case "envAmpLowTime":
					bytevalue = viewVCE_envs.unscaleAmpTimeValue(value);
					vce.Envelopes[osc - 1].AmpEnvelope.Table[((eleIndex - 1) * 4) + 2] = bytevalue;
					break;
				case "envAmpUpTime":
					bytevalue = viewVCE_envs.unscaleAmpTimeValue(value);
					vce.Envelopes[osc - 1].AmpEnvelope.Table[((eleIndex - 1) * 4) + 3] = bytevalue;
					break;
			}
		}
		console.log("env ele change " + ele.id + " rawval: " + ele.value + " -> " + value + " -> " + bytevalue);
		//console.log("in onchange - bytevalue: " + bytevalue + " " + typeof(bytevalue))

		let message = {
			"name": funcName,
			"payload": {
				"Osc": osc,
				"Index": eleIndex,
				"Value": bytevalue
			}
		};
		astilectron.sendMessage(message, function (message) {
			//console.log(funcName + " returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the value
				index.errorNotification(message.payload);
				return false;
			} else {
				viewVCE_envs.envChartUpdate(osc, selectedEnv, false);
				viewVCE_voice.sendToCSurface(ele, ele.id, bytevalue);
			}
		});
		return true;

	},

	changeEnvPoints: function (whichEnv, increment) {
		var eleIndex;
		var envOscSelectEle = document.getElementById("envOscSelect");
		var osc = parseInt(envOscSelectEle.value, 10); // one-based osc index
		var envEnvSelectEle = document.getElementById("envEnvSelect");
		var selectedEnv = parseInt(envEnvSelectEle.value, 10);
		var envs = vce.Envelopes[osc - 1];

		console.log("#points changed: " + whichEnv + " increment: " + increment);

		var changed = false;
		if (whichEnv === 'freq') {
			newlen = envs.FreqEnvelope.NPOINTS + increment;
			if (newlen >= 1 && newlen <= 16) {
				if (increment === -1 && envs.FreqEnvelope.ENVTYPE != 1 && (envs.FreqEnvelope.LOOPPT > newlen || envs.FreqEnvelope.SUSTAINPT > newlen)) {
					index.errorNotification("Cannot remove envelope point with SUSTAIN/LOOP marker.  Remove the SUSTAIN/LOOP marker before trying to remove points.");
					return;
				}
				envs.FreqEnvelope.NPOINTS = newlen;
				changed = true;
			}
		} else {
			newlen = envs.AmpEnvelope.NPOINTS + increment;
			if (newlen >= 1 && newlen <= 16) {
				if (increment === -1 && envs.AmpEnvelope.ENVTYPE != 1 && (envs.AmpEnvelope.LOOPPT > newlen || envs.AmpEnvelope.SUSTAINPT > newlen)) {
					index.errorNotification("Cannot remove envelope point with SUSTAIN/LOOP marker.  Remove the SUSTAIN/LOOP marker before trying to remove points.");
					return;
				}
				envs.AmpEnvelope.NPOINTS = newlen;
				changed = true;
			}
		}
		if (changed) {
			let message = {
				"name": "setOscEnvLengths",
				"payload": {
					"Osc": osc,
					"FreqLength": envs.FreqEnvelope.NPOINTS,
					"AmpLength": envs.AmpEnvelope.NPOINTS
				}
			};
			astilectron.sendMessage(message, function (message) {
				//console.log("setOscEnvLengths returned: " + JSON.stringify(message));
				// Check error
				if (message.name === "error") {
					// failed - dont change the value
					index.errorNotification(message.payload);
					return false;
				} else {
					viewVCE_envs.envChartUpdate(osc, selectedEnv, true);
				}
			});
		}
	},

	uncompressEnvelopes: function () {
		// the first time we evaluate this vce, the envelopes may be compressed.  To make it easier to add/remove 
		// filters in the editor, we rewrite the envelopes arrays such each has the max amount of elements and each are initialized
		// as SYNHCS does.
		if (vce.Extra["uncompressedEnvelopes"] != undefined) {
			// no need to do it again
			return;
		}
		// only need to worry about the number of oscillators in the Envelopes table; any addition osc's added will automatically
		// fill in "full length" envelopes (but use the length of the array not the current value of VOITAB lowered the number of osc's)
		for (i = 0; i < vce.Envelopes.length; i++) {
			const FULL_LENGTH = 16 * 4; // 16 rows, each with 4 values
			for (j = vce.Envelopes[i].FreqEnvelope.Table.length; j < FULL_LENGTH; j++) {
				vce.Envelopes[i].FreqEnvelope.Table.push(0);
			}
			for (j = vce.Envelopes[i].AmpEnvelope.Table.length; j < FULL_LENGTH; j++) {
				vce.Envelopes[i].AmpEnvelope.Table.push(0);
			}
		}
		vce.Extra.uncompressedEnvelopes = true;
	},

	changeTimeScale: function(val) {
		this.chart.options.scales.xAxes [0].type = val;
		this.chart.update();
	},

	changeFreqScale: function(val) {
		this.chart.options.scales.yAxes [0].type = val;
		this.chart.update();
	},

	changeTimeZoom: function(val) {
		var div = document.getElementById('envZoomDiv');
		div.style.width = val;
		//this.chart.update();
	},

	// XREF: needs to match the dataset order inside envChangeUpdate()
	valFieldNameByDatasetIdx : ["envFreqLowVal", "envFreqUpVal", "envAmpLowVal", "envAmpUpVal"],
	timeFieldNameByDatasetIdx : ["envFreqLowTime", "envFreqUpTime", "envAmpLowTime", "envAmpUpTime"],

	envChartUpdate: function (oscNum, envNum, animate) {
		viewVCE_envs.supressOnchange = true;

		viewVCE_envs.uncompressEnvelopes();

		var envCopySelectEle = document.getElementById("envCopySelect");
		// remove old options:
		while (envCopySelectEle.firstChild) {
			envCopySelectEle.removeChild(envCopySelectEle.firstChild);
		}
		// hide the copy selector for All or cases where there are no filters, or when we're not in voicing mode
		$('#envCopySelectDiv').hide();
		if (viewVCE_voice.voicingMode) {
			$('#envCopySelectDiv').show();
			// populate options in the select with only "other" osc (i.e. "this" osc should be not shown or at least unselectable)

			// first element is empty to avoid confusing the user if they havent selected something:
			var option = document.createElement("option");
			option.value = -1;
			option.innerHTML = "";
			envCopySelectEle.appendChild(option);

			for (i = 0; i < vce.Head.VOITAB; i++) {
				if ((i + 1) != oscNum) {
					var option = document.createElement("option");
					option.value = i + 1;
					option.innerHTML = i + 1;
					envCopySelectEle.appendChild(option);
				}
			}
		}

		var oscIndex = oscNum - 1;
		var envelopes = vce.Envelopes[oscIndex];

		var pointStyleMetadata = [
			// order needs to match the dataset array
			{color: 2,
				loopPt : -1,
				repeatPt : -1,
				sustainPt : -1
			},
			{color: 3,
				loopPt : -1,
				repeatPt : -1,
				sustainPt : -1
			},
			{color: 0,
				loopPt : -1,
				repeatPt : -1,
				sustainPt : -1
			},
			{color: 1,
				loopPt : -1,
				repeatPt : -1,
				sustainPt : -1
			},
		];
		function annotatePointStyle(ctx) {
			var styleMeta = filteredPointStyleMetadata[ctx.datasetIndex]
			console.log("annotate point ctx ",ctx)
			var img = new Image(14,14);
			if (styleMeta.sustainPt == ctx.dataIndex) {
				img.src  = `static/images/loopS-${styleMeta.color}.png`;
				return img;
			} else if (styleMeta.loopPt == ctx.dataIndex) {
				img.src  = `static/images/loopL-${styleMeta.color}.png`;
				return img;
			} else if (styleMeta.repeatPt == ctx.dataIndex) {
				img.src  = `static/images/loopR-${styleMeta.color}.png`;
				return img;
			}
			return 'circle'
		};

		// XREF: order needs to match the valFieldNameByDatasetIdx and timeFieldNameByDatasetIdx above
		let freqLowIdx = 0;
		let freqUpIdx = 1;
		let ampLowIdx = 2;
		let ampUpIdx = 3;
		let datasets = [
			{
				label: "Freq Low",
				xAxisID: "time-axis",
				yAxisID: "freq-axis",
				fill: false,
				lineTension: 0,
				pointRadius: 2,
				pointHitRadius: 5,
				pointStyle: annotatePointStyle,
				showLine: true,
				borderWidth: 3,
				backgroundColor: chartColors[2],
				borderColor: chartColors[2],
				data: []
			},
			{
				label: "Freq Up",
				xAxisID: "time-axis",
				yAxisID: "freq-axis",
				fill: false,
				lineTension: 0,
				pointRadius: 2,
				pointHitRadius: 5,
				pointStyle: annotatePointStyle,
				showLine: true,
				borderWidth: 3,
				backgroundColor: chartColors[3],
				borderColor: chartColors[3],
				data: []
			},
			{
				label: "Amp Low",
				xAxisID: "time-axis",
				yAxisID: "amp-axis",
				fill: false,
				lineTension: 0,
				pointRadius: 2,
				pointHitRadius: 5,
				pointStyle: annotatePointStyle,
				showLine: true,
				borderWidth: 3,
				backgroundColor: chartColors[0],
				borderColor: chartColors[0],
				data: []
			},
			{
				label: "Amp Up",
				xAxisID: "time-axis",
				yAxisID: "amp-axis",
				fill: false,
				lineTension: 0,
				pointRadius: 2,
				pointHitRadius: 5,
				pointStyle: annotatePointStyle,
				showLine: true,
				borderWidth: 3,
				backgroundColor: chartColors[1],
				borderColor: chartColors[1],
				data: []
			},
		];

		viewVCE_voice.sendToCSurface(null, `num-freq-env-points`, envelopes.FreqEnvelope.NPOINTS);
		viewVCE_voice.sendToCSurface(null, `num-amp-env-points`, envelopes.AmpEnvelope.NPOINTS);

		// clear old values:
		$('#envTable td.val input').val('');
		$('#envTable td.total span').html('');
		// clear the loop points
		$(`#envTable select option[value='']`).prop('selected', true);
		$(`#envTable select option[value='L']`).prop('selected', false);
		$(`#envTable select option[value='S']`).prop('selected', false);
		$(`#envTable select option[value='R']`).prop('selected', false);

		// fill in freq env data:

		// scaling algorithms derived from DISVAL: in OSCDSP.Z80

		var totalTimeLow = 0;
		var totalTimeUp = 0;
		var lastFreqLow = 0;
		var lastFreqUp = 0;
		var lastAmpLow = 0;
		var lastAmpUp = 0;

		for (i = 0; i < 16; i++) {
			// completely hide the rows for rows not used by either envelope
			var tr = $('#envTable tbody tr:eq(' + i + ')');
			if (i < Math.max(envelopes.FreqEnvelope.NPOINTS, envelopes.AmpEnvelope.NPOINTS)) {
				tr.show();
			} else {
				tr.hide();
			}
		}
		if (viewVCE_voice.voicingMode) {
			$(".listplusminus div").show();
		} else {
			$(".listplusminus div").hide();
		}
		viewVCE_voice.sendToCSurface(null, `num-freq-env-points`, envelopes.FreqEnvelope.NPOINTS);
		viewVCE_voice.sendToCSurface(null, `num-amp-env-points`, envelopes.AmpEnvelope.NPOINTS);

		// only show accelleration values if type1 envelope
		if (envelopes.FreqEnvelope.ENVTYPE === 1) {
			$('.type1accel div.Freq').show();
			$('#accelFreqLow').val(envelopes.FreqEnvelope.SUSTAINPT);
			$('#accelFreqUp').val(envelopes.FreqEnvelope.LOOPPT);
			viewVCE_voice.sendToCSurface(null, `freq-env-accel-visible`, 1);
			if (animate) {
				viewVCE_voice.sendToCSurface(null, `accelFreqLow`, envelopes.FreqEnvelope.SUSTAINPT);
				viewVCE_voice.sendToCSurface(null, `accelFreqUp`, envelopes.FreqEnvelope.LOOPPT);
			}
		} else {
			$('.type1accel div.Freq').hide();
			viewVCE_voice.sendToCSurface(null, `freq-env-accel-visible`, 0);
			if (animate) {
				viewVCE_voice.sendToCSurface(null, `accelFreqLow`, 0);
				viewVCE_voice.sendToCSurface(null, `acceFreqUp`, 0);
			}
		}
		// only show accelleration values if type1 envelope
		if (envelopes.AmpEnvelope.ENVTYPE === 1) {
			$('.type1accel div.Amp').show();
			$('#accelAmpLow').val(envelopes.AmpEnvelope.SUSTAINPT);
			$('#accelAmpUp').val(envelopes.AmpEnvelope.LOOPPT);
			viewVCE_voice.sendToCSurface(null, `amp-env-accel-visible`, 1);
			if (animate) {
				viewVCE_voice.sendToCSurface(null, `accelAmpLow`, envelopes.AmpEnvelope.SUSTAINPT);
				viewVCE_voice.sendToCSurface(null, `accelAmpUp`, envelopes.AmpEnvelope.LOOPPT);
			}
		} else {
			$('.type1accel div.Amp').hide();
			viewVCE_voice.sendToCSurface(null, `amp-env-accel-visible`, 0);
			if (animate) {
				viewVCE_voice.sendToCSurface(null, `accelAmpLow`, 0);
				viewVCE_voice.sendToCSurface(null, `accelAmpUp`, 0);
			}
		}

		for (i = envelopes.FreqEnvelope.NPOINTS; i < 16; i++) {
			// hide unused rows
			var tr = $('#envTable tbody tr:eq(' + i + ')');

			$(`#envFreqLoop\\[${i + 1}\\]`).hide();
			$(`#envFreqLowVal\\[${i + 1}\\]`).hide();
			$(`#envFreqUpVal\\[${i + 1}\\]`).hide();
			$(`#envFreqLowTime\\[${i + 1}\\]`).hide();
			$(`#envFreqUpTime\\[${i + 1}\\]`).hide();
			if (animate) {
				viewVCE_voice.sendToCSurface(null, `envFreqLowVal[${i + 1}]`, 0);
				viewVCE_voice.sendToCSurface(null, `envFreqUpVal[${i + 1}]`, 0);
				viewVCE_voice.sendToCSurface(null, `envFreqLowTime[${i + 1}]`, 0);
				viewVCE_voice.sendToCSurface(null, `envFreqUpTime[${i + 1}]`, 0);
			}
		}
		for (i = 0; i < envelopes.FreqEnvelope.NPOINTS; i++) {
			var tr = $('#envTable tbody tr:eq(' + i + ')');

			$(`#envFreqLoop\\[${i + 1}\\]`).show();
			$(`#envFreqLowVal\\[${i + 1}\\]`).show();
			$(`#envFreqUpVal\\[${i + 1}\\]`).show();
			$(`#envFreqLowTime\\[${i + 1}\\]`).show();
			$(`#envFreqUpTime\\[${i + 1}\\]`).show();

			// table is logically in groups of 4
			var freqLow = viewVCE_envs.scaleFreqEnvValue(envelopes.FreqEnvelope.Table[i * 4 + 0]);
			var freqUp = viewVCE_envs.scaleFreqEnvValue(envelopes.FreqEnvelope.Table[i * 4 + 1]);
			var timeLow = viewVCE_envs.scaleFreqTimeValue(envelopes.FreqEnvelope.Table[i * 4 + 2], i == 0);
			var timeUp = viewVCE_envs.scaleFreqTimeValue(envelopes.FreqEnvelope.Table[i * 4 + 3], i == 0);

			if (animate) {
				viewVCE_voice.sendToCSurface(null, `envFreqLowVal[${i + 1}]`, envelopes.FreqEnvelope.Table[i * 4 + 0]);
				viewVCE_voice.sendToCSurface(null, `envFreqUpVal[${i + 1}]`, envelopes.FreqEnvelope.Table[i * 4 + 1]);
				viewVCE_voice.sendToCSurface(null, `envFreqLowTime[${i + 1}]`, envelopes.FreqEnvelope.Table[i * 4 + 2]);
				viewVCE_voice.sendToCSurface(null, `envFreqUpTime[${i + 1}]`, envelopes.FreqEnvelope.Table[i * 4 + 3]);
			}

			if (i == 0) {
				// first row's time values are fixed at zero (the entries in the table are used for wave/kprop markers)
				timeLow = 0;
				timeUp = 0;
			}
			lastFreqLow = freqLow;
			lastFreqUp = freqUp;
			totalTimeLow += timeLow;
			totalTimeUp += timeUp;

			datasets[freqLowIdx].data.push({ x: totalTimeLow, y: freqLow });
			datasets[freqUpIdx].data.push({ x: totalTimeUp, y: freqUp });

			document.getElementById(`envFreqLowVal[${i + 1}]`).value = freqLow;
			document.getElementById(`envFreqUpVal[${i + 1}]`).value = freqUp;
			if (i !== 0) {
				document.getElementById(`envFreqLowTime[${i + 1}]`).value = timeLow;
				document.getElementById(`envFreqUpTime[${i + 1}]`).value = timeUp;
			}
			document.getElementById(`envFreqTotLowTime[${i + 1}]`).innerHTML = totalTimeLow;
			document.getElementById(`envFreqTotUpTime[${i + 1}]`).innerHTML = totalTimeUp;

			if (envelopes.FreqEnvelope.ENVTYPE != 1) {
				if (envelopes.FreqEnvelope.SUSTAINPT == (i + 1)) {
					$(`#envFreqLoop\\[${i + 1}\\] option[value='S']`).prop('selected', true);
					pointStyleMetadata[freqLowIdx].sustainPt = i;
					pointStyleMetadata[freqUpIdx].sustainPt = i;
				}
				if (envelopes.FreqEnvelope.LOOPPT == (i + 1)) {
					var v = envelopes.FreqEnvelope.ENVTYPE == 3 ? 'L' : 'R'
					$(`#envFreqLoop\\[${i + 1}\\] option[value='${v}']`).prop('selected', true);
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
		var maxTotalTime = Math.max(totalTimeLow, totalTimeUp);

		totalTimeLow = 0;
		totalTimeUp = 0;


		for (i = envelopes.FreqEnvelope.NPOINTS; i < 16; i++) {
			// hide unused rows

			$(`#envAmpLoop\\[${i + 1}\\]`).hide();
			$(`#envAmpLowVal\\[${i + 1}\\]`).hide();
			$(`#envAmpUpVal\\[${i + 1}\\]`).hide();
			$(`#envAmpLowTime\\[${i + 1}\\]`).hide();
			$(`#envAmpUpTime\\[${i + 1}\\]`).hide();

			if (animate) {
				viewVCE_voice.sendToCSurface(null, `envAmpLowVal[${i + 1}]`, 0);
				viewVCE_voice.sendToCSurface(null, `envAmpUpVal[${i + 1}]`, 0);
				viewVCE_voice.sendToCSurface(null, `envAmpLowTime[${i + 1}]`, 0);
				viewVCE_voice.sendToCSurface(null, `envAmpUpTime[${i + 1}]`, 0);
			}
		}

		// Amp envelopes have an implicit start point at time zero, value zero
		datasets[ampLowIdx].data.push({ x: 0, y: 0 });
		datasets[ampUpIdx].data.push({ x: 0, y: 0 });

		for (i = 0; i < envelopes.AmpEnvelope.NPOINTS; i++) {
			var tr = $('#envTable tbody tr:eq(' + i + ')');

			$(`#envAmpLoop\\[${i + 1}\\]`).show();
			$(`#envAmpLowVal\\[${i + 1}\\]`).show();
			$(`#envAmpUpVal\\[${i + 1}\\]`).show();
			$(`#envAmpLowTime\\[${i + 1}\\]`).show();
			$(`#envAmpUpTime\\[${i + 1}\\]`).show();

			// table is logically in groups of 4.
			// "j" accounts for the difference in column index due to the
			// row-spanning separators (only in i==0):
			j = (i == 0) ? 11 : 9;

			//	    console.dir(tr);
			//	    console.dir(tr.find('td:eq(' +(j+0)+ ')'));
			var isLast = (i + 1) >= envelopes.AmpEnvelope.NPOINTS;
			var ampLow = viewVCE_envs.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i * 4 + 0]);
			var ampUp = viewVCE_envs.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i * 4 + 1]);
			var timeLow = viewVCE_envs.scaleAmpTimeValue(envelopes.AmpEnvelope.Table[i * 4 + 2]);
			var timeUp = viewVCE_envs.scaleAmpTimeValue(envelopes.AmpEnvelope.Table[i * 4 + 3]);

			if (animate) {
				viewVCE_voice.sendToCSurface(null, `envAmpLowVal[${i + 1}]`, envelopes.AmpEnvelope.Table[i * 4 + 0]);
				viewVCE_voice.sendToCSurface(null, `envAmpUpVal[${i + 1}]`, envelopes.AmpEnvelope.Table[i * 4 + 1]);
				viewVCE_voice.sendToCSurface(null, `envAmpLowTime[${i + 1}]`, envelopes.AmpEnvelope.Table[i * 4 + 2]);
				viewVCE_voice.sendToCSurface(null, `envAmpUpTime[${i + 1}]`, envelopes.AmpEnvelope.Table[i * 4 + 3]);
			}

			lastAmpLow = ampLow;
			lastAmpUp = ampUp;
			totalTimeLow += timeLow;
			totalTimeUp += timeUp;

			datasets[ampLowIdx].data.push({ x: totalTimeLow, y: ampLow });
			datasets[ampUpIdx].data.push({ x: totalTimeUp, y: ampUp });

			//	    console.dir(datasets[ampLowIdx]);

			document.getElementById(`envAmpLowVal[${i + 1}]`).value = ampLow;
			document.getElementById(`envAmpUpVal[${i + 1}]`).value = ampUp;
			document.getElementById(`envAmpLowTime[${i + 1}]`).value = timeLow;
			document.getElementById(`envAmpUpTime[${i + 1}]`).value = timeUp;
			document.getElementById(`envAmpTotLowTime[${i + 1}]`).innerHTML = totalTimeLow;
			document.getElementById(`envAmpTotUpTime[${i + 1}]`).innerHTML = totalTimeUp;

			if (isLast) {
				document.getElementById(`envAmpLowVal[${i + 1}]`).disabled = true;
				document.getElementById(`envAmpUpVal[${i + 1}]`).disabled = true;
			} else if (viewVCE_voice.voicingMode) {
				document.getElementById(`envAmpLowVal[${i + 1}]`).disabled = false;
				document.getElementById(`envAmpUpVal[${i + 1}]`).disabled = false;
			}

			if (envelopes.AmpEnvelope.ENVTYPE != 1) {
				if (envelopes.AmpEnvelope.SUSTAINPT == (i + 1)) {
					$(`#envAmpLoop\\[${i + 1}\\] option[value='S']`).prop('selected', true);
					// we draw an extra point for amp curve - so the index of the loop point is i+1
					pointStyleMetadata[ampLowIdx].sustainPt = i+1;
					pointStyleMetadata[ampUpIdx].sustainPt = i+1;
				}
				if (envelopes.AmpEnvelope.LOOPPT == (i + 1)) {
					var v = envelopes.AmpEnvelope.ENVTYPE == 3 ? 'L' : 'R'
					$(`#envAmpLoop\\[${i + 1}\\] option[value='${v}']`).prop('selected', true);
					if (v === 'L') {
						// we draw an extra point for amp curve - so the index of the loop point is i+1
						pointStyleMetadata[ampLowIdx].loopPt = i+1;
						pointStyleMetadata[ampUpIdx].loopPt = i+1;
					} else {
						pointStyleMetadata[ampLowIdx].repeatPt = i+1;
						pointStyleMetadata[ampUpIdx].repeatPt = i+1;
					}
				}
			}
		}

		maxTotalTime = Math.max(maxTotalTime, totalTimeLow);
		maxTotalTime = Math.max(maxTotalTime, totalTimeUp);

		/*
	need to confirm the actual behavior of the envelopes to determine if this
	visualization makes sense
		datasets[freqLowIdx].data.push({x: maxTotalTime, y: lastFreqLow});
		datasets[freqUpIdx].data.push( {x: maxTotalTime,  y: lastFreqUp});
		datasets[ampLowIdx].data.push({x: maxTotalTime, y: lastAmpLow});
		datasets[ampUpIdx].data.push( {x: maxTotalTime,  y: lastAmpUp});
		*/

		//	console.dir(datasets);

		var animation_duration = animate ? 1000 : 0;

		var filteredPointStyleMetadata = [];
		var filteredDatasets = [];
		if (envNum < 0) {
			// all of them:
			filteredDatasets = datasets;
			filteredPointStyleMetadata = pointStyleMetadata;
		} else {
			filteredDatasets.push(datasets[envNum])
			filteredPointStyleMetadata.push(pointStyleMetadata[envNum])
		}
		var ctx = document.getElementById('envChart').getContext('2d');
		if (viewVCE_envs.chart != null) {
			viewVCE_envs.chart.destroy();
		}
		var timeAxisType = document.getElementById('timeScale').value;
		var freqAxisType = document.getElementById('freqScale').value;

		viewVCE_envs.chart = new Chart(ctx, {

			type: 'scatter',
			data: {
				//		labels: ['','','','','','','','','','', '','','','','','','','','','', '','','','','','','','','','','',''],
				datasets: filteredDatasets
			},

			// Configuration options go here
			options: {
				animation: {
					duration: animation_duration
				},
				tooltips: {
					mode: 'index',
				},
				hover: {
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
						scaleLabel: {
							display: true,
							labelString: "Time (ms)"
						},

						// after a pan, the graph sometimes shows poorly formatted values for min and max (apparently bypasses the tick callback?)
						//Workaround by just not displaying them
						afterTickToLabelConversion: function(scaleInstance) {
							// set the first and last tick to null so it does not display
							// note, ticks[0] is the last tick and ticks[length - 1] is the first
							scaleInstance.ticks[0] = null;
							scaleInstance.ticks[scaleInstance.ticks.length - 1] = null;

							// need to do the same thing for this similiar array which is used internally
							//scaleInstance.ticksAsNumbers[0] = null;
							//scaleInstance.ticksAsNumbers[scaleInstance.ticksAsNumbers.length - 1] = null;
						},

						ticks: {
							precision: 2,
							callback: function (value, index, values) {
								if (value >= 1.0) {
									return value.toFixed(0);
								} else {
									v = value.toFixed(2);
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
					yAxes: [{
						position: 'left',
						id: 'freq-axis',
						type: freqAxisType,
						gridLines: {
							color: '#666',
							display: true,
							drawBorder: true,
							drawOnChartArea: false
						},
						scaleLabel: {
							display: true,
							labelString: "Frequency (Hz)"
						},

						// after a pan, the graph sometimes shows poorly formatted values for min and max (apparently bypasses the tick callback?)
						//Workaround by just not displaying them
						afterTickToLabelConversion: function(scaleInstance) {
							// set the first and last tick to null so it does not display
							// note, ticks[0] is the last tick and ticks[length - 1] is the first
							scaleInstance.ticks[0] = null;
							scaleInstance.ticks[scaleInstance.ticks.length - 1] = null;

							// need to do the same thing for this similiar array which is used internally
							//scaleInstance.ticksAsNumbers[0] = null;
							//scaleInstance.ticksAsNumbers[scaleInstance.ticksAsNumbers.length - 1] = null;
						},

						ticks: {
							precision: 2,
							callback: function (value, index, values) {
								// don't use scientific notation
								if (value >= 1.0) {
									return value.toFixed(0);
								} else {
									v = value.toFixed(2);
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
						scaleLabel: {
							display: true,
							labelString: "Amplitude dB"
						},

						// after a pan, the graph sometimes shows poorly formatted values for min and max (apparently bypasses the tick callback?)
						//Workaround by just not displaying them
						afterTickToLabelConversion: function(scaleInstance) {
							// set the first and last tick to null so it does not display
							// note, ticks[0] is the last tick and ticks[length - 1] is the first
							scaleInstance.ticks[0] = null;
							scaleInstance.ticks[scaleInstance.ticks.length - 1] = null;

							// need to do the same thing for this similiar array which is used internally
							//scaleInstance.ticksAsNumbers[0] = null;
							//scaleInstance.ticksAsNumbers[scaleInstance.ticksAsNumbers.length - 1] = null;
						},

						ticks: {
							color: '#eee',
							display: true
						}
					}],
				},
				responsive: true,
				maintainAspectRatio: false,

				dragData: viewVCE_voice.voicingMode,
				dragDataRound: 0,
				dragX: true,

				dragOptions: {
					showTooltip: true
				},

				onDragStart: function(e, element) {
					viewVCE_envs.dragging = true;
					//console.log('onDragStart: ', envNum, e, element)
					// constrain amp curve:  first point fixed at 0,0
					//
					// HACK: also if only a single point in the dataset (common for freq envs), don't allow drag.
					//This works around an as yet undiagnosed bug that cases vlaues to go to floating point numbers smaller than zero and cofuse the auto-scaling.
					if (viewVCE_envs.chart.data.datasets[element._datasetIndex].data.length==1 ||
						((envNum < 0 && element._datasetIndex >= 2)||(envNum>=2)) && element._index === 0) {
						// can't move the first amp point
						dragOldValue.x = viewVCE_envs.chart.data.datasets[element._datasetIndex].data[element._index].x;
						dragOldValue.y = viewVCE_envs.chart.data.datasets[element._datasetIndex].data[element._index].y;
						//console.log("ondragStart: freeze 0th amp",element._datasetIndex,element._index,dragOldValue)
						viewVCE_envs.chart.update(0);
					}
				},
				onDrag: function(e, datasetIndex, index, value) {
					if (viewVCE_envs.chart.data.datasets[datasetIndex].data.length==1 ||
						((envNum < 0 && datasetIndex >= 2)||(envNum>=2)) && index === 0) {
						// can't move the first amp point
						viewVCE_envs.chart.data.datasets[datasetIndex].data[index].x = dragOldValue.x;
						viewVCE_envs.chart.data.datasets[datasetIndex].data[index].y = dragOldValue.y;
						//console.log("ondrag: freeze 0th amp",datasetIndex,index,dragOldValue)
						viewVCE_envs.chart.update(0);
						return
					}
					e.target.style.cursor = 'grabbing'
					// time must stay between neighboring points:
					var min = index > 0 ? viewVCE_envs.chart.data.datasets[datasetIndex].data[index-1].x : 0;
					// if the last point, use the scale max
					var max = (index === (viewVCE_envs.chart.data.datasets[datasetIndex].data.length-1))
					    ? (viewVCE_envs.chart.scales['time-axis'].max + 1)
						: viewVCE_envs.chart.data.datasets[datasetIndex].data[index+1].x;
					// if this is a freq env, then the 0th point's x value is fixed at 0
					if (((envNum < 0 && datasetIndex < 2)||(envNum<2)) && index === 0) {
						min = -1;
						max = 2; // the clamping expression below subtracts or adds 1
					}
					//console.log('ondrag: ', envNum, datasetIndex, index, value, min, max)
					if (value.x >= max) {
						value.x = max-1;
						//console.log('onDrag: CLAMP ', datasetIndex, index, value)
						viewVCE_envs.chart.update();
					} else if (value.x <= min) {
						value.x = min+1;
						//console.log('onDrag: CLAMP ', datasetIndex, index, value)
						viewVCE_envs.chart.update();
					}
					viewVCE_envs.updateEnvFromGraphChange(datasetIndex, index, value, false)
				},
				onDragEnd: function(e, datasetIndex, index, value) {
					viewVCE_envs.dragging = false;
					e.target.style.cursor = 'default'
					//console.log('onDragEnd: ', datasetIndex, index, value)

					viewVCE_envs.updateEnvFromGraphChange(datasetIndex, index, value, true)
				},

				tooltips: {
					mode: 'index',
				},
				hover: {
					mode: 'index',
					intersect: true,
					onHover: function (e) {
						const point = this.getElementAtEvent(e)
						if (viewVCE_voice.voicingMode && point.length
							&& !(((envNum < 0 && point[0]._datasetIndex >= 2) || (envNum >= 2)) && point[0]._index === 0)) {
							e.target.style.cursor = 'grab';
						} else {
							e.target.style.cursor = 'default';
						}
					}
				},
				plugins: {
					// zoom plugin is only used by the env graphs
					zoom: {
						zoom: {
							enabled: false
						},
						pan: {
							enabled: true,
							mode: function ({chart}) {
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
	},

	updateEnvFromGraphChange: function (datasetIndex, index, value, fireOnChange) {
		// now reverse engineer the changed env values and update the corresponding point value or time
		// if x has changed, then the TIME value for both the point and the preceding point need to change
		// (since the env values are the delta-t from the previous point, not the absolute t of the point)
		// if y has changed, only its value needs to be updated.
		var newV = value.y
		var newT
		var nextNewT = undefined
		var fieldIndex = index + 1; // fields are 1-based

		if (datasetIndex >= 2) {
			// amp.  the first point in the env corresponds to the second point on the graph
			fieldIndex = fieldIndex - 1;
			if (index === 1) {
				newT = value.x
			} else {
				newT = (value.x
					- viewVCE_envs.chart.data.datasets[datasetIndex].data[index - 1].x);
			}
		} else {
			// freq.  the first point in the env corresponds to the first point on the graph
			if (index === 0) {
				newT = undefined;
			} else {
				newT = (value.x
					- viewVCE_envs.chart.data.datasets[datasetIndex].data[index - 1].x);
			}
		}
// last point case is common to both types of env
// if last point, there's no nextT
		if (index != viewVCE_envs.chart.data.datasets[datasetIndex].data.length - 1) {
			nextNewT = (viewVCE_envs.chart.data.datasets[datasetIndex].data[index + 1].x
				- value.x);
		}
		console.log("UPDATE VALUES", fieldIndex, newV, newT, nextNewT)

		function setValueAndFireOnchange(id, val) {
			ele = document.getElementById(id);
			ele.value = val;
			// don't run onchange during the drag - since we redraw the graph after sending data to the Synergy
			//(and that aborts the drag)
			if (fireOnChange) {
				// just call the function directly; faking the event in the browser is error prone
				viewVCE_envs.onchange(ele);
			}
		}

		setValueAndFireOnchange(`${viewVCE_envs.valFieldNameByDatasetIdx[datasetIndex]}[${fieldIndex}]`, newV);
		if (newT != undefined) {
			setValueAndFireOnchange(`${viewVCE_envs.timeFieldNameByDatasetIdx[datasetIndex]}[${fieldIndex}]`, newT);
		}
		if (nextNewT != undefined) {
			setValueAndFireOnchange(`${viewVCE_envs.timeFieldNameByDatasetIdx[datasetIndex]}[${fieldIndex + 1}]`, nextNewT);
		}
	}

};
