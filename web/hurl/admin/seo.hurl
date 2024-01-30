# Display otp form when the correct email is used
POST {{host}}/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/login.html']" exists


# Redirect to dashboard when the otp is correct
POST {{host}}/login.html
HX-Current-URL: http://localhost:8080/sso.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
email:hello@world.com
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/index.html"

# Display the seo
GET {{host}}/admin/seo.html
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//table/tbody/tr" count > 0

# Access the seo with huge page shoes empty message
GET {{host}}/admin/seo.html?page=100
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

