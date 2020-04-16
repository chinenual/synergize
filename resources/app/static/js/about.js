let about = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {

	    let json = {"name" : "getVersion"};
	    astilectron.sendMessage(json, function(message) {
		console.log("getVersion returned: " + message);
		document.getElementById("version").innerHTML = message.payload;
	    });
	    
            // Listen
            about.listen();
        })
    },
    listen: function() {
	console.log("index listening...")
        astilectron.onMessage(function(message) {
	    console.log("got message: " + message);
//            switch (message.name) {
//            }
        });
    }
};
