let viewVCE_filters = {
	chart: null,

	onchange: function (ele) {
		var id = ele.id;
		console.log("changed: " + id);

		var selectEle = document.getElementById("filterSelect");
		var filterIndex = selectEle.selectedIndex;
		var filterValue = vce.Head.FILTER[filterIndex];

		var param
		var args
		var filterPattern = /filter\[(\d+)\]/;
		if (ret = id.match(filterPattern)) {
			funcname = "setFILTEREle";
			filterValue = filterValue;
			index = parseInt(ret[1])
			value = parseInt(ele.value, 10);
		}
		//console.dir(vce);
		let message = {
			"name": "setFilterEle",
			"payload": {
				"FilterValue": filterValue,
				"Index": index,
				"Value": value
			}
		};
		astilectron.sendMessage(message, function (message) {
			console.log("setFilterEle returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the value
				index.errorNotification(message.payload);
				return false;
			} else {
				vce.Extra.uncompressedFilters[filterIndex][index - 1] = value;
				viewVCE_filters.filtersChartUpdate(selectEle.options[filterIndex].value, selectEle.options[filterIndex].innerHTML);
			}
		});
		return true;
	},

	uncompressFilters: function () {
		// the first time we evaluate this vce, the Filters may be compressed.  To make it easier to add/remove 
		// filters in the editor, we rewrite the Filters array such that the zeroth element is the A-filter 
		// (whether the voice uses an Afilter or not), and then each row from 1-16 for each oscillator's B-filter - 
		// again whether the osc uses one or not
		if (vce.Extra["uncompressedFilters"] != undefined) {
			// no need to do it again
			return;
		}

		var oldFilterIdx = 0;
		var newFilters = new Array(17);
		for (var i = 0; i <= vce.Head.VOITAB; i++) {
			if (vce.Head.FILTER[i] < 0) {
				// the zeroth entry in the compressed array is the a-filter
				console.log("copying zeroth filter as A-filter");
				newFilters[0] = vce.Filters[0];
				oldFilterIdx++;
				break;
			}
		}
		if (newFilters[0] == undefined) {
			newFilters[0] = new Array(32);
			for (i = 0; i < 32; i++) {
				newFilters[0][i] = 0;
			}
		}
		// now for b filters -
		for (var i = 0; i <= vce.Head.VOITAB; i++) {
			if (vce.Head.FILTER[i] > 0) {
				console.log("copying " + oldFilterIdx + " filter as B-filter to osc #" + (i + 1));
				newFilters[i + 1] = vce.Filters[oldFilterIdx];
				oldFilterIdx++;
			}
		}
		// prepopulate oscillators not yet in use "just in case":
		for (var osc = 1; osc <= 16; osc++) {
			if (newFilters[osc] == undefined) {
				console.log("dummy filter for osc #" + osc);
				newFilters[osc] = new Array(32);
				for (i = 0; i < 32; i++) {
					newFilters[osc][i] = 0;
				}
			}
		}
		vce.Extra.uncompressedFilters = newFilters;
		console.log("Uncompressed filters:");
		console.dir(vce);
	},

	init: function () {
		var selectEle = document.getElementById("filterSelect");
		// remove old options:
		while (selectEle.firstChild) {
			selectEle.removeChild(selectEle.firstChild);
		}


		viewVCE_filters.uncompressFilters();

		var filterNames = [];
		var filterValues = [];
		// check for a-filter
		for (i = 0; i <= vce.Head.VOITAB; i++) {
			if (vce.Head.FILTER[i] < 0) {
				// zero'th filter is the a-filter
				filterNames.push('Af');
				filterValues.push(0);
				break;
			}
		}
		// now the b-filters
		for (i = 0; i <= vce.Head.VOITAB; i++) {
			if (vce.Head.FILTER[i] > 0) {
				filterNames.push('Bf ' + vce.Head.FILTER[i]);
				filterValues.push(i+1);
			}
		}

	console.log("filter names: " + filterNames);

		// Option values are index into the vce.Head.FILTERS array (and one extra with value -1 for "All")
		if (filterNames.length > 1) {
			var option = document.createElement("option");
			option.value = -1;
			option.innerHTML = "All";
			selectEle.appendChild(option);
		}

		for (i = 0; i < filterNames.length; i++) {
			var option = document.createElement("option");
			option.value = filterValues[i];
			option.innerHTML = filterNames[i];
			selectEle.appendChild(option);
		}
		if (filterNames.length <= 0) {
			// no filters
			document.getElementById("filtersChart").style.display = "none";
		} else if (filterNames.length > 1) {
			// "All" == -1
			viewVCE_filters.filtersChartUpdate(-1, 'All');
		} else {
			// first filter
			viewVCE_filters.filtersChartUpdate(filterValues[0], filterNames[0]);
		}
	},

	filtersChartUpdate: function (filterIndex, filterName) {
		filterIndex = parseInt(filterIndex, 10);
		console.log("filtersChart init " + filterIndex);
		var datasets = [];

		if (filterIndex >= 0) {
			console.log("filter " + filterIndex + ": " + vce.Extra.uncompressedFilters[filterIndex]);
			$('#filterTable').show();
			$('#filterTable td.val input').each(function (i, obj) {
				var id = obj.id;
				// id is "filter[<n>]" - we need the <n> part
				var idxString = id.substring(7);
				var idx = parseInt(idxString, 10) - 1;
				obj.value = vce.Extra.uncompressedFilters[filterIndex][idx];
			});
			// match the color rotation below.  We don't allocate a color for an "unused" Bf. So this is
			// senselessly complicated. Go look at the FILTERS array and figure out which "compressed" index 
			// this was
			var color = chartColors[chartColors.length - 1];
			if (filterIndex > 0) {
				var colorIdx = 0
				for (var i = 0; i <= vce.Head.VOITAB; i++) {
					if (filterIndex == vce.Head.FILTER[i]) {
						color = chartColors[(colorIdx) % chartColors.length]
						break;
					} else if (vce.Head.FILTER[i] > 0) {
						colorIdx++;
					}
				}
			}
			datasets = [{
				fill: false,
				lineTension: 0,
				pointRadius: 0,
				pointHitRadius: 5,
				label: filterName,
				backgroundColor: color,
				borderColor: color,
				data: vce.Extra.uncompressedFilters[filterIndex]
			}];
		} else {
			// "all"
			$('#filterTable').hide();
			console.log("filter len : " + vce.Extra.uncompressedFilters.length);
			// only include the ones actually in use (see vce.Head.FILTER entry)
			for (i = 0; i <= vce.Head.VOITAB; i++) {
				// A-table if in use goes first:
				if (vce.Head.FILTER[i] < 0) {
					var filterName = "Af";
					// we use modulo to compute b-filter color - since A filter is "-1", use the last color in the table
					var color = chartColors[chartColors.length - 1]
					datasets.push(
						{
							fill: false,
							lineTension: 0,
							pointRadius: 0,
							pointHitRadius: 5,
							label: filterName,
							backgroundColor: color,
							borderColor: color,
							data: vce.Extra.uncompressedFilters[0] // ZEROTH ele is the a-filter
						}
					);
				}
			}
			var colorIdx = 0;
			for (i = 0; i <= vce.Head.VOITAB; i++) {
				// A-table if in use goes first:
				if (vce.Head.FILTER[i] > 0) {
					var filterName = "Bf " + (i + 1);
					var color = chartColors[colorIdx % chartColors.length];
					colorIdx++;
					datasets.push(
						{
							fill: false,
							lineTension: 0,
							pointRadius: 0,
							pointHitRadius: 5,
							label: filterName,
							backgroundColor: color,
							borderColor: color,
							data: vce.Extra.uncompressedFilters[i + 1] // b-filters by one-based osc value
						}
					);
				}
			}
			// now B-filters:

			console.log("datasets = " + JSON.stringify(datasets));
		}

		console.dir(datasets);

		var ctx = document.getElementById('filtersChart').getContext('2d');
		if (viewVCE_filters.chart != null) {
			viewVCE_filters.chart.destroy();
		}
		viewVCE_filters.chart = new Chart(ctx, {

			type: 'line',
			data: {
				labels: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''],
				datasets: datasets
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
						gridLines: {
							color: '#666',
							display: true,
							drawBorder: false,
							drawOnChartArea: false
						},
						scaleLabel: {
							display: true,
							labelString: "Frequency"
						},
						ticks: {
							color: '#666',
							display: true
						}
					}],
					yAxes: [{
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
