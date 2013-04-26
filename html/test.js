$(document).ready(function() {          // when document ready...

//var tileshape = '<svg xmlns="http://www.w3.org/2000/svg" version="1.1"> <polygon points="0,31 31,0 63,31 31,63" style="fill:lime;stroke:purple;stroke-width:1"/> </svg>'

//var tileshape = '<svg xmlns="http://www.w3.org/2000/svg" version="1.1"> <polygon points="0,19 23,0 47,19 23,39" style="fill:lime;stroke:purple;stroke-width:1"/> </svg>'
var tileshape = '<svg xmlns="http://www.w3.org/2000/svg" version="1.1"> <polygon points="0,19 23,0 47,19 23,39"/> </svg>'

//var tileshape = '<div class="isotile"><svg xmlns="http://www.w3.org/2000/svg" version="1.1"> <polygon points="0,31 31,0 63,31 31,63" style="fill:lime;stroke:purple;stroke-width:1"/> </svg></div>'

var jsonmap = $.getJSON( "map.json", function() {
  console.log( "success" );
  })
  .done(function( data ) {

      var mapwidth = data[0].length * 48
      var rowheight = 20
      var table='<div class="map">';
      /* loop over each object in the array to create rows*/
      $.each( data, function( index, item){
            /* add to html string started above*/
        table+='<div class="maprow" width = "' + mapwidth.toString() + '">';
        if (index % 2 == 1) {table+='<div></div>'}
        $.each( item, function( index, subitem){
          // table+='<td></td><td></td>';
          //table+='<td continent="' + subitem.continent+'">' + subitem.terrain+'</td><td></td>';
          //table+='<td continent="' + subitem.continent+'"></td><td></td>';
          table+='<div class="tile" continent="' + subitem.continent+'"></div><div></div>';
        });
        if (index % 2 == 0) {table+='<div></div>'}
        table+='</div>';       
      });
      table+='</div>';
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
//  $("td[continent]").html('&diams;');
//  $("td[continent]").css("font-size", "2em");

  // hard code land continents green and water continent blue
  $("td[continent]").css("background-color", "green");
  $('td[continent="6"]').css("background-color", "blue");

  //$("td[continent]").append(tileshape);
  $("div.tile").append(tileshape);
  })
  .fail(function() { console.log( "error" ); })
  .always(function() { console.log( "complete" ); });

   // perform other work here ...

    // Set another completion function for the request above
    jsonmap.complete(function() { console.log( "second complete" ); });


});
