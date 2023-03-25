include .env
export

CONTAINER_IMAGE_NAME := ${GCP_HOSTNAME}/${GCP_PROJECT_ID}/${GCP_TARGET_IMAGE_NAME}:latest

serve:
	go run cmd/server.go

docker-build:
	docker build . --tag ${CONTAINER_IMAGE_NAME}

docker-run:
	docker run --env-file=.env -p 8080:8080 ${CONTAINER_IMAGE_NAME}

docker-push:
	docker push ${CONTAINER_IMAGE_NAME}

gcloud-deploy:
	gcloud run deploy ${GCP_SERVICE_NAME} \
		--region ${GCP_REGION} \
		--image ${CONTAINER_IMAGE_NAME} \
		--allow-unauthenticated \
		--set-env-vars SERVICE_PORT=${SERVICE_PORT},SERVICE_PROJECT=${SERVICE_PROJECT}

gcloud-service-delete:
	gcloud run services delete ${GCP_SERVICE_NAME} \
		--region ${GCP_REGION} \
		--quiet

gcloud-container-delete:
	gcloud container images delete ${CONTAINER_IMAGE_NAME}\
		--quiet

gcloud-full-build: gcloud-clean docker-build docker-push gcloud-deploy

gcloud-clean: gcloud-service-delete gcloud-container-delete