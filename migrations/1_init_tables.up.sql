 CREATE TABLE IF NOT EXISTS People(
    Id integer primary key ,
    Name varchar(255) NOT NULL ,
    Surname varchar(255) NOT NULL ,
    Patronymic varchar(255)
 );

create INDEX if not exists idx_id_person ON People(Id);

CREATE  TABLE IF NOT EXISTS Car (
     Id INTEGER PRIMARY KEY,
     regNum  text NOT NULL,
     mark VARCHAR(255) NOT NULL,
     model VARCHAR(255) NOT NULL,
     year   INTEGER ,
     owner_id INTEGER not null ,
     CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES People (Id) On DELETE CASCADE ON UPDATE CASCADE
 );

create INDEX  if not exists idx_regnum_car ON Car(regNum);
