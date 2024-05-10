# dockonaut
Navigating Your Containers, One Stack at a Time!

### How to user docker image:
```
docker run -itd --privileged --rm -v "/var/run/docker.sock:/var/run/docker.sock" \
	-v "$PWD/config.json:$PWD/config.json" -v "$PWD/workspace:$PWD/workspace" \
	-v "$PWD/demo:$PWD/demo" -w "$PWD" dockonaut
```

you need $PWD path also on the container because if you are going to mount another volume within the dockonaut the path should match the host path.
