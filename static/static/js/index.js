function highlight() {
  var key;
  var val = $("input[name='key']").val();

  if (!!val) {
    var mark = new Mark($(".file-list"));
    for (var i = 0; i < val.length; i++) {
      var code = val.charCodeAt(i);
      if (code > 255) {
        mark.mark(val[i]);
      } else if (val[i] == " ") {
        mark.mark(key);
        key = "";
      } else {
        key += val[i];
      }
    }
  }
}

function showAll() {
  $(".show-all-files").click(function() {
    var infoHash = $(this).data("infohash");
    if ($(this).text() == "show all") {
      $(".file-item-" + infoHash).show();
      $(this).text("hide");
    } else {
      $(".file-item-" + infoHash).hide();
      $(this).text("show all");
    }
  });
}

$(document).ready(function() {
  NProgress.configure({showSpinner: false});
  NProgress.start();

  $(window).load(function() {
    NProgress.done();
  });

  highlight();
  showAll();
});
