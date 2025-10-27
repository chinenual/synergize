let shell = require('electron').shell

let about = {
    init: function () {
        // make sure external web links open in system browser - not the application:
        document.addEventListener('click', function (event) {
            if (event.target.tagName === 'A' && event.target.href.startsWith('http')) {
                event.preventDefault()
                shell.openExternal(event.target.href)
            }
        })
        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function () {

            let json = { "name": "getVersion" };
            astilectron.sendMessage(json, function (message) {
                console.log("getVersion returned: " + JSON.stringify(message));
                document.getElementById("version").innerHTML = message.payload.Version;
                if (message.payload.NewVersionAvailable) {
                    document.getElementById("updateAvailable").style.display = "block";
                }
            });

            // Listen
            about.listen();
        })
    },

    listen: function () {
        console.log("index listening...")
        astilectron.onMessage(function (message) {
            console.log("got message: " + JSON.stringify(message));
            //            switch (message.name) {
            //            }
        });
    }
};
