let index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();
	
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {

	    /*
	    json = {"name" : "getCWD"};
	    astilectron.sendMessage(json, function(message) {
		console.log("getCwd returned: " + message);
		document.getElementById("cwd").innerHTML = message.payload;
	    });
	    */
	    
            // Listen
            index.listen();
            // Explore default path
            index.explore();
        })

    },
    loadSYN(name, path) {
	if (confirm("Load Synergy state file " + path)) {
            let message = {"name": "loadSYN",
			   "payload": path};
            // Send message
            asticode.loader.show();
            astilectron.sendMessage(message, function(message) {
		// Init
		asticode.loader.hide();
		
		// Check error
		if (message.name === "error") {
                    asticode.notifier.error(message.payload);
                    return
		} else {
		    asticode.notifier.info("Successfully loaded " + name + " to Synergy")
		    return
		}
	    });
	}
    },
    loadCRT(name, path) {
	if (confirm("Load Voice Cartridge file " + path)) {
            let message = {"name": "loadCRT",
			   "payload": path};
            // Send message
            asticode.loader.show();
            astilectron.sendMessage(message, function(message) {
		// Init
		asticode.loader.hide();
		
		// Check error
		if (message.name === "error") {
                    asticode.notifier.error(message.payload);
                    return
		} else {
		    asticode.notifier.info("Successfully loaded " + name + " to Synergy")
		    return
		}
	    });
	}
    },
    loadVCE(name, path) {	
        let message = {"name": "loadVCE",
		       "payload": path};
        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
	    // Init
	    asticode.loader.hide();

	    // Check error
	    if (message.name === "error") {
                asticode.notifier.error(message.payload);
                return
	    }
	    vce = message.payload;
	    index.load("view.html", document.getElementById("content"));
	    view.refreshText();
	});
    },
    addFolder(name, path) {
        let div = document.createElement("div");
        div.className = "dir";
        div.onclick = function() { index.explore(path) };
	if (name == "..") name = "&lt;Parent&gt;";
        div.innerHTML = `<i class="fa fa-folder"></i><span>` + name + `</span>`;
        document.getElementById("dirs").appendChild(div)
    },
    addSYNFile(name, path) {
        let div = document.createElement("div");
        div.className = "file";
        div.onclick = function() { index.loadSYN(name,path) };
        div.innerHTML = `<i class="fa fa-file"></i><span>` + name + `</span>`;
        document.getElementById("SYNfiles").appendChild(div)
    },
    addCRTFile(name, path) {
        let div = document.createElement("div");
        div.className = "file";
        div.onclick = function() { index.loadCRT(name,path) };
        div.innerHTML = `<i class="fa fa-file"></i><span>` + name + `</span>`;
        document.getElementById("CRTfiles").appendChild(div)
    },
    addVCEFile(name, path) {
        let div = document.createElement("div");
        div.className = "file";
        div.onclick = function() { index.loadVCE(name,path) };
        div.innerHTML = `<i class="fa fa-file"></i><span>` + name + `</span>`;
        document.getElementById("VCEfiles").appendChild(div)
    },
    explore: function(path) {
        // Create message
        let message = {"name": "explore"};
        if (typeof path !== "undefined") {
            message.payload = path
        }

        // Send message
        asticode.loader.show();
        astilectron.sendMessage(message, function(message) {
            // Init
            asticode.loader.hide();

            // Check error
            if (message.name === "error") {
                asticode.notifier.error(message.payload);
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
		div.innerHTML = "<b>Cartridge Files (.CRT)</b>";
		document.getElementById("CRTfiles").appendChild(div);

		for (let i = 0; i < message.payload.CRTfiles.length; i++) {
                    index.addCRTFile(message.payload.CRTfiles[i].name, message.payload.CRTfiles[i].path);
		}
	    }

	    document.getElementById("SYNfiles").innerHTML = ""
	    if (message.payload.SYNfiles.length > 0) {
		let div = document.createElement("div")
		div.innerHTML = "<b>Synergy State (.SYN)</b>";
		document.getElementById("SYNfiles").appendChild(div);
		for (let i = 0; i < message.payload.SYNfiles.length; i++) {
                    index.addSYNFile(message.payload.SYNfiles[i].name, message.payload.SYNfiles[i].path);
		}
	    }
	    
	    document.getElementById("VCEfiles").innerHTML = ""
	    if (message.payload.VCEfiles.length > 0) {
		let div = document.createElement("div")
		div.innerHTML = "<b>Voice Files (.VCE)</b>";
		document.getElementById("VCEfiles").appendChild(div);
		for (let i = 0; i < message.payload.VCEfiles.length; i++) {
                    index.addVCEFile(message.payload.VCEfiles[i].name, message.payload.VCEfiles[i].path);
		}
	    }
	    /**
            if (typeof message.payload.files !== "undefined") {
                document.getElementById("files_panel").style.display = "block";
                let canvas = document.createElement("canvas");
                document.getElementById("files").append(canvas);
                new Chart(canvas, message.payload.files);
            } else {
                document.getElementById("files_panel").style.display = "none";
            }
	    **/
        })
    },
    updateConnectionStatus: function (status) {
	console.log("update status: " + status);
	document.getElementById("firmwareVersion").innerHTML = status;
    },
    fileDialog: function () {
	const {dialog} = require('electron').remote;

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
    runCOMTST: function() {
	document.getElementById("testProgress").innerHTML = "Running...";
	let message = {"name" : "runCOMTST"};
	astilectron.sendMessage(message, function(message) {
	    console.log("runCOMTST returned: " + message);
	    document.getElementById("testProgress").innerHTML = message.payload;
	});
    },
    load: function(url, element) {
	console.log("load " + url + " into " + element);
	req = new XMLHttpRequest();

	req.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
		element.innerHTML = req.responseText;
            }
	};

	req.open("GET", url, false);
	req.send(null); 
    },
    listen: function() {
	console.log("index listening...")
        astilectron.onMessage(function(message) {
            switch (message.name) {
	    case "updateConnectionStatus":
		index.updateConnectionStatus(message.payload);
		return {payload: "ok"};
            case "fileDialog":
                f = index.fileDialog(message.payload);
                return {payload: f};
                break;
            case "viewVCE":
		console.log("viewVCE: " + message.payload);
		vce = message.payload;
		index.load("view.html", document.getElementById("content"));
		view.refreshText();
		
                return {payload: "ok"};
                break;
            case "runDiag":
		console.log("runDiag: " + message.payload);
		vce = message.payload;
		index.load("diag.html", document.getElementById("content"));
		
                return {payload: "ok"};
                break;
            }
        });
    }
};
