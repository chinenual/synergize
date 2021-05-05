const { dialog } = require('electron').remote;

let shell = require('electron').shell

const DEBOUNCE_WAIT_SHORT = 50;
const DEBOUNCE_WAIT = 250;

let index = {
	init: function () {
		dx2syn.init();

		// make sure external web links open in system browser - not the application:
		document.addEventListener('click', function (event) {
			if (event.target.tagName === 'A' && event.target.href.startsWith('http')) {
				event.preventDefault()
				shell.openExternal(event.target.href)
			}
		})
		// Wait for astilectron to be ready
		document.addEventListener('astilectron-ready', function () {
			let message = {
				"name": "isHTTPDebug",
				"payload": ""
			};
			// Send message
			astilectron.sendMessage(message, function (message) {
				// Check error
				if (message.name === "error") {
					index.errorNotification(message.payload);
				} else {
					if (message.payload) {
						debug.configDebugContextMenu();
					}
				}
			});

			// init menus to default state
			index.updateConnectionStatus("","")
			// Listen
			index.listen();
			// Explore default path
			index.explore();
		})

	},
	checkInputElementValue: function (ele) {
		if (!ele.value.match(/^-?\d+$/)) {
			return undefined;
		}
		var result = parseInt(ele.value, 10);
		if (ele.hasAttribute("min")) {
			var min = parseInt(ele.getAttribute("min"), 10);
			if (result < min) result = min;
		}
		if (ele.hasAttribute("max")) {
			var max = parseInt(ele.getAttribute("max"), 10);
			if (result > max) result = max;
		}
		return result;
	},

	spinnerOn: function () {
		document.getElementById("spinner").style.display = "block";
	},
	spinnerOff: function () {
		document.getElementById("spinner").style.display = "none";
	},

	confirmDialog: function (message, successCallback) {
		if (typeof message != 'string') {
			message = JSON.stringify(message);
		}
		console.log("CONFIRM DIALOG: " + message)
		document.getElementById("confirmTitle").innerHTML = "Confirm";
		document.getElementById("confirmText").innerHTML = message;
		document.getElementById("confirmOKButton").onclick = successCallback;
		$('#confirmModal').modal({
			backdrop: "static" // clicking outside the dialog doesnt close the dialog
		});
	},

	errorNotification: function (message) {
		if (typeof message != 'string') {
			message = JSON.stringify(message);
		}
		console.log("ERROR NOTIFICATION: " + message)
		document.getElementById("alertTitle").innerHTML = "Error";
		document.getElementById("alertText").innerHTML = message;
		$('#alertModal').modal();
	},
	infoNotification: function (message) {
		if (typeof message != 'string') {
			message = JSON.stringify(message);
		}
		console.log("INFO NOTIFICATION: " + message)
		document.getElementById("alertTitle").innerHTML = "Info";
		document.getElementById("alertText").innerHTML = message;
		$('#alertModal').modal();
	},

	chooseZeroconfService: function (prompt1, choices1, prompt2, choices2, onCancel, onOK, onRescan) {
		console.log("chooseZeroconfService: " + prompt1 + " " + JSON.stringify(choices1) + " " + prompt2 + " " + JSON.stringify(choices2));
		if (prompt1 != null) {
			document.getElementById("chooseZeroconf1Prompt").innerHTML = prompt1;
			$('#zeroconf1Div').show();
		} else {
			document.getElementById("chooseZeroconf1Prompt").innerHTML = "";
			document.getElementById("chooseZeroconf1Items").innerHTML = "";
			$('#zeroconf1Div').hide();
		}
		if (prompt2 != null) {
			document.getElementById("chooseZeroconf2Prompt").innerHTML = prompt2;
			$('#zeroconf2Div').show();
		} else {
			document.getElementById("chooseZeroconf2Prompt").innerHTML = "";
			document.getElementById("chooseZeroconf2Items").innerHTML = "";
			$('#zeroconf2Div').hide();
		}
		var html = "";
		if (prompt1 != null && choices1  != null) {
			for (i = 0; i < choices1.length; i++) {
				var addr = ""
				if (choices1[i].Port != 0) {
					addr = ` (${choices1[i].HostName}:${choices1[i].Port})`
				}
				html = html + `
		    <div class="form-check">
                <input class="form-check-input" type="radio" name="chooseZeroconf1Radios" id="chooseZeroconf1Radio${i}" value="${i}" ${i == 0 ? "checked" : ""}>
                <label class="form-check-label" for="chooseZeroconf1Radio${i}">
                   ${choices1[i].InstanceName}${addr}
                </label>
			</div>`
				console.log("html now " + html);
			}
			document.getElementById("chooseZeroconf1Items").innerHTML = html;
			console.log("innerHTML now " + document.getElementById("chooseZeroconf1Items").innerHTML);
		}

		if (prompt2 != null && choices2  != null) {
			var html = "";
			for (i = 0; i < choices2.length; i++) {
				var addr = ""
				if (choices2[i].Port != 0) {
					addr = ` (${choices2[i].HostName}:${choices2[i].Port})`
				}
				html = html + `
		    <div class="form-check">
                <input class="form-check-input" type="radio" name="chooseZeroconf2Radios" id="chooseZeroconf2Radio${i}" value="${i}" ${i == 0 ? "checked" : ""}>
                <label class="form-check-label" for="chooseZeroconf2Radio${i}">
                   ${choices2[i].InstanceName}${addr}
                </label>
			</div>`
				console.log("html now " + html);
			}
			document.getElementById("chooseZeroconf2Items").innerHTML = html;
			console.log("innerHTML now " + document.getElementById("chooseZeroconf2Items").innerHTML);
		}

		var selected = 0;
		document.getElementById("chooseZeroconfCancelButton").onclick = function () {
			console.log("Cancelled");
			onCancel();
		};
		document.getElementById("chooseZeroconfOKButton").onclick = function () {
			var idx1 = null
			var selected1 = null
			var idx2 = null
			var selected2 = null
			if (choices1 != null && prompt1 != null) {
				idx1 = parseInt($('#chooseZeroconf1Items input:checked').val(), 10);
				selected1 = choices1[idx1];
			}
			if (choices2 != null && prompt2 != null) {
				idx2 = parseInt($('#chooseZeroconf2Items input:checked').val(), 10);
				selected2 = choices2[idx2];
			}
			console.log("Selected " + idx1 + " " + idx2);
			onOK(selected1, selected2);
		};
		document.getElementById("chooseZeroconfRescanButton").onclick = function () {
			console.log("Rescan");
			onRescan();
		};

		$('#chooseZeroconfModal').modal({
			backdrop: "static" // clicking outside the dialog doesnt close the dialog
		});
	},

	saveSYNDialog: function () {

		path = dialog.showSaveDialogSync({
			filters: [
				{ name: 'State', extensions: ['syn'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile', 'promptToCreate']
		});
		console.log("in fileDialog: " + path);

		if (path != undefined) {
			viewVCE_voice.connectSynergy(function () {

				let message = {
					"name": "saveSYN",
					"payload": path
				};
				// Send message
				index.spinnerOn();
				astilectron.sendMessage(message, function (message) {
					index.spinnerOff();
					// Check error
					if (message.name === "error") {
						index.errorNotification(message.payload);
					} else {
						index.infoNotification("Successfully saved Synergy state to " + path);
					}
					index.refreshConnectionStatus();
				});
			});
		}
	},
	loadSYNDialog: function () {
		path = dialog.showOpenDialogSync({
			filters: [
				{ name: 'State', extensions: ['syn'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile']
		});
		console.log("in fileDialog: " + path);
		if (path != undefined) {
			index.loadSYN(path[0], path[0]);
		}
	},
	loadCRTDialog: function () {
		path = dialog.showOpenDialogSync({
			filters: [
				{ name: 'Cartridge', extensions: ['crt'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile']
		});
		console.log("in fileDialog: " + path);
		if (path != undefined) {
			index.viewCRT(path[0], path[0]);
		}
	},
	loadVCEDialog: function () {
		path = dialog.showOpenDialogSync({
			filters: [
				{ name: 'Voice', extensions: ['vce'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile']
		});
		console.log("in fileDialog: " + path);
		if (path != undefined) {
			index.viewVCE(path[0], path[0]);
		}
	},
	saveVCEDialog: function () {
		path = dialog.showSaveDialogSync({
			filters: [
				{ name: 'Voice', extensions: ['vce'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile', 'promptToCreate']
		});
		console.log("in saveVCEDialog: " + path);
		if (path != undefined) {
			let message = {
				"name": "saveVCE",
				"payload": path
			};
			// Send message
			index.spinnerOn();
			astilectron.sendMessage(message, function (message) {
				index.spinnerOff();
				// Check error
				if (message.name === "error") {
					index.errorNotification(message.payload);
				} else {
					index.infoNotification("Successfully saved Synergy voice file to " + path);
				}
				index.refreshConnectionStatus();
			});
		}
	},
	loadSYN: function (name, path) {
		viewVCE_voice.connectSynergy(function () {
			index.confirmDialog("Load Synergy state file " + path, function () {
				let message = {
					"name": "loadSYN",
					"payload": path
				};
				// Send message
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
			});
		});
	},
	viewCRT: function (name, path) {
		if (viewVCE_voice.voicingMode) {
			index.errorNotification("Can't load a CRT file while in Voicing mode");
			return;
		}
		let message = {
			"name": "readCRT",
			"payload": path
		};
		astilectron.sendMessage(message, function (message) {
			// Check error
			if (message.name === "error") {
				index.errorNotification(message.payload);
				return
			}
			crt_path = path;
			crt_name = name;
			crt = message.payload;
			index.load("viewCRT.html", "content",
				function () {
					viewCRT.init();
				});
			index.refreshConnectionStatus();
		});
	},
	viewVCE: function (name, path) {
		console.log("index.viewVCE " + name + " " + path)
		if (viewVCE_voice.voicingMode) {
			index.confirmDialog("Loading voice file will overwrite any pending edits - continue?", function () {
				index.raw_viewVCE(name, path);
			});
		} else {
			index.raw_viewVCE(name, path);
		}
	},
	raw_viewVCE: function (name, path) {
		console.log("index.raw_viewVCE " + name + " " + path)
		var name = "readVCE";
		if (viewVCE_voice.voicingMode) {
			name = "loadVceVoicingMode"
		}
		let message = {
			"name": name,
			"payload": path
		};
		// Send message
		index.spinnerOn();
		astilectron.sendMessage(message, function (message) {
			index.spinnerOff();
			// Check error
			if (message.name === "error") {
				index.errorNotification(message.payload);
				return
			}
			vce = message.payload;

			crt_name = null;
			crt_path = null;
			index.load("viewVCE.html", "content",
				function () {
					viewVCE.init();
				});
			index.refreshConnectionStatus();
		});
	},
	viewVCESlot: function (slot) {
		vce = crt.Voices[slot];

		console.log("view voice slot " + slot + " : " + vce);
		index.load("viewVCE.html", "content",
			function () {
				viewVCE.init();
			});
	},
	addFolder: function (name, path) {
		let div = document.createElement("div");
		div.className = "dir";
		div.onclick = function () { index.explore(path) };
		if (name == "..") name = "&lt;Parent&gt;";
		div.innerHTML = `<i class="fa fa-folder"></i><span>` + name + `</span>`;
		document.getElementById("dirs").appendChild(div)
	},
	addSYNFile: function (name, path) {
		let div = document.createElement("div");
		div.className = "file";
		div.onclick = function () { index.loadSYN(name, path) };
		div.innerHTML = `<i class="fa fa-file"></i><span>` + name + `</span>`;
		document.getElementById("SYNfiles").appendChild(div)
	},
	addCRTFile: function (name, path) {
		let div = document.createElement("div");
		div.className = "file";
		div.onclick = function () { index.viewCRT(name, path) };
		div.innerHTML = `<i class="fa fa-file"></i><span>` + name + `</span>`;
		document.getElementById("CRTfiles").appendChild(div)
	},
	addVCEFile: function (name, path) {
		let div = document.createElement("div");
		div.className = "file";
		div.onclick = function () { index.viewVCE(name, path) };
		div.innerHTML = `<i class="fa fa-file"></i><span>` + name + `</span>`;
		document.getElementById("VCEfiles").appendChild(div)
	},
	explore: function (path) {
		// Create message
		let message = { "name": "explore" };
		if (typeof path !== "undefined") {
			message.payload = path
		}

		// Send message
		astilectron.sendMessage(message, function (message) {
			// Check error
			if (message.name === "error") {
				index.errorNotification(message.payload);
				return
			}

			// Process path
			document.getElementById("path").innerHTML = message.payload.path;

			// Process dirs
			document.getElementById("dirs").innerHTML = ""
			for (let i = 0; i < message.payload.dirs.length; i++) {
				index.addFolder(message.payload.dirs[i].name, message.payload.dirs[i].path);
			}

			document.getElementById("CRTfiles").innerHTML = ""
			if (message.payload.CRTfiles.length > 0) {
				let div = document.createElement("div")
				div.innerHTML = "<div class='horizSeparator'></div><b>Cartridge Files (.CRT)</b>";
				document.getElementById("CRTfiles").appendChild(div);

				for (let i = 0; i < message.payload.CRTfiles.length; i++) {
					index.addCRTFile(message.payload.CRTfiles[i].name, message.payload.CRTfiles[i].path);
				}
			}

			document.getElementById("SYNfiles").innerHTML = ""
			if (message.payload.SYNfiles.length > 0) {
				let div = document.createElement("div")
				div.innerHTML = "<div class='horizSeparator'></div><b>Synergy State (.SYN)</b>";
				document.getElementById("SYNfiles").appendChild(div);
				for (let i = 0; i < message.payload.SYNfiles.length; i++) {
					index.addSYNFile(message.payload.SYNfiles[i].name, message.payload.SYNfiles[i].path);
				}
			}

			document.getElementById("VCEfiles").innerHTML = ""
			if (message.payload.VCEfiles.length > 0) {
				let div = document.createElement("div")
				div.innerHTML = "<div class='horizSeparator'></div><b>Voice Files (.VCE)</b>";
				document.getElementById("VCEfiles").appendChild(div);
				for (let i = 0; i < message.payload.VCEfiles.length; i++) {
					index.addVCEFile(message.payload.VCEfiles[i].name, message.payload.VCEfiles[i].path);
				}
			}
		})
	},
	disconnectSynergy: function () {
		if (viewVCE_voice.voicingMode) {
			index.confirmDialog("Disconnecting the Synergy will will discard any pending edits. Are you sure?", function () {
				viewVCE_voice.raw_voicingModeOff(true);
			});
		} else {
			index.raw_disconnectSynergy();
		}
	},

	raw_disconnectSynergy: function() {
		let message = { "name": "disconnectSynergy" };
		index.spinnerOn();
		astilectron.sendMessage(message, function (message) {
			index.spinnerOff();
			if (message.name === "error") {
				index.errorNotification(message.payload);
				return
			} else {
				index.updateConnectionStatus(message.payload.SynergyName, message.payload.ControlSurfaceName);
				index.infoNotification("Disconnected Synergy");
				return
			}
		});
	},

	disconnectControlSurface: function () {
		let message = { "name": "disconnectControlSurface" };
		index.spinnerOn();
		astilectron.sendMessage(message, function (message) {
			index.spinnerOff();
			if (message.name === "error") {
				index.errorNotification(message.payload);
				return
			} else {
				index.updateConnectionStatus(message.payload.SynergyName, message.payload.ControlSurfaceName);
				index.infoNotification("Disconnected Control Surface");
				$('#disableControlSurfaceMenuItem').addClass('disabled');
				return
			}
		});
	},
	disableVRAM: function () {
		viewVCE_voice.connectSynergy(function () {
			let message = { "name": "disableVRAM" };
			index.spinnerOn();
			astilectron.sendMessage(message, function (message) {
				index.spinnerOff();
				if (message.name === "error") {
					index.errorNotification(message.payload);
				} else {
					index.infoNotification("Successfully disabled Synergy's VRAM")
				}
				index.refreshConnectionStatus();
			});
		});
	},
	refreshConnectionStatus: function () {
		let message = {
			"name": "getConnectionStatus",
			"payload": "DummyPayload"
		};
		// Send message
		console.log("refreshing connection status");
		astilectron.sendMessage(message, function (message) {
			// Check error
			if (message.name === "error") {
				index.errorNotification(message.payload);
			} else {
				index.updateConnectionStatus(message.payload.SynergyName, message.payload.ControlSurfaceName);
			}
		});
	},

	synergyName: null,
	controlSurfaceName: null,

	updateConnectionStatus: function (synergyName, csName) {
		index.synergyName = (synergyName == null || synergyName === "") ? null : synergyName;
		index.contronSurfaceName = (csName == null || csName === "") ? null : csName;

		console.log("update status: " + synergyName + " " + csName);
		document.getElementById("synergyName").innerHTML = synergyName;
		document.getElementById("controlSurfaceName").innerHTML = csName;
		if (synergyName === null || synergyName === "") {
			document.getElementById("synergyName").innerHTML = "not connected";
			$('#disconnectSynergyMenuItem').addClass('disabled');
			$('#connectSynergyMenuItem').removeClass('disabled');
			document.getElementById("connectButtonImg").src = `static/images/grey-button-off-full.png`;
		} else {
			$('#disconnectSynergyMenuItem').removeClass('disabled');
			$('#connectSynergyMenuItem').addClass('disabled');
			document.getElementById("connectButtonImg").src = `static/images/grey-button-on-full.png`;
		}
		if (csName === null || csName === "") {
			$('#controlSurfaceStatus').hide();
			$('#disconnectControlSurfaceMenuItem').addClass('disabled');
		} else {
			$('#controlSurfaceStatus').show();
			$('#disconnectControlSurfaceMenuItem').removeClass('disabled');
		}
	},

	checkVersion: function(synergyWasDisconnected, controlSurfaceWasDisconnected) {
		console.log("checkVersion " + synergyWasDisconnected + " " + controlSurfaceWasDisconnected);
		let message = { "name": "checkVersion",
			"payload": {
			"SynergyWasDisconnected" : synergyWasDisconnected,
				"ControlSurfaceWasDisconnected": controlSurfaceWasDisconnected
			}
		};
		astilectron.sendMessage(message, function (message) {
			if (message.name === "error") {
				index.errorNotification(message.payload);
				return
			} else {
				return
			}
		});
	},

	fileDialog: function () {
		files = dialog.showOpenDialogSync({
			//electron bug? filter files cause the dialog to look wonky
			filters: [
				{ name: 'Voice', extensions: ['vce'] },
				{ name: 'Cartridge', extensions: ['crt'] },
				{ name: 'State', extensions: ['syn'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile']
		});
		console.log("in fileDialog: " + files);
		return files;
	},
	runCOMTST: function () {
		viewVCE_voice.connectSynergy(function () {
			let message = { "name": "runCOMTST" };
			index.spinnerOn();
			astilectron.sendMessage(message, function (message) {
				index.spinnerOff();
				console.log("runCOMTST returned: " + JSON.stringify(message));
				// Check error
				if (message.name === "error") {
					index.errorNotification(message.payload);
				} else {
					index.infoNotification(message.payload);
				}
			});
		});
		index.refreshConnectionStatus();
	},

	load: function (url, eleId, callback) {
		console.log("load " + url + " into " + eleId + " " + $(('#' + eleId)));
		//		console.dir($(('#' + eleId)));
		$(("#" + eleId)).load(url, function () {
			console.log("loaded url " + url);
			if (callback != undefined) {
				element = document.getElementById(eleId);
				callback(element);
			}
		});

		/*
		timing bug - onreadystatechange fires before the DOM is ready to query
		
		element = document.getElementById(eleId);
		req = new XMLHttpRequest();
		
		req.onreadystatechange = function () {
			if (this.readyState == 4 && this.status == 200) {
				element.innerHTML = req.responseText;
				if (callback != undefined) {
					callback(element);
				}
			}
		};
		
		req.open("GET", url, false);
		req.send(null);
		*/
	},
	dropdownMenu: function (contentId) {
		//console.log("toggle display on " + contentId);
		document.getElementById(contentId).style.display = "block";
	},

	viewDiag: function () {
		index.load("diag.html", "content");
	},
	showAbout: function () {
		let message = { "name": "showAbout" };
		astilectron.sendMessage(message, function (message) {
			// nop
		});
	},
	showPreferences: function () {
		let message = { "name": "showPreferences" };
		//console.log("show preferences javascript");
		astilectron.sendMessage(message, function (message) {
			// nop
		});
	},
	showTunings: function () {
		viewVCE_voice.connectSynergy(function () {
			let message = {"name": "showTunings"};
			//console.log("showTunings javascript");
			astilectron.sendMessage(message, function (message) {
				// nop
			});
		});
	},

	// debounce a function separately for each "first" argument - we use this
	// with first argument being the input ele being debounced - this allows
	// each input to be independently debounced even if all using the same onchange function
	// Adapted from: https://github.com/lodash/lodash/issues/2403 and https://stackoverflow.com/a/28795512
	debounceFirstArg: function (func, wait = 0, options = {}) {
		var mem = _.memoize(function () {
			return _.debounce(func, wait, options)
		});
		return function () { mem.apply(this, arguments).apply(this, arguments) }
	},

	listen: function () {
		console.log("index listening...")
		astilectron.onMessage(function (message) {
			switch (message.name) {
				case "explore":
					index.explore(message.payload);
					return { payload: "ok" };
				case "updateConnectionStatus":
					index.updateConnectionStatus(message.payload.SynergyName, message.payload.ControlSurfaceName);
					return { payload: "ok" };
				case "fileDialog":
					f = index.fileDialog(message.payload);
					return { payload: f };
					break;
				case "viewVCE":
					console.log("viewVCE: " + JSON.stringify(message.payload));
					vce = message.payload;
					index.load("view.html", "content",
						function () {
							viewVCE.init();
						});

					return { payload: "ok" };
					break;
				case "runDiag":
					index.viewDiag();
					return { payload: "ok" };
					break;
				case "updateFromCSurface":
					valueString = viewVCE_voice.updateFromCSurface(message.payload)
					return { payload: valueString };

				case "dx2synAddProcessLog":
					console.log("dx2synAddProcessLog  - " + message.payload)
					dx2syn.addProcessLog(message.payload);
					return { payload: "ok" };

				case "dx2synFinish":
					console.log("dx2synFinish  - " + message.payload)
					dx2syn.finishConvert(message.payload);
					return { payload: "ok" };

			}
		});
	}
};

function inDropbtn(ele) {
	if (ele == null) {
		//console.log("ele is null");
		return false
	} else if (ele.classList.contains('dropbtn')) {
		//console.log("ele has dropbtn " + JSON.stringify(ele));
		return true;
	}
	//console.log("ele doesnt have dropbtn - try parent " + JSON.stringify(ele) + " " + JSON.stringify(ele.parentElement));
	return inDropbtn(ele.parentElement);
}

/* close dropdowns if user clicks outside the menu */
window.onclick = function (event) {
	if (!inDropbtn(event.target)) {
		var dropdowns = document.getElementsByClassName("dropdown-content");
		var i;
		for (i = 0; i < dropdowns.length; i++) {
			var openDropdown = dropdowns[i];
			if (openDropdown.style.display === "block") {
				//console.log("toggle display off " + openDropdown.id);
				openDropdown.style.display = "none";
			}
		}
	}
}


