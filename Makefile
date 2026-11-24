API_IMAGE    := finance-api:latest
WEBAPP_IMAGE := finance-webapp:latest
# Overridable: make deploy VITE_API_URL=http://api.finance.local
VITE_API_URL ?= http://finance.local/api

.PHONY: images import deploy status verify logs teardown

build-api:
	colima nerdctl build -- -t $(API_IMAGE) ./api

build-webapp:
	colima nerdctl build -- --build-arg VITE_API_URL=$(VITE_API_URL) -t $(WEBAPP_IMAGE) ./webapp

images: build-api build-webapp

# Push the locally built images into k3s's containerd (namespace k8s.io).
import:
	colima nerdctl save $(API_IMAGE) $(WEBAPP_IMAGE) | colima ssh -- sudo k3s ctr --address /run/containerd/containerd.sock -n k8s.io images import -

# Full pipeline: build both images, import them, apply the manifests.
deploy: images import
	kubectl apply -f k8s/

status:
	kubectl get pods,svc,pvc,ingressroute,middleware

verify:
	curl -s http://finance.local/ | head -3
	curl -s http://finance.local/api/accounts

logs:
	kubectl logs deploy/api --tail=50

teardown:
	kubectl delete -f k8s/