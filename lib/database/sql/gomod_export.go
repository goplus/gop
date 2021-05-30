// Package sql provide Go+ "database/sql" package, as "database/sql" package in Go.
package sql

import (
	context "context"
	sql "database/sql"
	driver "database/sql/driver"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execmColumnTypeName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.ColumnType).Name()
	p.Ret(1, ret0)
}

func execmColumnTypeLength(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*sql.ColumnType).Length()
	p.Ret(1, ret0, ret1)
}

func execmColumnTypeDecimalSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*sql.ColumnType).DecimalSize()
	p.Ret(1, ret0, ret1, ret2)
}

func execmColumnTypeScanType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.ColumnType).ScanType()
	p.Ret(1, ret0)
}

func execmColumnTypeNullable(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*sql.ColumnType).Nullable()
	p.Ret(1, ret0, ret1)
}

func execmColumnTypeDatabaseTypeName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.ColumnType).DatabaseTypeName()
	p.Ret(1, ret0)
}

func toType0(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execmConnPingContext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.Conn).PingContext(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmConnExecContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Conn).ExecContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0, ret1)
}

func execmConnQueryContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Conn).QueryContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0, ret1)
}

func execmConnQueryRowContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sql.Conn).QueryRowContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0)
}

func execmConnPrepareContext(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*sql.Conn).PrepareContext(toType0(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmConnRaw(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.Conn).Raw(args[1].(func(driverConn interface{}) error))
	p.Ret(2, ret0)
}

func execmConnBeginTx(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*sql.Conn).BeginTx(toType0(args[1]), args[2].(*sql.TxOptions))
	p.Ret(3, ret0, ret1)
}

func execmConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.Conn).Close()
	p.Ret(1, ret0)
}

func execmDBPingContext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.DB).PingContext(toType0(args[1]))
	p.Ret(2, ret0)
}

func execmDBPing(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.DB).Ping()
	p.Ret(1, ret0)
}

func execmDBClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.DB).Close()
	p.Ret(1, ret0)
}

func execmDBSetMaxIdleConns(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sql.DB).SetMaxIdleConns(args[1].(int))
	p.PopN(2)
}

func execmDBSetMaxOpenConns(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sql.DB).SetMaxOpenConns(args[1].(int))
	p.PopN(2)
}

func execmDBSetConnMaxLifetime(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*sql.DB).SetConnMaxLifetime(args[1].(time.Duration))
	p.PopN(2)
}

func execmDBStats(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.DB).Stats()
	p.Ret(1, ret0)
}

func execmDBPrepareContext(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*sql.DB).PrepareContext(toType0(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmDBPrepare(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*sql.DB).Prepare(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmDBExecContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.DB).ExecContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBExec(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.DB).Exec(args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBQueryContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.DB).QueryContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBQuery(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.DB).Query(args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execmDBQueryRowContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sql.DB).QueryRowContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0)
}

func execmDBQueryRow(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sql.DB).QueryRow(args[1].(string), args[2:]...)
	p.Ret(arity, ret0)
}

func execmDBBeginTx(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*sql.DB).BeginTx(toType0(args[1]), args[2].(*sql.TxOptions))
	p.Ret(3, ret0, ret1)
}

func execmDBBegin(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*sql.DB).Begin()
	p.Ret(1, ret0, ret1)
}

func execmDBDriver(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.DB).Driver()
	p.Ret(1, ret0)
}

func execmDBConn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*sql.DB).Conn(toType0(args[1]))
	p.Ret(2, ret0, ret1)
}

func execDrivers(_ int, p *gop.Context) {
	ret0 := sql.Drivers()
	p.Ret(0, ret0)
}

func execmIsolationLevelString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(sql.IsolationLevel).String()
	p.Ret(1, ret0)
}

func execNamed(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := sql.Named(args[0].(string), args[1])
	p.Ret(2, ret0)
}

func execmNullBoolScan(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.NullBool).Scan(args[1])
	p.Ret(2, ret0)
}

func execmNullBoolValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(sql.NullBool).Value()
	p.Ret(1, ret0, ret1)
}

func execmNullFloat64Scan(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.NullFloat64).Scan(args[1])
	p.Ret(2, ret0)
}

func execmNullFloat64Value(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(sql.NullFloat64).Value()
	p.Ret(1, ret0, ret1)
}

func execmNullInt32Scan(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.NullInt32).Scan(args[1])
	p.Ret(2, ret0)
}

func execmNullInt32Value(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(sql.NullInt32).Value()
	p.Ret(1, ret0, ret1)
}

func execmNullInt64Scan(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.NullInt64).Scan(args[1])
	p.Ret(2, ret0)
}

func execmNullInt64Value(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(sql.NullInt64).Value()
	p.Ret(1, ret0, ret1)
}

func execmNullStringScan(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.NullString).Scan(args[1])
	p.Ret(2, ret0)
}

func execmNullStringValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(sql.NullString).Value()
	p.Ret(1, ret0, ret1)
}

func execmNullTimeScan(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.NullTime).Scan(args[1])
	p.Ret(2, ret0)
}

func execmNullTimeValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(sql.NullTime).Value()
	p.Ret(1, ret0, ret1)
}

func execOpen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := sql.Open(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func toType1(v interface{}) driver.Connector {
	if v == nil {
		return nil
	}
	return v.(driver.Connector)
}

func execOpenDB(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sql.OpenDB(toType1(args[0]))
	p.Ret(1, ret0)
}

func toType2(v interface{}) driver.Driver {
	if v == nil {
		return nil
	}
	return v.(driver.Driver)
}

func execRegister(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	sql.Register(args[0].(string), toType2(args[1]))
	p.PopN(2)
}

func execiResultLastInsertId(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(sql.Result).LastInsertId()
	p.Ret(1, ret0, ret1)
}

func execiResultRowsAffected(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(sql.Result).RowsAffected()
	p.Ret(1, ret0, ret1)
}

func execmRowScan(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sql.Row).Scan(args[1:]...)
	p.Ret(arity, ret0)
}

func execmRowsNext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.Rows).Next()
	p.Ret(1, ret0)
}

func execmRowsNextResultSet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.Rows).NextResultSet()
	p.Ret(1, ret0)
}

func execmRowsErr(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.Rows).Err()
	p.Ret(1, ret0)
}

func execmRowsColumns(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*sql.Rows).Columns()
	p.Ret(1, ret0, ret1)
}

func execmRowsColumnTypes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*sql.Rows).ColumnTypes()
	p.Ret(1, ret0, ret1)
}

func execmRowsScan(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sql.Rows).Scan(args[1:]...)
	p.Ret(arity, ret0)
}

func execmRowsClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.Rows).Close()
	p.Ret(1, ret0)
}

func execiScannerScan(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(sql.Scanner).Scan(args[1])
	p.Ret(2, ret0)
}

func execmStmtExecContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Stmt).ExecContext(toType0(args[1]), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execmStmtExec(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Stmt).Exec(args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execmStmtQueryContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Stmt).QueryContext(toType0(args[1]), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execmStmtQuery(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Stmt).Query(args[1:]...)
	p.Ret(arity, ret0, ret1)
}

func execmStmtQueryRowContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sql.Stmt).QueryRowContext(toType0(args[1]), args[2:]...)
	p.Ret(arity, ret0)
}

func execmStmtQueryRow(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sql.Stmt).QueryRow(args[1:]...)
	p.Ret(arity, ret0)
}

func execmStmtClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.Stmt).Close()
	p.Ret(1, ret0)
}

func execmTxCommit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.Tx).Commit()
	p.Ret(1, ret0)
}

func execmTxRollback(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*sql.Tx).Rollback()
	p.Ret(1, ret0)
}

func execmTxPrepareContext(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*sql.Tx).PrepareContext(toType0(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execmTxPrepare(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*sql.Tx).Prepare(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmTxStmtContext(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*sql.Tx).StmtContext(toType0(args[1]), args[2].(*sql.Stmt))
	p.Ret(3, ret0)
}

func execmTxStmt(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*sql.Tx).Stmt(args[1].(*sql.Stmt))
	p.Ret(2, ret0)
}

func execmTxExecContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Tx).ExecContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0, ret1)
}

func execmTxExec(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Tx).Exec(args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execmTxQueryContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Tx).QueryContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0, ret1)
}

func execmTxQuery(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0, ret1 := args[0].(*sql.Tx).Query(args[1].(string), args[2:]...)
	p.Ret(arity, ret0, ret1)
}

func execmTxQueryRowContext(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sql.Tx).QueryRowContext(toType0(args[1]), args[2].(string), args[3:]...)
	p.Ret(arity, ret0)
}

func execmTxQueryRow(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := args[0].(*sql.Tx).QueryRow(args[1].(string), args[2:]...)
	p.Ret(arity, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("database/sql")

func init() {
	I.RegisterFuncs(
		I.Func("(*ColumnType).Name", (*sql.ColumnType).Name, execmColumnTypeName),
		I.Func("(*ColumnType).Length", (*sql.ColumnType).Length, execmColumnTypeLength),
		I.Func("(*ColumnType).DecimalSize", (*sql.ColumnType).DecimalSize, execmColumnTypeDecimalSize),
		I.Func("(*ColumnType).ScanType", (*sql.ColumnType).ScanType, execmColumnTypeScanType),
		I.Func("(*ColumnType).Nullable", (*sql.ColumnType).Nullable, execmColumnTypeNullable),
		I.Func("(*ColumnType).DatabaseTypeName", (*sql.ColumnType).DatabaseTypeName, execmColumnTypeDatabaseTypeName),
		I.Func("(*Conn).PingContext", (*sql.Conn).PingContext, execmConnPingContext),
		I.Func("(*Conn).PrepareContext", (*sql.Conn).PrepareContext, execmConnPrepareContext),
		I.Func("(*Conn).Raw", (*sql.Conn).Raw, execmConnRaw),
		I.Func("(*Conn).BeginTx", (*sql.Conn).BeginTx, execmConnBeginTx),
		I.Func("(*Conn).Close", (*sql.Conn).Close, execmConnClose),
		I.Func("(*DB).PingContext", (*sql.DB).PingContext, execmDBPingContext),
		I.Func("(*DB).Ping", (*sql.DB).Ping, execmDBPing),
		I.Func("(*DB).Close", (*sql.DB).Close, execmDBClose),
		I.Func("(*DB).SetMaxIdleConns", (*sql.DB).SetMaxIdleConns, execmDBSetMaxIdleConns),
		I.Func("(*DB).SetMaxOpenConns", (*sql.DB).SetMaxOpenConns, execmDBSetMaxOpenConns),
		I.Func("(*DB).SetConnMaxLifetime", (*sql.DB).SetConnMaxLifetime, execmDBSetConnMaxLifetime),
		I.Func("(*DB).Stats", (*sql.DB).Stats, execmDBStats),
		I.Func("(*DB).PrepareContext", (*sql.DB).PrepareContext, execmDBPrepareContext),
		I.Func("(*DB).Prepare", (*sql.DB).Prepare, execmDBPrepare),
		I.Func("(*DB).BeginTx", (*sql.DB).BeginTx, execmDBBeginTx),
		I.Func("(*DB).Begin", (*sql.DB).Begin, execmDBBegin),
		I.Func("(*DB).Driver", (*sql.DB).Driver, execmDBDriver),
		I.Func("(*DB).Conn", (*sql.DB).Conn, execmDBConn),
		I.Func("Drivers", sql.Drivers, execDrivers),
		I.Func("(IsolationLevel).String", (sql.IsolationLevel).String, execmIsolationLevelString),
		I.Func("Named", sql.Named, execNamed),
		I.Func("(*NullBool).Scan", (*sql.NullBool).Scan, execmNullBoolScan),
		I.Func("(NullBool).Value", (sql.NullBool).Value, execmNullBoolValue),
		I.Func("(*NullFloat64).Scan", (*sql.NullFloat64).Scan, execmNullFloat64Scan),
		I.Func("(NullFloat64).Value", (sql.NullFloat64).Value, execmNullFloat64Value),
		I.Func("(*NullInt32).Scan", (*sql.NullInt32).Scan, execmNullInt32Scan),
		I.Func("(NullInt32).Value", (sql.NullInt32).Value, execmNullInt32Value),
		I.Func("(*NullInt64).Scan", (*sql.NullInt64).Scan, execmNullInt64Scan),
		I.Func("(NullInt64).Value", (sql.NullInt64).Value, execmNullInt64Value),
		I.Func("(*NullString).Scan", (*sql.NullString).Scan, execmNullStringScan),
		I.Func("(NullString).Value", (sql.NullString).Value, execmNullStringValue),
		I.Func("(*NullTime).Scan", (*sql.NullTime).Scan, execmNullTimeScan),
		I.Func("(NullTime).Value", (sql.NullTime).Value, execmNullTimeValue),
		I.Func("Open", sql.Open, execOpen),
		I.Func("OpenDB", sql.OpenDB, execOpenDB),
		I.Func("Register", sql.Register, execRegister),
		I.Func("(Result).LastInsertId", (sql.Result).LastInsertId, execiResultLastInsertId),
		I.Func("(Result).RowsAffected", (sql.Result).RowsAffected, execiResultRowsAffected),
		I.Func("(*Rows).Next", (*sql.Rows).Next, execmRowsNext),
		I.Func("(*Rows).NextResultSet", (*sql.Rows).NextResultSet, execmRowsNextResultSet),
		I.Func("(*Rows).Err", (*sql.Rows).Err, execmRowsErr),
		I.Func("(*Rows).Columns", (*sql.Rows).Columns, execmRowsColumns),
		I.Func("(*Rows).ColumnTypes", (*sql.Rows).ColumnTypes, execmRowsColumnTypes),
		I.Func("(*Rows).Close", (*sql.Rows).Close, execmRowsClose),
		I.Func("(Scanner).Scan", (sql.Scanner).Scan, execiScannerScan),
		I.Func("(*Stmt).Close", (*sql.Stmt).Close, execmStmtClose),
		I.Func("(*Tx).Commit", (*sql.Tx).Commit, execmTxCommit),
		I.Func("(*Tx).Rollback", (*sql.Tx).Rollback, execmTxRollback),
		I.Func("(*Tx).PrepareContext", (*sql.Tx).PrepareContext, execmTxPrepareContext),
		I.Func("(*Tx).Prepare", (*sql.Tx).Prepare, execmTxPrepare),
		I.Func("(*Tx).StmtContext", (*sql.Tx).StmtContext, execmTxStmtContext),
		I.Func("(*Tx).Stmt", (*sql.Tx).Stmt, execmTxStmt),
	)
	I.RegisterFuncvs(
		I.Funcv("(*Conn).ExecContext", (*sql.Conn).ExecContext, execmConnExecContext),
		I.Funcv("(*Conn).QueryContext", (*sql.Conn).QueryContext, execmConnQueryContext),
		I.Funcv("(*Conn).QueryRowContext", (*sql.Conn).QueryRowContext, execmConnQueryRowContext),
		I.Funcv("(*DB).ExecContext", (*sql.DB).ExecContext, execmDBExecContext),
		I.Funcv("(*DB).Exec", (*sql.DB).Exec, execmDBExec),
		I.Funcv("(*DB).QueryContext", (*sql.DB).QueryContext, execmDBQueryContext),
		I.Funcv("(*DB).Query", (*sql.DB).Query, execmDBQuery),
		I.Funcv("(*DB).QueryRowContext", (*sql.DB).QueryRowContext, execmDBQueryRowContext),
		I.Funcv("(*DB).QueryRow", (*sql.DB).QueryRow, execmDBQueryRow),
		I.Funcv("(*Row).Scan", (*sql.Row).Scan, execmRowScan),
		I.Funcv("(*Rows).Scan", (*sql.Rows).Scan, execmRowsScan),
		I.Funcv("(*Stmt).ExecContext", (*sql.Stmt).ExecContext, execmStmtExecContext),
		I.Funcv("(*Stmt).Exec", (*sql.Stmt).Exec, execmStmtExec),
		I.Funcv("(*Stmt).QueryContext", (*sql.Stmt).QueryContext, execmStmtQueryContext),
		I.Funcv("(*Stmt).Query", (*sql.Stmt).Query, execmStmtQuery),
		I.Funcv("(*Stmt).QueryRowContext", (*sql.Stmt).QueryRowContext, execmStmtQueryRowContext),
		I.Funcv("(*Stmt).QueryRow", (*sql.Stmt).QueryRow, execmStmtQueryRow),
		I.Funcv("(*Tx).ExecContext", (*sql.Tx).ExecContext, execmTxExecContext),
		I.Funcv("(*Tx).Exec", (*sql.Tx).Exec, execmTxExec),
		I.Funcv("(*Tx).QueryContext", (*sql.Tx).QueryContext, execmTxQueryContext),
		I.Funcv("(*Tx).Query", (*sql.Tx).Query, execmTxQuery),
		I.Funcv("(*Tx).QueryRowContext", (*sql.Tx).QueryRowContext, execmTxQueryRowContext),
		I.Funcv("(*Tx).QueryRow", (*sql.Tx).QueryRow, execmTxQueryRow),
	)
	I.RegisterVars(
		I.Var("ErrConnDone", &sql.ErrConnDone),
		I.Var("ErrNoRows", &sql.ErrNoRows),
		I.Var("ErrTxDone", &sql.ErrTxDone),
	)
	I.RegisterTypes(
		I.Type("ColumnType", reflect.TypeOf((*sql.ColumnType)(nil)).Elem()),
		I.Type("Conn", reflect.TypeOf((*sql.Conn)(nil)).Elem()),
		I.Type("DB", reflect.TypeOf((*sql.DB)(nil)).Elem()),
		I.Type("DBStats", reflect.TypeOf((*sql.DBStats)(nil)).Elem()),
		I.Type("IsolationLevel", reflect.TypeOf((*sql.IsolationLevel)(nil)).Elem()),
		I.Type("NamedArg", reflect.TypeOf((*sql.NamedArg)(nil)).Elem()),
		I.Type("NullBool", reflect.TypeOf((*sql.NullBool)(nil)).Elem()),
		I.Type("NullFloat64", reflect.TypeOf((*sql.NullFloat64)(nil)).Elem()),
		I.Type("NullInt32", reflect.TypeOf((*sql.NullInt32)(nil)).Elem()),
		I.Type("NullInt64", reflect.TypeOf((*sql.NullInt64)(nil)).Elem()),
		I.Type("NullString", reflect.TypeOf((*sql.NullString)(nil)).Elem()),
		I.Type("NullTime", reflect.TypeOf((*sql.NullTime)(nil)).Elem()),
		I.Type("Out", reflect.TypeOf((*sql.Out)(nil)).Elem()),
		I.Type("RawBytes", reflect.TypeOf((*sql.RawBytes)(nil)).Elem()),
		I.Type("Result", reflect.TypeOf((*sql.Result)(nil)).Elem()),
		I.Type("Row", reflect.TypeOf((*sql.Row)(nil)).Elem()),
		I.Type("Rows", reflect.TypeOf((*sql.Rows)(nil)).Elem()),
		I.Type("Scanner", reflect.TypeOf((*sql.Scanner)(nil)).Elem()),
		I.Type("Stmt", reflect.TypeOf((*sql.Stmt)(nil)).Elem()),
		I.Type("Tx", reflect.TypeOf((*sql.Tx)(nil)).Elem()),
		I.Type("TxOptions", reflect.TypeOf((*sql.TxOptions)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("LevelDefault", qspec.Int, sql.LevelDefault),
		I.Const("LevelLinearizable", qspec.Int, sql.LevelLinearizable),
		I.Const("LevelReadCommitted", qspec.Int, sql.LevelReadCommitted),
		I.Const("LevelReadUncommitted", qspec.Int, sql.LevelReadUncommitted),
		I.Const("LevelRepeatableRead", qspec.Int, sql.LevelRepeatableRead),
		I.Const("LevelSerializable", qspec.Int, sql.LevelSerializable),
		I.Const("LevelSnapshot", qspec.Int, sql.LevelSnapshot),
		I.Const("LevelWriteCommitted", qspec.Int, sql.LevelWriteCommitted),
	)
}
