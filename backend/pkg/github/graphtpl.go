package github

import "text/template"

var loginTpl = template.Must(template.New("login").Parse(`
{
  viewer {
    id
    login
    name
  }
  organization(login: {{ .Org | printf "%q" }}) {
    id
    name
    repository(name: {{ .Repo | printf "%q" }}) {
      id
      name
    }
  }
}`))

var listIssueTpl = template.Must(template.New("listIssues").Parse(`
{
  organization(login: {{ .Org | printf "%q" }}) {
    repository(name: {{ .Repo | printf "%q" }}) {
      issues(filterBy: {}, orderBy: {field: CREATED_AT, direction: ASC}, {{ .IssuePaging }}) {
        nodes {
          body
          title
          bodyHTML
          bodyText
          author {
            login
            url
          }
          assignees(first: 100) {
            totalCount
            nodes {
              id
              login
              name
              url
            }
          }
          closed
          closedAt
          updatedAt
          labels(first: 100) {
            totalCount
            nodes {
              color
              description
              id
              name
              url
            }
          }
          id
          url
          state
          number
          createdAt
          comments(last: 100) {
            totalCount
            pageInfo {
              hasNextPage
              hasPreviousPage
              endCursor
              startCursor
            }
            nodes {
              id
              author {
                login
                url
              }
              body
              bodyHTML
              bodyText
              url
              updatedAt
              createdAt
            }
          }
        }
        pageInfo {
          endCursor
          hasNextPage
          hasPreviousPage
          startCursor
        }
        totalCount
      }
      labels(first: 100) {
        totalCount
        nodes {
          color
          description
          id
          name
          url
        }
      }
    }
  }
}`))
