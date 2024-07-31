curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbiI6IjcwMGE1YWQ3LWE4NzctNGFhNC1iMDllLTFlMWMzYmQzMTYwNyIsInVzZXJuYW1lIjoiYWRtaW4iLCJleHAiOjE3MTQ1NTY4MDUsIm5iZiI6MTcxNDU1Njc4NSwiaWF0IjoxNzE0NTU2Nzg1fQ.bP5YJxHhDz7u96hmsTVWZf2_eEtsHwv4VJZaCCauRQc" http://localhost:8080/admin

### Add the go toolchain 


go release version should be added 
https://github.com/golang/go/issues/62278


###

multiple time calling for secret unlock in dashboard



## From Admin 

### To see all the users

/users/signed-up    
SELECT id,name,username, COALESCE(encryption_key, '') AS "publicKey" FROM users where signed_up = true
