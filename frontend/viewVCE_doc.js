import { viewVCE } from './viewVCE';

import $ from 'jquery';
Object.assign(window, { $: $, jQuery: $ });

export let viewVCE_doc = {

	init: function () {
		if (viewVCE.vce.Extra.Doc != null) {
			document.querySelector('#doctext').innerHTML = viewVCE.vce.Extra.Doc;
		} else {
			document.querySelector('#doctext').innerHTML = '';
		}
	}
};
