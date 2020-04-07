let view = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();
	
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            view.listen();
	    
        })
    },
    listen: function() {
	console.log("view listening...")
        astilectron.onMessage(function(message) {
            switch (message.name) {
            case "foobar":
                break;
            }
        });
    }
};
