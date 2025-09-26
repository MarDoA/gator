# gator
gator is a simple feed aggregator 

Installation

Requirements:
-Go 1.23.4+
-PostgreSQL 10+

Installation:
-using go install
go install github.com/MarDoA/gator@latest

you’ll need to create a configuration file named .gatorconfig.json in your home directory
with your string connection in it
example:
{
  "db_url": "postgres://user:password@localhost:5432/gator_db?sslmode=disable"
}

replace user,password and gator_db

