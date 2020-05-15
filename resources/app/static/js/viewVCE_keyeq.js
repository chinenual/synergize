let viewVCE_keyeq = {
    keyEqCurve: function(keq) {
	var result = [];
	// y = -24..6
	// x = 0..23
	for (v = 0; v < keq.length; v++) {
	    result[v] = keq[v];
	}
	return result;
    },
    
    init: function() {
	console.log("keyEqChart init");
	var propData = viewVCE_keyeq.keyEqCurve(vce.Head.VEQ);
	
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
		    backgroundColor: chartColors[0],
		    borderColor: chartColors[0],
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
    }    
};
