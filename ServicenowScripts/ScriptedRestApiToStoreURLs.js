(function process(/*RESTAPIRequest*/ request, /*RESTAPIResponse*/ response) {

    //gs.log("URLShortener Log: "+request.body.data.short_key+" - "+request.body.data.url);
    var stURLEntry = new GlideRecord("u_url_shortener");
    stURLEntry.initialize();
    stURLEntry.u_chunk = request.body.data.short_key;
    stURLEntry.u_url = request.body.data.url;
    stURLEntry.insert();

})(request, response);