let prefs = {
    init: function() {
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
	    
	    let message = {"name": "getPreferences",
			   "payload": ""};
	    var preferences;
	    
	    // Send message
            astilectron.sendMessage(message, function(message) {		
		// Check error
		if (message.name === "error") {
                    asticode.notifier.error(message.payload);
                    return
		}
		preferences = message.payload
		console.log("loaded preferences: " + preferences)

		document.getElementById("serialPort").value = preferences.SerialPort;
		document.getElementById("serialBaud").value = preferences.SerialBaud;
		document.getElementById("libraryPath").value = preferences.LibraryPath;
	    });
	    
            // Listen
            prefs.listen();
        })

    },
    serialPortDialog: function (ele, defaultValue) {
	const {dialog} = require('electron').remote;
	console.log("in fileDialog default: " + defaultValue);

	file = dialog.showOpenDialogSync({
            properties: ['openFile'],
	    title: "Select Serial Port",
	    buttonLabel: "Select",
	    defaultPath: defaultValue
	});
	console.log("in fileDialog: " + file);
	if (file != undefined && ele != undefined && ele != null) {
	    ele.value = file;
	    /* HACK: onchange and oninput dont fire because we use onclick
	     * this works, but would be "cleaner" to have something like
	     * onchange monitoring the changes. 
	     */
	    prefs.save();
	}
	return file;
    },
    libraryFolderDialog: function (ele, defaultValue) {
	const {dialog} = require('electron').remote;

	folder = dialog.showOpenDialogSync({
            properties: ['openDirectory'],
	    defaultPath: defaultValue,
	    title: "Choose Voice Library path"
	});
	console.log("in folderDialog: " + folder);
	if (folder != undefined && ele != undefined && ele != null) {
	    ele.value = folder;
	    /* HACK: onchange and oninput dont fire because we use onclick
	     * this works, but would be "cleaner" to have something like
	     * onchange monitoring the changes. 
	     */
	    prefs.save();
	}
	return folder;
    },

    save: function() {
	let message = {"name" : "savePreferences",
		       "payload": {
			   "SerialPort" : document.getElementById("serialPort").value,
			   "SerialBaud" : parseInt(document.getElementById("serialBaud").value,10),
			   "LibraryPath" : document.getElementById("libraryPath").value
		       }};
	    // Send message
            astilectron.sendMessage(message, function(message) {		
		// Check error
		if (message.name === "error") {
                    asticode.notifier.error(message.payload);
                    return
		}
	    });
    },
    listen: function() {
	console.log("prefs listening...")
        astilectron.onMessage(function(message) {
	    console.log("unexpected msg: " + message);
        });
    }
};
