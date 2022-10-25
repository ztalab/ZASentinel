### Thank you for your interest in ZASentinel

ZASentinel helps organizations improve information security by providing a better and simpler way to protect access to private network resources and applications.

The information contained in this document is intended to provide customers and potential customers with transparency and confidence about the security situation, practices and processes of ZASentinel. Although this document represents the security of ZASentinel when it was last updated, ensuring information security is a process of continuous improvement, and we will not stop the evaluation and support.

This document is divided into two main parts. The first part describes the security architecture of ZASentinel from the perspective of products, and the second part describes the general information security practices of ZASentinel.

### Zero Trust Access (ZTA)

ZASentinel's architecture is designed according to ZTA's main philosophy:

- Every request for private resources should be authenticated and authorized based on real-time context information, such as user identity, device status and other attributes such as time and location.
- Users should only have access to the minimum necessary resources (principle of least privilege).
- Events should be recorded in detail to support monitoring, detection, investigation and other analysis.

Generally speaking, these principles are embodied in our products in the following ways:

- No ZASentinel component can independently decide to allow traffic to flow to secure resources on remote networks. And user data flow authorization is always confirmed by running multiple checks by multiple components. In addition, the user data flow and the user authentication flow are handled by separate components and require separate authentication checks.
- We entrust user authentication to a third-party authentication provider, which creates a separation of concerns and provides an additional layer of security.
- We support extensive logging, provide enhanced visibility, and help administrators monitor, troubleshoot, and investigate the activities of the entire network and terminals, including unauthorized activities, as well as those of individuals who are authorized to access resources but may exceed their business authority.
- User data streams from the client to the server are encrypted end-to-end. Even though they can pass through relays (components residing in the infrastructure controlled by ZASentinel), ZASentinel cannot decrypt such data. (Encrypted transmission will never have intermediate termination/recovery on the relay. )
- Our client-relay-server, all components communicate with each other by upgrading to websocket protocol, which better supports the transfer of custom parameters and facilitates the verification of data security and legality.
- Our relay is mounted under the CDN network, which can effectively prevent DDOS attacks, and better optimize network lines to improve access speed. The network CDN supports Websocket protocol and seamlessly integrates with ZASentinel network.
- Our client application is designed to always run in the background.

From the customer's point of view, our product architecture provides additional security advantages in the following aspects:

- Hide the network -ZASentinel will not publicly expose the access point to your network to promote secure remote access to resources inside the network. This means that, if necessary, the customer network can remain invisible to the public Internet, and it will not be explored by potential attackers.
- Support granular access control at port level -Compared with VPN that provides access to the whole network, this helps to reduce the explosion radius of intrusion.
- Centralized access management for all private resources  -This helps IT team to audit and maintain access list to cope with personnel changes.
- Record all network activities centrally -Identity index log and analysis can provide insight into who is doing what in the whole enterprise network.
- Better availability for administrators -Solutions that are easier to configure and maintain are less prone to errors and more effective.
- Provide better usability for end users -This is an easy-to-deploy, always-on solution, which does not hinder users, promotes the adoption of VPN clients, and avoids common problems of VPN clients, where users shut down their internet connections because they interrupt or slow down.
- Provide a safe and fast network line -As a public basic service, the relay is mounted under the CDN network to optimize the network access line and effectively prevent DDOS attacks.

### Architecture & Component Overview

This section provides a high-level overview of the main components that make up the ZASentinel architecture and how they interact. The architecture of ZASentinel is in our document, and we strongly recommend that you read this article in order to have a more comprehensive understanding of the security foundation of ZASentinel.

![en-1](https://user-images.githubusercontent.com/52234994/165201473-fbb91967-269d-4986-84d1-df6e0e5775e0.png)

ZASentinel protects access to customers' remote network resources. ZASentinel consists of four main components, which together ensure that only authenticated users can access the resources they have access to. After ZASentinel is fully configured, the end result is that authorized users can connect to any resource without knowing the underlying network configuration or even which remote network the resource resides on. The four main components are:

- Controller-central coordination component of -ZASentinel. The controller is a multi-tenant component operated by ZASentinel. It performs a number of duties, including storing the client configuration information management console managed based on web to register and verify the server, issuing the signature authorization to the client that successfully sends the connection request, entrusting the user authentication to the identity provider, and distributing the signed access control list to the client and the server.
- Client-A software application installed on the end user's device. For users' requests for protected resources, the client acts as a combined authentication and authorization agent through ZASentinel. Routing and authorization decisions take place at the edge of the client.
- Server-A software component designed to be deployed on a device behind a remote network firewall. The server only starts the outbound connection with the controller and relay, and establishes communication with the client through these connections.
- Relay-geographically distributed components operated by ZASentinel. The relay serves as the registration point of the server and the connection point of the client wishing to establish a connection to the server. The relay is mounted under the CDN network to resist DDOS attacks.

Client and server components are located on the equipment and infrastructure controlled by customers, while controller and relay components are located in the infrastructure controlled by ZASentinel.

The following diagram illustrates how various components exchange information in a secure way, so as to connect from the client to a specific resource on a remote network:

<img src="https://user-images.githubusercontent.com/52234994/165201495-4125f8fa-381f-4fd9-89c8-31bce64e15bb.png" alt="image-20220402145306066" style="zoom: 60%;" />

1.The server registers itself with the geographically nearest relay.

The relay does not receive any information about the server except the randomly generated ID and its signature signed by the controller. Only when the signature of the server is signed by the same controller, the relay allows the server to register.

2.The client and the server respectively receive the traffic forwarding permission list.

Allow lists to be specific to each component. The client's allow list corresponds to the content that users are allowed to access, and the server's list covers the resources that administrators have configured. Two allow lists must be signed by the same controller.

3.The client is authenticated by a third-party identity provider, providing additional in-depth protection.
<img src="https://user-images.githubusercontent.com/52234994/165201503-f322fe86-9c3c-404f-a6f4-0a2a3ee5d896.png" alt="image-20220402145100815" style="zoom:60%;" />

4.The client initiates a TLS connection with a single end-to-end certificate lock to the requested server.

The relay only promotes this connection, but can't "see" any such data flow.

5.The server verifies whether the client request is signed by the same controller.

6.The client verifies whether the server signature matches the signature provided by the controller.

7.The server verifies whether the target address is in its allowed list.

8.Once established, the traffic flows to the destination through the encrypted TLS tunnel. DNS lookup and routing are forwarded from the client and executed by the server on the target network.

9.The relay can be mounted under the CDN network
The relay is mounted on the CDN network, and the CDN supports websocket protocol, which better supports the transmission of Header parameters. In terms of security, the CDN can naturally resist DDOS attacks, and effectively optimize network lines and increase traffic communication speed.

### Project analysis

This application is divided into three components: client, relay and server. The application relies on its own certificate to start and distinguishes different types according to the certificate extension field:

**Client:**

A client application (or client for short) is a software component installed on the user's device. The role of the client is to act as a combined authentication and authorization agent for users' requests for private resources. The client is responsible for establishing an encrypted tunnel with the corresponding server to access protected resources.

The following are the certificate extension fields:

```json
{
  "attrs": {
    "type": "client", 						
    "port": 48080,								
    "uuid":"a0b9238",
    "name":"client1", 
    "relays": [										
      {
        "uuid": "dcc509a",
        "name": "relay1",
        "addr": "relay.zsnb.xyz", 
        "port": 443,			 				
        "sort":1
      }
    ],
    "server": {
      "uuid": "879dea2",
      "name": "server1",
      "addr": "server.zsnb.xyz", 	
      "out_port": 443	   					
    },
   "target": {
      "host": "127.0.0.1", 		
      "port": 3306 								
    }
  }
}
```

**Server:**

The server is deployed behind the firewall of the dedicated remote network.

The responsibilities of the server are mainly:

- **Receiving client/relay connections**. Receive client requests and verify client access resources.
- **Proxy access resources**. Verify that the client requests resources, and the agent accesses private network resources.

The following are the certificate extension fields:

```json
{
  "attrs": {
    "type": "server", 				
    "host":"server.zsnb.xyz", 
    "port": 5091,							
    "out_port": 443,					
    "name": "server1",				
    "uuid": "server1",
    "resources": [					
      {
        "uuid": "879dea2", 
        "name": "mysql",			
        "type": "cidr",				
        "host": "127.0.0.1/16", 
        "port": 3306					
      }
    ]
  }
}

```

**Relay:**

Relay is the simplest component of the Zero Access architecture. It is mainly to receive communication policies, ACLs, etc. issued under the control surface to verify client requests.

The basic responsibilities of Relay are:

- As the registration point of the server. When initializing the server, the server information is registered to the relay, which allows the client to connect to the appropriate server, which contains any private network or other specific information shared.
- As the relay server of the client connection. Verify the legality of the request from the client and the communication policy issued under the control surface. After verification, the traffic will be transmitted.
- As a client connected to the server. After verifying that the client request is legal, relay allows the client to be connected directly to the requested server to establish a fixed communication tunnel.

The following are the certificate extension fields:

```json
{
  "attrs": {
    "type": "relay", 				
    "address":"relay.com", 	
    "port": 5091,						
    "name":"relay1",
    "uuid":"a0b9238",
  }
}
```

### Building

```
$ git clone git@github.com:ztalab/ZASentinel.git
$ cd ZASentinel
$ make release
```

You can set GOOS and GOARCH environment variables to allow Go to cross-compile alternative platforms.

The resulting binaries will be in the bin folder:

```
$ tree bin
bin
├── backend
```

Configure the certificate in the configuration file (configs/config.toml) or the environment variable to start.

```
bin/backend -c configs/config.toml help
```
