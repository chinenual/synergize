let viewVCE_envs = {

	chart: null,

	init: function () {
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

		viewVCE_envs.envChartUpdate(1, -1)
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
		return '' + viewVCE_envs.scaleAmpEnvValue(v);
	},
	TextToAmpEnvValue(v) {
		return '' + viewVCE_envs.unscaleAmpEnvValue(parseInt(v, 10));
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
		return '' + viewVCE_envs.unscaleFreqTimeValue(parseInt(v, 10));
	},


	// Freq Time values:
	//   as displayed: 0 .. 6576
	//   byte range:   0x0 .. 0x54 (0 .. 84)
	scaleAmpTimeValue: function (v) {
		// See OSCDSP.Z80 DISVAL: DVAL20: which is:
		// if (v < 39) return v;
		// return viewVCE_envs.scaleViaRtab((v * 2) - 54);
		return viewVCE_envs.ampTimeScale[v];

	},
	unscaleAmpTimeValue: function (v) {
		// fixme: linear search is brute force - but the list is short - performance is "ok" as is...
		for (var i = 0; i < viewVCE_envs.ampTimeScale.length; i++) {
			if (viewVCE_envs.ampTimeScale[i] >= v) {
				return i;
			}
		}
		// shouldnt happen!
		return viewVCE_envs.ampTimeScale.length - 1;
	},
	AmpTimeValueToText(v) {
		return '' + viewVCE_envs.scaleAmpTimeValue(v);
	},
	TextToAmpTimeValue(v) {
		return '' + viewVCE_envs.unscaleAmpTimeValue(parseInt(v, 10));
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
				arr: [[0, 0], [10, 10], [15, 15], [16, 25], [54, 2071], [75, 23436], [76, 26306], [77, 29528]],
				name: "freqTimeValue",
				func: viewVCE_envs.scaleFreqTimeValue,
			},
			{
				arr: [[0, 0], [20, 20], [40, 40], [41, 45], [54, 205], [75, 2325], [76, 2609], [83, 5859], [84, 6576]],
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

	logConversionFunctions: function () {
		var ok = true;
		var arr = [];
		for (var i = -61; i <= 63; i++) {
			var scaled = viewVCE_envs.scaleFreqEnvValue(i);
			arr.push(scaled);
		}
		console.log(" FreqEnvValue -16..63: " + JSON.stringify(arr));

		arr = [];
		for (var i = 55; i <= 127; i++) {
			var scaled = viewVCE_envs.scaleAmpEnvValue(i);
			arr.push(scaled);
		}
		console.log(" AmpEnvValue 55..127: " + JSON.stringify(arr));
		arr = [];
		for (var i = 0; i <= 84; i++) {
			var scaled = viewVCE_envs.scaleFreqTimeValue(i);
			arr.push(scaled);
		}
		console.log(" FreqTimeValue 0..84: " + JSON.stringify(arr));
		arr = [];
		for (var i = 0; i <= 84; i++) {
			var scaled = viewVCE_envs.scaleAmpTimeValue(i);
			arr.push(scaled);
		}
		console.log(" AmpTimeValue 0..84: " + JSON.stringify(arr));
	},

	supressOnChange: false,

	onchange: function (ele) {
		if (viewVCE.supressOnchange) { return; }
		if (viewVCE_envs.supressOnchange) { return; }

		var eleIndex;
		var envOscSelectEle = document.getElementById("envOscSelect");
		var osc = parseInt(envOscSelectEle.value, 10); // one-based osc index
		var envEnvSelectEle = document.getElementById("envOscSelect");
		var selectedEnv = parseInt(envEnvSelectEle.value, 10);

		console.log("env ele change " + ele.id + " " + ele.value);
		var pattern = /([A-Za-z]+)\[(\d+)\]/;
		var funcName;
		var eleIndex;
		var value;
		if (ret = ele.id.match(pattern)) {
			fieldType = ret[1];
			funcName = 'set' + fieldType.charAt(0).toUpperCase() + fieldType.slice(1);
			eleIndex = parseInt(ret[2])
			value = parseInt(ele.value, 10);
			// now scale the value to the byte value the synergy wants to see:
			switch (fieldType) {
				case "envFreqLowVal":
				case "envFreqUpVal":
					value = viewVCE_envs.unscaleFreqEnvValue(value);
					break;
				case "envFreqLowTime":
				case "envFreqUpTime":
					value = viewVCE_envs.unscaleFreqTimeValue(value);
					break;
				case "envAmpLowVal":
				case "envAmpUpVal":
					value = viewVCE_envs.unscaleAmpEnvValue(value);
					break;
				case "envAmpLowTime":
				case "envAmpUpTime":
					value = viewVCE_envs.unscaleAmpTimeValue(value);
					break;
			}
		}

		let message = {
			"name": funcName,
			"payload": {
				"Osc": osc,
				"Index": eleIndex,
				"Value": value
			}
		};
		astilectron.sendMessage(message, function (message) {
			console.log(funcName + " returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the value
				index.errorNotification(message.payload);
				return false;
			} else {
				viewVCE_envs.envChartUpdate(osc, selectedEnv);
			}
		});
		return true;

	},

	changeEnvPoints: function (whichEnv, increment) {
		var eleIndex;
		var envOscSelectEle = document.getElementById("envOscSelect");
		var osc = parseInt(envOscSelectEle.value, 10); // one-based osc index
		var envEnvSelectEle = document.getElementById("envOscSelect");
		var selectedEnv = parseInt(envEnvSelectEle.value, 10);
		var envs = vce.Envelopes[osc - 1];

		var changed = false;
		if (whichEnv === 'freq') {
			newlen = envs.FreqEnvelope.NPOINTS + increment;
			if (newlen >= 1 && newlen <= 16) {
				envs.FreqEnvelope.NPOINTS = newlen;
				changed = true;
			}
		} else {
			newlen = envs.AmpEnvelope.NPOINTS + increment;
			if (newlen >= 1 && newlen <= 16) {
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
				console.log("setOscEnvLengths returned: " + JSON.stringify(message));
				// Check error
				if (message.name === "error") {
					// failed - dont change the value
					index.errorNotification(message.payload);
					return false;
				} else {
					viewVCE_envs.envChartUpdate(osc, selectedEnv);
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

	envChartUpdate: function (oscNum, envNum) {
		viewVCE_envs.supressOnchange = true;

		viewVCE_envs.uncompressEnvelopes();

		var oscIndex = oscNum - 1;
		var envelopes = vce.Envelopes[oscIndex];

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
				showLine: true,
				borderWidth: 3,
				backgroundColor: chartColors[1],
				borderColor: chartColors[1],
				data: []
			},
		];

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

		console.log("env freq npt " + envelopes.FreqEnvelope.NPOINTS);
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
		for (i = envelopes.FreqEnvelope.NPOINTS; i < 16; i++) {
			// hide unused rows
			var tr = $('#envTable tbody tr:eq(' + i + ')');
			console.log("hide row " + i);

			$(`#envFreqLoop\\[${i + 1}\\]`).hide();
			$(`#envFreqLowVal\\[${i + 1}\\]`).hide();
			$(`#envFreqUpVal\\[${i + 1}\\]`).hide();
			$(`#envFreqLowTime\\[${i + 1}\\]`).hide();
			$(`#envFreqUpTime\\[${i + 1}\\]`).hide();

		}
		for (i = 0; i < envelopes.FreqEnvelope.NPOINTS; i++) {
			var tr = $('#envTable tbody tr:eq(' + i + ')');
			console.log("show row " + i);

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

			console.log(`freq raw row ${i + 1} ${envelopes.FreqEnvelope.Table[i * 4 + 0]} ${envelopes.FreqEnvelope.Table[i * 4 + 1]} ${envelopes.FreqEnvelope.Table[i * 4 + 2]} ${envelopes.FreqEnvelope.Table[i * 4 + 3]}`);
			console.log(`freq scaled row ${i + 1} ${freqLow} ${freqUp} ${timeLow} ${timeUp}`);


			if (envelopes.FreqEnvelope.SUSTAINPT == (i + 1)) {
				$(`#envFreqLoop\\[${i + 1}\\] option[value='S']`).prop('selected', true);
			}
			if (envelopes.FreqEnvelope.LOOPPT == (i + 1)) {
				$(`#envFreqLoop\\[${i + 1}\\] option[value='L']`).prop('selected', true);
			}
		}
		var maxTotalTime = Math.max(totalTimeLow, totalTimeUp);

		console.log("env amp npt " + envelopes.AmpEnvelope.NPOINTS);
		totalTimeLow = 0;
		totalTimeUp = 0;


		for (i = envelopes.FreqEnvelope.NPOINTS; i < 16; i++) {
			// hide unused rows
			console.log("hide row " + i);

			$(`#envAmpLoop\\[${i + 1}\\]`).hide();
			$(`#envAmpLowVal\\[${i + 1}\\]`).hide();
			$(`#envAmpUpVal\\[${i + 1}\\]`).hide();
			$(`#envAmpLowTime\\[${i + 1}\\]`).hide();
			$(`#envAmpUpTime\\[${i + 1}\\]`).hide();

		}
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
			console.log("j " + j);
			//	    console.dir(tr);
			//	    console.dir(tr.find('td:eq(' +(j+0)+ ')'));
			var isLast = (i + 1) >= envelopes.AmpEnvelope.NPOINTS;
			var ampLow = viewVCE_envs.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i * 4 + 0]);
			var ampUp = viewVCE_envs.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i * 4 + 1]);
			var timeLow = viewVCE_envs.scaleAmpTimeValue(envelopes.AmpEnvelope.Table[i * 4 + 2]);
			var timeUp = viewVCE_envs.scaleAmpTimeValue(envelopes.AmpEnvelope.Table[i * 4 + 3]);
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
			} else {
				document.getElementById(`envAmpLowVal[${i + 1}]`).disabled = false;
				document.getElementById(`envAmpUpVal[${i + 1}]`).disabled = false;
			}

			if (envelopes.AmpEnvelope.SUSTAINPT == (i + 1)) {
				$(`#envAmpLoop\\[${i + 1}\\] option[value='S']`).prop('selected', true);
			}
			if (envelopes.AmpEnvelope.LOOPPT == (i + 1)) {
				$(`#envAmpLoop\\[${i + 1}\\] option[value='L']`).prop('selected', true);
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

		var filteredDatasets = [];
		if (envNum < 0) {
			// all of them:
			filteredDatasets = datasets;
		} else {
			filteredDatasets.push(datasets[envNum])
		}
		var ctx = document.getElementById('envChart').getContext('2d');
		if (viewVCE_envs.chart != null) {
			viewVCE_envs.chart.destroy();
		}
		viewVCE_envs.chart = new Chart(ctx, {

			type: 'line',
			data: {
				//		labels: ['','','','','','','','','','', '','','','','','','','','','', '','','','','','','','','','','',''],
				datasets: filteredDatasets
			},

			// Configuration options go here
			options: {
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
						type: 'logarithmic',
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
						ticks: {
							callback: function (value, index, values) {
								// don't use scientific notation
								return value;
							},
							color: '#666',
							display: true
						}
					}],
					yAxes: [{
						position: 'left',
						id: 'freq-axis',
						type: 'logarithmic',
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
						ticks: {
							callback: function (value, index, values) {
								// don't use scientific notation
								return value;
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
						ticks: {
							color: '#eee',
							display: true
						}
					}],
				},
				responsive: false,
				maintainAspectRatio: false
			}
		});
		viewVCE_envs.supressOnchange = false;
	}

};
