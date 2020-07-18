let crt = {};
let crt_path = "";
let crt_name = "";

let viewCRT = {
	editMode: false,

	makeSlotOnclick(slot) {
		return function () {
			index.viewVCESlot(slot);
		}
	},
	editCRT: function (name, path) {
		if (viewCRT.editMode) {
			index.confirmDialog("Disabling edit mode will discard any pending edits. Are you sure?", function () {
				viewCRT.editMode = false;
				document.getElementById("editCRTButtonImg").src = `static/images/red-button-off-full.png`;
				$("#saveCRTButton").hide();
			});
		} else {
			viewCRT.editMode = true;
			document.getElementById("editCRTButtonImg").src = `static/images/red-button-on-full.png`;
			$("#saveCRTButton").show();
		}
	},
	saveCRT: function (name, path) {
	},
	loadCRT: function (name, path) {
		index.confirmDialog("Load Voice Cartridge file " + path, function () {
			let message = {
				"name": "loadCRT",
				"payload": path
			};
			index.spinnerOn();
			astilectron.sendMessage(message, function (message) {
				index.spinnerOff();
				// Check error
				if (message.name === "error") {
					index.errorNotification(message.payload);
				} else {
					index.infoNotification("Successfully loaded " + name + " to Synergy")
				}
				index.refreshConnectionStatus();
			});
		}
		);
	},
	init: function () {
		viewCRT.editMode = false;
		document.getElementById("editCRTButtonImg").src = `static/images/red-button-off-full.png`;
		$("#saveCRTButton").hide();

		console.log("view CRT " + crt_name)
		document.getElementById("crt_path").innerHTML = crt_name;
		// clear everything
		for (i = 0; i < 24; i++) {
			var ele = document.getElementById("crt_voicename_" + (i + 1));
			ele.innerHTML = "";
			ele.onclick = function () {
				// nop
			};
		}
		// set active voices
		for (i = 0; i < crt.Voices.length; i++) {
			var ele = document.getElementById("crt_voicename_" + (i + 1));
			ele.innerHTML = crt.Voices[i].Head.VNAME;
			ele.onclick = viewCRT.makeSlotOnclick(i);
		}
	}
};
