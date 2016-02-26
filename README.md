# azkube


## Overview
Tool used to deploy and bootstrap a Kubernetes cluster in Azure.


## Usage

### From source
```
go get github.com/colemickens/azkube
$GOPATH/bin/azkube --tenant-id="{your tenant id}" --subscription-id="{subscription id}"
```

### Docker
```
docker run -it colemickens/azkube:latest /azkube --tenant-id="{your tenant id}" --subscription-id="{subscription id}"
```


## Motivations
1. Existing shell script was fragile.
2. azure-xplat-cli changes out from underneath of us, and is slow, and doesn't handle errors well.
3. Need a tool and process to create service principals and configure them appropriately.
4. Need a tool to consume scripts, interweave ARM template variables/parameters, and output a "deployable" template.


## Future Improvements
1. Introduce concept of "flavors" with a defined interface that we hand off
2. Introduce an ubuntu 16.04 LTS (beta) flavor


## Notes
1. The user who executes the application must have permission to provision additional applications.
2. The resulting "templates" are fully parameterized and generic. They can be uploaded and used by others.


## Todo
1. create a way to bypass needing to create more applications - this requires a user present, or requires giving the runner Company User Admin permissions which is not encouraged.
2. create a way of specifying the ssh key to use
