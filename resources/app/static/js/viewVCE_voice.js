let viewVCE_voice = {
	voicingMode: false,

	timbreProportionCurve: function (center, sensitivity) {
		console.log("timbre prop " + center + " " + sensitivity)
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
		console.log("amp prop " + center + " " + sensitivity)
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

	onchange: function (ele, updater) {
		var id = ele.id;
		console.log("changed: " + id);

		var param
		var args
		var waveKeyPattern = /wk([A-Z]+)\[(\d+)\]/;
		var oscPattern = /([A-Z]+)\[(\d+)\]/;
		var headPattern = /([A-Z]+)/;
		if (ret = id.match(waveKeyPattern)) {
			param = ret[1]
			osc = parseInt(ret[2], 10)
			if (param == "WAVE") {
				funcname = "setOscWAVE"
				args = [osc, ele.value == "Sin" ? 0 : 1];
			} else {
				funcname = "setOscKEYPROP"
				args = [osc, ele.value ? 1 : 0];
			}
		} else if (ret = id.match(oscPattern)) {
			param = ret[1];
			osc = parseInt(ret[2], 10)
			args = [osc, parseInt(ele.value, 10)]

			console.log("changed: " + id + " param: " + param + " osc: " + osc);
			vce.Envelopes[osc - 1][param] = ele.value;
			funcname = "setVoiceByte"

		} else if (ret = id.match(headPattern)) {
			param = id;
			args = [parseInt(ele.value, 10)]
			vce.Head[param] = ele.value;
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
				console.log("updater: " + updater);
				// Check error
				if (message.name === "error") {
					// failed - dont change the boolean
					index.errorNotification(message.payload);
					return false;
				} else if (updater != undefined) {
					updater();
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
		console.log("vceEnv" + JSON.stringify(vce.Envelopes));
		console.log("vceFilters" + JSON.stringify(vce.Filters));

		// populate new ones:
		for (osc = 0; osc <= vce.Head.VOITAB; osc++) {
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

			// FIXME: assumes envelopes are sorted in oscillator order
			var patchByte = vce.Envelopes[osc].FreqEnvelope.OPTCH;
			var patchInhibitAddr = (patchByte & 0x20) != 0;
			var patchInhibitF0 = (patchByte & 0x04) != 0;
			var patchOutputDSR = ((patchByte & 0xc0) >> 6);
			var patchAdderInDSR = ((patchByte & 0x18) >> 3);
			var patchFOInputDSR = (patchByte & 0x03);

			console.log(osc + " patch byte: " + patchByte + "\n" +
				" patchInhibitAddr : " + patchInhibitAddr + "\n" +
				" patchInhibitF0   : " + patchInhibitF0 + "\n" +
				" patchOutputDSR   : " + patchOutputDSR + "\n" +
				" patchAdderInDSR  : " + patchAdderInDSR + "\n" +
				" patchFOInputDSR  : " + patchFOInputDSR + "\n");

			//--- Patch F
			td = document.createElement("td");
			if (!patchInhibitF0) {
				td.innerHTML = patchFOInputDSR + 1;
			}
			tr.appendChild(td);

			//--- Patch A
			td = document.createElement("td");
			if (!patchInhibitAddr) {
				td.innerHTML = patchAdderInDSR + 1;
			}
			tr.appendChild(td);

			//--- Patch O
			td = document.createElement("td");
			td.innerHTML = patchOutputDSR + 1;
			tr.appendChild(td);

			//--- Hrm
			td = document.createElement("td");
			td.innerHTML = `<input type="number" class="vceEdit vceNum" id="OHARM[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this)" value="${vce.Envelopes[osc].FreqEnvelope.OHARM}" 
			min="-11" max="30"
			disabled/>`;
			tr.appendChild(td);

			//--- Detn
			td = document.createElement("td");
			td.innerHTML = `<input type="number" class="vceEdit vceNum" id="FDETUN[${osc + 1}]" 
			onchange="viewVCE_voice.onchange(this)" value="${vce.Envelopes[osc].FreqEnvelope.FDETUN}" 
			min="-127" max="128"
			disabled/>`;
			tr.appendChild(td);

			var waveByte = vce.Envelopes[osc].FreqEnvelope.Table[0][3];
			var wave = ((waveByte & 0x1) == 0) ? 'Sin' : 'Tri';
			var keyprop = ((waveByte & 0x10) == 0) ? true : false;

			//--- Wave
			td = document.createElement("td");
			td.innerHTML = wave;
			td.innerHTML = `<select class="vceEdit" id="wkWAVE[${osc + 1}]" value="${keyprop}" 
			onchange="viewVCE_voice.onchange(this)" disabled/>
			<option value="Sin">Sin</option>
			<option value="Tri">Tri</option>
			</select>
			`;
			tr.appendChild(td);

			//--- Key
			td = document.createElement("td");
			// can't use disabled attr - bootstrap styling hides it - use javascript hack to make it readonnly
			td.innerHTML = `<input type="checkbox" id="wkKEYPROP[${osc + 1}]" value="true" 
			${keyprop ? " checked " : ""} 
			onchange="viewVCE_voice.voicingMode ? viewVCE_voice.onchange(this) : (this.checked=!this.checked)"/>`;
			tr.appendChild(td);

			//--- Flt
			td = document.createElement("td");
			var filter = vce.Head.FILTER[osc];
			td.innerHTML =
				(filter == 0) ? ''
					: (filter > 0) ? ('Bf ' + filter) : ('Af ' + -filter);
			tr.appendChild(td);

			tbody.appendChild(tr);
		}
	},

	changePatchType: function (newIndex) {
		let message = {
			"name": "setPatchType",
			"payload": parseInt(newIndex, 10)
		};
		//index.spinnerOn();
		astilectron.sendMessage(message, function (message) {
			//index.spinnerOff();
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

	toggleVoicingMode: function (mode) {
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
			//			viewVCE_voice.voicingModeVisuals();
		});
	},

	voicingModeVisuals: function () {
		console.log("voicingModeVisuals " + viewVCE_voice.voicingMode);
		$('.vceEdit').prop('disabled', !viewVCE_voice.voicingMode);
		if (viewVCE_voice.voicingMode) {
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

	init: function () {
		console.log("vceVoiceTab init");

		viewVCE_voice.patchTable();

		if (crt_name == null) {
			document.getElementById("backToCRT").hidden = true;
		} else {
			document.getElementById("backToCRT").hidden = false;
		}

		document.getElementById("vce_crt_name").innerHTML = crt_name;
		document.getElementById("vce_name").innerHTML = vce.Head.VNAME;

		document.getElementById("name").innerHTML = vce.Head.VNAME;
		document.getElementById("nOsc").value = vce.Head.VOITAB + 1;
		document.getElementById("keysPlayable").innerHTML = Math.floor(32 / (vce.Head.VOITAB + 1));
		viewVCE_voice.updateVibType();
		document.getElementById("VIBRAT").value = vce.Head.VIBRAT;
		document.getElementById("VIBDEL").value = vce.Head.VIBDEL;
		document.getElementById("VIBDEP").value = vce.Head.VIBDEP;
		document.getElementById("APVIB").value = vce.Head.APVIB;

		document.getElementById("VTRANS").value = vce.Head.VTRANS;
		document.getElementById("VACENT").value = vce.Head.VACENT;
		document.getElementById("VASENS").value = vce.Head.VASENS;
		document.getElementById("VTCENT").value = vce.Head.VTCENT;
		document.getElementById("VTSENS").value = vce.Head.VTSENS;
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

		console.log("ampData: " + JSON.stringify(ampData));
		console.log("timbreData: " + JSON.stringify(timbreData));

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
	}
};
