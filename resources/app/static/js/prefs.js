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
                    index.errorNotification(message.payload);
                    return
		}
		var os = message.payload.Os
		preferences = message.payload.Preferences
		
		console.log("loaded preferences: " + os + ": " + preferences)

		document.getElementById("serialPort").value = preferences.SerialPort;
		document.getElementById("serialBaud").value = preferences.SerialBaud;
		document.getElementById("libraryPath").value = preferences.LibraryPath;
		if (os === "darwin") {
		    // add an onclick handler to popup a file dialog
		    var ele = document.getElementById("serialPort");
		    ele.onclick = function() {
			prefs.serialPortDialog(this,this.value);
		    }
		} else {
		    // on windows, use a straight text box 
		}
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
	}
	return folder;
    },

    cancelAndClose: function() {
	let message = {"name" : "cancelPreferences",
		       payload: ""};
	    // Send message
            astilectron.sendMessage(message, function(message) {		
		// Check error
		if (message.name === "error") {
                    index.errorNotification(message.payload);
		}
	    });	
    },
    
    saveAndClose: function() {
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
                    index.errorNotification(message.payload);
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