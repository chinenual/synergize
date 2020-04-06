let index = {
    about: function(html) {
	window.open("about.html","About Synergize","width=400,height=300,scrollbars=0,resizable=0");
    },
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
        astilectron.onMessage(function(message) {
            switch (message.name) {
            case "about":
		
                index.about(message.payload);
                return {payload: "payload"};
                break;
            }
        });
    }
};
