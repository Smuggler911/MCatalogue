package postgresql

import (
	"MCatalogue/internal/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sagikazarmark/slog-shim"
)

type Repository struct {
	db *sql.DB
}

func New(dbInfo string) (*Repository, error) {
	const op = "repository.postgresql.New"
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	fmt.Println()
	return &Repository{db: db}, nil
}

func (r *Repository) Stop() error {
	return r.db.Close()
}

func (r *Repository) RunMigrations() error {
	const op = "repository.postgresql.RunMigrations"
	queryCreate :=
		`
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
  `
	queryAddValue :=
		`
INSERT INTO people(id,name,surname,patronymic)
VALUES
   (1,'John','Carlos','Gustavos'),
   (2,'Кирилл','Фризов','Владимирович'),
   (3,'Мария','Филипова','Филиповна'),
   (4,'Сергей','Наутилус','Григорьевич');

INSERT INTO car(id,regnum,mark,model,year,owner_id)
VALUES
    (1,'X123XX159','Lada','Vesta',2003,1),
    (2,'X123XX259','Lada','Granta',2004,2),
    (3,'X123XX459','Lada','Kalina',2013,1),
    (4,'X123XX269','Lada','X-RAY',2020,4),
    (5,'X133XX189','Lamborghini','Hurricane',20013,3),
    (6,'X153XX169','Ford','Focus',20020,4);

`
	_, err := r.db.Exec(queryCreate)
	if err != nil {
		slog.Error("Error creating tables  operation(:%s): %s", op, err)
		return err

	}

	_, err = r.db.Exec(queryAddValue)
	if err != nil {
		slog.Error("Error adding values operation(:%s): %s", op, err)
		return err

	}

	slog.Info("tables created")
	slog.Info("values added")

	return nil

}

func (r *Repository) GetAllCarData(columnName string, param string, limit int, offset int) ([]model.Car, error) {
	const op = "repository.postgresql.getAllCarData"

	var cars []model.Car

	sqlQuery := fmt.Sprintf(`SELECT * FROM car 
    WHERE %s=$1
    ORDER BY year DESC 
    LIMIT $2 OFFSET $3`, columnName)

	personQuery := `SELECT name,surname,patronymic FROM people where id = $1`

	carRows, err := r.db.Query(sqlQuery, param, limit, offset)
	if err != nil {
		slog.Error("Error querying  rows  (operation: %s): %v", op, err)
		return []model.Car{}, err
	}

	for carRows.Next() {
		var car model.Car
		var person model.People

		err = carRows.Scan(&car.Id, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.OwnerId)
		if err != nil {
			slog.Error("Error scanning rows  (operation: %s): %v", op, err)
			return []model.Car{}, err
		}

		personRows := r.db.QueryRow(personQuery, car.OwnerId)

		err = personRows.Scan(&person.Name, &person.Surname, &person.Patronymic)
		if err != nil {
			slog.Error("Error  scanning person rows (operation: %s): %v", op, err)
			return []model.Car{}, err
		}
		car.Owner = person

		cars = append(cars, car)
	}

	return cars, nil
}

func (r *Repository) DeleteCarRow(id int) error {
	const op = "repository.postgresql.deleteCarRow"

	carQuery := `DELETE FROM car WHERE id = $1`
	_, err := r.db.Exec(carQuery, id)
	if err != nil {
		slog.Error("Error deleting car (operation: %s): %v", op, err)
		return err
	}

	return nil
}

func (r *Repository) EditCarRow(id int, car model.Car) error {
	const op = "repository.postgresql.editCarRow"

	updateQuery := `
    Update  
        car
    Set 
      mark=
      	CASE 
          WHEN $1= '' 
          	THEN mark
            ELSE $1
      	  END ,
      model=
        CASE
          	WHEN $2=''
            	THEN model 
            	ELSE $2
         	END ,
      owner_id=
        CASE 
          WHEN $3=0 
              THEN owner_id
              ELSE $3
          END ,
      year=
        CASE 
           	WHEN $4=0 
           	 	THEN year
           	    ELSE $4
         	END ,
	  regNum=
	  	CASE 
	     	WHEN $5=''
    			THEN regNum
    			ELSE $5
      		END 
   	  WHERE 
        id = $6
`
	_, err := r.db.Exec(updateQuery, car.Mark, car.Model, car.OwnerId, car.Year, car.RegNum, id)
	if err != nil {
		slog.Error("Error updating car (operation: %s): %v", op, err)
		return fmt.Errorf("error updating car (operation: %s): %w", op, err)
	}
	return nil
}
