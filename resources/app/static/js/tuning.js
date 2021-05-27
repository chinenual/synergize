let tuning = {
	init: function() {
		console.log("tuning.init()")
		console.log("tuning table",document.getElementById("tuningTable"))
		document.getElementById("tuningTable").hidden = true;

		// Wait for astilectron to be ready
		document.addEventListener('astilectron-ready', function () {
			console.log("tuning astilectron-ready")

			let message = {
				"name": "getTuningParams",
				"payload": ""
			};
			// Send message
			astilectron.sendMessage(message, function (message) {
				// Check error
				if (message.name === "error") {
					index.errorNotification(message.payload);
					return
				}
				var prefs = message.payload

				//console.log("loaded preferences: " + os + ": " + JSON.stringify(preferences))
				document.getElementById("useStandard").checked = prefs.UseStandardTuning;
				document.getElementById("useStandardKBM").checked = prefs.UseStandardKeyboardMapping;
				document.getElementById("sclPath").value = prefs.SCLPath;
				document.getElementById("kbmPath").value = prefs.KBMPath;
				document.getElementById("middleNote").value = prefs.MiddleNote;

				document.getElementById("referenceNote").value = prefs.ReferenceNote;
				document.getElementById("referenceFreq").value = prefs.ReferenceFrequency;

				tuning.toggle();

			});

			// Listen
			tuning.listen();
		})
	},

	kbmFileDialog: function (ele, defaultValue) {
		folder = dialog.showOpenDialogSync({
			filters: [
				{ name: 'Keyboard Mapping', extensions: ['kbm'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile'],
			defaultPath: defaultValue,
			title: "Choose Scala Keyboard Mapping (KBM) path"
		});
		//console.log("in folderDialog: " + folder);
		if (folder != undefined && ele != undefined && ele != null) {
			ele.value = folder;
		}
		return folder;
	},

	sclFileDialog: function (ele, defaultValue) {
		folder = dialog.showOpenDialogSync({
			filters: [
				{ name: 'Scale', extensions: ['scl'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile'],
			defaultPath: defaultValue,
			title: "Choose Scala Scale (SCL) path"
		});
		//console.log("in folderDialog: " + folder);
		if (folder != undefined && ele != undefined && ele != null) {
			ele.value = folder;
		}
		return folder;
	},

	showFrequencyTable: function () {
		console.log("showFrequencyTable")
		let message = {
			"name": "getTuningFrequencies",
			"payload": tuning.createParamPayload()
		}
		astilectron.sendMessage(message, function (message) {
			// Check error
			if (message.name === "error") {
				console.log("error: ", message)
				index.errorNotification(message.payload);
			}
			var result = message.payload
			tuning.buildTuningDisplay(result)
		});
	},

	buildScaleTable: function(tones) {
		var tableEle = $('<table class="valTable"/>')
		var tableHeadEle = $("<thead/>")
		tableEle.append(tableHeadEle)
		var rowEle = $(tableHeadEle[0].insertRow(-1));
		var cell;
		cell = $('<th>#</th>')
		rowEle.append(cell);
		cell = $('<th class="val">Definition</th>')
		rowEle.append(cell);
		cell = $('<th class="val">Cents</th>')
		rowEle.append(cell);
		cell = $('<th class="val">Ratio</th>')
		rowEle.append(cell);
		var tableBodyEle = $("<tbody/>")
		tableEle.append(tableBodyEle)
		for (row = -1; row < tones.length; row++) {
			rowEle = $(tableBodyEle[0].insertRow(-1));
			if (row < 0) {
				// implicit "root" position of the scale
				cell  = $('<td>1</td>')
				rowEle.append(cell);
				cell  = $('<td class="val">1/1</td>')
				rowEle.append(cell);
				cell  = $('<td class="val">0</td>')
				rowEle.append(cell);
				cell  = $('<td class="val">1</td>')
				rowEle.append(cell);
			} else {
				cell = $(`<td>${row+2}</td>`)
				rowEle.append(cell);
				cell = $(`<td class="val">${tones[row].StringRep}</td>`);
				rowEle.append(cell);
				cell = $(`<td class="val">${tones[row].Cents.toFixed(3)}</td>`);
				rowEle.append(cell);
				cell = $(`<td class="val">${tones[row].FloatValue.toFixed(3)}</td>`);
				rowEle.append(cell);
			}
		}
		var scaleTableDiv = $('#scaleTableDiv');
		scaleTableDiv.html("");
		scaleTableDiv.append(tableEle);
	},

	buildFrequencyTable: function(freqs, scalePos) {
		// Table is 4 columns of 32 rows

		var tableEle = $('<table class="valTable"/>')
		var tableHeadEle = $("<thead/>")
		tableEle.append(tableHeadEle)
		var rowEle = $(tableHeadEle[0].insertRow(-1));
		for (var col = 0; col < 4; col++) {
			var cell;
			if (col > 0) {
				//cell  = $('<td rowSpan="32" style="padding:10px;"/>')
				//rowEle.append(cell);
				cell  = $('<td rowSpan="32" style="border-left:1px solid #666;padding:10px;"/>');
				rowEle.append(cell);
				//cell  = $('<td rowSpan="32" style="border-left:1px solid #666;padding:10px;"/>');
				//rowEle.append(cell);
			}
			cell = $('<th>Key</th>')
			rowEle.append(cell);
			cell = $('<th class="val">Degree</th>')
			rowEle.append(cell);
			cell = $('<th class="val">Hz</th>')
			rowEle.append(cell);
		}
		var tableBodyEle = $("<tbody/>")
		tableEle.append(tableBodyEle)

		midiName = function(i) {
			var names = ["C","C#","D","D#","E","F","F#","G","G#","A","A#","B"];
			var octave = Math.trunc(i / 12) - 1;
			return names[i%12] + octave
		}
		for (row = 0; row < 32; row++) {
			rowEle = $(tableBodyEle[0].insertRow(-1));

			for (col = 0; col < 4; col++) {
				note = col * 32 + row;
				freq = freqs[note].toFixed(3);

				if (col > 0 && row == 0) {
					//cell  = $('<td rowSpan="32" style="padding:10px;"/>')
					//rowEle.append(cell);
					cell  = $('<td rowSpan="32" style="border-left:1px solid #666;padding:10px;"/>');
					rowEle.append(cell);
					//cell  = $('<td rowSpan="32" style="border-left:1px solid #666;padding:10px;"/>');
					//rowEle.append(cell);
				}
				cell = $('<td/>')
				cell.html(note + "&nbsp;("+midiName(note)+")");
				rowEle.append(cell);
				cell = $('<td class="val"/>');
				cell.html(scalePos[note]+1);
				rowEle.append(cell);
				cell = $('<td class="val"/>');
				cell.html(freq);
				rowEle.append(cell);
			}
		}
		var freqTableDiv = $('#freqTableDiv');
		freqTableDiv.html("");
		freqTableDiv.append(tableEle);
	},

	toggle: function () {
		console.log("toggle")
		var useStandardChecked = document.getElementById("useStandard").checked;

		if (useStandardChecked) {
			document.getElementById("useStandardKBM").checked = true;
		}
		var useStandardKBMChecked = document.getElementById("useStandardKBM").checked;

		document.getElementById("sclPath").disabled = useStandardChecked;
		document.getElementById("useStandardKBM").disabled = useStandardChecked;

		document.getElementById("kbmPath").disabled = useStandardChecked ||  useStandardKBMChecked;

		document.getElementById("middleNote").disabled = (!useStandardChecked) && (!useStandardKBMChecked);

		document.getElementById("referenceNote").disabled = (!useStandardChecked) && (!useStandardKBMChecked);
		document.getElementById("referenceFreq").disabled = (!useStandardChecked) &&  (!useStandardKBMChecked);

	},

	createParamPayload: function () {
		return {
			"UseStandardTuning"         : document.getElementById("useStandard").checked,
			"UseStandardKeyboardMapping": document.getElementById("useStandardKBM").checked,
			"SCLPath"                   : document.getElementById("sclPath").value,
			"KBMPath"                   : document.getElementById("kbmPath").value,
			"MiddleNote"                : parseInt(document.getElementById("middleNote").value, 10),
			"ReferenceNote"             : parseInt(document.getElementById("referenceNote").value, 10),
			"ReferenceFrequency"        : parseFloat(document.getElementById("referenceFreq").value)
		}
	},

	cancelAndClose: function () {
		let message = {
			"name": "cancelTunings",
			payload: ""
		};
		// Send message
		astilectron.sendMessage(message, function (message) {
			// Check error
			if (message.name === "error") {
				index.errorNotification(message.payload);
			}
		});
	},

	buildTuningDisplay: function(result) {
		console.log("buildTuningDisplay",result)

		tuning.buildScaleTable(result.Tones);
		tuning.buildFrequencyTable(result.Frequencies, result.ScalePos);
		console.log("tuning table",document.getElementById("tuningTable"))
		document.getElementById("tuningTable").hidden = false;
	},

	sendToSynergy: function () {
		let message = {
			"name": "sendTuningToSynergy",
			"payload": tuning.createParamPayload()
		};
		//console.log("saveAndClose: " + message);
		// Send message
		astilectron.sendMessage(message, function (message) {
			// Check error
			if (message.name === "error") {
				index.errorNotification(message.payload);
			} else {
				var result = message.payload
				tuning.buildTuningDisplay(result)
				index.infoNotification("Success!");
			}
		});
	},

	listen: function () {
		console.log("tuning listening...")
		astilectron.onMessage(function (message) {
			console.log("unexpected msg: " + JSON.stringify(message));
		});
	}
};
