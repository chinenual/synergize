let vce = {};
let vceFiltersChart = null;

let viewVCE = {
    init: function() {
	Chart.defaults.global.defaultFontColor = 'white';
	Chart.defaults.global.defaultFontSize  = 14;
    },
    
    keyPropCurve: function(kprop) {
	var result = [];
	// y = 0..32
	// x = 0..23
	for (v = 0; v < kprop.length; v++) {
	    result[v] = kprop[v];
	}
	return result;
    },
    
    keyEqCurve: function(keq) {
	var result = [];
	// y = -24..6
	// x = 0..23
	for (v = 0; v < keq.length; v++) {
	    result[v] = keq[v];
	}
	return result;
    },
    
    timbreProportionCurve: function(center, sensitivity) {
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
	    var p = (center * 2) - 15 + (v/10 * sensitivity) - (2* sensitivity);
	    if (p > 31) p = 31;
	    if (p < 0)  p = 0;
	    result[v] = p;
	}
	return result;
    },
    
    ampProportionCurve: function(center, sensitivity) {
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
	    var p = (((v/10 * sensitivity)/2 - sensitivity) / 2) + center - 24
	    if (p > 6) p = 6;
	    if (p < -24)  p = -24;
	    result[v] = p+25;
	}
	return result;
    },
    
    patchTable: function() {
	var tbody = document.getElementById("patchTbody");
	// remove old rows:
	while (tbody.firstChild) {
	    tbody.removeChild(tbody.firstChild);
	}
	console.log("vceEnv" + vce.Envelopes);
	console.log("vceFilters" + vce.Filters);
	
	// populate new ones:
	for (osc = 0; osc <= vce.Head.VOITAB; osc++) {
	    var tr = document.createElement("tr");
	    var td = document.createElement("td");

	    //--- OSC
	    td.innerHTML = osc+1; // Osc
	    tr.appendChild(td);
	    
	    // FIXME: assumes envelopes are sorted in oscillator order
	    var patchByte = vce.Envelopes[osc].FreqEnvelope.OPTCH;
	    var patchInhibitAddr = (patchByte & 0x20) != 0;
	    var patchInhibitF0   = (patchByte & 0x04) != 0;
	    var patchOutputDSR   = ((patchByte & 0xc0) >> 6);
	    var patchAdderInDSR  = ((patchByte & 0x18) >> 3);
	    var patchFOInputDSR  = (patchByte & 0x03);

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
	    td.innerHTML = vce.Envelopes[osc].FreqEnvelope.OHARM
	    tr.appendChild(td);
	    
	    //--- Detn
	    td = document.createElement("td");
	    td.innerHTML = vce.Envelopes[osc].FreqEnvelope.FDETUN
	    tr.appendChild(td);

	    var waveByte = vce.Envelopes[osc].FreqEnvelope.Table[0][3];
	    var wave = ((waveByte & 0x7) == 0) ? 'Sine' : 'Triangle';
	    var keyprop = ((waveByte & 0x8) == 0) ? 'Key' : '';
	    
	    //--- Wave
	    td = document.createElement("td");
	    td.innerHTML = wave;
	    tr.appendChild(td);
	    
	    //--- Key
	    td = document.createElement("td");
	    td.innerHTML = keyprop;	    
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

    keyEqChart: function() {
	console.log("keyEqChart init");
	var propData = viewVCE.keyEqCurve(vce.Head.VEQ);
	
	var ctx = document.getElementById('keyEqChart').getContext('2d');
	var chart = new Chart(ctx, {
	    
	    type: 'line',
	    data: {
		labels: ['','','','','','','','','','','','','','','','','','','','','','','',''],
		//labels: ['','','',''],
		datasets: [{
		    fill: false,
		    lineTension: 0,
		    pointRadius: 0,
		    label: 'Key Equalization',
		    // https://www.color-hex.com/color-palette/89750
		    backgroundColor: 'rgb(255,215,0)',
		    borderColor: 'rgb(255,215,0)',
		    data: propData
		}]
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
	    }
	});
    },

    keyPropChart: function() {
	console.log("keyPropChart init");
	var propData = viewVCE.keyPropCurve(vce.Head.KPROP);
	
	var ctx = document.getElementById('keyPropChart').getContext('2d');
	var chart = new Chart(ctx, {
	    
	    type: 'line',
	    data: {
		labels: ['','','','','','','','','','','','','','','','','','','','','','','',''],
		datasets: [{
		    fill: false,
		    lineTension: 0,
		    pointRadius: 0,
		    label: 'Key Proportion',
		    // https://www.color-hex.com/color-palette/89750
		    backgroundColor: 'rgb(255,215,0)',
		    borderColor: 'rgb(255,215,0)',
		    data: propData
		}]
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
		maintainAspectRatio: false
	    }
	});
    },

    filtersChartInit: function() {
	var selectEle = document.getElementById("filterSelect");
	// remove old options:
	while (selectEle.firstChild) {
	    selectEle.removeChild(selectEle.firstChild);
	}
	var filter0name = '';
	for (i = 0; i < vce.Head.FILTER.length; i++) {
	    if (vce.Head.FILTER[i] != 0) {
		if (vce.Head.FILTER[i] > 0) {
		    // a B-filter
		    var option= document.createElement("option");
		    option.value=vce.Head.FILTER[i]-1;
		    option.innerHTML = "Bf " + vce.Head.FILTER[i];		    
		    if (i==0) filter0name = "Bf " + vce.Head.FILTER[i];
		    selectEle.appendChild(option);
		} else {
		    // an A-filter
		    var option= document.createElement("option");
		    option.value=(-vce.Head.FILTER[i])-1;
		    option.innerHTML = "Af " + (-vce.Head.FILTER[i]);
		    if (i==0) filter0name = "Af " + (-vce.Head.FILTER[i]);
		    selectEle.appendChild(option);
		}
	    }
	}
	if (selectEle.firstChild == null) {
	    // no filters
	    document.getElementById("filtersChart").style.display="none";
	} else {
	    viewVCE.filtersChartUpdate(0,filter0name);
	}
    },

    filtersChartUpdate: function(filterIndex, filterName) {
	console.log("filtersChart init " + filterIndex);
	var filterData = vce.Filters[filterIndex];

	$('#filterTable td.val').each(function(i, obj) {
	    var id=obj.id;
	    // id is "ft<n>" - we need the <n> part
	    var idxString=id.substring(2);
	    var idx=parseInt(idxString,10)-1;
	    obj.innerHTML = filterData[idx];
	});
				
	var ctx = document.getElementById('filtersChart').getContext('2d');
	if (vceFiltersChart != null) {
	    vceFiltersChart.destroy();
	}
	vceFiltersChart = new Chart(ctx, {
	    
	    type: 'line',
	    data: {
		labels: ['','','','','','','','','','', '','','','','','','','','','', '','','','','','','','','','','',''],
		datasets: [{
		    fill: false,
		    lineTension: 0,
		    pointRadius: 0,
		    label: filterName,
		    // https://www.color-hex.com/color-palette/89750
		    backgroundColor: 'rgb(255,215,0)',
		    borderColor: 'rgb(255,215,0)',
		    data: filterData
		}]
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
    },

    refreshText: function () {
	console.log("vceVoiceTab init");
	
	viewVCE.keyPropChart();
	viewVCE.keyEqChart();
	viewVCE.filtersChartInit();
	viewVCE.patchTable();
		
	if (crt_name == null) {
	    document.getElementById("backToCRT").hidden = true;
	} else {
	    document.getElementById("backToCRT").hidden = false;
	}

	document.getElementById("vce_name").innerHTML = vce.Head.VNAME;
	
	document.getElementById("name").innerHTML = vce.Head.VNAME;
	document.getElementById("nOsc").innerHTML = vce.Head.VOITAB + 1;
	document.getElementById("keysPlayable").innerHTML = Math.floor(32 / (vce.Head.VOITAB + 1));
	document.getElementById("vibType").innerHTML = (vce.Head.VIBDEL >= 0) ? "Sine" : "Random";
	document.getElementById("vibRate").innerHTML = vce.Head.VIBRAT;
	document.getElementById("vibDelay").innerHTML = vce.Head.VIBDEL;
	document.getElementById("vibDepth").innerHTML = vce.Head.VIBDEP;
	document.getElementById("transpose").innerHTML = vce.Head.VTRANS;
	var i;
	var count = 0;
	for (i = 0; i < vce.Head.FILTER.length; i++) {
	    if (vce.Head.FILTER[i] != 0) {
		count++;
	    }
	}
	document.getElementById("nFilter").innerHTML = count;
	
	Chart.defaults.global.defaultFontColor = 'white';
	Chart.defaults.global.defaultFontSize  = 14;
	
	var ampData = viewVCE.ampProportionCurve(vce.Head.VACENT, vce.Head.VASENS);
	var timbreData = viewVCE.timbreProportionCurve(vce.Head.VTCENT, vce.Head.VTSENS);
	
	console.log("ampData: " + ampData);
	console.log("timbreData: " + timbreData);
	
	var ctx = document.getElementById('velocityChart').getContext('2d');
	var chart = new Chart(ctx, {
	    
	    type: 'line',
	    data: {
		labels: ['','','','','','','','','','','','','','','','','','','','','','','','','','','','','','','',''],
		//labels: ['','','',''],
		datasets: [{
		    fill: false,
		    lineTension: 0,
		    pointRadius: 0,
		    label: 'Amplitude',
		    // https://www.color-hex.com/color-palette/89750
		    backgroundColor: 'rgb(255,215,0)',
		    borderColor: 'rgb(255,215,0)',
		    data: ampData
		},{
		    fill: false,
		    lineTension: 0,
		    pointRadius: 0,
		    label: 'Timbre',
		    backgroundColor: 'rgb(79,156,244)',
		    borderColor: 'rgb(79,156,244)',
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
			    callback: function(dataLabel, index) { return ''; }
			}
		    }],
		},
		responsive: false,
		maintainAspectRatio: false
	    }
	});
    }
};
