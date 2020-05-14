let vce = {};
let vceFiltersChart = null;

// https://www.color-hex.com/color-palette/89750
let chartColors=[
    'rgb(225,215,0)',
    'rgb(79,156,244)',
    'rgb(62,244,6)',
    'rgb(5,82,244)',
    'rgb(4,169,24)'
];

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

    envChartInit: function() {
	var selectEle = document.getElementById("envOscSelect");
	// remove old options:
	while (selectEle.firstChild) {
	    selectEle.removeChild(selectEle.firstChild);
	}

	for (i = 0; i <= vce.Head.VOITAB; i++) {
	    var option= document.createElement("option");
	    option.value="" + (i+1);
	    option.innerHTML = "" + (i+1);
	    selectEle.appendChild(option);
	}

	viewVCE.envChartUpdate(1)
    },

    // SYNHCS COMMON.Z80 FTAB:
    ftab: [
	0,2,4,6,8,10,12,14,	//Reference frequency table
	15,16,17,18,19,20,21,22,
	24,25,27,28,30,32,34,36,
	38,40,43,45,48,51,54,57,
	61,64,68,72,76,81,86,91,
	96,102,108,115,122,129,137,145,
	153,163,172,183,193,205,217,230,
	244,258,274,290,307,326,345,366,
	387,411,435,461,488,517,548,581,
	615,652,691,732,775,822,870,922,
	977,1035,1097,1162,1231,1304,1382,1464,
	1551,1644,1741,1845,1955,2071,2194,2325,
	2463,2609,2765,2929,3103,3288,3483,3691,
	3910,4143,4389,4650,4926,5219,5530,5859,
	6207,6576,6967,7382,7820,8286,8778,9300,
	9853,10439,11060,11718,12414,13153,13935,14764],

    scaleViaRtab(v) {
	return v;// TODO
    },
    
    scaleFreqEnvValue: function(v) {
	// See OSCDSP.Z80 DISVAL: DVAL10:
	return v; // TODO
    },
    
    scaleAmpEnvValue: function(v,last) {
	// See OSCDSP.Z80 DISVAL: DVAL30:
	console.log("scaleAmpEnvValue " + v + " " + last);	
	if (last) return 0;
	return Math.max(0, v-55);
    },
    
    scaleFreqTimeValue: function(v) {
	// See OSCDSP.Z80 DISVAL: 
	if (v < 15) return v;
	return viewVCE.scaleViaRtab((2*v)-14);
    },
    
    scaleAmpTimeValue: function(v) {
	// See OSCDSP.Z80 DISVAL: DVAL20:
	if (v < 23) return v;
	return viewVCE.scaleViaRtab(((v*2)-54)*2);
    },
    
    envChartUpdate: function(oscNum) {
	var oscIndex = oscNum-1;
	var envelopes = vce.Envelopes[oscIndex];
	
	// clear old values:

	$('#envTable td.val').each(function(i, obj) {
	    obj.innerHTML = '';
	});
	// fill in freq env data:

	// scaling algorithms derived from DISVAL: in OSCDSP.Z80
	
	console.log("env freq npt " + envelopes.FreqEnvelope.NPOINTS);
	for (i = 0; i < envelopes.FreqEnvelope.NPOINTS; i++) {
	    var tr = $('#envTable tbody tr:eq(' + i + ')');

	    // table is logically in groups of 4
	    tr.find('td:eq(2)')
		.html(viewVCE.scaleFreqEnvValue(envelopes.FreqEnvelope.Table[i*4 + 0])); 
	    tr.find('td:eq(3)')
		.html(viewVCE.scaleFreqEnvValue(envelopes.FreqEnvelope.Table[i*4 + 1])); 
	    tr.find('td:eq(4)')
		.html(viewVCE.scaleFreqTimeValue(envelopes.FreqEnvelope.Table[i*4 + 2])); 
	    tr.find('td:eq(5)')
		.html(viewVCE.scaleFreqTimeValue(envelopes.FreqEnvelope.Table[i*4 + 3])); 
	    if (envelopes.AmpEnvelope.SUSTAINPT == (i+1)) {
		tr.find('td:eq(1)').html("S&gt;");
	    }
	    if (envelopes.AmpEnvelope.LOOPNPT == (i+1)) {
		tr.find('td:eq(1)').html("L&gt;");
	    }
	}

	console.log("env amp npt " + envelopes.AmpEnvelope.NPOINTS);
	for (i = 0; i < envelopes.AmpEnvelope.NPOINTS; i++) {
	    var tr = $('#envTable tbody tr:eq(' + i + ')');

	    // table is logically in groups of 4.
	    // "j" accounts for the difference in column index due to the
	    // row-spanning separators (only in i==0):
	    j = (i == 0) ? 11 : 9;
	    console.log("j " + j);
	    console.dir(tr);
	    console.dir(tr.find('td:eq(' +(j+0)+ ')'));
	    tr.find('td:eq(' +(j+0)+ ')')
		.html(viewVCE.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i*4 + 0],
					       (i+1) >= envelopes.AmpEnvelope.NPOINTS)); 
	    tr.find('td:eq(' +(j+1)+ ')')
		.html(viewVCE.scaleAmpEnvValue(envelopes.AmpEnvelope.Table[i*4 + 1],
					       (i+1) >= envelopes.AmpEnvelope.NPOINTS)); 
	    tr.find('td:eq(' +(j+2)+ ')')
		.html(viewVCE.scaleAmpTimeValue(envelopes.AmpEnvelope.Table[i*4 + 2])); 
	    tr.find('td:eq(' +(j+3)+ ')')
		.html(viewVCE.scaleAmpTimeValue(envelopes.AmpEnvelope.Table[i*4 + 3]));
	    if (envelopes.AmpEnvelope.SUSTAINPT == (i+1)) {
		tr.find('td:eq(' +(j-1)+ ')').html("S&gt;");
	    }
	    if (envelopes.AmpEnvelope.LOOPNPT == (i+1)) {
		tr.find('td:eq(' +(j-1)+ ')').html("L&gt;");
	    }
	}



    },
	
    filtersChartInit: function() {
	var selectEle = document.getElementById("filterSelect");
	// remove old options:
	while (selectEle.firstChild) {
	    selectEle.removeChild(selectEle.firstChild);
	}
	var filter0name = '';
	var idx=0;
	for (i = 0; i < vce.Head.FILTER.length; i++) {
	    if (vce.Head.FILTER[i] != 0) {
		if (idx==0 && vce.Head.FILTER[0] != 0) {
		    var option= document.createElement("option");
		    option.value=-1;
		    option.innerHTML = "All";
		    selectEle.appendChild(option);
		}
		if (vce.Head.FILTER[i] > 0) {
		    // a B-filter
		    var option= document.createElement("option");
		    option.value=vce.Head.FILTER[i]-1;
		    option.innerHTML = "Bf " + vce.Head.FILTER[i];		    
		    if (idx==0) filter0name = "Bf " + vce.Head.FILTER[i];
		    selectEle.appendChild(option);
		    idx++;
		} else {
		    // an A-filter
		    var option= document.createElement("option");
		    option.value=(-vce.Head.FILTER[i])-1;
		    option.innerHTML = "Af " + (-vce.Head.FILTER[i]);
		    if (idx==0) filter0name = "Af " + (-vce.Head.FILTER[i]);
		    selectEle.appendChild(option);
		    idx++;
		}
	    }
	}
	if (selectEle.firstChild == null) {
	    // no filters
	    document.getElementById("filtersChart").style.display="none";
	} else {
	    // "All" == -1
	    viewVCE.filtersChartUpdate(-1);
	}
    },

    filtersChartUpdate: function(filterIndex) {
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
	    var filterName = $('#filterSelect option').eq(filterIndex+1).html();
	    datasets = [{
		    fill: false,
		    lineTension: 0,
		    pointRadius: 0,
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
			label: filterName,
			backgroundColor: chartColors[i % chartColors.length],
			borderColor: chartColors[i % chartColors.length],
			data: vce.Filters[i]
		    }
		);
	    }
	    console.log("datasets = " + datasets);
	}
	
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
    },

    refreshText: function () {
	console.log("vceVoiceTab init");
	
	viewVCE.keyPropChart();
	viewVCE.keyEqChart();
	viewVCE.envChartInit();
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
