# Geth as a Service

 - # To run this service

 - - Just type this prefix with file names in order to command line ``` (microk8s) kubectl apply -f (file names in order)```

 - - You must wait the jobs to finish its job before apply next file.

 - - To monitor Job's or the other components' status >>  ``` kubectl -n geth get all ```

 - - Congrats! You have a geth as a service.

 - # To monitor ethstats-poa

 - - ``` http://10.152.183.131:3000 ``` in your local machine.

 - # To use blockchain from somewhere (Such as REMIX IDE)

 - - local web3 as ``` http://10.152.183.133:8545``` again in your local machine.

 - # If your a dev,

 - - Use the service geth ```10.152.183.132``` with ports ``` 8545,8546,8547 ``` to connect miner nodes.
 
 
