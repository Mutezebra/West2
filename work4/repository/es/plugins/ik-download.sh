#!/bin/bash

cd "$(dirname "$0")" || exit

if [ -d "./ik" ]; then
  return 0
fi

# 下载ik分词器
url="https://github.com/infinilabs/analysis-ik/releases/download/v8.9.0/elasticsearch-analysis-ik-8.9.0.zip"

if ! wget $url; then
  echo "download failed"
  return 1
fi

if ! command -v unzip >/dev/null 2>&1; then
  echo "lack of command unzip"
  return 1
fi

# 解压
unzip elasticsearch-analysis-ik-8.9.0.zip -d ./ik/
return 0