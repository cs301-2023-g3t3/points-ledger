
### things to note:
1. need to change go version to 1.18 to work in lambda
2. need to re-build and zip each time there are code changes made, before redeploying using terraform



### to build and zip app to deploy 
```
GOOS=linux GOARCH=amd64 go build -o main *.go
```

```
zip main.zip main
```


### inside /terraform directory

Init terraform
```
terraform init
```

View planned changes
```
terraform plan
```

Apply changes
```
terraform apply
```

destroy all created infrastructure
``` 
terraform destroy
```