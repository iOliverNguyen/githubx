let apikey = localStorage.getItem('githubx-apikey');
let Q = Query(apikey);
let clientID = generateID();

function setupSidebar() {
  const appHTML = `
<div class="x-nav" v-bind:class="{'x-loading':isLoading}">
    <div class="x-title">Conversations</div>
    <a class="x-settings" v-on:click="settings"></a>
    <div class="x-list">
        <div class="x-item x-active"><span>all</span></div>
        <div class="x-item"><span>my team</span></div>
        <div class="x-item"><span>followed</span></div>
        <div class="x-item"><span>assigned</span></div>
    </div>
</div> 
<div class="x-conv">
<div class="x-conv-list">
    <div class="x-conv-item" v-for="$is in issues" :key="$is.number" v-on:click="loadURL($event, $is.url)">
      <div class="x-header">
        <div class="x-title">
          <div class="x-time">{{formatTime($is.lastChangedAt)}}</div>
          <a href="#" v-bind:href="$is.url">{{$is.title}}</a>
        </div>
      </div>
      <div class="x-info">#{{$is.number}}</div>
      <div class="x-body">{{issueLatestText($is)}}</div>
    </div>
</div>
</div>
`;

  function buildSidebar() {
    let $app = document.createElement('div');
    $app.id = 'x-sidebar';
    $app.innerHTML = appHTML;
    $app.querySelector('.x-settings').innerHTML = iconGear;
    document.body.appendChild($app);

    app = new Vue({
      el: '#x-sidebar',
      data: {
        isLoading: false,
        issues: []
      },
      methods: {
        formatTime: formatTime,
        issueLatestText: issueLatestText,
        loadURL: loadURL,
        settings: async () => {
          let _apikey = prompt('Input GitHub API KEY:');
          if (!_apikey || !_apikey.trim()) return

          let resp = await Query(_apikey).request('/api/Authorize', {});
          console.log('authorized as', resp)

          apikey = _apikey
          localStorage.setItem('githubx-apikey', apikey)
          Q = Query(apikey)
          loadDiscussions();
        },
      }
    });
  }

  function issueLatestText(is) {
    let lastCmt = is.comments && is.comments[is.comments.length-1];
    if (lastCmt) {
      return lastCmt.bodyText;
    }
    return is.bodyText;
  }

  async function init() {
    let resp = await Q.request('/api/Authorize', {});
    console.log('authorized as', resp);
  }

  async function loadDiscussions() {
    let resp = await Q.request('/api/ListIssues', {})
    console.log('issues', resp);

    app.issues = resp.issues;
  }

  // const githubAjax = '#js-repo-pjax-container, div[itemtype="http://schema.org/SoftwareSourceCode"] main, [data-pjax-container]';

  async function loadURL(event, url) {
    if (!url.startsWith('https://github.com/')) return;
    event.preventDefault();

    try {
      this.isLoading = true;
      window.history.replaceState( {} , 'GitHub', url );
      let resp = await fetch(url);
      let html = await resp.text();
      this.isLoading = false;
      let idx0 = html.indexOf('<main ');
      let idx1 = html.indexOf('</main>');
      if (idx0 < 0 || idx1 < 0 ) {
        console.error('can not detect main component');
        window.location = url;
        return
      }
      html = html.substring(idx0, idx1);
      $('main').innerHTML = html;

    } catch(e) {
      window.location = url;
    }
  }

  function processEvents(events) {
    if (events.length === 0) return;
    let issues = app.issues;
    for (let e of events) {
      if (e.type !== 'issue') continue;
      let is = e.data;
      let idx;
      for (let i = 0; i < issues.length; i++) {
        if (issues[i].id === is.id) { idx = i; break }
      }
      if (idx === undefined) {
        issues.unshift(e.data);
        continue;
      }
      issues = issues.slice(0,idx).concat(issues.slice(idx+1));
      issues.unshift(e.data);
    }
    app.issues = issues;
    console.log('processed events', issues);
  }

  let retryAfter = 1000;
  async function startListening() {
    try {
      let resp = await Q.request('/api/poll?id=' + clientID);
      console.log('events', resp);
      setTimeout(startListening, 4);
      processEvents(resp);

    } catch(e) {
      setTimeout(startListening, retryAfter);
      retryAfter = retryAfter * 2;
      if (retryAfter > 60000) retryAfter = 1000;
    }
  }

  async function start() {
      await init();
      await loadDiscussions();
      startListening();
  }

  buildSidebar();
  start();
}
