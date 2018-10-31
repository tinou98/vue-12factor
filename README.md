# vue-12factor [![](https://images.microbadger.com/badges/image/tinou98/vue-12factor.svg)](https://microbadger.com/images/tinou98/vue-12factor "Get your own image badge on microbadger.com")
## A container that helps to resolve the run-time configuration for a pre-built Web application.

As specified by [12 factor application](https://12factor.net/config), a strict separation of config and code is required.

But in pre-built Web application (like vue.js) code is compiled and every referecnce to environement variable are resolved at compile time.

This container allow you to use both : run-time and build-time configuration :
In your code use `{{ .MY_ENV_VARIABLE }}` to inject run time varaible that will be replaced on container startup.

## Example: Make the title change on env variable
`index.html`
```html
<!DOCTYPE html>
<html lang="fr">
  <head>
    <title>{{ .TITLE }}</title>
  </head>
  <body>
    <h1>{{ .TITLE }}</h1>
  </body>
</html>
```

`Dockerfile`
```Dockerfile
# Build application
FROM node:alpine as build-stage

WORKDIR /app

# Install package as a separate layer for improved docker caching
COPY package*.json ./
RUN npm install

# Build the application
COPY . .
RUN npm run build


# Runnable image
FROM tinou98/vue-12factor

COPY --from=build-stage /app/dist /srv/http
CMD ["js/*.js", "index.html"] # List file that contain run-time variable
```

Then build and run the image:
```shell
docker build -t vue-12factor-example .
docker run -it --rm -p8080:80 -e TITLE="Changeable title" vue-12factor-example
```

Now a server is running on port 8080, serving the page with injected environement.
