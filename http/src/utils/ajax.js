"use strict";


function POST(uri, payload, success, failure) {
	$.ajax({
		type: "POST",
		url: uri,
		dataType: "json",
		contentType: "application/json",
		data: JSON.stringify(payload),
		success: success,
		failure: failure
	})
}

export { POST };
