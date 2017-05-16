# Build & dispatch
Build trans service, and dispatch it to vagrant virtual boxex

		$ ./packageAndDispatch.sh

# Run
Run trans service on virtual boxex

    $ cd ../../vagrant

    $ vagrant ssh app1
    $ cd /opt/service
    $ nohup ./transService >> transService.log 2>&1 &
    $ exit

    $ vagrant ssh app2
    $ cd /opt/service
    $ nohup ./transService >> transService.log 2>&1 &
    $ exit

    $ vagrant ssh app3
    $ cd /opt/service
    $ nohup ./transService >> transService.log 2>&1 &
    $ exit
