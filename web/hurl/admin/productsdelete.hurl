# Get the last id 
GET {{host}}/admin/products.html
HTTP 200
[Asserts]
body not contains "No results found." 
[Captures]
id: xpath "normalize-space((//tr)[2]/td)"

# Destroy the product
POST {{host}}/admin/products/{{id}}/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//td[text()='{{id}}']" not exists
