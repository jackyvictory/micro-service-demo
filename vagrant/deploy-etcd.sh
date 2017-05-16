#!/bin/bash
# http://www.itdadao.com/articles/c15a1255268p0.html
# to run the script: bash deploy-etcd.sh
set -x
set -e

#更改这里的IP, 只支持部署3个节点etcd集群
declare -A NODE_MAP=(["etcd0"]="192.168.99.40" ["etcd1"]="192.168.99.50" ["etcd2"]="192.168.99.60")

etcd::download()
{
    ETCD_VER=v3.1.6    #指定要安装的版本号
    DOWNLOAD_URL=https://github.com/coreos/etcd/releases/download
    if ! [ -f ${PWD}/etcd-${ETCD_VER}-linux-amd64.tar.gz ]; then
      curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o ${PWD}/etcd-${ETCD_VER}-linux-amd64.tar.gz
    fi
    mkdir -p ${PWD}/temp-etcd && tar xzvf ${PWD}/etcd-${ETCD_VER}-linux-amd64.tar.gz -C ${PWD}/temp-etcd --strip-components=1
}

etcd::config()
{
    local node_index=$1

cat <<EOF >${PWD}/${node_index}.conf
export ETCD_NAME=${node_index}
export ETCD_DATA_DIR="/opt/etcd/data"
export ETCD_INITIAL_ADVERTISE_PEER_URLS="http://${NODE_MAP[${node_index}]}:2380"
export ETCD_LISTEN_PEER_URLS="http://${NODE_MAP[${node_index}]}:2380"
export ETCD_LISTEN_CLIENT_URLS="http://${NODE_MAP[${node_index}]}:2379,http://127.0.0.1:2379"
export ETCD_ADVERTISE_CLIENT_URLS="http://${NODE_MAP[${node_index}]}:2379"
export ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster-378"
export ETCD_INITIAL_CLUSTER="etcd0=http://${NODE_MAP['etcd0']}:2380,etcd1=http://${NODE_MAP['etcd1']}:2380,etcd2=http://${NODE_MAP['etcd2']}:2380"
export ETCD_INITIAL_CLUSTER_STATE="new"
export ETCDCTL_API=3
# export ETCD_DISCOVERY=""
# export ETCD_DISCOVERY_SRV=""
# export ETCD_DISCOVERY_FALLBACK="proxy"
# export ETCD_DISCOVERY_PROXY=""
#
# export ETCD_CA_FILE=""
# export ETCD_CERT_FILE=""
# export ETCD_KEY_FILE=""
# export ETCD_PEER_CA_FILE=""
# export ETCD_PEER_CERT_FILE=""
# export ETCD_PEER_KEY_FILE=""
EOF
}

etcd::gen_unit()
{
cat <<EOF >${PWD}/etcd.service
[Unit]
Description=Etcd Server
After=network.target

[Service]
Type=notify
WorkingDirectory=/opt/etcd/data
EnvironmentFile=-/opt/etcd/conf/etcd.conf
ExecStart=/usr/bin/etcd
Restart=always
RestartSec=8s
LimitNOFILE=40000

[Install]
WantedBy=multi-user.target
EOF
}

SSH_OPTS="-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -oLogLevel=ERROR -C"
etcd::scp()
{
  local host="$1"
  local src="$2"
  local dst="$3"
  scp -r ${SSH_OPTS} ${src[*]} "${host}:${dst}"
}
etcd::ssh()
{
  local host="$1"
  shift
  ssh ${SSH_OPTS} -t "${host}" "$@" >/dev/null 2>&1
}
etcd::ssh_nowait()
{
  local host="$1"
  shift
  ssh ${SSH_OPTS} "${host}" "source /opt/etcd/conf/etcd.conf; nohup $@"
}

etcd::deploy()
{

    for key in ${!NODE_MAP[@]}
    do
        etcd::config $key
        etcd::ssh "root@${NODE_MAP[$key]}" "mkdir -p /opt/etcd/data /opt/etcd/conf /opt/etcd/log"
        # etcd::ssh "root@${NODE_MAP[$key]}" "pkill etcd"
        etcd::scp "root@${NODE_MAP[$key]}" "${key}.conf" "/opt/etcd/conf/etcd.conf"
        # etcd::scp "root@${NODE_MAP[$key]}" "etcd.service" "/usr/lib/systemd/system"
        etcd::scp "root@${NODE_MAP[$key]}" "${PWD}/temp-etcd/etcd ${PWD}/temp-etcd/etcdctl" "/usr/bin"
        etcd::ssh "root@${NODE_MAP[$key]}" "chmod 755 /usr/bin/etcd*"
        etcd::ssh "root@${NODE_MAP[$key]}" "chown vagrant:vagrant /usr/bin/etcd*"
        etcd::ssh "root@${NODE_MAP[$key]}" "chown -R vagrant:vagrant /opt/etcd"
        # etcd::ssh_nowait "root@${NODE_MAP[$key]}" "systemctl daemon-reload && systemctl enable etcd && nohup systemctl start etcd"
        etcd::ssh_nowait "vagrant@${NODE_MAP[$key]}" "etcd > /opt/etcd/log/etcd.log 2>&1 &"
    done

}

etcd::clean()
{
  for key in ${!NODE_MAP[@]}
  do
    rm -f ${PWD}/${key}.conf
  done

  rm -fr ${PWD}/temp-etcd

  rm -f ${PWD}/etcd.service
}


etcd::download
etcd::gen_unit
etcd::deploy
etcd::clean

echo -e "部署完毕, 执行 etcdctl cluster-health ,检测是否OK"
