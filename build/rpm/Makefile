# 通过容器构建各个系统各个版本的rpm包
VERSION ?= $(shell git describe --tags --always --dirty=".dev")
# ================ go版本配置 ================
GO_VERSION ?= 1.15
GO_BASE_IMAGE ?= golang
GO_IMAGE ?= $(GO_BASE_IMAGE):$(GO_VERSION)

# 将项目打包的tgz文件放入rpmbuild/SOURCES
tgz ?= mkdir -p rpmbuild/SOURCES  && if [ ! -d "../tgz" ]; then echo tgz文件不存在创建tgz包;$(MAKE) -C ../ tgz && cp -f ../tgz/*tar.gz rpmbuild/SOURCES;fi

# 根据各个系统构建编译环境的容器
BUILD ?= DOCKER_BUILDKIT=1 \
	docker build \
	$(BUILD_IMAGE_FLAG) \
	--build-arg GO_IMAGE=$(GO_IMAGE) \
	-t rpmbuild-$@ \
	-f $@/Dockerfile \
	.

# ================ 配置构建的specs文件 ================
SPEC_FILES ?= http2socks.spec 
SPECS ?= $(addprefix SPECS/, $(SPEC_FILES))
# 在各个系统
RPMBUILD_FLAG ?= -ba \
	--define '_version ${VERSION}' \
	$(SPECS)
# 在容器里运行rpmbuild打包生成rpm文件
RUN ?= docker run \
	-v $(CURDIR)/rpmbuild/RPMS:/root/rpmbuild/RPMS \
	-v $(CURDIR)/rpmbuild/SRPMS:/root/rpmbuild/SRPMS \
	-v $(CURDIR)/rpmbuild/SOURCES:/root/rpmbuild/SOURCES \
	rpmbuild-$@ $(RPMBUILD_FLAG)

# ================ 目标操作系统配置 ================
CENTOS_RELEASES ?= centos-8
DISTROS := $(CENTOS_RELEASES)

.PHONY: help
help: ## 显示make的目标
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf " \033[36m%-20s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## 删除rpmbuild的包和中间产生的文件
	rm -rf rpmbuild

.PHONY: $(DISTROS)
$(DISTROS):
	@echo "================ 构建$@ ================"
	$(tgz)
	$(BUILD)
	$(RUN)
	@echo "================ 构建$@完成 ================"

.PHONY: rpm
rpm: centos ## 构建使用rpm系统的包
	
.PHONY: centos
centos: $(CENTOS_RELEASES) ## 构建centos的rpm包

.PHONY: centos-8
centos-8: ## 构建centos8的rpm包
