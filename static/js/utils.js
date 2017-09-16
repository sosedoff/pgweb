if (!Array.prototype.forEach) {
  // Simplified iterator for browsers without forEach support
  Array.prototype.forEach = function(cb) {
    if (typeof this.length != 'number') return;
    if (typeof callback != 'function') return;

    for (var i = 0; i < this.length; i++) cb(this[i]);
  }
}

function copyToClipboard(text) {
  const element = document.createElement("textarea");
  element.style.display = "none;"
  element.value = text;

  document.body.appendChild(element);
  element.focus();
  element.setSelectionRange(0, element.value.length);

  document.execCommand("copy");
  document.body.removeChild(element);
}

function guid() {
  function s4() { return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1); }
  return [s4(), s4(), "-", s4(), "-", s4(), "-", s4(), "-", s4(), s4(), s4()].join("");
}

