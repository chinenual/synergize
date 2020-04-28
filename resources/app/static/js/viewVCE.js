let vceHead = {};

let viewVCE = {
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
    refreshText: function () {
	if (crt_name == null) {
	    document.getElementById("backToCRT").hidden = true;
	} else {
	    document.getElementById("backToCRT").hidden = false;
	}
	
	document.getElementById("name").innerHTML = vceHead.VNAME;
	document.getElementById("nOsc").innerHTML = vceHead.VOITAB + 1;
	document.getElementById("keysPlayable").innerHTML = Math.floor(32 / (vceHead.VOITAB + 1));
	document.getElementById("vibType").innerHTML = (vceHead.VIBDEL >= 0) ? "Sine" : "Random";
	document.getElementById("vibRate").innerHTML = vceHead.VIBRAT;
	document.getElementById("vibDelay").innerHTML = vceHead.VIBDEL;
	document.getElementById("vibDepth").innerHTML = vceHead.VIBDEP;
	document.getElementById("transpose").innerHTML = vceHead.VTRANS;
	var i;
	var count = 0;
	for (i = 0; i < vceHead.FILTER.length; i++) {
	    if (vceHead.FILTER[i] != 0) {
		count++;
	    }
	}
	document.getElementById("nFilter").innerHTML = count;
	
	Chart.defaults.global.defaultFontColor = 'white';
	Chart.defaults.global.defaultFontSize  = 14;
	
	var ampData = viewVCE.ampProportionCurve(vceHead.VACENT, vceHead.VASENS);
	var timbreData = viewVCE.timbreProportionCurve(vceHead.VTCENT, vceHead.VTSENS);
	
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
