import {UIService} from '/bindings/github.com/chinenual/synergize';
import * as wails from '@wailsio/runtime';

export let about = {
  init: function() {
    // make sure external web links open in system browser - not the
    // application:
    document.addEventListener('click', function(event) {
      if (event.target.tagName === 'A' &&
          event.target.href.startsWith('http')) {
        event.preventDefault()
        wails.Browser.OpenURL(event.target.href);
      }
    })
    let val = ['', false];
    try {
      val = UIService.GetVersion()
    } catch (exc) {
        // ignore
        console.log("exc",exc);
    }
    let version = val[0];
    let update = val[1];
    console.log('getVersion returned: ' + version);
    document.getElementById('version').innerHTML = version;
    if (update) {
      document.getElementById('updateAvailable').style.display = 'block';
    }
  }
};
