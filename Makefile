 DOCKERHUB_USERNAME ?= muarwi
 APPLICATION_NAME ?= go-simple-microservice
 build:
	
	docker build -t ${DOCKERHUB_USERNAME}/${APPLICATION_NAME}-backend -f ./backend/Dockerfile ./backend
	docker build -t ${DOCKERHUB_USERNAME}/${APPLICATION_NAME}-frontend -f ./frontend/Dockerfile ./frontend
push:
	docker login -u ${DOCKERHUB_USERNAME}
	docker push ${DOCKERHUB_USERNAME}/${APPLICATION_NAME}-backend
	docker push ${DOCKERHUB_USERNAME}/${APPLICATION_NAME}-frontend