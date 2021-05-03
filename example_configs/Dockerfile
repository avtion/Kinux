FROM centos:7.9.2009
LABEL maintainer="广东第二师范学院计算机学院"

RUN [ -e /etc/yum.conf ] && sed -i '/tsflags=nodocs/d' /etc/yum.conf || true

# 重新安装所有软件用于完整man信息
RUN yum -y reinstall "*" && yum clean all

RUN yum makecache

# 安装软件
RUN yum -y install \
            kernel \
            e2fsprogs \
            sos \
            crash \
            strace \
            ltrace \
            tcpdump \
            abrt \
            pcp \
            systemtap \
            perf \
            bc \
            blktrace \
            btrfs-progs \
            ethtool \
            file \
            findutils \
            gcc \
            gdb \
            git \
            glibc-common \
            glibc-utils \
            hwloc \
            iotop \
            iproute \
            iputils \
            less \
            pciutils \
            ltrace \
            mailx \
            man-db \
            nc \
            netsniff-ng \
            net-tools \
            numactl \
            numactl-devel \
            passwd \
            perf \
            procps-ng \
            psmisc \
            screen \
            strace \
            sysstat \
            systemtap-client \
            tar \
            tcpdump \
            vim-enhanced \
            xauth \
            which \
            ostree \
            rpm-ostree \
            gdb-gdbserver \
            vim-minimal \
            bash-completion \
            subscription-manager \
            python-rhsm \
            rootfiles \
            yum-utils \
            nano \
            mysql \
            python3 \
            nmap \
            lsof \
            bash-completion \
            nodejs \
            && yum clean all

RUN adduser kinux --uid 1000
USER kinux

CMD ["tail", "-f", "/dev/null"]