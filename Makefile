ns:
	kubectl apply -f ./k8s/ns.yaml



dep: ns
	kubectl apply -f ./k8s/dep.yaml -n tst


srv: dep
	kubectl apply -f ./k8s/serv.yaml -n tst