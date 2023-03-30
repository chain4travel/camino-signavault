#!/bin/bash

#
# This script will drop the database and run all migrations. It takes one argument, the name of the database to reset.
# The script must be executed from the root of the project.
#

mysql -uroot -ppassword -e"DROP DATABASE IF EXISTS $1"
mysql -uroot -ppassword -e"CREATE DATABASE $1"

migrate -source file://db/migrations -database "mysql://root:password@tcp(127.0.0.1:3306)/$1" up
