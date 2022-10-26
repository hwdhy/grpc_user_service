rm *.pem

#创建秘钥
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C = CN/ST = hunan/L = changsha/O = tech school/OU =Education/CN = *.techschool.guru/emailAddress = testschool.guru@gmail.com"

echo "CA'S self-singned certificate"
openssl x509 -in ca-cert.pem -noout -text

openssl req  -newkey rsa:4096 -nodes  -keyout server-key.pem -out server-req.pem -subj "/C = CN/ST = Tle de France/L = Paris/O = PC Book/OU =Computer/CN = *.pcbook.com/emailAddress = pcbook@gmail.com"

openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

openssl x509 -in server-cert.pem -noout -text


openssl req  -newkey rsa:4096 -nodes  -keyout client-key.pem -out client-req.pem -subj "/C = CN/ST = Tle de France/L = Paris/O = PC Book/OU =Computer/CN = *.pcclient.com/emailAddress = client@gmail.com"

openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf

openssl x509 -in client-cert.pem -noout -text