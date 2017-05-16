# Install vagrant

  To be continued...

# Start vagrant virtual boxex

    $ vagrant reload

# Setup NTP both for virtual boxex and localhost

  To be continued...

# etcd3
## Reference

    http://www.infoq.com/cn/articles/etcd-interpretation-application-scenario-implement-principle
    http://cizixs.com/2016/08/02/intro-to-etcd

## Install etcd3 cluster on virtual boxex

* Reference

      https://coreos.com/etcd/docs/latest/dl_build.html
      http://www.itdadao.com/articles/c15a1255268p0.html

* Installation & Startup

  Note: maybe you need to setup some dependencies before run deploy-etcd.sh, please help yourself...

      $ vagrant ssh app1
      $ cd /vagrant
      $ bash deploy-etcd.sh

* Check etcd status

      $ etcdctl -version
      $ etcdctl cluster-health

# Prometheus
Monitoring service,  ingest and record etcd's metrics

* Reference

      https://coreos.com/etcd/docs/latest/op-guide/monitoring.html
      https://prometheus.io/

* Installation & Startup

      $ vagrant ssh app1
      $ sudo mkdir /opt/prometheus
      $ sudo chown vagrant:vagrant /opt/prometheus
      $ cd /vagrant
      $ PROMETHEUS_VERSION="1.3.1"
      $ wget https://github.com/prometheus/prometheus/releases/download/v$PROMETHEUS_VERSION/prometheus-$PROMETHEUS_VERSION.linux-amd64.tar.gz -O /tmp/prometheus-$PROMETHEUS_VERSION.linux-amd64.tar.gz
      $ tar -xvzf /tmp/prometheus-$PROMETHEUS_VERSION.linux-amd64.tar.gz --directory /opt/prometheus --strip-components=1
      $ /opt/prometheus/prometheus -version
      $ cat > /opt/prometheus/test-etcd.yaml <<EOF
      global:
      scrape_interval: 10s
      scrape_configs:
      - job_name: test-etcd
        static_configs:
        - targets: ['192.168.99.40:2379', '192.168.99.50:2379', '192.168.99.60:2379']
      EOF
      $ cat /opt/prometheus/test-etcd.yaml
      $ nohup /opt/prometheus/prometheus \
        -config.file /opt/prometheus/test-etcd.yaml \
        -web.listen-address ":9090" \
      -storage.local.path "test-etcd.data" >> /opt/prometheus/test-etcd.log  2>&1 &

* Check prometheus status

      http://192.168.99.40:9090/

# Grafana
The analytics platform for all your metrics

* Reference
      https://coreos.com/etcd/docs/latest/op-guide/monitoring.html
      https://grafana.com/grafana/download
      http://docs.grafana.org/installation/debian/

* Installation & Startup

      $ vagrant ssh app1
      $ cd /vagrant
      $ wget https://s3-us-west-2.amazonaws.com/grafana-releases/release/grafana_4.2.0_amd64.deb 
      $ sudo dpkg -i grafana_4.2.0_amd64.deb
      $ sudo service grafana-server start

* Check grafana status

    default user / pwd: admin / *****

      http://192.168.99.40:3000

* Configuration for grafana to monitor etcd cluster by metrics from Prometheus

    *Data Sources* - *Add data source*

      Name: test-etcd
      Type: Prometheus
      URL: http://192.168.99.40:9090
      Access: proxy

    *Dashboards* - *Import*

    paste Json: [default etcd dashboard template]( https://coreos.com/etcd/docs/latest/op-guide/grafana.json)
