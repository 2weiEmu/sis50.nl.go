package main

import (
	"database/sql"

	"github.com/gorilla/websocket"
	"github.com/mattn/go-sqlite3"
)

// NOTE: to make the LSP shut up
var sqlite3Conn sqlite3.SQLiteConn

type Context struct {
	DB             *sql.DB
	Upgrader       websocket.Upgrader
	Websocket_List []*websocket.Conn
}

func CreateInitialContext(database_conn_parameters string) (Context, error) {

	ctx := Context{}

	db, err := sql.Open("sqlite3", database_conn_parameters)
	if err != nil {
		return Context{}, err
	}
	ctx.DB = db

	return ctx, nil
}

func CleanupContext(context Context) {
	context.DB.Close()
}
