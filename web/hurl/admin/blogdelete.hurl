# Destroy the article
POST {{host}}/admin/blog/5/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//td[text()='5']" not exists

# Destroy the article does not work for cms
POST {{host}}/admin/blog/3/delete.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The item cannot be deleted." 

