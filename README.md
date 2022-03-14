# Golang-Authenication-Middleware

Created middleware using golang to demonstrate working of authentication. Using this middleware, Person can be authenticated on middleware itself.
Here are the instructions to use this middleware.

##### 1) Clone project on your system and run command mentioned below.
```
go mod download
```

##### 2) Run project using command below.
```
 go run .
```

##### 3) To genrate Token, try to reach out 
```
http://localhost:8080/genToken
```

##### 4) To test whether token is generate
```
http://localhost:8080/homepage
```

##### 5) To clear token, try to reach out
```
http://localhost:8080/clearToken
```

