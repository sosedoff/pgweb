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
