function copyEmail(elementID) {
  const element = document.getElementById(elementID);
  if (!element) return;

  //copy email
  const copyText = element.innerText;

  navigator.clipboard
    .writeText(copyText)
    .then(() => {
      console.log("Email copied");
    })
    .catch((err) => {
      console.error("Failed to copy email: ", err);
    });
}
