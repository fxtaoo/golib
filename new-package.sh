#!/usr/bin/env bash
# 创建新 package

while true ; do
  read -rp "package 名称：" package_name
  if [[ ! -d $package_name ]] ; then
    break;
  fi
  echo -e "重新输入 $package_name 该 package 以存在!\n"
done


read -rp "package 信息：" intro

content="//t $intro
package $package_name

"

# 创建目录与文件
(mkdir ./$package_name && cd $package_name && \
echo -e "$content" > ./${package_name}.go
echo -e "$content" > ./${package_name}_test.go
echo -e "# ${package_name} \n ${intro}" > ./README.md)

echo -e "\n[${package_name}](./${package_name}/README.md) ${intro}">> ./README.md




