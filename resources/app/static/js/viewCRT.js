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
				viewCRT.reinit();
			});
		} else {
			viewVCE_voice.connectSynergy(function () {
				viewCRT.editMode = true;
				viewCRT.reinit();
			});
		}
	},

	add: function (slot) {
		path = dialog.showOpenDialogSync({
			filters: [
				{ name: 'Voice', extensions: ['vce'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile']
		});
		console.log("in fileDialog: " + path);
		if (path != undefined) {
			let message = {
				"name": "crtEditAddVoice",
				"payload": {
					"VcePath": path[0],
					"Slot": slot,
					"Crt": crt
				}
			};
			console.dir(message.payload);
			astilectron.sendMessage(message, function (message) {
				// Check error
				if (message.name === "error") {
					index.errorNotification(message.payload);
					return
				}
				console.dir(message.payload);
				crt = message.payload;
				viewCRT.reinit();
				index.refreshConnectionStatus();
			});
		}
	},

	clear: function (slot) {
		var ele = document.getElementById("crt_voicename_" + slot);
		ele.innerHTML = "";
		ele.onclick = function () {
			// nop
		};
		crt.Voices[slot - 1] = null;
	},

	loadCRT: function () {

		if (path != undefined) {
			viewVCE_voice.connectSynergy(function () {

				let message = {
					"name": "crtEditLoadCRT",
					"payload": {
						Crt: crt
					}
				};
				// Send message
				index.spinnerOn();
				astilectron.sendMessage(message, function (message) {
					index.spinnerOff();
					// Check error
					if (message.name === "error") {
						index.errorNotification(message.payload);
					} else {
						index.infoNotification("Successfully loaded CRT to Synergy");
					}
					index.refreshConnectionStatus();
				});
			});
		}
	},

	saveCRT: function (name, path) {
		path = dialog.showSaveDialogSync({
			filters: [
				{ name: 'Cartridge', extensions: ['CRT'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile', 'promptToCreate']
		});
		console.log("in fileDialog: " + path);

		if (path != undefined) {
			viewVCE_voice.connectSynergy(function () {

				let message = {
					"name": "crtEditSaveCRT",
					"payload": {
						Path: path,
						Crt: crt
					}
				};
				// Send message
				index.spinnerOn();
				astilectron.sendMessage(message, function (message) {
					index.spinnerOff();
					// Check error
					if (message.name === "error") {
						index.errorNotification(message.payload);
					} else {
						index.infoNotification("Successfully saved CRT to " + path);
					}
					index.refreshConnectionStatus();
				});
			});
		}
	},

	viewLoadedCRT: function () {
		index.load("viewCRT.html", "content",
			function () {
				viewCRT.reinit();
			});
	},

	reinit: function () {
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
			if (crt.Voices[i] != null) {
				ele.innerHTML = crt.Voices[i].Head.VNAME;
				ele.onclick = viewCRT.makeSlotOnclick(i);
			}
		}

		if (viewCRT.editMode) {
			document.getElementById("editCRTButtonImg").src = `static/images/red-button-on-full.png`;
			$("#saveCRTButtonDiv").show();
			$(".crtSlotAddButton").show();
			$(".crtSlotClearButton").show();
		} else {
			document.getElementById("editCRTButtonImg").src = `static/images/red-button-off-full.png`;
			$("#saveCRTButtonDiv").hide();
			$(".crtSlotAddButton").hide();
			$(".crtSlotClearButton").hide();
		}

	},

	init: function () {
		viewCRT.editMode = false;
		viewCRT.reinit();
	}
};
