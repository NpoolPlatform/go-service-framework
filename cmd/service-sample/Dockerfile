FROM centos:7

COPY ServiceSample.viper.yaml /etc/ServiceSample/ServiceSample.viper.yaml
COPY service-sample /usr/local/bin/service-sample

CMD ["service-sample", "run"]
