# Checkliste
## Erledigt:

- Login: Als Benutzer möchte nur ich meine ToDo sehen
- ToDo erstellen: Als Benutzer möchte ich ein neues ToDo/Aufgabe mit Titel und Text erstellen
  und in der Datenbank speichern.
- ToDo auflisten: Als Benutzer möchte ich alle meine ToDo in einer Liste sehen.
- ToDo löschen: Als Benutzer möchte ich ein ToDo aus der Liste löschen, wenn ich es nicht
  mehr benötige.
- Erledigt: Als Benutzer möchte ich mein ToDo als erledigt markieren.
- ToDo teilen: Als Benutzer möchte ich mein ToDo mit einem anderen Benutzer teilen, damit
  dieser auch meine ToDo sehen und als Erledigt markieren kann

## Fehlt:

- ToDo aktualisieren: Als Benutzer möchte ich ein bestehendes ToDo ändern können.
- ToDo verschieben Als Benutzer möchte ich meine ToDo neu anordnen.
- ToDo Kategorien: Als Benutzer möchte ich meine ToDo Kategorisieren
- HTTPS
- Verschlüsselung von Passwörtern
- Verschlüsselung von DB Credentials

# Database

mysql database

- Create Database todoapi

```
CREATE DATABASE `todoapi` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
```

- Create Table todos

``` 
CREATE TABLE `todos` (
  `IdTodos` int NOT NULL AUTO_INCREMENT,
  `TodosName` varchar(45) DEFAULT NULL,
  `TodosDone` tinyint(1) DEFAULT NULL,
  `TodosText` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`IdTodos`)
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

- Create Table users

```
CREATE TABLE `users` (
  `IdUsers` int NOT NULL AUTO_INCREMENT,
  `Username` varchar(45) DEFAULT NULL,
  `Password` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`IdUsers`)
) ENGINE=InnoDB AUTO_INCREMENT=68 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

- Create Table todoowners

```
CREATE TABLE `todoowners` (
  `IdTodoOwners` int NOT NULL AUTO_INCREMENT,
  `IdOfOwner` int DEFAULT NULL,
  `IdOfTodo` int DEFAULT NULL,
  PRIMARY KEY (`IdTodoOwners`)
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```
