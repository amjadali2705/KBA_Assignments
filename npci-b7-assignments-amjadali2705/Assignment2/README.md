# Assignment 2

Developed a Hyperledger Fabric blockchain solution for managing vehicle insurance claims post-accident. 

# Hyperledger Fabric network using minifabric

To build a network using minifabric, first we have to create a spec.yaml.

After that we have to perform following commands:

Step 1: Start the network
```
minifab netup -s couchdb -e true -i 2.4.8 -o police.insuranceclaim.com
```

Step 2: Create the Channel
```
minifab create -c insuranceclaimchannel
```

Step 3: Join the peers to the channel
```
minifab join -c insuranceclaimchannel
```


Step 4: Anchor update
```
minifab anchorupdate
```

Step 5: Profile Generation
```
minifab profilegen
```




