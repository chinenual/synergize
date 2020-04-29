let crt = {};
let crt_path = "";

let viewCRT = {
    refreshText: function () {
	document.getElementById("crt_path").innerHTML = crt_path;
	for (i = 0; i < 24; i++) {
	    document.getElementById("crt_voicename_" + (i+1)).innerHTML = "";
	}
	for (i = 0; i < crt.Voices.length; i++) {
	    document.getElementById("crt_voicename_" + (i+1)).innerHTML = crt.Voices[i].VNAME;
	}
    }
};
