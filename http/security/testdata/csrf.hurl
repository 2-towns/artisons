# A POST request without HX-Request returns an error
POST {{host}}/otp
HX-Current-Url: http://localhost:8080/otp
[FormParams]
email: user@artisons.me
HTTP 400
[Asserts]
body contains "Your are not authorized to access to this page. This will be reported."