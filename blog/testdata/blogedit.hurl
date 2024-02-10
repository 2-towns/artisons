# Editing an article without title shows an error 
POST {{host}}/admin/blog/2/edit
HX-Request: true
[MultipartFormData]
title: 
description: 100 % coton pour les couleurs unies
status: online
image_1:file,image.jpg;
slug: blog-edit-{{time}}
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
header "HX-Retarget" == "#title-error" 
header "HX-Reswap" == "innerHTML show:#title-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing an article without description shows an error 
POST {{host}}/admin/blog/2/edit
HX-Request: true
[MultipartFormData]
title: Great article
description:
status: online
image_1:file,image.jpg;
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
header "HX-Retarget" == "#description-error" 
header "HX-Reswap" == "innerHTML show:#description-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing an article with a bad picture shows an error 
POST {{host}}/admin/blog/2/edit
HX-Request: true
[MultipartFormData]
title: Great article
description: Great description
status: online
image:file,image.txt;
slug:great-article-{{time}}
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
header "HX-Retarget" == "#image-error" 
header "HX-Reswap" == "innerHTML show:#image-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing an article with existing slug show an error
POST {{host}}/admin/blog/2/edit
HX-Request: true
[MultipartFormData]
title: Great article
lang: en
description: Great description
status: online
image:file,image.jpg;
slug: mangez-de-lail
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
header "HX-Retarget" == "#slug-error" 
header "HX-Reswap" == "innerHTML show:#slug-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing an article with correct data works
POST {{host}}/admin/blog/2/edit
HX-Request: true
[MultipartFormData]
title: Great great article 
description: Great description
status: online
image:file,image.jpg;
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/blog"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"
