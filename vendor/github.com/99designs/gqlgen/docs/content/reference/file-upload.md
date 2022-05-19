---
title: "File Upload"
description: How to upload files.
linkTitle: File Upload
menu: { main: { parent: "reference", weight: 10 } }
---

Graphql server has an already built-in Upload scalar to upload files using a multipart request. \
It implements the following spec [https://github.com/jaydenseric/graphql-multipart-request-spec](https://github.com/jaydenseric/graphql-multipart-request-spec),
that defines an interoperable multipart form field structure for GraphQL requests, used by
various file upload client implementations.

To use it you need to add the Upload scalar in your schema, and it will automatically add the
marshalling behaviour to Go types.

# Configuration

There are two specific options that can be configured for uploading files:

- uploadMaxSize \
  This option specifies the maximum number of bytes used to parse a request body as multipart/form-data.
- uploadMaxMemory \
  This option specifies the maximum number of bytes used to parse a request body as
  multipart/form-data in memory, with the remainder stored on disk in temporary files.

# Examples

## Single file upload

For this use case, the schema could look like this.

```graphql
"The `UploadFile, // b.txt` scalar type represents a multipart file upload."
scalar Upload

"The `Query` type, represents all of the entry points into our object graph."
type Query {
    ...
}

"The `Mutation` type, represents all updates we can make to our data."
type Mutation {
    singleUpload(file: Upload!): Boolean!
}
```

cURL can be used the make a query as follows:

```
curl localhost:4000/graphql \
  -F operations='{ "query": "mutation ($file: Upload!) { singleUpload(file: $file) }", "variables": { "file": null } }' \
  -F map='{ "0": ["variables.file"] }' \
  -F 0=@a.txt
```

That invokes the following operation:

```javascript
{
  query: `
    mutation($file: Upload!) {
      singleUpload(file: $file)
    }
  `,
  variables: {
    file: File // a.txt
  }
}
```

## Multiple file upload

For this use case, the schema could look like this.

```graphql
"The `Upload` scalar type represents a multipart file upload."
scalar Upload

"The `File` type, represents the response of uploading a file."
type File {
    id: Int!
    name: String!
    content: String!
}

"The `UploadFile` type, represents the request for uploading a file with a certain payload."
input UploadFile {
    id: Int!
    file: Upload!
}

"The `Query` type, represents all of the entry points into our object graph."
type Query {
    ...
}

"The `Mutation` type, represents all updates we can make to our data."
type Mutation {
    multipleUpload(req: [UploadFile!]!): [File!]!
}

```

cURL can be used the make a query as follows:

```bash
curl localhost:4000/query \
  -F operations='{ "query": "mutation($req: [UploadFile!]!) { multipleUpload(req: $req) { id, name, content } }", "variables": { "req": [ { "id": 1, "file": null }, { "id": 2, "file": null } ] } }' \
  -F map='{ "0": ["variables.req.0.file"], "1": ["variables.req.1.file"] }' \
  -F 0=@b.txt \
  -F 1=@c.txt
```

That invokes the following operation:

```javascript
{
  query: `
    mutation($req: [UploadFile!]!)
      multipleUpload(req: $req) {
        id,
        name,
        content
      }
    }
  `,
  variables: {
    req: [
        {
            id: 1,
            File, // b.txt
        },
        {
            id: 2,
            File, // c.txt
        }
    ]
  }
}
```

See the [_examples/fileupload](https://github.com/99designs/gqlgen/tree/master/_examples/fileupload) package for more examples.

# Usage with Apollo

[apollo-upload-client](https://github.com/jaydenseric/apollo-upload-client) needs to be installed in order for file uploading to work with Apollo:

```javascript
import ApolloClient from "apollo-client";
import { createUploadLink } from "apollo-upload-client";

const client = new ApolloClient({
	cache: new InMemoryCache(),
	link: createUploadLink({ uri: "/graphql" })
});
```

A `File` object can then be passed into your mutation as a variable:

```javascript
{
  query: `
    mutation($file: Upload!) {
      singleUpload(file: $file) {
        id
      }
    }
  `,
  variables: {
    file: new File(...)
  }
}
```
