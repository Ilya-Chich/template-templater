const menu = document.querySelector(".menu-btn");
const nav = document.querySelector("#site-nav");
if (menu && nav) {
  menu.addEventListener("click", () => {
    const open = menu.getAttribute("aria-expanded") === "true";
    menu.setAttribute("aria-expanded", open ? "false" : "true");
    nav.hidden = open;
  });
}
