$(document).ready(function() {          // when document ready...

  //$(".tile").addClass("fog");
  //$("div:contains('0x8')").addClass("fog");
  //$("div:contains('0x8')").css("color", "blue");
  //$("div.visible:contains('0x8')").css("color", "blue");
  //$("div.visible:contains('0x8')").css("background-color", "darkblue");
  //$("div.visible:contains('0x6')").css("background-color", "blue");
  //$("div.visible:contains('0x0')").css("background-color", "lightblue");

  //$("td.visible:contains('0x6')").css("background-color", "blue");
  $("td.visible:contains('0x0')").css("color", "white");

  $("td.visible:contains('0xdd')").css("background-color", "darkblue");
  $("td.visible:contains('0xcc')").css("background-color", "blue");
  $("td.visible:contains('0xbb')").css("background-color", "lightblue");
  $("td.visible:contains('0x0')").css("background-color", "cornsilk");
  $("td.visible:contains('0x33')").css("background-color", "white");
  $("td.visible:contains('0x11')").css("background-color", "sandybrown");
  $("td.visible:contains('0x62')").css("background-color", "chocolate");

  $("td.visible:contains('0xdd')").css("color", "darkblue");
  $("td.visible:contains('0xcc')").css("color", "blue");
  $("td.visible:contains('0xbb')").css("color", "lightblue");
  $("td.visible:contains('0x0')").css("color", "cornsilk");
  $("td.visible:contains('0x33')").css("color", "white");
  $("td.visible:contains('0x11')").css("color", "sandybrown");
  $("td.visible:contains('0x62')").css("color", "chocolate");

  $("td.fog").css("color", "black");

  // SVG selectors. Note that the fill affects the text fields, too. Add '"polygon"' to children to leave the text alone: ...children("polygon").css...
  // and also delete the folowing line which hides all the text fields
  $("text").css("display", "none");
  $("text:contains('0xdd')").parent().children().css("fill", "darkblue");
  $("text:contains('0xcc')").parent().children().css("fill", "mediumblue");
  $("text:contains('0xbb')").parent().children().css("fill", "blue");
  $("text:contains('0x0')").parent().children().css("fill", "cornsilk");
  $("text:contains('0x33')").parent().children().css("fill", "white");
  $("text:contains('0x62')").parent().children().css("fill", "chocolate");
  $("text:contains('0x11')").parent().children().css("fill", "sandybrown");

});
