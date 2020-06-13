(function() {
  const graphURL = 'https://api.github.com/graphql'

  function Query(token) {
    return {
      query: async (query) => {
        return fetch(graphURL, {
          method: 'POST',
          headers: {
            Authorization: 'Bearer ' + token
          },
          body: JSON.stringify({query: query})
        });
      }
    }
  }

  window.Query = Query;
})();
