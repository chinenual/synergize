import {UIService} from '/bindings/github.com/chinenual/synergize';
import {Chart} from 'chart.js';
import _ from 'lodash';

import {index} from './index';
import { viewVCE } from './viewVCE';
import { viewVCE_chartdrag } from './viewVCE_chartdrag';
import {viewVCE_voice} from './viewVCE_voice';

import $ from 'jquery';
Object.assign(window, { $: $, jQuery: $ });


export let viewVCE_filters = {
  chart: null,

  onchangeGain: function(ele) {
    let gain;
    if (ele.id.match(/Plus/)) {
      gain = 0.10;  // 10% up
    } else {
      gain = -0.10;  // 10% down
    }
    console.log(' gain ' + gain);

    for (let i = 1; i <= 32; i++) {
      let input = document.getElementById(`flt[${i}]`)
      if (input) {
        let oldval = parseInt(input.value, 10)
        let delta = oldval * gain;
        let newval = Math.round(oldval + delta);
        console.log(
            ' gain ' + gain + ' ' + input.id + ' ' + oldval + ' ' + delta +
            ' ' + newval);
        input.value = '' + newval;
        // update the chart after the last element
        viewVCE_filters.onchange(input, i === 32);
      }
    }
  },

  onchange: function(ele, updateChart) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    viewVCE_filters.deb_onchange(ele, updateChart);
  },

  deb_onchange: null,  // initialized during init()

  raw_onchange: async function(ele, updateChart) {
    if (viewVCE
            .supressOnchange) { /*console.log("raw viewVCE.suppressOnChange");*/
      return;
    }

    let id = ele.id;

    let value = index.checkInputElementValue(ele);
    if (value === undefined) {
      return;
    }

    let eleIndex;
    let selectEle = document.getElementById('filterSelect');
    let filterIndex = parseInt(
        selectEle.value, 10);  // index into the uncompressedFilters array
    let filterName = selectEle.options[selectEle.selectedIndex].innerHTML;
    let filterValue = viewVCE.vce.Head.FILTER[filterIndex];

    console.log('filter ele change ' + filterIndex + ' val: ' + filterValue);

    let filterPattern = /flt\[(\d+)\]/;
    let ret = id.match(filterPattern);
    if (ret) {
      eleIndex = parseInt(ret[1])
    }
    try {
      await UIService.SetFilterEle(filterIndex, eleIndex, value);
    } catch (exc) {
      index.errorNotification(exc);
      return false;
    }
    viewVCE.vce.Extra.uncompressedFilters[filterIndex][eleIndex - 1] = value;
    if (updateChart) {
      viewVCE_filters.filtersChartUpdate(filterIndex, filterName, false);
    }
    viewVCE_voice.sendToCSurface(ele, id, value);

    return true;
  },

  uncompressFilters: function() {
    // the first time we evaluate this viewVCE.vce, the Filters may be
    // compressed.  To make it easier to add/remove filters in the editor, we
    // rewrite the Filters array such that the zeroth element is the A-filter
    // (whether the voice uses an Afilter or not), and then each row from 1-16
    // for each oscillator's B-filter - again whether the osc uses one or not
    if (viewVCE.vce.Extra['uncompressedFilters'] != undefined) {
      // no need to do it again
      return;
    }

    let oldFilterIdx = 0;
    let newFilters = new Array(17);
    for (let i = 0; i <= viewVCE.vce.Head.VOITAB; i++) {
      if (viewVCE.vce.Head.FILTER[i] < 0) {
        // the zeroth entry in the compressed array is the a-filter
        newFilters[0] = viewVCE.vce.Filters[0];
        oldFilterIdx++;
        break;
      }
    }
    if (newFilters[0] == undefined) {
      newFilters[0] = new Array(32);
      for (let i = 0; i < 32; i++) {
        newFilters[0][i] = 0;
      }
    }
    // now for b filters -
    for (let i = 0; i <= viewVCE.vce.Head.VOITAB; i++) {
      if (viewVCE.vce.Head.FILTER[i] > 0) {
        newFilters[i + 1] = viewVCE.vce.Filters[oldFilterIdx];
        oldFilterIdx++;
      }
    }
    // prepopulate oscillators not yet in use "just in case":
    for (let osc = 1; osc <= 16; osc++) {
      if (newFilters[osc] == undefined) {
        newFilters[osc] = new Array(32);
        for (let i = 0; i < 32; i++) {
          newFilters[osc][i] = 0;
        }
      }
    }
    viewVCE.vce.Extra.uncompressedFilters = newFilters;
  },

  filterNames: [],
  filterValues: [],

  init: function(incrementalUpdate) {
    // console.log('--- start viewVCE_filters init ' + incrementalUpdate);
    if (viewVCE_filters.deb_onchange == null) {
      // viewVCE_filters.deb_onchange = _.debounce(viewVCE_filters.raw_onchange,
      // DEBOUNCE_WAIT_SHORT);
      viewVCE_filters.deb_onchange = viewVCE_filters.raw_onchange;
    }
    if (viewVCE_filters.deb_copyFrom == null) {
      viewVCE_filters.deb_copyFrom =
          _.debounce(viewVCE_filters.raw_copyFrom, index.DEBOUNCE_WAIT);
    }

    let selectEle = document.getElementById('filterSelect');
    // remove old options:
    while (selectEle.firstChild) {
      selectEle.removeChild(selectEle.firstChild);
    }


    viewVCE_filters.uncompressFilters();

    viewVCE_filters.filterNames = [];
    viewVCE_filters.filterValues = [];
    // check for a-filter
    for (let i = 0; i <= viewVCE.vce.Head.VOITAB; i++) {
      if (viewVCE.vce.Head.FILTER[i] < 0) {
        // zero'th filter is the a-filter
        viewVCE_filters.filterNames.push('Af');
        viewVCE_filters.filterValues.push(0);
        break;
      }
    }
    // now the b-filters
    for (let i = 0; i <= viewVCE.vce.Head.VOITAB; i++) {
      if (viewVCE.vce.Head.FILTER[i] > 0) {
        // FIXME: naming/numbering can be confusing.  For example,
        // INTERNAL/CATHERG voice uses one B-filter - for osc#3.  So we name it
        // "Bf1",but it shows as index "3" in the select. Should we name it Bf3
        // to match the osc?
        viewVCE_filters.filterNames.push('Bf ' + viewVCE.vce.Head.FILTER[i]);
        viewVCE_filters.filterValues.push(i + 1);
      }
    }

    // Option values are index into the viewVCE.vce.Head.FILTERS array (and one
    // extra with value -1 for "All")
    if (viewVCE_filters.filterNames.length > 1) {
      let option = document.createElement('option');
      option.value = -1;
      option.innerHTML = 'All';
      selectEle.appendChild(option);
    }

    for (let i = 0; i < viewVCE_filters.filterNames.length; i++) {
      let option = document.createElement('option');
      option.value = viewVCE_filters.filterValues[i];
      option.innerHTML = viewVCE_filters.filterNames[i];
      selectEle.appendChild(option);
    }
    document.getElementById('filtersChart').style.display = 'block';
    document.getElementById('filterTable').style.display = 'block';
    if (viewVCE_filters.filterNames.length <= 0) {
      // no filters
      document.getElementById('filtersChart').style.display = 'none';
      document.getElementById('filterTable').style.display = 'none';
      document.querySelector('#filterCopySelectDiv').style.display = 'none';
      for (let i = 1; i <= 32; i++) {
        viewVCE_voice.sendToCSurface(null, `flt[${i}]`, 0);
      }
    } else if (viewVCE_filters.filterNames.length > 1) {
      // "All" == -1
      viewVCE_filters.filtersChartUpdate(-1, 'All', true);
    } else {
      // first filter
      viewVCE_filters.filtersChartUpdate(
          viewVCE_filters.filterValues[0], viewVCE_filters.filterNames[0],
          true);
    }
    // console.log('--- finish viewVCE_filters init');
  },

  copyFrom: function(filterIndex, filterName) {
    if (viewVCE.supressOnchange) { /*console.log("viewVCE.suppressOnChange");*/
      return;
    }
    viewVCE_filters.deb_copyFrom(filterIndex, filterName);
  },

  deb_copyFrom: null,  // initialized during init()

  raw_copyFrom: async function(fromFilterIndex, fromFilterName) {
    fromFilterIndex = parseInt(fromFilterIndex, 10);

    if (fromFilterIndex < 0) {
      return;
    }
    let filterSelectEle = document.getElementById('filterSelect');
    let toFilterIndex =
        filterSelectEle.options[filterSelectEle.selectedIndex].value;
    toFilterIndex = parseInt(toFilterIndex, 10);

    index.spinnerOn();
    try {
      await UIService.SetFilterArray(
          toFilterIndex,
          viewVCE.vce.Extra.uncompressedFilters[fromFilterIndex]);
    } catch (exc) {
      index.errorNotification(exc);
      return false;
    }

    for (let i = 0; i < 32; i++) {
      viewVCE.vce.Extra.uncompressedFilters[toFilterIndex][i] =
          viewVCE.vce.Extra.uncompressedFilters[fromFilterIndex][i]
    }
    viewVCE_filters.filtersChartUpdate(
        filterSelectEle.options[filterSelectEle.selectedIndex].value,
        filterSelectEle.options[filterSelectEle.selectedIndex].innerHTML, true);

    return true;
  },

  filtersChartUpdate: function(filterIndex, filterName, animate) {
    filterIndex = parseInt(filterIndex, 10);
    let datasets = [];

    console.log('Filter update: index:' + filterIndex + ' name:' + filterName);

    let filterCopySelectEle = document.getElementById('filterCopySelect');
    // remove old options:
    while (filterCopySelectEle.firstChild) {
      filterCopySelectEle.removeChild(filterCopySelectEle.firstChild);
    }
    // hide the copy selector for All or cases where there are no filters, or
    // when we're not in voicing mode
    document.querySelector('#filterCopySelectDiv').style.display = 'none';

    if (filterIndex >= 0) {
      if (viewVCE_voice.voicingMode) {
        document.querySelector('#filterCopySelectDiv').style.display = 'block';
        // populate options in the select with only "other" filters (i.e. "this"
        // filter should be not shown or at least unselectable)

        // first element is empty to avoid confusing the user if they havent
        // selected something:
        let option = document.createElement('option');
        option.value = -1;
        option.innerHTML = '';
        filterCopySelectEle.appendChild(option);

        for (let i = 0; i < viewVCE_filters.filterNames.length; i++) {
          if (viewVCE_filters.filterValues[i] >= 0 &&
              viewVCE_filters.filterValues[i] != filterIndex) {
            let option = document.createElement('option');
            option.value = viewVCE_filters.filterValues[i];
            option.innerHTML = viewVCE_filters.filterNames[i];
            filterCopySelectEle.appendChild(option);
          }
        }
      }

      document.querySelector('#filterTable').style.display = 'block';
      document.querySelectorAll('#filterTable td.val input').forEach((obj, i) => {
        let id = obj.id;
        // id is "flt[<n>]" - we need the <n> part
        let idxString = id.substring(4);
        let idx = parseInt(idxString, 10) - 1;
        obj.value = viewVCE.vce.Extra.uncompressedFilters[filterIndex][idx];
        if (animate) {
          viewVCE_voice.sendToCSurface(
              obj, id, viewVCE.vce.Extra.uncompressedFilters[filterIndex][idx]);
        }
      });
      // match the color rotation below.  We don't allocate a color for an
      // "unused" Bf. So this is senselessly complicated. Go look at the FILTERS
      // array and figure out which "compressed" index this was
      let color = viewVCE.chartColors[viewVCE.chartColors.length - 1];
      if (filterIndex > 0) {
        let colorIdx = 0
        for (let i = 0; i <= viewVCE.vce.Head.VOITAB; i++) {
          if (filterIndex == viewVCE.vce.Head.FILTER[i]) {
            color =
                viewVCE.chartColors[(colorIdx) % viewVCE.chartColors.length];
            break;
          } else if (viewVCE.vce.Head.FILTER[i] > 0) {
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
        data: viewVCE.vce.Extra.uncompressedFilters[filterIndex]
      }];
    } else {
      // "all"

      // if (animate) {
      //  set MIDI control surface to all zeros:
      for (let i = 1; i <= 32; i++) {
        let id = 'flt[' + i + ']';
        let ele = document.getElementById(id);
        viewVCE_voice.sendToCSurface(ele, id, 0);
      }
      //}

      document.querySelector('#filterTable').style.display = 'none';
      // only include the ones actually in use (see viewVCE.vce.Head.FILTER
      // entry)
      for (let i = 0; i <= viewVCE.vce.Head.VOITAB; i++) {
        // A-table if in use goes first:
        if (viewVCE.vce.Head.FILTER[i] < 0) {
          let filterName = 'Af';
          // we use modulo to compute b-filter color - since A filter is "-1",
          // use the last color in the table
          let color = viewVCE.chartColors[viewVCE.chartColors.length - 1]
          datasets.push({
            fill: false,
            lineTension: 0,
            pointRadius: 0,
            pointHitRadius: 5,
            label: filterName,
            backgroundColor: color,
            borderColor: color,
            data: viewVCE.vce.Extra
                      .uncompressedFilters[0]  // ZEROTH ele is the a-filter
          });
        }
      }
      let colorIdx = 0;
      for (let i = 0; i <= viewVCE.vce.Head.VOITAB; i++) {
        // A-table if in use goes first:
        if (viewVCE.vce.Head.FILTER[i] > 0) {
          let filterName = 'Bf ' + (i + 1);
          let color =
              viewVCE.chartColors[colorIdx % viewVCE.chartColors.length];
          colorIdx++;
          datasets.push({
            fill: false,
            lineTension: 0,
            pointRadius: 0,
            pointHitRadius: 5,
            label: filterName,
            backgroundColor: color,
            borderColor: color,
            data: viewVCE.vce.Extra
                      .uncompressedFilters[i + 1]  // b-filters by one-based osc
                                                   // value
          });
        }
      }
      // now B-filters:
    }

    //		console.dir(datasets);

    let ctx = document.getElementById('filtersChart').getContext('2d');
    if (viewVCE_filters.chart != null) {
      viewVCE_filters.chart.destroy();
    }

    let animation_duration = animate ? 1000 : 0;

    viewVCE_filters.chart = new Chart(ctx, {

      type: 'line',
      data: {
        labels: [
          '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '',
          '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', ''
        ],
        datasets: datasets
      },

      // Configuration options go here
      options: {
        animation: {duration: animation_duration},

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
            scaleLabel: {display: true, labelString: 'Frequency'},
            ticks: {color: '#666', display: true}
          }],
          yAxes: [{
            id: 'y-axis',
            gridLines: {
              color: '#666',
              display: true,
              drawBorder: false,
              drawOnChartArea: true
            },
            scaleLabel: {display: true, labelString: 'dB'},
            ticks: {color: '#eee', display: true}
          }],
        },
        responsive: false,
        maintainAspectRatio: false,
        plugins: {
          // zoom plugin is only used by the env graphs
          zoom: {zoom: {enabled: false}, pan: {enabled: false}}
        }
      }
    });
    if (viewVCE_voice.voicingMode && filterIndex >= 0) {
      viewVCE_chartdrag.init(
          viewVCE_filters.chart, viewVCE_filters.onchange, 'flt', 0, 31, -64,
          63);
    }
  }
};
