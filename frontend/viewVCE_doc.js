import { $ } from 'jquery';
import { viewVCE } from './viewVCE';

export let viewVCE_doc = {

	init: function () {
		if (viewVCE.vce.Extra.Doc != null) {
			$('#doctext').html(viewVCE.vce.Extra.Doc)
		} else {
			$('#doctext').html('')
		}
	}
};
