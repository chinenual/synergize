let viewVCE_keyprop = {

    keyPropCurve: function(kprop) {
	var result = [];
	// y = 0..32
	// x = 0..23
	for (v = 0; v < kprop.length; v++) {
	    result[v] = kprop[v];
	}
	return result;
    },
    
    init: function() {
	console.log("keyPropChart init");
	var propData = viewVCE_keyprop.keyPropCurve(vce.Head.KPROP);
	
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
    }
    
};