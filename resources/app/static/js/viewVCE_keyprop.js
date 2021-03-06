let viewVCE_keyprop = {
	chart: null,

	keyPropCurve: function (kprop) {
		var result = [];
		// y = 0..32
		// x = 0..23
		for (v = 0; v < kprop.length; v++) {
			result[v] = kprop[v];
		}
		return result;
	},

	onchange: function (ele, updateChart) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		viewVCE_keyprop.deb_onchange(ele, updateChart);
	},

	deb_onchange: null, // initialized during init()

	raw_onchange: function (ele, updateChart) {
		if (viewVCE.supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/ return; }
		var id = ele.id;

		var value = index.checkInputElementValue(ele);
		if (value === undefined) {
			return;
		}

		console.log("changed: " + id + " val: " + value);

		var eleIndex;
		var pattern = /keyprop\[(\d+)\]/;
		if (ret = id.match(pattern)) {
			eleIndex = parseInt(ret[1])
		}
		let message = {
			"name": "setVoiceKPROPEle",
			"payload": {
				"Index": eleIndex,
				"Value": value
			}
		};
		astilectron.sendMessage(message, function (message) {
			//console.log("setVoiceKPROPEle returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the value
				index.errorNotification(message.payload);
				return false;
			} else {
				vce.Head.KPROP[eleIndex - 1] = value;
				if (updateChart) {
					viewVCE_keyprop.init(true);
				}
				viewVCE_voice.sendToCSurface(ele, id, value);
			}
		});
		return true;
	},

	init: function (incrementalUpdate) {
		console.log('--- start viewVCE_keyprop init ' + incrementalUpdate);
		if (viewVCE_keyprop.deb_onchange == null) {
			//viewVCE_keyprop.deb_onchange = index.debounceFirstArg(viewVCE_keyprop.raw_onchange, DEBOUNCE_WAIT_SHORT);
			viewVCE_keyprop.deb_onchange = viewVCE_keyprop.raw_onchange;
		}

		var propData = viewVCE_keyprop.keyPropCurve(vce.Head.KPROP);

		$('#keyPropTable td.val input').each(function (i, obj) {
			var id = obj.id;

			// id is "keyprop[<n>]" - we need the <n> part
			var idxString = id.substring(8);
			var idx = parseInt(idxString, 10) - 1;

			obj.value = '' + propData[idx];

			if (!incrementalUpdate) {
			   viewVCE_voice.sendToCSurface(obj, id, propData[idx]);
			}
		});
		if (viewVCE_keyprop.chart != null) {
			viewVCE_keyprop.chart.destroy();
		}
		var ctx = document.getElementById('keyPropChart').getContext('2d');
		viewVCE_keyprop.chart = new Chart(ctx, {

			type: 'line',

			data: {

				labels: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''],
				datasets: [{
					//					mouseover: function (e) {
					//						console.log("mouseover: " + JSON.stringify(e));
					//					},
					//					mousemove: function (e) {
					//						console.log("mousemove: " + JSON.stringify(e));
					//					},
					//					mouseout: function (e) {
					//						console.log("mouseout: " + JSON.stringify(e));
					//					},
					//					click: function (e) {
					//						console.log("click: " + JSON.stringify(e));
					//					},

					fill: false,
					lineTension: 0,
					pointRadius: 0,
					label: 'Key Proportion',
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
							color: '#666',
							display: true
						}
					}],
					yAxes: [{
						id: 'y-axis',
						gridLines: {
							color: '#666',
							display: true,
							drawBorder: false,
							drawOnChartArea: true
						},
						scaleLabel: {
							display: true,
							labelString: "Bounds"
						},
						ticks: {
							min: 0,
							max: 32,
							stepSize: 4,
							color: '#eee',
							display: true
						}
					}],
				},
				responsive: false,
				maintainAspectRatio: false,
				plugins: {
					// zoom plugin is only used by the env graphs
					zoom: {
						zoom: {
							enabled: false
						},
						pan: {
							enabled: false
						}
					}
				}
			}
		});

		if (viewVCE_voice.voicingMode) {
			viewVCE_chartdrag.init(viewVCE_keyprop.chart, viewVCE_keyprop.onchange, 'keyprop', 0, 23, 0, 32);
		}
		//console.log('--- finish viewVCE_keyprop init');
	}

};
