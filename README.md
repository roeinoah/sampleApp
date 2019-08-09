# README.md

Prepare the kubernetes cluster on AWS using terraform : 

 - connect to terraform server

1.install kubectl
  install aws-iam-authenticator
  install eksctl
  
2. copy vpc.tf and bucket.tf (includes EKS configuration) to your terraform server , and run "terraform apply"

3. wait for the deployment to finish.

4. when it's done run the following commands : 
    
      "aws eks --region eu-central-1 update-kubeconfig --name terraform-eks-demo"
  
      "terraform output config_map_aws_auth >> aws-auth-cm.yaml"
      
      "kubectl apply -f aws-auth-cm.yaml"
      
5. wait for all the kuberenetes nodes to come up using : kubectl get nodes --watch

6. copy the follwoing yaml files into the server : jenkins-deployment.yaml , jenkins-service.yaml .

7. run the following commands : 
  kubectl create -f jenkins-deployment.yaml
  kubectl create -f jenkins-service.yaml
  
8. wait for the jenkins deployment to be ready .

9. connect to jenkins and configure it with admin user and add docker plugin.

10. connect Jenkins to github repository to pull the source code and the docker file.

11. create new item on jenkins and add github as a repository , add build docker image step from the dockerfile we have just pulled from the repository.

12. post build : sudo docker run sudo docker run --rm -it gos3cve

