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
