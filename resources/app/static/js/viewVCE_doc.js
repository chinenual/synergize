let viewVCE_doc = {

	init: function () {
		if (vce.Extra.Doc != null) {
			$('#doctext').html(vce.Extra.Doc)
		} else {
			$('#doctext').html('')
		}
	}
};
