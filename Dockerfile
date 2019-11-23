FROM centos:centos7

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone

COPY docker/go-monitor /var/docker/go-monitor
RUN mkdir -p /var/docker/go-monitor/log

WORKDIR /var/docker/go-monitor
CMD [ "bin/monitor", "-c", "configs/monitor.json" ]
