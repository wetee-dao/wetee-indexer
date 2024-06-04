# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

tag=`date "+%Y-%m-%d-%H_%M"`

# 构建镜像
ego-go build -o ./bin/indexer ./server.go

docker build -t wetee/indexer:$tag .

docker push wetee/indexer:$tag

export WETEE_INDEXER_IMAGE=wetee/indexer:$tag
echo '' > ./hack/indexer.yaml
envsubst < ./hack/indexer-temp.yaml > ./hack/indexer.yaml

# # 部署镜像
kubectl create -f ./hack/indexer.yaml
