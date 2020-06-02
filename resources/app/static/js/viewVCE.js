let vce = {};
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
    init: function () {
        Chart.defaults.global.defaultFontColor = 'white';
        Chart.defaults.global.defaultFontSize = 14;

        viewVCE_voice.init();
        viewVCE_keyprop.init();
        viewVCE_keyeq.init();
        viewVCE_envs.init();
        viewVCE_filters.init();

        viewVCE_voice.voicingModeVisuals();
    }
}

