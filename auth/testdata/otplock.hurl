# Generate an otp for lock user
POST {{host}}/otp
HX-Request: true
[FormParams]
email: lock@artisons.me
HTTP 200
[Asserts]
xpath "//form[@hx-post='/login']" exists

# Displays an error message when requesting new OTP too quickly
POST {{host}}/otp
HX-Request: true
[FormParams]
email: lock@artisons.me
HTTP 200
[Asserts]
body contains "You need to wait before asking another OTP."

# Displays an error for first tentative
POST {{host}}/login
HX-Current-Url: http://localhost:8080/sso
HX-Request: true
[FormParams]
email: lock@artisons.me
otp: 1
otp: 2 
otp: 3 
otp: 4 
otp: 5
otp: 6 
HTTP 200
[Asserts]
body contains "The OTP does not match."

# Displays an error for second tentative
POST {{host}}/login
HX-Current-Url: http://localhost:8080/sso
HX-Request: true
[FormParams]
email: lock@artisons.me
otp: 1
otp: 2 
otp: 3 
otp: 4 
otp: 5
otp: 6 
HTTP 200
[Asserts]
body contains "The OTP does not match."

# Cancel the tentative after the third tentative
POST {{host}}/login
HX-Current-Url: http://localhost:8080/sso
HX-Request: true
[FormParams]
email: lock@artisons.me
otp: 1
otp: 2 
otp: 3 
otp: 4 
otp: 5
otp: 6 
HTTP 200
[Asserts]
body contains "You reached the max tentatives."