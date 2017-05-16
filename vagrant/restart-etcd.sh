set -x
set -e

#更改这里的IP, 只支持部署3个节点etcd集群
declare -A NODE_MAP=(["etcd0"]="192.168.99.40" ["etcd1"]="192.168.99.50" ["etcd2"]="192.168.99.60")

SSH_OPTS="-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null -oLogLevel=ERROR -C"

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
etcd::restart()
{

    for key in ${!NODE_MAP[@]}
    do
        etcd::ssh "vagrant@${NODE_MAP[$key]}" "pkill etcd"
        etcd::ssh_nowait "vagrant@${NODE_MAP[$key]}" "etcd > /opt/etcd/log/etcd.log 2>&1 &"
    done

}

etcd::restart

echo -e "启动完毕, 执行 etcdctl cluster-health ,检测是否OK"
