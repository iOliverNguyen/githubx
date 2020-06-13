const api_key = 'b5aea74fadebbcd8964a9cd3b7df145b2c51dabd';

function setupSidebar() {
  let app;

  const appHTML = `
<div class="x-nav">
    <div class="x-title"></div>
    <div class="x-list">
        <div class="x-item">A</div>
        <div class="x-item">B</div>
    </div>
</div> 
<div class="x-conv-list">
    <div class="x-conv-item">
      {{ message }}
    </div>
</div>
`;

  const discussionQuery = `
{
  organization(login: "etopvn") {
    name
    teams(first: 10) {
      totalCount
      edges {
        node {
          name
          discussions(last: 10) {
            nodes {
              title
              body
              author {
                avatarUrl
                login
              }
              comments(last: 10) {
                nodes {
                  author {
                    avatarUrl
                    login
                  }
                  body
                  bodyText
                  url
                }
              }
              bodyText
              url
              updatedAt
              createdAt
              commentsUrl
            }
          }
        }
      }
    }
  }
}
`

  function buildSidebar() {
    let $app = document.createElement('div');
    $app.id = 'x-sidebar';
    $app.innerHTML = appHTML;
    document.body.appendChild($app);

    app = new Vue({
      el: '#x-sidebar',
      data: {
        message: 'hello !'
      },
      methods: {

      }
    });
  }

  const Q = Query(api_key);

  async function loadDiscussions() {
    let resp = await Q.query(discussionQuery);
    if (resp.status !== 200) {
      throw new Error('can not load discussions');
    }
    let json = await resp.json();
    console.log(json);
    let discussions = JSON.parse(await resp.json());
    console.log(discussions);
  }

  buildSidebar();
  loadDiscussions();
  return app;
}
