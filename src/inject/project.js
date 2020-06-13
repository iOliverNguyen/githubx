function setupProjectCtrl() {

  function buildProjectCtrl($milestones) {
    let mapMilestones = {};
    for (let i = 0; i < $milestones.length; i++) {
      let $ms = $milestones[i];
      let attr = $ms.getAttribute('data-card-filter');
      let ms = attr.replace('milestone:', '');
      if (ms.indexOf('"') === 0) ms = JSON.parse(ms);
      mapMilestones[ms] = true;
    }

    let milestones = Object.keys(mapMilestones).sort();
    milestones = milestones.filter((ms) => {
      let q = 'article';
      q += '[data-card-milestone*=' + JSON.stringify(JSON.stringify(ms.toLowerCase())) + ']';
      q += '[data-card-state*=' + JSON.stringify(JSON.stringify('open')) + ']';
      console.log('q', q, $(q));
      return $(q);
    });

    let buttons = milestones.map((ms) => {
      let btn = document.createElement('button');
      btn.innerText = ms;
      btn.classList.add("x-btn", "btn", "btn-block", "btn-outline");
      btn.onclick = (event) => {
        $('button.issues-reset-query').click();
        let $input = $('input.js-card-filter-input');
        $input.value = 'milestone:' + JSON.stringify(ms.toLowerCase());

        let e = document.createEvent('Event');
        e.initEvent('keypress');
        e.which = e.keyCode = 13;
        $input.dispatchEvent(e);
      };
      return btn;
    });

    $ctrl = document.createElement('div');
    $ctrl.classList.add('x-ctrl', 'd-flex');
    buttons.forEach((btn) => $ctrl.appendChild(btn));
    return $ctrl;
  }

  function setup() {
    let $milestones = [...($$('button[data-card-filter^="milestone:"]'))];
    $msCtrl = buildProjectCtrl($milestones);

    let $header = $('.project-header');
    let $headerCtrl = $header.querySelector('.project-header-controls');
    $header.insertBefore($msCtrl, $headerCtrl);
  }

  setup();
}
