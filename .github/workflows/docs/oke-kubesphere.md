# KubeSphere on Oracle OKE

This guide walks you through the steps of deploying KubeSphere on Oracle Kubernetes Engine.
Create a Kubernetes Cluster

    A standard Kubernetes cluster in OKE is a prerequisite of installing KubeSphere. Go to the navigation menu and refer to the image below to create a cluster.

    oke-cluster

    In the pop-up window, select Quick Create and click Launch Workflow.

    oke-quickcreate

    Note
    In this example, Quick Create is used for demonstration which will automatically create all the resources necessary for a cluster in Oracle Cloud. If you select Custom Create, you need to create all the resources (such as VCN and LB Subnets) by yourself.

    Next, you need to set the cluster with basic information. Here is an example for your reference. When you finish, click Next.

    set-basic-info

    Note
        To install KubeSphere 3.3 on Kubernetes, your Kubernetes version must be v1.19.x, v1.20.x, v1.21.x, * v1.22.x, * v1.23.x, and * v1.24.x. For Kubernetes versions with an asterisk, some features of edge nodes may be unavailable due to incompatability. Therefore, if you want to use edge nodes, you are advised to install Kubernetes v1.21.x or earlier.
        It is recommended that you should select Public for Visibility Type, which will assign a public IP address for every node. The IP address can be used later to access the web console of KubeSphere.
        In Oracle Cloud, a Shape is a template that determines the number of CPUs, amount of memory, and other resources that are allocated to an instance. VM.Standard.E2.2 (2 CPUs and 16G Memory) is used in this example. For more information, see Standard Shapes.
        3 nodes are included in this example. You can add more nodes based on your own needs especially in a production environment.

    Review cluster information and click Create Cluster if no adjustment is needed.

    create-cluster

    After the cluster is created, click Close.

    cluster-ready

    Make sure the Cluster Status is Active and click Access Cluster.

    access-cluster

    In the pop-up window, select Cloud Shell Access to access the cluster. Click Launch Cloud Shell and copy the code provided by Oracle Cloud.

    cloud-shell-access

    In Cloud Shell, paste the command so that we can execute the installation command later.

    cloud-shell-oke

    Warning
    If you do not copy and execute the command above, you cannot proceed with the steps below.

Install KubeSphere on OKE

    Install KubeSphere using kubectl. The following commands are only for the default minimal installation.

    kubectl apply -f https://github.com/kubesphere/ks-installer/releases/download/v3.3.1/kubesphere-installer.yaml

    kubectl apply -f https://github.com/kubesphere/ks-installer/releases/download/v3.3.1/cluster-configuration.yaml

    Inspect the logs of installation:

    kubectl logs -n kubesphere-system $(kubectl get pod -n kubesphere-system -l 'app in (ks-install, ks-installer)' -o jsonpath='{.items[0].metadata.name}') -f

    When the installation finishes, you can see the following message:

    #####################################################
    ###              Welcome to KubeSphere!           ###
    #####################################################

    Console: http://10.0.10.2:30880
    Account: admin
    Password: P@88w0rd

    NOTES：
      1. After logging into the console, please check the
        monitoring status of service components in
        the "Cluster Management". If any service is not
        ready, please wait patiently until all components 
        are ready.
      2. Please modify the default password after login.

    #####################################################
    https://kubesphere.io             20xx-xx-xx xx:xx:xx

Access KubeSphere Console

Now that KubeSphere is installed, you can access the web console of KubeSphere either through NodePort or LoadBalancer.

    Check the service of KubeSphere console through the following command:

    kubectl get svc -n kubesphere-system

    The output may look as below. You can change the type to LoadBalancer so that the external IP address can be exposed.

    console-nodeport

    Tip
    It can be seen above that the service ks-console is being exposed through a NodePort, which means you can access the console directly via NodeIP:NodePort (the public IP address of any node is applicable). You may need to open port 30880 in firewall rules.

    Execute the command to edit the service configuration.

    kubectl edit svc ks-console -o yaml -n kubesphere-system

    Navigate to type and change NodePort to LoadBalancer. Save the configuration after you finish.

    change-svc-type

    Execute the following command again and you can see the IP address displayed as below.

    kubectl get svc -n kubesphere-system

    console-service

    Log in to the console through the external IP address with the default account and password (admin/P@88w0rd). In the cluster overview page, you can see the dashboard.
**************************************************
Collecting installation results ...
#####################################################
###              Welcome to KubeSphere!           ###
#####################################################

Console: http://10.0.10.12:30880
Account: admin
Password: P@88w0rd
NOTES：
  1. After you log into the console, please check the
     monitoring status of service components in
     "Cluster Management". If any service is not
     ready, please wait patiently until all components 
     are up and running.
  2. Please change the default password after login.

#####################################################
https://kubesphere.io             2022-12-12 09:07:53
#####################################################
