# artifactory-replication

ARTIFACT_TYPE: "binary" for binary artifacts replication, "docker" for docker images

JUMP_HOST_NAME: jumphost hostname

JUMP_HOST_USER: username to connect to jumphost

JUMP_HOST_KEY: key to use when connecting to jumphost

JUMP_HOST_DESTINATION: destination host and port to proxy

JUMP_HOST_LOCAL_PORT: local port to use when connecting through jumphost

# env variables for docker registry to docker registry replication:

SOURCE_REGISTRY: source docker registry to sync from

DESTINATION_REGISTRY: destination docker registry to sync to

DESTINATION_REGISTRY_TYPE: aws, google or alicloud

DOCKER_TAG: replicate only specific tag for all images in source repo

IMAGE_FILTER: image path prefix

SOURCE_USER: source registry user, if needed

SOURCE_PASSWORD: source registry password, if needed

DESTINATION_USER: destination registry user, if needed

DESTINATION_PASSWORD: destination registry password, if needed

DOCKER_CLEAN: clean destination repos

DOCKER_CLEAN_KEEP_TAGS: clean oldest N tags from destination registry, 10 if not specified

SOURCE_PROD_REGISTRY: source prod registry, to exclude images from cleanup

SOURCE_PROD_REGISTRY_USER: user for prod registry

SOURCE_PROD_REGISTRY_PASSWORD: password for prod registry


# env variables for artifactory binary to S3 replication

SOURCE_REGISTRY: source artifactory binary repo to sync from

DESTINATION_REGISTRY: destination binary registry name to sync to

DESTINATION_REGISTRY_TYPE: destination registry type

IMAGE_FILTER: image path repository, recursive copy not supported, specify inmost directory

SOURCE_USER: source repository user, if needed

SOURCE_PASSWORD: source repository password, if needed

AWS_ACCESS_KEY_ID: aws access key id for user with write access to s3 bucket

AWS_SECRET_ACCESS_KEY: aws secret access key for user with write access to s3 bucket

AWS_REGION: aws region where bucket is located

HELM_CDN_DOMAIN: domain name for cdn to use in helm charts