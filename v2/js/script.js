$(function () {
  $("input").change(function () {
    var label = $(this).parent().find("span");
    if (typeof this.files != "undefined") {
      // fucking IE
      if (this.files.length == 0) {
        label.removeClass("withFile").text(label.data("default"));
      } else {
        var file = this.files[0];
        var name = file.name;
        var size = (file.size / 1048576).toFixed(3); //size in mb
        label.addClass("withFile").text(name + " (" + size + "mb)");
      }
    } else {
      var name = this.value.split("\\");
      label.addClass("withFile").text(name[name.length - 1]);
    }
    return false;
  });
});

function toggle_div_fun(div_id, btn_id, button_labels) {
  let first_label, second_label;

  first_label = button_labels[0];
  if (typeof button_labels[0] === "undefined") {
    first_label = "Открыть форму";
  }

  second_label = button_labels[1];
  if (typeof button_labels[1] === "undefined") {
    second_label = "Закрыть форму";
  }

  var divelement = document.getElementById(div_id);
  var button = document.getElementById(btn_id);

  if (divelement.style.height == "0px" || divelement.style.height == "") {
    divelement.style.height = divelement.scrollHeight + "px";
    button.textContent = second_label;
  } else {
    divelement.style.height = "0px";
    button.textContent = first_label;
  }
}
