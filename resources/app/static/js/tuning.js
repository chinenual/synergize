let tuning = {
	init: function() {
		console.log("tuning.init()")
		document.getElementById("freqTableDiv").hidden = true;
	},

	showFrequencyTable: function () {
		console.log("showFrequencyTable")

		// Table is 8 columns of 16 rows
		var tableEle = $('<table class="valTable"/>')
		var tableHeadEle = $("<thead/>")
		tableEle.append(tableHeadEle)
		var rowEle = $(tableHeadEle[0].insertRow(-1));
		for (var col = 0; col < 8; col++) {
			var cell;
			if (col > 0) {
				cell  = $('<td rowSpan="16" style="padding:10px;"/>')
				rowEle.append(cell);
				cell  = $('<td rowSpan="16" style="border-left:1px solid #666;padding:10px;"/>');
				rowEle.append(cell);
			}
			cell = $('<th>Note</th>')
			rowEle.append(cell);
			cell = $('<th class="val">Hz</th>')
			rowEle.append(cell);
		}
		var tableBodyEle = $("<tbody/>")
		tableEle.append(tableBodyEle)
		for (row = 0; row < 16; row++) {
			rowEle = $(tableBodyEle[0].insertRow(-1));

			for (col = 0; col < 8; col++) {
				note = col * 16 + row;
				freq = (note*12.0).toFixed(1) // fake

				if (col > 0 && row == 0) {
					cell  = $('<td rowSpan="16" style="padding:10px;"/>')
					rowEle.append(cell);
					cell  = $('<td rowSpan="16" style="border-left:1px solid #666;padding:10px;"/>');
					rowEle.append(cell);
				}
				cell = $('<td/>')
				cell.html(note);
				rowEle.append(cell);
				cell = $('<td class="val"/>');
				cell.html(freq);
				rowEle.append(cell);
			}
		}
		var freqTableDiv = $('#freqTableDiv');
		freqTableDiv.html("");
		freqTableDiv.append(tableEle);
		document.getElementById("freqTableDiv").hidden = false;
	},

	toggle: function () {
		console.log("toggle")
		var useStandardChecked = document.getElementById("useStandard").checked;

		if (useStandardChecked) {
			document.getElementById("useStandardKBM").checked = true;
		}
		var useStandardKBMChecked = document.getElementById("useStandardKBM").checked;

		document.getElementById("freqTableDiv").hidden = true;

		document.getElementById("sclPath").disabled = useStandardChecked;
		document.getElementById("useStandardKBM").disabled = useStandardChecked;

		document.getElementById("kbmPath").disabled = useStandardChecked ||  useStandardKBMChecked;

		document.getElementById("referenceNote").disabled = (!useStandardChecked) && (!useStandardKBMChecked);
		document.getElementById("referenceFreq").disabled = (!useStandardChecked) &&  (!useStandardKBMChecked);

	},
};
