# Destroy the tag
POST {{host}}/admin/tags/crazy/delete
HX-Request: true
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
xpath "//td[text()='crazy']" not exists
