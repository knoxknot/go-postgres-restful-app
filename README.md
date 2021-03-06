# A Restful Boostore App with a Database Datastore

Again I will be showcasing this on an ubuntu operating system. Thus we begin by installing and starting postgresql

`echo "y" | sudo apt install postgresql postgresql-contrib`  
`sudo service postgresql start`

**Creating Initial Operations on the Database**  
`sudo -u postgres psql`  # log in with default postgres user and run the below commands  
<code>
 CREATE DATABASE bookstore CONNECTION LIMIT 10;  <br>
 CREATE USER developer WITH PASSWORD 'p2ssW0rd'; <br>
 ALTER ROLE developer SET client_encoding TO 'utf8'; <br>
 ALTER ROLE developer SET default_transaction_isolation TO 'read committed'; <br>
 ALTER ROLE developer SET timezone TO 'UTC'; <br>
 GRANT ALL PRIVILEGES ON DATABASE bookstore TO developer; <br>
 \q
</code>  
**RUN** commands below to generate the sql statement and import data into database  
`sudo -u postgres psql 'postgres://developer:p2ssW0rd@localhost:5432/bookstore' < bookstore.sql`


#### Testing with Curl for Raw Data   

- curl -X GET localhost:8080/api/v1/books
- curl -X GET localhost:8080/api/v1/books/show?isbn=978-1470184841
-	curl -i -X POST -d "isbn=978-1470184841&title=Metamorphosis&author=Franz Kafka&price=5.90" localhost:8080/api/v1/books/create
-	curl -i -X DELETE localhost:8080/api/v1/books/delete?isbn=978-1470184841
- curl -i -X PUT -d "isbn=978-1470184841&title=Metamorphosis&author=Frank Kafka&price=7.99" localhost:8080/api/v1/books/update?isbn=978-1470184841