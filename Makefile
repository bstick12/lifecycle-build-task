.PHONY: all create cf-tiny cf-bionic cf-cflinuxfs3

all: create

create: tiny bionic cflinuxfs3

tiny:
	docker build . -t bstick12/pack-lifecycle-resource:tiny --build-arg BUILDER=cloudfoundry/cnb:tiny --build-arg DOCKER_USERNAME=${DOCKER_USERNAME} --build-arg DOCKER_PASSWORD=${DOCKER_PASSWORD}

bionic:
	docker build . -t bstick12/pack-lifecycle-resource:bionic --build-arg BUILDER=cloudfoundry/cnb:bionic --build-arg DOCKER_USERNAME=${DOCKER_USERNAME} --build-arg DOCKER_PASSWORD=${DOCKER_PASSWORD}

cflinuxfs3:
	docker build . -t bstick12/pack-lifecycle-resource:cflinuxfs3 --build-arg BUILDER=cloudfoundry/cnb:cflinuxfs3 --build-arg DOCKER_USERNAME=${DOCKER_USERNAME} --build-arg DOCKER_PASSWORD=${DOCKER_PASSWORD}
