# Mesos Slave Attribute Appender

The **mesos-slave-attribute-appender** is for those of us who have torn our hair out writing veeeery loooong Prometheus queries because the [mesosphere/mesos-exporter](https://github.com/mesosphere/mesos_exporter) metrics don't contain the Mesos slave attributes.  Currently, the mesos-exporter only provides the option to export a separate metric with this information, which doesn't help with our long queries much.

That's where mesos-slave-attribute-appender comes in! Just run it on all your slave nodes (e.g. as a systemd unit), and set up your Prometheus to scrape this thing for your slave metrics instead. It appends the attributes from the `mesos-slave-common` file as new labels onto the end of each metric.

# Usage
You will need to pass it the following environment variables:

- `PRIVATE_IPV4`, the IP address of the Mesos slave.
- `MESOS_SLAVE_COMMON_PATH`, the location on the slave where the `mesos-slave-common` file (containing the slave attributes) can be found.
- `MESOS_EXPORTER_PORT`, the port that `mesosphere/mesos-exporter` is running on.
- `PORT`, the port that the mesos-slave-attribute-appender will run on (defaults to `19001`--make sure your mesos-exporter is not using the same port.)
