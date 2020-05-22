let viewVCE_filters = {
    
    init: function() {
	var selectEle = document.getElementById("filterSelect");
	// remove old options:
	while (selectEle.firstChild) {
	    selectEle.removeChild(selectEle.firstChild);
	}

	var filterNames = [];
	// check for a-filter
	for (i = 0; i < vce.Head.FILTER.length; i++) {
	    if (vce.Head.FILTER[i] < 0) {
		// zero'th filter is the a-filter
		filterNames.push('Af');
		break;
	    }
	}
	// now the b-filters
	for (i = 0; i < vce.Head.FILTER.length; i++) {
	    if (vce.Head.FILTER[i] > 0) {
		filterNames.push('Bf ' + vce.Head.FILTER[i]);
		break;
	    }
	}

	if (filterNames.length > 1) {
	    var option= document.createElement("option");
	    option.value=-1;
	    option.innerHTML = "All";
	    selectEle.appendChild(option);
	}
	
	for (i = 0; i < filterNames.length; i++) {
	    var option= document.createElement("option");
	    option.value=i;
	    option.innerHTML = filterNames[i];
	    selectEle.appendChild(option);
	}
	if (filterNames.length <= 0) {
	    // no filters
	    document.getElementById("filtersChart").style.display="none";
	} else if (filterNames.length > 1) {
	    // "All" == -1
	    viewVCE_filters.filtersChartUpdate(-1, 'All');
	} else {
	    // first filter
	    viewVCE_filters.filtersChartUpdate(0, filterNames[0]);
	}
    },

    filtersChartUpdate: function(filterIndex, filterName) {
	filterIndex = parseInt(filterIndex,10);
	console.log("filtersChart init " + filterIndex);
	var datasets = [];

	if (filterIndex >= 0) {
	    console.log("filter " + filterIndex + ": " + vce.Filters[filterIndex]);
	    $('#filterTable').show();
	    $('#filterTable td.val').each(function(i, obj) {
		var id=obj.id;
		// id is "ft<n>" - we need the <n> part
		var idxString=id.substring(2);
		var idx=parseInt(idxString,10)-1;
		obj.innerHTML = vce.Filters[filterIndex][idx];
	    });
	    datasets = [{
		fill: false,
		lineTension: 0,
		pointRadius: 0,
		pointHitRadius: 5,
		label: filterName,
		backgroundColor: chartColors[filterIndex % chartColors.length],
		borderColor: chartColors[filterIndex % chartColors.length],
		data: vce.Filters[filterIndex]
	    }];
	} else {
	    // "all"
	    $('#filterTable').hide();
	    console.log("filter len : " + vce.Filters.length);
	    for (i = 0; i < vce.Filters.length; i++) {
		console.log("filter " + i + ": " + vce.Filters[i]);
		var filterName = $('#filterSelect option').eq(i+1).html();
		datasets.push(
		    {
			fill: false,
			lineTension: 0,
			pointRadius: 0,
			pointHitRadius: 5,
			label: filterName,
			backgroundColor: chartColors[i % chartColors.length],
			borderColor: chartColors[i % chartColors.length],
			data: vce.Filters[i]
		    }
		);
	    }
	    console.log("datasets = " + JSON.stringify(datasets));
	}

	console.dir(datasets);
	
	var ctx = document.getElementById('filtersChart').getContext('2d');
	if (vceFiltersChart != null) {
	    vceFiltersChart.destroy();
	}
	vceFiltersChart = new Chart(ctx, {
	    
	    type: 'line',
	    data: {
		labels: ['','','','','','','','','','', '','','','','','','','','','', '','','','','','','','','','','',''],
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
