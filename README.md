[![pipeline status](https://api.travis-ci.org/bityuan/bityuan.svg?branch=master)](https://travis-ci.org/bityuan/bityuan/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bityuan/bityuan)](https://goreportcard.com/report/github.com/bityuan/bityuan)

# 敢么平行公链系统

#### 编译

```
git clone https://github.com/GM-Publicchain/gm $GOPATH/src/github.com/GM-Publicchain/gm
cd $GOPATH/src/github.com/GM-Publicchain/gm
go build -i -o gm
go build -i -o gm-cli github.com/GM-Publicchain/gm/cli
```

#### 运行

拷贝编译好的gm, gm-cli, gm.toml这三个文件置于同一个文件夹下，执行：

```
./gm -f gm.toml
```


