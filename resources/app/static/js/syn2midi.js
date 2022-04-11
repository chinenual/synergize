let syn2midi = {

	init: function() {

	},

	lastTempo : 120.0,

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
		}
		return path;
	},

	convertDialog: function() {

		document.getElementById("tempo").value = syn2midi.lastTempo;
		document.getElementById("syn2midiConvertButton").onclick = function() {
			if (document.getElementById("synPath").value != undefined
				&& document.getElementById("tempo").value != undefined) {
				syn2midi.lastTempo = document.getElementById("tempo").value
				syn2midi.runConvert(document.getElementById("synPath").value,
					document.getElementById("tempo").value);
			}
		};

		$('#syn2midiModal').modal({
			backdrop: "static" // clicking outside the dialog doesnt close the dialog
		});
	},

	runConvert: function(path, tempo) {
		console.log("runConvert: " + path + " " + tempo)
		{
			let message = {
				"name": "syn2midi",
				"payload": {
					"Path": path,
					"Tempo":parseFloat(tempo)
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
