# Destroy the article
POST {{host}}/admin/blog/3/delete
HX-Request: true
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
xpath "//td[text()='3']" not exists

# Destroy the article does not work for cms
POST {{host}}/admin/blog/4/delete
HX-Request: true
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
[Asserts]
xpath "//div[@class='alert alert-danger box']" exists 
xpath "normalize-space(//p[@class='alert-message']/text())" == "The item cannot be deleted." 

