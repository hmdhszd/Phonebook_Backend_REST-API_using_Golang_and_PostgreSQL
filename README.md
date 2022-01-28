
## Run Program

Get a copy of the program:

```bash
  git clone https://github.com/hmdhszd/Phonebook_Backend_REST-API_using_Golang_and_PostgreSQL.git
```

Go to the project directory:

```bash
  cd Phonebook_Backend_REST-API_using_Golang_and_PostgreSQL
```

Build image and start the Backend app and PostgreSQL database:

```bash
  docker-compose up -d --build --force-recreate
```

Check if they are running properly:

```bash
  docker-compose ps
```

The output should be like this:

```bash
        Name                    Command            State            Ports         
----------------------------------------------------------------------------------
phonebook_backend_api   ./Phonebook_Backend_REST   Up      0.0.0.0:8000->8000/tcp,
                        -A ...                             :::8000->8000/tcp      
phonebookdb             docker-entrypoint.sh       Up      0.0.0.0:5432->5432/tcp,
                        postgres                           :::5432->5432/tcp  

```


## Connect to database and add a table


Go inside of the container:

```bash
  docker exec -it phonebookdb /bin/sh
```

Then connect to PostgreSQL database:

```bash
  psql -U postgres
```

Now, add a table:

```bash
CREATE TABLE contacts (
  id SERIAL,
  PhoneNumber VARCHAR(250) NOT NULL,
  FullName VARCHAR(250) NOT NULL,
  Address VARCHAR(250) NOT NULL,
  Email VARCHAR(250) NOT NULL,
  PRIMARY KEY (id)
);
```

And add some data into the table:

```bash
  
INSERT INTO contacts (
  PhoneNumber,
  FullName,
  Address,
  Email
)
VALUES
    ('00989121234567', 'Hamid Hosseinzadeh','No.59, Iran, Tehran','hamid@hszd.ir'),
    ('0011234567890', 'Brad Pit','No.24, Iran, Yazd','brad@pit.com'),
    ('0010982847492', 'George Clooney','No.100, USA, LA','george@clooney.com'),
    ('0029848289238', 'Leonardo DiCaprio','9255 Sunset Blvd., Suite 615, West Hollywood, California 90069 United States','Leonardo@DiCaprio.com');

```

At the end, check if the data has been added successfully:

```bash
  SELECT * FROM contacts;
```

The output should be like this:

```bash
 id |  phonenumber   |      fullname      |                                   address                                    |         email         
----+----------------+--------------------+------------------------------------------------------------------------------+-----------------------
  1 | 00989121234567 | Hamid Hosseinzadeh | No.59, Iran, Tehran                                                          | hamid@hszd.ir
  2 | 0011234567890  | Brad Pit           | No.24, Iran, Yazd                                                            | brad@pit.com
  3 | 0010982847492  | George Clooney     | No.100, USA, LA                                                              | george@clooney.com
  4 | 0029848289238  | Leonardo DiCaprio  | 9255 Sunset Blvd., Suite 615, West Hollywood, California 90069 United States | Leonardo@DiCaprio.com
(4 rows)
```
## Now, we can test our API using Postman:

## Add a new contact

![App Screenshot](https://github.com/hmdhszd/Phonebook_Backend_REST-API_using_Golang_and_PostgreSQL/blob/master/Screenshots/add-new-contact.png)

## Get all contacts

![App Screenshot](https://github.com/hmdhszd/Phonebook_Backend_REST-API_using_Golang_and_PostgreSQL/blob/master/Screenshots/get-all-contacts.png)

## Search in contacts

![App Screenshot](https://github.com/hmdhszd/Phonebook_Backend_REST-API_using_Golang_and_PostgreSQL/blob/master/Screenshots/search-by-name.png)

## Delete a contact by it's name

![App Screenshot](https://github.com/hmdhszd/Phonebook_Backend_REST-API_using_Golang_and_PostgreSQL/blob/master/Screenshots/delete-by-name.png)

## Delete a contact by it's number

![App Screenshot](https://github.com/hmdhszd/Phonebook_Backend_REST-API_using_Golang_and_PostgreSQL/blob/master/Screenshots/delete-by-number.png)

## Delete all contacts

![App Screenshot](https://github.com/hmdhszd/Phonebook_Backend_REST-API_using_Golang_and_PostgreSQL/blob/master/Screenshots/delete-all.png)

