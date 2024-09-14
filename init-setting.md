# init-setting


## 1. minikube 
1. 클러스터 구축
minikube start --memory 6000 --cpus 2

## 2. 타겟 어플리케이션
1. nginx
kubectl create deployment nginx --image=nginx --replicas=1


## 3. chaos-operator, chaos-runner, litmus-go
1. litmus-go
(1) 설치
- git clone -b distributed-tracing --single-branch https://github.com/namkyu1999/litmus-go.git
- 다운로드 이후 .git 삭제
(2) 수정
- pkg/utils/utils.go #6
  - **변경 전**: OTELExporterOTLPEndpoint = "OTEL_EXPORTER_OTLP_ENDPOINT"
  - **변경 후**: OTELExporterOTLPEndpoint = "otel-collector.observability.svc.cluster.local:4317"

2. chaos-runner
(1) 설치
- git clone -b distributed-tracing --single-branch https://github.com/namkyu1999/chaos-runner.git
- 다운로드 이후 .git 삭제
(2) 수정
- pkg/telemetry/otel.go #19
  - **변경 전**: const OTELExporterOTLPEndpoint = OTEL_EXPORTER_OTLP_ENDPOINT
  - **변경 후**: const OTELExporterOTLPEndpoint = "otel-collector.observability.svc.cluster.local:4317"
- Dockerfile 수정
  - **변경 전**: 
    - ADD ./chaos-runner /chaos-runner
      WORKDIR /chaos-runner
  - **변경 후**: 
    - ADD ./chaos-runner /chaos-runner
      ADD ./litmus-go /litmus-go
      WORKDIR /chaos-runner
      RUN go mod edit -replace github.com/litmuschaos/litmus-go=../litmus-go
      RUN go mod tidy
      (3) 도커 이미지 생성
- otel-litmus(240915) 디렉토리에서 진행
- docker build -t suhyen/chaos-runner:v01 -f chaos-runner/build/Dockerfile .
- docker push suhyen/chaos-runner:v01

3. chaos-operator
(1) 설치
- git clone -b distributed-tracing --single-branch https://github.com/namkyu1999/chaos-operator.git
(2) 수정
- pkg/telemetry/otel.go #18
  - **변경 전**: const OTELExporterOTLPEndpoint = OTEL_EXPORTER_OTLP_ENDPOINT
  - **변경 후**: const OTELExporterOTLPEndpoint = "otel-collector.observability.svc.cluster.local:4317"
  (3) 도커 이미지 생성
- chaos-operator 디렉토리에서 진행
- docker build -t suhyen/chaos-operator:v01 -f build/Dockerfile .
- docker push suhyen/chaos-operator:v01


## 4. Litmus
1. Litmus 설치 및 접속
(1) helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/
(2) helm repo list
(3) kubectl create ns litmus
(4) helm install chaos litmuschaos/litmus \
--namespace=litmus \
--set portal.frontend.service.type=NodePort \
--set mongodb.image.registry=ghcr.io/zcube \
--set mongodb.image.repository=bitnami-compat/mongodb \
--set mongodb.image.tag=6.0.5  
(5) kubectl get all -n litmus 
(6) minikube service chaos-litmus-frontend-service -n litmus --url
    (admin/litmus) -> (admin/Litmus1@)

2. Litmus Environments 설정
(1) + New Environment 
- Environment Name: local
- Environment Type: Production 
(2) + Enable Chaos
- Name: local
- Chaos Components Installation: Cluster-wide access
- Installation Location(Namespace): litmus 
- Service Account Name: litmus 
- Kubernetes Setup Instructions: Download
(3) local-litmus-chaos-enable.yml 수정
- local-litmus-chaos-enable.yml 다운로드 위치로 이동  
- #652: subscriber server addr 변경
  - kubectl get service chaos-litmus-frontend-service -n litmus
  - **변경 전**: SERVER_ADDR: http://127.0.0.1:54772/api/query
  - **변경 후**: SERVER_ADDR: http://10.106.97.98:9091/api/query
- #1043: chaos-operator 이미지 변경
  - **변경 전**: image: litmuschaos.docker.scarf.sh/litmuschaos/chaos-operator:3.10.0
  - **변경 후**: image: suhyen/chaos-operator:v01
- #1054: chaos-runner 이미지 변경
  - **변경 전**: value: litmuschaos.docker.scarf.sh/litmuschaos/chaos-runner:3.10.0
  - **변경 후**: value: suhyen/chaos-runner:v01
- kubectl apply -f local-litmus-chaos-enable.yml
- kubectl get all -n litmus
- CONNECTED 뜰 때까지 기다리기

3. Litmus Resilience Probe 설정
(1) + New Probe (Command)  
- Name: nginx-probe
- Timeout: 3s
- Interval: 3s 
- Attempt: 3
- Command: kubectl get pods --all-namespaces | grep nginx | grep Running | wc -l 
- Type: Int 
- Comparison Criteria: > 
- Value: 0

4. Litmus Experiment 실행
- experiment name: pod-delete-1
- namespace: default
- component: deployment
- label: app=nginx
- probe: nginx-probe (EOT)


## 5. Otel-collector, jaeger, prometheus
(1) kubectl create ns observability
(2) kubectl apply -f init-setting.yaml
(3) kubectl get all -n observability
(4) minikube service jaeger -n observability --url



