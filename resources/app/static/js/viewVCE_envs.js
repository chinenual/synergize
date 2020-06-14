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

	// SYNHCS COMMON.Z80 FTAB:
	ftab: [
		0, 2, 4, 6, 8, 10, 12, 14,	//Reference frequency table
		15, 16, 17, 18, 19, 20, 21, 22,
		24, 25, 27, 28, 30, 32, 34, 36,
		38, 40, 43, 45, 48, 51, 54, 57,
		61, 64, 68, 72, 76, 81, 86, 91,
		96, 102, 108, 115, 122, 129, 137, 145,
		153, 163, 172, 183, 193, 205, 217, 230,
		244, 258, 274, 290, 307, 326, 345, 366,
		387, 411, 435, 461, 488, 517, 548, 581,
		615, 652, 691, 732, 775, 822, 870, 922,
		977, 1035, 1097, 1162, 1231, 1304, 1382, 1464,
		1551, 1644, 1741, 1845, 1955, 2071, 2194, 2325,
		2463, 2609, 2765, 2929, 3103, 3288, 3483, 3691,
		3910, 4143, 4389, 4650, 4926, 5219, 5530, 5859,
		6207, 6576, 6967, 7382, 7820, 8286, 8778, 9300,
		9853, 10439, 11060, 11718, 12414, 13153, 13935, 14764],

	scaleViaRtab(v) {
		if (v <= 0) return 0;
		if (v >= viewVCE_envs.ftab.length) return viewVCE_envs.ftab[viewVCE_envs.ftab.length - 1];
		return viewVCE_envs.ftab[v - 1];
	},

	// Freq values:
	//   as displayed: -61 .. 63
	//   byte range:   0xc3 .. 0x3f
	scaleFreqEnvValue: function (v) {
		// See OSCDSP.Z80 DISVAL: DVAL10:
		return v; // TODO
	},

	// Amp values:
	//   as displayed: 0 .. 72
	//   byte range:   0x37 .. 0x7f
	scaleAmpEnvValue: function (v, last) {
		// See OSCDSP.Z80 DISVAL: DVAL30:
		if (last) return 0;
		return Math.max(0, v - 55);
	},

	// Freq Time values:
	//   as displayed: 0 .. 29528
	//   byte range:   0x0 .. 0x54
	scaleFreqTimeValue: function (v, first) {
		// See OSCDSP.Z80 DISVAL:
		if (first) return 0;
		if (v < 15) return v;
		return viewVCE_envs.scaleViaRtab((2 * v) - 14);
	},

	// Freq Time values:
	//   as displayed: 0 .. 6576
	//   byte range:   0x0 .. 0x54
	scaleAmpTimeValue: function (v) {
		// See OSCDSP.Z80 DISVAL: DVAL20:
		if (v < 39) return v;
		return viewVCE_envs.scaleViaRtab((v * 2) - 54);
	},

	envChartUpdate: function (oscNum, envNum) {
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

		for (i = 0; i < envelopes.FreqEnvelope.NPOINTS; i++) {
			var tr = $('#envTable tbody tr:eq(' + i + ')');

			// table is logically in groups of 4
			var freqLow = viewVCE_envs.scaleFreqEnvValue(envelopes.FreqEnvelope.Table[i * 4 + 0]);
			var freqUp = viewVCE_envs.scaleFreqEnvValue(envelopes.FreqEnvelope.Table[i * 4 + 1]);
			var timeLow = viewVCE_envs.scaleFreqTimeValue(envelopes.FreqEnvelope.Table[i * 4 + 2], i == 0);
			var timeUp = viewVCE_envs.scaleFreqTimeValue(envelopes.FreqEnvelope.Table[i * 4 + 3], i == 0);

			lastFreqLow = freqLow;
			lastFreqUp = freqUp;
			totalTimeLow += timeLow;
			totalTimeUp += timeUp;

			datasets[freqLowIdx].data.push({ x: totalTimeLow, y: freqLow });
			datasets[freqUpIdx].data.push({ x: totalTimeUp, y: freqUp });

			document.getElementById(`envFreqLowVal[${i + 1}]`).value = freqLow;
			document.getElementById(`envFreqUpVal[${i + 1}]`).value = freqUp;
			document.getElementById(`envFreqLowTime[${i + 1}]`).value = timeLow;
			document.getElementById(`envFreqUpTime[${i + 1}]`).value = timeUp;
			document.getElementById(`envFreqTotLowTime[${i + 1}]`).innerHTML = totalTimeLow;
			document.getElementById(`envFreqTotUpTime[${i + 1}]`).innerHTML = totalTimeUp;

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
		for (i = 0; i < envelopes.AmpEnvelope.NPOINTS; i++) {
			var tr = $('#envTable tbody tr:eq(' + i + ')');

			// table is logically in groups of 4.
			// "j" accounts for the difference in column index due to the
			// row-spanning separators (only in i==0):
			j = (i == 0) ? 11 : 9;
			console.log("j " + j);
			//	    console.dir(tr);
			//	    console.dir(tr.find('td:eq(' +(j+0)+ ')'));
			var ampLow =
				viewVCE_envs.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i * 4 + 0],
					(i + 1) >= envelopes.AmpEnvelope.NPOINTS);
			var ampUp =
				viewVCE_envs.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i * 4 + 1],
					(i + 1) >= envelopes.AmpEnvelope.NPOINTS);
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

	}

};
