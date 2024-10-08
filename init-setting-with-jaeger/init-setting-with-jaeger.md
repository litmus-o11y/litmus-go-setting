# init-setting-with-jaeger

### 1. 타겟 애플리케이션
kubectl create deployment nginx --image=nginx --replicas=1  

### 2. Cert manager 설치  
참고: https://cert-manager.io/docs/installation/kubectl/  
(1) kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.yaml  
(2) watch -n 1 "kubectl get all -n default"  

### 3. Jaeger operator 설치
참고: https://www.jaegertracing.io/docs/1.60/operator/   
(1) kubectl create namespace observability   
(2) kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.60.0/jaeger-operator.yaml -n observability   
(3) watch -n 1 "kubectl get all -n observability"  
(4) kubectl apply -f simplest.yaml -n observability   
(5) watch -n 1 "kubectl get all -n observability"  
(6) minikube service simplest-query -n observability --url   

### 4. Litmus 설치
(1) kubectl create ns litmus   
(2) helm install chaos litmuschaos/litmus \
--namespace=litmus \
--set portal.frontend.service.type=NodePort \
--set mongodb.image.registry=ghcr.io/zcube \
--set mongodb.image.repository=bitnami-compat/mongodb \
--set mongodb.image.tag=6.0.5   
(3) watch -n 1 "kubectl get all -n litmus"   
(4) minikube service chaos-litmus-frontend-service -n litmus --url    
(admin/litmus) -> (admin/Litmus1@)   

### 5. Litmus Environments
(1) Name: local   
(2) 환경 yaml: subscriber server addr 변경
- kubectl get service chaos-litmus-frontend-service -n litmus  
- #652 수정  
  SERVER_ADDR: http://10.97.52.58:9091/api/query   

(3) cd ~/Downloads   
(4) kubectl apply -f local-litmus-chaos-enable.yml   

### 6. Litmus Resilience Probes (Command)
- Name: nginx-probe   
- 3s / 3s / 3   
- kubectl get pods --all-namespaces | grep nginx | grep Running | wc -l   
- Int / > / 0   

### 7. Litmus Chaos Experiments
(1) probe 설정    
- nginx-probe   
- default / deployment / app=nginx   
- EOT   

(2) 실험 yaml 수정  
- #140 수정   
  image: lak9348/go-runner:t20   
- #168 추가   
  name: OTEL_EXPORTER_OTLP_ENDPOINT     
  value: "simplest-collector.observability.svc.cluster.local:4317"      
