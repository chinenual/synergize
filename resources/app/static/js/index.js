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
