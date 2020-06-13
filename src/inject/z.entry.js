function init() {
  let path = location.pathname;
  if (/[a-z0-9-_/]+\/projects\/[0-9]+/.exec(path)) setTimeout(setupProjectCtrl, 3000);

  setupSidebar();
}

init();
