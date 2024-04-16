IMAGE_NAME = work
CONTAINER_NAME = work4

build_app_image:
	docker build -t $(IMAGE_NAME):latest -f ./docker-build/Dockerfile .

build:
	docker compose --env-file=./docker-build/docker-compose-env.env --file=./docker-build/docker-compose.yaml --profile release up -d

hz_gen:
	hz update -idl idl/user.proto
	hz update -idl idl/video.proto
	hz update -idl idl/interact.proto
	hz update -idl idl/relation.proto
	hz update -idl idl/test.proto