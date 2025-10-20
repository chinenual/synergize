let syn2midi = {

	init: function() {
	},

	lastTempo : 120.0,
	lastRaw : false,
	lastMaxClock : 2 * 60,

	synpathDialog: function (ele) {
		path = dialog.showOpenDialogSync({
			title: 'Choose SYN file to convert to MIDI',
			filters: [
				{ name: 'SYN', extensions: ['syn'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile']
		});
		//console.log("in folderDialog: " + folder);
		if (path != undefined && ele != undefined && ele != null) {
			ele.value = path;

			syn2midi.getSynSequencerState(path[0]);
		}
		return path;
	},

	modeChange: function() {
		raw = !document.getElementById("syn2midiMode").checked;
		console.log("mode change " + raw);
		if (raw) {
			$("#syn2midiMaxClockDiv").hide();
			$("#syn2midiTrackModeDiv").hide();
		} else {
			$("#syn2midiMaxClockDiv").show();
			$("#syn2midiTrackModeDiv").show();
		}
	},

	convertDialog: function() {

		document.getElementById("syn2midiTempo").value = syn2midi.lastTempo;
		document.getElementById("syn2midiMode").checked = !syn2midi.lastRaw;
		document.getElementById("syn2midiMaxClock").value = syn2midi.lastMaxClock;
		syn2midi.modeChange();
		document.getElementById("syn2midiConvertButton").onclick = function() {
			if (document.getElementById("syn2midiPath").value != undefined
				&& document.getElementById("syn2midiTempo").value != undefined) {
				buttons = [
					parseInt(document.getElementById("syn2midiTrack1").value,10),
					parseInt(document.getElementById("syn2midiTrack2").value,10),
					parseInt(document.getElementById("syn2midiTrack3").value,10),
					parseInt(document.getElementById("syn2midiTrack4").value,10)
				]
				syn2midi.lastTempo = document.getElementById("syn2midiTempo").value
				syn2midi.runConvert(document.getElementById("syn2midiPath").value,
					parseInt(document.getElementById("syn2midiTempo").value,10),
					!document.getElementById("syn2midiMode").checked,
					parseInt(document.getElementById("syn2midiMaxClock").value,10),
					buttons);
			}
		};

		$('#syn2midiModal').modal({
			backdrop: "static" // clicking outside the dialog doesnt close the dialog
		});
	},

	getSynSequencerState: function(path) {
		console.log("getSynSequencerState: " + path )
		{
			let message = {
				"name": "getSynSequencerState",
				"payload": {
					"Path": path
				}
			};
			// Send message
			console.log("call getSynSequencerState: " + path );
			astilectron.sendMessage(message, function (message) {
				// Check error
				if (message.name === "error") {
					index.errorNotification(message.payload);
				} else {
					document.getElementById("syn2midiTrack1").value = "" + message.payload.TrackButtons[0]
					document.getElementById("syn2midiTrack2").value = "" + message.payload.TrackButtons[1]
					document.getElementById("syn2midiTrack3").value = "" + message.payload.TrackButtons[2]
					document.getElementById("syn2midiTrack4").value = "" + message.payload.TrackButtons[3]
				}
			});
		}
	},

	runConvert: function(path, tempo, raw, maxClockSeconds, buttons) {
		console.log("runConvert: " + path + " " + tempo)
		{
			let message = {
				"name": "syn2midi",
				"payload": {
					"Path": path,
					"Tempo":parseFloat(tempo),
					"Raw": raw,
					"MaxClockSeconds": maxClockSeconds,
					"TrackButtons": buttons
				}
			};
			// Send message
			console.log("call syn2midi: " + path + " " + tempo);
			astilectron.sendMessage(message, function (message) {
				// Check error
				if (message.name === "error") {
					index.errorNotification(message.payload);
				} else {
					index.infoNotification("Successfully converted Synergy sequence data to " + path + ".mid");
				}
			});
		}
	},

};
