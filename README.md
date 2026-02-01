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

10. connect Jenkins to github repository via webook to pull the source code and the docker file.

11. copy the file called deployment.yaml to your server (the one which we will be using to deploy the app itself)

12. copy the jenkinsfile into jenkins , it will pull the source code + the dockerfile and build it , afterwards it will push it to     docker registry and finally it will deploy the app.



test
test
test
test
test
