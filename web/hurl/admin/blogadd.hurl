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

# Adding an article without title shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: 
lang: en
description: 100 % coton pour les couleurs unies
status: online
image_1:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#title-error" 
header "HX-Reswap" == "innerHTML show:#title-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The title is invalid." 

# Adding an article without description shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description:
status: online
image_1:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#description-error" 
header "HX-Reswap" == "innerHTML show:#description-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The description is invalid." 

# Adding an article without status shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: 
image_1:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#status-error" 
header "HX-Reswap" == "innerHTML show:#status-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The status is invalid." 

# Adding an article with bad status shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: hello
image_1:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#status-error" 
header "HX-Reswap" == "innerHTML show:#status-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The status is invalid."

# Adding an article without lang shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: 
description: Great description
status: online
image_1:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#lang-error" 
header "HX-Reswap" == "innerHTML show:#lang-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The lang is invalid." 

# Adding an article with bad lang shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: !!!
description: Great description
status: online
image_1:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Retarget" == "#lang-error" 
header "HX-Reswap" == "innerHTML show:#lang-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The lang is invalid."

# Adding an article without file shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: online
HTTP 200
[Asserts]
header "HX-Retarget" == "#image-error" 
header "HX-Reswap" == "innerHTML show:#image-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The image is required." 

# Adding an article with a bad picture shows an error 
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: online
image:file,image.txt;
HTTP 200
[Asserts]
header "HX-Retarget" == "#image-error" 
header "HX-Reswap" == "innerHTML show:#image-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The image is invalid." 

# Adding an article with correct data works
POST {{host}}/admin/blog/add.html
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: online
image:file,image.jpg;
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/blog.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"
