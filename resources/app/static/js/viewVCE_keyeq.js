let viewVCE_keyeq = {
	chart: null,

	keyEqCurve: function (keq) {
		var result = [];
		// y = -24..6
		// x = 0..23
		for (v = 0; v < keq.length; v++) {
			result[v] = keq[v];
		}
		return result;
	},

	onchange: function (ele) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		viewVCE_keyeq.deb_onchange(ele);
	},

	deb_onchange: null, // initialized during init()

	raw_onchange: function (ele) {
		if (viewVCE.supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/ return; }
		var value = index.checkInputElementValue(ele);
		if (value == undefined) {
			return;
		}

		var id = ele.id;
		console.log("changed: " + id + " val: " + ele.value);

		var eleIndex;
		var pattern = /keyeq\[(\d+)\]/;
		if (ret = id.match(pattern)) {
			eleIndex = parseInt(ret[1])
		}
		let message = {
			"name": "setVoiceVEQEle",
			"payload": {
				"Index": eleIndex,
				"Value": value
			}
		};
		astilectron.sendMessage(message, function (message) {
			//console.log("setVoiceVEQEle returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the value
				index.errorNotification(message.payload);
				return false;
			} else {
				vce.Head.VEQ[eleIndex - 1] = value;
				viewVCE_keyeq.init(true);
				viewVCE_voice.sendToCSurface(ele, id, value);
			}
		});
		return true;
	},

	init: function (incrementalUpdate) {
		//console.log('--- start viewVCE_keyeq init ' + incrementalUpdate);
		if (viewVCE_keyeq.deb_onchange == null) {
			//viewVCE_keyeq.deb_onchange = _.debounce(viewVCE_keyeq.raw_onchange, DEBOUNCE_WAIT_SHORT);
			viewVCE_keyeq.deb_onchange = viewVCE_keyeq.raw_onchange;
		}

		var propData = viewVCE_keyeq.keyEqCurve(vce.Head.VEQ);

		$('#keyEqTable td.val input').each(function (i, obj) {
			var id = obj.id;
			// id is "keyeq[<n>]" - we need the <n> part
			var idxString = id.substring(6);
			var idx = parseInt(idxString, 10) - 1;

			obj.value = propData[idx];

			if (!incrementalUpdate) {
				viewVCE_voice.sendToCSurface(obj, id, propData[idx]);
			}
		});
		if (viewVCE_keyeq.chart != null) {
			viewVCE_keyeq.chart.destroy();
		}

		var ctx = document.getElementById('keyEqChart').getContext('2d');
		viewVCE_keyeq.chart = new Chart(ctx, {

			type: 'line',
			data: {
				labels: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''],
				//labels: ['','','',''],
				datasets: [{
					fill: false,
					lineTension: 0,
					pointRadius: 0,
					label: 'Key Equalization',
					backgroundColor: chartColors[0],
					borderColor: chartColors[0],
					data: propData
				}]
			},

			// Configuration options go here
			options: {
				animation: {
					duration: 0
				},

				tooltips: {
					mode: 'index',
				},
				hover: {
					mode: 'index',
				},
				scales: {
					xAxes: [{
						id: 'x-axis',
						gridLines: {
							color: '#666',
							display: true,
							drawBorder: false,
							drawOnChartArea: false
						},
						scaleLabel: {
							display: true,
							labelString: "Key"
						},
						ticks: {
							min: 0,
							max: 23,
							color: '#eee',
							display: true
						}
					}],
					yAxes: [{
						id: 'y-axis',
						grid: {
							color: '#666'
						},
						gridLines: {
							color: '#666',
							display: true,
							drawBorder: false,
							drawOnChartArea: true
						},
						scaleLabel: {
							display: true,
							labelString: "dB"
						},
						ticks: {
							min: -24,
							max: 8,
							stepSize: 4,
							color: '#eee',
							display: true
						}
					}],
				},
				responsive: false,
				maintainAspectRatio: false
			},

		});

		if (viewVCE_voice.voicingMode) {
			viewVCE_chartdrag.init(viewVCE_keyeq.chart, viewVCE_keyeq.onchange, 'keyeq', 0, 23, -24, 8);
		}

		//console.log('--- finish viewVCE_keyeq init');
	}
};
