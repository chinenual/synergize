let viewVCE_voice = {
	voicingMode: false,

	timbreProportionCurve: function (center, sensitivity) {
		var result = [];
		if (sensitivity == 0) {
			for (v = 0; v < 32; v++) {
				result[v] = center;
			}
			return result;
		}
		// center = 0..32
		// sensitivity = 1..31
		for (v = 0; v < 32; v++) {
			// this appears to be what the z80 code is doing for timbre PROPC:
			var p = (center * 2) - 15 + ((v / 10) * sensitivity) - (2 * sensitivity);
			if (p > 31) p = 31;
			if (p < 0) p = 0;
			result[v] = p;
		}
		return result;
	},

	ampProportionCurve: function (center, sensitivity) {
		var result = [];
		if (sensitivity == 0) {
			for (v = 0; v < 32; v++) {
				result[v] = center;
			}
			return result;
		}
		// center = 0..32
		// sensitivity = 1..31
		for (v = 0; v < 32; v++) {
			// this appears to be what the z80 code is doing for timbre PROPC:
			var p = ((((v / 10) * sensitivity) / 2.0 - sensitivity) / 2.0) + (center - 24)
			if (p > 6) p = 6;
			if (p < -24) p = -24;
			result[v] = p + 25;
		}
		return result;
	},

	SOLO: [],
	MUTE: [],

	toggleOsc: function (ele) {
		console.log("toggle " + ele.id);
		var oscPattern = /([A-Z]+)\[(\d+)\]/;
		if (ret = ele.id.match(oscPattern)) {
			param = ret[1];
			osc = parseInt(ret[2], 10); /* 1-based */

			state = ele.classList.contains("on");
			viewVCE_voice[param][osc - 1] = !state;
			ele.classList.toggle('on');

			let message = {
				"name": "setOscSolo",
				"payload": {
					"Mute": viewVCE_voice.MUTE,
					"Solo": viewVCE_voice.SOLO,
				}
			};
			astilectron.sendMessage(message, function (message) {
				console.log("setOscSolo returned: " + JSON.stringify(message));
				// Check error
				if (message.name === "error") {
					// failed - dont change the boolean
					index.errorNotification(message.payload);
					return false;
				}
			});
		}
	},

	filterChanged: function (ele) {
		var id = ele.id;
		console.log("filterChanged: " + id + " val: " + ele.value);

		var filterPattern = /FILTER\[(\d+)\]/;
		if (ret = id.match(filterPattern)) {
			osc = parseInt(ret[1])
		} else {
			console.log("ERROR: filterCHanged called with bad ele " + ele);
		}
		var filterValue;
		if (ele.value == "") {
			filterValue = 0;
		} else {
			filterValue = parseInt(ele.value, 10);
		}
		vce.Head.FILTER[osc - 1] = filterValue;
		viewVCE_filters.init(true);
	},

	OHARMToText: function (str) {
		var newStr;
		var val = parseInt(str, 10);
		if (val < 0) {
			newStr = "s" + (-val);
		} else if (val == 31) {
			newStr = "dc";
		} else {
			newStr = str;
		}
		//console.log("   OHARMToText(" + str + ") returns " + newStr);
		return newStr;
	},

	TextToOHARM: function (str) {
		var newStr;
		if (str === "dc") {
			newStr = '31';
		} else if (ret = str.match(/s(\d+)/)) {
			val = parseInt(ret[1], 10);
			newStr = '' + (-val);
		} else if (ret = str.match(/\d+/)) {
			newStr = str;
		} else {
			/// error! 
			console.log("ERROR: TextToOHARM cant decode " + str);
			newStr = str;
		}
		//console.log("   TextToOHARM(" + str + ") returns " + newStr);
		return newStr;
	},

	FDETUNToText: function (str) {
		/* THIS IS UGLY.  The Synergy uses a non-linear mapping of this byte to a 1/30Hz increment, 
		* with 5 positive values reserved for "random" settings. Ick ick ick.  
		*
		* See D.DTN routine in VOIDSP.Z80 in the SYNHCS sourcecode.
		*/
		var newStr;
		var val = parseInt(str, 10);
		if (val > 58) {
			// CASE A
			newStr = "ran" + (val - 58);
		} else if (val >= -32 && val <= 32) {
			// CASE B
			val = val * 3;
			newStr = '' + val;
		} else if (val > 0) {
			// CASE C
			val = ((val * 2) - 32) * 3;
			newStr = '' + val;
		} else { // negative
			// CASE D
			val = ((val * 2) + 32) * 3;
			newStr = '' + val;
		}
		//console.log("   FDETUNToText(" + str + ") returns " + newStr);
		return newStr;
	},

	TextToFDETUN: function (str) {
		// See FDETUNToText.  This "reverses" that attrocity
		var newStr;
		if (ret = str.match(/ran(\d+)/)) {
			// CASE A
			var val = parseInt(ret[1], 10);
			val += 58;
			newStr = '' + val;
		} else {
			var val = parseInt(str, 10);
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
		//console.log("   TextToFDETUN(" + str + ") returns " + newStr);
		return newStr;
	},

	NullablePatchRegisterToText: function (str) {
		if (str.trim() === '') {
			return '';
		}
		var val = parseInt(str, 10);
		if (val === 0) {
			return ''
		}
		return '' + val;
	},

	TextToNullablePatchRegister: function (str) {
		if (str.trim() === '') {
			return '0';
		}
		var val = parseInt(str, 10);
		if (val === 0) {
			return '0'
		}
		return '' + val;
	},

	testConversionFunctions: function () {
		var ok = true;
		for (var i = 0; i < 5; i++) {
			var str = viewVCE_voice.NullablePatchRegisterToText('' + i);
			var reverseStr = viewVCE_voice.TextToNullablePatchRegister(str);
			if (('' + i) != reverseStr) {
				ok = false;
				console.log("ERROR: PatchReg " + i + " totext: " + str + " reversed to " + reverseStr)
			}
		}
		for (var i = -11; i <= 31; i++) {
			var str = viewVCE_voice.OHARMToText('' + i);
			var reverseStr = viewVCE_voice.TextToOHARM(str);
			if (('' + i) != reverseStr) {
				ok = false;
				console.log("ERROR: OHARM " + i + " totext: " + str + " reversed to " + reverseStr)
			}
		}
		for (var i = -63; i <= 63; i++) {
			var str = viewVCE_voice.FDETUNToText('' + i);
			var reverseStr = viewVCE_voice.TextToFDETUN(str);
			if (('' + i) != reverseStr) {
				ok = false;
				console.log("ERROR: FDETUN " + i + " totext: " + str + " reversed to " + reverseStr)
			}
		}
		//if (val >= (-32 * 3) && val <= (32 * 3)) {
		//	// CASE B
		//} else if (val > 0) {
		//	// CASE C
		//} else {
		//	// CASE D
		//}
		// test that the FDETUN rounding does the right thing: 
		var v = viewVCE_voice.FDETUNToText(viewVCE_voice.TextToFDETUN('20'));
		if ('21' != v) {
			ok = false;
			console.log("ERROR: FDETUN CASE B " + 20 + " rounded to " + v + " - expected " + 21)
		}
		v = viewVCE_voice.FDETUNToText(viewVCE_voice.TextToFDETUN('247'));
		if ('246' != v) {
			ok = false;
			console.log("ERROR: FDETUN CASE C " + 247 + " rounded to " + v + " - expected " + 246)
		}
		v = viewVCE_voice.FDETUNToText(viewVCE_voice.TextToFDETUN('-205'));
		if ('-204' != v) {
			ok = false;
			console.log("ERROR: FDETUN CASE D " + -205 + " rounded to " + v + " - expected " + -204)
		}

		console.log("viewVCE_voice.testConversionFunctions: " + (ok ? "PASS" : "FAIL"));
		return ok;
	},

	onchangeDSR: function (param, osc /*1-based*/, value) {
		osc = parseInt(osc, 10);

		console.log("onchangeDSR: " + param + "[" + osc + "] == " + value);

		// displayed values are 1-based, bit values in the patch byte are 0-based
		var patchFOInputDSR = document.getElementById(`patchFOInputDSR[${osc}]`).value;
		var patchAdderInDSR = document.getElementById(`patchAdderInDSR[${osc}]`).value;
		var patchOutputDSR = document.getElementById(`patchOutputDSR[${osc}]`).value;

		// XREF: patch byte encode/decode
		var patchInhibitAddr = patchAdderInDSR == '' ? true : false;
		var patchInhibitF0 = patchFOInputDSR == '' ? true : false;
		var patchByte = 0;
		patchByte |= ((parseInt(patchFOInputDSR, 10) - 1) & 0x03);
		patchByte |= (((parseInt(patchAdderInDSR, 10) - 1) << 3) & 0x18);
		patchByte |= (((parseInt(patchOutputDSR, 10) - 1) << 6) & 0xc0);
		if (patchInhibitAddr) {
			patchByte |= 0x20;
		}
		if (patchInhibitF0) {
			patchByte |= 0x04;
		}

		console.log(osc + " old patch byte: " + vce.Envelopes[osc - 1].FreqEnvelope.OPTCH + "\n" +
			" new patch byte: " + patchByte + "\n" +
			" patchInhibitAddr : " + patchInhibitAddr + "\n" +
			" patchInhibitF0   : " + patchInhibitF0 + "\n" +
			" patchOutputDSR   : " + patchOutputDSR + "\n" +
			" patchAdderInDSR  : " + patchAdderInDSR + "\n" +
			" patchFOInputDSR  : " + patchFOInputDSR + "\n");

		let message = {
			"name": "setPatchByte",
			"payload": {
				"Osc": osc,
				"Value": patchByte
			}
		};
		astilectron.sendMessage(message, function (message) {
			console.log("setPatchByte returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the boolean
				index.errorNotification(message.payload);
				return false;
			}
			vce.Envelopes[osc - 1].FreqEnvelope.OPTCH = patchByte;
			viewVCE_voice.patchTable(); // in case the patch diagram changes due to the edit
			viewVCE_voice.voicingModeVisuals();
		});


	},

	onchange: function (ele, updater, valueConverter) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		viewVCE_voice.deb_onchange(ele, updater, valueConverter);
	},

	deb_onchange: null,

	raw_onchange: function (ele, updater, valueConverter) {
		if (viewVCE.supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/ return; }
		var id = ele.id;
		console.log("changed: " + id + ", new value: " + ele.value);

		if (valueConverter == undefined) {
			// identity function if none passed
			valueConverter = function (v) { return v };
		}

		var value
		var param
		var args
		var filterPattern = /FILTER\[(\d+)\]/;
		var dsrPattern = /([0-9A-Za-z]+DSR)\[(\d+)\]/;
		var waveKeyPattern = /wk([A-Z]+)\[(\d+)\]/;
		var oscPattern = /([A-Z]+)\[(\d+)\]/;
		var headPattern = /([A-Z]+)/;
		if (id === "VNAME") {
			param = "VNAME"
			funcname = "setVNAME";
			args = ele.value;
		} else if (ret = id.match(dsrPattern)) {
			param = ret[1]
			osc = ret[2]
			value = parseInt(valueConverter(ele.value), 10);
			return viewVCE_voice.onchangeDSR(param, osc, value);
		} else if (ret = id.match(filterPattern)) {
			param = "FILTER"
			funcname = "setOscFILTER";
			osc = parseInt(ret[1])
			value = parseInt(valueConverter(ele.value), 10);
			args = [osc, value];
		} else if (ret = id.match(waveKeyPattern)) {
			param = ret[1]
			osc = parseInt(ret[2], 10)
			if (param == "WAVE") {
				funcname = "setOscWAVE"
				value = ele.value == "Sin" ? 0 : 1;
				args = [osc, value];
			} else {
				funcname = "setOscKEYPROP"
				value = ele.checked ? 1 : 0;
				args = [osc, value];
			}
		} else if (ret = id.match(oscPattern)) {
			param = ret[1];
			osc = parseInt(ret[2], 10)
			value = parseInt(valueConverter(ele.value), 10)
			args = [osc, value]

			console.log("changed: " + id + " param: " + param + " osc: " + osc);
			vce.Envelopes[osc - 1][param] = valueConverter(ele.value);
			funcname = "setVoiceByte"

		} else if (ret = id.match(headPattern)) {
			param = id;
			value = parseInt(valueConverter(ele.value), 10)
			args = [value]
			vce.Head[param] = valueConverter(ele.value);
			funcname = "setVoiceByte"
		}
		//console.dir(vce);
		if (param != null) {
			let message = {
				"name": funcname,
				"payload": {
					"Param": param,
					"Args": args
				}
			};
			astilectron.sendMessage(message, function (message) {
				console.log(funcname + " returned: " + JSON.stringify(message));
				// Check error
				if (message.name === "error") {
					// failed - dont change the boolean
					index.errorNotification(message.payload);
					return false;
				} else {
					viewVCE_voice.sendToCSurface(ele, ele.id, value);
					if (updater != undefined) {
						console.log("updater: " + updater);
						updater(ele);
					}
				}
			});

		}
		return true;
	},

	patchTable: function () {
		document.getElementById("patchType").value = vce.Extra.PatchType;

		var tbody = document.getElementById("patchTbody");
		// remove old rows:
		while (tbody.firstChild) {
			tbody.removeChild(tbody.firstChild);
		}

		var outRegisters = [[], [], [], []];
		var freqDAG = "";

		viewVCE_voice.sendToCSurface(null, `num-osc`, vce.Head.VOITAB + 1);
		/*
		for (osc = vce.Head.VOITAB + 1; osc < 16; osc++) {
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

		// populate new ones:
		for (osc = 0; osc <= vce.Head.VOITAB; osc++) {
			viewVCE_voice.sendToCSurface(null, `osc-enabled[${osc + 1}]`, 1);

			var tr = document.createElement("tr");
			var td = document.createElement("td");

			//--- OSC
			td.innerHTML = osc + 1; // Osc
			// Mute 
			span = document.createElement("span");
			span.innerHTML = `&nbsp;&nbsp;<span onclick="viewVCE_voice.toggleOsc(this)" class="vceEditToggleText" id="MUTE[${osc + 1}]">M</span>`;
			td.append(span);
			// Solo
			span = document.createElement("span");
			span.innerHTML = `&nbsp;<span onclick="viewVCE_voice.toggleOsc(this)" class="vceEditToggleText" id="SOLO[${osc + 1}]">S</span>`;
			td.append(span);

			tr.appendChild(td);

			// XREF: patch byte encode/decode
			// FIXME: assumes envelopes are sorted in oscillator order
			var patchByte = vce.Envelopes[osc].FreqEnvelope.OPTCH;
			var patchInhibitAddr = (patchByte & 0x20) != 0;
			var patchInhibitF0 = (patchByte & 0x04) != 0;
			var patchOutputDSR = ((patchByte & 0xc0) >> 6);
			var patchAdderInDSR = ((patchByte & 0x18) >> 3);
			var patchFOInputDSR = (patchByte & 0x03);

			//			console.log(osc + " patch byte: " + patchByte + "\n" +
			//				" patchInhibitAddr : " + patchInhibitAddr + "\n" +
			//				" patchInhibitF0   : " + patchInhibitF0 + "\n" +
			//				" patchOutputDSR   : " + patchOutputDSR + "\n" +
			//				" patchAdderInDSR  : " + patchAdderInDSR + "\n" +
			//				" patchFOInputDSR  : " + patchFOInputDSR + "\n");

			// compute the DAG based on current register usage:
			if (!patchInhibitF0) {
				var modulatingOscs = outRegisters[patchFOInputDSR];
				for (var i = 0; i < modulatingOscs.length; i++) {
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
			td = document.createElement("td");
			var reg = 0;
			if (!patchInhibitF0) {
				reg = patchFOInputDSR + 1;
			} else {
				reg = 0;
			}
			if (osc == 0) {
				// the first osc's Freq DSR can't be altered - render as a disabled input control so we can get its value in the onchange
				// function without any special casing
				td.innerHTML = `<div class="spinwrapper"><input type="text" class="vceNum vceEditDisabled" id="patchFOInputDSR[${osc + 1}]" 
				value="${viewVCE_voice.NullablePatchRegisterToText('' + reg)}" 
				disabled/></div>`;
			} else {
				td.innerHTML = `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinNullablePatchReg" id="patchFOInputDSR[${osc + 1}]" 
				onchange="viewVCE_voice.onchange(this,undefined,viewVCE_voice.TextToNullablePatchRegister)" value="${viewVCE_voice.NullablePatchRegisterToText('' + reg)}" 
				min="0" max="4"
				disabled/></div>`;
			}
			tr.appendChild(td);

			//--- Patch A
			td = document.createElement("td");
			if (!patchInhibitAddr) {
				reg = patchAdderInDSR + 1;
			} else {
				reg = 0;
			}
			td.innerHTML = `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinNullablePatchReg" id="patchAdderInDSR[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this,undefined,viewVCE_voice.TextToNullablePatchRegister)" value="${viewVCE_voice.NullablePatchRegisterToText('' + reg)}" 
			min="0" max="4"
			disabled/></div>`;
			tr.appendChild(td);

			//--- Patch O
			td = document.createElement("td");
			td.innerHTML = `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinPlain" id="patchOutputDSR[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this,undefined,undefined)" value="${patchOutputDSR + 1}" 
			min="1" max="4"
			disabled/></div>`;
			tr.appendChild(td);

			//--- Hrm
			// HACK: we wrap these elements in <div> of a fixed size to keep the up/down buttons positioned properly.
			// Someone more skilled in the ways of CSS would surely have a cleaner solution.
			td = document.createElement("td");
			td.innerHTML = `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinOHARM" id="OHARM[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this,undefined,viewVCE_voice.TextToOHARM)" value="${viewVCE_voice.OHARMToText(vce.Envelopes[osc].FreqEnvelope.OHARM)}" 
			min="-11" max="31"
			disabled/></div>`;
			tr.appendChild(td);
			viewVCE_voice.sendToCSurface(null, `OHARM[${osc + 1}]`, vce.Envelopes[osc].FreqEnvelope.OHARM + 11)

			//--- Detn
			td = document.createElement("td");
			td.innerHTML = `<div class="spinwrapper"><input type="text" class="vceEdit vceNum spinFDETUN" id="FDETUN[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this,undefined,viewVCE_voice.TextToFDETUN)" value="${viewVCE_voice.FDETUNToText(vce.Envelopes[osc].FreqEnvelope.FDETUN)}" 
			min="-63" max="63"
			disabled/></div>`;
			tr.appendChild(td);
			viewVCE_voice.sendToCSurface(null, `FDETUN[${osc + 1}]`, vce.Envelopes[osc].FreqEnvelope.FDETUN + 63)

			var waveByte = vce.Envelopes[osc].FreqEnvelope.Table[3];
			var wave = ((waveByte & 0x1) == 0) ? 'Sin' : 'Tri';
			var keyprop = ((waveByte & 0x10) == 0) ? false : true;

			//--- Wave
			td = document.createElement("td");
			td.innerHTML = wave;
			td.innerHTML = `<select class="vceEdit" id="wkWAVE[${osc + 1}]" value="${wave}" 
			onchange="viewVCE_voice.onchange(this)" disabled/>
			<option ${wave == 'Sin' ? "selected" : ""} value="Sin">Sin</option>
			<option ${wave == 'Tri' ? "selected" : ""} value="Tri">Tri</option>
			</select>
			`;
			tr.appendChild(td);
			viewVCE_voice.sendToCSurface(null, `wkWAVE[${osc + 1}]`, wave == 'Sin' ? 0 : 1);

			//--- Key
			td = document.createElement("td");
			// can't use disabled attr - bootstrap styling hides it - use javascript hack to make it readonnly
			td.innerHTML = `<input type="checkbox" id="wkKEYPROP[${osc + 1}]" value="true" 
			${keyprop ? " checked " : ""} 
			onchange="viewVCE_voice.voicingMode ? viewVCE_voice.onchange(this) : (this.checked=!this.checked)"/>`;
			tr.appendChild(td);
			viewVCE_voice.sendToCSurface(null, `wkKEYPROP[${osc + 1}]`, keyprop ? 1 : 0);

			//--- Flt
			td = document.createElement("td");
			var filter = vce.Head.FILTER[osc];
			td.innerHTML =
				(filter == 0) ? ''
					: (filter > 0) ? ('Bf ' + filter) : ('Af ' + -filter);
			td = document.createElement("td");
			td.innerHTML = wave;
			td.innerHTML = `<select class="vceEdit" id="FILTER[${osc + 1}]" value="${filter}" 
					onchange="viewVCE_voice.onchange(this,viewVCE_voice.filterChanged)" disabled/>
					<option ${filter == 0 ? "selected" : ""} value=""></option>
					<option ${filter < 0 ? "selected" : ""} value="-1">Af</option>
					<option ${filter > 0 ? "selected" : ""} value="${osc + 1}">Bf</option>
					</select>
					`;
			tr.appendChild(td);
			viewVCE_voice.sendToCSurface(null, `FILTER[${osc + 1}]`, filter == 0 ? 0 : filter < 0 ? 1 : 2);

			tbody.appendChild(tr);
		}
		// final row is the plus/minus buttons
		{
			let temp = document.createElement('template');
			temp.innerHTML = `<tr class="listplusminus" id="oscPlusMinus" style="display:none;">
							    <td colspan="9">
			    					<div style="margin-top: 5px; float: left;">
				    					<input id="delOsc" type='button' value='-'
										    onclick='viewVCE_voice.setNumOscillators(vce.Head.VOITAB)' />
									    <input id="addOsc" type='button' value='+'
									   	    onclick='viewVCE_voice.setNumOscillators(vce.Head.VOITAB+2)' />
								    </div>
								</td>
							</tr>`;
			tbody.appendChild(temp.content.firstChild);
		}

		console.log("freqDAG: " + freqDAG);
		// Generate the patch diagram:
		var patchDiagramCanvas = document.getElementById('patchDiagram');
		// nomnoml is confused by leading spaces on directives lines, so...:
		var patchDiagramSource =
			`
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
		console.log("nomnoml src: " + patchDiagramSource);
		nomnoml.draw(patchDiagramCanvas, patchDiagramSource);


	},

	changePatchType: function (newIndex) {
		let message = {
			"name": "setPatchType",
			"payload": parseInt(newIndex, 10)
		};
		index.spinnerOn();
		astilectron.sendMessage(message, function (message) {
			index.spinnerOff();
			console.log("setPatchType returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the boolean
				index.errorNotification(message.payload);
			} else {
				for (i = 0; i < vce.Envelopes.length; i++) {
					vce.Envelopes[i].FreqEnvelope.OPTCH = message.payload[i];
				}
				vce.Extra.PatchType = parseInt(newIndex, 0);
				viewVCE.init();
			}
			index.refreshConnectionStatus();
		});
	},

	setNumOscillators: function (newNum) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		viewVCE_voice.deb_setNumOscillators(newNum);
	},

	deb_setNumOscillators: null,

	raw_setNumOscillators: function (newNum) {
		if (viewVCE.supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/ return; }
		console.log("setNumOscillators: " + newNum);
		if (newNum < 1 || newNum > 16) {
			return;
		}

		let message = {
			"name": "setNumOscillators",
			"payload": {
				"NumOsc": parseInt(newNum, 10),
				"PatchType": parseInt(document.getElementById("patchType").value, 10)
			}
		};
		//index.spinnerOn();
		astilectron.sendMessage(message, function (message) {
			//index.spinnerOff();
			console.log("setNumOscillators returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the boolean
				index.errorNotification(message.payload);
			} else {
				/// now the tricky part - update the in memory version of vce to reflect what just happened:
				vce.Head.VOITAB = newNum - 1
				var oldLength = vce.Envelopes.length

				if (vce.Head.VOITAB <= 0) {
					$("#delOsc").addClass('disabled');
				} else {
					$("#delOsc").removeClass('disabled');
				}
				if (vce.Head.VOITAB >= 15) {
					$("#addOsc").addClass('disabled');
				} else {
					$("#addOsc").removeClass('disabled');
				}
				if (newNum <= oldLength) {
					// nothing to do - just ignored the extra envelopes
				} else {
					for (i = oldLength; i < newNum; i++) {
						// copy the envelope template into the vce:
						// abuse JSON to do a deep copy:
						vce.Envelopes[i] = JSON.parse(JSON.stringify(message.payload.EnvelopeTemplate))
						// overwrite the default patch type 
						vce.Envelopes[i].OPTCH = message.payload.PatchBytes[i]
					}
				}
				viewVCE.init();
			}
			index.refreshConnectionStatus();
		});

	},

	toggleVoicingMode(mode) {
		if (!mode) {
			index.confirmDialog("Disabling Voicing Mode will discard any pending edits. Are you sure?", function () {
				viewVCE_voice.raw_toggleVoicingMode(mode);
			});
		} else {
			viewVCE_voice.raw_toggleVoicingMode(mode);
		}
	},

	raw_toggleVoicingMode: function (mode) {
		console.log(`VoicingMode ${mode ? 'on' : 'off'}`);
		let message = {
			"name": "toggleVoicingMode",
			"payload": mode
		};
		index.spinnerOn();
		astilectron.sendMessage(message, function (message) {
			index.spinnerOff();
			console.log("toggleVoiceMode returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the boolean
				index.errorNotification(message.payload);
			} else {
				viewVCE_voice.voicingMode = mode;
				if (message.payload != null) {
					if (mode) {
						vce = message.payload;

						crt_name = null;
						crt_path = null;
						index.load("viewVCE.html", "content",
							function (ele) {
								viewVCE.init();
							});
					} else {
						// if we just disabled voicing, clear the VCE view
						document.getElementById("content").innerHTML = "";
					}
				}
				index.infoNotification(`Voicing mode ${mode ? 'enabled' : 'disabled'}`);
			}
			index.refreshConnectionStatus();

			if (mode) {
				// reset the SOLO arrays
				for (osc = 0; osc < 16; osc++) {
					viewVCE_voice.MUTE[osc] = false;
					viewVCE_voice.SOLO[osc] = false;
					$('.vceEditToggle').removeClass('on');
				}
			}
			viewVCE_voice.voicingModeVisuals();
		});
	},

	voicingModeVisuals: function () {
		var mode = viewVCE_voice.voicingMode;
		//mode = true; // For debugging and CSS tweaking: force edit controls to be visible

		// XREF: converter mappings : these are duplicated in updateFromCSurface
		if (mode) {
			// CSS for styling the buttons when disabled is HARD.  So avoid it.
			$('.vceNum.spinNullablePatchReg').TouchSpin({
				verticalbuttons: true,
				verticalup: '\u25b4', //'\u25b2',
				verticaldown: '\u25be', //'\u25bc',
				buttonup_txt: '\u25b4', //'\u25b2',
				buttondown_txt: '\u25be', //'\u25bc',
				callback_before_calculation: function (value) {
					return viewVCE_voice.TextToNullablePatchRegister(value);
				},
				callback_after_calculation: function (value) {
					return viewVCE_voice.NullablePatchRegisterToText(value);
				}
			}); $('.vceNum.spinOHARM').TouchSpin({
				verticalbuttons: true,
				verticalup: '\u25b4', //'\u25b2',
				verticaldown: '\u25be', //'\u25bc',
				buttonup_txt: '\u25b4', //'\u25b2',
				buttondown_txt: '\u25be', //'\u25bc',
				callback_before_calculation: function (value) {
					return viewVCE_voice.TextToOHARM(value);
				},
				callback_after_calculation: function (value) {
					return viewVCE_voice.OHARMToText(value);
				}
			});
			$('.vceNum.spinFDETUN').TouchSpin({
				verticalbuttons: true,
				verticalup: '\u25b4', //'\u25b2',
				verticaldown: '\u25be', //'\u25bc',
				buttonup_txt: '\u25b4', //'\u25b2',
				buttondown_txt: '\u25be', //'\u25bc',
				callback_before_calculation: function (value) {
					return viewVCE_voice.TextToFDETUN(value);
				},
				callback_after_calculation: function (value) {
					return viewVCE_voice.FDETUNToText(value);
				}
			});
			$('.vceNum.spinAmpEnv').TouchSpin({
				verticalbuttons: true,
				verticalup: '\u25b4', //'\u25b2',
				verticaldown: '\u25be', //'\u25bc',
				buttonup_txt: '\u25b4', //'\u25b2',
				buttondown_txt: '\u25be', //'\u25bc',
				callback_before_calculation: function (value) {
					return viewVCE_envs.TextToAmpEnvValue(value);
				},
				callback_after_calculation: function (value) {
					return viewVCE_envs.AmpEnvValueToText(value);
				}
			});
			$('.vceNum.spinAmpTime').TouchSpin({
				verticalbuttons: true,
				verticalup: '\u25b4', //'\u25b2',
				verticaldown: '\u25be', //'\u25bc',
				buttonup_txt: '\u25b4', //'\u25b2',
				buttondown_txt: '\u25be', //'\u25bc',
				callback_before_calculation: function (value) {
					return viewVCE_envs.TextToAmpTimeValue(value);
				},
				callback_after_calculation: function (value) {
					return viewVCE_envs.AmpTimeValueToText(value);
				}
			});
			$('.vceNum.spinFreqTime').TouchSpin({
				verticalbuttons: true,
				verticalup: '\u25b4', //'\u25b2',
				verticaldown: '\u25be', //'\u25bc',
				buttonup_txt: '\u25b4', //'\u25b2',
				buttondown_txt: '\u25be', //'\u25bc',
				callback_before_calculation: function (value) {
					return viewVCE_envs.TextToFreqTimeValue(value);
				},
				callback_after_calculation: function (value) {
					return viewVCE_envs.FreqTimeValueToText(value);
				}
			});
			// plain number variant:
			$('.vceNum.spinPLAIN').TouchSpin({
				verticalbuttons: true,
				verticalup: '\u25b4', //'\u25b2',
				verticaldown: '\u25be', //'\u25bc',
				buttonup_txt: '\u25b4', //'\u25b2',
				buttondown_txt: '\u25be', //'\u25bc',
			});
			// make any plain-text spans align:
			$('.spinNOSPIN').addClass("spinNOSPIN-Enabled");

		}

		// Load/Save menu items get disabled/enabled:
		if (mode) {
			$('#disableVRAMMenuItem').addClass('disabled');
			$('#loadCRTMenuItem').addClass('disabled');
			$('#saveVCEMenuItem').removeClass('disabled');
			$('#oscPlusMinus').show();
		} else {
			$('#disableVRAMMenuItem').removeClass('disabled');
			$('#loadCRTMenuItem').removeClass('disabled');
			$('#saveVCEMenuItem').addClass('disabled');
			$('#oscPlusMinus').hide();
		}
		if (vce.Head.VOITAB <= 0) {
			$("#delOsc").addClass('disabled');
		} else {
			$("#delOsc").removeClass('disabled');
		}
		if (vce.Head.VOITAB >= 15) {
			$("#addOsc").addClass('disabled');
		} else {
			$("#addOsc").removeClass('disabled');
		}

		$('.vceEdit').prop('disabled', !mode);
		if (mode) {
			$('.vceEditToggleText').show();
		} else {
			$('.vceEditToggleText').hide();
		}
		document.getElementById("voiceModeButtonImg").src = `static/images/red-button-${viewVCE_voice.voicingMode ? 'on' : 'off'}-full.png`;
	},

	chart: null,

	updateChart: function () {
		var ampData = viewVCE_voice.ampProportionCurve(vce.Head.VACENT, vce.Head.VASENS);
		var timbreData = viewVCE_voice.timbreProportionCurve(vce.Head.VTCENT, vce.Head.VTSENS);

		viewVCE_voice.chart.data.datasets[0].data = ampData;
		viewVCE_voice.chart.data.datasets[1].data = timbreData;
		viewVCE_voice.chart.update();
	},

	updateVibType: function () {
		document.getElementById("vibType").innerHTML = (vce.Head.VIBDEP >= 0) ? "Sine" : "Random";
	},

	init: function (incrementalUpdate) {
		console.log('--- start viewVCE_voice init');

		if (viewVCE_voice.deb_onchange == null) {
			viewVCE_voice.deb_onchange = index.debounceFirstArg(viewVCE_voice.raw_onchange, 50);
		}
		if (viewVCE_voice.deb_setNumOscillators == null) {
			viewVCE_voice.deb_setNumOscillators = _.debounce(viewVCE_voice.raw_setNumOscillators, 250);
		}

		$('#vceTabs a[href="#vceVoiceTab"]').on('shown.bs.tab', function (e) {
			viewVCE_voice.sendToCSurface(null, "voice-tab", 1);
		});
		$('#vceTabs a[href="#vceEnvsTab"]').on('shown.bs.tab', function (e) {
			viewVCE_voice.sendToCSurface(null, "freq-envelopes-tab", 1);
		});
		$('#vceTabs a[href="#vceFiltersTab"]').on('shown.bs.tab', function (e) {
			viewVCE_voice.sendToCSurface(null, "filters-tab", 1);
		});
		$('#vceTabs a[href="#vceKeyEqTab"]').on('shown.bs.tab', function (e) {
			viewVCE_voice.sendToCSurface(null, "keyeq-tab", 1);
		});
		$('#vceTabs a[href="#vceKeyPropTab"]').on('shown.bs.tab', function (e) {
			viewVCE_voice.sendToCSurface(null, "keyprop-tab", 1);
		});

		viewVCE_voice.patchTable();

		console.log("view VCE, CRT:" + crt_name + ", VCE: " + vce.Head.VNAME)

		if (crt_name == null) {
			document.getElementById("backToCRT").hidden = true;
		} else {
			document.getElementById("backToCRT").hidden = false;
		}

		document.getElementById("nOsc").innerHTML = vce.Head.VOITAB + 1;
		document.getElementById("keysPlayable").innerHTML = Math.floor(32 / (vce.Head.VOITAB + 1));
		viewVCE_voice.updateVibType();
		document.getElementById("VIBRAT").value = vce.Head.VIBRAT;
		viewVCE_voice.sendToCSurface(document.getElementById("VIBRAT"), "VIBRAT", vce.Head.VIBRAT)
		document.getElementById("VIBDEL").value = vce.Head.VIBDEL;
		viewVCE_voice.sendToCSurface(document.getElementById("VIBDEL"), "VIBDEL", vce.Head.VIBDEL)
		document.getElementById("VIBDEP").value = vce.Head.VIBDEP;
		viewVCE_voice.sendToCSurface(document.getElementById("VIBDEP"), "VIBDEP", vce.Head.VIBDEP)
		document.getElementById("APVIB").value = vce.Head.APVIB;
		viewVCE_voice.sendToCSurface(document.getElementById("APVIB"), "APVIB", vce.Head.APVIB)

		document.getElementById("VTRANS").value = vce.Head.VTRANS;
		viewVCE_voice.sendToCSurface(document.getElementById("VTRANS"), "VTRANS", vce.Head.VTRANS)
		document.getElementById("VACENT").value = vce.Head.VACENT;
		viewVCE_voice.sendToCSurface(document.getElementById("VACENT"), "VACENT", vce.Head.VACENT)
		document.getElementById("VASENS").value = vce.Head.VASENS;
		viewVCE_voice.sendToCSurface(document.getElementById("VASENS"), "VASENS", vce.Head.VASENS)
		document.getElementById("VTCENT").value = vce.Head.VTCENT;
		viewVCE_voice.sendToCSurface(document.getElementById("VTCENT"), "VTCENT", vce.Head.VTCENT)
		document.getElementById("VTSENS").value = vce.Head.VTSENS;
		viewVCE_voice.sendToCSurface(document.getElementById("VTSENS"), "VTSENS", vce.Head.VTSENS)

		var i;
		var count = 0;
		for (i = 0; i < vce.Head.FILTER.length; i++) {
			if (vce.Head.FILTER[i] != 0) {
				count++;
			}
		}
		document.getElementById("nFilter").innerHTML = count;

		Chart.defaults.global.defaultFontColor = 'white';
		Chart.defaults.global.defaultFontSize = 14;

		var ampData = viewVCE_voice.ampProportionCurve(vce.Head.VACENT, vce.Head.VASENS);
		var timbreData = viewVCE_voice.timbreProportionCurve(vce.Head.VTCENT, vce.Head.VTSENS);

		var ctx = document.getElementById('velocityChart').getContext('2d');
		if (viewVCE_voice.chart != null) {
			// kill off the old chart so we dont get conflicts
			viewVCE_voice.chart.destroy();
		}
		viewVCE_voice.chart = new Chart(ctx, {

			type: 'line',
			data: {
				labels: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''],
				//labels: ['','','',''],
				datasets: [{
					fill: false,
					lineTension: 0,
					pointRadius: 0,
					pointHitRadius: 5,
					label: 'Amplitude',
					backgroundColor: chartColors[0],
					borderColor: chartColors[0],
					data: ampData
				}, {
					fill: false,
					lineTension: 0,
					pointRadius: 0,
					pointHitRadius: 5,
					label: 'Timbre',
					backgroundColor: chartColors[1],
					borderColor: chartColors[1],
					data: timbreData
				}]
			},

			// Configuration options go here
			options: {
				animation: {
					duration: 0
				},
				scales: {
					xAxes: [{
						gridlines: {
							color: '#666',
							display: true,
							drawBorder: true,
							drawOnChartArea: false
						},
						scaleLabel: {
							display: true,
							labelString: "Velocity"
						},
						ticks: {
							display: true
						}
					}],
					yAxes: [{
						gridlines: {
							color: '#666',
							display: true,
							drawBorder: true,
							drawOnChartArea: false
						},
						scaleLabel: {
							display: true,
							labelString: "Proportion"
						},
						ticks: {
							display: true,
							max: 33,
							min: 0,
							callback: function (dataLabel, index) { return ''; }
						}
					}],
				},
				responsive: false,
				maintainAspectRatio: false
			}
		});

		document.getElementById("vce_crt_name").innerHTML = crt_name;
		// do this last to help the uitest to not start testing too soon
		document.getElementById("vce_name").innerHTML = vce.Head.VNAME;
		document.getElementById("VNAME").value = vce.Head.VNAME.replace(/ +$/g, ''); // trim trailing spaces for editing
		console.log('--- finish viewVCE_voice init');
	},

	updateFromCSurface: function (payload) {
		if (!viewVCE_voice.voicingMode) {
			// ignore unless we're voicing
		}

		// special handling for switching tabs -- fieldnames end with "-tab":
		if (payload.Field.search("-tab") >= 0) {
			if (payload.Field === "voice-tab") {
				$('#vceTabs a[href="#vceVoiceTab"]').tab('show');
			} else if (payload.Field === "freq-envelopes-tab") {
				$('#vceTabs a[href="#vceEnvsTab"]').tab('show');
			} else if (payload.Field === "amp-envelopes-tab") {
				$('#vceTabs a[href="#vceEnvsTab"]').tab('show');
			} else if (payload.Field === "filters-tab") {
				$('#vceTabs a[href="#vceFiltersTab"]').tab('show');
			} else if (payload.Field === "keyeq-tab") {
				$('#vceTabs a[href="#vceKeyEqTab"]').tab('show');
			} else if (payload.Field === "keyprop-tab") {
				$('#vceTabs a[href="#vceKeyPropTab"]').tab('show');
			}
			return;
		}


		var ele = document.getElementById(payload.Field)
		if (ele === undefined || ele === null) {
			console.log("updateFromCSurface " + payload.Field + " element not found");
			return
		}
		var value = payload.Value


		// when using MIDI
		// value comes in unscaled (it's a 0-based MIDI value).  
		// Use the min value on the input control to correct for an offset and then use the text 
		// conversion function attached the the touchspin (if any) to turn that into a string)
		//
		// when using OSC, the value is the direct Synergy byte value - no offset
		/*
		if (ele.hasAttribute("min")) {
			var min = parseInt(ele.getAttribute("min"), 10);
			value = value + min;
		}
		*/

		// XREF: converter mappings: it would be nicer to directly query the input element to determine what 
		// sort of touchspin callbacks are associate, if any.  But its not obvous how to do that, so
		// I duplicate some logic here
		var converter = function (value) {
			return value;
		};
		if (ele.classList.contains("spinFreqTime")) {
			converter = viewVCE_envs.FreqTimeValueToText;
		} else if (ele.classList.contains('spinNullablePatchReg')) {
			converter = viewVCE_voice.NullablePatchRegisterToText(value);
		} else if (ele.classList.contains('spinOHARM')) {
			converter = viewVCE_voice.OHARMToText;
		} else if (ele.classList.contains('spinFDETUN')) {
			converter = viewVCE_voice.FDETUNToText;
		} else if (ele.classList.contains('spinAmpEnv')) {
			converter = viewVCE_envs.AmpEnvValueToText;
		} else if (ele.classList.contains('spinAmpTime')) {
			converter = viewVCE_envs.AmpTimeValueToText;
		} else if (ele.classList.contains('spinFreqTime')) {
			converter = viewVCE_envs.FreqTimeValueToText;
		}

		var valueString = converter("" + value)

		//console.log("  updateFromCSurface " + payload.Field + "was " + ele.value);
		if (ele.disabled) {
			//console.log("   disabled!");
			return;
		}
		if (ele.display == "none") {
			//console.log("   hidden/disabled!");
			return;
		}
		if (ele.nodeName == "SELECT") {
			// cycle through the options in each click
			var options = ele.options
			//			console.log("cycle SELECT: currently " + options.selectedIndex + " len: " + options.length);
			var i = options.selectedIndex + 1
			if (i >= options.length) {
				i = 0
			}
			options.selectedIndex = i;
			//			console.log("cycle SELECT: now " + options.selectedIndex);
			valueString = options[i].text;
		} else if (ele.nodeName == "SPAN") {
			// SOLO/MUTE buttons
			ele.onclick();

			valueString - ele.classList.contains("on") ? "ON" : "OFF"
			return valueString; // don't trigger non-existent onchange()
		} else if (ele.type == "checkbox") {
			ele.checked = value == 0 ? "" : "checked";
		} else if (ele.type == "text") {
			ele.value = valueString;
		}
		//console.log("  updateFromCSurface " + payload.Field + "NOW " + ele.value);
		ele.onchange();
		return valueString;
	},

	sendToCSurface: function (ele, field, value) {
		// pass ele==null to force the value to just be sent without scaling

		// when using MIDI
		// value comes in scaled to "synergy byte" value - but MIDI values are always 0 based.  
		// when using OSC, no need to scale the offset
		// Use the min value on the input control to correct for an offset 
		/*
		if (ele != null && ele.hasAttribute("min")) {
			var min = parseInt(ele.getAttribute("min"), 10);
			value = value - min;
		}
		*/

		let message = {
			"name": "sendToCSurface",
			"payload": {
				Field: field,
				Value: parseInt(value, 10)
			}
		};
		astilectron.sendMessage(message, function (message) {
			//console.log("sendToCSurface returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				index.errorNotification(message.payload);
			}
		});
	}
};
