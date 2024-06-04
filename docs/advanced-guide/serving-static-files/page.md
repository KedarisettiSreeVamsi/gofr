# Serving Static Files using GoFr

Often, we require to serve static content be it a default profile image or a static website. We want to have a mechanism to serve those content without having a hassel of implementing it from scratch.

GoFr provides a default mechanism where if a public folder is available in the directory of the application, it automatically provides an endpoint with `/public/<filename>`, here filename refers to the file we want to get static content to be served. 

Example Project folder utilizing public endpoint:

```
project_folder
|
|---config
|       .env
|---public
|       <img1>.jpeg
|       <img2>.png
|       <img3>.jpeg
|   main.go
|   main_test.go
```

main.go code:

```go
package main

import "gofr.dev/pkg/gofr"

func main(){
    app := gofr.New()
    app.Run()
}

```

Additionally,if we want to serve more static endpoints, we have a dedicated function called `AddStaticFiles()` which takes 2 parameters endpoint and the filepath of the static folder which we want to serve.

Providing an example below along with File System Example:

```
project_folder
|
|---config
|       .env
|---public
|       <img1>.jpeg
|       <img2>.png
|       <img3>.jpeg
|---static
|       |---css
|       |       main.css
|       |---js
|       |       main.js
|       |   index.html
|   main.go
|   main_test.go
```


main.go file:

```go

package main

import "gofr.dev/pkg/gofr"

func main(){
    app := gofr.New()
    app.AddStaticFiles("static","./static")
    app.Run()
}

```

In the above example, both endpoints `/public` and `/static` are available for the app to render the static content.