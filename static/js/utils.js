if (!Array.prototype.forEach) {
  // Simplified iterator for browsers without forEach support
  Array.prototype.forEach = function(cb) {
    if (typeof this.length != 'number') return;
    if (typeof callback != 'function') return;

    for (var i = 0; i < this.length; i++) cb(this[i]);
  }
}

async function copyToClipboard(text) {
  try {
    await navigator.clipboard.writeText(text);
  } catch (err) {
    console.error("Error copying to clipboard: ", err);
  }
}

function guid() {
  function s4() { return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1); }
  return [s4(), s4(), "-", s4(), "-", s4(), "-", s4(), "-", s4(), s4(), s4()].join("");
}

