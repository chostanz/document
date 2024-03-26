# backend
Base URL = http://localhost:1234

# list all form
//you can change it according your data stored at document_ms
- Dampak Analisa -> filtered by document_code = 'DA' //you can change it at the service.go
- ITCM -> filtered by document_code = 'ITCM' 
- Berita Acara -> filtered by document_code = 'BA'

# auth
//you can change it at middleware, according to your data stored at role_ms 
Role required:  
  Member -> middleware with role_code = 'M'
  Admin -> middleware with role_code = 'A'
  Superadmin -> middleware with role_code = 'SA'