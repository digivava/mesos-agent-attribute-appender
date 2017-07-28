FROM alpine:3.4

COPY bin/mesos-slave-attribute-appender /bin/mesos-slave-attribute-appender

EXPOSE 19001

CMD ["/bin/mesos-slave-attribute-appender"]
