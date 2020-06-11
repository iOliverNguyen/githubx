function buildCtrl($milestones) {
  let mapMilestones = {};
  for (let i = 0; i < $milestones.length; i++) {
    let $ms = $milestones[i];
    let attr = $ms.getAttribute('data-card-filter');
    let ms = JSON.parse(attr.replace('milestone:',''));
    if (mapMilestones[ms]) continue;
    mapMilestones[ms] = $ms;
  }

  let milestones = Object.keys(mapMilestones).sort();
  let buttons = milestones.map((ms) => {
    let btn = document.createElement('button');
    btn.innerText = ms;
    btn.onclick = (event) => {
      mapMilestones[ms].click();
    };
    return btn;
  });

  $ctrl = document.createElement('div');
  $ctrl.classList.add('d-flex');
  buttons.forEach((btn) => $ctrl.appendChild(btn));
  return $ctrl;
}

function init() {
  let $milestones = [...(document.querySelectorAll('[data-card-filter^="milestone:"]'))];
  $msCtrl = buildCtrl($milestones);

  let $header = document.querySelector('.project-header');
  let $headerCtrl = $header.querySelector('.project-header-controls');
  $header.insertBefore($msCtrl, $headerCtrl);
}

setTimeout(init, 3000);
