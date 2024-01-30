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

# Adding an article without title shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: 
lang: en
description: 100 % coton pour les couleurs unies
status: online
image:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#title-error" 
header "HX-Reswap" == "innerHTML show:#title-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding an article without description shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description:
status: online
image:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#description-error" 
header "HX-Reswap" == "innerHTML show:#description-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding an article without file shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: online
slug:blog-without-file-{{time}}
HTTP 200
[Asserts]
header "HX-Retarget" == "#image-error" 
header "HX-Reswap" == "innerHTML show:#image-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding an article with a bad picture shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: online
image:file,image.txt;
slug:blog-with-bad-picture-{{time}}
HTTP 200
[Asserts]
header "HX-Retarget" == "#image-error" 
header "HX-Reswap" == "innerHTML show:#image-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding an article with existing slug show an error
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: online
image:file,image.jpg;
slug: mangez-de-lail
HTTP 200
[Asserts]
header "HX-Retarget" == "#slug-error" 
header "HX-Reswap" == "innerHTML show:#slug-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Adding an article with correct data works
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: online
image:file,image.jpg;
slug: blog-{{time}}
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/blog.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"
