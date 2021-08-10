execDir=~/go/src/k8s.io/code-generator

"${execDir}"/generate-groups.sh all github.com/vince15dk/k8s-operator-nhncloud/pkg/client github.com/vince15dk/k8s-operator-nhncloud/pkg/apis nhncloud.com:v1beta1 --output-base "/Users/nhn/Desktop/Linux/Go/k8s-operator-nhncloud" --go-header-file "${execDir}"/hack/boilerplate.go.txt

controller-gen paths=github.com/vince15dk/k8s-operator-nhncloud/pkg/apis/nhncloud.com/v1beta1 crd:trivialVersions=true crd:crdVersions=v1 output:crd:artifacts:config=manifests

kl create secret generic nhn-token --from-literal=tenantId=12b2e9b9caac4247887fe8501492d6c0 --from-literal=userName=sukjoo.kim@nhn.com --from-literal=password='nhn!@%$0'