(function process(/*RESTAPIRequest*/ request, /*RESTAPIResponse*/ response) {

    var chunk = request.body.data.short_key;
    var stURLEntry = new GlideRecord("u_url_shortener");
    stURLEntry.addQuery("u_chunk", chunk);
    stURLEntry.query();
    if (stURLEntry.next()) {
        response.setContentType('application/json');
        response.setStatus(200);
        var writer = response.getStreamWriter();
        var res = {
            "originalURL": stURLEntry.u_url.toString()
        };
        writer.writeString(JSON.stringify(res));
    }

})(request, response);