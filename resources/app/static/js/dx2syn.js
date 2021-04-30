let dx2syn = {
	convertFileDialog: function() {
		path = dialog.showOpenDialogSync({
			title: "Choose DX7 Sysex",
			filters: [
				{ name: 'DX Sysex', extensions: ['syx','sysx'] },
				{ name: 'All Files', extensions: ['*'] }],
			properties: ['openFile', 'openDirectory']
		});
		console.log("in convertFileDialog: " + path);
		if (path != undefined) {
			dx2syn.runConvert(path[0]);
		}
	},

	convertFolderDialog: function() {
		path = dialog.showOpenDialogSync({
			title: "Choose Folder containing DX7 Sysex's",
			properties: ['openDirectory']
		});
		console.log("in convertFolderDialog: " + path);
		if (path != undefined) {
			dx2syn.runConvert(path[0]);
		}
	},

	runConvert: function(path) {
		console.log("runConvert: " + path)
		document.getElementById("subprocessTitle").innerHTML = "Convert DX7 Sysex";
		document.getElementById("logOutput").innerHTML = '';
		document.getElementById("subprocessCloseButton").setAttribute("disabled", "disabled");
		document.getElementById("subprocessCancelButton").removeAttribute("disabled");

		document.getElementById("subprocessCancelButton").onclick = function() {
			console.log("SAW CANCEL");
			let message = {
				"name": "dx2synCancel",
				"payload": "DummyPayload"
			};
			// Send message
			console.log("call dx2synCancel: " + path);
			astilectron.sendMessage(message, function (message) {
				// Check error
				document.getElementById("subprocessCloseButton").removeAttribute("disabled");
				document.getElementById("subprocessCancelButton").setAttribute("disabled", "disabled");
				if (message.name === "error") {
					index.errorNotification(message.payload);
				}
			});
		};
		$('#subprocessModal').modal({
			backdrop: "static" // clicking outside the dialog doesnt close the dialog
		});
		{
			let message = {
				"name": "dx2synStart",
				"payload": {"Path": path}
			};
			// Send message
			console.log("call dx2synStart: " + path);
			astilectron.sendMessage(message, function (message) {
				// Check error
				if (message.name === "error") {
					index.errorNotification(message.payload);
				}
			});
		}
	},

	finishConvert: function(msg) {
		dx2syn.addProcessLog("\n" + msg);
		document.getElementById("subprocessCloseButton").removeAttribute("disabled");
		document.getElementById("subprocessCancelButton").setAttribute("disabled", "disabled");
	},

	addProcessLog: function(msgs) {
		html = msgs.replaceAll('\n','<br/>\n');
		if (html[html.length-1] != '\n') {
			html += "<br/>\n";
		}
		document.getElementById("logOutput").innerHTML =
			document.getElementById("logOutput").innerHTML + html;
		$('#subprocessModal').modal('handleUpdate')
	},
};
