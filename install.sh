#! /bin/sh
# golang project installer  by usher.yue  for  ops doubi~.....
# golang environment installed to  current parent's  directory .
# example: install.sh linux32/linux64/darwin32/darwin64
if [ $# == 0 ];then
  echo  'please input platform..... '
  exit 2
fi 
case $1 in  
   "linux32")
		GOPACKAGE_LINK="http://golangtc.com/static/go/1.6/go1.6.linux-386.tar.gz"
		PKGNAME=gopkg_linux32
	;;
	"linux64")
		GOPACKAGE_LINK="http://www.golangtc.com/static/go/1.6/go1.6.linux-amd64.tar.gz"
		PKGNAME=gopkg_linux64
	;;
	"darwin64")
		GOPACKAGE_LINK="http://golangtc.com/static/go/1.6/go1.6.darwin-amd64.tar.gz"
		PKGNAME=gopkg_darwin64
	;;
	*)
	echo "os platform not  recognize ..... "
	exit 2 
	;;
esac
echo $GOPACKAGE_LINK
echo $PKGNAME
TMP=$2
GOPATH=$3
GOROOT=$4
INSTALL=$5
GOINSTALLROOT=${GOROOT}bin
echo  'start installing golang package project......'
#此处是目录创建
if [ ! -d "./tmp" ];then
echo "Create Download Floder On  Current Directory....."
mkdir -p $TMP 
fi
echo "begin download golang package..."
if [ ! -d $GOROOT ];then
echo "正在创建$GOROOT......"
mkdir -p $GOROOT
fi
if [ ! -d $GOPATH ];then
echo "正在创建$GOPATH......"
mkdir -p $GOPATH
fi
echo "generate object file directory....."
mkdir -p "${GOPATH}/src"
mkdir -p "${GOPATH}/bin"
mkdir -p "${GOPATH}/pkg"
#download file
if [ ! -f ${TMP}${PKGNAME}  ];then
wget -P $TMP $GOPACKAGE_LINK -O ${PKGNAME} && mv ${PKGNAME}   $TMP
fi
#unzip
tar zxvf ${TMP}${PKGNAME}  -C $INSTALL
#下面开始初始化环境变量
#export path
export GOPATH 
export GOROOT
export PATH=$PATH:$GOINSTALLROOT
echo $GOPATH 
echo $GOROOT
echo $PATH

#change work directory
sourcecode=$(pwd) && cp -fr $sourcecode ${GOPATH}src
#获取目录名字
directory=$(echo $sourcecode | awk -F "/" '{print $NF}')

#change work directory
echo $directory
cd ${GOPATH}src/${directory}
go build -a -v  -o ${directory}
cp ${directory} ${sourcecode}/${directory}
cd ${sourcecode}
echo "project build complete......"