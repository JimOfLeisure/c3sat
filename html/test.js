$(document).ready(function() {          // when document ready...

var jsonmap = $.getJSON( "map.json", function() {
  console.log( "success" );
  })
  .done(function( data ) {

      var table='<table>';
      /* loop over each object in the array to create rows*/
      $.each( data, function( index, item){
            /* add to html string started above*/
        table+='<tr>';
        if (index % 2 == 1) {table+='<td>.</td>'}
        $.each( item, function( index, subitem){
          table+='<td continent="' + subitem.continent+'">' + subitem.terrain+'</td><td></td>';
        });
        if (index % 2 == 0) {table+='<td></td>'}
        table+='</tr>';       
      });
      table+='</table>';
/* insert the html string*/
      $("#map").html( table );      

    console.log( "second success" );

// color all non-pad tiles green
  $("td[continent]").css("background-color", "green");

// coloring tiles based on terrain
  $("td:contains('0')").css("color", "white");

  $("td:contains('0')").css("background-color", "cornsilk");
  $("td:contains('221')").css("background-color", "darkblue");
  $("td:contains('204')").css("background-color", "blue");
  $("td:contains('187')").css("background-color", "lightblue");
  $("td:contains('51')").css("background-color", "white");
  $("td:contains('17')").css("background-color", "sandybrown");
  $("td:contains('98')").css("background-color", "chocolate");

  $("td:contains('0')").css("color", "cornsilk");
  $("td:contains('221')").css("color", "darkblue");
  $("td:contains('204')").css("color", "blue");
  $("td:contains('187')").css("color", "lightblue");
  $("td:contains('51')").css("color", "white");
  $("td:contains('17')").css("color", "sandybrown");
  $("td:contains('98')").css("color", "chocolate");

// replace text with unicode
  $("td[continent]").html('&diams;');
  $("td[continent]").css("font-size", "2em");

  })
  .fail(function() { console.log( "error" ); })
  .always(function() { console.log( "complete" ); });

   // perform other work here ...

    // Set another completion function for the request above
    jsonmap.complete(function() { console.log( "second complete" ); });


});
