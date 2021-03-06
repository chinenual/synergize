let prefs = {
        init: function () {
            // Wait for astilectron to be ready
            document.addEventListener('astilectron-ready', function () {

                let message = {
                    "name": "getPreferences",
                    "payload": ""
                };
                var preferences;

                // Send message
                astilectron.sendMessage(message, function (message) {
                    // Check error
                    if (message.name === "error") {
                        index.errorNotification(message.payload);
                        return
                    }
                    var os = message.payload.Os
                    preferences = message.payload.Preferences

                    //console.log("loaded preferences: " + os + ": " + JSON.stringify(preferences))

                    document.getElementById("useSerial").checked = preferences.UseSerial ? "checked" : "";
                    document.getElementById("serialPort").value = preferences.SerialPort;
                    document.getElementById("serialBaud").value = preferences.SerialBaud;
                    document.getElementById("flowControl").checked = preferences.SerialFlowControl;

                    document.getElementById("libraryPath").value = preferences.LibraryPath;

                    document.getElementById("useOsc").checked = preferences.UseOsc ? "checked" : "";
                    document.getElementById("oscAutoConfig").checked = preferences.OscAutoConfig ? "checked" : "";
                    document.getElementById("oscPort").value = preferences.OscPort;
                    document.getElementById("oscCSurfaceAddress").value = preferences.OscCSurfaceAddress;
                    document.getElementById("oscCSurfacePort").value = preferences.OscCSurfacePort;
                    prefs.toggleOsc();

                    if (os === "darwin") {
                        /*
                          * as nice as this would be, macos hides the /dev directory
                          * from the UI dialogs.  Best to just let folk type it in...
                          *

                        // add an onclick handler to popup a file dialog
                        var ele = document.getElementById("serialPort");
                        ele.onclick = function() {
                        prefs.serialPortDialog(this,this.value);
                        }

                        *
                        */
                    } else {
                        // on windows, use a straight text box
                    }
                });

                // Listen
                prefs.listen();
            })

        },

        toggleOsc: function () {
            var useSerialChecked = document.getElementById("useSerial").checked;
            var useOscChecked = document.getElementById("useOsc").checked;
            var autoChecked = document.getElementById("oscAutoConfig").checked;


            document.getElementById("serialPort").disabled = (!useSerialChecked);
            document.getElementById("serialBaud").disabled = (!useSerialChecked);
            document.getElementById("flowControl").disabled = (!useSerialChecked);

            document.getElementById("oscPort").disabled = (!useOscChecked);
            document.getElementById("oscAutoConfig").disabled = (!useOscChecked);
            document.getElementById("oscCSurfaceAddress").disabled = (!useOscChecked) || autoChecked;
            document.getElementById("oscCSurfacePort").disabled = (!useOscChecked) || autoChecked;
        },

        serialPortDialog: function (ele, defaultValue) {
                //console.log("in fileDialog default: " + defaultValue);

                file = dialog.showOpenDialogSync({
                    properties: ['openFile'],
                    title: "Select Serial Port",
                    buttonLabel: "Select",
                    defaultPath: defaultValue
                });
                //console.log("in fileDialog: " + file);
                if (file != undefined && ele != undefined && ele != null) {
                    ele.value = file;
                }
                return file;
            }

        ,
        libraryFolderDialog: function (ele, defaultValue) {
            folder = dialog.showOpenDialogSync({
                properties: ['openDirectory'],
                defaultPath: defaultValue,
                title: "Choose Voice Library path"
            });
            //console.log("in folderDialog: " + folder);
            if (folder != undefined && ele != undefined && ele != null) {
                ele.value = folder;
            }
            return folder;
        }
        ,

        cancelAndClose: function () {
            let message = {
                "name": "cancelPreferences",
                payload: ""
            };
            // Send message
            astilectron.sendMessage(message, function (message) {
                // Check error
                if (message.name === "error") {
                    index.errorNotification(message.payload);
                }
            });
        }
        ,

        saveAndClose: function () {
            let message = {
                "name": "savePreferences",
                "payload": {
                    "UseSerial": document.getElementById("useSerial").checked,
                    "SerialPort": document.getElementById("serialPort").value,
                    "SerialBaud": parseInt(document.getElementById("serialBaud").value, 10),
                    "SerialFlowControl": document.getElementById("flowControl").checked,
                    "LibraryPath": document.getElementById("libraryPath").value,
                    "UseOsc": document.getElementById("useOsc").checked,
                    "OscAutoConfig": document.getElementById("oscAutoConfig").checked,
                    "OscPort": parseInt(document.getElementById("oscPort").value, 10),
                    "OscCSurfaceAddress": document.getElementById("oscCSurfaceAddress").value,
                    "OscCSurfacePort": parseInt(document.getElementById("oscCSurfacePort").value, 10),
                }
            };
            //console.log("saveAndClose: " + message);
            // Send message
            astilectron.sendMessage(message, function (message) {
                // Check error
                if (message.name === "error") {
                    index.errorNotification(message.payload);
                }
            });
        }
        ,
        listen: function () {
            console.log("prefs listening...")
            astilectron.onMessage(function (message) {
                console.log("unexpected msg: " + JSON.stringify(message));
            });
        }
    }
;
