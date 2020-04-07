let index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();
	
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();
	    
            // Explore default path
	    //            index.explore();
	    
	    
	    
        })
    },
    listen: function() {
	console.log("index listening...")
        astilectron.onMessage(function(message) {
            switch (message.name) {
            case "about":
//		
//                index.about(message.payload);
//                return {payload: "payload"};
                break;
            }
        });
    }
};
