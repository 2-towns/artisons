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
header "HX-Redirect" == "/admin"

# Display the dashboard
GET {{host}}/admin
HTTP 200
[Asserts]
xpath "//select[@name='days'][@hx-get='/admin']" exists
xpath "//script[@id='visits']" exists
xpath "//script[@id='unique-visits']" exists
xpath "//script[@id='pageviews']" exists
xpath "//script[@id='bounce-rates']" exists
xpath "//script[@id='orders']" exists
xpath "//script[@id='orders-counts']" exists
xpath "//script[@src='/public/vendor/admin/chartist.js']" exists
xpath "//script[@src='/public/vendor/admin/chartist-plugin-tooltip.js']" exists
xpath "//input[@id='days'][@value='7']" exists

# Change the pagination
GET {{host}}/admin?days=14
HX-Request: true
HTTP 200
[Asserts]
xpath "//select[@name='days'][@hx-get='/admin']" not exists
xpath "//script[@src='/public/vendor/admin/chartist.js']" not exists
xpath "//script[@src='/public/vendor/admin/chartist-plugin-tooltip.js']" not exists
xpath "//script[@id='visits']" exists
xpath "//script[@id='unique-visits']" exists
xpath "//script[@id='pageviews']" exists
xpath "//script[@id='bounce-rates']" exists
xpath "//script[@id='orders']" exists
xpath "//script[@id='orders-counts']" exists
xpath "//input[@id='days'][@value='14']" exists