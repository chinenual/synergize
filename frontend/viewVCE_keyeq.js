import {Chart} from 'chart.js';

import { index } from './index';
import { viewVCE } from './viewVCE';
import { viewVCE_chartdrag } from './viewVCE_chartdrag';
import {viewVCE_voice} from './viewVCE_voice';
import {UIService} from '/bindings/github.com/chinenual/synergize';

import $ from 'jquery';
Object.assign(window, { $: $, jQuery: $ });

export let viewVCE_keyeq = {
  chart: null,

  keyEqCurve: function(keq) {
    let result = [];
    // y = -24..6
    // x = 0..23
    for (let v = 0; v < keq.length; v++) {
      result[v] = keq[v];
    }
    return result;
  },

  onchange: function(ele, updateChart) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    viewVCE_keyeq.deb_onchange(ele, updateChart);
  },

  deb_onchange: null,  // initialized during init()

  raw_onchange: async function(ele, updateChart) {
    if (viewVCE
            .supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/
      return;
    }
    let value = index.checkInputElementValue(ele);
    if (value == undefined) {
      return;
    }

    let id = ele.id;
    console.log('changed: ' + id + ' val: ' + ele.value);

    let eleIndex;
    let pattern = /keyeq\[(\d+)\]/;
    let ret = id.match(pattern);
    if (ret) {
      eleIndex = parseInt(ret[1])
    }
    try {
      await UIService.SetVoiceVEQEle(eleIndex, value);
    } catch (exc) {
      // failed - dont change the value
      index.errorNotification(exc);
      return false;
    }
    viewVCE.vce.Head.VEQ[eleIndex - 1] = value;
    if (updateChart) {
      viewVCE_keyeq.init(true);
    }
    viewVCE_voice.sendToCSurface(ele, id, value);
    return true;
  },

  init: function(incrementalUpdate) {
    // console.log('--- start viewVCE_keyeq init ' + incrementalUpdate);
    if (viewVCE_keyeq.deb_onchange == null) {
      // viewVCE_keyeq.deb_onchange = _.debounce(viewVCE_keyeq.raw_onchange,
      // DEBOUNCE_WAIT_SHORT);
      viewVCE_keyeq.deb_onchange = viewVCE_keyeq.raw_onchange;
    }

    let propData = viewVCE_keyeq.keyEqCurve(viewVCE.vce.Head.VEQ);

    $('#keyEqTable td.val input').each(function(i, obj) {
      let id = obj.id;
      // id is "keyeq[<n>]" - we need the <n> part
      let idxString = id.substring(6);
      let idx = parseInt(idxString, 10) - 1;

      obj.value = propData[idx];

      if (!incrementalUpdate) {
        viewVCE_voice.sendToCSurface(obj, id, propData[idx]);
      }
    });
    if (viewVCE_keyeq.chart != null) {
      viewVCE_keyeq.chart.destroy();
    }

    let ctx = document.getElementById('keyEqChart').getContext('2d');
    viewVCE_keyeq.chart = new Chart(ctx, {

      type: 'line',
      data: {
        labels: [
          '', '', '', '', '', '', '', '', '', '', '', '',
          '', '', '', '', '', '', '', '', '', '', '', ''
        ],
        // labels: ['','','',''],
        datasets: [{
          fill: false,
          lineTension: 0,
          pointRadius: 0,
          label: 'Key Equalization',
          backgroundColor: viewVCE.chartColors[0],
          borderColor: viewVCE.chartColors[0],
          data: propData
        }]
      },

      // Configuration options go here
      options: {
        animation: {duration: 0},

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
            scaleLabel: {display: true, labelString: 'Key'},
            ticks: {min: 0, max: 23, color: '#eee', display: true}
          }],
          yAxes: [{
            id: 'y-axis',
            grid: {color: '#666'},
            gridLines: {
              color: '#666',
              display: true,
              drawBorder: false,
              drawOnChartArea: true
            },
            scaleLabel: {display: true, labelString: 'dB'},
            ticks: {min: -24, max: 8, stepSize: 4, color: '#eee', display: true}
          }],
        },
        responsive: false,
        maintainAspectRatio: false,
        plugins: {
          // zoom plugin is only used by the env graphs
          zoom: {zoom: {enabled: false}, pan: {enabled: false}}
        }
      },

    });

    if (viewVCE_voice.voicingMode) {
      viewVCE_chartdrag.init(
          viewVCE_keyeq.chart, viewVCE_keyeq.onchange, 'keyeq', 0, 23, -24, 8);
    }

    // console.log('--- finish viewVCE_keyeq init');
  }
};
