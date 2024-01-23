# Display otp form when the correct email is used
POST {{host}}/auth/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/auth/login.html']" exists
[Captures]
glue: xpath "normalize-space(//input[@id='glue']/@value)"

# Redirect to dashboard when the otp is correct
POST {{host}}/auth/login.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
glue: {{glue}}
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
GET {{host}}/admin/filters.html
HTTP 200
[Asserts]
body not contains "No results found." 
[Captures]
key: xpath "normalize-space((//tr)[1]/td)"

# Destroy the filter
POST {{host}}/admin/filters/{{key}}/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//td[text()='{{key}}']" not exists

# Destroy the filter non editable show error
POST {{host}}/admin/filters/colors/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The filter cannot be editable." 

