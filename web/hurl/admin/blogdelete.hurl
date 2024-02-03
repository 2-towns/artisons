# Get the last id 
GET {{host}}/admin/blog.html
HTTP 200
[Asserts]
body not contains "No results found." 
[Captures]
id: xpath "normalize-space((//tr)[2]/td)"

# Destroy the article
POST {{host}}/admin/blog/{{id}}/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//td[text()='{{id}}']" not exists

# Destroy the article does not work for cms
POST {{host}}/admin/blog/5/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The item cannot be deleted." 

