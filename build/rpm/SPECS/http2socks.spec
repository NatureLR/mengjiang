# 构建rpm包的文件
%global debug_package %{nil}

Name:           http2socks
Version:        %{_version}
Release:        1%{?dist}
Summary:        用于创建GO项目的脚手架

Group:          Application/WebServer
License:        Apache 2.0
URL:            http://www.baidu.com
Source0:        %{name}.tar.gz

# 构建依赖
BuildRequires:  git
BuildRequires:  make

# 详细描述
%description



# 构建之前执行的脚本，一般为解压缩
%prep

# %setup 不加任何选项，仅将软件包打开。
# %setup -a 切换目录前，解压指定 Source 文件，例如 "-a 0" 表示解压 "Source0"
# %setup -n newdir 将软件包解压在newdir目录。
# %setup -c 解压缩之前先产生目录。
# %setup -b num 将第 num 个 source 文件解压缩。
# %setup -D 解压前不删除目录
# %setup -T 不使用default的解压缩操作。
# %setup -q 不显示解包过程
# %setup -T -b 0 将第 0 个源代码文件解压缩。
# %setup -c -n newdir 指定目录名称 newdir，并在此目录产生 rpm 套件。
%setup -q -c -n src -a 0

# 编译脚本
%build

cd http2socks && make -C build build

# 检查
%check

http2socks/bin/http2socks version

# 安装阶段需要做的
%install

install -D  -p  -m 0755 ${RPM_BUILD_DIR}/src/http2socks/bin/http2socks ${RPM_BUILD_ROOT}%{_bindir}/http2socks
install -D -m 0644 ${RPM_BUILD_DIR}/src/http2socks/build/systemd/http2socks.service ${RPM_BUILD_ROOT}%{_unitdir}/http2socks.service

# 说明%{buildroot}中那些文件和目录需要打包到rpm中
%files

%{_bindir}/http2socks
%{_unitdir}/http2socks.service

# 变更记录
%changelog
