let viewVCE_voice = {
	voicingMode: false,

	timbreProportionCurve: function (center, sensitivity) {
		var result = [];
		if (sensitivity == 0) {
			for (v = 0; v < 32; v++) {
				result[v] = center;
			}
		}
		// center = 0..32
		// sensitivity = 1..31
		for (v = 0; v < 32; v++) {
			// this appears to be what the z80 code is doing for timbre PROPC:
			// HACK: the division by 10 is a hack - makes the graphs "look" ok but i can't account for why
			var p = (center * 2) - 15 + (v / 10 * sensitivity) - (2 * sensitivity);
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
		}
		// center = 0..32
		// sensitivity = 1..31
		for (v = 0; v < 32; v++) {
			// this appears to be what the z80 code is doing for timbre PROPC:
			// HACK: the division by 10 is a hack - makes the graphs "look" ok but i can't account for why
			var p = (((v / 10 * sensitivity) / 2 - sensitivity) / 2) + center - 24
			if (p > 6) p = 6;
			if (p < -24) p = -24;
			result[v] = p + 25;
		}
		return result;
	},

	onchange: function (ele) {
		var id = ele.id;
		console.log("changed: " + id);

		var oscPattern = /([A-Z]+)\[(\d+)\]/;
		if (ret = id.match(oscPattern)) {
			param = ret[1];
			osc = ret[2];
			console.log("changed: " + id + " param: " + param + " osc: " + osc);

			let message = {
				"name": "setVoiceByte",
				"payload": {
					"Param": param,
					"Args": [parseInt(osc, 10), parseInt(ele.value, 10)]
				}
			};
			astilectron.sendMessage(message, function (message) {
				console.log("setVoiceByte returned: " + JSON.stringify(message));
				// Check error
				if (message.name === "error") {
					// failed - dont change the boolean
					index.errorNotification(message.payload);
				}
			});

		}
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
			td.innerHTML = `<input type="number" class="vceEdit vceNum" id="OHARM[${osc + 1}]" onchange="viewVCE_voice.onchange(this)" value="${vce.Envelopes[osc].FreqEnvelope.OHARM}" disabled/>`;
			tr.appendChild(td);

			//--- Detn
			td = document.createElement("td");
			td.innerHTML = `<input type="number" class="vceEdit vceNum" id="FDETUN[${osc + 1}]" onchange="viewVCE_voice.onchange(this)" value="${vce.Envelopes[osc].FreqEnvelope.FDETUN}" disabled/>`;
			tr.appendChild(td);

			var waveByte = vce.Envelopes[osc].FreqEnvelope.Table[0][3];
			var wave = ((waveByte & 0x7) == 0) ? 'Sin' : 'Tri';
			var keyprop = ((waveByte & 0x8) == 0) ? 'Key' : '';

			//--- Wave
			td = document.createElement("td");
			td.innerHTML = wave;
			td.innerHTML = `<select class="vceEdit" id="wave[${osc + 1}]" value="${keyprop}" disabled/>
			<option value="Sin">Sin</option>
			<option value="Tri">Tri</option>
			</select>
			`;
			tr.appendChild(td);

			//--- Key
			td = document.createElement("td");
			td.innerHTML = `<input type="text" class="vceEdit vceNum" id="keyprop[${osc + 1}]" value="${keyprop}" disabled/>`;
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
					vce = message.payload;

					crt_name = null;
					crt_path = null;
					index.load("viewVCE.html", document.getElementById("content"));
					viewVCE.init();
				}
				index.infoNotification(`Voicing mode ${mode ? 'enabled' : 'disabled'}`);
			}
			index.refreshConnectionStatus();

			$('.vceEdit').prop('disabled', !viewVCE_voice.voicingMode);
			document.getElementById("voiceModeButtonImg").src = `static/images/red-button-${viewVCE_voice.voicingMode ? 'on' : 'off'}-full.png`;
		});
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
		document.getElementById("vibType").value = (vce.Head.VIBDEL >= 0) ? "Sine" : "Random";
		document.getElementById("vibRate").value = vce.Head.VIBRAT;
		document.getElementById("vibDelay").value = vce.Head.VIBDEL;
		document.getElementById("vibDepth").value = vce.Head.VIBDEP;
		document.getElementById("transpose").value = vce.Head.VTRANS;
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
		var chart = new Chart(ctx, {

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
