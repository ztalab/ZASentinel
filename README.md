# ZASentinel Whitepaper



## Thank you for your interest in ZASentinel

ZASentinel helps organizations improve information security by providing a better and easier way to secure access to private network resources and applications. However, in order to deliver on this commitment, we recognize that we must first keep our own security practices in order. Clients entrust us with the custody of their most sensitive and confidential information, so information security is at the heart of our business and a priority that is always taken seriously by everyone at ZASentinel.

The information contained in this document is intended to provide customers and potential customers with transparency and confidence regarding ZASentinel's security posture, practices and processes. While this document represents the security of ZASentinel at the time of the last update, ensuring information security is a process of continuous improvement and we do not stop evaluating and supporting.

This document is divided into two main parts. The first part describes ZASentinel's security architecture from a product perspective, and the second part describes ZASentinel's general information security practices.

If you have any unanswered questions, you can contact us using the contact details at the end of this document.

## About Us

Founded in 2021, ZASentinel's goal is to help organizations easily implement Zero Trust-based secure access solutions without sacrificing security, availability or performance. ZASentinel believes it should be possible to "work anywhere".

ZASentinel is operated by a global team of experienced professionals in engineering, product management, customer support and business operations with a long track record of successfully delivering leading enterprise technology solutions. Our team comes from Huawei, Yuer.com, Kidswant and other companies.

### Do you want to use ZASentinel?

ZASentinel is trusted by some of the world's most innovative and fastest-growing companies. Most of our customers face similar challenges: on the one hand, they strive to expand their teams quickly, and on the other hand, they have to deal with security and operational problems caused by relying on outdated or cumbersome technologies. Our customers come from all walks of life, including strictly regulated industries such as financial services, health care and legal services.

### What services do we provide?

ZASentinel provides a secure access platform, replacing traditional VPN with a modern identity-first network solution that integrates enterprise-level security and consumer-level user experience. It can be set up in less than 15 minutes and integrated with all major cloud and identity providers. ZASentinel helps enterprises switch to SASE architecture by linking each network event to an identity (user, device and service), thus providing enterprises with unparalleled control and visibility over the entire network activities.

ZASentinel is delivered as a software as a service (SaaS) product with downloadable software components installed on end users and other devices.

##  Product safety

### 1.1. Product architecture and design 

#### Overview of the method

ZASentinel was originally designed to provide a safe solution for modern "work everywhere" employees. The modern labor force is becoming more and more dispersed:

- More and more users work in remote locations such as home and public places;
- User devices are scattered in desktop and mobile devices, work release devices and personal devices;
- The application is being migrated internally to the cloud or provided by a third-party SaaS provider;
- Networks are increasingly cloud-based or managed by third parties.

Therefore, the traditional fixed boundary method based on IP address cannot meet today's security or availability needs. ZASentinel's answer is a method based on the concept of zero-trust access, which we call identity-first networks.

#### Identity-first network

The premise behind ZASentinel's identity first network model is to rethink the basic assumptions of the network to meet the needs of modern labor. The decentralization and abstraction of users, devices, applications and networks (for example, through infrastructure such as Terraform, that is, code frameworks) means that different methods are needed.

Our identity-first approach starts with asking a simple question: Should network requests be allowed to leave the device? And if so, whose identity should it depend on? Because there is no clearly authorized identity, the network connection will never allow access to your network, so there is no longer a question of who the network connection belongs to and why it is authorized.

This connection between network request initiation and identity makes it very easy to understand the activities within the network. The days of whitelisting IP addresses are gone forever.

Maintain complex subnet allocations and VLAN segments, or manually piece together network events across physical networks. With ZASentinel, you can provide users with the flexibility they need to work anywhere, while maintaining confidence that each network connection will be authorized and reviewed according to the user's identity.

#### Zero trust access(ZTA)

The architecture of ZASentinel is designed according to the main concepts of ZTA:

- Each request for private resources should be authenticated and authorized based on real-time contextual information, such as user identity, device status, and other attributes such as time and location.
- Users should only have access to the minimum necessary resources (minimum privilege principle).
- Incidents should be recorded in detail to support monitoring, detection, investigation and other analysis.

In summary, these principles are reflected in our products in the following ways:

- No ZASentinel component can independently determine the security resources that allow traffic to flow to remote networks. User and data flow authorizations are always confirmed by running multiple checks through multiple components. In addition, user data streams and user authentication streams are handled by separate components and require separate verification checks.
- We delegate user authentication to a third-party identity authentication provider (IdP), which creates a separation of concerns and provides an additional layer of security.
- We support extensive logging, provide enhanced visibility, and help administrators monitor, troubleshoot and investigate the activities of the entire network and terminal, including unauthorized activities, as well as the activities of individuals authorized to access resources that may exceed their business authority.
- User data streams from client to server are end-to-end encrypted, and ZASentinel cannot decrypt such data even if they can be relayed (components resident in the infrastructure controlled by ZASentinel). Encrypted transmission will never terminate/restore in the middle of the relay.)
- Our client-relay-server, between components, by upgrading to the websocket protocol communication, we better support custom parameter passing and facilitate the security and legitimacy verification of various data.
- Our relay terminal is mounted under the CDN network, which can effectively prevent DDOS attacks, better optimize network lines and improve access speed. CDN network supports the Websocket protocol and seamlessly integrates with ZASentinel networks.
- Our client application is designed to run in the background all the time.

From the customer's perspective, our product architecture provides additional security advantages in the following aspects:

- **Hide network** - ZASentinel does not expose access points to your network publicly to facilitate secure remote access to internal resources of the network. This means that if necessary, the customer network can remain invisible to the public Internet and will not be explored by potential attackers.
- **Support port-level fine access control** - This helps reduce the explosion radius of intrusion compared with VPNs that provide network-wide access.
- **Centralized access management of all private resources** - This helps IT teams review and maintain access lists for staffing changes
- ** Centrally record all network activities** - Identity index logs and analysis can gain an in-depth understanding of who is doing what the whole enterprise network.
- **Provide better availability for administrators** - Solutions that are easier to configure and maintain are less error-prone and more effective.
- **Better availability for end users** - This is an easy-to-deploy and always online solution that does not hinder users, promotes the adoption of VPN clients, and avoids common problems of VPN clients. In VPN clients, users will turn them off because of interrupting or slowing down their Internet connections.
- **Provide safe and fast network route** - As a public basic service, the relay is mounted under the CDN network to optimize the network access route and effectively prevent DDOS attacks.

#### Architecture & Component Overview

This section highly summarizes the main components that make up the ZASentinel architecture and how they interact. The architecture of ZASentinel is in our file, and we strongly recommend that you read this article to have a more comprehensive understanding of the security foundation of ZASentinel.

![cn-1](./images/en-1.png)

ZASentinel protects access to customers' remote network resources. ZASentinel consists of four main components that together ensure that only authenticated users can access the resources they have access to. After ZASentinel is fully configured, the final result is that authorized users can connect to any resource without knowing the underlying network configuration, or even knowing which remote network the resources reside on. The four main components are:

- Controller - the central coordination component of ZASentinel. The controller is a multi-tenant component operated by ZASentinel. It performs a number of responsibilities, including storing web-based customer configuration information, management console registration and verification server, issuing signature authorization to clients who successfully send connection requests, and delegating user verification to An identity provider, and a list of access control that distributes signatures to clients and servers.
- Client - Software applications installed on end-user devices. For users' requests for protected resources, the client acts as the authentication and authorization agent of the combination through ZASentinel. Network routing and authorization decisions occur on the edge of the client.
- Server - Software components designed to be deployed on devices behind the remote network firewall. The server only starts outbound connections to the controller and the repeater, through which communication is established with the client.
- Relay - a geographically distributed component operated by ZASentinel. Relay acts as the registration point of the server and the connection point of the client who wants to establish a connection to the server. The relay is mounted on the CDN network to resist DDOS attacks.

Client and server components are located on the equipment and infrastructure controlled by customers, and controllers and relay components are located in the infrastructure controlled by ZASentinel.

The following figure shows how various components exchange information in a secure way to facilitate connection from clients to specific resources on remote networks:

<img src="./images/en-2.png" alt="image-20220402145306066" style="zoom: 60%;" />

1. **The server registers itself to the geographically nearest relay. **

   Except for the randomly generated ID and its signature signed by the controller, the relay does not receive any information about the server. Relay only allows server registration if the server's signature is signed by the same controller.

2. **Clients and servers each receive a list of traffic forwarding permissions.** 

   Allow the list to be specific to each component. The allowable list of the client corresponds to the content allowed to be accessed by users, and the scope of the list of the server is the resources already configured by the administrator. Both allowable lists must be signed by the same controller.

3. **The client is authenticated by a third-party identity provider, providing additional depth protection.**
   <img src="./images/en-3.png" alt="image-20220402145100815" style="zoom:60%;" />

4. **The client initiates a single end-to-end certificate-locked TLS connection to the requested server. **

   Relay only promotes this connection, and cannot "see" any such data stream.

5. **The server verifies whether the client request is signed by the same controller.**

6. **The client verifies whether the server signature matches the signature provided by the controller.**

7. **The server verifies whether the destination address is in its allowable list. **

8. **Once established, the traffic flows to the destination through the encrypted TLS tunnel. **

   DNS lookup and routing are forwarded from the client and executed by the server on the target network.

9. **The relay can be mounted on the CDN network.**
   The relay mounts to the CDN network. CDN supports the websocket protocol and better supports Header parameter transfer. In terms of security, CDN can naturally resist DDOS attacks, and effectively optimizes network lines to increase traffic communication speed.

### 1.2. Customer data

#### What customer data do we process?

The main customer data types processed by ZASentinel are:

- User details (such as email addresses, names and group members, but excluding passwords, because ZASentinel delegates authentication to third-party identity providers);
- Infrastructure information (such as network details, resource details and access control lists); and
- Logs (such as crashes and bug reports for diagnosis and troubleshooting). The ZASentinel component also records events that allow customers to monitor user activity (for example, user login and token requests).

User network traffic to ZASentinel-protected resources can be transmitted encrypted through ZASentinel relay. The flow of this carrying data is relayed on an instantaneous basis. Relay does not store traffic or any network-identifiable information.

ZASentinel can also process customer data submitted by customers related to customer support requests. This may include configuration data, error logs, and other information that the customer decides to provide to ZASentinel and other information needed to diagnose technical problems.

#### Data Separation

Customer data is logically isolated according to the customer tenant ID in ZASentinel's system.

#### Data confidentiality

According to our contract with customers, ZASentinel regards customer data as customer confidential information. Ownership of customer data is retained by the customer.

### 1.3. Product performance and scalability

Service reliability is the core aspect of information security. This section introduces service reliability related to performance and scalability. For information on service reliability in usability, please refer to Infrastructure & Physical Security below.

We have designed our infrastructure and software to ensure that ZASentinel performs well, even when the usage of individual customers or our entire customer base increases. Our main ways to achieve this goal include:

- Eliminate the return/tumber problem - The traffic routed through ZASentinel adopts a more direct route, rather than routing all traffic through a central gateway away from the beginning and end, thus reducing user latency and organizational bandwidth usage. ZASentinel clients automatically and intelligently connect to ZASentinel controllers and repeaters, which provide the best performance according to the user's physical location at that time and the resources they need to connect.
- Support split tunnel - Organizes that any user traffic that is not routed through ZASentinel will completely bypass ZASentinel and be handled independently by the user device. This reduces the sending of unnecessary traffic through additional relay segments.
- Load balancing - ZASentinel handles multiple levels of load balancing. ZASentinel's controllers and relay terminals are distributed in different locations and geographical regions. As part of infrastructure planning, our goal is to strategically allocate them to reduce latency and provide load balancing in areas with high expected traffic loads. For example, in the high-flow area, we can add additional controllers and repeaters and balance the load between them. The delay is further reduced by using the same IaaS provider used by the customer to host controllers and relays (for example, in Alibaba Cloud, Huawei Cloud and GCP). On the client, customers can install multiple servers in the same remote network, and ZASentinel will automatically handle the load balancing between servers to meet access requests for specific networks.

- Processing extensions for customers - The traditional network access security model requires organizations to deploy and maintain their own security infrastructure, such as VPN gateways. Scaling up will disproportionately increase management costs and occupy resources that can be used for other programs. ZASentinel alleviated the IT department's concerns about expansion.
- Distributed authorization processing - Authorization processing workloads are distributed, such as at the ZASentinel client level, which helps improve overall performance rather than being concentrated in one location.

##  Information security program

### 2.1. Overview

This section contains a summary of the ZASentinel information security plan. ZASentinel retains a written set of information security plans, policies and procedures, which are reviewed at least once a year and supplemented by periodic risk reviews, which will be integrated into the continuous development of our information security plans.

### 2.2.  **Administrative & Organizational Security** 

#### **Governance & Responsibility** 

The chief technology officer of ZASentinel is mainly responsible for ZASentinel's information security plan. ZASentinel also has an interdisciplinary security team responsible for implementing, reviewing and maintaining its information security plan, including senior managers.

However, as a security project, ZASentinel believes that security is a common problem and therefore a shared responsibility of our entire organization. For example, all our engineers are required to regard safety as an essential part of their work, and they will not simply delegate all responsibilities to other colleagues who are more concerned about safety.

#### People  **Security** 

##### Employee Background Checks

All new employees are subject to background checks by third-party suppliers specializing in background checks. Background checks include crime and global watch list surveys. Conduct background checks for new employees in accordance with local practices and laws.

##### Employee confidentiality obligations

All employees and independent contractors must sign contracts to protect customer information and other proprietary information as confidential information. Our employee manual and security training emphasize the importance of maintaining customer data confidentiality.

##### Offboarding Employees

We have documented the process we followed to ensure that when employees leave the company, their access is cancelled in a timely manner, and any company assets they own are returned or destroyed (as appropriate).

##### Workforce safety training

All employees are required to receive information security training during the onboarding period. In addition, employees will receive regular safety review training to keep abreast of safety best practices. Ask employees to review and confirm that they have reviewed the company's information security policies and procedures, including acceptable IT resource use policies.

#### Vendor management

We work with suppliers and service providers that help us provide services to customers, such as email service providers, payment processors and infrastructure providers. Some of these suppliers have the ability to access or store customer data, and we have taken various measures to ensure that they only do so in an appropriate and safe way.

First of all, we choose suppliers based on our experience in working with suppliers, their reputation and their assessment of how much they meet our business needs. Secondly, we conduct due diligence on potential suppliers, including a risk assessment of their safety situation. Third, when we hire new suppliers, we ensure that they are signed in writing, including appropriate terms related to confidentiality, security, privacy and service levels. Generally speaking, suppliers are only allowed to use the data provided to them (including the data of our customers) for the purpose of providing services to ZASentinel.

#### **Incident Response Plans & Breach Notification** 

ZASentinel has an event response process that covers many types of events, including security and availability events. This is a series of documented processes involving multiple teams, such as security, engineering, customer support, legal, communication and management.

The accident response process covers all aspects from the initial response, investigation, notification, mitigation and remediation.

If there is a safety incident affecting the customer, we will notify the affected customer in accordance with our legal obligations. Please note that it is not always possible to notify the customer immediately in the event of a safety accident, because it sometimes takes time to conduct an appropriate investigation to determine what happened.

Customers can report any actual or suspicious incidents to us through our customer support channels.

#### Insurance

ZASentinel has a standard set of insurance policies, and we evaluate its coverage as suitable for our business. The policy includes general business liability, worker compensation, self (including network liability), and DGO insurance.

### 2.3. Application security

#### Data Protection & Access Control

##### Access control

We provide users' access to the system according to the user role and minimum permission principle.

Use ZASentinel to protect access to private network resources in production and other environments. Authentication is performed by the single sign-on system of our identity provider and multi-factor authentication is enabled. With ZASentinel, we can also finely control user access rights at the resource level (rather than the network level) according to the minimum permission principle, and apply different security policies according to the authentication of users, devices and the context of the requested access.

Internal enterprise applications use SSO and MFA for authentication whenever possible, and enforce minimum password complexity requirements.

We have also automated the production environment deployment process, which means that users do not have the right to change the production environment manually or directly. Developers cannot directly access the database containing customer data. Developers usually do not have or do not need to access the production environment server through SSH.

##### Access monitoring

We use ZASentinel and other log systems to monitor access to various systems and all aspects of the ZASentinel infrastructure.

##### CDN Network

We carefully select CDN service providers and test the security and access speed of relevant CDN services.

CDN networks can effectively resist DDOS attacks and effectively optimize network access lines.

CDN network supports the Websocket protocol and seamlessly integrates with ZASentinel networks.

##### Data encryption

Use industry-standard encryption protocols to encrypt in transit and static customer data.

During transmission, client application communication is protected by TLS/SSL (HTTPS) connection.

Staticly, customer data is stored in a database managed by Alibaba Cloud Platform, which uses AES-256 or better standard encryption and symmetric keys. The data key itself is encrypted with the master key stored in the secure key library and is changed regularly.

We do not use any customized or proprietary encryption frameworks or implementations. Please note that ZASentinel does not store any customer passwords.

##### Data backup

We automate daily backups of customer databases. For disaster recovery purposes, backups are stored for a limited period of time and tested regularly.

##### Data deletion

Permanently and securely delete customer data at the request of the customer and in accordance with any contractual commitments made to the customer.

#### Software Development Methodology & Testing

##### Software development

All software code written must be reviewed by a second person. In addition, ZASentinel also performs internal and third-party security tests, as described below.

Developers usually do not have access to production systems or data. Customer data is not used for testing.

We usually notify customers of major updates to software components that can be downloaded. Smaller updates, such as user interface adjustments, will be released regularly without clear notifications. When the latest stable version of the ZASentinel software is available, we recommend that the customer upgrade to that version.

##### Internal security test

ZASentinel uses various tools to statically analyze the code and report problems - both our proprietary code and vulnerabilities in third-party libraries. Repair detected vulnerabilities in a timely manner according to our vulnerability management policy.

##### Third-party security tests

ZASentinel cooperates with Hacker House, a well-known third-party security expert, to conduct regular security tests on its applications. Hacker House's testing activities have expanded from penetration testing to application security assurance and product analysis, including:

- Component-by-component analysis of ZASentinel in the "white box" environment
- Reverse engineering, runtime and static analysis of each component to ensure that the engineering design complies with best practice safety guidelines; 
- Perform automated stress testing, manual vulnerability discovery, and runtime and source code review
- Conduct threat modeling.

##### Penetration test request

We allow customers to penetrate our systems under certain circumstances. Customers must obtain the approval of our security team in advance and inform us in advance of the time and scope of the penetration test, and may need to sign an agreement covering such testing activities. For more information, please contact your account manager.

### 2.4. Infrastructure and physical security

#### Device security

All end users' laptops and desktops are required to install anti-virus/anti-malware software and enable full disk encryption.

#### Infrastructure change management

Every change (including infrastructure change) proposed to our production environment must be approved, and every such change and corresponding approval should be recorded. Our CI/CD pipeline supply infrastructure will be automatically changed after approval.

#### Confidential management

We use a commercial confidential management system provided by a major supplier to store secrets, such as authentication tokens, passwords, API credentials and certificates. Keys rotate regularly.

#### Server reinforcement

We use Google Cloud to provide pre-reforced server infrastructure. We mainly interact with the server by deploying Docker containers coordinated with Kubernetes.

#### Network segmentation

Our production network is divided into different regions according to the security level. Each environment has its own subnet, and only network policies based on predefined allowable lists allow internal communication.

#### Personnel safety

ZASentinel employees are mentored and trained to ensure that their physical working environment remains safe (whether at home, airport, coffee shop or other remote location) and that any work equipment is properly protected, including when not in use.

#### Usability & Flexibility

When our services need access to mission-critical network resources, service availability is crucial. We ensure very high service availability by:

- Use world-class infrastructure providers - we use the Alibaba Cloud platform to host our core components. The technology that powers Alibaba Cloud is used by Valley Alibaba to support its own applications, which are used by billions of people.
- The infrastructure of ZASentinel, which uses multiple geographically separated data centers, hosts multiple physically separated Alibaba cloud platform data centers to achieve redundancy. This helps balancing the load and reduces the risk of natural disasters in the environment and other specific locations.
- Implement fault-tolerant and redundant infrastructure - our services are provided by multiple data centers, which reflect each other's capabilities. If one data center has usability problems, other data centers will automatically bear the load.
- Provide the ability to resist DDOS attacks - we have implemented some measures to reduce the risk of DDOS attacks.
- 24/7 monitoring - We use various automation tools to monitor our services 24/7 and remind us of any service availability problems.

#### Disaster Recovery & Business Continuity Planning

ZASentinel has a written DRP/BCP. Our goal is to ensure that customers have access to our services whenever they need it.

### About ZASentinel, etc.

#### Contact us.

Email: corerman@gmail.com

##### Homo sapiens

1403, 14th Floor, Guyang Century Building, Xuanwu District, Nanjing City, Jiangsu Province