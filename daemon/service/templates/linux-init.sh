#!/bin/bash
#
# pmzero       Startup script for the pmzero
#
# description: pmzero is the easiest process manager
# processname: pmzero
# config: /etc/pmzero/pmzero.conf
# config: /etc/pmzero/envs
# pidfile: /var/run/pmzero.pid

SCRIPT="/usr/bin/pmzero daemon"
RUNAS=root

PIDFILE=/var/run/pmzero.pid
LOGFILE=/var/log/pmzero.log

start() {
  if [ -f /var/run/$PIDNAME ] && kill -0 $(cat /var/run/$PIDNAME); then
    echo 'Service already running' >&2
    return 1
  fi
  echo 'Starting service…' >&2
  local CMD="$SCRIPT &> \"$LOGFILE\" & echo \$!"
  su -c "$CMD" $RUNAS > "$PIDFILE"
  echo 'Service started' >&2
}

stop() {
  if [ ! -f "$PIDFILE" ] || ! kill -0 $(cat "$PIDFILE"); then
    echo 'Service not running' >&2
    return 1
  fi
  echo 'Stopping service…' >&2
  kill -15 $(cat "$PIDFILE") && rm -f "$PIDFILE"
  echo 'Service stopped' >&2
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart)
    stop
    start
    ;;
  *)
    echo "Usage: $0 {start|stop|restart}"
esac
