let vce = null;
let vceFilterNames = null;

// https://www.color-hex.com/color-palette/89750
//let chartColors=[
//    'rgb(225,215,0)',
//    'rgb(79,156,244)',
//    'rgb(62,244,6)',
//    'rgb(5,82,244)',
//    'rgb(4,169,24)'
//];

// XREF: match yellow - see base.css
// https://htmlcolors.com/palette/26/google-palette
let chartColors = [
    'rgb(244,180,0)', // golden
    'rgb(66,133,244)', // blue
    'rgb(219,68,55)', // redish
    'rgb(15,157,88)', // green
    'rgb(255,255,255)' // white
];

let viewVCE = {
    // flag to prevent programatic voice changes from triggering onchange updates to the Synergy
    supressOnchange: false,

    init: function () {
        // no onchange events while we update input and text for the new voice
        viewVCE.supressOnchange = true;
        //console.log('--- start viewVCE init');
        Chart.defaults.global.defaultFontColor = 'white';
        Chart.defaults.global.defaultFontSize = 14;

        viewVCE_keyprop.init(false);
        viewVCE_keyeq.init(false);
        viewVCE_envs.init(false);
        viewVCE_filters.init(false);
        viewVCE_doc.init();
        viewVCE_voice.init(false); // do this last since the uitest uses existance of VNAME to indicate that the pages are loaded and ready to roll

        viewVCE_voice.voicingModeVisuals();

        // back to normal:
        viewVCE.supressOnchange = false;
        viewVCE_voice.sendToCSurface(null, "voice-tab", 1);
        //console.log('--- finish viewVCE init');
    }
}

