#!/bin/bash
#!/bin/bash #NOPRINT
export NAME="docker-parallelcoind" #NOPRINT
export DATADIR="`pwd`" #NOPRINT
# source $DATADIR/config #NOPRINT
echo "Loading command aliases... Type 'halp' to see available commands" #NOPRINT
### HALP! How to control your $NAME docker container
alias      dkr="sudo docker"
         ### [ shortcut to run docker with sudo ]
alias   .where="echo $DATADIR"
         ### [ show where the current instance activated by init.sh lives ]
alias      .cd="cd $DATADIR"
         ### [ change working directory to instance folder ]
alias     .run="sudo docker run --network=\"host\" -v $DATADIR/data:/root/.parallelcoin -d=true -p 11047:11047 -p 11048:11048 -p 21047:21047 -p 21048:21048 --device /dev/fuse --cap-add SYS_ADMIN --security-opt apparmor:unconfined --name $NAME $NAME"
         ### [ start up the container (after building, to restart. for a'.stop'ed container, use '.start') ]
alias   .start="sudo docker start $NAME"
         ### [ start the container that was previously '.stop'ed ]
alias    .stop="sudo docker stop $NAME"
         ### [ stop the container, start it again with '.start' ]
alias   .enter="sudo docker exec -it $NAME bash"
         ### [ open a shell inside the container ]
# alias     .log="sudo tail -f $DATADIR/data/steemd.log"
         ### [ show the current output from the primary process in the container ]
alias   .build="sudo docker build -t $NAME $DATADIR"
         ### [ build the container from the Dockerfile ]
alias      .rm="sudo docker rm $NAME"
         ### [ remove the current container (for rebuilding) ]
alias .editdkr="nano $DATADIR/Dockerfile"
         ### [ edit the Dockerfile ]
alias  .editsh="nano $DATADIR/init.sh;source $DATADIR/init.sh"
         ### [ edit init.sh with nano then reload ]
alias  halp="sed 's/\$NAME/$NAME/g' $DATADIR/init.sh|sed 's#\$DATADIR#$DATADIR#g'|grep -v NOPRINT|sed 's/alias //g'|sed 's/=\"/     \"/g'|sed 's/#/>/g'"
######### hit the 'q' key to exit help viewer <<<<<<<<<