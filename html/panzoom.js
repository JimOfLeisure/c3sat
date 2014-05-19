$(document).ready(function() {          // when document ready...

  /* js load code yoinked from http://stackoverflow.com/questions/14068031/embedding-external-svg-in-html-for-javascript-manipulation */
  /*
  xhr = new XMLHttpRequest();
  xhr.open("GET","map.svg",false);
  //xhr.open("GET","Mexican_states_by_population_2013.svg",false);
  // Following line is just to be on the safe side;
  // not needed if your server delivers SVG with correct MIME type
  xhr.overrideMimeType("image/svg+xml");
  xhr.send("");
  // add panzoom class to svg
  // not working xhr.responseXML.documentElement.className = xhr.responseXML.documentElement.className + "panzoom"
  document.getElementById("map")
    .appendChild(xhr.responseXML.documentElement);
    */

  // add panzoom class to svg
  $("svg").addClass("panzoom");

  // Series of buttons to hide various layers / elements
  $("button#baseter").click(function(){
      $(".tilebaseterrain").toggle();
  });

  $("button#ovrter").click(function(){
      $(".overlayterrain").toggle();
  });

  $("button#debuginfo").click(function(){
      $(".whatsthis").toggle();
  });

  $("button#fogofwar").click(function(){
      $(".fog").toggle();
  });

  $("button#bgrect").click(function(){
      $(".mapEdge").toggle();
  });

//$(".panzoom").panzoom();
// Don't allow map to be smaller than viewport
/*
  $('.panzoom').panzoom({
    contain: 'invert',
    minScale: 1
  });
  */

        (function() {
          var $section = $('body');
          var $panzoom = $section.find('.panzoom').panzoom({
            $zoomIn: $section.find(".zoom-in"),
            $zoomOut: $section.find(".zoom-out"),
            $zoomRange: $section.find(".zoom-range"),
            $reset: $section.find(".reset")
          });
          $panzoom.parent().on('mousewheel.focal', function( e ) {
            e.preventDefault();
            var delta = e.delta || e.originalEvent.wheelDelta;
            var zoomOut = delta ? delta < 0 : e.originalEvent.deltaY > 0;
            $panzoom.panzoom('zoom', zoomOut, {
              increment: 0.1,
              animate: false,
              focal: e
            });
          });
        })();

});
