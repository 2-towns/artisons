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
email: hello@world.com
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/index.html"

# Get the last id 
GET {{host}}/admin/tags.html
HTTP 200
[Asserts]
body not contains "No results found." 
[Captures]
key: xpath "normalize-space((//tr)[1]/td)"

# Destroy the tag
POST {{host}}/admin/tags/{{key}}/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//td[text()='{{key}}']" not exists
