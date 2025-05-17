 DOCKERHUB_USERNAME ?= muarwi
 APPLICATION_NAME ?= go-simple-microservice
 build:
	docker build -t ${DOCKERHUB_USERNAME}/${APPLICATION_NAME}-backend -f ./backend/Dockerfile ./backend
	docker build -t ${DOCKERHUB_USERNAME}/${APPLICATION_NAME}-frontend -f ./frontend/Dockerfile ./frontend
push:
	docker login -u ${DOCKERHUB_USERNAME}
	docker push ${DOCKERHUB_USERNAME}/${APPLICATION_NAME}-backend
	docker push ${DOCKERHUB_USERNAME}/${APPLICATION_NAME}-frontend
kubernetes:
	kubectl apply -f k8s/00_namespace.yaml
	kubectl apply -f k8s/01_postgres_storage.yaml
	kubectl apply -f k8s/02_postgres_credentials.yaml
	kubectl apply -f k8s/03_postgres_deployment.yaml
	kubectl apply -f k8s/04_postgres_service.yaml
	kubectl apply -f k8s/05_backend-deployment.yaml
	kubectl apply -f k8s/06_backend-service.yaml
	kubectl apply -f k8s/07_frontend-deployment.yaml
	kubectl apply -f k8s/08_frontend-service.yaml
	kubectl apply -f k8s/09_ingress.yaml
