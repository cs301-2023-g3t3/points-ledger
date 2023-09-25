
### build and zip app to deploy
```
GOOS=linux GOARCH=amd64 go build -o main *.go
```

```
zip main.zip main
```


### inside terraform directory

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