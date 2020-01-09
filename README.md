# `lifecycle-build` task

A Concourse task for building [OCI
images](https://github.com/opencontainers/image-spec) using the [Cloud Native Buildpacks](https://buildpacks.io/) [lifecycle](https://github.com/buildpacks/lifecycle) 

## Usage

There is a task implementation for each of the Cloud Foundry CNB builders
* cloudfoundry/cnb:tiny -> bstick12/lifecycle-build-task:tiny
* cloudfoundry/cnb:bionic -> bstick12/lifecycle-build-task:bionic
* cloudfoundry/cnb:cflinuxfs3 -> bstick12/lifecycle-build-task:cflinuxfs3

You can see the implementations at [lifecycle-build-task](https://hub.docker.com/r/bstick12/lifecycle-build-task/tags)

You provide the task configuration as follows:

### `image_resource`

Point your task as the `lifecycle-build-task` image resource

```yaml
 image_resource:
    type: registry-image
    source:
      repository: bstick12/lifecycle-build-task
      # Change this if you want to use a different CNB builder e.g. bionic, cflinuxfs3
      tag: tiny
```

### `params`

Next, any of the following parameters can be specified:

* `REPOSITORY`:  Required. The name of the repository that you want to build e.g bstick12/built-with-cnb-lifecycle-task
* `TAG`: Optional. The tag value for the image. Defaults to `latest`
* `REGISTRY_USERNAME`: Required. The username to authenticate with when pushing
* `REGISTRY_PASSWORD`: Required. The password to use when authenticating
* `CONTEXT`: Required. The directory name of the input source you want to build
* `DEBUG`: Optional. Log debug output of the task. Defaults to `false`
* `BUILD_ENV` (default empty): a hash of environment variables that will be passed to the CNB builder

### `inputs`

There are no required inputs - your task should just list each artifact it
needs as an input. Typically this is in close correlation with `$CONTEXT`:

```yaml
params:
  CONTEXT: git-my-source

inputs:
- name: git-my-source
```

### `run`

Your task should run the `lifecycle-build-task` executable:

```yaml
run:
  path: lifecycle-build-task
```

## Example

The example below builds the [goflake-server](https://github.com/bstick12/goflake-server) using the `cloudfoundry/cnb:tiny` builder

```yaml
resources:
- name: git-goflake-server
  type: git
  icon: "github-circle"
  source:
    uri: https://github.com/bstick12/goflake-server

jobs:
- name: build-goflake-server-tiny
  plan:
  - get: git-goflake-server
    trigger: true
  - task: goflake-server-lifecycle
    config:
      platform: linux
      # The image resource configuration for the lifecycle-build task
      image_resource:
        type: registry-image
        source:
          repository: bstick12/lifecycle-build-task
          # Change this if you want to use a different CNB builder e.g. bionic, cflinuxfs3
          tag: tiny
      params:
        REPOSITORY: bstick12/goflake-server-lifecycle
        TAG: tiny
        CONTEXT: git-goflake-server
        REGISTRY_USERNAME: ((docker-hub-username))
        REGISTRY_PASSWORD: ((docker-hub-password))
      inputs:
      - name: git-goflake-server
      run:
        path: lifecycle-build-task
```


## Developing

There is `Makefile` provided to build and test the `lifecycle-build-task` locally
