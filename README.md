# go-simple-microservices
Simple microservices in Go for demo purposes

## Front-end

```bash
cd frontend/
go mod init frontend-service
go mod tidy
```

## Front-end

```bash
cd backend/
go mod init backend-service
go get .
go mod tidy
```

## Troubleshoot

```bash
kubectl run multitool --image=wbitt/network-multitool -n go-microservice-app
kubectl exec -it multitool -n -n go-microservice-app -- bash
```