$(document).ready(function() {          // when document ready...

// Starting with getJSON example yoinked from http://api.jquery.com/jQuery.getJSON/

// Assign handlers immediately after making the request,
// and remember the jsonmap object for this request
var jsonmap = $.getJSON( "map.json", function() {
  console.log( "success" );
  })
  .done(function( data ) {
     $("#map").html("<table>") 
//     $.each( data, function( i, data ) {
     $.each( data, function( ) {
    $("#map").append(this.length);
       $("#map table").append("<tr/>");
       $.each( this, function( ) {
         $("#map table tr:last").append("<td>" + this.terrain + "</td>");
//         $("#map table tr td").append(this.continent);
//         $( "<div>" + this.continent + "</div>" ).appendTo( "#map" );
//       });
//     });
//    for (var i=0, len=data.length; i < len; i++) {
//        console.log(data[i]);
      })
     })
    console.log( "second success" );
  })
  .fail(function() { console.log( "error" ); })
  .always(function() { console.log( "complete" ); });
   
   // perform other work here ...
    
    // Set another completion function for the request above
    jsonmap.complete(function() { console.log( "second complete" ); });

});
