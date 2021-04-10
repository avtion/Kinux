FROM centos/tools
LABEL maintainer="广东第二师范学院计算机学院"
# 拷贝上海时间
ADD /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 使用清华大学镜像源
RUN sed -e 's|^mirrorlist=|#mirrorlist=|g' \
             -e 's|^#baseurl=http://mirror.centos.org|baseurl=https://mirrors.tuna.tsinghua.edu.cn|g' \
             -i.bak \
             /etc/yum.repos.d/CentOS-*.repo
RUN yum makecache

# 安装应用程序
RUN yum -y install nano \
                   mysql \
                   python3 \
                   nmap \
                   lsof \
                   bash-completion

CMD ["tail", "-f", "/dev/null"]