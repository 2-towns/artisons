# Get the last id 
GET {{host}}/admin/filters
[Cookies]
wsid: 444444
HTTP  200
[Asserts]
body not contains "No results found." 
[Captures]
key: xpath "normalize-space((//tr)[1]/td)"

# Destroy the filter
POST {{host}}/admin/filters/{{key}}/delete
HX-Request: true
[Cookies]
wsid: 444444
HTTP  200
[Asserts]
xpath "//td[text()='{{key}}']" not exists

# Destroy the filter non editable show error
POST {{host}}/admin/filters/colors/delete
HX-Request: true
[Cookies]
wsid: 444444
HTTP  200
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The filter cannot be editable." 

