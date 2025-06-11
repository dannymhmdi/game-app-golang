# what is cloud native means?
It means high access to infinite resourses.
Imagine you have an application which user upload-download files if your hard goes fool
what do you do? but when you use a cloud service like AWS it has a s3 or minIO service for storing file and 
you dont need to worry about file storing infrastructure this service is high available means AWS itself make backup
of files.

## AWS definition:
Cloud native is the software approach of building, deploying, and managing modern applications in cloud computing environments. Modern companies want to build highly scalable, flexible, and resilient applications that they can update quickly to meet customer demands. To do so, they use modern tools and techniques that inherently support application development on cloud infrastructure. These cloud-native technologies support fast and frequent changes to applications without impacting service delivery, providing adopters with an innovative, competitive advantage.

### Independence:
Each part of cloudnative applications are independent to each other ,this means also you can
develop , deploy and maintain them individually.

### Resiliency:
A well-designed cloud native application is able to survive and stay online even in the event of an infrastructure outage.
 #### Key Principles for Surviving Infrastructure Outages
1. Microservices Architecture

Decomposing the app into loosely coupled services ensures that a failure in one component doesn’t bring down the entire system.

Each service can be scaled, updated, and recovered independently.

2. Auto-Scaling & Self-Healing

Kubernetes (or managed services like EKS/GKE/AKS) can automatically restart failed containers or spin up new instances.

3. Chaos Engineering & Fault Tolerance Testing

Proactively testing failures (e.g., killing instances, network partitions) using tools like Chaos Monkey or Gremlin ensures the system can recover gracefully.

### No Downtime
Imagine if you want to run your database service you have one instance but when run it by cloud provider it consider
3 instance one master and 2 slaves when master goes down use one of slaves to prevent downtime.

