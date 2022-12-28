TAG=${1:-latest}

echo "TAG=$TAG"

echo "docker build -t 10.30.30.22:9080/isyscore/engine/isc-json-engine:$TAG ."

docker build -t 10.30.30.36:9080/isyscore/engine/isc-json-engine:$TAG .

echo "docker push 10.30.30.22:9080/isyscore/engine/isc-json-engine:$TAG"

docker push 10.30.30.36:9080/isyscore/engine/isc-json-engine:$TAG

echo "push success"
