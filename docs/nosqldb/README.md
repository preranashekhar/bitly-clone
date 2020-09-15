### NoSQL DB

- This DB is used as lookup cache by the redirect server when looking up for a long url given a short url. This is an implementation a 5 node database. This DB is deployed as a 5 pod service on AWS EKS.

### Setting up the 5 pod service

This deployment requires pod to pod communication for data replication. There might be a way to establish this through hostname but I went ahead with pod IP communication. This might not be ideal in other scenarios, however this project did not expect us to handle
pod restarts and hence it was OK to use IP address communication. I modified our nosql implementation to communicate based on IP which is checked in at https://github.com/nguyensjsu/cmpe281-preranashekhar/tree/aws-deployment

- How to setup the initial configuration

From the `kubectl get pods --namespace nosql -o wide command`, we will get to know on which nodes the pods are deployed. I then ssh'ed into
the nodes from jumpbox, obtained the IPs of the PODs and ran a command like this


This will ensure that all pods have the required configuration to enable pod to pod communication


### Note about Persistent Volume

I tried a lot to get persistent volume working for my cluster. However it was quite challenging for below reasons.

 - The NoSql code puts both code files and data files in the current working directory. When we mount a persistent volume
from the pod's node onto the container's current working directory in the pod, the file gets overwritten. This is also
an anti-pattern when mounting volumes
 - Since, pod to node assignment is arbitrary we need to have 5 different persistent volumes and 5 different persistent
volume claims. This was challenging to manage as well.

I have attached ymls I used to try.

### Screenshots

![Screen Shot 2020-05-02 at 11 33 24 PM](https://user-images.githubusercontent.com/55044852/80907613-55916e00-8ccd-11ea-817a-155d1859d563.png)

