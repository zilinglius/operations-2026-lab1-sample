# 选做 O1 使用。镜像里不重新编译：直接放入 CI 已经构建好的二进制，
# 这样镜像里跑的就是制品里的那个文件（build once, deploy many，L-05 第 11 页）。
#
# 架构：镜像是什么架构由放进去的二进制决定。Actions 的 ubuntu-24.04 runner 是 x64，发布的制品和镜像都是 linux/amd64；
# 在 arm64 主机（如 Apple 芯片的 Mac）上运行要加 --platform linux/amd64，由 Docker 模拟运行，较慢。
#
# 构建上下文是解压后的制品目录，而不是仓库根目录：
#   tar -xzf dist/enroll-conflict_1.0.0_linux_amd64.tar.gz -C /tmp
#   docker build -f Dockerfile \
#     --build-arg VERSION=1.0.0 --build-arg COMMIT=<提交号> --build-arg SOURCE=<仓库地址> \
#     -t enroll-conflict:1.0.0 /tmp/enroll-conflict_1.0.0_linux_amd64
FROM scratch

ARG VERSION=dev
ARG COMMIT=unknown
ARG SOURCE=

# source 标签让 GitHub Container Registry 把镜像关联到对应仓库。
LABEL org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.revision=$COMMIT \
      org.opencontainers.image.source=$SOURCE

COPY enroll-conflict /enroll-conflict
COPY BUILD_INFO /BUILD_INFO
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/enroll-conflict"]
