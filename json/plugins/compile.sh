echo "构建 ./$1.so"

go build -gcflags "all=-N -l" -ldflags "-s -w" -o ./"$1"_tmp.so -buildmode=plugin ./"$1".go

size() {
  stat -c %s "$1" | tr -d '\n'
}
# shellcheck disable=SC2046
echo "构建完成，动态链接库大小" + `size "./$1_tmp.so"`

echo "执行压缩命令"
upx ./"$1"_tmp.so -o ./"$1".so

rm ./"$1"_tmp.so

echo "压缩完成，动态链接库大小" + `size "./"$1".so"`