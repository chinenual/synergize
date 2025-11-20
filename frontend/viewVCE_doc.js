import { viewVCE } from './viewVCE';

export let viewVCE_doc = {

	init: function () {
		if (viewVCE.vce.Extra.Doc != null) {
			document.querySelector('#doctext').innerHTML = viewVCE.vce.Extra.Doc;
		} else {
			document.querySelector('#doctext').innerHTML = '';
		}
	}
};
