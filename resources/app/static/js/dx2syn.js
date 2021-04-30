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
		document.getElementById("subprocessCancelButton").onclick = function() {
			console.log("SAW CANCEL");
		};
		$('#subprocessModal').modal({
			backdrop: "static" // clicking outside the dialog doesnt close the dialog
		});
		for (i=0;i<1000;i++) {
			dx2syn.addDxSubprocessLog(`msg${i}\n`);

		}
	},

	addDxSubprocessLog: function(msgs) {
		html = msgs.replaceAll('\n','<br/>\n');
		document.getElementById("logOutput").innerHTML =
			document.getElementById("logOutput").innerHTML + html;
		$('#subprocessModal').modal('handleUpdate')
	},
};
