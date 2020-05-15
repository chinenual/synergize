let crt = {};
let crt_path = "";
let crt_name = "";

let viewCRT = {
    makeSlotOnclick(slot) {
	return function() {
	    index.viewVCESlot(slot);
	}
    },
    init: function () {
	document.getElementById("crt_path").innerHTML = crt_name;
	// clear everything
	for (i = 0; i < 24; i++) {
	    var ele = document.getElementById("crt_voicename_" + (i+1));
	    ele.innerHTML = "";
	    ele.onclick = function() {
		// nop
	    };
	}
	// set active voices
	for (i = 0; i < crt.Voices.length; i++) {
	    var ele = document.getElementById("crt_voicename_" + (i+1));
	    ele.innerHTML = crt.Voices[i].Head.VNAME;
	    ele.onclick = viewCRT.makeSlotOnclick(i);
	}
    }
};
