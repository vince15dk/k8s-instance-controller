execDir=~/go/src/k8s.io/code-generator

"${execDir}"/generate-groups.sh all github.com/vince15dk/k8s-operator-nhncloud/pkg/client github.com/vince15dk/k8s-operator-nhncloud/pkg/apis nhncloud.com:v1beta1 --output-base "/Users/nhn/Desktop/Linux/Go/k8s-operator-nhncloud" --go-header-file "${execDir}"/hack/boilerplate.go.txt

controller-gen paths=github.com/vince15dk/k8s-operator-nhncloud/pkg/apis/nhncloud.com/v1beta1 crd:trivialVersions=true crd:crdVersions=v1 output:crd:artifacts:config=manifests