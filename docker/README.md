# Docker Client API


## Develop with Docker Engine API

Install [Go SDK](https://docs.docker.com/engine/api/sdk/#go-sdk)

```shell
$ go get github.com/docker/docker/client
```

## Setup MongoDB Docker Container For testing

### Pull MongoDB Image

```shell
$docker pull mongo:5.0.3
```

```go
func Run(image string, ports []string, env []string) (remove func(), containerInfo ContainerInfo, err error) {

	client, err := client.NewClientWithOpts()

	if err != nil {
		panic(err)
	}

	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}
	for _, port := range ports {
		exposedPorts[nat.Port(port)] = struct{}{}
		portBindings[nat.Port(port)] = []nat.PortBinding{{
			HostIP:   "127.0.0.1",
			HostPort: "0",
		}}
	}

	resp, err := client.ContainerCreate(context.Background(), &container.Config{
		Image:        image,
		ExposedPorts: exposedPorts,
		Env:          env,
	}, &container.HostConfig{
		PortBindings: portBindings,
	}, nil, nil, "")

	if err != nil {
		panic(err)
	}

	containerID := resp.ID

	err = client.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	remove = func() {
		defer client.Close()

		err := client.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		if err != nil {
			panic(err)
		}
	}

	info, err := client.ContainerInspect(context.Background(), containerID)

	containerInfo.Hosts = make(map[string][]Host, len(ports))

	for _, port := range ports {
		for _, portBinding := range info.NetworkSettings.Ports[nat.Port(port)] {
			containerInfo.Hosts[port] = append(containerInfo.Hosts[port], Host{IP: portBinding.HostIP, Port: portBinding.HostPort})
		}
	}

	containerInfo.Status = info.State.Running

	return
}
```


# Reference

- https://docs.docker.com/engine/api
- https://docs.docker.com/engine/api/sdk/examples/
- https://github.com/moby/moby/tree/master/client