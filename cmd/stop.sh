    	#!/bin/sh
    	PID=`ps -ef | grep chatsvr | grep -v grep | awk '{print $2}'`
    	if [ "" != "$PID" ]; then
    	 echo "killing $PID"
    	 kill -9 $PID
    	else
    	 echo "chatsvr not running!"
    	fi
