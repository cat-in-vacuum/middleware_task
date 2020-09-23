### Test task for Middleware Inc

#### Usage ####

    go build cmd/cmd.go --port=5005

#### Endpoints ####

    /api/v1/notifications

**request example**

POST:

    {"URL": "https://jsonplaceholder.typicode.com/posts/1"}
    {"URL": "https://jsonplaceholder.typicode.com/posts/need404"}
    {"URL": "unsopportedschema(://jsonplaceholder.typicode.com/posts/3"}

**response**

    {"url":"unsopportedschema(://jsonplaceholder.typicode.com/posts/3","error":"parse unsopportedschema(://jsonplaceholder.typicode.com/posts/3: invalid URI for request"}
    {"url":"https://jsonplaceholder.typicode.com/posts/need404","code":"404 Not Found","error":"non_200_resp_status"}
    {"url":"https://jsonplaceholder.typicode.com/posts/1","body":{"body":"qdent occaecati excepturi optio ... reprehenderit","userId":1},"code":"200 OK"}