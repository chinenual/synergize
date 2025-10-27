// Generic drag support for the simple filter, keyeq and keyprop style graphs.
// Can't use for envelopes since those are a lot more complicated.

//ChartJs event handler attaching events to chart canvas
var viewVCE_chartdrag = {
    registeredCharts: [],

    deregister: function (chart) {
        console.log("deregister: ", chart, viewVCE_chartdrag.registeredCharts[chart])
        cleanup = viewVCE_chartdrag.registeredCharts[chart];
        if (cleanup != null) {
            cleanup(chart);
        }
        viewVCE_chartdrag.registeredCharts[chart] = null;
    },

    //Call init with a ChartJs Chart instance to apply mouse and touch events to its canvas.
    init: function (chart, onchange, fieldName, xMin, xMax, yMin, yMax) {
            viewVCE_chartdrag.dragging = false;
            viewVCE_chartdrag.points = [];

            viewVCE_chartdrag.deregister(chart);

            var state = {
                onchange: onchange,
                fieldName: fieldName,
                chart: chart,
                xMin: xMin,
                xMax: xMax,
                yMin: yMin,
                yMax: yMax,
            }
            //Event handler for event types subscribed
            var evtHandler =
                function myeventHandler(evt) {
                    var cancel = false;
                    switch (evt.type) {
                        case "mousedown":
                        case "touchstart":
                            cancel = viewVCE_chartdrag.onDragStart(evt, state);
                            break;
                        case "mousemove":
                        case "touchmove":
                            cancel = viewVCE_chartdrag.onDrag(evt, state);
                            break;
                        case "mouseup":
                        case "touchend":
                            cancel = viewVCE_chartdrag.onDragEnd(evt, state);
                            break;
                        case "mouseout":
                            console.log("mouseout");
                            cancel = viewVCE_chartdrag.onDragEnd(evt, state);
                        default:
                        //handleDefault(evt);
                    }

                    if (cancel) {
                        //Prevent the event e from bubbling up the DOM
                        if (evt.cancelable) {
                            state.dragging = false;
                            if (evt.preventDefault) {
                                evt.preventDefault();
                            }
                            if (evt.cancelBubble != null) {
                                evt.cancelBubble = true;
                            }
                        }
                    }
                };
            //Events to subscribe to
            var events = ['mousedown', 'touchstart', 'mousemove', 'touchmove', 'mouseup', 'touchend'];

            //Subscribe events
            events.forEach(function (evtName) {
                chart.chart.canvas.addEventListener(evtName, evtHandler);
            })

            var cleanup = function(chart) {
                events.forEach(function (evtName) {
                    chart.chart.canvas.removeEventListener(evtName, evtHandler);
                })
            }
            viewVCE_chartdrag.registeredCharts[chart] = cleanup;

    }

    ,

    scaleValues: function (e, state) {
        var idx = Math.round(state.chart.scales['x-axis'].getValueForPixel(e.layerX));
        idx = Math.min(idx, state.xMax);
        idx = Math.max(idx, state.xMin);

        var value = Math.round(state.chart.scales['y-axis'].getValueForPixel(e.layerY));
        value = Math.min(value, state.yMax);
        value = Math.max(value, state.yMin);

        return {
            idx: idx,
            value: value
        };
    }
    ,


    undefined,
    dragging
:
undefined,

    onDragStart
:

function (e, state) {
    viewVCE_chartdrag.points = [];
    console.log("onDragStart: ", e, state, viewVCE_chartdrag.scaleValues(e, state))

    viewVCE_chartdrag.dragging = true;
    var val = viewVCE_chartdrag.scaleValues(e, state)
    viewVCE_chartdrag.updateField(state, val.idx, val.value)
    return true;

}

,


onDrag: function (e, state) {

    if (viewVCE_chartdrag.dragging) {
        //console.log("onDrag: ", e, state, viewVCE_chartdrag.scaleValues(e, state))
        var val = viewVCE_chartdrag.scaleValues(e, state)
        viewVCE_chartdrag.updateField(state, val.idx, val.value)
    }
    return true;
}
,


onDragEnd: function (e, state) {
    viewVCE_chartdrag.dragging = false;
    console.log("onDragEnd: ", e, state, viewVCE_chartdrag.scaleValues(e, state))
    var val = viewVCE_chartdrag.scaleValues(e, state)
    viewVCE_chartdrag.updateField(state, val.idx, val.value)
    /*        function sleep(seconds){
                console.log("pause "+seconds+"s...");
                var waitUntil = new Date().getTime() + seconds*1000;
                while(new Date().getTime() < waitUntil) true;
            }
            sleep(2);
    */
    for (i = 0; i < viewVCE_chartdrag.points.length; i++) {
        if (viewVCE_chartdrag.points[i] != undefined) {
            ele = document.getElementById(`${state.fieldName}[${i + 1}]`);
            state.onchange(ele,false)
        }
    }
}
,

updateField: function (state, idx, value) {
    viewVCE_chartdrag.points[idx] = value;
    state.chart.data.datasets[0].data[idx] = value;
    state.chart.update(0);
    ele = document.getElementById(`${state.fieldName}[${idx + 1}]`, value);
    ele.value = value;
}
}

