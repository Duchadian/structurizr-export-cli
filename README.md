
# structurizr-export-cli

A small cli for exporting views as PNG images from [Structurizr](https://structurizr.com/) using a Chrome browser.
It will automatically look for the views defined in your `workspace.dsl` file and export those to PNG images.
Like [Structurizr Puppeteer](https://github.com/structurizr/puppeteer), this connects to the Structurizr instance and exports images from there directly, no conversions in between.
Very much a work in progress, not recommended for users that experience discomfort at seeing Golang stack traces.

Heavily borrows from [Structurizr Puppeteer](https://github.com/structurizr/puppeteer) for much of its logic.
This tool can do a subset of what Structurizr Puppeteer does, but without the pain of NodeJS dependencies.
So far, it has only been tested on [Structurizr Lite](https://docs.structurizr.com/lite). 
Any other version of Structurizr might or might not work with this tool, no guarantees are given.

## Usage/Examples
Download the latest release from the [releases page](https://github.com/Duchadian/structurizr-export-cli/releases).


### Local run

For this example, the assumption is that you have a Structurizr Lite instance running with Docker:
```shell
# assumed Structurizr setup
docker run --rm -p 8080:8080 -d --name structurizr -v <folder with your workspace.json and workspace.dsl>:/usr/local/structurizr structurizr/lite
```

This tool could then be run as follows:
```shell
./structurizr-export-cli http://localhost:8080
```

Once run, it will download Chrome (if it is not available already), and start a remotely controllable instance.
This instance cannot be headless (i.e. invisible), because the images do not get loaded otherwise. 
It then loops over the views and exports them to the `export` directory. 
The directory is configurable with the `--export-dir` flag.

### Remote Run

There are situations in which the tool cannot be run against a local Chrome instance (e.g. CI). 
In these situations, running a non-headless Chrome browser is usually difficult. 
The best approach is to run a pre-made `rod` container that already has all the dependencies to do this properly:
```shell
# example rod container
docker run --rm -d --name rod -p 7317:7317 ghcr.io/go-rod/rod
```

the cli can then be configured to use this instance:
```shell
./structurizr-export-cli --rod-remote=ws://<your rod container>:7317 <your structurizr url> 
```

Keep in mind that your Structurizr url needs to be resolvable _from the rod container_. 
`localhost` will likely not work.
An example Gitlab job to demonstrate what configuration _will_ work (assuming you have the cli in your project root):

```yaml
extract_diagrams:
  stage: extract_diagrams
  image: docker:19
  variables:
    DOCKER_DRIVER: overlay2
    DOCKER_HOST: tcp://docker:2375
    DOCKER_TLS_CERTDIR: ""
    FF_NETWORK_PER_BUILD: "true"
  services:
    - name: docker:dind
      alias: docker
  script:
    - docker run --rm -d --name structurizr -p 8080:8080 -v "$(pwd):/usr/local/structurizr" structurizr/lite
    - docker run --rm -d --name rod -p 7317:7317 ghcr.io/go-rod/rod
    - sleep 5
    - ./structurizr-export-cli --rod-remote=ws://docker:7317 http://$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' structurizr):8080
    - docker rm --force structurizr
    - docker rm --force rod
  artifacts:
    paths:
      - "export/*.png"
    untracked: false
    when: on_success
```

## TODO
- add authentication option
- add SVG export option
- clean up error messages, especially with remote rod
- add automated tests
- CI releases
