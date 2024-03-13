# Node

[Daggerverse](https://daggerverse.dev/mod/github.com/Dudesons/daggerverse/node)
![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)

A Nodejs module

## Features

* Expose basic functions like testing, transpile, clean, select the right package manager (more with `dagger function -m github.com/Dudesons/daggerverse/node`)
* 2 Lazy functions:
   * `with-auto-setup` which will extract information from the project
     * detect if tests are present
     * detect if it's package or not
     * detect if lint command is available
     * detect the package manager
     * Information like name, version, engine version ...
   * `pipeline`: Ideally call after `with-auto-setup`, this function will execute all the pipeline from the source to a package / docker image

## Prerequisite for lazy functions

The module expects to find some information in the `package.json`:
 * Fields: 
   * `.name`
   * `.version`
   * `.engines.node`
 * Scripts:
   * `test` (required if test are find): expect a command to run tests
   * `build` (required): expect to command to build / transpile the code
   * `clean` (required): cleanup the project 
   * `lint` (optional): expect a command to check the lint not fixing it

## Examples

### Lazy method

This example will execute the whole pipeline which consist in:
 * install dependencies
 * test the project
 * transpile it
 * push the image to the ttl.sh registry

```go
dag.
    Node().
    WithAutoSetup(
        "testdata-fastify",
        testDataSrc.Directory("myapi"),
    ).
    Pipeline(
        ctx,
      NodePipelineOpts{
         DryRun: true,
         TTL:    "5m",
         IsOci:  true,
      },
    )
```

```shell
dagger call -m "github.com/Dudesons/daggerverse/node" \
  with-auto-setup --pipeline-id="testdata-myapi" --src=../testdata/node/myapi/ \
  pipeline --dry-run=true --ttl=5m --is-oci=true
```

In this example we are building a npm package:
```go
dag.
   Node().
   WithAutoSetup(
       "testdata-lib",
       testDataSrc.Directory("mylib"),
   ).
   Pipeline(
       ctx,
       NodePipelineOpts{
           DryRun:        true,
           PackageDevTag: "beta",
       },
   )
```

### Create a package

```go
dag.
   Node().
   WithPipelineID("testdata-mylib").
   WithVersion("20.9.0").
   WithSource(testDataSrc.Directory("mylib")).
   WithNpm().
   Install().
   Test().
   Build().
   Publish(NodePublishOpts{DryRun: true, DevTag: "beta"}).
   Do(ctx)
```

```shell
dagger call -m "github.com/Dudesons/daggerverse/node" \
  with-pipeline-id --pipeline-id="testdata-lib" \
  with-version --version=20.9.0 \
  with-source --src=../testdata/node/mylib/ \
  with-npm \
  install \
  test \
  build \
  publish --dry-run=true --dev-tag=beta \
  do
```


### Test + Transpilation
```go
dag.
   Node().
   WithPipelineID("testdata-fastify").
   WithVersion("20.9.0").
   WithSource(<Directory with the source code>).
   WithNpm().
   Install().
   Test().
   Build().
   Do(ctx)
```

```shell
dagger call -m "github.com/Dudesons/daggerverse/node" \
  with-pipeline-id --pipeline-id="testdata-fastify" \
  with-version --version=20.9.0 \
  with-source --src=../testdata/node/myapi/ \
  with-npm \
  install \
  test \
  build \
  do
```

### Open a shell or node console

```shell
dagger call -m "github.com/Dudesons/daggerverse/node" \
  with-pipeline-id --pipeline-id="testdata-fastify" \
  with-version --version=20.9.0 \
  with-source --src=../testdata/node/myapi/ \
  with-npm \
  install \
  shell
```

output:
```shell
/opt/app # ls -la                                                                                                                                                                                                                                                                                          
total 208                                                                                                                                                                                                                                                                                                  
drwxr-xr-x    1 root     root          4096 Mar  6 06:41 .                                                                                                                                                                                                                                                 
drwxr-xr-x    1 root     root          4096 Mar  6 06:41 ..                                                                                                                                                                                                                                                
-rw-rw-r--    1 root     root           307 Mar  4 05:16 mockData.ts                                                                                                                                                                                                                                       
drwxr-xr-x  204 root     root         12288 Mar  5 08:22 node_modules                                                                                                                                                                                                                                      
-rw-rw-r--    1 root     root        168920 Mar  5 22:59 package-lock.json                                                                                                                                                                                                                                 
-rw-rw-r--    1 root     root           765 Mar  5 08:12 package.json                                                                                                                                                                                                                                      
drwxrwxr-x    2 root     root          4096 Mar  4 06:40 src                                                                                                                                                                                                                                               
drwxrwxr-x    2 root     root          4096 Mar  4 05:16 test                                                                                                                                                                                                                                              
-rw-rw-r--    1 root     root          1297 Mar  4 06:38 tsconfig.json                                                                                                                                                                                                                                     
/opt/app # npm run test                                                                                                                                                                                                                                                                                    
                                                                                                                                                                                                                                                                                                           
> test                                                                                                                                                                                                                                                                                                     
> vitest run                                                                                                                                                                                                                                                                                               
                                                                                                                                                                                                                                                                                                           
                                                                                                                                                                                                                                                                                                           
 RUN  v1.3.1 /opt/app                                                                                                                                                                                                                                                                                      
                                                                                                                                                                                                                                                                                                           
 ✓ test/app.test.ts (3)                                                                                                                                                                                                                                                                                    
   ✓ with HTTP injection                                                                                                                                                                                                                                                                                   
   ✓ with a running server                                                                                                                                                                                                                                                                                 
   ✓ with axios                                                                                                                                                                                                                                                                                            
                                                                                                                                                                                                                                                                                                           
 Test Files  1 passed (1)                                                                                                                                                                                                                                                                                  
      Tests  3 passed (3)                                                                                                                                                                                                                                                                                  
   Start at  06:42:26                                                                                                                                                                                                                                                                                      
   Duration  687ms (transform 59ms, setup 0ms, collect 246ms, tests 44ms, environment 0ms, prepare 83ms)                                                                                                                                                                                                   
                                                                                                                                                                                                                                                                                                           
/opt/app # 
```

```shell
dagger call -m "github.com/Dudesons/daggerverse/node" \
  with-pipeline-id --pipeline-id="testdata-fastify" \
  with-version --version=20.9.0 \
  with-source --src=../testdata/node/myapi/ \
  with-npm \
  install \
  shell --cmd=node
```

output:
```shell
Welcome to Node.js v20.9.0.                                                                                                                                                                                                                                                                                
Type ".help" for more information.                                                                                                                                                                                                                                                                         
>      
```

**more example in the `/ci/node.go`**

## To Do

- [ ] Add more package manager
- [ ] Add the deployment to a bucket for static files or expose the dist folder
- [ ] Improve documentation
