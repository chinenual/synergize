let viewVCE_filters = {
	chart: null,

	onchange: function (ele) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		viewVCE_filters.deb_onchange(ele);
	},

	deb_onchange: null, // initialized during init()

	raw_onchange: function (ele) {
		if (viewVCE.supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/ return; }

		var id = ele.id;

		var value = index.checkInputElementValue(ele);
		if (value === undefined) {
			return;
		}

		var eleIndex;
		var selectEle = document.getElementById("filterSelect");
		var filterIndex = parseInt(selectEle.value, 10); // index into the uncompressedFilters array
		var filterName = selectEle.options[selectEle.selectedIndex].innerHTML;
		var filterValue = vce.Head.FILTER[filterIndex];

		console.log("filter ele change " + filterIndex + " val: " + filterValue);

		var param
		var args
		var filterPattern = /flt\[(\d+)\]/;
		if (ret = id.match(filterPattern)) {
			funcname = "setFILTEREle";
			eleIndex = parseInt(ret[1])
		}
		//console.dir(vce);
		let message = {
			"name": "setFilterEle",
			"payload": {
				"UiFilterIndex": filterIndex,
				"Index": eleIndex,
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
				vce.Extra.uncompressedFilters[filterIndex][eleIndex - 1] = value;
				viewVCE_filters.filtersChartUpdate(filterIndex, filterName, false);
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
				newFilters[i + 1] = vce.Filters[oldFilterIdx];
				oldFilterIdx++;
			}
		}
		// prepopulate oscillators not yet in use "just in case":
		for (var osc = 1; osc <= 16; osc++) {
			if (newFilters[osc] == undefined) {
				newFilters[osc] = new Array(32);
				for (i = 0; i < 32; i++) {
					newFilters[osc][i] = 0;
				}
			}
		}
		vce.Extra.uncompressedFilters = newFilters;
	},

	filterNames: [],
	filterValues: [],

	init: function () {
		console.log('--- start viewVCE_filters init');
		if (viewVCE_filters.deb_onchange == null) {
			//viewVCE_filters.deb_onchange = _.debounce(viewVCE_filters.raw_onchange, 250);
			viewVCE_filters.deb_onchange = viewVCE_filters.raw_onchange;
		}
		if (viewVCE_filters.deb_copyFrom == null) {
			viewVCE_filters.deb_copyFrom = _.debounce(viewVCE_filters.raw_copyFrom, 250);
		}

		var selectEle = document.getElementById("filterSelect");
		// remove old options:
		while (selectEle.firstChild) {
			selectEle.removeChild(selectEle.firstChild);
		}


		viewVCE_filters.uncompressFilters();

		viewVCE_filters.filterNames = [];
		viewVCE_filters.filterValues = [];
		// check for a-filter
		for (i = 0; i <= vce.Head.VOITAB; i++) {
			if (vce.Head.FILTER[i] < 0) {
				// zero'th filter is the a-filter
				viewVCE_filters.filterNames.push('Af');
				viewVCE_filters.filterValues.push(0);
				break;
			}
		}
		// now the b-filters
		for (i = 0; i <= vce.Head.VOITAB; i++) {
			if (vce.Head.FILTER[i] > 0) {
				// FIXME: naming/numbering can be confusing.  For example, INTERNAL/CATHERG voice uses
				// one B-filter - for osc#3.  So we name it "Bf1",but it shows as index "3" in the select. 
				// Should we name it Bf3 to match the osc?
				viewVCE_filters.filterNames.push('Bf ' + vce.Head.FILTER[i]);
				viewVCE_filters.filterValues.push(i + 1);
			}
		}

		// Option values are index into the vce.Head.FILTERS array (and one extra with value -1 for "All")
		if (viewVCE_filters.filterNames.length > 1) {
			var option = document.createElement("option");
			option.value = -1;
			option.innerHTML = "All";
			selectEle.appendChild(option);
		}

		for (i = 0; i < viewVCE_filters.filterNames.length; i++) {
			var option = document.createElement("option");
			option.value = viewVCE_filters.filterValues[i];
			option.innerHTML = viewVCE_filters.filterNames[i];
			selectEle.appendChild(option);
		}
		document.getElementById("filtersChart").style.display = "block";
		document.getElementById("filterTable").style.display = "block";
		if (viewVCE_filters.filterNames.length <= 0) {
			// no filters
			document.getElementById("filtersChart").style.display = "none";
			document.getElementById("filterTable").style.display = "none";
			$('#filterCopySelectDiv').hide();
		} else if (viewVCE_filters.filterNames.length > 1) {
			// "All" == -1
			viewVCE_filters.filtersChartUpdate(-1, 'All', true);
		} else {
			// first filter
			viewVCE_filters.filtersChartUpdate(viewVCE_filters.filterValues[0], viewVCE_filters.filterNames[0], true);
		}
		console.log('--- finish viewVCE_filters init');
	},

	copyFrom: function (filterIndex, filterName) {
		if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/ return; }
		viewVCE_filters.deb_copyFrom(filterIndex, filterName);
	},

	deb_copyFrom: null, // initialized during init()

	raw_copyFrom: function (fromFilterIndex, fromFilterName) {
		fromFilterIndex = parseInt(fromFilterIndex, 10);

		if (fromFilterIndex < 0) {
			return;
		}
		var filterSelectEle = document.getElementById("filterSelect");
		var toFilterIndex = filterSelectEle.options[filterSelectEle.selectedIndex].value;
		toFilterIndex = parseInt(toFilterIndex, 10);

		let message = {
			"name": "setFilterArray",
			"payload": {
				"UiFilterIndex": toFilterIndex,
				"Values": vce.Extra.uncompressedFilters[fromFilterIndex]
			}
		};
		index.spinnerOn();
		astilectron.sendMessage(message, function (message) {
			index.spinnerOff();
			console.log("setFilterArray returned: " + JSON.stringify(message));
			// Check error
			if (message.name === "error") {
				// failed - dont change the value
				index.errorNotification(message.payload);
				return false;
			} else {
				for (i = 0; i < 32; i++) {
					vce.Extra.uncompressedFilters[toFilterIndex][i] = vce.Extra.uncompressedFilters[fromFilterIndex][i]
				}
				viewVCE_filters.filtersChartUpdate(filterSelectEle.options[filterSelectEle.selectedIndex].value,
					filterSelectEle.options[filterSelectEle.selectedIndex].innerHTML, true);
			}
		});
		return true;

	},

	filtersChartUpdate: function (filterIndex, filterName, animate) {
		filterIndex = parseInt(filterIndex, 10);
		var datasets = [];

		console.log("Filter update: index:" + filterIndex + " name:" + filterName);

		var filterCopySelectEle = document.getElementById("filterCopySelect");
		// remove old options:
		while (filterCopySelectEle.firstChild) {
			filterCopySelectEle.removeChild(filterCopySelectEle.firstChild);
		}
		// hide the copy selector for All or cases where there are no filters, or when we're not in voicing mode
		$('#filterCopySelectDiv').hide();

		if (filterIndex >= 0) {
			if (viewVCE_voice.voicingMode) {
				$('#filterCopySelectDiv').show();
				// populate options in the select with only "other" filters (i.e. "this" filter should be not shown or at least unselectable)

				// first element is empty to avoid confusing the user if they havent selected something:
				var option = document.createElement("option");
				option.value = -1;
				option.innerHTML = "";
				filterCopySelectEle.appendChild(option);

				for (i = 0; i < viewVCE_filters.filterNames.length; i++) {
					if (viewVCE_filters.filterValues[i] >= 0 && viewVCE_filters.filterValues[i] != filterIndex) {
						var option = document.createElement("option");
						option.value = viewVCE_filters.filterValues[i];
						option.innerHTML = viewVCE_filters.filterNames[i];
						filterCopySelectEle.appendChild(option);
					}
				}
			}

			$('#filterTable').show();
			$('#filterTable td.val input').each(function (i, obj) {
				var id = obj.id;
				// id is "flt[<n>]" - we need the <n> part
				var idxString = id.substring(4);
				var idx = parseInt(idxString, 10) - 1;
				obj.value = vce.Extra.uncompressedFilters[filterIndex][idx];

				viewVCE_voice.sendToMIDI(obj, id, vce.Extra.uncompressedFilters[filterIndex][idx]);
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

			// set MIDI control surface to all zeros:
			for (i = 1; i <= 32; i++) {
				var id = "flt[" + i + "]";
				var ele = document.getElementById(id);
				viewVCE_voice.sendToMIDI(ele, id, 0);
			}

			$('#filterTable').hide();
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

		}

		//		console.dir(datasets);

		var ctx = document.getElementById('filtersChart').getContext('2d');
		if (viewVCE_filters.chart != null) {
			viewVCE_filters.chart.destroy();
		}

		var animation_duration = animate ? 1000 : 0;

		viewVCE_filters.chart = new Chart(ctx, {

			type: 'line',
			data: {
				labels: ['', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''],
				datasets: datasets
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
