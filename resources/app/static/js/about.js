let about = {
    setVersion: function(versionString) {
	console.log("top of about " + versionString);
	document.getElementById("version").innerHTML = versionString;
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();
	
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            about.listen();
        })
    },
    listen: function() {
	console.log("index listening...")
        astilectron.onMessage(function(message) {
            switch (message.name) {
            case "setVersion":
		
                about.setVersion(message.payload);
                return {payload: "payload"};
                break;
            }
        });
    }
};
