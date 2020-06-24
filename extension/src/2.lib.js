;(function() {
  function Query(token) {
    return {
      buildURL: (path) => baseURL + path,
      query: async (query) => {
        return fetch(graphURL, {
          method: 'POST',
          headers: {
            Authorization: 'Bearer ' + token
          },
          body: JSON.stringify({query: query})
        });
      },
      request: async (path, body, type = 'json') => {
        let resp = await fetch(baseURL + path, {
          method: 'POST',
          headers: {
            Authorization: 'Bearer ' + token
          },
          body: JSON.stringify(body)
        });
        if (resp.status !== 200) throw new Error(await resp.text());
        return resp.json();
      }
    }
  }

  window.Query = Query;
})();

function formatTime(time) {
  let t = new Date(time);
  let now = new Date();
  if (t.toDateString() === now.toDateString()) {
    return t.toTimeString().substr(0, 5);
  }
  if (now - t >= 0 && now - t <= 6*60*60*24*1000) {
    return t.toDateString().substr(0,3);
  }
  return t.getDay() + '/' + t.getMonth() + '/' + (t.getFullYear()+'').substr(2,2);
}

function generateID() {
  let a = new Uint32Array(2);
  crypto.getRandomValues(a);
  return ('' + a[0] + a[1]).substr(2);
}
